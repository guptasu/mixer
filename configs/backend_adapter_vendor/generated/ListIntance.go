package listtemplate

import istio_mixer_v1_config_descriptor "istio.io/api/mixer/v1/config/descriptor"

/*

This is the struct that will be passed to the adapters.

This struct is constructed based on the evaluated value from a Constructor.

For example:

constructors:
- name: MyListConstructor
  type: MyBlacklistCheckerType
  params:
    checkExpression: source.ip

types:
- name: MyBlacklistCheckerType
  template: ListTemplate
  params:
    blacklist: true
    checkExpression: STRING


The resulting ListIntance would look like:

SampleListInstance = ListIntance {
  // consider source.ip evaluated to 'foo.bar.com' during request time.
  CheckExpression = "foo.bar.com"
  Template = &MyListTemplate
}

MyListTemplate = ListTemplate {
  blacklist = true
  checkExpression = STRING
}

*/
type ListIntance struct {
  Template *ListTemplate
  CheckExpression interface{}
}
