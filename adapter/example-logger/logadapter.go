import "istio.io/autogen/mixer/aspect/log"

func Register(r Registrar) {
    r.registerLog(&adapter{})
}

struct adapter {}

func (adapter) Close() error { return nil }
func (adapter)  NewProcessor(Config, []*Template) (Processor, error) {
    return &processor{}, nil
}

struct processor {
    // state here
}

func (processor) Close() error { return nil }

func (p processor) Process([]*Instance) error {
    // do work logging the instances
    return nil
}