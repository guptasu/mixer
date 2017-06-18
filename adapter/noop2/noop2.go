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

package noop2

import (
	"fmt"

	"github.com/gogo/protobuf/types"
	"github.com/golang/protobuf/proto"

	adapter "istio.io/mixer/pkg/adapter"
	metric "istio.io/mixer/pkg/templates/metric"
)

type (
	Adapter struct{}
	builder struct{}
)

func (builder) Name() string { return "newnoop" }
func (builder) Description() string {
	return "An adapter that does nothing, just echos the calls made from mixer"
}
func (builder) Close() error { return nil }

func (builder) ValidateConfig(msg proto.Message) error {
	fmt.Println("ValidateConfig called with input", msg)
	return nil
}

func (builder) Configure(msg proto.Message) error {
	fmt.Println("Configure called with input", msg)
	return nil
}

func (builder) DefaultConfig() proto.Message { return &types.Empty{} }

func (builder) ConfigureMetric(typeParams map[string]*metric.Type) error {
	fmt.Println("ConfigureMetric in noop Adapter called with", typeParams)
	return nil
}

func (builder) ReportMetric(instances []*metric.Instance) error {
	fmt.Println("ReportMetric in noop Adapter called with", instances)
	return nil
}

func Register(r adapter.Registrar2) {
	r.RegisterMetricProcessor(builder{})
}