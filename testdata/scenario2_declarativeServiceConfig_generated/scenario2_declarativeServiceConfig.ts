/// <reference path="TypesFromAspectDescriptors.ts"/>

/// <reference path="WellKnownAttribs.ts"/>


function report(attributes: Attributes) {
  if (true) {
    RecordRequestCountInMyAspect1({
      value: 1,
      source: attributes.SourceName !== undefined ? attributes.SourceName :
                                                    'unknown',
      target: attributes.TargetName !== undefined ? attributes.TargetName :
                                                    'unknown',
      method: attributes.ApiMethod !== undefined ? attributes.ApiMethod :
                                                   'unknown',
      response_code:
          attributes.ResponseCode !== undefined ? attributes.ResponseCode : 200,
      service: attributes.ApiName !== undefined ? attributes.ApiName : 'unknown'
    })
  }
}
function check(attributes) {
  // TODO
}
function quota(attributes) {
  // TODO
}
