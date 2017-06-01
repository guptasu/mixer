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
	"fmt"
	adapter_cnfg "istio.io/mixer/pkg/adapter/config"
	"istio.io/mixer/pkg/adapterManager/wrappertypes"
	"istio.io/mixer/pkg/templates/mymetric/generated"
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
func dispatchMyMetricInstance(handler adapter_cnfg.Handler, vals []*wrappertypes.InstanceMakerInfo) interface{} {
	fmt.Println(vals)

	fmt.Println("about to call handler with instance", handler, vals)
	//handler.(mymetric.MyMetricProcessor).ProcessMyMetric(constructMyMetricInstance(vals))
	return nil
}

func constructMyMetricInstance(instanceMakers []*wrappertypes.InstanceMakerInfo) []mymetric.Instance {
	myMetricInstances := make([]mymetric.Instance, 0)
	for _, _ = range instanceMakers {
		////////////
		// TODO add all the expression evaluation and construction of metricInstance.
		///////////
		//params.
		metricInstance := mymetric.Instance{TypeName: "instantiating this type"}
		myMetricInstances = append(myMetricInstances, metricInstance)
	}
	return myMetricInstances
}
