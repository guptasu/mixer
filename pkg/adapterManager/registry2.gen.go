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

// !!!!!!!!!!!!!!!!!!!!! WARNING !!!!!!!!!!!!!!!!!!!!!!!!
// THIS IS AUTO GENERATED FILE - SIMULATED - HAND WRITTEN

package adapterManager

import (
	"fmt"

	"github.com/golang/glog"

	"istio.io/mixer/pkg/adapter"
	adpCnfg "istio.io/mixer/pkg/adapter/config"
	"istio.io/mixer/pkg/templates/mymetric/generated"
)


// registry implements pkg/adapter/Registrar2.
// registry is initialized in the constructor and is immutable thereafter.
// All registered builders must have unique names per aspect kind.
// It also implements builders that manager uses.
type registry2 struct {
	handlersByName map[string]*adpCnfg.Handler
}

// newRegistry returns a new Builder registry.
func newRegistry2(builders []adapter.RegisterFn2) *registry2 {
	r := &registry2{make(map[string]*adpCnfg.Handler)}
	for idx, builder := range builders {
		glog.V(3).Infof("Registering [%d] %#v", idx, builder)
		builder(r)
	}

	// ensure interfaces are satisfied.
	// should be compiled out.
	var _ adapter.Registrar2 = r
	var _ builderFinder2 = r
	return r
}

func (r *registry2) FindHandler(name string) (b adpCnfg.Handler, found bool) {
	if bi, found := r.handlersByName[name]; !found {
		return nil, false
	} else {
		return *bi, true
	}
}

func (r *registry2) insertHandler(b adpCnfg.Handler) {
	bi := r.handlersByName[b.Name()]
	if bi == nil {
		bi = &b
		r.handlersByName[b.Name()] = bi
	} else if *bi != b {
		// panic only if 2 different handler objects are trying to identify by the
		// same Name.  2nd registration is ok so long as old and the new are same
		msg := fmt.Errorf("duplicate registration for '%s' : old = %v new = %v", b.Name(), bi, b)
		glog.Error(msg)
		panic(msg)
	}


	if glog.V(1) {
		glog.Infof("Registered %s", b.Name())
	}
}

/////////// ALL THE BELOW CODE IS GENERATED FROM TEMPLATES //////////////////

func (r *registry2) RegisterMyMetricProcessor(b mymetric.MyMetricProcessor) {
	r.insertHandler(b)
}
