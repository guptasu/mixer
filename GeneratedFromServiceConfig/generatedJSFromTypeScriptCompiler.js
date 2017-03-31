//-----------------CallBack Method Declaration-----------------
// This method gets injected at runtime. Need this declaration to make
// TypeScript happy
var CallBackFromUserScript_go = function (name, val) { };
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
function RecordRequestCount(val) {
    CallBackFromUserScript_go('metrics', { descriptorName: 'request_count', value: val });
}
function RecordRequestLatency(val) {
    CallBackFromUserScript_go('metrics', { descriptorName: 'request_latency', value: val });
}
var Attributes = (function () {
    function Attributes(attribs) {
        // Fill the set of attribues that are part of the call (data is available
        // inside the attribs).
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
        if (attribs.Get('response.latency')[1]) {
            this.ResponseLatency = attribs.Get('response.latency')[0];
        }
        if (attribs.Get('api.method')[1]) {
            this.ApiMethod = attribs.Get('api.method')[0];
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
        RecordRequestCount({
            value: attributes.ResponseLatency !== undefined ?
                attributes.ResponseLatency :
                100,
            service: attributes.ApiName !== undefined ? attributes.ApiName : 'one',
            source: attributes.SourceName !== undefined ? attributes.SourceName :
                'one',
            target: attributes.TargetName !== undefined ? attributes.TargetName :
                'one',
            method: attributes.ApiMethod !== undefined ? attributes.ApiMethod : 'one',
            response_code: attributes.ResponseHttpCode !== undefined ?
                attributes.ResponseHttpCode :
                111
        });
        RecordRequestLatency({
            value: attributes.ResponseLatency !== undefined ?
                attributes.ResponseLatency :
                200,
            source: attributes.SourceName !== undefined ? attributes.SourceName :
                'two',
            target: attributes.TargetName !== undefined ? attributes.TargetName :
                'two',
            method: attributes.ApiMethod !== undefined ? attributes.ApiMethod : 'two',
            response_code: attributes.ResponseHttpCode !== undefined ?
                attributes.ResponseHttpCode :
                222,
            service: attributes.ApiName !== undefined ? attributes.ApiName : 'two'
        });
    }
    if (attributes.SourceName == 'foo') {
        RecordRequestLatency({
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
        RecordRequestCount({
            value: 400,
            target: attributes.TargetName !== undefined ? attributes.TargetName :
                'four',
            method: attributes.ApiMethod !== undefined ? attributes.ApiMethod :
                'four',
            response_code: attributes.ResponseHttpCode !== undefined ?
                attributes.ResponseHttpCode :
                444,
            service: attributes.ApiName !== undefined ? attributes.ApiName : 'four',
            source: attributes.SourceName !== undefined ? attributes.SourceName :
                'four'
        });
    }
}
function check(attributes) {
    // TODO
}
function quota(attributes) {
    // TODO
}
