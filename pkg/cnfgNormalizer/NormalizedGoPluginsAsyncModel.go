// Copyright 2017 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cnfgNormalizer

import (
	//"fmt"
	"plugin"
	"istio.io/mixer/pkg/attribute"
	"istio.io/mixer/pkg/config"
)

type NormalizedGoPluginsAsyncModel struct {
	report        plugin.Symbol
	plugin        *plugin.Plugin
	goPackagePath string
}


// invoked at runtime
func (n NormalizedGoPluginsAsyncModel) Evalaute(requestBag *attribute.MutableBag,
	callBack func(kind string, val interface{})) [][]interface{} {
	report, _ := n.plugin.Lookup("Report")
	attr := constructAttributesForGoPlugin(requestBag)
	aspectNameToMethodMapping := report.(func(map[string]interface{}) map[string]string)(attr)
	//fmt.Println(aspectNameToMethodMapping)

	l := len(aspectNameToMethodMapping)
	//fmt.Printf("## length of aspects found = %d\n", len(aspectNameToMethodMapping))
	result := make([][]interface{}, l)
	resultChan := make(chan []interface{}, l)
	for aspectID, mtdForEval := range aspectNameToMethodMapping {
		go func() {
			mtd, _ := n.plugin.Lookup(mtdForEval)
			evalData := mtd.(func(map[string]interface{}) interface{})(attr)
			//fmt.Println(evalData)
			resultChan <- []interface{}{aspectID, map[string]interface{}{"descriptorName": "request_count", "value": evalData}}
			//result = append(result, innerValue)
			//done <- true
		}()
	}
	for i := 0; i < l; i++ {
		select {
		case res := <-resultChan:
			result[i] = res
		}
	}

	//fmt.Println("LENGTH OF THE EVALUATED DATA", len(result))
	return result
}

func createNormalizedGoPluginConfigAsyncModel(goPackagePath string) config.NormalizedConfig {
	p, _ := plugin.Open(goPackagePath)
	return NormalizedGoPluginsAsyncModel{goPackagePath: goPackagePath, plugin: p}
}
