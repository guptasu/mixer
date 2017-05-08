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
	"istio.io/mixer/pkg/cnfgNormalizer/typeScriptGenerator"
	pb "istio.io/mixer/pkg/config/proto"
)

type NormalizedGoPlugins struct {
	report        plugin.Symbol
	plugin        *plugin.Plugin
	goPackagePath string
}

type CnftToGopackageNormalizer struct {
	normalizedConfig config.NormalizedConfig
}


func constructAttributesForGoPlugin(requestBag *attribute.MutableBag) map[string]interface{} {
	attribs := make(map[string]interface{})
	for _, attribName := range requestBag.Names() {
		attribs[typeScriptGenerator.DotCaseToCamelCase(attribName)], _ = requestBag.Get(attribName)
	}
	return attribs
}

// invoked at runtime
func (n NormalizedGoPlugins) Evalaute(requestBag *attribute.MutableBag,
	callBack func(kind string, val interface{})) [][]interface{} {
	report, _ := n.plugin.Lookup("Report")
	result := report.(func(map[string]interface{}) [][]interface{})(constructAttributesForGoPlugin(requestBag))


	return result
}

func (n CnftToGopackageNormalizer) Normalize(sc *pb.ServiceConfig, fileLocation string) config.NormalizedConfig {
	// NOT IMPLEMENTED..
	return nil
}

func (n CnftToGopackageNormalizer) ReloadNormalizedConfigFile(fileLocation string) config.NormalizedConfig {
	goPackagePath := fileLocation
	n.normalizedConfig = createGoPackageNormalizedConfig(goPackagePath)
	return n.normalizedConfig
}

func createGoPackageNormalizedConfig(goPackagePath string) config.NormalizedConfig {
	p, _ := plugin.Open(goPackagePath)
	return NormalizedGoPlugins{goPackagePath: goPackagePath, plugin: p}
}
