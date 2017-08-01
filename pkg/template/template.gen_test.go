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

package template

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"

	pb "istio.io/api/mixer/v1/config/descriptor"
	adptConfig "istio.io/mixer/pkg/adapter/config"
	adpTmpl "istio.io/mixer/pkg/adapter/template"
	sample_check "istio.io/mixer/template/sample/check"
	sample_quota "istio.io/mixer/template/sample/quota"
	sample_report "istio.io/mixer/template/sample/report"
)

func TestGeneratedFields(t *testing.T) {
	for _, tst := range []struct {
		tmplName  string
		ctrCfg    proto.Message
		variety   adpTmpl.TemplateVariety
		bldrName  string
		hndlrName string
	}{
		{
			tmplName:  sample_report.TemplateName,
			ctrCfg:    &sample_report.ConstructorParam{},
			variety:   adpTmpl.TEMPLATE_VARIETY_REPORT,
			bldrName:  "istio.io/mixer/template/sample/report.SampleProcessorBuilder",
			hndlrName: "istio.io/mixer/template/sample/report.SampleProcessor",
		},
		{
			tmplName:  sample_check.TemplateName,
			ctrCfg:    &sample_check.ConstructorParam{},
			variety:   adpTmpl.TEMPLATE_VARIETY_CHECK,
			bldrName:  "istio.io/mixer/template/sample/check.SampleProcessorBuilder",
			hndlrName: "istio.io/mixer/template/sample/check.SampleProcessor",
		},
		{
			tmplName:  sample_quota.TemplateName,
			ctrCfg:    &sample_quota.ConstructorParam{},
			variety:   adpTmpl.TEMPLATE_VARIETY_QUOTA,
			bldrName:  "istio.io/mixer/template/sample/quota.QuotaProcessorBuilder",
			hndlrName: "istio.io/mixer/template/sample/quota.QuotaProcessor",
		},
	} {
		t.Run(tst.tmplName, func(t *testing.T) {
			if !reflect.DeepEqual(SupportedTmplInfo[tst.tmplName].CtrCfg, tst.ctrCfg) {
				t.Errorf("SupportedTmplInfo[%s].CtrCfg = %T, want %T", tst.tmplName, SupportedTmplInfo[tst.tmplName].CtrCfg, tst.ctrCfg)
			}
			if SupportedTmplInfo[tst.tmplName].Variety != tst.variety {
				t.Errorf("SupportedTmplInfo[%s].Variety = %v, want %v", tst.tmplName, SupportedTmplInfo[tst.tmplName].Variety, tst.variety)
			}
			if SupportedTmplInfo[tst.tmplName].BldrName != tst.bldrName {
				t.Errorf("SupportedTmplInfo[%s].BldrName = %v, want %v", tst.tmplName, SupportedTmplInfo[tst.tmplName].BldrName, tst.bldrName)
			}
			if SupportedTmplInfo[tst.tmplName].HndlrName != tst.hndlrName {
				t.Errorf("SupportedTmplInfo[%s].HndlrName = %v, want %v", tst.tmplName, SupportedTmplInfo[tst.tmplName].HndlrName, tst.hndlrName)
			}
		})
	}
}

type badHandler struct{}

func (h badHandler) Close() error                                         { return nil }
func (h badHandler) Build(cnfg proto.Message) (adptConfig.Handler, error) { return nil, nil }

type reportHandler struct{}

func (h reportHandler) Close() error                                           { return nil }
func (h reportHandler) ReportSample(instances []*sample_report.Instance) error { return nil }
func (h reportHandler) Build(cnfg proto.Message) (adptConfig.Handler, error)   { return nil, nil }
func (h reportHandler) ConfigureSample(map[string]*sample_report.Type) error   { return nil }

type checkHandler struct{}

func (h checkHandler) Close() error { return nil }
func (h checkHandler) CheckSample(instance []*sample_check.Instance) (bool, adptConfig.CacheabilityInfo, error) {
	return true, adptConfig.CacheabilityInfo{}, nil
}
func (h checkHandler) Build(cnfg proto.Message) (adptConfig.Handler, error) { return nil, nil }
func (h checkHandler) ConfigureSample(map[string]*sample_check.Type) error  { return nil }

