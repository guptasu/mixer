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
