/// <reference path="TypesFromAspectDescriptors.ts"/>

/// <reference path="WellKnownAttribs.ts"/>


function report(attributes: Attributes) {
    RecordRequestCountInPrometheusReportingAllMetrics({
        value: attributes.ResponseLatency !== undefined ?
            attributes.ResponseLatency :
            100,
        method: attributes.ApiMethod !== undefined ? attributes.ApiMethod : 'one',
        response_code: attributes.ResponseHttpCode !== undefined ?
            attributes.ResponseHttpCode :
            111,
        service: attributes.ApiName !== undefined ? attributes.ApiName : 'one',
        source: attributes.SourceName !== undefined ? attributes.SourceName :
            'one',
        target: attributes.TargetName !== undefined ? attributes.TargetName :
            'JUST EDITED one'
    })
}
function check(attributes) {
  // TODO
}
function quota(attributes) {
  // TODO
}