type quotaHandler struct{}

func (h quotaHandler) Close() error                                           { return nil }
func (h quotaHandler) ReportSample(instances []*sample_report.Instance) error { return nil }
func (h quotaHandler) Build(cnfg proto.Message) (adptConfig.Handler, error)   { return nil, nil }
func (h quotaHandler) ConfigureQuota(map[string]*sample_quota.Type) error     { return nil }

func TestHandlerSupportsTemplate(t *testing.T) {
	for _, tst := range []struct {
		tmplName string
		hndlr    adptConfig.Handler
		result   bool
	}{
		{
			tmplName: sample_report.TemplateName,
			hndlr:    badHandler{},
			result:   false,
		},
		{
			tmplName: sample_report.TemplateName,
			hndlr:    reportHandler{},
			result:   true,
		},
		{
			tmplName: sample_check.TemplateName,
			hndlr:    badHandler{},
			result:   false,
		},
		{
			tmplName: sample_check.TemplateName,
			hndlr:    checkHandler{},
			result:   true,
		},
		{
			tmplName: sample_quota.TemplateName,
			hndlr:    badHandler{},
			result:   false,
		},
		{
			tmplName: sample_quota.TemplateName,
			hndlr:    quotaHandler{},
			result:   true,
		},
	} {
		t.Run(tst.tmplName, func(t *testing.T) {
			c := SupportedTmplInfo[tst.tmplName].HandlerSupportsTemplate(tst.hndlr)
			if c != tst.result {
				t.Errorf("SupportedTmplInfo[%s].HandlerSupportsTemplate(%T) = %t, want %t", tst.tmplName, tst.hndlr, c, tst.result)
			}
		})
	}
}

func TestBuilderSupportsTemplate(t *testing.T) {
	for _, tst := range []struct {
		tmplName  string
		hndlrBldr adptConfig.HandlerBuilder
		result    bool
	}{
		{
			tmplName:  sample_report.TemplateName,
			hndlrBldr: badHandler{},
			result:    false,
		},
		{
			tmplName:  sample_report.TemplateName,
			hndlrBldr: reportHandler{},
			result:    true,
		},
		{
			tmplName:  sample_check.TemplateName,
			hndlrBldr: badHandler{},
			result:    false,
		},
		{
			tmplName:  sample_check.TemplateName,
			hndlrBldr: checkHandler{},
			result:    true,
		},
		{
			tmplName:  sample_quota.TemplateName,
			hndlrBldr: badHandler{},
			result:    false,
		},
		{
			tmplName:  sample_quota.TemplateName,
			hndlrBldr: quotaHandler{},
			result:    true,
		},
	} {
		t.Run(tst.tmplName, func(t *testing.T) {
			c := SupportedTmplInfo[tst.tmplName].SupportsTemplate(tst.hndlrBldr)
			if c != tst.result {
				t.Errorf("SupportedTmplInfo[%s].SupportsTemplate(%T) = %t, want %t", tst.tmplName, tst.hndlrBldr, c, tst.result)
			}
		})
	}
}

