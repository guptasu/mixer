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
	"reflect"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/wrappers"

	pb "istio.io/mixer/pkg/config/proto"
	"istio.io/mixer/pkg/expr"
	tmpl "istio.io/mixer/pkg/template"
)

type fakeTmplRepo struct {
	tmptExists      bool
	inferTypeError  error
	inferTypeResult proto.Message
}

func newFakeTmplRepo(inferError error, inferResult proto.Message, tmplFound bool) tmpl.Repository {
	return fakeTmplRepo{inferTypeError: inferError, inferTypeResult: inferResult, tmptExists: tmplFound}
}
func (t fakeTmplRepo) GetTemplateInfo(template string) (tmpl.Info, bool) {
	return tmpl.Info{
		InferTypeFn:    func(interface{}, tmpl.TypeEvalFn) (proto.Message, error) { return t.inferTypeResult, t.inferTypeError },
		CnstrDefConfig: nil,
	}, t.tmptExists
}

func TestInferTypes(t *testing.T) {
	tests := []struct {
		name         string
		constructors map[string]*pb.Constructor
		tmplRepo     tmpl.Repository
		want         map[string]proto.Message
		wantError    string
	}{
		{
			name:         "SingleCnstr",
			constructors: map[string]*pb.Constructor{"inst1": {"inst1", "tpml1", nil}},
			tmplRepo:     newFakeTmplRepo(nil, &wrappers.Int32Value{Value: 1}, true),
			want:         map[string]proto.Message{"inst1": &wrappers.Int32Value{Value: 1}},
			wantError:    "",
		},
		{
			name: "MultipleCnstr",
			constructors: map[string]*pb.Constructor{"inst1": {"inst1", "tpml1", nil},
				"inst2": {"inst2", "tpml1", nil}},
			tmplRepo:  newFakeTmplRepo(nil, &wrappers.Int32Value{Value: 1}, true),
			want:      map[string]proto.Message{"inst1": &wrappers.Int32Value{Value: 1}, "inst2": &wrappers.Int32Value{Value: 1}},
			wantError: "",
		},
		{
			name:         "InvalidTmplInCnstr",
			constructors: map[string]*pb.Constructor{"inst1": {"inst1", "INVALIDTMPLNAME", nil}},
			tmplRepo:     newFakeTmplRepo(fmt.Errorf("invalid tmpl"), nil, false),
			want:         nil,
			wantError:    "is not registered",
		},
		{
			name:         "ErrorDuringTypeInfr",
			constructors: map[string]*pb.Constructor{"inst1": {"inst1", "tpml1", nil}},
			tmplRepo:     newFakeTmplRepo(fmt.Errorf("error during type infer"), nil, true),
			want:         nil,
			wantError:    "cannot infer type information",
		},
	}
	ex, _ := expr.NewCEXLEvaluator(expr.DefaultCacheSize)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hc := handlerConfigurer{typeChecker: ex, tmplRepo: tt.tmplRepo}
			v, err := hc.inferTypes(tt.constructors)
			if tt.wantError == "" {
				if err != nil {
					t.Errorf("got err %v\nwant <nil>", err)
				}
				if !reflect.DeepEqual(tt.want, v) {
					t.Errorf("got %v\nwant %v", v, tt.want)
				}
			} else {
				if err == nil || !strings.Contains(err.Error(), tt.wantError) {
					t.Errorf("got error %v\nwant %v", err, tt.wantError)
				}
			}

		})
	}
}

