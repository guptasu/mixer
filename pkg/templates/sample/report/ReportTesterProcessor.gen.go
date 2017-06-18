package istio_mixer_adapter_sample_report

import "istio.io/mixer/pkg/adapter/config"

type Instance struct {
	name       string
	Value      interface{}
	Dimensions map[string]interface{}
}

type SampleProcessor interface {
	config.Handler
	ConfigureSample(map[string]*Type /*Constructor:instance_name to Type mapping. Note type name will not be passed at all*/) error
	ReportSample([]*Instance /*The type is inferred from the Instance.name and the mapping of instance to types passed during the config time*/) error
}
