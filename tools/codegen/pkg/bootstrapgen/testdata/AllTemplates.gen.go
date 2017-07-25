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

// THIS FILE IS AUTOMATICALLY GENERATED.

package template

import (
	"github.com/gogo/protobuf/proto"

	"istio.io/mixer/pkg/template"
)

var (
	SupportedTmplInfo = map[string]template.Info{

		foo_bar_mylistchecker.TemplateName: {
			InferTypeFn:     inferTypeForList,
			CnstrDefConfig:  &foo_bar_mylistchecker.ConstructorParam{},
			ConfigureTypeFn: configureTypeForList,
		},

		istio_mixer_adapter_metric.TemplateName: {
			InferTypeFn:     inferTypeForMetric,
			CnstrDefConfig:  &istio_mixer_adapter_metric.ConstructorParam{},
			ConfigureTypeFn: configureTypeForMetric,
		},

		istio_mixer_adapter_quota.TemplateName: {
			InferTypeFn:     inferTypeForQuota,
			CnstrDefConfig:  &istio_mixer_adapter_quota.ConstructorParam{},
			ConfigureTypeFn: configureTypeForQuota,
		},

		istio_mixer_adapter_log.TemplateName: {
			InferTypeFn:     inferTypeForLog,
			CnstrDefConfig:  &istio_mixer_adapter_log.ConstructorParam{},
			ConfigureTypeFn: configureTypeForLog,
		},
	}
)

/////////////////////// Start generated code for template List ///////////////////////
func supportsListBuilder(hndlrBuilder adptConfig.HandlerBuilder) bool {
	_, ok := hndlrBuilder.(foo_bar_mylistchecker.ListProcessorBuilder)
	return ok
}
func inferTypeForList(cp proto.Message, tEvalFn template.TypeEvalFn) (proto.Message, error) {
	var err error

	cpb := cp.(*foo_bar_mylistchecker.ConstructorParam)

	infrdType := &foo_bar_mylistchecker.Type{}

	return infrdType, nil
}

/////////////////////// End generated code for template List ///////////////////////

/////////////////////// Start generated code for template Metric ///////////////////////
func supportsMetricBuilder(hndlrBuilder adptConfig.HandlerBuilder) bool {
	_, ok := hndlrBuilder.(istio_mixer_adapter_metric.MetricProcessorBuilder)
	return ok
}
func inferTypeForMetric(cp proto.Message, tEvalFn template.TypeEvalFn) (proto.Message, error) {
	var err error

	cpb := cp.(*istio_mixer_adapter_metric.ConstructorParam)

	infrdType := &istio_mixer_adapter_metric.Type{}

	return infrdType, nil
}

/////////////////////// End generated code for template Metric ///////////////////////

/////////////////////// Start generated code for template Quota ///////////////////////
func supportsQuotaBuilder(hndlrBuilder adptConfig.HandlerBuilder) bool {
	_, ok := hndlrBuilder.(istio_mixer_adapter_quota.QuotaProcessorBuilder)
	return ok
}
func inferTypeForQuota(cp proto.Message, tEvalFn template.TypeEvalFn) (proto.Message, error) {
	var err error

	cpb := cp.(*istio_mixer_adapter_quota.ConstructorParam)

	infrdType := &istio_mixer_adapter_quota.Type{}

	return infrdType, nil
}

/////////////////////// End generated code for template Quota ///////////////////////

/////////////////////// Start generated code for template Log ///////////////////////
func supportsLogBuilder(hndlrBuilder adptConfig.HandlerBuilder) bool {
	_, ok := hndlrBuilder.(istio_mixer_adapter_log.LogProcessorBuilder)
	return ok
}
func inferTypeForLog(cp proto.Message, tEvalFn template.TypeEvalFn) (proto.Message, error) {
	var err error

	cpb := cp.(*istio_mixer_adapter_log.ConstructorParam)

	infrdType := &istio_mixer_adapter_log.Type{}

	return infrdType, nil
}

/////////////////////// End generated code for template Log ///////////////////////
