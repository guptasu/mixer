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

package config

import "istio.io/mixer/pkg/adapter/config"
import (
	pb "istio.io/mixer/pkg/config/proto"
	type_dispatcher_generated "istio.io/mixer/pkg/config/generated"
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

func (t *typeConfigDispatcher) GetTypeToTemplateMapping() map[string]string {
	allTypesToTemplate := make(map[string]string)
	for templateName, types := range t.templateToTypes {
		for k,_ := range types {
			allTypesToTemplate[k] = templateName
		}
	}
	return allTypesToTemplate
}

func CreateTypeConfigDispatcher() (*typeConfigDispatcher, error) {
	templateToTypes := make(map[string]map[string]interface{})
	return &typeConfigDispatcher{templateToTypes: templateToTypes}, nil
}

// return all types associated with a template map[typeName]typeParam
func (t typeConfigDispatcher) GetTypesForTemplate(templateName string) map[string]interface{} {
	return t.templateToTypes[templateName]
}

func (t *typeConfigDispatcher) configureTypes(handler config.Handler) error {
	type_dispatcher_generated.ConfigureTypes(handler, *t)
	return nil;
}
