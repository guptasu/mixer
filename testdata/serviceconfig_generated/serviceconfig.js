//-----------------CallBack Method Declaration-----------------
// This method gets injected at runtime. Need this declaration to make
// TypeScript happy
var CallBackFromUserScript_go = function (aspectName, val) { };
//-----------------All Types Declaration-----------------
var RequestCount = (function () {
    function RequestCount() {
    }
    return RequestCount;
}());
var RequestLatency = (function () {
    function RequestLatency() {
    }
    return RequestLatency;
}());
function RecordRequestCountInPrometheusReportingAllMetrics(val) {
    CallBackFromUserScript_go('prometheus_reporting_all_metrics', { descriptorName: 'request_count', value: val });
}
function RecordRequestCountInPrometheusReportingJustReqLatency(val) {
    CallBackFromUserScript_go('prometheus_reporting_just_req_latency', { descriptorName: 'request_count', value: val });
}
function RecordRequestCountInPrometheusReportingJustReqCount(val) {
    CallBackFromUserScript_go('prometheus_reporting_just_req_count', { descriptorName: 'request_count', value: val });
}
function RecordRequestLatencyInPrometheusReportingAllMetrics(val) {
    CallBackFromUserScript_go('prometheus_reporting_all_metrics', { descriptorName: 'request_latency', value: val });
}
function RecordRequestLatencyInPrometheusReportingJustReqLatency(val) {
    CallBackFromUserScript_go('prometheus_reporting_just_req_latency', { descriptorName: 'request_latency', value: val });
}
function RecordRequestLatencyInPrometheusReportingJustReqCount(val) {
    CallBackFromUserScript_go('prometheus_reporting_just_req_count', { descriptorName: 'request_latency', value: val });
}
var Attributes = (function () {
    function Attributes(attribs) {
        // Fill the set of attribues that are part of the call (data is available
        // inside the attribs).
        if (attribs.Get('response.latency')[1]) {
            this.ResponseLatency = attribs.Get('response.latency')[0];
        }
        if (attribs.Get('api.method')[1]) {
            this.ApiMethod = attribs.Get('api.method')[0];
        }
        if (attribs.Get('target.name')[1]) {
            this.TargetName = attribs.Get('target.name')[0];
        }
        if (attribs.Get('api.name')[1]) {
            this.ApiName = attribs.Get('api.name')[0];
        }
        if (attribs.Get('source.name')[1]) {
            this.SourceName = attribs.Get('source.name')[0];
        }
        if (attribs.Get('response.http.code')[1]) {
            this.ResponseHttpCode = attribs.Get('response.http.code')[0];
        }
    }
    return Attributes;
}());
function ConstructAttributes(attr) {
    return new Attributes(attr);
}
/// <reference path="TypesFromAspectDescriptors.ts"/>
/// <reference path="WellKnownAttribs.ts"/>
function report(attributes) {
    if (attributes.SourceName == 'test') {
        RecordRequestCountInPrometheusReportingAllMetrics({
            value: attributes.ResponseLatency !== undefined ?
                attributes.ResponseLatency :
                100,
            method: attributes.ApiMethod !== undefined ? attributes.ApiMethod : 'one',
            response_code: attributes.ResponseHttpCode !== undefined ?
                attributes.ResponseHttpCode :
                111,
            service: attributes.ApiName !== undefined ? attributes.ApiName : 'one',
            source: attributes.SourceName !== undefined ? attributes.SourceName :
                'one',
            target: attributes.TargetName !== undefined ? attributes.TargetName :
                'one'
        });
        RecordRequestLatencyInPrometheusReportingAllMetrics({
            value: attributes.ResponseLatency !== undefined ?
                attributes.ResponseLatency :
                2000,
            service: attributes.ApiName !== undefined ? attributes.ApiName : 'two',
            source: attributes.SourceName !== undefined ? attributes.SourceName :
                'two',
            target: attributes.TargetName !== undefined ? attributes.TargetName :
                'two',
            method: attributes.ApiMethod !== undefined ? attributes.ApiMethod : 'two',
            response_code: attributes.ResponseHttpCode !== undefined ?
                attributes.ResponseHttpCode :
                222
        });
    }
    if (attributes.SourceName == 'foo') {
        RecordRequestLatencyInPrometheusReportingJustReqLatency({
            value: attributes.ResponseLatency !== undefined ?
                attributes.ResponseLatency :
                300,
            response_code: attributes.ResponseHttpCode !== undefined ?
                attributes.ResponseHttpCode :
                333,
            service: attributes.ApiName !== undefined ? attributes.ApiName : 'three',
            source: attributes.SourceName !== undefined ? attributes.SourceName :
                'three',
            target: attributes.TargetName !== undefined ? attributes.TargetName :
                'three',
            method: attributes.ApiMethod !== undefined ? attributes.ApiMethod :
                'three'
        });
        RecordRequestCountInPrometheusReportingJustReqCount({
            value: 400,
            service: attributes.ApiName !== undefined ? attributes.ApiName : 'four',
            source: attributes.SourceName !== undefined ? attributes.SourceName :
                'four',
            target: attributes.TargetName !== undefined ? attributes.TargetName :
                'four',
            method: attributes.ApiMethod !== undefined ? attributes.ApiMethod :
                'four',
            response_code: attributes.ResponseHttpCode !== undefined ?
                attributes.ResponseHttpCode :
                444
        });
    }
}
function check(attributes) {
    // TODO
}
function quota(attributes) {
    // TODO
}
