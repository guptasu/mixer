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
	"crypto/sha1"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/gogo/protobuf/types"
	"github.com/golang/protobuf/proto"

	dpb "istio.io/api/mixer/v1/config/descriptor"
	"istio.io/mixer/pkg/adapter"
	"istio.io/mixer/pkg/adapter/config"
	listcheckerpb "istio.io/mixer/pkg/aspect/config"
	"istio.io/mixer/pkg/attribute"
	"istio.io/mixer/pkg/config/descriptor"
	pb "istio.io/mixer/pkg/config/proto"
	"istio.io/mixer/pkg/expr"
	tmpl "istio.io/mixer/pkg/template"
)

type fakeVFinder struct {
	ada   map[string]adapter.ConfigValidator
	hbi   map[string]*adapter.BuilderInfo
	asp   map[Kind]AspectValidator
	kinds KindSet
}

func (f *fakeVFinder) FindAdapterValidator(name string) (adapter.ConfigValidator, bool) {
	v, found := f.ada[name]
	return v, found
}

func (f *fakeVFinder) FindBuilderInfo(name string) (*adapter.BuilderInfo, bool) {
	v, found := f.hbi[name]
	return v, found
}

func (f *fakeVFinder) FindAspectValidator(kind Kind) (AspectValidator, bool) {
	v, found := f.asp[kind]
	return v, found
}

func (f *fakeVFinder) AdapterToAspectMapperFunc(builder string) KindSet {
	return f.kinds
}

type lc struct {
	ce *adapter.ConfigErrors
}

func (m *lc) DefaultConfig() (c adapter.Config) {
	return &listcheckerpb.ListsParams{}
}

// ValidateConfig determines whether the given configuration meets all correctness requirements.
func (m *lc) ValidateConfig(c adapter.Config) *adapter.ConfigErrors {
	return m.ce
}

type ac struct {
	ce *adapter.ConfigErrors
}

func (*ac) DefaultConfig() AspectParams {
	return &listcheckerpb.ListsParams{}
}

// ValidateConfig determines whether the given configuration meets all correctness requirements.
func (a *ac) ValidateConfig(AspectParams, expr.TypeChecker, descriptor.Finder) *adapter.ConfigErrors {
	return a.ce
}

type configTable struct {
	cerr     *adapter.ConfigErrors
	ada      map[string]adapter.ConfigValidator
	hbi      map[string]*adapter.BuilderInfo
	asp      map[Kind]AspectValidator
	nerrors  int
	selector string
	strict   bool
	cfg      string
}

func newVfinder(ada map[string]adapter.ConfigValidator, asp map[Kind]AspectValidator,
	hbi map[string]*adapter.BuilderInfo) *fakeVFinder {
	var kinds KindSet
	for k := range asp {
		kinds = kinds.Set(k)
	}
	return &fakeVFinder{ada: ada, hbi: hbi, asp: asp, kinds: kinds}
}

func fakeConfigureHandler(actions []*pb.Action, constructors map[string]*pb.Constructor,
	handlers map[string]*HandlerBuilderInfo, tmplRepo tmpl.Repository, expr expr.TypeChecker, df expr.AttributeDescriptorFinder) error {
	return nil
}

func TestConfigValidatorError(t *testing.T) {
	var ct *adapter.ConfigErrors
	evaluator := newFakeExpr()
	cerr := ct.Appendf("url", "Must have a valid URL")

	tests := []*configTable{
		{nil,
			map[string]adapter.ConfigValidator{
				"denyChecker": &lc{},
				"metrics2":    &lc{},
			}, nil,
			nil, 0, "service.name == “*”", false, ConstGlobalConfig},
		{nil,
			map[string]adapter.ConfigValidator{
				"metrics":  &lc{},
				"metrics2": &lc{},
			}, nil,
			nil, 1, "service.name == “*”", false, ConstGlobalConfig},
		{nil, nil, nil,
			map[Kind]AspectValidator{
				MetricsKind: &ac{},
				QuotasKind:  &ac{},
			},
			0, "service.name == “*”", false, sSvcConfig},
		{nil, nil, nil,
			map[Kind]AspectValidator{
				MetricsKind: &ac{},
				QuotasKind:  &ac{},
			},
			1, "service.name == “*”", true, sSvcConfig},
		{cerr, nil, nil,
			map[Kind]AspectValidator{
				QuotasKind: &ac{ce: cerr},
			},
			2, "service.name == “*”", false, sSvcConfig},
		{ct.Append("/:metrics", unknownValidator("metrics")),
			nil, nil, nil, 2, "\"\"", false, sSvcConfig},
	}

	for idx, tt := range tests {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {

			var ce *adapter.ConfigErrors
			mgr := newVfinder(tt.ada, tt.asp, tt.hbi)
			p := newValidator(mgr.FindAspectValidator, mgr.FindAdapterValidator, mgr.FindBuilderInfo, fakeConfigureHandler, nil,
				mgr.AdapterToAspectMapperFunc, tt.strict, evaluator)
			if tt.cfg == sSvcConfig {
				ce = p.validateServiceConfig(globalRulesKey, fmt.Sprintf(tt.cfg, tt.selector), false)
			} else {
				ce = p.validateAdapters(keyAdapters, tt.cfg)
			}
			cok := ce == nil
			ok := tt.nerrors == 0

			if ok != cok {
				t.Fatalf("got %t, want %t", cok, ok)
			}
			if ce == nil {
				return
			}

			if len(ce.Multi.Errors) != tt.nerrors {
				t.Fatalf("got %s, want %s", ce.Error(), tt.cerr.Error())
			}
		})
	}
}

