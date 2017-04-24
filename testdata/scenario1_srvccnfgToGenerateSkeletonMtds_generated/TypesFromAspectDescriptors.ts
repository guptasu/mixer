
//-----------------CallBack Method Declaration-----------------
// This method gets injected at runtime. Need this declaration to make
// TypeScript happy
var __interal__callback_fn = function(aspectName: string, val: any) {};

//-----------------All Types Declaration-----------------
class RequestCount {
  value: number;
  target: string;
  service: string;
  method: string;
  response_code: number;
  source: string;
}
class RequestLatency {
  value: number;
  service: string;
  method: string;
  response_code: number;
  source: string;
  target: string;
}

function RecordRequestCountInMyLocalMetricReporter(val: RequestCount){
    __interal__callback_fn("MyLocalMetricReporter",
                           {descriptorName : "request_count", value : val})}

function RecordRequestLatencyInMyLocalMetricReporter(val: RequestLatency) {
  __interal__callback_fn("MyLocalMetricReporter",
                         {descriptorName : "request_latency", value : val})
}
