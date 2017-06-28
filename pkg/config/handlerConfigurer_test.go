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
	"strings"
	"testing"

	pb "istio.io/mixer/pkg/config/proto"
)

func TestDedupeAndGroupInstances(t *testing.T) {
	tests := []struct {
		name         string
		actions      []*pb.Action
		constructors map[string]*pb.Constructor
		handlers     map[string]*HandlerBuilderInfo
		result       map[string]instancesByTemplate
		eError       string
	}{
		{
			name:         "SimpleNoDedupeNeeded",
			handlers:     map[string]*HandlerBuilderInfo{"h1": nil},
			constructors: map[string]*pb.Constructor{"i1": {"i1", "t1", nil}},
			actions:      []*pb.Action{{"h1", []string{"i1"}}},
			result: map[string]instancesByTemplate{
				"h1": {map[string][]string{"t1": {"i1"}}},
			},
			eError: "",
		},
		{
			name:     "DedupeAcrossActions",
			handlers: map[string]*HandlerBuilderInfo{"h1": nil},
			constructors: map[string]*pb.Constructor{
				"repeatInst": {"repeatInst", "t1", nil},
				"i2":         {"i2", "t1", nil}},
			actions: []*pb.Action{
				{"h1", []string{"repeatInst"}},
				{"h1", []string{"repeatInst", "i2"}}},
			result: map[string]instancesByTemplate{
				"h1": {map[string][]string{"t1": {"repeatInst", "i2"}}},
			},
			eError: "",
		},
		{
			name:     "DedupeWithinAction",
			handlers: map[string]*HandlerBuilderInfo{"h1": nil},
			constructors: map[string]*pb.Constructor{
				"repeatInst": {"repeatInst", "t1", nil},
				"i2":         {"i2", "t1", nil}},
			actions: []*pb.Action{
				{"h1", []string{"repeatInst", "repeatInst"}},
				{"h1", []string{"i2"}}},
			result: map[string]instancesByTemplate{
				"h1": {map[string][]string{"t1": {"repeatInst", "i2"}}},
			},
			eError: "",
		},
		{
			name:     "BadInstanceRef",
			handlers: map[string]*HandlerBuilderInfo{"h1": nil},
			constructors: map[string]*pb.Constructor{
				"i2": {"i2", "t1", nil}},
			actions: []*pb.Action{
				{"h1", []string{"badRefToInst"}},
			},
			result: nil,
			eError: "unable to find an a constructor with instance name 'badRefToInst'",
		},
		{
			name:     "BadHandlerRef",
			handlers: map[string]*HandlerBuilderInfo{},
			constructors: map[string]*pb.Constructor{
				"i2": {"i2", "t1", nil}},
			actions: []*pb.Action{
				{"badHandlerRef", []string{"i2"}},
			},
			result: nil,
			eError: "unable to find a configured handler with name 'badHandlerRef'",
		},
		{
			name:     "MultipleTemplates",
			handlers: map[string]*HandlerBuilderInfo{"h1": nil},
			constructors: map[string]*pb.Constructor{
				"i1tA": {"i1tA", "tA", nil},
				"i2tA": {"i2tA", "tA", nil},

				"i3tB": {"i3tB", "tB", nil},
				"i4tB": {"i4tB", "tB", nil},
				"i5tB": {"i5tB", "tB", nil},
			},
			actions: []*pb.Action{
				{"h1", []string{"i2tA", "i4tB", "i5tB", "i1tA", "i3tB"}},
			},
			result: map[string]instancesByTemplate{
				"h1": {map[string][]string{"tA": {"i2tA", "i1tA"}, "tB": {"i3tB", "i5tB", "i4tB"}}},
			},
			eError: "",
		},
		{
			name:     "UnionAcrossActionsWithMultipleTemplates",
			handlers: map[string]*HandlerBuilderInfo{"h1": nil, "h2": nil},
			constructors: map[string]*pb.Constructor{
				"i1tA": {"i1tA", "tA", nil},
				"i2tA": {"i2tA", "tA", nil},

				"i3tB": {"i3tB", "tB", nil},
				"i4tB": {"i4tB", "tB", nil},
				"i5tB": {"i5tB", "tB", nil},
			},
			actions: []*pb.Action{
				{"h1", []string{"i1tA", "i3tB"}},
				{"h1", []string{"i2tA", "i4tB", "i5tB"}},
				{"h2", []string{"i2tA", "i4tB", "i5tB", "i1tA", "i3tB"}},
			},
			result: map[string]instancesByTemplate{
				"h1": {map[string][]string{"tA": {"i2tA", "i1tA"}, "tB": {"i3tB", "i5tB", "i4tB"}}},
				"h2": {map[string][]string{"tA": {"i2tA", "i1tA"}, "tB": {"i3tB", "i5tB", "i4tB"}}},
			},
			eError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := dedupeAndGroupInstancesByTemplate(tt.actions, tt.constructors, tt.handlers)
			if tt.eError == "" {
				if err != nil {
					t.Errorf("got err %v\nwant <nil>", err)
				}
				if !equals(tt.result, v) {
					t.Errorf("got %v\nwant %v", v, tt.result)
				}
			} else {
				if err == nil || !strings.Contains(err.Error(), tt.eError) {
					t.Errorf("got error %v\nwant %v", err, tt.eError)
				}
			}

		})
	}

}

func equals(expected map[string]instancesByTemplate, actual map[string]instancesByTemplate) bool {
	if len(expected) != len(actual) {
		return false
	}

	for k, exTmplCnstMap := range expected {
		var actTmplCnstMap instancesByTemplate
		var ok bool
		if actTmplCnstMap, ok = actual[k]; !ok {
			return false
		}

		var exInstNamesByTmpl = exTmplCnstMap.instancesNamesByTemplate
		var actInstNamesByTmpl = actTmplCnstMap.instancesNamesByTemplate
		if len(exInstNamesByTmpl) != len(actInstNamesByTmpl) {
			return false
		}

		for exTmplName, exInsts := range exInstNamesByTmpl {
			var actInsts []string
			var ok bool
			if actInsts, ok = actInstNamesByTmpl[exTmplName]; !ok {
				return false
			}

			if len(exInsts) != len(actInsts) {
				return false
			}

			for _, exInst := range exInsts {
				if !contains(actInsts, exInst) {
					return false
				}
			}
		}
	}
	return true
}
