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

// THIS FILE IS AUTOMATICALLY GENERATED.

package sample

import (
	"context"
	"errors"
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/glog"

	"istio.io/api/mixer/v1/config/descriptor"
	"istio.io/mixer/pkg/adapter"
	adptTmpl "istio.io/mixer/pkg/adapter/template"
	"istio.io/mixer/pkg/attribute"
	"istio.io/mixer/pkg/expr"
	"istio.io/mixer/pkg/template"

	"istio.io/mixer/template/sample/check"

	"istio.io/mixer/template/sample/quota"

	"istio.io/mixer/template/sample/report"

	"time"
)

var (
	SupportedTmplInfo = map[string]template.Info{

		istio_mixer_adapter_sample_check.TemplateName: {
			Name:               istio_mixer_adapter_sample_check.TemplateName,
			Impl:               "istio.mixer.adapter.sample.check",
			CtrCfg:             &istio_mixer_adapter_sample_check.InstanceParam{},
			Variety:            adptTmpl.TEMPLATE_VARIETY_CHECK,
			BldrInterfaceName:  istio_mixer_adapter_sample_check.TemplateName + "." + "HandlerBuilder",
			HndlrInterfaceName: istio_mixer_adapter_sample_check.TemplateName + "." + "Handler",
			BuilderSupportsTemplate: func(hndlrBuilder adapter.Builder2) bool {
				_, ok := hndlrBuilder.(istio_mixer_adapter_sample_check.HandlerBuilder2)
				return ok
			},
			HandlerSupportsTemplate: func(hndlr adapter.Handler) bool {
				_, ok := hndlr.(istio_mixer_adapter_sample_check.Handler)
				return ok
			},
			InferType: func(cp proto.Message, tEvalFn template.TypeEvalFn) (proto.Message, error) {
				var err error = nil
				cpb := cp.(*istio_mixer_adapter_sample_check.InstanceParam)
				infrdType := &istio_mixer_adapter_sample_check.Type{}

				if cpb.CheckExpression == "" {
					return nil, errors.New("expression for field CheckExpression cannot be empty")
				}
				if t, e := tEvalFn(cpb.CheckExpression); e != nil || t != istio_mixer_v1_config_descriptor.STRING {
					if e != nil {
						return nil, fmt.Errorf("failed to evaluate expression for field CheckExpression: %v", e)
					}
					return nil, fmt.Errorf("error type checking for field CheckExpression: Evaluated expression type %v want %v", t, istio_mixer_v1_config_descriptor.STRING)
				}

				for _, v := range cpb.StringMap {
					if t, e := tEvalFn(v); e != nil || t != istio_mixer_v1_config_descriptor.STRING {
						if e != nil {
							return nil, fmt.Errorf("failed to evaluate expression for field StringMap: %v", e)
						}
						return nil, fmt.Errorf("error type checking for field StringMap: Evaluated expression type %v want %v", t, istio_mixer_v1_config_descriptor.STRING)
					}
				}

				_ = cpb
				return infrdType, err
			},
			SetType: func(types map[string]proto.Message, builder *adapter.Builder2) {
				// Mixer framework should have ensured the type safety.
				castedBuilder := (*builder).(istio_mixer_adapter_sample_check.HandlerBuilder2)
				castedTypes := make(map[string]*istio_mixer_adapter_sample_check.Type, len(types))
				for k, v := range types {
					// Mixer framework should have ensured the type safety.
					v1 := v.(*istio_mixer_adapter_sample_check.Type)
					castedTypes[k] = v1
				}
				castedBuilder.SetSampleTypes(castedTypes)
			},

			ProcessCheck: func(ctx context.Context, instName string, inst proto.Message, attrs attribute.Bag,
				mapper expr.Evaluator, handler adapter.Handler) (adapter.CheckResult, error) {
				castedInst := inst.(*istio_mixer_adapter_sample_check.InstanceParam)

				CheckExpression, err := mapper.Eval(castedInst.CheckExpression, attrs)

				if err != nil {
					msg := fmt.Sprintf("failed to eval CheckExpression for instance '%s': %v", instName, err)
					glog.Error(msg)
					return adapter.CheckResult{}, errors.New(msg)
				}

				StringMap, err := template.EvalAll(castedInst.StringMap, attrs, mapper)

				if err != nil {
					msg := fmt.Sprintf("failed to eval StringMap for instance '%s': %v", instName, err)
					glog.Error(msg)
					return adapter.CheckResult{}, errors.New(msg)
				}

				_ = castedInst

				instance := &istio_mixer_adapter_sample_check.Instance{
					Name: instName,

					CheckExpression: CheckExpression.(string),

					StringMap: func(m map[string]interface{}) map[string]string {
						res := make(map[string]string, len(m))
						for k, v := range m {
							res[k] = v.(string)
						}
						return res
					}(StringMap),
				}
				return handler.(istio_mixer_adapter_sample_check.Handler).HandleSample(ctx, instance)
			},
		},

		istio_mixer_adapter_sample_quota.TemplateName: {
			Name:               istio_mixer_adapter_sample_quota.TemplateName,
			Impl:               "istio.mixer.adapter.sample.quota",
			CtrCfg:             &istio_mixer_adapter_sample_quota.InstanceParam{},
			Variety:            adptTmpl.TEMPLATE_VARIETY_QUOTA,
			BldrInterfaceName:  istio_mixer_adapter_sample_quota.TemplateName + "." + "HandlerBuilder",
			HndlrInterfaceName: istio_mixer_adapter_sample_quota.TemplateName + "." + "Handler",
			BuilderSupportsTemplate: func(hndlrBuilder adapter.Builder2) bool {
				_, ok := hndlrBuilder.(istio_mixer_adapter_sample_quota.HandlerBuilder2)
				return ok
			},
			HandlerSupportsTemplate: func(hndlr adapter.Handler) bool {
				_, ok := hndlr.(istio_mixer_adapter_sample_quota.Handler)
				return ok
			},
			InferType: func(cp proto.Message, tEvalFn template.TypeEvalFn) (proto.Message, error) {
				var err error = nil
				cpb := cp.(*istio_mixer_adapter_sample_quota.InstanceParam)
				infrdType := &istio_mixer_adapter_sample_quota.Type{}

				infrdType.Dimensions = make(map[string]istio_mixer_v1_config_descriptor.ValueType, len(cpb.Dimensions))
				for k, v := range cpb.Dimensions {
					if infrdType.Dimensions[k], err = tEvalFn(v); err != nil {
						return nil, err
					}
				}

				for _, v := range cpb.BoolMap {
					if t, e := tEvalFn(v); e != nil || t != istio_mixer_v1_config_descriptor.BOOL {
						if e != nil {
							return nil, fmt.Errorf("failed to evaluate expression for field BoolMap: %v", e)
						}
						return nil, fmt.Errorf("error type checking for field BoolMap: Evaluated expression type %v want %v", t, istio_mixer_v1_config_descriptor.BOOL)
					}
				}

				_ = cpb
				return infrdType, err
			},
			SetType: func(types map[string]proto.Message, builder *adapter.Builder2) {
				// Mixer framework should have ensured the type safety.
				castedBuilder := (*builder).(istio_mixer_adapter_sample_quota.HandlerBuilder2)
				castedTypes := make(map[string]*istio_mixer_adapter_sample_quota.Type, len(types))
				for k, v := range types {
					// Mixer framework should have ensured the type safety.
					v1 := v.(*istio_mixer_adapter_sample_quota.Type)
					castedTypes[k] = v1
				}
				castedBuilder.SetQuotaTypes(castedTypes)
			},

			ProcessQuota: func(ctx context.Context, quotaName string, inst proto.Message, attrs attribute.Bag,
				mapper expr.Evaluator, handler adapter.Handler, args adapter.QuotaRequestArgs) (adapter.QuotaResult2, error) {
				castedInst := inst.(*istio_mixer_adapter_sample_quota.InstanceParam)

				Dimensions, err := template.EvalAll(castedInst.Dimensions, attrs, mapper)

				if err != nil {
					msg := fmt.Sprintf("failed to eval Dimensions for instance '%s': %v", quotaName, err)
					glog.Error(msg)
					return adapter.QuotaResult2{}, errors.New(msg)
				}

				BoolMap, err := template.EvalAll(castedInst.BoolMap, attrs, mapper)

				if err != nil {
					msg := fmt.Sprintf("failed to eval BoolMap for instance '%s': %v", quotaName, err)
					glog.Error(msg)
					return adapter.QuotaResult2{}, errors.New(msg)
				}

				instance := &istio_mixer_adapter_sample_quota.Instance{
					Name: quotaName,

					Dimensions: Dimensions,

					BoolMap: func(m map[string]interface{}) map[string]bool {
						res := make(map[string]bool, len(m))
						for k, v := range m {
							res[k] = v.(bool)
						}
						return res
					}(BoolMap),
				}

				return handler.(istio_mixer_adapter_sample_quota.Handler).HandleQuota(ctx, instance, args)
			},
		},

		istio_mixer_adapter_sample_report.TemplateName: {
			Name:               istio_mixer_adapter_sample_report.TemplateName,
			Impl:               "istio.mixer.adapter.sample.report",
			CtrCfg:             &istio_mixer_adapter_sample_report.InstanceParam{},
			Variety:            adptTmpl.TEMPLATE_VARIETY_REPORT,
			BldrInterfaceName:  istio_mixer_adapter_sample_report.TemplateName + "." + "HandlerBuilder",
			HndlrInterfaceName: istio_mixer_adapter_sample_report.TemplateName + "." + "Handler",
			BuilderSupportsTemplate: func(hndlrBuilder adapter.Builder2) bool {
				_, ok := hndlrBuilder.(istio_mixer_adapter_sample_report.HandlerBuilder2)
				return ok
			},
			HandlerSupportsTemplate: func(hndlr adapter.Handler) bool {
				_, ok := hndlr.(istio_mixer_adapter_sample_report.Handler)
				return ok
			},
			InferType: func(cp proto.Message, tEvalFn template.TypeEvalFn) (proto.Message, error) {
				var err error = nil
				cpb := cp.(*istio_mixer_adapter_sample_report.InstanceParam)
				infrdType := &istio_mixer_adapter_sample_report.Type{}

				if cpb.Value == "" {
					return nil, errors.New("expression for field Value cannot be empty")
				}
				if infrdType.Value, err = tEvalFn(cpb.Value); err != nil {
					return nil, err
				}

				infrdType.Dimensions = make(map[string]istio_mixer_v1_config_descriptor.ValueType, len(cpb.Dimensions))
				for k, v := range cpb.Dimensions {
					if infrdType.Dimensions[k], err = tEvalFn(v); err != nil {
						return nil, err
					}
				}

				if cpb.Int64Primitive == "" {
					return nil, errors.New("expression for field Int64Primitive cannot be empty")
				}
				if t, e := tEvalFn(cpb.Int64Primitive); e != nil || t != istio_mixer_v1_config_descriptor.INT64 {
					if e != nil {
						return nil, fmt.Errorf("failed to evaluate expression for field Int64Primitive: %v", e)
					}
					return nil, fmt.Errorf("error type checking for field Int64Primitive: Evaluated expression type %v want %v", t, istio_mixer_v1_config_descriptor.INT64)
				}

				if cpb.BoolPrimitive == "" {
					return nil, errors.New("expression for field BoolPrimitive cannot be empty")
				}
				if t, e := tEvalFn(cpb.BoolPrimitive); e != nil || t != istio_mixer_v1_config_descriptor.BOOL {
					if e != nil {
						return nil, fmt.Errorf("failed to evaluate expression for field BoolPrimitive: %v", e)
					}
					return nil, fmt.Errorf("error type checking for field BoolPrimitive: Evaluated expression type %v want %v", t, istio_mixer_v1_config_descriptor.BOOL)
				}

				if cpb.DoublePrimitive == "" {
					return nil, errors.New("expression for field DoublePrimitive cannot be empty")
				}
				if t, e := tEvalFn(cpb.DoublePrimitive); e != nil || t != istio_mixer_v1_config_descriptor.DOUBLE {
					if e != nil {
						return nil, fmt.Errorf("failed to evaluate expression for field DoublePrimitive: %v", e)
					}
					return nil, fmt.Errorf("error type checking for field DoublePrimitive: Evaluated expression type %v want %v", t, istio_mixer_v1_config_descriptor.DOUBLE)
				}

				if cpb.StringPrimitive == "" {
					return nil, errors.New("expression for field StringPrimitive cannot be empty")
				}
				if t, e := tEvalFn(cpb.StringPrimitive); e != nil || t != istio_mixer_v1_config_descriptor.STRING {
					if e != nil {
						return nil, fmt.Errorf("failed to evaluate expression for field StringPrimitive: %v", e)
					}
					return nil, fmt.Errorf("error type checking for field StringPrimitive: Evaluated expression type %v want %v", t, istio_mixer_v1_config_descriptor.STRING)
				}

				for _, v := range cpb.Int64Map {
					if t, e := tEvalFn(v); e != nil || t != istio_mixer_v1_config_descriptor.INT64 {
						if e != nil {
							return nil, fmt.Errorf("failed to evaluate expression for field Int64Map: %v", e)
						}
						return nil, fmt.Errorf("error type checking for field Int64Map: Evaluated expression type %v want %v", t, istio_mixer_v1_config_descriptor.INT64)
					}
				}

				if cpb.TimeStamp == "" {
					return nil, errors.New("expression for field TimeStamp cannot be empty")
				}
				if t, e := tEvalFn(cpb.TimeStamp); e != nil || t != istio_mixer_v1_config_descriptor.TIMESTAMP {
					if e != nil {
						return nil, fmt.Errorf("failed to evaluate expression for field TimeStamp: %v", e)
					}
					return nil, fmt.Errorf("error type checking for field TimeStamp: Evaluated expression type %v want %v", t, istio_mixer_v1_config_descriptor.TIMESTAMP)
				}

				if cpb.Duration == "" {
					return nil, errors.New("expression for field Duration cannot be empty")
				}
				if t, e := tEvalFn(cpb.Duration); e != nil || t != istio_mixer_v1_config_descriptor.DURATION {
					if e != nil {
						return nil, fmt.Errorf("failed to evaluate expression for field Duration: %v", e)
					}
					return nil, fmt.Errorf("error type checking for field Duration: Evaluated expression type %v want %v", t, istio_mixer_v1_config_descriptor.DURATION)
				}

				_ = cpb
				return infrdType, err
			},
			SetType: func(types map[string]proto.Message, builder *adapter.Builder2) {
				// Mixer framework should have ensured the type safety.
				castedBuilder := (*builder).(istio_mixer_adapter_sample_report.HandlerBuilder2)
				castedTypes := make(map[string]*istio_mixer_adapter_sample_report.Type, len(types))
				for k, v := range types {
					// Mixer framework should have ensured the type safety.
					v1 := v.(*istio_mixer_adapter_sample_report.Type)
					castedTypes[k] = v1
				}
				castedBuilder.SetSampleTypes(castedTypes)
			},

			ProcessReport: func(ctx context.Context, insts map[string]proto.Message, attrs attribute.Bag, mapper expr.Evaluator, handler adapter.Handler) error {
				var instances []*istio_mixer_adapter_sample_report.Instance
				for name, inst := range insts {
					md := inst.(*istio_mixer_adapter_sample_report.InstanceParam)

					Value, err := mapper.Eval(md.Value, attrs)

					if err != nil {
						msg := fmt.Sprintf("failed to eval Value for instance '%s': %v", name, err)
						glog.Error(msg)
						return errors.New(msg)
					}

					Dimensions, err := template.EvalAll(md.Dimensions, attrs, mapper)

					if err != nil {
						msg := fmt.Sprintf("failed to eval Dimensions for instance '%s': %v", name, err)
						glog.Error(msg)
						return errors.New(msg)
					}

					Int64Primitive, err := mapper.Eval(md.Int64Primitive, attrs)

					if err != nil {
						msg := fmt.Sprintf("failed to eval Int64Primitive for instance '%s': %v", name, err)
						glog.Error(msg)
						return errors.New(msg)
					}

					BoolPrimitive, err := mapper.Eval(md.BoolPrimitive, attrs)

					if err != nil {
						msg := fmt.Sprintf("failed to eval BoolPrimitive for instance '%s': %v", name, err)
						glog.Error(msg)
						return errors.New(msg)
					}

					DoublePrimitive, err := mapper.Eval(md.DoublePrimitive, attrs)

					if err != nil {
						msg := fmt.Sprintf("failed to eval DoublePrimitive for instance '%s': %v", name, err)
						glog.Error(msg)
						return errors.New(msg)
					}

					StringPrimitive, err := mapper.Eval(md.StringPrimitive, attrs)

					if err != nil {
						msg := fmt.Sprintf("failed to eval StringPrimitive for instance '%s': %v", name, err)
						glog.Error(msg)
						return errors.New(msg)
					}

					Int64Map, err := template.EvalAll(md.Int64Map, attrs, mapper)

					if err != nil {
						msg := fmt.Sprintf("failed to eval Int64Map for instance '%s': %v", name, err)
						glog.Error(msg)
						return errors.New(msg)
					}

					TimeStamp, err := mapper.Eval(md.TimeStamp, attrs)

					if err != nil {
						msg := fmt.Sprintf("failed to eval TimeStamp for instance '%s': %v", name, err)
						glog.Error(msg)
						return errors.New(msg)
					}

					Duration, err := mapper.Eval(md.Duration, attrs)

					if err != nil {
						msg := fmt.Sprintf("failed to eval Duration for instance '%s': %v", name, err)
						glog.Error(msg)
						return errors.New(msg)
					}

					instances = append(instances, &istio_mixer_adapter_sample_report.Instance{
						Name: name,

						Value: Value,

						Dimensions: Dimensions,

						Int64Primitive: Int64Primitive.(int64),

						BoolPrimitive: BoolPrimitive.(bool),

						DoublePrimitive: DoublePrimitive.(float64),

						StringPrimitive: StringPrimitive.(string),

						Int64Map: func(m map[string]interface{}) map[string]int64 {
							res := make(map[string]int64, len(m))
							for k, v := range m {
								res[k] = v.(int64)
							}
							return res
						}(Int64Map),

						TimeStamp: TimeStamp.(time.Time),

						Duration: Duration.(time.Duration),
					})
					_ = md
				}

				if err := handler.(istio_mixer_adapter_sample_report.Handler).HandleSample(ctx, instances); err != nil {
					return fmt.Errorf("failed to report all values: %v", err)
				}
				return nil
			},
		},
	}
)
