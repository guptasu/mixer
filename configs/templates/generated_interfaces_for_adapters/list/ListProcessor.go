package foo_bar_mylistchecker

type Instance struct {
  name            string
  CheckExpression interface{}
}

type ListProcessor interface {
  ConfigureList(types map[string]*Type /*Constructor:instance_name to Type mapping. Note type name will not be passed at all*/) error
  CheckList(instances []*Instance /*The type is inferred from the Instance.name and the mapping of instance to types passed during the config time*/) (bool, error)
}
