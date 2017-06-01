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

package wrappertypes

import "istio.io/mixer/pkg/adapter/config"

// TODO: the two types should be different.
type InstanceToReportFnDispatcher func(handler config.Handler, instanceMakers []*InstanceMakerInfo) interface{}
type InstanceToCheckFnDispatcher func(handler config.Handler, instanceMakers []*InstanceMakerInfo) interface{}

type InstanceMakerInfo struct {
	TypeName         string
	TemplateName     string
	Params           interface{}
}
