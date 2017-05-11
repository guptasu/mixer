package main

import (
	"istio.io/mixer/pkg/attribute"
	"istio.io/mixer/pkg/adapter"
)
/////////// USER WRITTEN CODE ///////////
func ConstructRequestCountForPrometheusReportingAllMetrics(attributesBag *attribute.MutableBag) RequestCount {
	reqCount := RequestCount{}
	if attributesBag.WellKnownAttributes.Response.ResponseCode != 0 {
		reqCount.Value = attributesBag.WellKnownAttributes.Response.ResponseCode
	} else {
		reqCount.Value = 1001
	}

	if attributesBag.WellKnownAttributes.Source.SourceName != "" {
		reqCount.Source = attributesBag.WellKnownAttributes.Source.SourceName
	} else {
		reqCount.Source = "one1"
	}

	if attributesBag.WellKnownAttributes.Source.SourceName != "" {
		reqCount.Target = attributesBag.WellKnownAttributes.Source.SourceName
	} else {
		reqCount.Target = "one1"
	}

	if attributesBag.WellKnownAttributes.Response.ResponseCode != 0 {
		reqCount.ResponseCode = attributesBag.WellKnownAttributes.Response.ResponseCode
	} else {
		reqCount.ResponseCode = 1231111
	}

	if x,y := attributesBag.Get("foo.bar"); y {
		reqCount.Method = x.(string)
	} else {
		reqCount.Method = "one1"
	}

	if attributesBag.WellKnownAttributes.Request.RequestMethod != "" {
		reqCount.Service = attributesBag.WellKnownAttributes.Request.RequestMethod
	} else {
		reqCount.Service = "one1myservice"
	}



	return reqCount
}

func Report(attributes *attribute.MutableBag) [][]interface{} {
	var result = CreateReportResult()

	if true {
		result.InsertRequestCountForPrometheusReportingAllMetrics0(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
	}
	return result.result
}

/////////// COMMON GENERTED STUB CODE. NOT WRITTEN BY USER///////////
type RequestCount struct {
	Value        int64 `m:"value"`
	Service      string `m:"service"`
	Method       string `m:"method"`
	ResponseCode int64 `m:"response_code"`
	Source       string `m:"source"`
	Target       string `m:"target"`
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

