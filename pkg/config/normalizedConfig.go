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

package config

import (
	"github.com/robertkrimen/otto"
	"istio.io/mixer/pkg/attribute"
	configpb "istio.io/mixer/pkg/config/proto"
	"fmt"
)

type NormalizedConfig struct {
	JavaScript string
}

func (n *NormalizedConfig) Evalaute(requestBag *attribute.MutableBag,
	aspectFinder func(cfgs []*configpb.Combined, kind string, adapterName string, val interface{}) (*configpb.Combined, error),
	callBack func(kind string, adapterImplName string, val interface{})) {
	vm := otto.New()
	vm.Set("CallBackFromUserScript_go", callBack)
	vm.Run(n.JavaScript)
	checkFn, _ := vm.Get("report")
	_, errFromJS := checkFn.Call(otto.NullValue(), requestBag)
	if errFromJS != nil {
		fmt.Println(errFromJS)
	}
}
