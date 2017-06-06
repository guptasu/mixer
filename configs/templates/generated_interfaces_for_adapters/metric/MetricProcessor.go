package metric

import "istio.io/mixer/configs/templates/metric"

type Instance struct {
  Value      interface{}
  Dimensions map[string]interface{}
}

type MetricProcessor interface {
  ConfigureMetric(types map[string]*istio_mixer_adapter_metric.Type /*typeName to Type mapping*/) error
  ReportMetric(instances map[string]*Instance /*typeName to Instance (generated from Constructor) mapping*/) (error)
}