func TestConvertHandlerParamsErrors(t *testing.T) {

	tTable := []struct {
		params         interface{}
		defaultCnfg    proto.Message
		validateConfig adapter.ValidateConfig
		errorStr       string
	}{
		{
			map[string]interface{}{
				"check_expression": "src.ip",
				"blacklist":        "true (this should be bool)",
			},
			&listcheckerpb.ListsParams{},
			func(c proto.Message) error { return nil },
			"failed to decode",
		},
		{
			map[string]interface{}{
				"check_expression": "src.ip",
				"blacklist":        true,
				"wrongextrafield":  true,
			},
			&listcheckerpb.ListsParams{},
			func(c proto.Message) error { return nil },
			"failed to unmarshal",
		},
		{
			map[string]interface{}{
				"check_expression": "src.ip",
				"blacklist":        true,
			},
			&listcheckerpb.ListsParams{},
			func(c proto.Message) error { return errors.New("handler config validation failed") },
			"handler config validation failed",
		},
	}

	for _, tt := range tTable {
		t.Run(tt.errorStr, func(t *testing.T) {
			_, ce := convertHandlerParams(
				&adapter.BuilderInfo{
					DefaultConfig:  tt.defaultCnfg,
					ValidateConfig: tt.validateConfig,
				}, "myhandlerCnfgBlock", tt.params, true)

			if !strings.Contains(ce.Error(), tt.errorStr) {
				t.Errorf("got %s, want %s\n", ce.Error(), tt.errorStr)
			}
		})
	}
}

func TestHandlerConfigValidator(t *testing.T) {
	evaluator := newFakeExpr()

	tests := []*configTable{
		{
			nil,
			nil,
			map[string]*adapter.BuilderInfo{
				"fooHandlerAdapter": {
					DefaultConfig:          &types.Empty{},
					ValidateConfig:         func(c proto.Message) error { return nil },
					CreateHandlerBuilderFn: func() config.HandlerBuilder { return nil },
				},
			},
			nil, 0, "service.name == “*”", false, ConstGlobalConfig,
		},
		{
			nil,
			nil,
			map[string]*adapter.BuilderInfo{ /*Empty lookup. Should cause error, Adapter not found*/ },
			nil, 1, "service.name == “*”", false, ConstGlobalConfig,
		},
		{
			nil,
			nil,
			map[string]*adapter.BuilderInfo{
				"fooHandlerAdapter": {
					DefaultConfig:          &types.Empty{},
					ValidateConfig:         func(c proto.Message) error { return nil },
					CreateHandlerBuilderFn: func() config.HandlerBuilder { return nil },
				},
			},
			nil, 1, "service.name == “*”", false, duplicateCnstrs,
		},
	}

	for idx, tt := range tests {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {

			var ce *adapter.ConfigErrors

			mgr := newVfinder(tt.ada, tt.asp, tt.hbi)
			p := newValidator(mgr.FindAspectValidator, mgr.FindAdapterValidator, mgr.FindBuilderInfo, fakeConfigureHandler, nil,
				mgr.AdapterToAspectMapperFunc, tt.strict, evaluator)
			ce = p.validateHandlers(tt.cfg)
			cok := ce == nil
			ok := tt.nerrors == 0

			if ok != cok {
				t.Fatalf("got %t, want %t", cok, ok)
			}
			if ce == nil {
				return
			}

			if len(ce.Multi.Errors) != tt.nerrors {
				t.Fatalf("got %s, want %s", ce.Error(), tt.cerr.Error())
			}
		})
	}
}

