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

package adapterManager

import (
	"istio.io/mixer/pkg/config/proto"
	"istio.io/mixer/pkg/adapterManager/wrappertypes"
	inst_dispatcher "istio.io/mixer/pkg/adapterManager/generated"
)


type RuntimeInstanceDispatcher struct {
	constructors map[string]*wrappertypes.InstanceMakerInfo
}

// Should be instantiated during config time.
func CreateRuntimeInstanceDispatcher(
	typeToTmplt map[string]string,
	constructors []*istio_mixer_v1_config.Constructor) *RuntimeInstanceDispatcher {


	result := make(map[string]*wrappertypes.InstanceMakerInfo)

	// TODO VALIDATION...
	for _, cnstr := range constructors {
		templateName := typeToTmplt[cnstr.Type]
		result[cnstr.Name] = &wrappertypes.InstanceMakerInfo{
			TypeName:         cnstr.Type,       // This is one of the parameters we pass as part of the Instance. It references to
							    // one of the configured Types (Types are configured during config time).
			Params:           cnstr.Params,     // Constructor.Params
			TemplateName:     templateName,     // Template is needed to group the instances before passing to Adapters
		}
	}
	return &RuntimeInstanceDispatcher{constructors: result}
}

func (r *RuntimeInstanceDispatcher) DispatchToHandler(actions []*istio_mixer_v1_config.Action) {
	// TODO add go routines here to fan out.

	for _, action := range actions {
		// group instances into bucket for each template types.
		constructorsGroupedPerTemplate := r.getConstructorsGroupedPerTemplate(action.Instances)

		for templateName, instanceMakerInfos := range constructorsGroupedPerTemplate {
			c,r := inst_dispatcher.GetDispatchMethod(templateName)
			if c != nil {
				c(nil, instanceMakerInfos)
			}else {
				r(nil, instanceMakerInfos)
			}
		}
	}
}

func (r *RuntimeInstanceDispatcher) getConstructorsGroupedPerTemplate(instanceNames []string) map[string][]*wrappertypes.InstanceMakerInfo {
	constGroupsPerTemplate := make(map[string][]*wrappertypes.InstanceMakerInfo)
	for _, instName := range instanceNames {
		if _, ok := constGroupsPerTemplate[r.getTemplate(instName)]; !ok {
			constGroupsPerTemplate[r.getTemplate(instName)] = make([]*wrappertypes.InstanceMakerInfo, 0)
		}
		val := r.constructors[instName]
		//r.constructors[instName].createInstanceFn()
		constGroupsPerTemplate[r.getTemplate(instName)] = append(constGroupsPerTemplate[r.getTemplate(instName)], val)
	}
	return constGroupsPerTemplate
}

func (r *RuntimeInstanceDispatcher) getTemplate(instName string) string {
	return r.constructors[instName].TemplateName
}
