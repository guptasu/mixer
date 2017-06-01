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

package config

import "istio.io/mixer/pkg/adapter/config"
import (
	pb "istio.io/mixer/pkg/config/proto"
	"istio.io/mixer/pkg/templates/mymetric/generated"
	"istio.io/mixer/pkg/templates/mymetric/generated/config"
)

type typeConfigDispatcher struct {
	templateToTypes map[string]map[string]interface{}
}

func (t *typeConfigDispatcher) appendTypes(types []*pb.Type) error {
	for _, typ := range types {
		if _, ok := t.templateToTypes[typ.Template]; !ok {
			t.templateToTypes[typ.Template] = make(map[string]interface{})
		}
		t.templateToTypes[typ.Template][typ.Name] = typ.Params
	}
	// TODO guptasu: validate the types too. ensure they can be casted to appropriate types.
	return nil
}

func CreateTypeConfigDispatcher() (*typeConfigDispatcher, error) {
	templateToTypes := make(map[string]map[string]interface{})
	return &typeConfigDispatcher{templateToTypes: templateToTypes}, nil
}

func (t *typeConfigDispatcher) configureTypes(handler config.Handler) error {
	var x interface{} = handler
	if casted, ok := x.(mymetric.MyMetricProcessor); ok {
		defaultTyp := &foo_bar_mymetric.Type{}
		result := make(map[string]foo_bar_mymetric.Type)
		for k,v := range t.templateToTypes["foo.bar.mymetric.MyMetric"] {
			err := decode(v, defaultTyp, false)
			if err != nil {
				panic(err)
			}
			result[k] = *defaultTyp
		}
		casted.ConfigureMyMetric(result)
	}

	return nil
}
