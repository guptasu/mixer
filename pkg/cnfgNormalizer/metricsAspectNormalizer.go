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
	aconfig "istio.io/mixer/pkg/aspect/config"
	"fmt"
)

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

func getMetricAspectAllDeclarations() string {
	return `
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
