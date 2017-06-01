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
	"istio.io/api/mixer/v1/config"
	//"istio.io/mixer/pkg/templates/mymetric/generated/config"
	"fmt"
)

type CreateInstance func([]*InstanceMakerInfo) interface{}

type InstanceMakerInfo struct {
	typeName         string
	templateName     string
	createInstanceFn CreateInstance
	params           interface{}
}

type RuntimeInstanceDispatcher struct {
	constructors map[string]*InstanceMakerInfo
}

// Should be instantiated during config time.
func CreateRuntimeInstanceDispatcher(
	types []istio_mixer_v1_config.Type,
	constructors []istio_mixer_v1_config.Constructor) {

	typeToTmplt := make(map[string]string)
	for _, typ := range types {
		typeToTmplt[typ.Name] = typ.Template
	}

	result := make(map[string]*InstanceMakerInfo)

	// TODO VALIDATION...
	for _, cnstr := range constructors {
		result[cnstr.Name] = &InstanceMakerInfo{
			templateName:     typeToTmplt[cnstr.Type],
			typeName:         cnstr.Type,
			createInstanceFn: tmplToInstanceCreator[typeToTmplt[cnstr.Type]],
			params:           cnstr.Params,
		}
	}
}

func (r *RuntimeInstanceDispatcher) getTemplate(instName string) string {
	return r.constructors[instName].templateName
}

func (r *RuntimeInstanceDispatcher) dispatchToHandler(actions []istio_mixer_v1_config.Action) {

	for _, action := range actions {
		//handlerName := action.Handler
		instanceNames := action.Instances

		// group instances into bucket for each template types.
		constGroupsPerTemplate := make(map[string][]*InstanceMakerInfo)
		for _, instName := range instanceNames {
			if _, ok := constGroupsPerTemplate[r.getTemplate(instName)]; !ok {
				constGroupsPerTemplate[r.getTemplate(instName)] = make([]*InstanceMakerInfo, 0)
			}
			val := r.constructors[instName]
			//r.constructors[instName].createInstanceFn()
			constGroupsPerTemplate[r.getTemplate(instName)] = append(constGroupsPerTemplate[r.getTemplate(instName)], val)
		}
	}
}

func (r *RuntimeInstanceDispatcher) evaluateConstructorParams(params interface{}) interface{} {
	return nil
}

/////////// ALL THE BELOW CODE IS GENERATED FROM TEMPLATES //////////////////

var tmplToInstanceCreator = map[string]CreateInstance{
	"foo.bar.mymetric.MyMetric": CreateMyMetricInstance,
}

func CreateMyMetricInstance(vals []*InstanceMakerInfo) interface{} {

	fmt.Println(vals)
	return nil
}
