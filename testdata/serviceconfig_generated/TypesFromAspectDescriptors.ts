
//-----------------CallBack Method Declaration-----------------
// This method gets injected at runtime. Need this declaration to make
// TypeScript happy
var CallBackFromUserScript_go = function(aspectName: string, val: any) {};

//-----------------All Types Declaration-----------------
class RequestCount {
  value: number;
  source: string;
  target: string;
  service: string;
  method: string;
  response_code: number;
}
class RequestLatency {
  value: number;
  source: string;
  target: string;
  service: string;
  method: string;
  response_code: number;
}

function RecordRequestCountInPrometheusReportingAllMetrics(val: RequestCount){
    CallBackFromUserScript_go(
        'prometheus_reporting_all_metrics',
        {descriptorName: 'request_count', value: val})}

function RecordRequestCountInPrometheusReportingJustReqLatency(
    val: RequestCount) {
  CallBackFromUserScript_go(
      'prometheus_reporting_just_req_latency',
      {descriptorName: 'request_count', value: val})
}

function RecordRequestCountInPrometheusReportingJustReqCount(
    val: RequestCount) {
  CallBackFromUserScript_go(
      'prometheus_reporting_just_req_count',
      {descriptorName: 'request_count', value: val})
}

function RecordRequestLatencyInPrometheusReportingAllMetrics(
    val: RequestLatency) {
  CallBackFromUserScript_go(
      'prometheus_reporting_all_metrics',
      {descriptorName: 'request_latency', value: val})
}

function RecordRequestLatencyInPrometheusReportingJustReqLatency(
    val: RequestLatency) {
  CallBackFromUserScript_go(
      'prometheus_reporting_just_req_latency',
      {descriptorName: 'request_latency', value: val})
}

function RecordRequestLatencyInPrometheusReportingJustReqCount(
    val: RequestLatency) {
  CallBackFromUserScript_go(
      'prometheus_reporting_just_req_count',
      {descriptorName: 'request_latency', value: val})
}
