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
	//v,_ := ToMap(val, "m")

	innerValue := []interface{}{"aspectName0", map[string]interface{}{"descriptorName": "request_count", "value": val}}
	r.result = append(r.result, innerValue)
}