func TestCacheHandlerBuilder(t *testing.T) {
	evaluator := newFakeExpr()

	var globalConfig = `
subject: "namespace:ns"
revision: "2022"
handlers:
  - name: fooHandler
    adapter: fooHandlerAdapter
`

	const testSupportedTemplate adapter.SupportedTemplates = "testSupportedTemplate"
	tests := []*configTable{
		{
			hbi: map[string]*adapter.BuilderInfo{
				"fooHandlerAdapter": {
					DefaultConfig:          &types.Empty{},
					ValidateConfig:         func(c proto.Message) error { return nil },
					CreateHandlerBuilderFn: func() config.HandlerBuilder { return nil },
					SupportedTemplates:     []adapter.SupportedTemplates{testSupportedTemplate},
				},
			},
			cfg:     globalConfig,
			nerrors: 0,
		},
		{
			hbi:     map[string]*adapter.BuilderInfo{ /*Empty lookup. Should cause error, Adapter not found*/ },
			cfg:     globalConfig,
			nerrors: 1,
		},
	}

	for idx, tt := range tests {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			mgr := newVfinder(tt.ada, tt.asp, tt.hbi)
			p := newValidator(mgr.FindAspectValidator, mgr.FindAdapterValidator, mgr.FindBuilderInfo, fakeConfigureHandler, nil,
				mgr.AdapterToAspectMapperFunc, tt.strict, evaluator)
			_ = p.validateHandlers(tt.cfg)
			v, ok := p.handlerBuilderByName["fooHandler"]
			if tt.nerrors > 0 && ok {
				t.Fatalf("got p.handlerBuilderByName[\"fooHandler\"] = %v, want: <nil>", v)

			}
			if tt.nerrors == 0 {
				if !ok {
					t.Fatalf("got p.handlerBuilderByName[\"fooHandler\"] = %v, want: <nil>", v)

				} else if !reflect.DeepEqual(p.handlerBuilderByName["fooHandler"].supportedTemplates, []adapter.SupportedTemplates{testSupportedTemplate}) {
					t.Fatalf("got p.handlerBuilderByName[\"fooHandler\"]: %v, Expected: =%v",
						p.handlerBuilderByName["fooHandler"].supportedTemplates, []adapter.SupportedTemplates{testSupportedTemplate})
				}
			}
		})
	}
}

type fakeGoodHandlerBuilder struct{}

func (f fakeGoodHandlerBuilder) Build(cnfg proto.Message) (config.Handler, error) { return nil, nil }

type fakeBadHandlerBuilder struct{}

func (f fakeBadHandlerBuilder) Build(cnfg proto.Message) (config.Handler, error) {
	return nil, errors.New("build failed")
}

type fakePanicHandlerBuilder struct{}

func (f fakePanicHandlerBuilder) Build(cnfg proto.Message) (config.Handler, error) {
	panic("panic from handler Build method")
}

func getConfigureHandlerFn(err error) ConfigureHandlerFn {
	return func(actions []*pb.Action, constructors map[string]*pb.Constructor,
		handlers map[string]*HandlerBuilderInfo, tmplRepo tmpl.Repository, expr expr.TypeChecker, df expr.AttributeDescriptorFinder) error {
		return err
	}
}

func TestBuildHandlers(t *testing.T) {
	var hbgood config.HandlerBuilder = fakeGoodHandlerBuilder{}
	var hbb config.HandlerBuilder = fakeBadHandlerBuilder{}
	var hpb config.HandlerBuilder = fakePanicHandlerBuilder{}
	tests := []struct {
		configureHandler     ConfigureHandlerFn
		handlerBuilderByName map[string]*HandlerBuilderInfo
		expectedError        string
		buildPanic           bool
	}{
		{
			getConfigureHandlerFn(nil),
			map[string]*HandlerBuilderInfo{
				"foo": {
					handlerBuilder: &hbgood, handlerCnfg: &pb.Handler{Params: &types.Empty{}},
				},
			},
			"",
			false,
		},
		{
			getConfigureHandlerFn(errors.New("some error during configuration")),
			map[string]*HandlerBuilderInfo{
				"foo": {
					handlerBuilder: &hbgood, handlerCnfg: &pb.Handler{Params: &types.Empty{}},
				},
			},
			"some error during configuration",
			false,
		},
		{
			getConfigureHandlerFn(nil),
			map[string]*HandlerBuilderInfo{
				"foo": {
					handlerBuilder: &hbb, handlerCnfg: &pb.Handler{Params: &types.Empty{}},
				},
			},
			"failed to build a handler instance",
			false,
		},
		{
			getConfigureHandlerFn(nil),
			map[string]*HandlerBuilderInfo{
				"foo": {
					handlerBuilder: &hpb, handlerCnfg: &pb.Handler{Params: &types.Empty{}},
				},
			},
			"",
			true,
		},
	}

	for idx, tt := range tests {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			p := newValidator(nil, nil, nil, tt.configureHandler, nil, nil, true, nil)
			p.handlerBuilderByName = tt.handlerBuilderByName
			err := p.buildHandlers()

			if tt.buildPanic {
				if !tt.handlerBuilderByName["foo"].isBroken {
					t.Error("The handler should be marked as broken.")
				}
			} else {
				if tt.expectedError == "" {
					if len(p.validated.handlerByName) == len(tt.handlerBuilderByName) {
						for k := range tt.handlerBuilderByName {
							if _, ok := p.validated.handlerByName[k]; !ok {
								t.Fatalf("got validated.handlerByName[%s] = nil, want !nil", k)
							}
						}
					} else {
						t.Fatalf("got: validated.handlerByName %v, want: '%v'", p.validated.handlerByName, tt.handlerBuilderByName)
					}
				} else {
					if !strings.Contains(err.Error(), tt.expectedError) {
						t.Fatalf("got %s, want %s", tt.expectedError, err.Error())
					}
				}
			}

		})
	}
}

