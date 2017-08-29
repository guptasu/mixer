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
	"os"
	"path"
	"testing"
	"time"

	adp "istio.io/mixer/adapter"
	"istio.io/mixer/adapter/noop"
	"istio.io/mixer/pkg/adapter"
	adaptManager "istio.io/mixer/pkg/adapterManager"
	"istio.io/mixer/pkg/aspect"
	"istio.io/mixer/pkg/config"
	"istio.io/mixer/pkg/config/store"
	"istio.io/mixer/pkg/expr"
	"istio.io/mixer/pkg/il/evaluator"
	"istio.io/mixer/pkg/pool"
	mixerRuntime "istio.io/mixer/pkg/runtime"
	"istio.io/mixer/pkg/template"
	e2eTmpl "istio.io/mixer/test/e2e/template"
	"context"
	"istio.io/mixer/pkg/attribute"
)

const (
	globalCnfg = `
subject: namespace:ns
revision: "2022"
handlers:
  - name: fooHandler
    adapter: fakeHandler

manifests:
  - name: istio-proxy
    revision: "1"
    attributes:
      source.name:
        value_type: STRING
      target.name:
        value_type: STRING
      response.count:
        value_type: INT64
      attr.bool:
        value_type: BOOL
      attr.string:
        value_type: STRING
      attr.double:
        value_type: DOUBLE
      attr.int64:
        value_type: INT64
`
	srvConfig = `
subject: namespace:ns
action_rules:
- selector: target.name == "*"
  actions:
  - handler: fooHandler
    instances:
    - fooInstance

instances:
- name: fooInstance
  template: "report"
  params:
    value: "2"
    dimensions:
      source: source.name
      target_ip: target.name
`
)

func testConfigFlow(
	t *testing.T,
	declarativeSrvcCnfgFilePath string,
	declaredGlobalCnfgFilePath string,
	configStore2URL string,
) {
	// TODO replace
	useIL := false
	globalActualHandlerCallInfoToValidate = make(map[string]interface{})
	apiPoolSize := 1024
	adapterPoolSize := 1024
	configIdentityAttribute := "target.service"
	identityDomainAttribute := "svc.cluster.local"
	loopDelay := time.Second * 5
	singleThreadedGoRoutinePool := false
	configDefaultNamespace := "istio-config-default"
	gp := getGoRoutinePool(apiPoolSize, singleThreadedGoRoutinePool)
	adapterGP := getAdapterGoRoutinePool(adapterPoolSize, singleThreadedGoRoutinePool)

	adptInfos := []adapter.InfoFn{GetFakeHndlrBuilderInfo}
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

	adapters := []adapter.InfoFn{GetFakeHndlrBuilderInfo}
	adapterMap := adp.InventoryMap(adapters)
	store2, err := store.NewRegistry2(config.Store2Inventory()...).NewStore2(configStore2URL)
	if err != nil {
		t.Fatalf("Failed to connect to the configuration2 server. %v", err)
	}
	dispatcher, err = mixerRuntime.New(eval, gp, adapterGP,
		configIdentityAttribute, configDefaultNamespace,
		store2, adapterMap, e2eTmpl.SupportedTmplInfo,
	)
	if err != nil {
		t.Fatalf("Failed to create runtime dispatcher. %v", err)
	}

	adapterMgr := adaptManager.NewManager(
		[]adapter.RegisterFn{
			noop.Register,
		},
		aspect.Inventory(),
		eval,
		gp,
		adapterGP,
	)

	store, err := config.NewCompatFSStore(declaredGlobalCnfgFilePath, declarativeSrvcCnfgFilePath)
	if err != nil {
		t.Errorf("NewCompatFSStore failed: %v", err)
		return
	}

	configManager := config.NewManager(eval, adapterMgr.AspectValidatorFinder, adapterMgr.BuilderValidatorFinder, adptInfos,
		adapterMgr.SupportedKinds,
		template.NewRepository(e2eTmpl.SupportedTmplInfo),
		store,
		loopDelay,
		configIdentityAttribute, identityDomainAttribute)

	if useIL {
		configManager.Register(ilEval)
	}
	configManager.Register(adapterMgr)
	configManager.Start()

	// validate globalActualHandlerCallInfoToValidate
	if len(globalActualHandlerCallInfoToValidate) != 2 {
		t.Errorf("got call count %d\nwant %d", len(globalActualHandlerCallInfoToValidate), 2)
	}

	if globalActualHandlerCallInfoToValidate["ConfigureSampleReport"] == nil || globalActualHandlerCallInfoToValidate["Build"] == nil {
		t.Errorf("got call info as : %v. \nwant calls %s and %s to have been called", globalActualHandlerCallInfoToValidate, "ConfigureSample", "Build")
	}

	requestBag := attribute.GetMutableBag(nil)
	requestBag.Set(configIdentityAttribute, identityDomainAttribute)
	err = dispatcher.Report(context.TODO(), requestBag, )
	if err != nil {
		t.Errorf("Report call failed. %v", err)
	}
}
func getAdapterGoRoutinePool(adapterPoolSize int, singleThreadedGoRoutinePool bool) *pool.GoroutinePool {
	adapterGP := pool.NewGoroutinePool(adapterPoolSize, singleThreadedGoRoutinePool)
	adapterGP.AddWorkers(adapterPoolSize)
	defer adapterGP.Close()
	return adapterGP
}
func getGoRoutinePool(apiPoolSize int, singleThreadedGoRoutinePool bool) *pool.GoroutinePool {
	gp := pool.NewGoroutinePool(apiPoolSize, singleThreadedGoRoutinePool)
	gp.AddWorkers(apiPoolSize)
	gp.AddWorkers(apiPoolSize)
	defer gp.Close()
	return gp
}

func TestConfigFlow(t *testing.T) {
	sc, gsc := getCnfgs(srvConfig, globalCnfg)
	testConfigFlow(t, sc.Name(), gsc.Name(), "fs://"+path.Dir(sc.Name()))
	_ = os.Remove(sc.Name())
	_ = os.Remove(gsc.Name())
}
