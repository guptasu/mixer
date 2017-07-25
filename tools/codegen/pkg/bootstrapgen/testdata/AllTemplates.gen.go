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
	adptConfig "istio.io/mixer/pkg/adapter/config"

	"istio.io/mixer/template/sample/report" // XXXXXXXXXXX
)

/////////////////////// Start generated code for template List ///////////////////////
func supportsListBuilder(hndlrBuilder adptConfig.HandlerBuilder) bool {
	_, ok := hndlrBuilder.(foo_bar_mylistchecker.ListProcessorBuilder)
	return ok
}

/////////////////////// End generated code for template List ///////////////////////

/////////////////////// Start generated code for template Metric ///////////////////////
func supportsMetricBuilder(hndlrBuilder adptConfig.HandlerBuilder) bool {
	_, ok := hndlrBuilder.(istio_mixer_adapter_metric.MetricProcessorBuilder)
	return ok
}

/////////////////////// End generated code for template Metric ///////////////////////

/////////////////////// Start generated code for template Quota ///////////////////////
func supportsQuotaBuilder(hndlrBuilder adptConfig.HandlerBuilder) bool {
	_, ok := hndlrBuilder.(istio_mixer_adapter_quota.QuotaProcessorBuilder)
	return ok
}

/////////////////////// End generated code for template Quota ///////////////////////

/////////////////////// Start generated code for template Log ///////////////////////
func supportsLogBuilder(hndlrBuilder adptConfig.HandlerBuilder) bool {
	_, ok := hndlrBuilder.(istio_mixer_adapter_log.LogProcessorBuilder)
	return ok
}

/////////////////////// End generated code for template Log ///////////////////////
