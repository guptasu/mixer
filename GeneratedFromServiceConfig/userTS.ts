/// <reference path="typeDefs.ts"/>

/// <reference path="attribs.ts"/>


function report(attributes: Attributes) {
  if (attributes.SourceName == 'test') {
    RecordRequestCount({
      value: attributes.ResponseLatency !== undefined ?
          attributes.ResponseLatency :
          100,
      target: attributes.TargetName !== undefined ? attributes.TargetName :
                                                    'one',
      method: attributes.ApiMethod !== undefined ? attributes.ApiMethod : 'one',
      response_code: attributes.ResponseHttpCode !== undefined ?
          attributes.ResponseHttpCode :
          111,
      service: attributes.ApiName !== undefined ? attributes.ApiName : 'one',
      source: attributes.SourceName !== undefined ? attributes.SourceName :
                                                    'one'
    })

    RecordRequestLatency({
      value: attributes.ResponseLatency !== undefined ?
          attributes.ResponseLatency :
          200,
      method: attributes.ApiMethod !== undefined ? attributes.ApiMethod : 'two',
      response_code: attributes.ResponseHttpCode !== undefined ?
          attributes.ResponseHttpCode :
          222,
      service: attributes.ApiName !== undefined ? attributes.ApiName : 'two',
      source: attributes.SourceName !== undefined ? attributes.SourceName :
                                                    'two',
      target: attributes.TargetName !== undefined ? attributes.TargetName :
                                                    'two'
    })
  }

  if (attributes.SourceName == 'foo') {
    RecordRequestLatency({
      value: attributes.ResponseLatency !== undefined ?
          attributes.ResponseLatency :
          300,
      source: attributes.SourceName !== undefined ? attributes.SourceName :
                                                    'three',
      target: attributes.TargetName !== undefined ? attributes.TargetName :
                                                    'three',
      method: attributes.ApiMethod !== undefined ? attributes.ApiMethod :
                                                   'three',
      response_code: attributes.ResponseHttpCode !== undefined ?
          attributes.ResponseHttpCode :
          333,
      service: attributes.ApiName !== undefined ? attributes.ApiName : 'three'
    })

    RecordRequestCount({
      value: 400,
      response_code: attributes.ResponseHttpCode !== undefined ?
          attributes.ResponseHttpCode :
          444,
      service: attributes.ApiName !== undefined ? attributes.ApiName : 'four',
      source: attributes.SourceName !== undefined ? attributes.SourceName :
                                                    'four',
      target: attributes.TargetName !== undefined ? attributes.TargetName :
                                                    'four',
      method: attributes.ApiMethod !== undefined ? attributes.ApiMethod : 'four'
    })
  }
}
function check(attributes) {
  // TODO
}
function quota(attributes) {
  // TODO
}