type fakeTemplateRepo struct{ templateConstructorParamMap map[string]proto.Message }

func newFakeTemplateRepo(templateConstructorParamMap map[string]proto.Message) tmpl.Repository {
	return fakeTemplateRepo{templateConstructorParamMap: templateConstructorParamMap}
}

func (t fakeTemplateRepo) GetTemplateInfo(template string) (tmpl.Info, bool) {
	if t.templateConstructorParamMap == nil {
		return tmpl.Info{}, false
	}
	if v, ok := t.templateConstructorParamMap[template]; ok {
		return tmpl.Info{
			CnstrDefConfig: v,
		}, true
	}
	return tmpl.Info{}, false
}

func TestValidateRulesConfig(t *testing.T) {
	const sSvcConfigValid = `
subject: namespace:ns
revision: "2022"
action_rules:
- selector: target.service == "*"
  actions:
  - handler: somehandler
    instances:
    - RequestCountByService
`
	const sSvcConfigNestedValid = `
subject: namespace:ns
revision: "2022"
action_rules:
- selector: target.service == "*"
  actions:
  - handler: somehandler
    instances:
    - RequestCountByService
  rules:
  - selector: target.service == "*"
    actions:
    - handler: somehandler
      instances:
      - RequestCountByService
`
	const sSvcConfigMissingHandler = `
subject: namespace:ns
revision: "2022"
action_rules:
- selector: target.service == "*"
  actions:
  - instances:
    - RequestCountByService
`
	const sSvcConfigNestedMissingHandler = `
subject: namespace:ns
revision: "2022"
action_rules:
- selector: target.ip == "*"
  actions:
  - handler: somehandler
    instances:
    - RequestCountByService
  rules:
  - selector: source.ip == "*"
    actions:
    - instances:
      - RequestCountByService
`
	const sSvcConfigInvalidSelector = `
subject: namespace:ns
revision: "2022"
action_rules:
- selector: == * == invalid
  actions:
  - handler: somehandler
    instances:
    - RequestCountByService
`

	const tmpl1 adapter.SupportedTemplates = "tmp1"
	evaluator := newFakeExpr()
	tests := []struct {
		cfg        string
		nerrors    int
		cnstrMap   map[string]*pb.Constructor
		handlerMap map[string]*HandlerBuilderInfo
		cerr       []string
		nActions   int
		expr       expr.TypeChecker
	}{
		{
			sSvcConfigValid,
			0,
			map[string]*pb.Constructor{"RequestCountByService": {Template: "tmp1"}},
			map[string]*HandlerBuilderInfo{"somehandler": {supportedTemplates: []adapter.SupportedTemplates{tmpl1}}},
			nil,
			1,
			evaluator,
		},
		{
			sSvcConfigNestedValid,
			0,
			map[string]*pb.Constructor{"RequestCountByService": {Template: "tmp1"}},
			map[string]*HandlerBuilderInfo{"somehandler": {supportedTemplates: []adapter.SupportedTemplates{tmpl1}}},
			nil,
			2,
			evaluator,
		},
		{
			sSvcConfigValid,
			2,
			map[string]*pb.Constructor{},
			map[string]*HandlerBuilderInfo{},
			[]string{"handler not specified or is invalid", "instance 'RequestCountByService' is not defined"},
			0,
			evaluator,
		},
		{
			sSvcConfigMissingHandler,
			1,
			map[string]*pb.Constructor{"RequestCountByService": {Template: "tmp1"}},
			map[string]*HandlerBuilderInfo{"somehandler": {}},
			[]string{"handler not specified or is invalid"},
			0,
			evaluator,
		},
		{
			sSvcConfigNestedMissingHandler,
			1,
			map[string]*pb.Constructor{"RequestCountByService": {Template: "tmp1"}},
			map[string]*HandlerBuilderInfo{"somehandler": {supportedTemplates: []adapter.SupportedTemplates{tmpl1}}},
			[]string{"handler not specified or is invalid"},
			1,
			evaluator,
		},
		{
			sSvcConfigInvalidSelector,
			1,
			map[string]*pb.Constructor{"RequestCountByService": {Template: "tmp1"}},
			map[string]*HandlerBuilderInfo{"somehandler": {supportedTemplates: []adapter.SupportedTemplates{tmpl1}}},
			[]string{"bad expression"},
			1, /*even if the selector is wrong the action is correct*/
			&fakeExpr{err: errors.New("bad expression")},
		},
		{
			sSvcConfigValid,
			1,
			map[string]*pb.Constructor{"RequestCountByService": {Template: "TemplateHandlerNotSupport"}},
			map[string]*HandlerBuilderInfo{"somehandler": {supportedTemplates: []adapter.SupportedTemplates{tmpl1}}},
			[]string{"handler does not support the template"},
			1,
			evaluator,
		},
	}

	tdf := newFakeTemplateRepo(map[string]proto.Message{"FooTemplate": &types.Empty{}})
	for idx, tt := range tests {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {

			var ce *adapter.ConfigErrors

			p := newValidator(nil, nil, nil, nil, tdf, nil, true, tt.expr)
			ce = p.validateRulesConfig(globalRulesKey, tt.cfg, tt.cnstrMap, tt.handlerMap)

			cok := ce == nil
			ok := tt.nerrors == 0

			if ok != cok {
				t.Errorf("got %t, want %t", cok, ok)
			}
			if len(p.actions) != tt.nActions {
				t.Errorf("got len(p.actions)=%d, want %d", len(p.actions), tt.nActions)
			}
			if ce == nil {
				return
			}

			if !containErrors(ce.Multi.Errors, tt.cerr) {
				t.Fatalf("got '%v' want %v", ce.Error(), tt.cerr)
			}
		})
	}
}

