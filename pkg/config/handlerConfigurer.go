// Copyright 2017 Istio Authors
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

import (
	"fmt"

	pb "istio.io/mixer/pkg/config/proto"
)

func configureHandlers(actions []*pb.Action, constructors map[string]*pb.Constructor,
	handlers map[string]*HandlerBuilderInfo) error {
	_, err := groupHandlerInstancesByTemplate(actions, constructors, handlers)
	// TODO Add type inference and dispatch to adapter code.
	return err
}

type instancesByTemplate struct {
	instancesNamesByTemplate map[string][]string
}

func (t *instancesByTemplate) insertInstance(instName string, tmplName string) {
	// TODO validate the tmplName and if handler supports the template
	instsPerTmpl, alreadyPresent := t.instancesNamesByTemplate[tmplName]
	if !alreadyPresent {
		t.instancesNamesByTemplate[tmplName] = make([]string, 0)
	}

	// Add the instance only if does not already exists
	if !contains(instsPerTmpl, instName) {
		t.instancesNamesByTemplate[tmplName] = append(t.instancesNamesByTemplate[tmplName], instName)
	}
}

func newInstancesByTemplateMapping() instancesByTemplate {
	return instancesByTemplate{make(map[string][]string)}
}

func groupHandlerInstancesByTemplate(actions []*pb.Action, constructors map[string]*pb.Constructor,
	handlers map[string]*HandlerBuilderInfo) (map[string]instancesByTemplate, error) {
	result := make(map[string]instancesByTemplate)

	for _, action := range actions {
		hName := action.GetHandler()
		if _, ok := handlers[hName]; !ok {
			return nil, fmt.Errorf("unable to find a configured handler with name '%s' referenced in action %v", hName, action)
		}

		tmplCnstrsMapping, alreadyPresent := result[hName]
		if !alreadyPresent {
			tmplCnstrsMapping = newInstancesByTemplateMapping()
			result[hName] = tmplCnstrsMapping
		}

		for _, iName := range action.GetInstances() {
			var cnstr *pb.Constructor
			var ok bool
			if cnstr, ok = constructors[iName]; !ok {
				return nil, fmt.Errorf("unable to find an a constructor with instance name '%s' "+
					"referenced in action %v", iName, action)
			}

			tmplCnstrsMapping.insertInstance(iName, cnstr.GetTemplate())

		}
	}
	return result, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
