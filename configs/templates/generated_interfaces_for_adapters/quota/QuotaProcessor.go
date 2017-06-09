package quota

import (
  "istio.io/mixer/configs/templates/quota"
  "istio.io/mixer/bazel-mixer/configs/templates/quota"
)

type Instance struct {
  name       string
  Dimensions map[string]interface{}
}

type QuotaProcessor interface {
  ConfigureList(types map[string]*istio_mixer_adapter_quota.Type /*Constructor:instance_name to Type mapping. Note type name will not be passed at all*/) error
  CheckList(instances []*Instance /*The type is inferred from the Instance.name and the mapping of instance to types passed during the config time*/) (bool, error)
}
