
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

function RecordRequestCountInAspectOne(val: RequestCount){
    CallBackFromUserScript_go("AspectOne",
                              {descriptorName : "request_count", value : val})}

function RecordRequestCountInAspectTwo(val: RequestCount) {
  CallBackFromUserScript_go("AspectTwo",
                            {descriptorName : "request_count", value : val})
}

function RecordRequestLatencyInAspectOne(val: RequestLatency) {
  CallBackFromUserScript_go("AspectOne",
                            {descriptorName : "request_latency", value : val})
}

function RecordRequestLatencyInAspectTwo(val: RequestLatency) {
  CallBackFromUserScript_go("AspectTwo",
                            {descriptorName : "request_latency", value : val})
}
