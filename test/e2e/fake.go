// Copyright 2016 Istio Authors
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

package e2e

import (
	"github.com/gogo/protobuf/types"

	"istio.io/mixer/pkg/adapter"
	reportTmpl "istio.io/mixer/test/e2e/template/report"
)

type (
	fakeHndlr struct {
		callData map[string]interface{}
	}
	fakeHndlrBldr struct{
		callData map[string]interface{}
	}
)

func (f fakeHndlrBldr) Build(cnfg adapter.Config, _ adapter.Env) (adapter.Handler, error) {
	f.callData["Build"] = cnfg
	fakeHndlrObj := fakeHndlr{}
	return fakeHndlrObj, nil
}
func (f fakeHndlrBldr) ConfigureSampleReportHandler(typeParams map[string]*reportTmpl.Type) error {
	f.callData["ConfigureSampleReport"] = typeParams
	return nil
}
func (f fakeHndlr) HandleSampleReport(instances []*reportTmpl.Instance) error {
	f.callData["HandleSampleReport"] = instances
	return nil
}
func (f fakeHndlr) Close() error {
	f.callData["Close"] = nil
	return nil
}

type spyAdapter struct {
	behavior *AdptBehavior
	builderCallData *builderCallData
}

type AdptBehavior struct {
	name string
}

type builderCallData struct {
	data map[string]interface{}
}

func (s *spyAdapter) getFakeHndlrBldrInfoFn() adapter.InfoFn {
	return func() adapter.BuilderInfo {
		return adapter.BuilderInfo{
			Name:                 s.behavior.name,
			Description:          "",
			SupportedTemplates:   []string{reportTmpl.TemplateName},
			CreateHandlerBuilder: func() adapter.HandlerBuilder { return fakeHndlrBldr{callData: s.builderCallData.data} },
			DefaultConfig:        &types.Empty{},
			ValidateConfig: func(msg adapter.Config) *adapter.ConfigErrors {
				return nil
			},
		}
	}
}