func TestValidateConstructorConfigs(t *testing.T) {
	const sSvcConfigInvalidTemplate = `
subject: namespace:ns
revision: "2022"
constructors:
- instanceName: RequestCountByService
  template: invalidtemplate
  params:
    value: 1
    dimensions:
      source: origin.service
      target_ip: destination.ip
`
	const sSvcConfigValidParams = `
subject: namespace:ns
revision: "2022"
constructors:
- instanceName: RequestCountByService
  template: FooTemplate
  params:
`
	const sSvcConfigExtraParamFields = `
subject: namespace:ns
revision: "2022"
constructors:
- instanceName: RequestCountByService
  template: FooTemplate
  params:
    check_expression: src.ip
    blacklist: true
`
	const sSvcConfigDuplicateCnstrs = `
subject: namespace:ns
revision: "2022"
constructors:
- instanceName: RequestCountByService
  template: FooTemplate
  params:
- instanceName: RequestCountByService
  template: FooTemplate
  params:
`

	evaluator := newFakeExpr()
	tests := []struct {
		cfg     string
		nerrors int
		tdf     tmpl.Repository
		cerr    []string
	}{
		{
			sSvcConfigInvalidTemplate,
			1,
			newFakeTemplateRepo(map[string]proto.Message{"FooTemplate": &types.Empty{}}),
			[]string{"is not a registered"},
		},
		{
			sSvcConfigValidParams,
			0,
			newFakeTemplateRepo(map[string]proto.Message{"FooTemplate": &types.Empty{}}),
			nil,
		},
		{
			sSvcConfigExtraParamFields,
			1,
			newFakeTemplateRepo(map[string]proto.Message{"FooTemplate": &types.Empty{}}),
			[]string{"failed to decode constructor params"},
		},
		{
			sSvcConfigDuplicateCnstrs,
			1,
			newFakeTemplateRepo(map[string]proto.Message{"FooTemplate": &types.Empty{}}),
			[]string{"duplicate constructors with same instanceNames "},
		},
	}

	for idx, tt := range tests {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {

			var ce *adapter.ConfigErrors

			p := newValidator(nil, nil, nil, nil, tt.tdf, nil, true, evaluator)
			ce = p.validateConstructorConfigs(globalRulesKey, tt.cfg)

			cok := ce == nil
			ok := tt.nerrors == 0

			if ok != cok {
				t.Errorf("got %t want %t ", cok, ok)
			}
			if ce == nil {
				return
			}

			if !containErrors(ce.Multi.Errors, tt.cerr) {
				t.Fatalf("got: %v, want: '%v' ", ce.Error(), tt.cerr)
			}
		})
	}
}

