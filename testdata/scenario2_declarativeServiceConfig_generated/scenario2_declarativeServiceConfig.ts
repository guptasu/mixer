/// <reference path="TypesFromAspectDescriptors.ts"/>

/// <reference path="WellKnownAttribs.ts"/>


function report(attributes: Attributes): ReportResult {
  var result = new ReportResult();

  if (true) {
    result.InsertRequestCountForMyAspect1(
        ConstructRequestCountForMyAspect1(attributes))
  }

  return result;
}
function check(attributes) {
  // TODO
}
function quota(attributes) {
  // TODO
}
