package listtemplate

import istio_mixer_v1_config_descriptor "istio.io/api/mixer/v1/config/descriptor"

/*

This is the struct that users will write within the constructor

In the below constructor, everything under the params should match the ListParam struct

constructors:
- name: MyListConstructor
  type: MyBlacklistCheckerType
  params:
    checkExpression: source.ip

*/
type ListParam struct {
  // string expression written in the yaml file.
  checkExpression string
}
