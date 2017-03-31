
class Attributes {
  // All the well known attribute names.
  ResponseLatency: number;
  ApiMethod: string;
  TargetName: string;
  ApiName: string;
  SourceName: string;
  ResponseHttpCode: number;

  constructor(attribs: any) {
    // Fill the set of attribues that are part of the call (data is available
    // inside the attribs).


    if (attribs.Get('response.latency')[1]) {
      this.ResponseLatency = attribs.Get('response.latency')[0]
    }

    if (attribs.Get('api.method')[1]) {
      this.ApiMethod = attribs.Get('api.method')[0]
    }

    if (attribs.Get('target.name')[1]) {
      this.TargetName = attribs.Get('target.name')[0]
    }

    if (attribs.Get('api.name')[1]) {
      this.ApiName = attribs.Get('api.name')[0]
    }

    if (attribs.Get('source.name')[1]) {
      this.SourceName = attribs.Get('source.name')[0]
    }

    if (attribs.Get('response.http.code')[1]) {
      this.ResponseHttpCode = attribs.Get('response.http.code')[0]
    }
  }
}
function ConstructAttributes(attr: any): Attributes {
  return new Attributes(attr)
}