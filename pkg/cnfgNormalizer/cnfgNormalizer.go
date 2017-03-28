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
	"istio.io/mixer/pkg/config"
	aconfig "istio.io/mixer/pkg/aspect/config"
	"fmt"
	configpb "istio.io/mixer/pkg/config/proto"
	"github.com/robertkrimen/otto"
	"istio.io/mixer/pkg/attribute"
)

type NormalizedConfigJS struct {
	JavaScript string
}

func (n NormalizedConfigJS) Evalaute(requestBag *attribute.MutableBag,
	aspectFinder func(cfgs []*configpb.Combined, kind string, adapterName string, val interface{}) (*configpb.Combined, error),
	callBack func(kind string, adapterImplName string, val interface{})) {
	vm := otto.New()
	vm.Set("CallBackFromUserScript_go", callBack)
	vm.Run(n.JavaScript)
	checkFn, _ := vm.Get("report")
	_, errFromJS := checkFn.Call(otto.NullValue(), requestBag)
	if errFromJS != nil {
		fmt.Println(errFromJS)
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
	var injectionMethods string
	for _, aspectRule := range vd.GetValidatedSC().GetRules() {
		var tmpReportMethodStr string
		for _, aspect := range aspectRule.GetAspects() {
			if aspect.Kind == "metrics" {
				invocation, wrapperMtdToIngest := GetJSInvocationForMetricAspect(aspect.Params.(*aconfig.MetricsParams), aspect.Adapter)
				tmpReportMethodStr = tmpReportMethodStr + invocation
				// HACK. skipping duplication incertion of mtd.
				if len(injectionMethods) == 0 {
					injectionMethods = injectionMethods + wrapperMtdToIngest
				}
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
	var userScript = userJSAllCode + "\n" +
		injectionMethods + "\n" +
		getAllDeclarations() + "\n" +
		getCallbackDeclaration()

	fmt.Println(userScript)
	err := ioutil.WriteFile("/tmp/generatedJSFromConfig.js", []byte(userScript), 0644)
	if err != nil {
		panic(err)
	}

	return NormalizedConfigJS{JavaScript: userScript}
}

func getCallbackDeclaration() string {
	return `
	var CallBackFromUserScript_go;
	`
}

func getAllDeclarations() string {
	return getMetricAspectAllDeclarations()
}