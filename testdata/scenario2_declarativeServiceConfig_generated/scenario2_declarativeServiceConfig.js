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
    CallBackFromUserScript_go("prometheus_reporting_all_metrics", { descriptorName: "request_count", value: val });
}
function RecordRequestLatencyInPrometheusReportingAllMetrics(val) {
    CallBackFromUserScript_go("prometheus_reporting_all_metrics", { descriptorName: "request_latency", value: val });
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
        if (attribs.Get('response.code')[1]) {
            this.ResponseCode = attribs.Get('response.code')[0];
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
    if (true) {
        RecordRequestCountInPrometheusReportingAllMetrics({
            value: 1,
            target: attributes.TargetName !== undefined ? attributes.TargetName
                : "unknown",
            method: attributes.ApiMethod !== undefined ? attributes.ApiMethod
                : "unknown",
            response_code: attributes.ResponseCode !== undefined ? attributes.ResponseCode : 200,
            service: attributes.ApiName !== undefined ? attributes.ApiName
                : "unknown",
            source: attributes.SourceName !== undefined ? attributes.SourceName
                : "unknown"
        });
    }
}
function check(attributes) {
    // TODO
}
function quota(attributes) {
    // TODO
}
