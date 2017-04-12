
class Attributes {
  // All the well known attribute names.
  ResponseCode: number;
  ResponseLatency: number;
  ApiMethod: string;
  TargetName: string;
  ApiName: string;
  SourceName: string;

  constructor(attribs: any) {
    // Fill the set of attribues that are part of the call (data is available
    // inside the attribs).


    if (attribs['response.code'] !== undefined) {
      this.ResponseCode = attribs['response.code']
    }

    if (attribs['response.latency'] !== undefined) {
      this.ResponseLatency = attribs['response.latency']
    }

    if (attribs['api.method'] !== undefined) {
      this.ApiMethod = attribs['api.method']
    }

    if (attribs['target.name'] !== undefined) {
      this.TargetName = attribs['target.name']
    }

    if (attribs['api.name'] !== undefined) {
      this.ApiName = attribs['api.name']
    }

    if (attribs['source.name'] !== undefined) {
      this.SourceName = attribs['source.name']
    }
  }
}