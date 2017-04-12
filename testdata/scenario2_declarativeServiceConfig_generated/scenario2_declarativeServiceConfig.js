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
function RecordRequestLatencyInPrometheusReportingAllMetrics(val) {
    CallBackFromUserScript_go('prometheus_reporting_all_metrics', { descriptorName: 'request_latency', value: val });
}
var Attributes = (function () {
    function Attributes(attribs) {
        // Fill the set of attribues that are part of the call (data is available
        // inside the attribs).
        if (attribs['response.latency'] !== undefined) {
            this.ResponseLatency = attribs['response.latency'];
        }
        if (attribs['api.method'] !== undefined) {
            this.ApiMethod = attribs['api.method'];
        }
        if (attribs['target.name'] !== undefined) {
            this.TargetName = attribs['target.name'];
        }
        if (attribs['api.name'] !== undefined) {
            this.ApiName = attribs['api.name'];
        }
        if (attribs['source.name'] !== undefined) {
            this.SourceName = attribs['source.name'];
        }
        if (attribs['response.code'] !== undefined) {
            this.ResponseCode = attribs['response.code'];
        }
    }
    return Attributes;
}());
/// <reference path="TypesFromAspectDescriptors.ts"/>
/// <reference path="WellKnownAttribs.ts"/>
function report(attributes) {
    if (true) {
        RecordRequestCountInPrometheusReportingAllMetrics({
            value: 1,
            response_code: attributes.ResponseCode !== undefined ? attributes.ResponseCode : 200,
            service: attributes.ApiName !== undefined ? attributes.ApiName :
                'unknown',
            source: attributes.SourceName !== undefined ? attributes.SourceName :
                'unknown',
            target: attributes.TargetName !== undefined ? attributes.TargetName :
                'unknown',
            method: attributes.ApiMethod !== undefined ? attributes.ApiMethod :
                'unknown'
        });
    }
}
function check(attributes) {
    // TODO
}
function quota(attributes) {
    // TODO
}
