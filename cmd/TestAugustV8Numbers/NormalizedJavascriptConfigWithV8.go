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

package main


import (
	"istio.io/mixer/pkg/attribute"
	"istio.io/mixer/pkg/cnfgNormalizer/typeScriptGenerator"
	"github.com/augustoroman/v8"
	"fmt"
	"istio.io/mixer/pkg/config"
	pb "istio.io/mixer/pkg/config/proto"
)

type NormalizedJavascriptConfigWithV8 struct {
	v8Context *v8.Context
	reportMtd *v8.Value
}

// invoked at runtime
func (n NormalizedJavascriptConfigWithV8) Evalaute(requestBag *attribute.MutableBag,
	callBack func(kind string, val interface{})) [][]interface {} {
	attribsV8Value, err := n.v8Context.Create(typeScriptGenerator.ConstructAttributesForJS(requestBag))
	if err != nil {
		panic(fmt.Sprintf("ERROR constructing/binding attribs object", err))

	}
	val, errFromJS := n.reportMtd.Call(nil, attribsV8Value)
	if errFromJS != nil {
		panic(fmt.Sprintf("ERROR FROM JS with v8 engine", errFromJS))
	}

	vresult2,_ := val.Get("result")
	//fmt.Printf("*** value=%v, type=%T\n", vresult2, vresult2)
	var returnVal [][]interface {}
	//vresult2.GetIndex()

	for i := 0 ; i < 50; i++ {
		//fmt.Println("111-", i)
		k, err := vresult2.GetIndex(i)
		//fmt.Println("111-", k)
		if err == nil {
			var objToInsert []interface{}
			p,_ := k.GetIndex(0)
			objToInsert = append(objToInsert, p.String())
			m,_ := k.GetIndex(1)
			//t, _ := m.Get("descriptorName");
			//fmt.Println("fffGGG", t)
			objToInsert = append(objToInsert, transformToMap(m))

			returnVal = append(returnVal, objToInsert)
		} else {
			panic(err)
		}
	}

	//fmt.Println("**** RETURN FROM EVAL")
	//fmt.Println(returnVal)
	return returnVal
}

func transformToMap(value *v8.Value) map[string]interface{} {
	final := make(map[string]interface{})
	tmp, _ := value.Get("descriptorName")
	final["descriptorName"] = tmp.String()
	vValue,_ := value.Get("value")

	valFinal := make(map[string]interface{})
	tmp, e := vValue.Get("value")
	if e == nil {
		valFinal["value"] = tmp
	}
	tmp, e = vValue.Get("source")
	if e == nil {
		valFinal["source"] = tmp
	}
	tmp, e = vValue.Get("target")
	if e == nil {
		valFinal["target"] = tmp
	}
	tmp, e = vValue.Get("service")
	if e == nil {
		valFinal["service"] = tmp
	}
	tmp, e = vValue.Get("method")
	if e == nil {
		valFinal["method"] = tmp
	}
	tmp, e = vValue.Get("response_code")
	if e == nil {
		valFinal["response_code"] = tmp
	}
	final["value"] = valFinal
	return final
}

type NormalizedJavascriptConfigNormalizerWithAugustV8 struct {
	normalizedConfig config.NormalizedConfig
}

func (n NormalizedJavascriptConfigNormalizerWithAugustV8) Normalize(sc *pb.ServiceConfig, fileLocation string) config.NormalizedConfig {

	typeDefTSCode := getPredefinedTypesForDescriptors(sc)

	attributeTypeDeclaration := getAttributesDeclaration()

	fileForTypesFromAspectDescriptors := "TypesFromAspectDescriptors.ts"
	fileForWellKnownAttribs := "WellKnownAttribs.ts"
	userTSAllCode := getUserTSCodeFile(sc, fileForTypesFromAspectDescriptors, fileForWellKnownAttribs)

	generatedJS := getJS(userTSAllCode, typeDefTSCode, attributeTypeDeclaration, fileForTypesFromAspectDescriptors, fileForWellKnownAttribs, fileLocation)

	n.normalizedConfig = createNormalizedConfigWithV8(generatedJS)
	return n.normalizedConfig
}

func (n NormalizedJavascriptConfigNormalizerWithAugustV8) ReloadNormalizedConfigFile(fileLocation string) config.NormalizedConfig {
	generatedJS := GenerateJsFromTypeScript(fileLocation)
	n.normalizedConfig = createNormalizedConfigWithV8(generatedJS)
	return n.normalizedConfig
}

func createNormalizedConfigWithV8(generatedJS string) config.NormalizedConfig {
	ctx := v8.NewIsolate().NewContext()

	_, err := ctx.Eval(generatedJS, "")
	if err != nil {
		fmt.Println("ERROR parsing JS", err)
	}
	reportMtd, err := ctx.Global().Get("report")
	if err != nil {
		fmt.Println("ERROR finding report method", err)
	}
	return NormalizedJavascriptConfigWithV8{v8Context: ctx, reportMtd: reportMtd}
}
