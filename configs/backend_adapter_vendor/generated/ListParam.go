package listtemplate

/*

This is the struct that users will write within the constructor

In the below constructor, everything under the params should match the ListParam struct

constructors:
- name: MyListCheckerConstructor
  type: global/types/GenericListCheckerType
  params:
    blacklist: true
    checkExpression: source.ip

*/
type ListParam struct {
  blacklist bool
  checkExpression string
}
