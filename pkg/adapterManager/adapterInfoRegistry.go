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

package adapterManager

import (
	"fmt"
	"strings"

	"github.com/golang/glog"

	"istio.io/mixer/pkg/adapter"
	"istio.io/mixer/pkg/adapter/config"
)

type adapterInfoRegistry struct {
	adapterInfosByName map[string]*adapter.BuilderInfo
}

type handlerBuilderValidator func(hndlrBuilder config.HandlerBuilder, t adapter.SupportedTemplates,
	handlerName string) (bool, string)

// newRegistry returns a new AdapterInfo registry.
func newRegistry2(adapterInfos []adapter.GetAdapterInfoFn, hndlrBldrValidator handlerBuilderValidator) *adapterInfoRegistry {
	r := &adapterInfoRegistry{make(map[string]*adapter.BuilderInfo)}
	for idx, builder := range adapterInfos {
		glog.V(3).Infof("Registering [%d] %#v", idx, builder)
		adptInfo := builder()
		if a, ok := r.adapterInfosByName[adptInfo.Name]; ok {
			// panic only if 2 different AdapterInfo objects are trying to identify by the
			// same Name.
			msg := fmt.Errorf("duplicate registration for '%s' : old = %v new = %v", a.Name, adptInfo, a)
			glog.Error(msg)
			panic(msg)
		} else {

			if ok, errMsg := doesBuilderSupportsTemplates(adptInfo, hndlrBldrValidator); !ok {
				// panic if an Adapter's HandlerBuilder does not implement interfaces that it says it wants to support.
				glog.Error(errMsg)
				panic(errMsg)
			}

			r.adapterInfosByName[adptInfo.Name] = &adptInfo
		}
	}
	return r
}

// AdapterInfoMap returns the known AdapterInfos, indexed by their names.
func AdapterInfoMap(handlerRegFns []adapter.GetAdapterInfoFn,
	hndlrBldrValidator handlerBuilderValidator) map[string]*adapter.BuilderInfo {
	return newRegistry2(handlerRegFns, hndlrBldrValidator).adapterInfosByName
}

// FindAdapterInfo returns the AdapterInfo object with the given name.
func (r *adapterInfoRegistry) FindAdapterInfo(name string) (b *adapter.BuilderInfo, found bool) {
	bi, found := r.adapterInfosByName[name]
	if !found {
		return nil, false
	}
	return bi, true
}

func doesBuilderSupportsTemplates(info adapter.BuilderInfo, hndlrBldrValidator handlerBuilderValidator) (bool, string) {
	handlerBuilder := info.CreateHandlerBuilderFn()
	resultMsgs := make([]string, 0)
	for _, t := range info.SupportedTemplates {
		if ok, errMsg := hndlrBldrValidator(handlerBuilder, t, info.Name); !ok {
			resultMsgs = append(resultMsgs, errMsg)
		}
	}
	if len(resultMsgs) != 0 {
		return false, strings.Join(resultMsgs, "\n")
	}
	return true, ""
}
