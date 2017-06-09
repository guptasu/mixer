package log

import "istio.io/mixer/configs/templates/log"

type Instance struct {
  name       string
  Dimensions map[string]interface{}
}

type LogProcessor interface {
  ConfigureLog(types map[string]*istio_mixer_adapter_log.Type /*Constructor:instance_name to Type mapping. Note type name will not be passed at all*/) error
  ReportLog(instances []*Instance /*The type is inferred from the Instance.name and the mapping of instance to types passed during the config time*/) (error)
}
