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
	"fmt"

	"github.com/robertkrimen/otto"

	"istio.io/mixer/pkg/attribute"
	"istio.io/mixer/pkg/config"
	"istio.io/mixer/pkg/cnfgNormalizer/typeScriptGenerator"
	pb "istio.io/mixer/pkg/config/proto"

)

type NormalizedJavascriptConfig struct {
	// JavaScript string
	VM        *otto.Otto
	reportMtd otto.Value
}

type NormalizedJavascriptConfigNormalizer struct {
	normalizedConfig config.NormalizedConfig
}

// invoked at runtime
func (n NormalizedJavascriptConfig) Evalaute(requestBag *attribute.MutableBag,
	callBack func(kind string, val interface{})) [][]interface{} {
	resultValue, errFromJS := n.reportMtd.Call(otto.NullValue(), typeScriptGenerator.ConstructAttributesForJS(requestBag))
	if errFromJS != nil {
		fmt.Println("ERROR FROM JS", errFromJS)
	}

	evaluatedData, _ := resultValue.Export()
	v := evaluatedData.(map[string]interface{})["result"]
	return v.([][]interface{})
}


func (n NormalizedJavascriptConfigNormalizer) Normalize(sc *pb.ServiceConfig, fileLocation string) config.NormalizedConfig {

	typeDefTSCode := getPredefinedTypesForDescriptors(sc)

	attributeTypeDeclaration := getAttributesDeclaration()

	fileForTypesFromAspectDescriptors := "TypesFromAspectDescriptors.ts"
	fileForWellKnownAttribs := "WellKnownAttribs.ts"
	userTSAllCode := getUserTSCodeFile(sc, fileForTypesFromAspectDescriptors, fileForWellKnownAttribs)

	generatedJS := getJS(userTSAllCode, typeDefTSCode, attributeTypeDeclaration, fileForTypesFromAspectDescriptors, fileForWellKnownAttribs, fileLocation)

	n.normalizedConfig = createNormalizedConfig(generatedJS)
	return n.normalizedConfig
}

func (n NormalizedJavascriptConfigNormalizer) ReloadNormalizedConfigFile(fileLocation string) config.NormalizedConfig {
	generatedJS := GenerateJsFromTypeScript(fileLocation)
	n.normalizedConfig = createNormalizedConfig(generatedJS)
	return n.normalizedConfig
}

func createNormalizedConfig(generatedJS string) config.NormalizedConfig {
	var vm *otto.Otto
	vm = otto.New()
	vm.Run(generatedJS)
	reportMtd, _ := vm.Get("report")
	return NormalizedJavascriptConfig{VM: vm, reportMtd: reportMtd}
}
