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
	fakeHndlr     struct{}
	fakeHndlrBldr struct{}
)

func (fakeHndlrBldr) Build(cnfg adapter.Config, _ adapter.Env) (adapter.Handler, error) {
	globalActualHandlerCallInfoToValidate["Build"] = cnfg
	fakeHndlrObj := fakeHndlr{}
	return fakeHndlrObj, nil
}
func (fakeHndlrBldr) ConfigureSampleReportHandler(typeParams map[string]*reportTmpl.Type) error {
	globalActualHandlerCallInfoToValidate["ConfigureSampleReport"] = typeParams
	return nil
}
func (fakeHndlr) HandleSampleReport(instances []*reportTmpl.Instance) error {
	globalActualHandlerCallInfoToValidate["HandleSampleReport"] = instances
	return nil
}
func (fakeHndlr) Close() error {
	globalActualHandlerCallInfoToValidate["Close"] = nil
	return nil
}
func GetFakeHndlrBuilderInfo() adapter.BuilderInfo {
	return adapter.BuilderInfo{
		Name:                 "fakeHandler",
		Description:          "",
		SupportedTemplates:   []string{reportTmpl.TemplateName},
		CreateHandlerBuilder: func() adapter.HandlerBuilder { return fakeHndlrBldr{} },
		DefaultConfig:        &types.Empty{},
		ValidateConfig: func(msg adapter.Config) *adapter.ConfigErrors {
			return nil
		},
	}
}
