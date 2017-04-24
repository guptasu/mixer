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
var ReportResult = (function () {
    function ReportResult() {
        this.result = [];
    }
    ReportResult.prototype.InsertRequestCountForMyAspect1 = function (val) {
        this.result.push(['MyAspect1', { descriptorName: 'request_count', value: val }]);
    };
    ReportResult.prototype.InsertRequestLatencyForMyAspect1 = function (val) {
        this.result.push(['MyAspect1', { descriptorName: 'request_latency', value: val }]);
    };
    ReportResult.prototype.Build = function () {
        return this.result;
    };
    return ReportResult;
}());
function ConstructRequestCountForMyAspect1(attributes) {
    return {
        value: 1,
        target: attributes.TargetName !== undefined ? attributes.TargetName :
            'unknown',
        method: attributes.ApiMethod !== undefined ? attributes.ApiMethod :
            'unknown',
        response_code: attributes.ResponseCode !== undefined ?
            attributes.ResponseCode :
            200,
        service: attributes.ApiName !== undefined ? attributes.ApiName :
            'unknown',
        source: attributes.SourceName !== undefined ? attributes.SourceName :
            'unknown'
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
    var result = new ReportResult();
    if (true) {
        result.InsertRequestCountForMyAspect1(ConstructRequestCountForMyAspect1(attributes));
    }
    return result;
}
function check(attributes) {
    // TODO
}
function quota(attributes) {
    // TODO
}
