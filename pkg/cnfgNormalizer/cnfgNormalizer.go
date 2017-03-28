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
	"os/exec"
	"io/ioutil"
	"github.com/golang/glog"
	"istio.io/mixer/pkg/config"
	aconfig "istio.io/mixer/pkg/aspect/config"
	"istio.io/mixer/pkg/expr"
	"fmt"
)


func Normalize(vd *config.Validated) *config.NormalizedConfig {
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
				invocation, wrapperMtdToIngest := GetJSInvocationForMetricAspectTest(aspect.Params.(*aconfig.MetricsParams), aspect.Adapter)
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
	var userScript = userJSAllCode + "\n" + injectionMethods + "\n" + getAllDeclarations()

	fmt.Println(userScript)
	err := ioutil.WriteFile("/tmp/generatedJSFromConfig.js", []byte(userScript), 0644)
	if err != nil {
		panic(err)
	}

	return &config.NormalizedConfig{JavaScript: userScript}
}
///////////////////// THIS SHOULD BELONG TO METRIC ASPECT MANAGER /////////////////

func GetJSInvocationForMetricAspect(metricsParams *aconfig.MetricsParams, adapterName string) (string, string) {
	params := metricsParams
	//var allMetricsStr bytes.Buffer
	var metricsStr bytes.Buffer
	for _, metric := range params.Metrics {
		var labelStr bytes.Buffer
		labelStr.WriteString(fmt.Sprintf("%s: %s", "value", getJSForExpression(metric.Value)))
		labelLen := len(metric.Labels)
		if labelLen != 0 {
			labelStr.WriteString(",\n")
		}
		for key, value := range metric.Labels {
			labelStr.WriteString(fmt.Sprintf(`				    %s: %s`, key, getJSForExpression(value)))
			labelLen--
			if labelLen != 0 {
				labelStr.WriteString(",\n")
			}
		}
		metricsStr.WriteString(fmt.Sprintf(`
			      "%s": {
			        %s
			      },`, metric.DescriptorName, labelStr.String()))

	}
	metricsStrBuilt := metricsStr.String()
	if metricsStrBuilt[len(metricsStrBuilt)-1] == ',' {
		metricsStrBuilt = metricsStrBuilt[0 : len(metricsStrBuilt)-1]
	}
	callStr := fmt.Sprintf(`
				%s({
	  			  %s
				})
	`, GetJSWrapperMethodNameForMetricAspect(adapterName), metricsStrBuilt)
	return callStr, GetJSWrapperMethodsToInjectForMetricAspect(adapterName, "metrics")
}

func GetJSInvocationForMetricAspectTest(metricsParams *aconfig.MetricsParams, adapterName string) (string, string) {
	err := exec.Command("tsc", "--outFile", "/tmp/genJSFromTS.js", "/tmp/tmp.eWerwGNd6f/greeter.ts").Run()
	if err != nil {
		fmt.Println("tst generation failed", err)
	}
	return `
    if (attributes.Get("source.name")[0] == "test") {
        {
            var requestCount1 = new RequestCount();
            requestCount1.value = 1;
            requestCount1.target = attributes.Get("target.name")[1] ? attributes.Get("target.name")[0] : "one";
            requestCount1.method = attributes.Get("api.method")[1] ? attributes.Get("api.method")[0] : "one";
            requestCount1.response_code = attributes.Get("response.http.code")[1] ? attributes.Get("response.http.code")[0] : 200;
            requestCount1.service = attributes.Get("api.name")[1] ? attributes.Get("api.name")[0] : "one";
            requestCount1.source = attributes.Get("source.name")[1] ? attributes.Get("source.name")[0] : "one";
            var result = new MetricAspectResponse();
            result.request_count = requestCount1;
            RecordToprometheus(result);
        }
    }
	`, `

function RecordToprometheus(val) {
    CallBackFromUserScript_go("metrics", "prometheus", val)
}
	`
}

func getAllDeclarations() string {
	return `
	var CallBackFromUserScript_go;
	var RequestCount = (function () {
		function RequestCount() {
		}
		return RequestCount;
	}());
	var RequestLatency = (function () {
	function RequestLatency() {
	}
	return RequestLatency;
	}());
	var MetricAspectResponse = (function () {
	function MetricAspectResponse() {
	}
	return MetricAspectResponse;
	}());
	`
}

func GetJSWrapperMethodNameForMetricAspect(adapterName string) string {
	return "RecordTo" + adapterName
}

func GetJSWrapperMethodsToInjectForMetricAspect(adapterName string, kindName string) string {

	methodName := GetJSWrapperMethodNameForMetricAspect(adapterName)
	var embeddedMethodsInUserScriptFmt = `
                function %s(val) {
                  CallBackFromUserScript_go("%s", "%s", val)
                }
`
	return fmt.Sprintf(embeddedMethodsInUserScriptFmt, methodName, kindName, adapterName)
}

func getJSForExpression(expression string) string{
	ex, err := expr.Parse(expression)
	var out string
	if err != nil {
		glog.Warning("Unable to parse : %s. %v. Setting expression to false", expression, err)
		out = "false"
	} else {
		condition, _ := EvalJSExpession(ex, expr.FuncMap(), "attributes.Get")
		out = condition
	}
	return out
}

// Eval evaluates the expression given an attribute bag and a function map.
func EvalJSExpession(e *expr.Expression, fMap map[string]expr.FuncBase, getPropMtdName string) (string, error) {
	if e.Const != nil {
		return e.Const.StrValue, nil
	}
	if e.Var != nil {
		return fmt.Sprintf(getPropMtdName + "(\"%s\")[0]", e.Var.Name), nil
	}

	fn := fMap[e.Fn.Name]
	if fn == nil {
		return "", fmt.Errorf("unknown function: %s", e.Fn.Name)
	}
	// may panic
	if e.Fn.Name == "EQ" {
		leftStr, _ := EvalJSExpession(e.Fn.Args[0], fMap, getPropMtdName)
		rightStr, _ := EvalJSExpession(e.Fn.Args[1], fMap, getPropMtdName)
		return fmt.Sprintf("%s == %s", leftStr, rightStr), nil
	}
	if e.Fn.Name == "OR" {
		//(age < 18) ? "Too young":"Old enough"
		allArgs := e.Fn.Args
		if len(allArgs) > 0 {
			chkIfExists := fmt.Sprintf(getPropMtdName+"(\"%s\")[1]", allArgs[0].Var.Name)
			leftexp, _ := EvalJSExpession(e.Fn.Args[0], fMap, getPropMtdName)
			rightexp, _ := EvalJSExpession(e.Fn.Args[1], fMap, getPropMtdName)
			return fmt.Sprintf("%s ? %s : %s", chkIfExists, leftexp, rightexp), nil
		}
	}
	return "", nil
}