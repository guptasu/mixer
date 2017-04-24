/// <reference path="TypesFromAspectDescriptors.ts"/>

/// <reference path="WellKnownAttribs.ts"/>

var internalStressTestClientName = "MyInternalStressTestClient"
var internalMethodRegex = /__internalmtd__/;
var internalSystemMethodFriendlyName = "InternalSystemMethod"

function report(attributes: Attributes) {
    if (attributes.SourceName != internalStressTestClientName) { // skip report for stress tests calls
        var reqCount = new RequestCount();
        reqCount.value = 1;

        reqCount.response_code = attributes.ResponseCode !== undefined ?
            attributes.ResponseCode : 200;

        reqCount.service = attributes.ApiName !== undefined ?
            attributes.ApiName : "unknown";

        if (attributes.ApiMethod !== undefined) {
            // different internal methods can all be reported as single method.
            reqCount.method = attributes.ApiMethod.search(internalMethodRegex) != -1 ?
                internalSystemMethodFriendlyName : attributes.ApiMethod;
        } else {
            reqCount.method = "unknown";
        }

        RecordRequestCountInMyLocalMetricReporter(reqCount);
    }
}

function check(attributes) {
  // TODO
}
function quota(attributes) {
  // TODO
}
