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

package adapterManager

import (
	adapter_cnfg "istio.io/mixer/pkg/adapter/config"
	"istio.io/mixer/pkg/adapterManager/wrappertypes"
	//"istio.io/mixer/pkg/expr"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"istio.io/mixer/pkg/attribute"
	"istio.io/mixer/pkg/templates/mymetric/generated"
	"istio.io/mixer/pkg/templates/mymetric/generated/config"
	"istio.io/mixer/pkg/expr"
	"istio.io/mixer/pkg/aspect"
)

// TODO: use correc type for kind
var TmplToMethodKind = map[string]string{
	"foo.bar.mymetric.MyMetric": "report",
}

var TmplToReportFnDispatcher = map[string]wrappertypes.InstanceToReportFnDispatcher{
	"foo.bar.mymetric.MyMetric": dispatchMyMetricInstance,
}

var TmplToCheckFnDispatcher = map[string]wrappertypes.InstanceToCheckFnDispatcher{}

func GetDispatchMethod(templateName string) (wrappertypes.InstanceToReportFnDispatcher, wrappertypes.InstanceToCheckFnDispatcher) {

	if kind, ok := TmplToMethodKind[templateName]; ok {
		if kind == "report" {
			if fn, ok := TmplToReportFnDispatcher[templateName]; ok {
				return fn, nil
			} else {
				// TODO..
				panic("should not happen as config should be already valid")
			}
		} else {
			if fn, ok := TmplToCheckFnDispatcher[templateName]; ok {
				return nil, fn
			} else {
				// TODO..
				panic("should not happen as config should be already valid")
			}
		}
	} else {
		panic("should not happen as config should be already valid")
	}
}

//////////////////////////// GENERATED FROM MYMETRIC TEMPLATE //////////////////////////////////////////////
func dispatchMyMetricInstance(handler adapter_cnfg.Handler, vals []*wrappertypes.InstanceMakerInfo, reqBag *attribute.MutableBag) interface{} {
	//_ = constructMyMetricInstance(vals)
	handler.(mymetric.MyMetricProcessor).ProcessMyMetric(constructMyMetricInstance(vals, reqBag))
	return nil
}

func constructMyMetricInstance(instanceMakers []*wrappertypes.InstanceMakerInfo, reqBag *attribute.MutableBag) []mymetric.Instance {
	myMetricInstances := make([]mymetric.Instance, 0)
	for _, instanceMaker := range instanceMakers {
		data := &foo_bar_mymetric.Constructor{}
		err := decode(instanceMaker.Params, data, false)
		if err != nil {
			panic(err)
		}
		ex, _ := expr.NewCEXLEvaluator(expr.DefaultCacheSize)
		v, _ := ex.Eval(data.Value, reqBag)
		//if err != nil {
		//	panic(err)
		//}
		dimensions,_ := aspect.EvalAll(data.Dimensions, reqBag, ex)
		metricInstance := mymetric.Instance{TypeName: instanceMaker.TypeName, Value: v, Dimensions:dimensions}
		myMetricInstances = append(myMetricInstances, metricInstance)
	}
	return myMetricInstances
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
