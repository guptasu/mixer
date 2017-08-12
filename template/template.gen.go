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
	"fmt"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	rpc "github.com/googleapis/googleapis/google/rpc"
	"github.com/hashicorp/go-multierror"

	"istio.io/api/mixer/v1/config/descriptor"
	"istio.io/mixer/pkg/adapter"
	adptTmpl "istio.io/mixer/pkg/adapter/template"
	"istio.io/mixer/pkg/attribute"
	"istio.io/mixer/pkg/expr"
	"istio.io/mixer/pkg/status"
	"istio.io/mixer/pkg/template"

	"istio.io/mixer/template/checknothing"

	"istio.io/mixer/template/listentry"

	"istio.io/mixer/template/logentry"

	"istio.io/mixer/template/metric"

	"istio.io/mixer/template/quota"

	"istio.io/mixer/template/reportnothing"
)

var (
	SupportedTmplInfo = map[string]template.Info{

		checknothing.TemplateName: {
			CtrCfg:    &checknothing.InstanceParam{},
			Variety:   adptTmpl.TEMPLATE_VARIETY_CHECK,
			BldrName:  "istio.io/mixer/template/checknothing.CheckNothingHandlerBuilder",
			HndlrName: "istio.io/mixer/template/checknothing.CheckNothingHandler",
			SupportsTemplate: func(hndlrBuilder adapter.HandlerBuilder) bool {
				_, ok := hndlrBuilder.(checknothing.CheckNothingHandlerBuilder)
				return ok
			},
			HandlerSupportsTemplate: func(hndlr adapter.Handler) bool {
				_, ok := hndlr.(checknothing.CheckNothingHandler)
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
				castedBuilder := (*builder).(checknothing.CheckNothingHandlerBuilder)
				castedTypes := make(map[string]*checknothing.Type, len(types))
				for k, v := range types {
					// Mixer framework should have ensured the type safety.
					v1 := v.(*checknothing.Type)
					castedTypes[k] = v1
				}
				return castedBuilder.ConfigureCheckNothingHandler(castedTypes)
			},

			ProcessCheck: func(ctx context.Context, instName string, inst proto.Message, attrs attribute.Bag, mapper expr.Evaluator,
				handler adapter.Handler) (rpc.Status, adapter.CacheabilityInfo) {
				var found bool
				var err error

				castedInst := inst.(*checknothing.InstanceParam)
				var instances []*checknothing.Instance

				instance := &checknothing.Instance{
					Name: instName,
				}
				_ = castedInst

				var cacheInfo adapter.CacheabilityInfo
				if found, cacheInfo, err = handler.(checknothing.CheckNothingHandler).HandleCheckNothing(ctx, instance); err != nil {
					return status.WithError(err), adapter.CacheabilityInfo{}
				}

				if found {
					return status.OK, cacheInfo
				}

				return status.WithPermissionDenied(fmt.Sprintf("%s rejected", instances)), adapter.CacheabilityInfo{}
			},
			ProcessReport: nil,
			ProcessQuota:  nil,
		},

		listentry.TemplateName: {
			CtrCfg:    &listentry.InstanceParam{},
			Variety:   adptTmpl.TEMPLATE_VARIETY_CHECK,
			BldrName:  "istio.io/mixer/template/listentry.ListEntryHandlerBuilder",
			HndlrName: "istio.io/mixer/template/listentry.ListEntryHandler",
			SupportsTemplate: func(hndlrBuilder adapter.HandlerBuilder) bool {
				_, ok := hndlrBuilder.(listentry.ListEntryHandlerBuilder)
				return ok
			},
			HandlerSupportsTemplate: func(hndlr adapter.Handler) bool {
				_, ok := hndlr.(listentry.ListEntryHandler)
				return ok
			},
			InferType: func(cp proto.Message, tEvalFn template.TypeEvalFn) (proto.Message, error) {
				var err error = nil
				cpb := cp.(*listentry.InstanceParam)
				infrdType := &listentry.Type{}

				if cpb.Value == "" {
					return nil, fmt.Errorf("expression for field Value cannot be empty")
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
				castedBuilder := (*builder).(listentry.ListEntryHandlerBuilder)
				castedTypes := make(map[string]*listentry.Type, len(types))
				for k, v := range types {
					// Mixer framework should have ensured the type safety.
					v1 := v.(*listentry.Type)
					castedTypes[k] = v1
				}
				return castedBuilder.ConfigureListEntryHandler(castedTypes)
			},

			ProcessCheck: func(ctx context.Context, instName string, inst proto.Message, attrs attribute.Bag, mapper expr.Evaluator,
				handler adapter.Handler) (rpc.Status, adapter.CacheabilityInfo) {
				var found bool
				var err error

				castedInst := inst.(*listentry.InstanceParam)
				var instances []*listentry.Instance

				Value, err := mapper.Eval(castedInst.Value, attrs)

				if err != nil {
					return status.WithError(err), adapter.CacheabilityInfo{}
				}

				instance := &listentry.Instance{
					Name: instName,

					Value: Value.(string),
				}
				_ = castedInst

				var cacheInfo adapter.CacheabilityInfo
				if found, cacheInfo, err = handler.(listentry.ListEntryHandler).HandleListEntry(ctx, instance); err != nil {
					return status.WithError(err), adapter.CacheabilityInfo{}
				}

				if found {
					return status.OK, cacheInfo
				}

				return status.WithPermissionDenied(fmt.Sprintf("%s rejected", instances)), adapter.CacheabilityInfo{}
			},
			ProcessReport: nil,
			ProcessQuota:  nil,
		},

		logentry.TemplateName: {
			CtrCfg:    &logentry.InstanceParam{},
			Variety:   adptTmpl.TEMPLATE_VARIETY_REPORT,
			BldrName:  "istio.io/mixer/template/logentry.LogEntryHandlerBuilder",
			HndlrName: "istio.io/mixer/template/logentry.LogEntryHandler",
			SupportsTemplate: func(hndlrBuilder adapter.HandlerBuilder) bool {
				_, ok := hndlrBuilder.(logentry.LogEntryHandlerBuilder)
				return ok
			},
			HandlerSupportsTemplate: func(hndlr adapter.Handler) bool {
				_, ok := hndlr.(logentry.LogEntryHandler)
				return ok
			},
			InferType: func(cp proto.Message, tEvalFn template.TypeEvalFn) (proto.Message, error) {
				var err error = nil
				cpb := cp.(*logentry.InstanceParam)
				infrdType := &logentry.Type{}

				infrdType.Labels = make(map[string]istio_mixer_v1_config_descriptor.ValueType, len(cpb.Labels))
				for k, v := range cpb.Labels {
					if infrdType.Labels[k], err = tEvalFn(v); err != nil {
						return nil, err
					}
				}

				if cpb.Severity == "" {
					return nil, fmt.Errorf("expression for field Severity cannot be empty")
				}
				if t, e := tEvalFn(cpb.Severity); e != nil || t != istio_mixer_v1_config_descriptor.STRING {
					if e != nil {
						return nil, fmt.Errorf("failed to evaluate expression for field Severity: %v", e)
					}
					return nil, fmt.Errorf("error type checking for field Severity: Evaluated expression type %v want %v", t, istio_mixer_v1_config_descriptor.STRING)
				}

				_ = cpb
				return infrdType, err
			},
			ConfigureType: func(types map[string]proto.Message, builder *adapter.HandlerBuilder) error {
				// Mixer framework should have ensured the type safety.
				castedBuilder := (*builder).(logentry.LogEntryHandlerBuilder)
				castedTypes := make(map[string]*logentry.Type, len(types))
				for k, v := range types {
					// Mixer framework should have ensured the type safety.
					v1 := v.(*logentry.Type)
					castedTypes[k] = v1
				}
				return castedBuilder.ConfigureLogEntryHandler(castedTypes)
			},

			ProcessReport: func(ctx context.Context, insts map[string]proto.Message, attrs attribute.Bag, mapper expr.Evaluator, handler adapter.Handler) rpc.Status {
				result := &multierror.Error{}
				var instances []*logentry.Instance

				castedInsts := make(map[string]*logentry.InstanceParam, len(insts))
				for k, v := range insts {
					v1 := v.(*logentry.InstanceParam)
					castedInsts[k] = v1
				}
				for name, md := range castedInsts {

					Labels, err := template.EvalAll(md.Labels, attrs, mapper)

					if err != nil {
						result = multierror.Append(result, fmt.Errorf("failed to eval Labels for instance '%s': %v", name, err))
						continue
					}

					Severity, err := mapper.Eval(md.Severity, attrs)

					if err != nil {
						result = multierror.Append(result, fmt.Errorf("failed to eval Severity for instance '%s': %v", name, err))
						continue
					}

					instances = append(instances, &logentry.Instance{
						Name: name,

						Labels: Labels,

						Severity: Severity.(string),
					})
					_ = md
				}

				if err := handler.(logentry.LogEntryHandler).HandleLogEntry(ctx, instances); err != nil {
					result = multierror.Append(result, fmt.Errorf("failed to report all values: %v", err))
				}

				err := result.ErrorOrNil()
				if err != nil {
					return status.WithError(err)
				}

				return status.OK
			},
			ProcessCheck: nil,
			ProcessQuota: nil,
		},

		metric.TemplateName: {
			CtrCfg:    &metric.InstanceParam{},
			Variety:   adptTmpl.TEMPLATE_VARIETY_REPORT,
			BldrName:  "istio.io/mixer/template/metric.MetricHandlerBuilder",
			HndlrName: "istio.io/mixer/template/metric.MetricHandler",
			SupportsTemplate: func(hndlrBuilder adapter.HandlerBuilder) bool {
				_, ok := hndlrBuilder.(metric.MetricHandlerBuilder)
				return ok
			},
			HandlerSupportsTemplate: func(hndlr adapter.Handler) bool {
				_, ok := hndlr.(metric.MetricHandler)
				return ok
			},
			InferType: func(cp proto.Message, tEvalFn template.TypeEvalFn) (proto.Message, error) {
				var err error = nil
				cpb := cp.(*metric.InstanceParam)
				infrdType := &metric.Type{}

				if cpb.Value == "" {
					return nil, fmt.Errorf("expression for field Value cannot be empty")
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

				_ = cpb
				return infrdType, err
			},
			ConfigureType: func(types map[string]proto.Message, builder *adapter.HandlerBuilder) error {
				// Mixer framework should have ensured the type safety.
				castedBuilder := (*builder).(metric.MetricHandlerBuilder)
				castedTypes := make(map[string]*metric.Type, len(types))
				for k, v := range types {
					// Mixer framework should have ensured the type safety.
					v1 := v.(*metric.Type)
					castedTypes[k] = v1
				}
				return castedBuilder.ConfigureMetricHandler(castedTypes)
			},

			ProcessReport: func(ctx context.Context, insts map[string]proto.Message, attrs attribute.Bag, mapper expr.Evaluator, handler adapter.Handler) rpc.Status {
				result := &multierror.Error{}
				var instances []*metric.Instance

				castedInsts := make(map[string]*metric.InstanceParam, len(insts))
				for k, v := range insts {
					v1 := v.(*metric.InstanceParam)
					castedInsts[k] = v1
				}
				for name, md := range castedInsts {

					Value, err := mapper.Eval(md.Value, attrs)

					if err != nil {
						result = multierror.Append(result, fmt.Errorf("failed to eval Value for instance '%s': %v", name, err))
						continue
					}

					Dimensions, err := template.EvalAll(md.Dimensions, attrs, mapper)

					if err != nil {
						result = multierror.Append(result, fmt.Errorf("failed to eval Dimensions for instance '%s': %v", name, err))
						continue
					}

					instances = append(instances, &metric.Instance{
						Name: name,

						Value: Value,

						Dimensions: Dimensions,
					})
					_ = md
				}

				if err := handler.(metric.MetricHandler).HandleMetric(ctx, instances); err != nil {
					result = multierror.Append(result, fmt.Errorf("failed to report all values: %v", err))
				}

				err := result.ErrorOrNil()
				if err != nil {
					return status.WithError(err)
				}

				return status.OK
			},
			ProcessCheck: nil,
			ProcessQuota: nil,
		},

		quota.TemplateName: {
			CtrCfg:    &quota.InstanceParam{},
			Variety:   adptTmpl.TEMPLATE_VARIETY_QUOTA,
			BldrName:  "istio.io/mixer/template/quota.QuotaHandlerBuilder",
			HndlrName: "istio.io/mixer/template/quota.QuotaHandler",
			SupportsTemplate: func(hndlrBuilder adapter.HandlerBuilder) bool {
				_, ok := hndlrBuilder.(quota.QuotaHandlerBuilder)
				return ok
			},
			HandlerSupportsTemplate: func(hndlr adapter.Handler) bool {
				_, ok := hndlr.(quota.QuotaHandler)
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
				castedBuilder := (*builder).(quota.QuotaHandlerBuilder)
				castedTypes := make(map[string]*quota.Type, len(types))
				for k, v := range types {
					// Mixer framework should have ensured the type safety.
					v1 := v.(*quota.Type)
					castedTypes[k] = v1
				}
				return castedBuilder.ConfigureQuotaHandler(castedTypes)
			},

			ProcessQuota: func(ctx context.Context, quotaName string, inst proto.Message, attrs attribute.Bag, mapper expr.Evaluator, handler adapter.Handler,
				qma adapter.QuotaRequestArgs) (rpc.Status, adapter.CacheabilityInfo, adapter.QuotaResult) {
				castedInst := inst.(*quota.InstanceParam)

				Dimensions, err := template.EvalAll(castedInst.Dimensions, attrs, mapper)

				if err != nil {
					msg := fmt.Sprintf("failed to eval Dimensions for instance '%s': %v", quotaName, err)
					glog.Error(msg)
					return status.WithInvalidArgument(msg), adapter.CacheabilityInfo{}, adapter.QuotaResult{}
				}

				instance := &quota.Instance{
					Name: quotaName,

					Dimensions: Dimensions,
				}

				var qr adapter.QuotaResult
				var cacheInfo adapter.CacheabilityInfo
				if qr, cacheInfo, err = handler.(quota.QuotaHandler).HandleQuota(ctx, instance, qma); err != nil {
					glog.Errorf("Quota allocation failed: %v", err)
					return status.WithError(err), adapter.CacheabilityInfo{}, adapter.QuotaResult{}
				}
				if qr.Amount == 0 {
					msg := fmt.Sprintf("Unable to allocate %v units from quota %s", qma.QuotaAmount, quotaName)
					glog.Warning(msg)
					return status.WithResourceExhausted(msg), adapter.CacheabilityInfo{}, adapter.QuotaResult{}
				}
				if glog.V(2) {
					glog.Infof("Allocated %v units from quota %s", qma.QuotaAmount, quotaName)
				}
				return status.OK, cacheInfo, qr
			},
			ProcessReport: nil,
			ProcessCheck:  nil,
		},

		reportnothing.TemplateName: {
			CtrCfg:    &reportnothing.InstanceParam{},
			Variety:   adptTmpl.TEMPLATE_VARIETY_REPORT,
			BldrName:  "istio.io/mixer/template/reportnothing.ReportNothingHandlerBuilder",
			HndlrName: "istio.io/mixer/template/reportnothing.ReportNothingHandler",
			SupportsTemplate: func(hndlrBuilder adapter.HandlerBuilder) bool {
				_, ok := hndlrBuilder.(reportnothing.ReportNothingHandlerBuilder)
				return ok
			},
			HandlerSupportsTemplate: func(hndlr adapter.Handler) bool {
				_, ok := hndlr.(reportnothing.ReportNothingHandler)
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
				castedBuilder := (*builder).(reportnothing.ReportNothingHandlerBuilder)
				castedTypes := make(map[string]*reportnothing.Type, len(types))
				for k, v := range types {
					// Mixer framework should have ensured the type safety.
					v1 := v.(*reportnothing.Type)
					castedTypes[k] = v1
				}
				return castedBuilder.ConfigureReportNothingHandler(castedTypes)
			},

			ProcessReport: func(ctx context.Context, insts map[string]proto.Message, attrs attribute.Bag, mapper expr.Evaluator, handler adapter.Handler) rpc.Status {
				result := &multierror.Error{}
				var instances []*reportnothing.Instance

				castedInsts := make(map[string]*reportnothing.InstanceParam, len(insts))
				for k, v := range insts {
					v1 := v.(*reportnothing.InstanceParam)
					castedInsts[k] = v1
				}
				for name, md := range castedInsts {

					instances = append(instances, &reportnothing.Instance{
						Name: name,
					})
					_ = md
				}

				if err := handler.(reportnothing.ReportNothingHandler).HandleReportNothing(ctx, instances); err != nil {
					result = multierror.Append(result, fmt.Errorf("failed to report all values: %v", err))
				}

				err := result.ErrorOrNil()
				if err != nil {
					return status.WithError(err)
				}

				return status.OK
			},
			ProcessCheck: nil,
			ProcessQuota: nil,
		},
	}
)