func TestInferTypeForSampleReport(t *testing.T) {
	for _, tst := range []inferTypeTest{
		{
			name: "SimpleValid",
			cnstrCnfg: `
value: response.size
dimensions:
  source: source.ip
  target: source.ip
`,
			cnstrParamPb:           &sample_report.ConstructorParam{},
			typeEvalRet:            pb.INT64,
			typeEvalError:          nil,
			expectedValueType:      pb.INT64,
			expectedDimensionsType: map[string]pb.ValueType{"source": pb.INT64, "target": pb.INT64},
			expectedErr:            "",
			willPanic:              false,
		},
		{
			name:         "NotValidConstructorParam",
			cnstrCnfg:    ``,
			cnstrParamPb: &empty.Empty{}, // cnstr type mismatch
			expectedErr:  "is not of type",
			willPanic:    true,
		},
		{
			name: "ErrorFromTypeEvaluator",
			cnstrCnfg: `
value: response.size
dimensions:
  source: source.ip
`,
			cnstrParamPb:  &sample_report.ConstructorParam{},
			typeEvalError: fmt.Errorf("some expression x.y.z is invalid"),
			expectedErr:   "some expression x.y.z is invalid",
		},
	} {
		t.Run(tst.name, func(t *testing.T) {
			cp := tst.cnstrParamPb
			_ = fillProto(tst.cnstrCnfg, cp)
			typeEvalFn := func(expr string) (pb.ValueType, error) { return tst.typeEvalRet, tst.typeEvalError }
			defer func() {
				r := recover()
				if tst.willPanic && r == nil {
					t.Errorf("Expected to recover from panic for %s, but recover was nil.", tst.name)
				} else if !tst.willPanic && r != nil {
					t.Errorf("got panic %v, expected success.", r)
				}
			}()
			cv, cerr := SupportedTmplInfo[sample_report.TemplateName].InferType(cp.(proto.Message), typeEvalFn)
			if tst.expectedErr == "" {
				if cerr != nil {
					t.Errorf("got err %v\nwant <nil>", cerr)
				}
				if tst.expectedValueType != cv.(*sample_report.Type).Value {
					t.Errorf("got inferTypeForSampleReport(\n%s\n).value=%v\nwant %v", tst.cnstrCnfg, cv.(*sample_report.Type).Value, tst.expectedValueType)
				}
				if len(tst.expectedDimensionsType) != len(cv.(*sample_report.Type).Dimensions) {
					t.Errorf("got len ( inferTypeForSampleReport(\n%s\n).dimensions) =%v \n want %v", tst.cnstrCnfg, len(cv.(*sample_report.Type).Dimensions), len(tst.expectedDimensionsType))
				}
				for a, b := range tst.expectedDimensionsType {
					if cv.(*sample_report.Type).Dimensions[a] != b {
						t.Errorf("got inferTypeForSampleReport(\n%s\n).dimensions[%s] =%v \n want %v", tst.cnstrCnfg, a, cv.(*sample_report.Type).Dimensions[a], b)
					}
				}

			} else {
				if cerr == nil || !strings.Contains(cerr.Error(), tst.expectedErr) {
					t.Errorf("got error %v\nwant %v", cerr, tst.expectedErr)
				}
			}
		})
	}
}

type inferTypeTest struct {
	name                   string
	cnstrCnfg              string
	cnstrParamPb           interface{}
	typeEvalRet            pb.ValueType
	typeEvalError          error
	expectedValueType      pb.ValueType
	expectedDimensionsType map[string]pb.ValueType
	expectedErr            string
	willPanic              bool
}

func TestInferTypeForSampleCheck(t *testing.T) {
	for _, tst := range []inferTypeTest{
		{
			name: "SimpleValid",
			cnstrCnfg: `
check_expression: response.size
`,
			cnstrParamPb:      &sample_check.ConstructorParam{},
			typeEvalRet:       pb.STRING,
			typeEvalError:     nil,
			expectedValueType: pb.STRING,
			expectedErr:       "",
			willPanic:         false,
		},
		{
			name:         "NotValidConstructorParam",
			cnstrCnfg:    ``,
			cnstrParamPb: &empty.Empty{}, // cnstr type mismatch
			willPanic:    true,
		},
	} {
		t.Run(tst.name, func(t *testing.T) {
			cp := tst.cnstrParamPb
			_ = fillProto(tst.cnstrCnfg, cp)
			typeEvalFn := func(expr string) (pb.ValueType, error) { return tst.typeEvalRet, tst.typeEvalError }
			defer func() {
				r := recover()
				if tst.willPanic && r == nil {
					t.Errorf("Expected to recover from panic for %s, but recover was nil.", tst.name)
				} else if !tst.willPanic && r != nil {
					t.Errorf("got panic %v, expected success.", r)
				}
			}()
			cv, cerr := SupportedTmplInfo[sample_check.TemplateName].InferType(cp.(proto.Message), typeEvalFn)
			if tst.willPanic {
				t.Error("Should not reach this statement due to panic.")
			}
			if tst.expectedErr == "" {
				if cerr != nil {
					t.Errorf("got err %v\nwant <nil>", cerr)
				}
				if tst.expectedValueType != cv.(*sample_check.Type).CheckExpression {
					t.Errorf("got inferTypeForSampleCheck(\n%s\n).value=%v\nwant %v", tst.cnstrCnfg, cv.(*sample_check.Type).CheckExpression, tst.expectedValueType)
				}
			} else {
				if cerr == nil || !strings.Contains(cerr.Error(), tst.expectedErr) {
					t.Errorf("got error %v\nwant %v", cerr, tst.expectedErr)
				}
			}
		})
	}
}

