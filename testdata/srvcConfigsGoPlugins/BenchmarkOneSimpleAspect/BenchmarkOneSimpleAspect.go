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
		reqCount.ResponseCode = 1231111
	}

	if x, y := attributesBag.Get("foo.bar"); y {
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
	result.InsertRequestCountForPrometheusReportingAllMetrics0(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
	return result.result
}
