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
	//"github.com/augustoroman/v8"
)

var (
	callbackMtdName        = "__interal__callback_fn"
	callbackMtdDeclaration = "var " + callbackMtdName + " = function(aspectName: string, val: any){};"
)

type NormalizedJavascriptConfig struct {
	// JavaScript string
	VM *otto.Otto
}

// invoked at runtime
func (n NormalizedJavascriptConfig) Evalaute(requestBag *attribute.MutableBag,
	callBack func(kind string, val interface{})) [][]interface {} {

	checkFn, _ := n.VM.Get("report")
	resultValue, errFromJS := checkFn.Call(otto.NullValue(), constructAttributesForJS(requestBag))
	if errFromJS != nil {
		fmt.Println("ERROR FROM JS", errFromJS)
	}

	evaluatedData,_ := resultValue.Export()
	v := evaluatedData.(map[string]interface{})["result"]
	return v.([][]interface {})
}

func CreateNormalizedJavascriptConfig(js string) NormalizedJavascriptConfig {

	var vm *otto.Otto
	vm = otto.New()
	vm.Run(js)

	return NormalizedJavascriptConfig{VM: vm}
}
/*
type NormalizedJavascriptConfigWithV8 struct {
	v8Context *v8.Context
	reportMtd *v8.Value
}

// invoked at runtime
func (n NormalizedJavascriptConfigWithV8) Evalaute(requestBag *attribute.MutableBag,
	callBack func(kind string, val interface{})) {

	tmpCallBack := func(in v8.CallbackArgs) (*v8.Value, error) {
		//fmt.Printf("Args: %s", in.Args)
		callBack(in.Arg(0).String(), in.Arg(1))
		return nil, nil
	}
	n.v8Context.Global().Set(callbackMtdName, n.v8Context.Bind(callbackMtdName, tmpCallBack))


	attribsV8Value, err := n.v8Context.Create(constructAttributesForJS(requestBag))
	if err != nil {
		fmt.Println("ERROR constructing/binding attribs object", err)
	}
	_, errFromJS := n.reportMtd.Call(nil, attribsV8Value)
	if errFromJS != nil {
		fmt.Println("ERROR FROM JS with v8 engine", errFromJS)
	}
}

func CreateNormalizedJavascriptConfigWithV8(js string) NormalizedJavascriptConfigWithV8 {

	ctx := v8.NewIsolate().NewContext()

	_, err := ctx.Eval(js, "")
	if err != nil {
		fmt.Println("ERROR parsing JS", err)
		//return nil
	}
	reportMtd,err := ctx.Global().Get("report")
	if err != nil {
		fmt.Println("ERROR finding report method", err)
		//return nil
	}
	return NormalizedJavascriptConfigWithV8{v8Context: ctx, reportMtd:reportMtd}
}
*/
type NormalizedJavascriptConfigNormalizer struct {
	normalizedConfig config.NormalizedConfig
}

func (n NormalizedJavascriptConfigNormalizer) Normalize(sc *pb.ServiceConfig, fileLocation string) config.NormalizedConfig {
	typeDefTSCode := getPredefinedTypesForDescriptors(sc)

	attributeTypeDeclaration := getAttributesDeclaration()

	fileForTypesFromAspectDescriptors := "TypesFromAspectDescriptors.ts"
	fileForWellKnownAttribs := "WellKnownAttribs.ts"
	userTSAllCode := getUserTSCodeFile(sc, fileForTypesFromAspectDescriptors, fileForWellKnownAttribs)

	generatedJS := getJS(userTSAllCode, typeDefTSCode, attributeTypeDeclaration, fileForTypesFromAspectDescriptors, fileForWellKnownAttribs, fileLocation)

	//n.normalizedConfig = CreateNormalizedJavascriptConfigWithV8(generatedJS)
	n.normalizedConfig = CreateNormalizedJavascriptConfig(generatedJS)
	return n.normalizedConfig
}

func (n NormalizedJavascriptConfigNormalizer) ReloadNormalizedConfigFile(fileLocation string) config.NormalizedConfig {
	generatedJS := GenerateJsFromTypeScript(fileLocation)
	//n.normalizedConfig = CreateNormalizedJavascriptConfigWithV8(generatedJS)
	n.normalizedConfig = CreateNormalizedJavascriptConfig(generatedJS)
	return n.normalizedConfig
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
		function report(attributes: Attributes) : ReportResult {
		    var result = new ReportResult();
		    %s
		    return result;
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
	err = exec.Command("clang-format", "-i", tempUserTSFile).Run()
	if err != nil {
		fmt.Println("cannot run clang-format", tempUserTSFile, err)
	}
	ioutil.WriteFile(tempTypeDefsTSFile, []byte(typeDefTSCode), 0644)
	err = exec.Command("clang-format", "-i", tempTypeDefsTSFile).Run()
	if err != nil {
		fmt.Println("cannot run clang-format", tempTypeDefsTSFile, err)
	}
	ioutil.WriteFile(tempAttribsDefsTSFile, []byte(attributeTypeDeclaration), 0644)
	err = exec.Command("clang-format", "-i", tempAttribsDefsTSFile).Run()
	if err != nil {
		fmt.Println("cannot run clang-format", tempAttribsDefsTSFile, err)
	}
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
