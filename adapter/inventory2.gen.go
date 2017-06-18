// Copyright 2017 Istio Authors
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
	"istio.io/mixer/pkg/adapter"
	"istio.io/mixer/adapter/noop2"
)

// Inventory2 returns the inventory of all available adapters.
func Inventory2() []adapter.RegisterFn2 {
	return []adapter.RegisterFn2{
		noop2.Register,
	}
}
