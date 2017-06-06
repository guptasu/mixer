package quota

import (
  "istio.io/mixer/configs/templates/quota"
  "istio.io/mixer/bazel-mixer/configs/templates/quota"
)

type Instance struct {
  Dimensions map[string]interface{}
}

type QuotaProcessor interface {
  ConfigureList(types map[string]*istio_mixer_adapter_quota.Type /*typeName to Type mapping*/) error
  CheckList(instances map[string]*Instance /*typeName to Instance (generated from Constructor) mapping*/) (bool, error)
}
