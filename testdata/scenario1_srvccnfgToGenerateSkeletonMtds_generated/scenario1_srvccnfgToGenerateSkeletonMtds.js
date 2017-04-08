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
function RecordRequestCountInAspectOne(val) {
    CallBackFromUserScript_go("AspectOne", { descriptorName: "request_count", value: val });
}
function RecordRequestCountInAspectTwo(val) {
    CallBackFromUserScript_go("AspectTwo", { descriptorName: "request_count", value: val });
}
function RecordRequestLatencyInAspectOne(val) {
    CallBackFromUserScript_go("AspectOne", { descriptorName: "request_latency", value: val });
}
function RecordRequestLatencyInAspectTwo(val) {
    CallBackFromUserScript_go("AspectTwo", { descriptorName: "request_latency", value: val });
}
var Attributes = (function () {
    function Attributes(attribs) {
        // Fill the set of attribues that are part of the call (data is available
        // inside the attribs).
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
        if (attribs.Get('response.latency')[1]) {
            this.ResponseLatency = attribs.Get('response.latency')[0];
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
    var reqcnt = new RequestCount();
    if (attributes.ResponseHttpCode !== undefined) {
        if (attributes.ResponseHttpCode >= 400) {
            reqcnt.response_code = 400;
        }
        else {
            reqcnt.response_code = attributes.ResponseHttpCode;
        }
    }
    else {
        reqcnt.response_code = 201;
    }
    reqcnt.value = 20;
    reqcnt.method = 'one';
    reqcnt.service = 'one';
    reqcnt.source = 'one';
    reqcnt.target = 'one';
    RecordRequestCountInAspectOne(reqcnt);
}
function check(attributes) {
    // TODO
}
function quota(attributes) {
    // TODO
}
