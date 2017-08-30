// Copyright 2016 Istio Authors
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

package e2e

import (
	"testing"
	//"time"

	adp "istio.io/mixer/adapter"
	//"istio.io/mixer/adapter/noop"
	"istio.io/mixer/pkg/adapter"
	//adaptManager "istio.io/mixer/pkg/adapterManager"
	//"istio.io/mixer/pkg/aspect"
	"istio.io/mixer/pkg/config"
	"istio.io/mixer/pkg/config/store"
	"istio.io/mixer/pkg/expr"
	"istio.io/mixer/pkg/il/evaluator"
	"istio.io/mixer/pkg/pool"
	mixerRuntime "istio.io/mixer/pkg/runtime"
	//"istio.io/mixer/pkg/template"
	"os"
	"path"
	"istio.io/mixer/pkg/template"
	"istio.io/mixer/pkg/attribute"
	"reflect"
	"fmt"
	"github.com/davecgh/go-spew/spew"
)

const (
	configIdentityAttribute = "target.service"
	identityDomainAttribute = "svc.cluster.local"
)

// fail fatal if dispatcher cannot be constructed
func getDispatcher(t *testing.T, configStore2URL string, adptInfos []adapter.InfoFn, tmplInfos map[string]template.Info) mixerRuntime.Dispatcher {
	// TODO replace
	useIL := false
	apiPoolSize := 1024
	adapterPoolSize := 1024
	//loopDelay := time.Second * 5
	singleThreadedGoRoutinePool := false
	configDefaultNamespace := "istio-config-default"
	gp := getGoRoutinePool(apiPoolSize, singleThreadedGoRoutinePool)
	adapterGP := getAdapterGoRoutinePool(adapterPoolSize, singleThreadedGoRoutinePool)
	adapterMap := adp.InventoryMap(adptInfos)
	eval, err := expr.NewCEXLEvaluator(expr.DefaultCacheSize)
	if err != nil {
		t.Errorf("Failed to create expression evaluator: %v", err)
	}
	var ilEval *evaluator.IL
	if useIL {
		ilEval, err = evaluator.NewILEvaluator(expr.DefaultCacheSize)
		if err != nil {
			t.Fatalf("Failed to create IL expression evaluator with cache size %d: %v", 1024, err)
		}
		eval = ilEval
	}
	var dispatcher mixerRuntime.Dispatcher

	store2, err := store.NewRegistry2(config.Store2Inventory()...).NewStore2(configStore2URL)
	if err != nil {
		t.Fatalf("Failed to connect to the configuration2 server. %v", err)
	}
	dispatcher, err = mixerRuntime.New(eval, gp, adapterGP,
		configIdentityAttribute, configDefaultNamespace,
		store2, adapterMap, tmplInfos,
	)
	if err != nil {
		t.Fatalf("Failed to create runtime dispatcher. %v", err)
	}

	return dispatcher
}

func getAdapterGoRoutinePool(adapterPoolSize int, singleThreadedGoRoutinePool bool) *pool.GoroutinePool {
	adapterGP := pool.NewGoroutinePool(adapterPoolSize, singleThreadedGoRoutinePool)
	adapterGP.AddWorkers(adapterPoolSize)
	return adapterGP
}
func getGoRoutinePool(apiPoolSize int, singleThreadedGoRoutinePool bool) *pool.GoroutinePool {
	gp := pool.NewGoroutinePool(apiPoolSize, singleThreadedGoRoutinePool)
	gp.AddWorkers(apiPoolSize)
	gp.AddWorkers(apiPoolSize)
	return gp
}

func getCnfgs(srvcCnfg, attrCnfg string) (dir string) {
	tmpDir := path.Join(os.TempDir(), "e2eStoreDir")
	os.MkdirAll(tmpDir, os.ModePerm)

	srvcCnfgFile, _ := os.Create(path.Join(tmpDir, "srvc.yaml"))
	globalCnfgFile, _ := os.Create(path.Join(tmpDir, "global.yaml"))

	_, _ = globalCnfgFile.Write([]byte(attrCnfg))
	_, _ = srvcCnfgFile.Write([]byte(srvcCnfg))

	_ = globalCnfgFile.Close()
	_ = srvcCnfgFile.Close()

	return tmpDir
}

// return adapterInfoFns + corresponding SkyAdapter object.
func cnstrAdapterInfos(adptBehaviors []AdptBehavior) ([]adapter.InfoFn, []*spyAdapter) {
	var adapterInfos []adapter.InfoFn = make([]adapter.InfoFn, 0)
	var spyAdapters []*spyAdapter = make([]*spyAdapter, 0)
	for _, b := range adptBehaviors {
		sa := newSpyAdapter(b)
		spyAdapters = append(spyAdapters, sa)
		adapterInfos = append(adapterInfos, sa.getAdptInfoFn())
	}
	return adapterInfos, spyAdapters
}

func getAttrBag(attribs map[string]interface{}) *attribute.MutableBag {
	requestBag := attribute.GetMutableBag(nil)
	requestBag.Set(configIdentityAttribute, identityDomainAttribute)
	for k, v := range attribs {
		requestBag.Set(k, v)
	}
	return requestBag
}


func cmpAndErr(msg string, t *testing.T, expt interface{}, actual interface{}) {
	a := InterfaceSlice(expt)
	b := InterfaceSlice(actual)
	if len(a) != len(b) {
		t.Errorf(fmt.Sprintf("Not equal -> %s.\nActual :\n%s\n\nExpected :\n%s",msg, spew.Sdump(actual), spew.Sdump(expt)))
		return
	}

	for _, x1 := range a {
		f := false
		for _, x2 := range b {
			if reflect.DeepEqual(x1, x2) {
				f = true
			}
		}
		if !f {
			t.Errorf(fmt.Sprintf("Not equal -> %s.\nActual :\n%s\n\nExpected :\n%s",msg, spew.Sdump(actual), spew.Sdump(expt)))
			return
		}
	}
	return
}


func InterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)

	ret := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}
