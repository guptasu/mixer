// Copyright 2017 the Istio Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stackdriver

import (
	"context"

	multierror "github.com/hashicorp/go-multierror"

	"istio.io/mixer/adapter/stackdriver/config"
	"istio.io/mixer/adapter/stackdriver/log"
	sdmetric "istio.io/mixer/adapter/stackdriver/metric"
	"istio.io/mixer/pkg/adapter"
	"istio.io/mixer/template/logentry"
	"istio.io/mixer/template/metric"
)

type (
	builder struct {
		m metric.HandlerBuilder2
		l logentry.HandlerBuilder2
	}

	handler struct {
		m metric.Handler
		l logentry.Handler
	}
)

var (
	_ metric.HandlerBuilder = &obuilder{}
	_ metric.Handler        = &handler{}

	_ logentry.HandlerBuilder = &obuilder{}
	_ logentry.Handler        = &handler{}
)

// GetInfo returns the BuilderInfo associated with this adapter implementation.
func GetInfo() adapter.BuilderInfo {
	return adapter.BuilderInfo{
		Name:        "stackdriver",
		Impl:        "istio.io/mixer/adapte/stackdriver",
		Description: "Publishes StackDriver metrics and logs.",
		SupportedTemplates: []string{
			metric.TemplateName,
			logentry.TemplateName,
		},
		CreateHandlerBuilder: func() adapter.HandlerBuilder {
			return &obuilder{&builder{m: sdmetric.NewBuilder(), l: log.NewBuilder()}}
		},
		DefaultConfig:  &config.Params{},
		ValidateConfig: func(msg adapter.Config) *adapter.ConfigErrors { return nil },
		NewBuilder:     func() adapter.Builder2 { return &builder{m: sdmetric.NewBuilder(), l: log.NewBuilder()} },
	}
}

func (b *builder) SetMetricTypes(metrics map[string]*metric.Type) {
	b.m.SetMetricTypes(metrics)
}

func (b *builder) SetLogEntryTypes(entries map[string]*logentry.Type) {
	b.l.SetLogEntryTypes(entries)
}
func (b *builder) SetAdapterConfig(c adapter.Config) {
	b.m.SetAdapterConfig(c)
	b.l.SetAdapterConfig(c)
}

func (b *builder) Validate() (ce *adapter.ConfigErrors) {
	mce := b.m.Validate()
	lce := b.l.Validate()

	ce = ce.Extend(mce)
	ce = ce.Extend(lce)
	return
}

// Build creates a stack driver handler object.
func (b *builder) Build(ctx context.Context, env adapter.Env) (adapter.Handler, error) {
	m, err := b.m.Build(ctx, env)
	if err != nil {
		return nil, err
	}
	mh, _ := m.(metric.Handler)

	l, err := b.l.Build(ctx, env)
	if err != nil {
		return nil, err
	}
	lh, _ := l.(logentry.Handler)

	return &handler{m: mh, l: lh}, nil
}

func (h *handler) Close() error {
	return multierror.Append(h.m.Close(), h.l.Close()).ErrorOrNil()
}

func (h *handler) HandleMetric(ctx context.Context, values []*metric.Instance) error {
	return h.m.HandleMetric(ctx, values)
}

func (h *handler) HandleLogEntry(ctx context.Context, values []*logentry.Instance) error {
	return h.l.HandleLogEntry(ctx, values)
}

type obuilder struct {
	b *builder
}

func (o *obuilder) SetMetricTypes(metrics map[string]*metric.Type) error {
	o.b.SetMetricTypes(metrics)
	return nil
}

func (o *obuilder) SetLogEntryTypes(entries map[string]*logentry.Type) error {
	o.b.SetLogEntryTypes(entries)
	return nil
}

// Build creates a stack driver handler object.
func (o *obuilder) Build(cfg adapter.Config, env adapter.Env) (adapter.Handler, error) {
	o.b.SetAdapterConfig(cfg)
	return o.b.Build(context.Background(), env)
}