func TestFullConfigValidator(tt *testing.T) {
	fe := newFakeExpr()
	ctable := []struct {
		cerr     *adapter.ConfigError
		ada      map[string]adapter.ConfigValidator
		asp      map[Kind]AspectValidator
		selector string
		strict   bool
		cfg      string
		exprErr  error
	}{
		{nil,
			map[string]adapter.ConfigValidator{
				"denyChecker": &lc{},
				"metrics":     &lc{},
				"listchecker": &lc{},
			},
			map[Kind]AspectValidator{
				DenialsKind: &ac{},
				MetricsKind: &ac{},
				ListsKind:   &ac{},
			},
			"service.name == “*”", false, sSvcConfig2, nil},
		{nil,
			map[string]adapter.ConfigValidator{
				"denyChecker": &lc{},
				"metrics":     &lc{},
				"listchecker": &lc{},
			},
			map[Kind]AspectValidator{
				DenialsKind: &ac{},
				MetricsKind: &ac{},
				ListsKind:   &ac{},
			},
			"", false, sSvcConfig2, nil},
		{&adapter.ConfigError{Field: "namedAdapter", Underlying: errors.New("lists//denychecker.2 not available")},
			map[string]adapter.ConfigValidator{
				"denyChecker": &lc{},
				"metrics":     &lc{},
				"listchecker": &lc{},
			},
			map[Kind]AspectValidator{
				DenialsKind: &ac{},
				MetricsKind: &ac{},
				ListsKind:   &ac{},
			},
			"", false, sSvcConfig3, nil},
		{&adapter.ConfigError{Field: ":Selector service.name == “*”", Underlying: errors.New("invalid expression")},
			map[string]adapter.ConfigValidator{
				"denyChecker": &lc{},
				"metrics":     &lc{},
				"listchecker": &lc{},
			},
			map[Kind]AspectValidator{
				DenialsKind: &ac{},
				MetricsKind: &ac{},
				ListsKind:   &ac{},
			},
			"service.name == “*”", false, sSvcConfig1, errors.New("invalid expression")},
	}
	for idx, ctx := range ctable {
		tt.Run(fmt.Sprintf("[%d]", idx), func(t *testing.T) {
			mgr := newVfinder(ctx.ada, ctx.asp, nil)
			fe.err = ctx.exprErr
			p := newValidator(mgr.FindAspectValidator, mgr.FindAdapterValidator, nil, fakeConfigureHandler, nil,
				mgr.AdapterToAspectMapperFunc, ctx.strict, fe)
			// ConstGlobalConfig only defines 1 adapter: denyChecker
			_, ce := p.validate(newFakeMap(ConstGlobalConfig, ctx.cfg))
			cok := ce == nil
			ok := ctx.cerr == nil
			if ok != cok {
				t.Fatalf("%d got %t, want %t ", idx, cok, ok)

			}
			if ce == nil {
				return
			}
			if len(ce.Multi.Errors) < 2 {
				t.Fatal("expected at least 2 errors reported")
			}
			if !strings.Contains(ce.Multi.Errors[1].Error(), ctx.cerr.Error()) {
				t.Fatalf("%d got: %#v\nwant: %#v\n", idx, ce.Multi.Errors[1].Error(), ctx.cerr.Error())
			}
		})
	}
}

// globalRulesKey this rule set applies to all requests
// so we create a well known key for it
var globalRulesKey = rulesKey{Scope: global, Subject: global}

func TestConfigParseError(t *testing.T) {
	mgr := &fakeVFinder{}
	evaluator := newFakeExpr()
	p := newValidator(mgr.FindAspectValidator, mgr.FindAdapterValidator, nil, fakeConfigureHandler, nil,
		mgr.AdapterToAspectMapperFunc, false, evaluator)
	ce := p.validateServiceConfig(globalRulesKey, "<config>  </config>", false)

	if ce == nil || !strings.Contains(ce.Error(), "error unmarshaling") {
		t.Error("Expected unmarshal Error", ce)
	}
	ce = p.validateAdapters("", "<config>  </config>")

	if ce == nil || !strings.Contains(ce.Error(), "error unmarshaling") {
		t.Error("Expected unmarshal Error", ce)
	}

	ce = p.validateHandlers("<config>  </config>")

	if ce == nil || !strings.Contains(ce.Error(), "error unmarshaling") {
		t.Error("Expected unmarshal Error", ce)
	}

	ce = p.validateDescriptors("", "<config>  </config>")

	if ce == nil || !strings.Contains(ce.Error(), "error unmarshaling") {
		t.Error("Expected unmarshal Error", ce)
	}
	_, ce = p.validate(map[string]string{
		keyGlobalServiceConfig: "<config>  </config>",
		keyAdapters:            "<config>  </config>",
		keyDescriptors:         "<config>  </config>",
		keyHandlers:            "<config>  </config>",
	})
	if ce == nil || !strings.Contains(ce.Error(), "error unmarshaling") {
		t.Error("Expected unmarshal Error", ce)
	}
}

