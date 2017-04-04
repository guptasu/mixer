/// <reference path="TypesFromAspectDescriptors.ts"/>

/// <reference path="WellKnownAttribs.ts"/>

function report(attributes: Attributes) {

    var reqcnt = new RequestCount();
    if (attributes.ResponseHttpCode !== undefined) {
        if (attributes.ResponseHttpCode >= 400) {
            reqcnt.response_code = 400
        } else {
            reqcnt.response_code = attributes.ResponseHttpCode
        }
    } else{
        reqcnt.response_code = 201
    }
    reqcnt.value = 20;
    reqcnt.method = 'one';
    reqcnt.service = 'one';
    reqcnt.source = 'one';
    reqcnt.target = 'one'

    RecordRequestCountInAspectOne(reqcnt)
}
function check(attributes) {
  // TODO
}
function quota(attributes) {
  // TODO
}
