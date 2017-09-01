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

package template

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

	"istio.io/mixer/template/checknothing"

	"istio.io/mixer/template/listentry"

	"istio.io/mixer/template/logentry"

	"istio.io/mixer/template/metric"

	"istio.io/mixer/template/quota"

	"istio.io/mixer/template/reportnothing"

	"time"
)

var (
	SupportedTmplInfo = map[string]template.Info{

		checknothing.TemplateName: {
			Name:               checknothing.TemplateName,
			Impl:               "checknothing",
			CtrCfg:             &checknothing.InstanceParam{},
			Variety:            adptTmpl.TEMPLATE_VARIETY_CHECK,
			BldrInterfaceName:  checknothing.TemplateName + "." + "HandlerBuilder",
			HndlrInterfaceName: checknothing.TemplateName + "." + "Handler",
			BuilderSupportsTemplate: func(hndlrBuilder adapter.HandlerBuilder) bool {
				_, ok := hndlrBuilder.(checknothing.HandlerBuilder)
				return ok
			},
			HandlerSupportsTemplate: func(hndlr adapter.Handler) bool {
				_, ok := hndlr.(checknothing.Handler)
				return ok
			},
			InferType: func(cp proto.Message, tEvalFn template.TypeEvalFn) (proto.Message, error) {
				var err error = nil
				cpb := cp.(*checknothing.InstanceParam)
				infrdType := &checknothing.Type{}

				_ = cpb
				return infrdType, err
			},
			ConfigureType: func(types map[string]proto.Message, builder *adapter.HandlerBuilder) error {
				// Mixer framework should have ensured the type safety.
				castedBuilder := (*builder).(checknothing.HandlerBuilder)
				castedTypes := make(map[string]*checknothing.Type, len(types))
				for k, v := range types {
					// Mixer framework should have ensured the type safety.
					v1 := v.(*checknothing.Type)
					castedTypes[k] = v1
				}
				return castedBuilder.ConfigureCheckNothingHandler(castedTypes)
			},

			ProcessCheck: func(ctx context.Context, instName string, inst proto.Message, attrs attribute.Bag,
				mapper expr.Evaluator, handler adapter.Handler) (adapter.CheckResult, error) {
				castedInst := inst.(*checknothing.InstanceParam)

				_ = castedInst

				instance := &checknothing.Instance{
					Name: instName,
				}
				return handler.(checknothing.Handler).HandleCheckNothing(ctx, instance)
			},
		},

		listentry.TemplateName: {
			Name:               listentry.TemplateName,
			Impl:               "listentry",
			CtrCfg:             &listentry.InstanceParam{},
			Variety:            adptTmpl.TEMPLATE_VARIETY_CHECK,
			BldrInterfaceName:  listentry.TemplateName + "." + "HandlerBuilder",
			HndlrInterfaceName: listentry.TemplateName + "." + "Handler",
			BuilderSupportsTemplate: func(hndlrBuilder adapter.HandlerBuilder) bool {
				_, ok := hndlrBuilder.(listentry.HandlerBuilder)
				return ok
			},
			HandlerSupportsTemplate: func(hndlr adapter.Handler) bool {
				_, ok := hndlr.(listentry.Handler)
				return ok
			},
			InferType: func(cp proto.Message, tEvalFn template.TypeEvalFn) (proto.Message, error) {
				var err error = nil
				cpb := cp.(*listentry.InstanceParam)
				infrdType := &listentry.Type{}

				if cpb.Value == "" {
					return nil, errors.New("expression for field Value cannot be empty")
				}
				if t, e := tEvalFn(cpb.Value); e != nil || t != istio_mixer_v1_config_descriptor.STRING {
					if e != nil {
						return nil, fmt.Errorf("failed to evaluate expression for field Value: %v", e)
					}
					return nil, fmt.Errorf("error type checking for field Value: Evaluated expression type %v want %v", t, istio_mixer_v1_config_descriptor.STRING)
				}

				_ = cpb
				return infrdType, err
			},
			ConfigureType: func(types map[string]proto.Message, builder *adapter.HandlerBuilder) error {
				// Mixer framework should have ensured the type safety.
				castedBuilder := (*builder).(listentry.HandlerBuilder)
				castedTypes := make(map[string]*listentry.Type, len(types))
				for k, v := range types {
					// Mixer framework should have ensured the type safety.
					v1 := v.(*listentry.Type)
					castedTypes[k] = v1
				}
				return castedBuilder.ConfigureListEntryHandler(castedTypes)
			},

			ProcessCheck: func(ctx context.Context, instName string, inst proto.Message, attrs attribute.Bag,
				mapper expr.Evaluator, handler adapter.Handler) (adapter.CheckResult, error) {
				castedInst := inst.(*listentry.InstanceParam)

				Value, err := mapper.Eval(castedInst.Value, attrs)

				if err != nil {
					msg := fmt.Sprintf("failed to eval Value for instance '%s': %v", instName, err)
					glog.Error(msg)
					return adapter.CheckResult{}, errors.New(msg)
				}

				_ = castedInst

				instance := &listentry.Instance{
					Name: instName,

					Value: Value.(string),
				}
				return handler.(listentry.Handler).HandleListEntry(ctx, instance)
			},
		},

		logentry.TemplateName: {
			Name:               logentry.TemplateName,
			Impl:               "logentry",
			CtrCfg:             &logentry.InstanceParam{},
			Variety:            adptTmpl.TEMPLATE_VARIETY_REPORT,
			BldrInterfaceName:  logentry.TemplateName + "." + "HandlerBuilder",
			HndlrInterfaceName: logentry.TemplateName + "." + "Handler",
			BuilderSupportsTemplate: func(hndlrBuilder adapter.HandlerBuilder) bool {
				_, ok := hndlrBuilder.(logentry.HandlerBuilder)
				return ok
			},
			HandlerSupportsTemplate: func(hndlr adapter.Handler) bool {
				_, ok := hndlr.(logentry.Handler)
				return ok
			},
			InferType: func(cp proto.Message, tEvalFn template.TypeEvalFn) (proto.Message, error) {
				var err error = nil
				cpb := cp.(*logentry.InstanceParam)
				infrdType := &logentry.Type{}

				infrdType.Variables = make(map[string]istio_mixer_v1_config_descriptor.ValueType, len(cpb.Variables))
				for k, v := range cpb.Variables {
					if infrdType.Variables[k], err = tEvalFn(v); err != nil {
						return nil, err
					}
				}

				if cpb.Timestamp == "" {
					return nil, errors.New("expression for field Timestamp cannot be empty")
				}
				if t, e := tEvalFn(cpb.Timestamp); e != nil || t != istio_mixer_v1_config_descriptor.TIMESTAMP {
					if e != nil {
						return nil, fmt.Errorf("failed to evaluate expression for field Timestamp: %v", e)
					}
					return nil, fmt.Errorf("error type checking for field Timestamp: Evaluated expression type %v want %v", t, istio_mixer_v1_config_descriptor.TIMESTAMP)
				}

				if cpb.Severity == "" {
					return nil, errors.New("expression for field Severity cannot be empty")
				}
				if t, e := tEvalFn(cpb.Severity); e != nil || t != istio_mixer_v1_config_descriptor.STRING {
					if e != nil {
						return nil, fmt.Errorf("failed to evaluate expression for field Severity: %v", e)
					}
					return nil, fmt.Errorf("error type checking for field Severity: Evaluated expression type %v want %v", t, istio_mixer_v1_config_descriptor.STRING)
				}

				if cpb.MonitoredResourceType == "" {
					return nil, errors.New("expression for field MonitoredResourceType cannot be empty")
				}
				if t, e := tEvalFn(cpb.MonitoredResourceType); e != nil || t != istio_mixer_v1_config_descriptor.STRING {
					if e != nil {
						return nil, fmt.Errorf("failed to evaluate expression for field MonitoredResourceType: %v", e)
					}
					return nil, fmt.Errorf("error type checking for field MonitoredResourceType: Evaluated expression type %v want %v", t, istio_mixer_v1_config_descriptor.STRING)
				}

				infrdType.MonitoredResourceDimensions = make(map[string]istio_mixer_v1_config_descriptor.ValueType, len(cpb.MonitoredResourceDimensions))
				for k, v := range cpb.MonitoredResourceDimensions {
					if infrdType.MonitoredResourceDimensions[k], err = tEvalFn(v); err != nil {
						return nil, err
					}
				}

				_ = cpb
				return infrdType, err
			},
			ConfigureType: func(types map[string]proto.Message, builder *adapter.HandlerBuilder) error {
				// Mixer framework should have ensured the type safety.
				castedBuilder := (*builder).(logentry.HandlerBuilder)
				castedTypes := make(map[string]*logentry.Type, len(types))
				for k, v := range types {
					// Mixer framework should have ensured the type safety.
					v1 := v.(*logentry.Type)
					castedTypes[k] = v1
				}
				return castedBuilder.ConfigureLogEntryHandler(castedTypes)
			},

			ProcessReport: func(ctx context.Context, insts map[string]proto.Message, attrs attribute.Bag, mapper expr.Evaluator, handler adapter.Handler) error {
				var instances []*logentry.Instance
				for name, inst := range insts {
					md := inst.(*logentry.InstanceParam)

					Variables, err := template.EvalAll(md.Variables, attrs, mapper)

					if err != nil {
						msg := fmt.Sprintf("failed to eval Variables for instance '%s': %v", name, err)
						glog.Error(msg)
						return errors.New(msg)
					}

					Timestamp, err := mapper.Eval(md.Timestamp, attrs)

					if err != nil {
						msg := fmt.Sprintf("failed to eval Timestamp for instance '%s': %v", name, err)
						glog.Error(msg)
						return errors.New(msg)
					}

					Severity, err := mapper.Eval(md.Severity, attrs)

					if err != nil {
						msg := fmt.Sprintf("failed to eval Severity for instance '%s': %v", name, err)
						glog.Error(msg)
						return errors.New(msg)
					}

					MonitoredResourceType, err := mapper.Eval(md.MonitoredResourceType, attrs)

					if err != nil {
						msg := fmt.Sprintf("failed to eval MonitoredResourceType for instance '%s': %v", name, err)
						glog.Error(msg)
						return errors.New(msg)
					}

					MonitoredResourceDimensions, err := template.EvalAll(md.MonitoredResourceDimensions, attrs, mapper)

					if err != nil {
						msg := fmt.Sprintf("failed to eval MonitoredResourceDimensions for instance '%s': %v", name, err)
						glog.Error(msg)
						return errors.New(msg)
					}

					instances = append(instances, &logentry.Instance{
						Name: name,

						Variables: Variables,

						Timestamp: Timestamp.(time.Time),

						Severity: Severity.(string),

						MonitoredResourceType: MonitoredResourceType.(string),

						MonitoredResourceDimensions: MonitoredResourceDimensions,
					})
					_ = md
				}

				if err := handler.(logentry.Handler).HandleLogEntry(ctx, instances); err != nil {
					return fmt.Errorf("failed to report all values: %v", err)
				}
				return nil
			},
		},

		metric.TemplateName: {
			Name:               metric.TemplateName,
			Impl:               "metric",
			CtrCfg:             &metric.InstanceParam{},
			Variety:            adptTmpl.TEMPLATE_VARIETY_REPORT,
			BldrInterfaceName:  metric.TemplateName + "." + "HandlerBuilder",
			HndlrInterfaceName: metric.TemplateName + "." + "Handler",
			BuilderSupportsTemplate: func(hndlrBuilder adapter.HandlerBuilder) bool {
				_, ok := hndlrBuilder.(metric.HandlerBuilder)
				return ok
			},
			HandlerSupportsTemplate: func(hndlr adapter.Handler) bool {
				_, ok := hndlr.(metric.Handler)
				return ok
			},
			InferType: func(cp proto.Message, tEvalFn template.TypeEvalFn) (proto.Message, error) {
				var err error = nil
				cpb := cp.(*metric.InstanceParam)
				infrdType := &metric.Type{}

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

				if cpb.MonitoredResourceType == "" {
					return nil, errors.New("expression for field MonitoredResourceType cannot be empty")
				}
				if t, e := tEvalFn(cpb.MonitoredResourceType); e != nil || t != istio_mixer_v1_config_descriptor.STRING {
					if e != nil {
						return nil, fmt.Errorf("failed to evaluate expression for field MonitoredResourceType: %v", e)
					}
					return nil, fmt.Errorf("error type checking for field MonitoredResourceType: Evaluated expression type %v want %v", t, istio_mixer_v1_config_descriptor.STRING)
				}

				infrdType.MonitoredResourceDimensions = make(map[string]istio_mixer_v1_config_descriptor.ValueType, len(cpb.MonitoredResourceDimensions))
				for k, v := range cpb.MonitoredResourceDimensions {
					if infrdType.MonitoredResourceDimensions[k], err = tEvalFn(v); err != nil {
						return nil, err
					}
				}

				_ = cpb
				return infrdType, err
			},
			ConfigureType: func(types map[string]proto.Message, builder *adapter.HandlerBuilder) error {
				// Mixer framework should have ensured the type safety.
				castedBuilder := (*builder).(metric.HandlerBuilder)
				castedTypes := make(map[string]*metric.Type, len(types))
				for k, v := range types {
					// Mixer framework should have ensured the type safety.
					v1 := v.(*metric.Type)
					castedTypes[k] = v1
				}
				return castedBuilder.ConfigureMetricHandler(castedTypes)
			},

			ProcessReport: func(ctx context.Context, insts map[string]proto.Message, attrs attribute.Bag, mapper expr.Evaluator, handler adapter.Handler) error {
				var instances []*metric.Instance
				for name, inst := range insts {
					md := inst.(*metric.InstanceParam)

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

					MonitoredResourceType, err := mapper.Eval(md.MonitoredResourceType, attrs)

					if err != nil {
						msg := fmt.Sprintf("failed to eval MonitoredResourceType for instance '%s': %v", name, err)
						glog.Error(msg)
						return errors.New(msg)
					}

					MonitoredResourceDimensions, err := template.EvalAll(md.MonitoredResourceDimensions, attrs, mapper)

					if err != nil {
						msg := fmt.Sprintf("failed to eval MonitoredResourceDimensions for instance '%s': %v", name, err)
						glog.Error(msg)
						return errors.New(msg)
					}

					instances = append(instances, &metric.Instance{
						Name: name,

						Value: Value,

						Dimensions: Dimensions,

						MonitoredResourceType: MonitoredResourceType.(string),

						MonitoredResourceDimensions: MonitoredResourceDimensions,
					})
					_ = md
				}

				if err := handler.(metric.Handler).HandleMetric(ctx, instances); err != nil {
					return fmt.Errorf("failed to report all values: %v", err)
				}
				return nil
			},
		},

		quota.TemplateName: {
			Name:               quota.TemplateName,
			Impl:               "quota",
			CtrCfg:             &quota.InstanceParam{},
			Variety:            adptTmpl.TEMPLATE_VARIETY_QUOTA,
			BldrInterfaceName:  quota.TemplateName + "." + "HandlerBuilder",
			HndlrInterfaceName: quota.TemplateName + "." + "Handler",
			BuilderSupportsTemplate: func(hndlrBuilder adapter.HandlerBuilder) bool {
				_, ok := hndlrBuilder.(quota.HandlerBuilder)
				return ok
			},
			HandlerSupportsTemplate: func(hndlr adapter.Handler) bool {
				_, ok := hndlr.(quota.Handler)
				return ok
			},
			InferType: func(cp proto.Message, tEvalFn template.TypeEvalFn) (proto.Message, error) {
				var err error = nil
				cpb := cp.(*quota.InstanceParam)
				infrdType := &quota.Type{}

				infrdType.Dimensions = make(map[string]istio_mixer_v1_config_descriptor.ValueType, len(cpb.Dimensions))
				for k, v := range cpb.Dimensions {
					if infrdType.Dimensions[k], err = tEvalFn(v); err != nil {
						return nil, err
					}
				}

				_ = cpb
				return infrdType, err
			},
			ConfigureType: func(types map[string]proto.Message, builder *adapter.HandlerBuilder) error {
				// Mixer framework should have ensured the type safety.
				castedBuilder := (*builder).(quota.HandlerBuilder)
				castedTypes := make(map[string]*quota.Type, len(types))
				for k, v := range types {
					// Mixer framework should have ensured the type safety.
					v1 := v.(*quota.Type)
					castedTypes[k] = v1
				}
				return castedBuilder.ConfigureQuotaHandler(castedTypes)
			},

			ProcessQuota: func(ctx context.Context, quotaName string, inst proto.Message, attrs attribute.Bag,
				mapper expr.Evaluator, handler adapter.Handler, args adapter.QuotaRequestArgs) (adapter.QuotaResult2, error) {
				castedInst := inst.(*quota.InstanceParam)

				Dimensions, err := template.EvalAll(castedInst.Dimensions, attrs, mapper)

				if err != nil {
					msg := fmt.Sprintf("failed to eval Dimensions for instance '%s': %v", quotaName, err)
					glog.Error(msg)
					return adapter.QuotaResult2{}, errors.New(msg)
				}

				instance := &quota.Instance{
					Name: quotaName,

					Dimensions: Dimensions,
				}

				return handler.(quota.Handler).HandleQuota(ctx, instance, args)
			},
		},

		reportnothing.TemplateName: {
			Name:               reportnothing.TemplateName,
			Impl:               "reportnothing",
			CtrCfg:             &reportnothing.InstanceParam{},
			Variety:            adptTmpl.TEMPLATE_VARIETY_REPORT,
			BldrInterfaceName:  reportnothing.TemplateName + "." + "HandlerBuilder",
			HndlrInterfaceName: reportnothing.TemplateName + "." + "Handler",
			BuilderSupportsTemplate: func(hndlrBuilder adapter.HandlerBuilder) bool {
				_, ok := hndlrBuilder.(reportnothing.HandlerBuilder)
				return ok
			},
			HandlerSupportsTemplate: func(hndlr adapter.Handler) bool {
				_, ok := hndlr.(reportnothing.Handler)
				return ok
			},
			InferType: func(cp proto.Message, tEvalFn template.TypeEvalFn) (proto.Message, error) {
				var err error = nil
				cpb := cp.(*reportnothing.InstanceParam)
				infrdType := &reportnothing.Type{}

				_ = cpb
				return infrdType, err
			},
			ConfigureType: func(types map[string]proto.Message, builder *adapter.HandlerBuilder) error {
				// Mixer framework should have ensured the type safety.
				castedBuilder := (*builder).(reportnothing.HandlerBuilder)
				castedTypes := make(map[string]*reportnothing.Type, len(types))
				for k, v := range types {
					// Mixer framework should have ensured the type safety.
					v1 := v.(*reportnothing.Type)
					castedTypes[k] = v1
				}
				return castedBuilder.ConfigureReportNothingHandler(castedTypes)
			},

			ProcessReport: func(ctx context.Context, insts map[string]proto.Message, attrs attribute.Bag, mapper expr.Evaluator, handler adapter.Handler) error {
				var instances []*reportnothing.Instance
				for name, inst := range insts {
					md := inst.(*reportnothing.InstanceParam)

					instances = append(instances, &reportnothing.Instance{
						Name: name,
					})
					_ = md
				}

				if err := handler.(reportnothing.Handler).HandleReportNothing(ctx, instances); err != nil {
					return fmt.Errorf("failed to report all values: %v", err)
				}
				return nil
			},
		},
	}
)
