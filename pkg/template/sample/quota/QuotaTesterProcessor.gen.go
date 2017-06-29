// Copyright 2017 Istio Authors.
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

package istio_mixer_adapter_sample_quota

import (
	"istio.io/mixer/pkg/adapter/config"
	"istio.io/mixer/pkg/adapter"
)

const TemplateName = "istio.mixer.adapter.sample.quota.Sample"

// Instance represent the runtime structure that will be passed to the CheckQuotaSample method in the handlers.
type Instance struct {
	Name       string
	Dimensions map[string]interface{}
}

// SampleProcessorBuilder represent the Go interface that HandlerBuilder must implement if it wants to be server
// template named `istio.mixer.adapter.sample.quota.Sample` and wants to be configured with configured Types
// for that template.
type SampleProcessorBuilder interface {
	config.HandlerBuilder
	ConfigureSample(map[string]*Type /*Constructor:instance_name to Type mapping. Note type name will not be passed at all*/) error
}

// SampleProcessor represent the Go interface that handlers must implement if it wants to process the template
// named `istio.mixer.adapter.sample.quota.Sample`
type SampleProcessor interface {
	config.Handler
	CheckQuotaSample([]*Instance, adapter.QuotaArgs) (adapter.QuotaResult, error)
}
