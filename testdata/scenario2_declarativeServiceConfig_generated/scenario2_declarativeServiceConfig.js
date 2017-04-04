/// <reference path="TypesFromAspectDescriptors.ts"/>
/// <reference path="WellKnownAttribs.ts"/>
function report(attributes) {
    if (true) {
        RecordRequestCountInPrometheusReportingAllMetrics({
            value: attributes.ResponseLatency !== undefined
                ? attributes.ResponseLatency
                : 100,
            source: attributes.SourceName !== undefined ? attributes.SourceName
                : "one",
            target: attributes.TargetName !== undefined ? attributes.TargetName
                : "one",
            method: attributes.ApiMethod !== undefined ? attributes.ApiMethod
                : "one",
            response_code: attributes.ResponseHttpCode !== undefined
                ? attributes.ResponseHttpCode
                : 111,
            service: attributes.ApiName !== undefined ? attributes.ApiName : "one"
        });
        RecordRequestLatencyInPrometheusReportingAllMetrics({
            value: attributes.ResponseLatency !== undefined
                ? attributes.ResponseLatency
                : 2000,
            service: attributes.ApiName !== undefined ? attributes.ApiName : "two",
            source: attributes.SourceName !== undefined ? attributes.SourceName
                : "two",
            target: attributes.TargetName !== undefined ? attributes.TargetName
                : "two",
            method: attributes.ApiMethod !== undefined ? attributes.ApiMethod
                : "two",
            response_code: attributes.ResponseHttpCode !== undefined
                ? attributes.ResponseHttpCode
                : 222
        });
    }
}
function check(attributes) {
    // TODO
}
function quota(attributes) {
    // TODO
}