func TestDedupeAndGroupInstances(t *testing.T) {
	tests := []struct {
		name         string
		actions      []*pb.Action
		constructors map[string]*pb.Constructor
		handlers     map[string]*HandlerBuilderInfo
		want         map[string]instancesByTemplate
		wantError    string
	}{
		{
			name:         "SimpleNoDedupeNeeded",
			handlers:     map[string]*HandlerBuilderInfo{"hndlr1": nil},
			constructors: map[string]*pb.Constructor{"i1": {"i1", "tpml1", nil}},
			actions:      []*pb.Action{{"hndlr1", []string{"i1"}}},
			want: map[string]instancesByTemplate{
				"hndlr1": {map[string][]string{"tpml1": {"i1"}}},
			},
			wantError: "",
		},
		{
			name:     "DedupeAcrossActions",
			handlers: map[string]*HandlerBuilderInfo{"hndlr1": nil},
			constructors: map[string]*pb.Constructor{
				"repeatInst": {"repeatInst", "tpml1", nil},
				"inst2":      {"inst2", "tpml1", nil}},
			actions: []*pb.Action{
				{"hndlr1", []string{"repeatInst"}},
				{"hndlr1", []string{"repeatInst", "inst2"}}},
			want: map[string]instancesByTemplate{
				"hndlr1": {map[string][]string{"tpml1": {"repeatInst", "inst2"}}},
			},
			wantError: "",
		},
		{
			name:     "DedupeWithinAction",
			handlers: map[string]*HandlerBuilderInfo{"hndlr1": nil},
			constructors: map[string]*pb.Constructor{
				"repeatInst": {"repeatInst", "tpml1", nil},
				"inst2":      {"inst2", "tpml1", nil}},
			actions: []*pb.Action{
				{"hndlr1", []string{"repeatInst", "repeatInst"}},
				{"hndlr1", []string{"inst2"}}},
			want: map[string]instancesByTemplate{
				"hndlr1": {map[string][]string{"tpml1": {"repeatInst", "inst2"}}},
			},
			wantError: "",
		},
		{
			name:     "BadInstanceRef",
			handlers: map[string]*HandlerBuilderInfo{"hndlr1": nil},
			constructors: map[string]*pb.Constructor{
				"inst2": {"inst2", "tpml1", nil}},
			actions: []*pb.Action{
				{"hndlr1", []string{"badRefToInst"}},
			},
			want:      nil,
			wantError: "unable to find an a constructor with instance name 'badRefToInst'",
		},
		{
			name:     "BadHandlerRef",
			handlers: map[string]*HandlerBuilderInfo{},
			constructors: map[string]*pb.Constructor{
				"inst2": {"inst2", "tpml1", nil}},
			actions: []*pb.Action{
				{"badHandlerRef", []string{"inst2"}},
			},
			want:      nil,
			wantError: "unable to find a configured handler with name 'badHandlerRef'",
		},
		{
			name:     "MultipleTemplates",
			handlers: map[string]*HandlerBuilderInfo{"hndlr1": nil},
			constructors: map[string]*pb.Constructor{
				"inst1tmplA": {"inst1tmplA", "tmplA", nil},
				"inst2tmplA": {"inst2tmplA", "tmplA", nil},

				"inst3tmplB": {"inst3tmplB", "tmplB", nil},
				"inst4tmplB": {"inst4tmplB", "tmplB", nil},
				"inst5tmplB": {"inst5tmplB", "tmplB", nil},
			},
			actions: []*pb.Action{
				{"hndlr1", []string{"inst2tmplA", "inst4tmplB", "inst5tmplB", "inst1tmplA", "inst3tmplB"}},
			},
			want: map[string]instancesByTemplate{
				"hndlr1": {map[string][]string{"tmplA": {"inst2tmplA", "inst1tmplA"}, "tmplB": {"inst3tmplB", "inst5tmplB", "inst4tmplB"}}},
			},
			wantError: "",
		},
		{
			name:     "UnionAcrossActionsWithMultipleTemplates",
			handlers: map[string]*HandlerBuilderInfo{"hndlr1": nil, "hndlr2": nil},
			constructors: map[string]*pb.Constructor{
				"inst1tmplA": {"inst1tmplA", "tmplA", nil},
				"inst2tmplA": {"inst2tmplA", "tmplA", nil},

				"inst3tmplB": {"inst3tmplB", "tmplB", nil},
				"inst4tmplB": {"inst4tmplB", "tmplB", nil},
				"inst5tmplB": {"inst5tmplB", "tmplB", nil},
			},
			actions: []*pb.Action{
				{"hndlr1", []string{"inst1tmplA", "inst3tmplB"}},
				{"hndlr1", []string{"inst2tmplA", "inst4tmplB", "inst5tmplB"}},
				{"hndlr2", []string{"inst2tmplA", "inst4tmplB", "inst5tmplB", "inst1tmplA", "inst3tmplB"}},
			},
			want: map[string]instancesByTemplate{
				"hndlr1": {map[string][]string{"tmplA": {"inst2tmplA", "inst1tmplA"}, "tmplB": {"inst3tmplB", "inst5tmplB", "inst4tmplB"}}},
				"hndlr2": {map[string][]string{"tmplA": {"inst2tmplA", "inst1tmplA"}, "tmplB": {"inst3tmplB", "inst5tmplB", "inst4tmplB"}}},
			},
			wantError: "",
		},
	}
	ex, _ := expr.NewCEXLEvaluator(expr.DefaultCacheSize)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			hc := handlerConfigurer{typeChecker: ex, tmplRepo: nil}
			v, err := hc.groupHandlerInstancesByTemplate(tt.actions, tt.constructors, tt.handlers)
			if tt.wantError == "" {
				if err != nil {
					t.Errorf("got err %v\nwant <nil>", err)
				}
				if !equals(tt.want, v) {
					t.Errorf("got %v\nwant %v", v, tt.want)
				}
			} else {
				if err == nil || !strings.Contains(err.Error(), tt.wantError) {
					t.Errorf("got error %v\nwant %v", err, tt.wantError)
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
