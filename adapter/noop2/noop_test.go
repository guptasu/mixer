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

// NOTE: This test will eventually be auto-generated so that it automatically supports all templates
//       known to Mixer. For now, it's manually curated.

import (
	"reflect"
	"testing"
	"time"

	rpc "github.com/googleapis/googleapis/google/rpc"

	"istio.io/mixer/pkg/adapter"
	"istio.io/mixer/template/checknothing"
	"istio.io/mixer/template/listentry"
	"istio.io/mixer/template/logentry"
	"istio.io/mixer/template/metric"
	"istio.io/mixer/template/quota"
	"istio.io/mixer/template/reportnothing"
)

func TestBasic(t *testing.T) {
	info := GetBuilderInfo()

	if !contains(info.SupportedTemplates, checknothing.TemplateName) ||
		!contains(info.SupportedTemplates, reportnothing.TemplateName) ||
		!contains(info.SupportedTemplates, listentry.TemplateName) ||
		!contains(info.SupportedTemplates, logentry.TemplateName) ||
		!contains(info.SupportedTemplates, metric.TemplateName) ||
		!contains(info.SupportedTemplates, quota.TemplateName) {
		t.Error("Didn't find all expected supported templates")
	}

	builder := info.CreateHandlerBuilder()
	cfg := info.DefaultConfig

	if err := info.ValidateConfig(cfg); err != nil {
		t.Errorf("Got error %v, expecting success", err)
	}

	checkNothingBuilder := builder.(checknothing.HandlerBuilder)
	if err := checkNothingBuilder.ConfigureCheckNothingHandler(nil); err != nil {
		t.Errorf("Got error %v, expecting success", err)
	}

	reportNothingBuilder := builder.(reportnothing.HandlerBuilder)
	if err := reportNothingBuilder.ConfigureReportNothingHandler(nil); err != nil {
		t.Errorf("Got error %v, expecting success", err)
	}

	listEntryBuilder := builder.(listentry.HandlerBuilder)
	if err := listEntryBuilder.ConfigureListEntryHandler(nil); err != nil {
		t.Errorf("Got error %v, expecting success", err)
	}

	logEntryBuilder := builder.(logentry.HandlerBuilder)
	if err := logEntryBuilder.ConfigureLogEntryHandler(nil); err != nil {
		t.Errorf("Got error %v, expecting success", err)
	}

	metricBuilder := builder.(metric.HandlerBuilder)
	if err := metricBuilder.ConfigureMetricHandler(nil); err != nil {
		t.Errorf("Got error %v, expecting success", err)
	}

	quotaBuilder := builder.(quota.HandlerBuilder)
	if err := quotaBuilder.ConfigureQuotaHandler(nil); err != nil {
		t.Errorf("Got error %v, expecting success", err)
	}

	handler, err := builder.Build(nil, nil)
	if err != nil {
		t.Errorf("Got error %v, expecting success", err)
	}

	var exptCheckRes = adapter.CheckResult{
		Status:        rpc.Status{Code: int32(rpc.OK)},
		ValidDuration: 1000000000 * time.Second,
		ValidUseCount: 1000000000,
	}
	var exptReportResult = adapter.ReportResult{
		Status: rpc.Status{Code: int32(rpc.OK)},
	}

	checkNothingHandler := handler.(checknothing.Handler)
	if result, err := checkNothingHandler.HandleCheckNothing(nil, nil); err != nil {
		t.Errorf("Got error %v, expecting success", err)
	} else if !reflect.DeepEqual(result, exptCheckRes) {
		t.Errorf("Got %v, expecting %v result", result, exptCheckRes)
	}

	reportNothingHandler := handler.(reportnothing.Handler)
	if result, err := reportNothingHandler.HandleReportNothing(nil, nil); err != nil {
		t.Errorf("Got error %v, expecting success", err)
	} else if !reflect.DeepEqual(result, exptReportResult) {
		t.Errorf("Got %v, expecting %v result", result, exptReportResult)
	}

	listEntryHandler := handler.(listentry.Handler)
	if result, err := listEntryHandler.HandleListEntry(nil, nil); err != nil {
		t.Errorf("Got error %v, expecting success", err)
	} else if !reflect.DeepEqual(result, exptCheckRes) {
		t.Errorf("Got %v, expecting %v result", result, exptCheckRes)
	}

	logEntryHandler := handler.(logentry.Handler)
	if result, err := logEntryHandler.HandleLogEntry(nil, nil); err != nil {
		t.Errorf("Got error %v, expecting success", err)
	} else if !reflect.DeepEqual(result, exptReportResult) {
		t.Errorf("Got %v, expecting %v result", result, exptReportResult)
	}

	metricHandler := handler.(metric.Handler)
	if result, err := metricHandler.HandleMetric(nil, nil); err != nil {
		t.Errorf("Got error %v, expecting success", err)
	} else if !reflect.DeepEqual(result, exptReportResult) {
		t.Errorf("Got %v, expecting %v result", result, exptReportResult)
	}

	if err = handler.Close(); err != nil {
		t.Errorf("Got error %v, expecting success", err)
	}

	exptQuotaResult := adapter.QuotaResult2{
		Status:        rpc.Status{Code: int32(rpc.OK)},
		ValidDuration: 1000000000 * time.Second,
		Amount:        100,
	}

	quotaHandler := handler.(quota.Handler)
	if result, err := quotaHandler.HandleQuota(nil, nil, adapter.QuotaRequestArgs{QuotaAmount: 100}); err != nil {
		t.Errorf("Got error %v, expecting success", err)
	} else if !reflect.DeepEqual(result, exptQuotaResult) {
		t.Errorf("Got %v, expecting %v result", result, exptQuotaResult)
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
