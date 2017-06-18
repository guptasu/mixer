package istio_mixer_adapter_metric

import "istio.io/mixer/pkg/adapter/config"

type Instance struct {
	name       string
	Value      interface{}
	Dimensions map[string]interface{}
}

type MetricProcessor interface {
	config.Handler
	ConfigureMetric(map[string]*Type /*Constructor:instance_name to Type mapping. Note type name will not be passed at all*/) error
	ReportMetric([]*Instance /*The type is inferred from the Instance.name and the mapping of instance to types passed during the config time*/) error
}
