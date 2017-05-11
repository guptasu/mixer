package main

import (
	"istio.io/mixer/pkg/attribute"
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
		reqCount.ResponseCode = 1231
	}

	if attributesBag.WellKnownAttributes.Request.RequestMethod != "" {
		reqCount.Method = attributesBag.WellKnownAttributes.Request.RequestMethod
	} else {
		reqCount.Method = "one1"
	}

	if attributesBag.WellKnownAttributes.Request.RequestMethod != "" {
		reqCount.Service = attributesBag.WellKnownAttributes.Request.RequestMethod
	} else {
		reqCount.Service = "one1"
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


func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics0(val RequestCount) {
	//v,_ := ToMap(val, "m")

	innerValue := []interface{}{"aspectName0", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

