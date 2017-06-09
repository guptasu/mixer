package metric

import "istio.io/mixer/configs/templates/metric"

type Instance struct {
  name       string
  Value      interface{}
  Dimensions map[string]interface{}
}

type MetricProcessor interface {
  ConfigureMetric(types map[string]*istio_mixer_adapter_metric.Type /*Constructor:instance_name to Type mapping. Note type name will not be passed at all*/) error
  ReportMetric(instances []*Instance /*The type is inferred from the Instance.name and the mapping of instance to types passed during the config time*/) (error)
}

