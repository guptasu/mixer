
//-----------------CallBack Method Declaration-----------------
// This method gets injected at runtime. Need this declaration to make
// TypeScript happy
var __interal__callback_fn = function(aspectName: string, val: any) {};

//-----------------All Types Declaration-----------------
class RequestCount {
  value: number;
  response_code: number;
  source: string;
  target: string;
  service: string;
  method: string;
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
    __interal__callback_fn("MyAspect1",
                           {descriptorName : "request_count", value : val})}

function RecordRequestLatencyInMyAspect1(val: RequestLatency) {
  __interal__callback_fn("MyAspect1",
                         {descriptorName : "request_latency", value : val})
}

function ConstructRequestCountForMyAspect1(attributes: Attributes) {
  return {
    value: 1,
        service: attributes.ApiName !== undefined ? attributes.ApiName
                                                  : "unknown",
        source: attributes.SourceName !== undefined ? attributes.SourceName
                                                    : "unknown",
        target: attributes.TargetName !== undefined ? attributes.TargetName
                                                    : "unknown",
        method: attributes.ApiMethod !== undefined ? attributes.ApiMethod
                                                   : "unknown",
        response_code: attributes.ResponseCode !== undefined
            ? attributes.ResponseCode
            : 200
  }
}
