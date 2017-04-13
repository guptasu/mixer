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
function RecordRequestCountInMyAspect1(val) {
    CallBackFromUserScript_go('MyAspect1', { descriptorName: 'request_count', value: val });
}
function RecordRequestLatencyInMyAspect1(val) {
    CallBackFromUserScript_go('MyAspect1', { descriptorName: 'request_latency', value: val });
}
function ConstructRequestCountForMyAspect1(attributes) {
    return {
        value: 1,
        source: attributes.SourceName !== undefined ? attributes.SourceName :
            'unknown',
        target: attributes.TargetName !== undefined ? attributes.TargetName :
            'unknown',
        method: attributes.ApiMethod !== undefined ?
            attributes.ApiMethod :
            attributes.ApiName !== undefined ? attributes.ApiMethod !== undefined ?
                attributes.ApiMethod :
                attributes.ApiName :
                'unknown',
        response_code: attributes.ResponseCode !== undefined ?
            attributes.ResponseCode :
            200,
        service: attributes.ApiName !== undefined ? attributes.ApiName :
            'unknown'
    };
}
var Attributes = (function () {
    function Attributes(attribs) {
        // Fill the set of attribues that are part of the call (data is available
        // inside the attribs).
        if (attribs['api.name'] !== undefined) {
            this.ApiName = attribs['api.name'];
        }
        if (attribs['source.name'] !== undefined) {
            this.SourceName = attribs['source.name'];
        }
        if (attribs['response.code'] !== undefined) {
            this.ResponseCode = attribs['response.code'];
        }
        if (attribs['response.latency'] !== undefined) {
            this.ResponseLatency = attribs['response.latency'];
        }
        if (attribs['api.method'] !== undefined) {
            this.ApiMethod = attribs['api.method'];
        }
        if (attribs['target.name'] !== undefined) {
            this.TargetName = attribs['target.name'];
        }
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
