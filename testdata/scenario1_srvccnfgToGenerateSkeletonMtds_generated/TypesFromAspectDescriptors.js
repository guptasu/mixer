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
