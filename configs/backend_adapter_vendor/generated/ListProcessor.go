package listtemplate

import istio_mixer_v1_config_descriptor "istio.io/api/mixer/v1/config/descriptor"

/*

This is the interface that adapters need to implement if they want to do check call based on the ListTemplate input.

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

The resulting ListTemplate would look like:
MyListTemplate = ListTemplate {
  blacklist = true
  checkExpression = STRING
}

Those are the two objects needed to invoke the adapters.
*/

type ListProcessor interface {
  Configure(templates []*ListTemplate) error
  Process(instances []*ListIntance) error
}
