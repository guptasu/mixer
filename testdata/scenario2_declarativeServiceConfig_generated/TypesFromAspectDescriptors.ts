
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

function RecordRequestCountInMyAspect1(val: RequestCount){
    CallBackFromUserScript_go(
        'MyAspect1', {descriptorName: 'request_count', value: val})}

function RecordRequestLatencyInMyAspect1(val: RequestLatency) {
  CallBackFromUserScript_go(
      'MyAspect1', {descriptorName: 'request_latency', value: val})
}
