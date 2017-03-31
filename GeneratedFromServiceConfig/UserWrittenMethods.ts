/// <reference path="TypesFromAspectDescriptors.ts"/>

/// <reference path="WellKnownAttribs.ts"/>


function report(attributes: Attributes) {
  if (attributes.SourceName == 'test') {
    RecordRequestCount({
      value: attributes.ResponseLatency !== undefined ?
          attributes.ResponseLatency :
          100,
      service: attributes.ApiName !== undefined ? attributes.ApiName : 'one',
      source: attributes.SourceName !== undefined ? attributes.SourceName :
                                                    'one',
      target: attributes.TargetName !== undefined ? attributes.TargetName :
                                                    'one',
      method: attributes.ApiMethod !== undefined ? attributes.ApiMethod : 'one',
      response_code: attributes.ResponseHttpCode !== undefined ?
          attributes.ResponseHttpCode :
          111
    })

    RecordRequestLatency({
      value: attributes.ResponseLatency !== undefined ?
          attributes.ResponseLatency :
          200,
      source: attributes.SourceName !== undefined ? attributes.SourceName :
                                                    'two',
      target: attributes.TargetName !== undefined ? attributes.TargetName :
                                                    'two',
      method: attributes.ApiMethod !== undefined ? attributes.ApiMethod : 'two',
      response_code: attributes.ResponseHttpCode !== undefined ?
          attributes.ResponseHttpCode :
          222,
      service: attributes.ApiName !== undefined ? attributes.ApiName : 'two'
    })
  }

  if (attributes.SourceName == 'foo') {
    RecordRequestLatency({
      value: attributes.ResponseLatency !== undefined ?
          attributes.ResponseLatency :
          300,
      response_code: attributes.ResponseHttpCode !== undefined ?
          attributes.ResponseHttpCode :
          333,
      service: attributes.ApiName !== undefined ? attributes.ApiName : 'three',
      source: attributes.SourceName !== undefined ? attributes.SourceName :
                                                    'three',
      target: attributes.TargetName !== undefined ? attributes.TargetName :
                                                    'three',
      method: attributes.ApiMethod !== undefined ? attributes.ApiMethod :
                                                   'three'
    })

    RecordRequestCount({
      value: 400,
      target: attributes.TargetName !== undefined ? attributes.TargetName :
                                                    'four',
      method: attributes.ApiMethod !== undefined ? attributes.ApiMethod :
                                                   'four',
      response_code: attributes.ResponseHttpCode !== undefined ?
          attributes.ResponseHttpCode :
          444,
      service: attributes.ApiName !== undefined ? attributes.ApiName : 'four',
      source: attributes.SourceName !== undefined ? attributes.SourceName :
                                                    'four'
    })
  }
}
function check(attributes) {
  // TODO
}
function quota(attributes) {
  // TODO
}
