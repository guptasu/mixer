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

package gRPCAdapter

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	metric "istio.io/mixer/pkg/templates/mymetric/generated/config"
)

type (
	Adapter struct{}
	builder struct{}
)

func (Adapter) ConfigureMetric(typeParams map[string]metric.Type) {
	fmt.Println("ConfigureMetricCalled with", typeParams)
}
func (Adapter) ProcessMetric(instances []metric.Instance) {
	fmt.Println("ConfigureMetricCalled with", instances)
}

func (builder) Name() string        { return "grpcAdapter" }
func (builder) Description() string { return "an adapter that does nothing" }
func (builder) Close() error        { return nil }

func (builder) ValidateConfig(msg proto.Message) error {
	fmt.Println("Handler level ValidateConfig called")
	return nil
}

func (builder) Configure(msg proto.Message) error {
	fmt.Println("Handler level Configure called passing adapter config as proto.Message")
	return nil
}
func (builder) DefaultConfig() proto.Message { return &types.Empty{} }
