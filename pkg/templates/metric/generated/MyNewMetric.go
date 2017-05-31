package mymetric

import mymetric "istio.io/mixer/pkg/templates/metric/generated/config"

type MetricProcessor interface {
	ConfigureMetric(typeParams map[string]mymetric.Type)
	ProcessMetric(metricInstances []mymetric.Instance)
}

