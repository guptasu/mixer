package prometheus

import (
	"fmt"

	"github.com/douglas-reid/mixer-config-experiments/metric"
)

type Adapter struct{}
type Processor struct{}

// this will likely need refinement in approach
func RegisterAdapter() metric.Adapter {
	return newAdapter()
}

// this is where adapter config would go
func newAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) NewProcessor(config metric.Config, templates []*metric.Template) (metric.Processor, error) {
	// do preprocessing / registration of templates here maybe
	return &Processor{}, nil
}

func (Adapter) Close() error { return nil }

func (p *Processor) Process(instances []*metric.Instance) error {
	for _, instance := range instances {
		// record values in appropriate metrics
		fmt.Println(instance.Value)
	}
	return nil
}

func (Processor) Close() error { return nil }
