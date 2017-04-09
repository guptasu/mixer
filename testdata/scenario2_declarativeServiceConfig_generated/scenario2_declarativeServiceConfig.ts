/// <reference path="TypesFromAspectDescriptors.ts"/>

/// <reference path="WellKnownAttribs.ts"/>

function report(attributes: Attributes) {

  if (true) {

    RecordRequestCountInPrometheusReportingAllMetrics({
      value : 1,
      target : attributes.TargetName !== undefined ? attributes.TargetName
                                                   : "unknown",
      method : attributes.ApiMethod !== undefined ? attributes.ApiMethod
                                                  : "unknown",
      response_code :
          attributes.ResponseCode !== undefined ? attributes.ResponseCode : 200,
      service : attributes.ApiName !== undefined ? attributes.ApiName
                                                 : "unknown",
      source : attributes.SourceName !== undefined ? attributes.SourceName
                                                   : "unknown"
    })
  }
}
function check(attributes) {
  // TODO
}
function quota(attributes) {
  // TODO
}
