package mymetric

import (
	mymetric "istio.io/mixer/pkg/templates/metric/generated/config"
	"io"
)

// TODO Does not belong here.
type Handler interface {
	io.Closer

	// Name returns the official name of the aspects produced by this builder.
	Name() string

	// Description returns a user-friendly description of the aspects produced by this builder.
	Description() string
}


type MetricProcessor interface {
	Handler
	ConfigureMetric(typeParams map[string]mymetric.Type)
	ProcessMetric(metricInstances []mymetric.Instance)
}