func TestDecoderError(t *testing.T) {
	err := decode(make(chan int), nil, true)
	if err == nil {
		t.Error("Expected json encode error")
	}
}

const ConstGlobalConfigValid = `
subject: "namespace:ns"
revision: "2022"
handlers:
  - name: fooHandler
    adapter: fooHandlerAdapter
adapters:
  - name: default
    kind: denials
    impl: denyChecker
    params:
      check_expression: src.ip
      blacklist: true
`
const ConstGlobalConfig = ConstGlobalConfigValid + `
      unknown_field: true
`

var duplicateCnstrs = `
subject: "namespace:ns"
revision: "2022"
handlers:
  - name: fooHandler
    adapter: fooHandlerAdapter
  - name: fooHandler
    adapter: fooHandlerAdapter
`

const sSvcConfig1 = `
subject: "namespace:ns"
revision: "2022"
rules:
- selector: service.name == “*”
  aspects:
  - kind: metrics
    params:
      metrics:
      - name: response_time_by_consumer
        value: metric_response_time
        metric_kind: DELTA
        labels:
        - key: target_consumer_id
`
const sSvcConfig2 = `
subject: namespace:ns
revision: "2022"
rules:
- selector: service.name == “*”
  aspects:
  - kind: lists
    inputs: {}
    params:
`
const sSvcConfig3 = `
subject: namespace:ns
revision: "2022"
rules:
- selector: service.name == “*”
  aspects:
  - kind: lists
    inputs: {}
    params:
    adapter: denychecker.2
`
const sSvcConfig = `
subject: namespace:ns
revision: "2022"
rules:
- selector: %s
  aspects:
  - kind: metrics
    adapter: ""
    inputs: {}
    params:
      check_expression: src.ip
      blacklist: true
      unknown_field: true
  rules:
  - selector: src.name == "abc"
    aspects:
    - kind: quotas
      adapter: ""
      inputs: {}
      params:
        check_expression: src.ip
        blacklist: true
`

type fakeExpr struct {
	err error
}

// newFakeExpr returns the basic
func newFakeExpr() *fakeExpr {
	return &fakeExpr{}
}

func UnboundVariable(vname string) error {
	return fmt.Errorf("unbound variable %s", vname)
}

// Eval evaluates given expression using the attribute bag
func (e *fakeExpr) Eval(mapExpression string, attrs attribute.Bag) (interface{}, error) {
	if v, found := attrs.Get(mapExpression); found {
		return v, nil
	}
	return nil, UnboundVariable(mapExpression)
}

// EvalString evaluates given expression using the attribute bag to a string
func (e *fakeExpr) EvalString(mapExpression string, attrs attribute.Bag) (string, error) {
	v, found := attrs.Get(mapExpression)
	if found {
		return v.(string), nil
	}
	return "", UnboundVariable(mapExpression)
}

// EvalPredicate evaluates given predicate using the attribute bag
func (e *fakeExpr) EvalPredicate(mapExpression string, attrs attribute.Bag) (bool, error) {
	v, found := attrs.Get(mapExpression)
	if found {
		return v.(bool), nil
	}
	return false, UnboundVariable(mapExpression)
}

func (e *fakeExpr) EvalType(string, expr.AttributeDescriptorFinder) (dpb.ValueType, error) {
	return dpb.VALUE_TYPE_UNSPECIFIED, e.err
}
func (e *fakeExpr) AssertType(string, expr.AttributeDescriptorFinder, dpb.ValueType) error {
	return e.err
}

func TestValidated_Clone(t *testing.T) {
	aa := map[adapterKey]*pb.Adapter{
		{AccessLogsKind, "n1"}: {},
	}

	hh := map[string]*HandlerInfo{
		"foo": nil,
	}

	rule := map[rulesKey]*pb.ServiceConfig{
		{"global", "global"}: {},
	}

	shas := map[string][sha1.Size]byte{
		keyGlobalServiceConfig: {},
	}

	adp := map[string]*pb.GlobalConfig{
		keyAdapters: {},
	}

	desc := map[string]*pb.GlobalConfig{
		keyDescriptors: {},
	}

	v := &Validated{
		adapterByName: aa,
		handlerByName: hh,
		rule:          rule,
		adapter:       adp,
		descriptor:    desc,
		numAspects:    1,
		shas:          shas,
	}

	v1 := v.Clone()

	if !reflect.DeepEqual(v, v1) {
		t.Errorf("got %v\nwant %v", v1, v)
	}
}

