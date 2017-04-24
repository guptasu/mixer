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
function RecordRequestCountInMyLocalMetricReporter(val) {
    __interal__callback_fn("MyLocalMetricReporter", { descriptorName: "request_count", value: val });
}
function RecordRequestLatencyInMyLocalMetricReporter(val) {
    __interal__callback_fn("MyLocalMetricReporter", { descriptorName: "request_latency", value: val });
}