func TestInferTypeForSampleQuota(t *testing.T) {
	for _, tst := range []inferTypeTest{
		{
			name: "SimpleValid",
			cnstrCnfg: `
dimensions:
  source: source.ip
  target: source.ip
  env: target.ip
`,
			cnstrParamPb:           &sample_quota.ConstructorParam{},
			typeEvalRet:            pb.STRING,
			typeEvalError:          nil,
			expectedValueType:      pb.STRING,
			expectedDimensionsType: map[string]pb.ValueType{"source": pb.STRING, "target": pb.STRING, "env": pb.STRING},
			expectedErr:            "",
			willPanic:              false,
		},
		{
			name:         "NotValidConstructorParam",
			cnstrCnfg:    ``,
			cnstrParamPb: &empty.Empty{}, // cnstr type mismatch
			expectedErr:  "is not of type",
			willPanic:    true,
		},
		{
			name: "ErrorFromTypeEvaluator",
			cnstrCnfg: `
dimensions:
  source: source.ip
`,
			cnstrParamPb:  &sample_quota.ConstructorParam{},
			typeEvalError: fmt.Errorf("some expression x.y.z is invalid"),
			expectedErr:   "some expression x.y.z is invalid",
		},
	} {
		t.Run(tst.name, func(t *testing.T) {
			cp := tst.cnstrParamPb
			_ = fillProto(tst.cnstrCnfg, cp)
			typeEvalFn := func(expr string) (pb.ValueType, error) { return tst.typeEvalRet, tst.typeEvalError }
			defer func() {
				r := recover()
				if tst.willPanic && r == nil {
					t.Errorf("Expected to recover from panic for %s, but recover was nil.", tst.name)
				} else if !tst.willPanic && r != nil {
					t.Errorf("got panic %v, expected success.", r)
				}
			}()
			cv, cerr := SupportedTmplInfo[sample_quota.TemplateName].InferType(cp.(proto.Message), typeEvalFn)
			if tst.expectedErr == "" {
				if cerr != nil {
					t.Errorf("got err %v\nwant <nil>", cerr)
				}
				if len(tst.expectedDimensionsType) != len(cv.(*sample_quota.Type).Dimensions) {
					t.Errorf("got len ( inferTypeForSampleReport(\n%s\n).dimensions) =%v \n want %v", tst.cnstrCnfg, len(cv.(*sample_quota.Type).Dimensions), len(tst.expectedDimensionsType))
				}
				for a, b := range tst.expectedDimensionsType {
					if cv.(*sample_quota.Type).Dimensions[a] != b {
						t.Errorf("got inferTypeForSampleReport(\n%s\n).dimensions[%s] =%v \n want %v", tst.cnstrCnfg, a, cv.(*sample_quota.Type).Dimensions[a], b)
					}
				}

			} else {
				if cerr == nil || !strings.Contains(cerr.Error(), tst.expectedErr) {
					t.Errorf("got error %v\nwant %v", cerr, tst.expectedErr)
				}
			}
		})
	}
}

func fillProto(cfg string, o interface{}) error {
	//data []byte, m map[string]interface{}, err error
	var m map[string]interface{}
	var data []byte
	var err error

	if err = yaml.Unmarshal([]byte(cfg), &m); err != nil {
		return err
	}

	if data, err = json.Marshal(m); err != nil {
		return err
	}

	err = yaml.Unmarshal(data, o)
	return err
}
