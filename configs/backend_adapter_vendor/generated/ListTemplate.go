package listtemplate

import istio_mixer_v1_config_descriptor "istio.io/api/mixer/v1/config/descriptor"

/*

This is the struct that users will write within the type

In the below types, everything under the params should match the ListTemplate struct

Example : two types defined:

types:
- name: MyBlacklistCheckerType
  template: ListTemplate
  params:
    blacklist: true
    checkExpression: STRING

- name: MyWhitelistCheckerType
  template: ListTemplate
  params:
    blacklist: false
    checkExpression: STRING

*/
type ListTemplate struct {
  blacklist bool
  checkExpression istio_mixer_v1_config_descriptor.ValueType
}
