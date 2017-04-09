/// <reference path="TypesFromAspectDescriptors.ts"/>
/// <reference path="WellKnownAttribs.ts"/>
function report(attributes) {
    if (true) {
        RecordRequestCountInPrometheusReportingAllMetrics({
            value: 1,
            response_code: attributes.ResponseCode !== undefined ? attributes.ResponseCode : 200,
            service: attributes.ApiName !== undefined ? attributes.ApiName
                : "unknown",
            source: attributes.SourceName !== undefined ? attributes.SourceName
                : "unknown",
            target: attributes.TargetName !== undefined ? attributes.TargetName
                : "unknown",
            method: attributes.ApiMethod !== undefined ? attributes.ApiMethod
                : "unknown"
        });
    }
}
function check(attributes) {
    // TODO
}
function quota(attributes) {
    // TODO
}
