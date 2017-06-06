package log

import "istio.io/mixer/configs/templates/log"

type Instance struct {
  Dimensions map[string]interface{}
}

type LogProcessor interface {
  ConfigureLog(types map[string]*istio_mixer_adapter_log.Type /*typeName to Type mapping*/) error
  ReportLog(instances map[string]*Instance /*typeName to Instance (generated from Constructor) mapping*/) (error)
}
