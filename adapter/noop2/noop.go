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

package noop2

// NOTE: This adapter will eventually be auto-generated so that it automatically supports all templates
//       known to Mixer. For now, it's manually curated.

import (
	"context"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/golang/protobuf/proto"

	"istio.io/mixer/pkg/adapter"
	"istio.io/mixer/template/checknothing"
	"istio.io/mixer/template/listentry"
	"istio.io/mixer/template/logentry"
	"istio.io/mixer/template/metric"
	"istio.io/mixer/template/quota"
	"istio.io/mixer/template/reportnothing"
	rpc "github.com/googleapis/googleapis/google/rpc"
)

type (
	handler struct{}
	builder struct{}
)

var _ checknothing.CheckNothingHandlerBuilder = builder{}
var _ checknothing.CheckNothingHandler = handler{}
var _ reportnothing.ReportNothingHandlerBuilder = builder{}
var _ reportnothing.ReportNothingHandler = handler{}
var _ listentry.ListEntryHandlerBuilder = builder{}
var _ listentry.ListEntryHandler = handler{}
var _ logentry.LogEntryHandlerBuilder = builder{}
var _ logentry.LogEntryHandler = handler{}
var _ metric.MetricHandlerBuilder = builder{}
var _ metric.MetricHandler = handler{}
var _ quota.QuotaHandlerBuilder = builder{}
var _ quota.QuotaHandler = handler{}

///////////////// Configuration Methods ///////////////

func (builder) Build(proto.Message, adapter.Env) (adapter.Handler, error) {
	return handler{}, nil
}

func (builder) ConfigureCheckNothingHandler(map[string]*checknothing.Type) error {
	return nil
}

func (builder) ConfigureReportNothingHandler(map[string]*reportnothing.Type) error {
	return nil
}

func (builder) ConfigureListEntryHandler(map[string]*listentry.Type) error {
	return nil
}

func (builder) ConfigureLogEntryHandler(map[string]*logentry.Type) error {
	return nil
}

func (builder) ConfigureMetricHandler(map[string]*metric.Type) error {
	return nil
}

func (builder) ConfigureQuotaHandler(map[string]*quota.Type) error {
	return nil
}

////////////////// Runtime Methods //////////////////////////

var checkResult = adapter.CheckResult {
	ValidDuration: 1000000000 * time.Second,
	ValidUseCount: 1000000000,
}

func (handler) HandleCheckNothing(context.Context, *checknothing.Instance) (adapter.CheckResult, error) {
	return checkResult, nil
}

func (handler) HandleReportNothing(context.Context, []*reportnothing.Instance) error {
	return nil
}

func (handler) HandleListEntry(context.Context, *listentry.Instance) (adapter.CheckResult, error) {
	return checkResult, nil
}

func (handler) HandleLogEntry(context.Context, []*logentry.Instance) error {
	return nil
}

func (handler) HandleMetric(context.Context, []*metric.Instance) error {
	return nil
}

func (handler) HandleQuota(ctx context.Context, _ *quota.Instance, args adapter.QuotaRequestArgs) (adapter.QuotaResult2, error) {
	return adapter.QuotaResult2 {
		Status: rpc.Status{Code: int32(rpc.OK)},
		ValidDuration: 1000000000 * time.Second,
		Amount:     args.QuotaAmount,
	},
		nil
}

func (handler) Close() error { return nil }

////////////////// Bootstrap //////////////////////////

// GetBuilderInfo returns the BuilderInfo associated with this adapter implementation.
func GetBuilderInfo() adapter.BuilderInfo {
	return adapter.BuilderInfo{
		Name:        "istio.io/mixer/adapters/noop",
		Description: "An adapter that does nothing",
		SupportedTemplates: []string{
			checknothing.TemplateName,
			reportnothing.TemplateName,
			listentry.TemplateName,
			logentry.TemplateName,
			metric.TemplateName,
			quota.TemplateName,
		},
		CreateHandlerBuilderFn: func() adapter.HandlerBuilder { return builder{} },
		DefaultConfig:          &types.Empty{},
		ValidateConfig:         func(msg proto.Message) error { return nil },
	}
}
