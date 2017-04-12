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
	"bytes"
	"fmt"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	aconfig "istio.io/mixer/pkg/aspect/config"
	"istio.io/mixer/pkg/attribute"
	"istio.io/mixer/pkg/config"
	pb "istio.io/mixer/pkg/config/proto"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	callbackMtdName        = "CallBackFromUserScript_go"
	callbackMtdDeclaration = "var " + callbackMtdName + " = function(aspectName: string, val: any){};"
)

type NormalizedJavascriptConfig struct {
	// JavaScript string
	VM *otto.Otto
}

type NormalizedJavascriptConfigNormalizer struct {
	normalizedJavascriptConfig NormalizedJavascriptConfig
}

func (n NormalizedJavascriptConfigNormalizer) Normalize(sc *pb.ServiceConfig, fileLocation string) config.NormalizedConfig {
	typeDefTSCode := getPredefinedTypesForDescriptors(sc)

	attributeTypeDeclaration := getAttributesDeclaration()

	fileForTypesFromAspectDescriptors := "TypesFromAspectDescriptors.ts"
	fileForWellKnownAttribs := "WellKnownAttribs.ts"
	userTSAllCode := getUserTSCodeFile(sc, fileForTypesFromAspectDescriptors, fileForWellKnownAttribs)

	generatedJS := getJS(userTSAllCode, typeDefTSCode, attributeTypeDeclaration, fileForTypesFromAspectDescriptors, fileForWellKnownAttribs, fileLocation)

	var vm *otto.Otto
	vm = otto.New()
	vm.Run(generatedJS)

	n.normalizedJavascriptConfig = NormalizedJavascriptConfig{VM: vm}
	return n.normalizedJavascriptConfig
}

func (n NormalizedJavascriptConfigNormalizer) ReloadNormalizedConfigFile(fileLocation string) config.NormalizedConfig {
	generatedJS := GenerateJsFromTypeScript(fileLocation)
	var vm *otto.Otto
	vm = otto.New()
	vm.Run(generatedJS)

	n.normalizedJavascriptConfig = NormalizedJavascriptConfig{VM: vm}
	return n.normalizedJavascriptConfig
}

// invoked at runtime
func (n NormalizedJavascriptConfig) Evalaute(requestBag *attribute.MutableBag,
	callBack func(kind string, val interface{})) {
	n.VM.Set(callbackMtdName, callBack)
	checkFn, _ := n.VM.Get("report")
	_, errFromJS := checkFn.Call(otto.NullValue(), constructAttributesForJS(requestBag))
	if errFromJS != nil {
		fmt.Println("ERROR FROM JS", errFromJS)
	}
}

func getUserTSCodeFile(sc *pb.ServiceConfig, imports ...string) string {
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

	for _, aspectRule := range sc.GetRules() {
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
	for _, fileToImport := range imports {
		importStringBuffer.WriteString(fmt.Sprintf("/// <reference path=\"%s\"/>\n\n", fileToImport))

	}
	userTSAllCode = importStringBuffer.String() + userTSAllCode

	return userTSAllCode

}

func getPredefinedTypesForDescriptors(sc *pb.ServiceConfig) string {
	return "" +
		"\n//-----------------CallBack Method Declaration-----------------\n" +
		"//This method gets injected at runtime. Need this declaration to make TypeScript happy\n" +
		callbackMtdDeclaration + "\n" +
		"\n//-----------------All Types Declaration-----------------\n" +
		getAllDeclarations(sc) + "\n"
}

func getJS(userTSAllCode string, typeDefTSCode string, attributeTypeDeclaration string, fileNameForTypesFromAspectDescriptors string, fileNameForWellKnownAttribs string, srvcConfigFileLocation string) string {
	generatedDirName := filepath.Join(filepath.Dir(srvcConfigFileLocation), getFileNameWithoutExt(srvcConfigFileLocation)+"_generated")
	tempTypeDefsTSFile := filepath.Join(generatedDirName, fileNameForTypesFromAspectDescriptors)
	tempAttribsDefsTSFile := filepath.Join(generatedDirName, fileNameForWellKnownAttribs)
	tempUserTSFile := filepath.Join(generatedDirName, getFileNameWithDifferentExt(srvcConfigFileLocation, ".ts"))
	err := os.MkdirAll(generatedDirName, os.ModePerm)
	if err != nil {
		fmt.Println("cannot create directory", generatedDirName)
	}
	ioutil.WriteFile(tempUserTSFile, []byte(userTSAllCode), 0644)
	_ = exec.Command("clang-format", "-i", tempUserTSFile).Run()
	ioutil.WriteFile(tempTypeDefsTSFile, []byte(typeDefTSCode), 0644)
	_ = exec.Command("clang-format", "-i", tempTypeDefsTSFile).Run()
	ioutil.WriteFile(tempAttribsDefsTSFile, []byte(attributeTypeDeclaration), 0644)
	_ = exec.Command("clang-format", "-i", tempAttribsDefsTSFile).Run()

	return GenerateJsFromTypeScript(tempUserTSFile)
}

func getFileNameWithDifferentExt(filePath string, ext string) string {
	tmp := filepath.Base(filePath)

	return tmp[0:len(tmp)-len(filepath.Ext(tmp))] + ext
}

func getFileNameWithoutExt(filePath string) string {
	tmp := filepath.Base(filePath)
	return tmp[0 : len(tmp)-len(filepath.Ext(tmp))]
}

func getOutFileNameWithDifferentExt(filePath string, ext string) string {
	tmp := filepath.Base(filePath)

	return filepath.Join(filepath.Dir(filePath), tmp[0:len(tmp)-len(filepath.Ext(tmp))]+ext)
}

func GenerateJsFromTypeScript(userTSFile string) string {
	tempGeneratedJSOutFile := getOutFileNameWithDifferentExt(userTSFile, ".js")
	//err := exec.Command("tsc", "--lib", "es7", "--outFile", tempGeneratedJSOutFile, userTSFile).Run()
	err := exec.Command("tsc", "--outFile", tempGeneratedJSOutFile, userTSFile).Run()
	if err != nil {
		fmt.Println("tst generation failed", err)
	}
	generatedJS, err := ioutil.ReadFile(tempGeneratedJSOutFile)
	if err != nil {
		fmt.Println("cannot read generated JS file", err)
	}
	return string(generatedJS)
}

func getAllDeclarations(sc *pb.ServiceConfig) string {
	allUserDeclaredMetricsAspectNames := make([]string, 0)
	allUserDeclaredMetricsAspects := make([]*pb.Aspect, 0)
	for _, aspectRule := range sc.GetRules() {
		for _, aspect := range aspectRule.GetAspects() {
			if aspect.Kind == "metrics" {
				allUserDeclaredMetricsAspects = append(allUserDeclaredMetricsAspects, aspect)
				allUserDeclaredMetricsAspectNames = append(allUserDeclaredMetricsAspectNames, aspect.Name)
			}
			// TODO... need to go nested inside the rules within the aspects
		}
	}
	return GetMetricAspectAllDeclarations(callbackMtdName, allUserDeclaredMetricsAspectNames, allUserDeclaredMetricsAspects)
}

func getAttributesDeclaration() string {
	return GetTypeFromAttributes()
}
