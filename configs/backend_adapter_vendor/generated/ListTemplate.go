package listtemplate

/*

This is the struct that users will write within the type

In the below types, everything under the params should match the ListTemplate struct

Example:

types:
- name: global/types/GenericListCheckerType
  template: global/template/ListTemplate
  params:

*/
type ListTemplate struct {
  // This is empty since for global/template/ListTemplate there are :
  // - no ValueType expressed_in_constructor fields.
  // - no fields other than expressed_in_constructor annotated.
  //
  // If any of the above were true, this struct would be non empty.
}
