package mymetric

import (
	mymetric "istio.io/mixer/pkg/templates/metric/generated/config"
	"istio.io/mixer/pkg/adapter/config"
)


type MetricProcessor interface {
	config.Handler
	ConfigureMetric(typeParams map[string]mymetric.Type)
	ProcessMetric(metricInstances []mymetric.Instance)
}

