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

package aspect

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	rpc "github.com/googleapis/googleapis/google/rpc"

	"istio.io/mixer/pkg/adapter"
	config2 "istio.io/mixer/pkg/adapter/config"
	"istio.io/mixer/pkg/attribute"
	"istio.io/mixer/pkg/config"
	"istio.io/mixer/pkg/config/descriptor"
	cpb "istio.io/mixer/pkg/config/proto"
	"istio.io/mixer/pkg/expr"
	"istio.io/mixer/pkg/template"
)

type (
	checkManager struct{}

	checkExecutor struct {
		hndlrName    string
		tmplName     string
		procDispatch template.ProcessCheckFn
		hndlr        config2.Handler
		ctrs         map[string]proto.Message // constructor name -> constructor params
	}
)

func newCheckManagerImpl() CheckManager {
	return &checkManager{}
}

func (m *checkManager) NewCheckExecutor(c *cpb.Combined, createAspect CreateAspectFunc, env adapter.Env, df descriptor.Finder, repo template.Repository) (CheckExecutor, error) {
	ctrs := make(map[string]proto.Message)
	for _, cstr := range c.Constructors {
		ctrs[cstr.InstanceName] = cstr.Params.(proto.Message)
	}

	tmpl := ""
	if len(c.Constructors) > 0 {
		// All constructors should have the same template name. The adapter manager would have already grouped them by
		// template.
		tmpl = c.Constructors[0].Template
	} else {
		return nil, fmt.Errorf("NewCheckExecutor instantiated with empty constructor list")
	}

	out, err := createAspect(env, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to construct check aspect with config '%v': %v", c, err)
	}

	v, ok := out.(config2.Handler)
	if !ok {
		return nil, fmt.Errorf("wrong aspect type returned after creation; expected MetricsAspect: %#v", out)
	}

	ti, _ := repo.GetTemplateInfo(tmpl)
	if b := ti.HandlerSupportsTemplate(v); !b {
		return nil, fmt.Errorf("Handler does not implement interface %s. "+
			"Therefore, it cannot support template %v", ti.HndlrName, tmpl)
	}

	return &checkExecutor{c.Builder.Name, tmpl, ti.ProcessCheck, v, ctrs}, nil
}

func (*checkManager) DefaultConfig() config.AspectParams { return nil }
func (*checkManager) ValidateConfig(c config.AspectParams, tc expr.TypeChecker, df descriptor.Finder) (ce *adapter.ConfigErrors) {
	return
}

func (*checkManager) Kind() config.Kind {
	return config.Undefined
}

func (w *checkExecutor) Execute(attrs attribute.Bag, mapper expr.Evaluator) rpc.Status {
	s, _ := w.procDispatch(w.ctrs, attrs, mapper, w.hndlr) // ignore Cacheability info for now.
	return s
}

func (w *checkExecutor) Close() error {
	return w.hndlr.Close()
}
