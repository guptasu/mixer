// Copyright 2017 Istio Authors.
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

package adapter

import (
	"istio.io/mixer/pkg/adapter/config"
)

// BuilderInfo describes the Adapter and provides a function to a Handler Builder method.
type BuilderInfo struct {
	// Name returns the official name of the adapter.
	Name string
	// Description returns a user-friendly description of the adapter.
	Description string
	// CreateHandlerBuilderFn is a function that creates a HandlerBuilder which implements Builders associated
	// with the SupportedTemplates.
	CreateHandlerBuilderFn CreateHandlerBuilder
	// SupportedTemplates expressess all the templates the Adapter wants to serve.
	SupportedTemplates []SupportedTemplates
}

// CreateHandlerBuilder is a function that creates a HandlerBuilder.
type CreateHandlerBuilder func() config.HandlerBuilder

// GetAdapterInfoFn returns an AdapterInfo object that Mixer will use to create HandlerBuilder
type GetAdapterInfoFn func() BuilderInfo
