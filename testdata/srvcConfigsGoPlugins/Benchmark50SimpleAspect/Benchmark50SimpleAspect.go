package main

import (
	"reflect"
	"fmt"
)
/////////// USER WRITTEN CODE ///////////
func ConstructRequestCountForPrometheusReportingAllMetrics(attributes map[string]interface{}) RequestCount {
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

func Report(attributes map[string]interface{}) [][]interface{} {
	var result = CreateReportResult()

	if true {
		result.InsertRequestCountForPrometheusReportingAllMetrics1(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics2(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics3(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics4(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics5(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics6(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics7(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics8(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics9(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics10(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics11(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics12(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics13(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics14(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics15(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics16(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics17(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics18(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics19(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics20(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics21(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics22(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics23(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics24(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics25(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics26(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics27(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics28(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics29(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics30(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics31(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics32(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics33(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics34(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics35(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics36(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics37(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics38(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics39(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics40(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics41(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics42(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics43(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics44(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics45(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics46(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics47(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics48(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
		result.InsertRequestCountForPrometheusReportingAllMetrics49(ConstructRequestCountForPrometheusReportingAllMetrics(attributes))
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


func ToMap(in interface{}, tag string) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ToMap only accepts structs; got %T", v)
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		if tagv := fi.Tag.Get(tag); tagv != "" {
			out[tagv] = v.Field(i).Interface()
		}
	}
	return out, nil
}


func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics0(val RequestCount) {
	innerValue := []interface{}{"aspectName0", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics1(val RequestCount) {
	innerValue := []interface{}{"aspectName1", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics2(val RequestCount) {
	innerValue := []interface{}{"aspectName2", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics3(val RequestCount) {
	innerValue := []interface{}{"aspectName3", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics4(val RequestCount) {
	innerValue := []interface{}{"aspectName4", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics5(val RequestCount) {
	innerValue := []interface{}{"aspectName5", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics6(val RequestCount) {
	innerValue := []interface{}{"aspectName6", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics7(val RequestCount) {
	innerValue := []interface{}{"aspectName7", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics8(val RequestCount) {
	innerValue := []interface{}{"aspectName8", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics9(val RequestCount) {
	innerValue := []interface{}{"aspectName9", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics10(val RequestCount) {
	innerValue := []interface{}{"aspectName10", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics11(val RequestCount) {
	innerValue := []interface{}{"aspectName11", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics12(val RequestCount) {
	innerValue := []interface{}{"aspectName12", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics13(val RequestCount) {
	innerValue := []interface{}{"aspectName13", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics14(val RequestCount) {
	innerValue := []interface{}{"aspectName14", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics15(val RequestCount) {
	innerValue := []interface{}{"aspectName15", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics16(val RequestCount) {
	innerValue := []interface{}{"aspectName16", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics17(val RequestCount) {
	innerValue := []interface{}{"aspectName17", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics18(val RequestCount) {
	innerValue := []interface{}{"aspectName18", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics19(val RequestCount) {
	innerValue := []interface{}{"aspectName19", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics20(val RequestCount) {
	innerValue := []interface{}{"aspectName20", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics21(val RequestCount) {
	innerValue := []interface{}{"aspectName21", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics22(val RequestCount) {
	innerValue := []interface{}{"aspectName22", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics23(val RequestCount) {
	innerValue := []interface{}{"aspectName23", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics24(val RequestCount) {
	innerValue := []interface{}{"aspectName24", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics25(val RequestCount) {
	innerValue := []interface{}{"aspectName25", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics26(val RequestCount) {
	innerValue := []interface{}{"aspectName26", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics27(val RequestCount) {
	innerValue := []interface{}{"aspectName27", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics28(val RequestCount) {
	innerValue := []interface{}{"aspectName28", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics29(val RequestCount) {
	innerValue := []interface{}{"aspectName29", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics30(val RequestCount) {
	innerValue := []interface{}{"aspectName30", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics31(val RequestCount) {
	innerValue := []interface{}{"aspectName31", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics32(val RequestCount) {
	innerValue := []interface{}{"aspectName32", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics33(val RequestCount) {
	innerValue := []interface{}{"aspectName33", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics34(val RequestCount) {
	innerValue := []interface{}{"aspectName34", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics35(val RequestCount) {
	innerValue := []interface{}{"aspectName35", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics36(val RequestCount) {
	innerValue := []interface{}{"aspectName36", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics37(val RequestCount) {
	innerValue := []interface{}{"aspectName37", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics38(val RequestCount) {
	innerValue := []interface{}{"aspectName38", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics39(val RequestCount) {
	innerValue := []interface{}{"aspectName39", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics40(val RequestCount) {
	innerValue := []interface{}{"aspectName40", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics41(val RequestCount) {
	innerValue := []interface{}{"aspectName41", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics42(val RequestCount) {
	innerValue := []interface{}{"aspectName42", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics43(val RequestCount) {
	innerValue := []interface{}{"aspectName43", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics44(val RequestCount) {
	innerValue := []interface{}{"aspectName44", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics45(val RequestCount) {
	innerValue := []interface{}{"aspectName45", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics46(val RequestCount) {
	innerValue := []interface{}{"aspectName46", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics47(val RequestCount) {
	innerValue := []interface{}{"aspectName47", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics48(val RequestCount) {
	innerValue := []interface{}{"aspectName48", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

func (r *ReportResult) InsertRequestCountForPrometheusReportingAllMetrics49(val RequestCount) {
	innerValue := []interface{}{"aspectName49", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}
