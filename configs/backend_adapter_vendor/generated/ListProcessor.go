package listtemplate

/*

This is the interface that adapters need to implement if they want to do check call based on the ListTemplate input.

For example:

constructors:
- name: MyListCheckerConstructor
  type: GenericListCheckerType
  params:
    blacklist: true
    checkExpression: source.ip

types:
- name: GenericListCheckerType
  template: ListTemplate
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

Those are the two objects needed to invoke the adapter's Configure and
Process calls.
*/

type ListProcessor interface {
  Configure(templates []*ListTemplate) error
  Process(instances []*ListIntance) error
}
