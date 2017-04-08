
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
    CallBackFromUserScript_go("prometheus_reporting_all_metrics",
                              {descriptorName : "request_count", value : val})}

function RecordRequestLatencyInPrometheusReportingAllMetrics(
    val: RequestLatency) {
  CallBackFromUserScript_go("prometheus_reporting_all_metrics",
                            {descriptorName : "request_latency", value : val})
}
