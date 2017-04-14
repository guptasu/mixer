//-----------------CallBack Method Declaration-----------------
// This method gets injected at runtime. Need this declaration to make
// TypeScript happy
var __interal__callback_fn = function (aspectName, val) { };
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
function RecordRequestCountInMyAspect1(val) {
    __interal__callback_fn("MyAspect1", { descriptorName: "request_count", value: val });
}
function RecordRequestLatencyInMyAspect1(val) {
    __interal__callback_fn("MyAspect1", { descriptorName: "request_latency", value: val });
}
function ConstructRequestCountForMyAspect1(attributes) {
    return {
        value: 1,
        method: attributes.ApiMethod !== undefined ? attributes.ApiMethod
            : "unknown",
        response_code: attributes.ResponseCode !== undefined
            ? attributes.ResponseCode
            : 200,
        service: attributes.ApiName !== undefined ? attributes.ApiName
            : "unknown",
        source: attributes.SourceName !== undefined ? attributes.SourceName
            : "unknown",
        target: attributes.TargetName !== undefined ? attributes.TargetName
            : "unknown"
    };
}
var Attributes = (function () {
    function Attributes() {
    }
    return Attributes;
}());
/// <reference path="TypesFromAspectDescriptors.ts"/>
/// <reference path="WellKnownAttribs.ts"/>
function report(attributes) {
    if (true) {
        RecordRequestCountInMyAspect1(ConstructRequestCountForMyAspect1(attributes));
    }
}
function check(attributes) {
    // TODO
}
function quota(attributes) {
    // TODO
}
