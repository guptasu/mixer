package mymetric

import (
	"istio.io/mixer/pkg/adapter/config"
	mymetricproto "istio.io/mixer/pkg/templates/mymetric/generated/config"
)

type MyMetricProcessor interface {
	config.Handler
	ConfigureMyMetric(typeParams map[string]mymetricproto.Type)
	ProcessMyMetric(metricInstances []Instance)
}

type Instance struct {
	TypeName  string
	Value      interface{}
	Dimensions map[string]interface{}
}
