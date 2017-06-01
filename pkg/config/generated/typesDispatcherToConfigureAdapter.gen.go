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
	"istio.io/mixer/pkg/templates/mymetric/generated"
	"istio.io/mixer/pkg/templates/mymetric/generated/config"
	"github.com/golang/protobuf/proto"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gogo/protobuf/jsonpb"
)

type TypeFinder interface {
	GetTypesForTemplate(templateName string) map[string]interface{}
}


func ConfigureTypes(handler config.Handler, t TypeFinder) error {
	var x interface{} = handler
	if casted, ok := x.(mymetric.MyMetricProcessor); ok {
		defaultTyp := &foo_bar_mymetric.Type{}
		result := make(map[string]foo_bar_mymetric.Type)
		for typeName, typeParam := range t.GetTypesForTemplate("foo.bar.mymetric.MyMetric") {
			err := decode(typeParam, defaultTyp, false)
			if err != nil {
				panic(err)
			}
			result[typeName] = *defaultTyp
		}
		casted.ConfigureMyMetric(result)
	}

	return nil
}

func decode(src interface{}, dst proto.Message, strict bool) error {
	ba, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("failed to marshal config into json: %v", err)
	}
	um := jsonpb.Unmarshaler{AllowUnknownFields: !strict}
	if err := um.Unmarshal(bytes.NewReader(ba), dst); err != nil {
		b2, _ := json.Marshal(dst)
		return fmt.Errorf("failed to unmarshal config <%s> into proto: %v %s", string(ba), err, string(b2))
	}
	return nil
}

