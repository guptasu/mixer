package mymetric

import (
	mymetric "istio.io/mixer/pkg/templates/mymetric/generated/config"
	"istio.io/mixer/pkg/adapter/config"
)


type MyMetricProcessor interface {
	config.Handler
	ConfigureMyMetric(typeParams map[string]mymetric.Type)
	ProcessMyMetric(metricInstances []mymetric.Instance)
}

