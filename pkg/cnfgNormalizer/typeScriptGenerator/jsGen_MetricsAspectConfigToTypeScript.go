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

package typeScriptGenerator

import (
	"bytes"
	"fmt"
	"strings"

	dpb "istio.io/api/mixer/v1/config/descriptor"
	aconfig "istio.io/mixer/pkg/aspect/config"
	pb "istio.io/mixer/pkg/config/proto"
)

///////////////////// THIS SHOULD BELONG TO METRIC ASPECT MANAGER /////////////////

var (
	desc = []*dpb.MetricDescriptor{
		{
			Name:        "request_count",
			Kind:        dpb.COUNTER,
			Value:       dpb.INT64,
			Description: "request count by source, target, service, and code",
			Labels: map[string]dpb.ValueType{
				"source":        dpb.STRING,
				"target":        dpb.STRING,
				"service":       dpb.STRING,
				"method":        dpb.STRING,
				"response_code": dpb.INT64,
			},
		},
		{
			Name:        "request_latency",
			Kind:        dpb.COUNTER,
			Value:       dpb.INT64,
			Description: "request latency by source, target, and service",
			Labels: map[string]dpb.ValueType{
				"source":        dpb.STRING,
				"target":        dpb.STRING,
				"service":       dpb.STRING,
				"method":        dpb.STRING,
				"response_code": dpb.INT64,
			},
		},
	}
)

func GenerateUserCodeForMetrics(metricsParams *aconfig.MetricsParams, aspectName string) string {
	params := metricsParams

	//var allMetricsStr bytes.Buffer
	var metricsStr bytes.Buffer
	for _, metric := range params.Metrics {
		methodToInvoke := getMethodNameFromDescriptorAndAspectName(metric.DescriptorName, aspectName)
		callStr := fmt.Sprintf(`
				result.%s(
	  			  %s(%s)
				)`, methodToInvoke, getMethodNameForConstructionOfDescriptor(metric.DescriptorName, aspectName), "attributes")
		metricsStr.WriteString(callStr + "\n")

	}
	return metricsStr.String()
}

func createObject(metric *aconfig.MetricsParams_Metric) bytes.Buffer {
	var labelStr bytes.Buffer
	labelStr.WriteString(fmt.Sprintf("%s: %s", "value", GetJSForExpression(metric.Value)))
	labelLen := len(metric.Labels)
	if labelLen != 0 {
		labelStr.WriteString(",\n")
	}
	for key, value := range metric.Labels {
		labelStr.WriteString(fmt.Sprintf(`%s: %s`, key, GetJSForExpression(value)))
		labelLen--
		if labelLen != 0 {
			labelStr.WriteString(",\n")
		}
	}
	return labelStr
}

func GetMetricAspectAllDeclarations(callbackMtdName string, declaredMetricAspectNames []string, declaredMetricAspects []*pb.Aspect) string {
	var metricsStr bytes.Buffer
	for _, metricDescriptor := range desc {
		metricsStr.WriteString(createClassForDescriptor(metricDescriptor))
	}
	metricsStr.WriteString(createReportResultClass(declaredMetricAspectNames))
	for _, metricAspect := range declaredMetricAspects {
		params := metricAspect.Params.(*aconfig.MetricsParams)
		for _, metric := range params.Metrics {
			metricsStr.WriteString(createDescrptorConstructionMethod(metric, metricAspect.Name) + "\n")
		}
	}
	return metricsStr.String()
}

func createReportResultClass(declaredMetricAspectNames []string) string {
	template := `
          class ReportResult {
              private result : [string,any][]

              constructor(){
                  this.result = [];
              }

              %s

              Build() {
                  return this.result;
              }
          }
	`
	var metricsStr bytes.Buffer
	for _, metricDescriptor := range desc {
		metricsStr.WriteString(createMethodForDescriptor(metricDescriptor, declaredMetricAspectNames))
	}

	return fmt.Sprintf(template, metricsStr.String())
}
func getMethodNameFromDescriptorAndAspectName(descriptorNameSnakeCase string, aspectNameSnakeCase string) string {
	return fmt.Sprintf("Insert%sFor%s", snake2UpperCamelCase(descriptorNameSnakeCase), snake2UpperCamelCase(aspectNameSnakeCase))
}

func getMethodNameForConstructionOfDescriptor(descriptorNameSnakeCase string, aspectNameSnakeCase string) string {
	return fmt.Sprintf("Construct%sFor%s", snake2UpperCamelCase(descriptorNameSnakeCase), snake2UpperCamelCase(aspectNameSnakeCase))
}

func createDescrptorConstructionMethod(metricDescriptorObject *aconfig.MetricsParams_Metric, aspectName string) string {
	methodNameForConstructor := getMethodNameForConstructionOfDescriptor(metricDescriptorObject.DescriptorName, aspectName)
	var fieldBuffer = createObject(metricDescriptorObject)
	return fmt.Sprintf(`
                                function %s(attributes: Attributes) {
	                          return {%s}
	                        }
				`, methodNameForConstructor, fieldBuffer.String())
}

func createClassForDescriptor(metricDescriptor *dpb.MetricDescriptor) string {
	var metricsStr bytes.Buffer
	var fieldsStr bytes.Buffer
	fieldsStr.WriteString(fmt.Sprintf("%s: %s;", "value", getJSType(metricDescriptor.Value)))
	for name, valType := range metricDescriptor.Labels {
		fieldsStr.WriteString(fmt.Sprintf("%s: %s;", name, getJSType(valType)))
	}
	metricsStr.WriteString(fmt.Sprintf("class %s {%s}\n", snake2UpperCamelCase(metricDescriptor.Name), fieldsStr.String()))
	return metricsStr.String()
}

func createMethodForDescriptor(metricDescriptor *dpb.MetricDescriptor, declaredMetricAspectNames []string) string {
	var metricsDescriptorFunction bytes.Buffer
	descriptorNameUpperCamel := snake2UpperCamelCase(metricDescriptor.Name)
	for _, aspectIDSnakeCaseAssumed := range declaredMetricAspectNames {
		metricsDescriptorFunction.WriteString(fmt.Sprintf(`
	%s(val: %s) {
	  this.result.push(["%s", {descriptorName: "%s", value: val}])
	}
	`, getMethodNameFromDescriptorAndAspectName(metricDescriptor.Name, aspectIDSnakeCaseAssumed),
			descriptorNameUpperCamel, aspectIDSnakeCaseAssumed, metricDescriptor.Name))
	}
	return metricsDescriptorFunction.String()
}

func snake2UpperCamelCase(s string) string {
	subStrs := strings.Split(s, "_")
	for i, subStr := range subStrs {
		subStrs[i] = strings.Title(subStr)
	}
	return strings.Join(subStrs, "")
}

func getJSType(valueType dpb.ValueType) string {
	return valueTypeToJSType[valueType]
}