func TestParseConfigKey(t *testing.T) {
	for _, tst := range []struct {
		input string
		key   *rulesKey
	}{
		{keyGlobalServiceConfig, &rulesKey{"global", "global"}},
		{"/scopes/global/subjects/global", nil},
		{"/SCOPES/global/subjects/global/rules", nil},
		{"/scopes/global/SUBJECTS/global/rules", nil},
	} {
		t.Run(tst.input, func(t *testing.T) {
			k := parseRulesKey(tst.input)
			if !reflect.DeepEqual(k, tst.key) {
				t.Errorf("got %s\nwant %s", k, tst.key)
			}
		})
	}
}

func TestUnknownKind(t *testing.T) {
	name := "DOESNOTEXITS"
	want := unknownKind(name)
	_, ce := convertAspectParams(nil, name, "abc", true, nil)
	if ce == nil {
		t.Errorf("got nil\nwant %s", want)
	}
	if !strings.Contains(ce.Error(), want.Error()) {
		t.Errorf("got %s\nwant %s", ce, want)
	}
}

type fakeCV struct {
	err string
	cfg adapter.Config
}

func (f *fakeCV) DefaultConfig() (c adapter.Config) { return f.cfg }

// ValidateConfig determines whether the given configuration meets all correctness requirements.
func (f *fakeCV) ValidateConfig(c adapter.Config) (ce *adapter.ConfigErrors) {
	return ce.Appendf("abc", f.err)
}

func TestConvertAdapterParamsErrors(t *testing.T) {
	for _, tt := range []struct {
		params interface{}
		cv     *fakeCV
	}{{map[string]interface{}{
		"check_expression": "src.ip",
		"blacklist":        "true (this should be bool)",
	}, &fakeCV{cfg: &listcheckerpb.ListsParams{},
		err: "failed to decode"}},
		{map[string]interface{}{
			"check_expression": "src.ip",
			"blacklist":        true,
		}, &fakeCV{cfg: &listcheckerpb.ListsParams{}, err: "failed to unmarshal"}},
	} {
		t.Run(tt.cv.err, func(t *testing.T) {
			_, ce := convertAdapterParams(func(name string) (adapter.ConfigValidator, bool) {
				return tt.cv, true
			}, "ABC", tt.params, true)

			if !strings.Contains(ce.Error(), tt.cv.err) {
				t.Errorf("got %s, want %s. ", ce.Error(), tt.cv.err)
			}
		})
	}
}

func TestValidateDescriptors(t *testing.T) {
	tests := []struct {
		name string
		in   string // string repr of config
		err  string
	}{
		{"empty", "", ""},
		{"valid", `
metrics:
  - name: request_count
    kind: COUNTER
    value: INT64
    description: request count by source, target, service, and code
    labels:
      source: 1 # STRING
      target: 1 # STRING
      service: 1 # STRING
      method: 1 # STRING
      response_code: 2 # INT64
quotas:
  - name: RequestCount
    rate_limit: true
logs:
  - name: accesslog.common
    display_name: Apache Common Log Format
    payload_format: TEXT
    log_template: '{{.originIP}}'
    labels:
      originIp: 6 # IP_ADDRESS
      sourceUser: 1 # STRING
      timestamp: 5 # TIMESTAMP
      method: 1 # STRING
      url: 1 # STRING
      protocol: 1 # STRING
      responseCode: 2 # INT64
      responseSize: 2 # INT64
monitored_resources:
  - name: pod
    labels:
      podName: 1
principals:
  - name: Bob
    labels:
      originIp: 6`, ""},
		{"invalid name", `
quotas:
  - name:
    rate_limit: true`, "name"},
	}

	for idx, tt := range tests {
		t.Run(fmt.Sprintf("[%d] %s", idx, tt.name), func(t *testing.T) {
			v := &validator{validated: &Validated{descriptor: make(map[string]*pb.GlobalConfig)}}
			if err := v.validateDescriptors(tt.name, tt.in); err != nil || tt.err != "" {
				if tt.err == "" {
					t.Fatalf("validateDescriptors = '%s', wanted no err", err.Error())
				} else if !strings.Contains(err.Error(), tt.err) {
					t.Fatalf("got: '%s', want errors containing the string '%s', ", err.Error(), tt.err)
				}
			}
		})
	}
}

func containErrors(actualErrors []error, expectedErrorStrs []string) bool {
	if len(actualErrors) != len(expectedErrorStrs) {
		return false
	}

	for _, exp := range expectedErrorStrs {
		match := false
		for _, actualError := range actualErrors {
			if strings.Contains(actualError.Error(), exp) {
				match = true
				break
			}
		}
		if !match {
			return false
		}
	}
	return true
}
