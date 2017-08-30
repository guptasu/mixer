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

	"context"
	"istio.io/mixer/pkg/adapter"
	reportTmpl "istio.io/mixer/test/e2e/template/report"
)

type (
	fakeHndlr struct {
		hndlrbehavior hndlrBehavior
		hndlrCallData *hndlrCallData
	}

	fakeBldr struct {
		bldrbehavior  builderBehavior
		bldrCallData  *bldrCallData

		// builder needs handler behavior and data to pass it along during Build call.
		hndlrbehavior hndlrBehavior
		hndlrCallData *hndlrCallData
	}

	spyAdapter struct {
		behavior      AdptBehavior
		bldrCallData  *bldrCallData
		hndlrCallData *hndlrCallData
	}

	AdptBehavior struct {
		name          string
		hndlrBehavior hndlrBehavior
		bldrBehavior  builderBehavior
	}

	hndlrBehavior struct {
		HandleSampleReport_error error
		HandleSampleReport_panic bool

		Close_error error
		Close_panic bool
	}

	hndlrCallData struct {
		HandleSampleReport_instances []*reportTmpl.Instance
		HandleSampleReport_invoked   bool

		Close_invoked bool
	}

	builderBehavior struct {
		ConfigureSampleReportHandler_err   error
		ConfigureSampleReportHandler_panic bool

		Build_err   error
		Build_panic bool
	}

	bldrCallData struct {
		ConfigureSampleReportHandler_invoked bool
		ConfigureSampleReportHandler_types   map[string]*reportTmpl.Type

		Build_invoked  bool
		Build_adptCnfg adapter.Config
	}
)

func (f fakeBldr) Build(cnfg adapter.Config, _ adapter.Env) (adapter.Handler, error) {
	f.bldrCallData.Build_invoked = true
	if f.bldrbehavior.Build_panic {
		panic("Build")
	}

	f.bldrCallData.Build_adptCnfg = cnfg
	hndlr := fakeHndlr{hndlrbehavior: f.hndlrbehavior, hndlrCallData: f.hndlrCallData}
	return hndlr, f.bldrbehavior.Build_err
}
func (f fakeBldr) ConfigureSampleReportHandler(typeParams map[string]*reportTmpl.Type) error {
	f.bldrCallData.ConfigureSampleReportHandler_invoked = true
	if f.bldrbehavior.ConfigureSampleReportHandler_panic {
		panic("ConfigureSampleReportHandler")
	}

	f.bldrCallData.ConfigureSampleReportHandler_types = typeParams
	return f.bldrbehavior.ConfigureSampleReportHandler_err
}

func (f fakeHndlr) HandleSampleReport(ctx context.Context, instances []*reportTmpl.Instance) error {
	f.hndlrCallData.HandleSampleReport_invoked = true
	if f.hndlrbehavior.HandleSampleReport_panic {
		panic("HandleSampleReport")
	}

	f.hndlrCallData.HandleSampleReport_instances = instances
	return f.hndlrbehavior.HandleSampleReport_error
}

func (f fakeHndlr) Close() error {
	f.hndlrCallData.Close_invoked = true
	if f.hndlrbehavior.Close_panic {
		panic("Close")
	}

	return f.hndlrbehavior.Close_error
}

func newSpyAdapter(b AdptBehavior) *spyAdapter {
	return &spyAdapter{behavior: b, bldrCallData: &bldrCallData{}, hndlrCallData: &hndlrCallData{}}
}

func (s *spyAdapter) getAdptInfoFn() adapter.InfoFn {
	return func() adapter.BuilderInfo {
		return adapter.BuilderInfo{
			Name:               s.behavior.name,
			Description:        "",
			SupportedTemplates: []string{reportTmpl.TemplateName},
			CreateHandlerBuilder: func() adapter.HandlerBuilder {
				return fakeBldr{
					bldrbehavior:  s.behavior.bldrBehavior,
					bldrCallData:  s.bldrCallData,
					hndlrbehavior: s.behavior.hndlrBehavior,
					hndlrCallData: s.hndlrCallData,
				}
			},
			DefaultConfig: &types.Empty{},
			ValidateConfig: func(msg adapter.Config) *adapter.ConfigErrors {
				return nil
			},
		}
	}
}
