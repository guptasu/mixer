package noopMetricKindAdapter

import (
	"github.com/gogo/protobuf/types"

	"istio.io/mixer/pkg/adapter"
)

type (
	factory struct {
		adapter.DefaultBuilder
	}
	metricAdapter struct {
	}
)

var (
	name = "noopMetricKindAdapter"
	desc = "Publishes metrics"
)

func Register(r adapter.Registrar) {
	r.RegisterMetricsBuilder(&factory{adapter.NewDefaultBuilder(name, desc, &types.Empty{})})
}

func (f *factory) NewMetricsAspect(env adapter.Env, cfg adapter.Config, metrics map[string]*adapter.MetricDefinition) (adapter.MetricsAspect, error) {
	return &metricAdapter{}, nil
}

func (p *metricAdapter) Record(vals []adapter.Value) error {
	return nil
}

func (*metricAdapter) Close() error { return nil }
