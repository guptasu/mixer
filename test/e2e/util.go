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
	"os"
	"io/ioutil"
	"path"
	"bytes"
)

const (
	configIdentityAttribute = "target.service"
	identityDomainAttribute = "svc.cluster.local"
)

// fail fatal if dispatcher cannot be constructed
func getDispatcher(t *testing.T, configStore2URL string, declaredGlobalCnfgFilePath string, declarativeSrvcCnfgFilePath string, adptInfos []adapter.InfoFn) mixerRuntime.Dispatcher {
	// TODO replace
	useIL := false
	apiPoolSize := 1024
	adapterPoolSize := 1024
	loopDelay := time.Second * 5
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
		t.Fatalf("NewCompatFSStore failed: %v", err)
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
	return dispatcher
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

func getCnfgs(srvcCnfg, attrCnfg string) (declarativeSrvcCnfg *os.File, declaredGlobalCnfg *os.File) {
	dir2, _ := ioutil.TempDir("e2eStoreDir", "")
	srvcCnfgFile, _ := os.Create(path.Join(dir2, "srvc.yaml"))
	globalCnfgFile, _ := os.Create(path.Join(dir2, "global.yaml"))

	_, _ = globalCnfgFile.Write([]byte(attrCnfg))
	_ = globalCnfgFile.Close()

	var srvcCnfgBuffer bytes.Buffer
	srvcCnfgBuffer.WriteString(srvcCnfg)

	_, _ = srvcCnfgFile.Write([]byte(srvcCnfgBuffer.String()))
	_ = srvcCnfgFile.Close()

	return srvcCnfgFile, globalCnfgFile
}
