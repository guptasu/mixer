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
	"io/ioutil"
	"os/exec"
	"istio.io/mixer/pkg/config"
	aconfig "istio.io/mixer/pkg/aspect/config"
	"fmt"
	"bytes"
	"github.com/robertkrimen/otto"
	"istio.io/mixer/pkg/attribute"
)

var (
	callbackMtdName = "CallBackFromUserScript_go"
	callbackMtdDeclaration = "var " + callbackMtdName + " = function(aspectName: string, val: any){};"
)

type NormalizedJavascriptConfig struct {
	JavaScript string
}

// invoked at runtime
func (n NormalizedJavascriptConfig) Evalaute(requestBag *attribute.MutableBag,
	callBack func(kind string, val interface{})) {
	vm := otto.New()
	vm.Run(n.JavaScript)
	vm.Set(callbackMtdName, callBack)

	attribConstructor, _ := vm.Get("ConstructAttributes")
	attributesFromJS, errFromJS := attribConstructor.Call(otto.NullValue(), requestBag)
	fmt.Println(attributesFromJS)

	checkFn, _ := vm.Get("report")
	_, errFromJS = checkFn.Call(otto.NullValue(), attributesFromJS)
	if errFromJS != nil {
		fmt.Println("ERROR FROM JS", errFromJS)
	}
}

// invoked at configuration time
func Normalize(vd *config.Validated) config.NormalizedConfig {

	typeDefTSCode := getPredefinedTypesForDescriptors(vd)

	attributeTypeDeclaration := getAttributesDeclaration()

	fileForTypesFromAspectDescriptors := "TypesFromAspectDescriptors.ts"
	fileForWellKnownAttribs := "WellKnownAttribs.ts"
	userTSAllCode := getUserTSCodeFile(vd, fileForTypesFromAspectDescriptors, fileForWellKnownAttribs);

	generatedJS := getJS(userTSAllCode, typeDefTSCode, attributeTypeDeclaration, fileForTypesFromAspectDescriptors, fileForWellKnownAttribs)

	return NormalizedJavascriptConfig{JavaScript: generatedJS}
}

func getUserTSCodeFile(vd *config.Validated, imports ...string) string {
	//vd.serviceConfig
	/*
	create methodStrs for each method (chk, report, quota)
	for each rule r:
	  create tmpMethodStrs for each method (chk, report, quota)
          for each aspect a:
	    generate JS k via Aspect Managers
	    add k to appropriate tmpMethodStrs.
	  for each non empty tmpMethodStrs {
	    insert : if (rule) {
	      tmpMethodStr
	    }
	  }
	  Do recurrs for each nested rules.
	*/
	var reportMethodStr string

	for _, aspectRule := range vd.GetValidatedSC().GetRules() {
		var tmpReportMethodStr string
		for _, aspect := range aspectRule.GetAspects() {
			if aspect.Kind == "metrics" {
				userCodeForMetricAspect := GenerateUserCodeForMetrics(aspect.Params.(*aconfig.MetricsParams), aspect.Name)
				tmpReportMethodStr = tmpReportMethodStr + userCodeForMetricAspect
			}
		}
		if len(tmpReportMethodStr) > 0 {
			//TODO FIX FOR COMPLEXT OPTIONS
			var ifStatementStr string
			if len(aspectRule.Selector) > 0 {
				expressionStr := aspectRule.Selector
				ifStatementStr = getJSForExpression(expressionStr)

			} else {
				ifStatementStr = "true"
			}
			ruleIfBlocked := fmt.Sprintf(`
			if (%s) {
			  %s
			}
			`, ifStatementStr, tmpReportMethodStr)

			reportMethodStr = reportMethodStr + ruleIfBlocked
		}
	}
	allJSMethodFormat := `
		function report(attributes: Attributes) {
		    %s
		}
		function check(attributes) {
		    // TODO
		}
		function quota(attributes) {
		    // TODO
		}
		`

	userTSAllCode := fmt.Sprintf(allJSMethodFormat, reportMethodStr)

	var importStringBuffer bytes.Buffer
	for _,fileToImport := range imports {
		importStringBuffer.WriteString(fmt.Sprintf("/// <reference path=\"%s\"/>\n\n", fileToImport))

	}
	userTSAllCode = importStringBuffer.String() + userTSAllCode

	return userTSAllCode

}

func getPredefinedTypesForDescriptors(vd *config.Validated) string {
	return ""+
		"\n//-----------------CallBack Method Declaration-----------------\n" +
		"//This method gets injected at runtime. Need this declaration to make TypeScript happy\n" +
		callbackMtdDeclaration + "\n" +
		"\n//-----------------All Types Declaration-----------------\n" +
		getAllDeclarations(vd) + "\n"
}

func getJS(userTSAllCode string, typeDefTSCode string, attributeTypeDeclaration string, fileNameForTypesFromAspectDescriptors string, fileNameForWellKnownAttribs string) string {
	tempTypeDefsTSFile := "/tmp/TSConversion/" + fileNameForTypesFromAspectDescriptors
	tempAttribsDefsTSFile := "/tmp/TSConversion/" + fileNameForWellKnownAttribs
	tempUserTSFile := "/tmp/TSConversion/UserWrittenMethods.ts"
	tempGeneratedJSFile := "/tmp/TSConversion/generatedJSFromTypeScriptCompiler.js"

	ioutil.WriteFile(tempUserTSFile, []byte(userTSAllCode), 0644)
	err := exec.Command("clang-format", "-i", tempUserTSFile).Run()
	ioutil.WriteFile(tempTypeDefsTSFile, []byte(typeDefTSCode), 0644)
	err = exec.Command("clang-format", "-i", tempTypeDefsTSFile).Run()
	ioutil.WriteFile(tempAttribsDefsTSFile, []byte(attributeTypeDeclaration), 0644)
	err = exec.Command("clang-format", "-i", tempAttribsDefsTSFile).Run()

	err = exec.Command("tsc", "--lib", "es7", "--outFile",  tempGeneratedJSFile, tempUserTSFile).Run()
	if err != nil {
		fmt.Println("tst generation failed", err)
	}
	generatedJS, err := ioutil.ReadFile(tempGeneratedJSFile)
	if err != nil {
		fmt.Println("cannot read generated JS file", err)
	}
	return string(generatedJS);
}

func getAllDeclarations(vd *config.Validated) string {
	allUserDeclaredMetricsAspectNames := make([]string, 0)
	for _, aspectRule := range vd.GetValidatedSC().GetRules() {
		for _, aspect := range aspectRule.GetAspects() {
			if aspect.Kind == "metrics" {
				allUserDeclaredMetricsAspectNames = append(allUserDeclaredMetricsAspectNames, aspect.Name)
			}
			// TODO... need to go nested inside the rules within the aspects
		}
	}
	return GetMetricAspectAllDeclarations(callbackMtdName, allUserDeclaredMetricsAspectNames)
}

func getAttributesDeclaration() string {
	return GetTypeFromAttributes()
}
