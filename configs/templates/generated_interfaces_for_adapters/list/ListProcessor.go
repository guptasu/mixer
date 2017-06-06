package list

import "istio.io/mixer/configs/templates/list"

type Instance struct {
  Blacklist       bool
  CheckExpression interface{}
}

type ListProcessor interface {
  ConfigureList(types map[string]*foo_bar_mylistchecker.Type /*typeName to Type mapping*/) error
  CheckList(instances map[string]*Instance /*typeName to Instance (generated from Constructor) mapping*/) (bool, error)
}

