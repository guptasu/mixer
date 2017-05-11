package main

import (
	"istio.io/mixer/pkg/adapter"
)

/////////// COMMON GENERTED STUB CODE. NOT WRITTEN BY USER///////////
type RequestCount struct {
	Value        int64
	Service      string
	Method       string
	ResponseCode int64
	Source       string
	Target       string
}

type ReportResult struct {
	result [][]interface{}
}

func CreateReportResult() *ReportResult {
	result := make([][]interface{}, 0)

	return &ReportResult{result: result}
}

func WrapRequestCountToAdapterReqObject(val RequestCount) *adapter.Value{
	a := adapter.Value{}
	a.MetricValue = val.Value
	a.Labels = make(map[string]interface{})
	a.Labels["method"] = val.Method
	a.Labels["response_code"] = val.ResponseCode
	a.Labels["service"] = val.Service
	a.Labels["source"] = val.Source
	a.Labels["target"] = val.Target
	return &a
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics0(val RequestCount) {
	// convert flattened RequestCount into structure to be passed to adapters
	a := WrapRequestCountToAdapterReqObject(val)

	innerValue := []interface{}{"aspectName0", map[string]interface{}{"descriptorName": "request_count", "value": a}}
	r.result = append(r.result, innerValue)
}

