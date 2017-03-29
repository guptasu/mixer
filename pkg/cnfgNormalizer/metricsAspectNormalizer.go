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
	dpb "istio.io/api/mixer/v1/config/descriptor"
	aconfig "istio.io/mixer/pkg/aspect/config"
	"strings"
)

///////////////////// THIS SHOULD BELONG TO METRIC ASPECT MANAGER /////////////////

var (
	valueTypeToJSType = map[dpb.ValueType]string{
		dpb.INT64:  "number",
		dpb.STRING: "string",
	}
	debug       = true
	desc = []*dpb.MetricDescriptor{
		{
			Name:        "request_count",
			Kind:        dpb.COUNTER,
			Value:       dpb.INT64,
			Description: "request count by source, target, service, and code",
			Labels: []*dpb.LabelDescriptor{
				{Name: "source", ValueType: dpb.STRING},
				{Name: "target", ValueType: dpb.STRING},
				{Name: "service", ValueType: dpb.STRING},
				{Name: "method", ValueType: dpb.STRING},
				{Name: "response_code", ValueType: dpb.INT64},
			},
		},
		{
			Name:        "request_latency",
			Kind:        dpb.COUNTER,
			Value:       dpb.INT64,
			Description: "request latency by source, target, and service",
			Labels: []*dpb.LabelDescriptor{
				{Name: "source", ValueType: dpb.STRING},
				{Name: "target", ValueType: dpb.STRING},
				{Name: "service", ValueType: dpb.STRING},
				{Name: "method", ValueType: dpb.STRING},
				{Name: "response_code", ValueType: dpb.INT64},
			},
		},
	}
)

func GenerateUserCodeForMetrics(metricsParams *aconfig.MetricsParams, adapterName string) string {
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
			labelStr.WriteString(fmt.Sprintf(`%s: %s`, key, getJSForExpression(value)))
			labelLen--
			if labelLen != 0 {
				labelStr.WriteString(",\n")
			}
		}
		methodToInvoke := getMethodNameForDescriptor(metric.DescriptorName)
		if debug {
			metricsStr.WriteString(fmt.Sprintf("\nconsole.log(\"invoking %s\")\n", methodToInvoke))
		}
		callStr := fmt.Sprintf(`
				%s({
	  			  %s
				})`, methodToInvoke, labelStr.String())
		metricsStr.WriteString(callStr + "\n")
		if debug {
			metricsStr.WriteString(fmt.Sprintf("\nconsole.log(\"Done %s\")\n", methodToInvoke))
		}

	}
	return metricsStr.String()
}

func getMethodNameForDescriptor(descriptorName string) string {
	//return "RecordTo" + adapterName'
	return "Record" + snake2UpperCamelCase(descriptorName)
}

func GetMetricAspectAllDeclarations(callbackMtdName string) string {
	return getTypesAndMethodsForDescriptors(callbackMtdName)

}

func getTypesAndMethodsForDescriptors(callbackMtdName string) string {
	var metricsStr bytes.Buffer
	for _, metricDescriptor := range desc {
		metricsStr.WriteString(createIndividualClass(metricDescriptor))
	}
	for _, metricDescriptor := range desc {
		metricsStr.WriteString(createIndividualMethod(metricDescriptor, callbackMtdName))
	}

	return metricsStr.String()
}

func createIndividualClass(metricDescriptor *dpb.MetricDescriptor) string {
	var metricsStr bytes.Buffer
	var fieldsStr bytes.Buffer
	fieldsStr.WriteString(fmt.Sprintf("%s: %s;", "value", getJSType(metricDescriptor.Value)))
	for _, label := range metricDescriptor.Labels {
		fieldsStr.WriteString(fmt.Sprintf("%s: %s;", label.Name, getJSType(label.ValueType)))
	}
	metricsStr.WriteString(fmt.Sprintf("class %s {%s}\n", snake2UpperCamelCase(metricDescriptor.Name), fieldsStr.String()))
	return metricsStr.String()
}

func getJSType(valueType dpb.ValueType) string {
	return valueTypeToJSType[valueType]
}

func createIndividualMethod(metricDescriptor *dpb.MetricDescriptor, callbackMtdname string) string {
	var metricsDescriptorFunction bytes.Buffer
	var fieldsStr bytes.Buffer
	fieldsStr.WriteString(fmt.Sprintf("%s: %s;", "value", getJSType(metricDescriptor.Value)))
	for _, label := range metricDescriptor.Labels {
		fieldsStr.WriteString(fmt.Sprintf("%s: %s;", label.Name, getJSType(label.ValueType)))
	}
	typeName := snake2UpperCamelCase(metricDescriptor.Name)
	metricsDescriptorFunction.WriteString(fmt.Sprintf(`
	function Record%s(val: %s) {
	  console.log("Record%s invoked");
	  %s("metrics", {descriptorName: "%s", value: val})
	  console.log("Record%s finished");
	}
	`, typeName, typeName, typeName, callbackMtdname, metricDescriptor.Name, typeName))
	return metricsDescriptorFunction.String()
}

func snake2UpperCamelCase(s string) string {
	subStrs := strings.Split(s, "_")
	for i, subStr := range subStrs {
		subStrs[i] = strings.Title(subStr)
	}
	return strings.Join(subStrs, "")
}
