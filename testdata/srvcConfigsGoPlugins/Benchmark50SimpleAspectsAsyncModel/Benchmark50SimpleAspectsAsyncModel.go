package main

/////////// USER WRITTEN CODE ///////////
func ConstructRequestCountForPrometheusReportingAllMetrics(attributes map[string]interface{}) interface{} {
	reqCount := RequestCount{}
	if attributes["ResponseCode"] == 0 {
		reqCount.Value = attributes["ResponseCode"].(int64)
	} else {
		reqCount.Value = 1001
	}

	if attributes["SourceName"] == "" {
		reqCount.Source = attributes["SourceName"].(string)
	} else {
		reqCount.Source = "one1"
	}

	if attributes["SourceName"] == "" {
		reqCount.Target = attributes["SourceName"].(string)
	} else {
		reqCount.Target = "one1"
	}

	if attributes["ResponseCode"] == 0 {
		reqCount.ResponseCode = attributes["ResponseCode"].(int64)
	} else {
		reqCount.ResponseCode = 1231
	}

	if attributes["ApiMethod"] == "" {
		reqCount.Method = attributes["ApiMethod"].(string)
	} else {
		reqCount.Method = "one1"
	}

	if attributes["ApiName"] == "" {
		reqCount.Service = attributes["ApiName"].(string)
	} else {
		reqCount.Service = "one1"
	}

	return reqCount
}

func Report(attributes map[string]interface{}) (aspectNameAndFunctionForEval map[string]string) {

	aspectNameAndFunctionForEval = make(map[string]string)

	if true {
		aspectNameAndFunctionForEval["aspectName0"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName1"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName2"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName3"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName4"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName5"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName6"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName7"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName8"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName9"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName10"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName11"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName12"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName13"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName14"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName15"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName16"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName17"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName18"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName19"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName20"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName21"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName22"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName23"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName24"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName25"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName26"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName27"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName28"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName29"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName30"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName31"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName32"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName33"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName34"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName35"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName36"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName37"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName38"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName39"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName40"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName41"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName42"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName43"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName44"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName45"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName46"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName47"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName48"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
		aspectNameAndFunctionForEval["aspectName49"] = "ConstructRequestCountForPrometheusReportingAllMetrics"
	}
	return aspectNameAndFunctionForEval
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
