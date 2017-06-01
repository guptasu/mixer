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

// !!!!!!!!!!!!!!!!!!!!! WARNING !!!!!!!!!!!!!!!!!!!!!!!!
// THIS IS AUTO GENERATED FILE - SIMULATED - HAND WRITTEN

package gRPCAdapter

import (
	"fmt"
	generated_adapter_registrar "istio.io/mixer/pkg/adapter/generated"
	foo_bar_mymetric "istio.io/mixer/pkg/templates/mymetric/generated/config"
	"istio.io/mixer/pkg/templates/mymetric/generated"
)

/////////// ALL THE BELOW CODE IS GENERATED FROM TEMPLATES //////////////////
func (builder) ConfigureMyMetric(typeParams map[string]foo_bar_mymetric.Type) {
	fmt.Println("ConfigureMyMetric called with", typeParams)
}


func (builder) ProcessMyMetric(instances []mymetric.Instance) {
	fmt.Println("ProcessMyMetric called with", instances)
}

// Register registers the no-op adapter as every aspect.
func Register(r generated_adapter_registrar.Registrar2) {
	r.RegisterMyMetricProcessor(builder{})
}
