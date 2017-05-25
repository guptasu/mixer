package listtemplate

/*

This is the struct that will be passed to the adapters.

This struct is constructed based on the evaluated value from a Constructor.

For example:

constructors:
- name: MyListCheckerConstructor
  type: global/types/GenericListCheckerType
  params:
    blacklist: true
    checkExpression: source.ip

types:
- name: global/types/GenericListCheckerType
  template: global/template/ListTemplate
  params:

The resulting ListIntance would look like:

SampleListInstance = ListIntance {
  // consider source.ip evaluated to 'foo.bar.com' during request time.
  Blacklist = true
  CheckExpression = "foo.bar.com"
  Template = &MyListTemplate
}

MyListTemplate = ListTemplate {
}

*/
type ListIntance struct {
  Template *ListTemplate
  Blacklist bool
  CheckExpression string
}

