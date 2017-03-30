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
	"github.com/robertkrimen/otto"
	"istio.io/mixer/pkg/attribute"
)

var (
	callbackMtdName = "CallBackFromUserScript_go"
	callbackMtdDeclaration = "var " + callbackMtdName + " = function(name: string, val: any){};"
)

type NormalizedConfigJS struct {
	JavaScript string
}

func (n NormalizedConfigJS) Evalaute(requestBag *attribute.MutableBag,
	callBack func(kind string, val interface{})) {
	vm := otto.New()
	vm.Run(n.JavaScript)
	vm.Set(callbackMtdName, callBack)
	checkFn, _ := vm.Get("report")
	_, errFromJS := checkFn.Call(otto.NullValue(), requestBag)
	if errFromJS != nil {
		fmt.Println("ERROR FROM JS", errFromJS)
	}
}

func Normalize(vd *config.Validated) config.NormalizedConfig {
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
				userCodeForMetricAspect := GenerateUserCodeForMetrics(aspect.Params.(*aconfig.MetricsParams), aspect.Adapter)
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
		function report(attributes) {
		    %s
		}
		function check(attributes) {
		    // TODO
		}
		function quota(attributes) {
		    // TODO
		}
		`

	userJSAllCode := fmt.Sprintf(allJSMethodFormat, reportMethodStr)
	var userScript = ""+
		"\n//-----------------CallBack Method Declaration-----------------\n" +
		"//This method gets injected at runtime. Need this declaration to make TypeScript happy\n" +
		callbackMtdDeclaration + "\n" +
		"\n//-----------------All Types Declaration-----------------\n" +
		getAllDeclarations() + "\n" +
		"\n//-----------------User Code-----------------\n" +
		userJSAllCode

	fmt.Println(userScript)

	tempTSFile := "/tmp/TSConversion/userTS.ts"
	tempGeneratedJSFile := "/tmp/TSConversion/generatedJSFromUserTS.js"

	ioutil.WriteFile(tempTSFile, []byte(userScript), 0644)
	err := exec.Command("clang-format", "-i", tempTSFile).Run()
	err = exec.Command("tsc", "--outFile", tempGeneratedJSFile, tempTSFile).Run()
	if err != nil {
		fmt.Println("tst generation failed", err)
	}
	generatedJS, err := ioutil.ReadFile(tempGeneratedJSFile)
	if err != nil {
		fmt.Println("cannot read generated JS file", err)
	}
	return NormalizedConfigJS{JavaScript: string(generatedJS)}
}

func getAllDeclarations() string {
	return GetMetricAspectAllDeclarations(callbackMtdName)
}