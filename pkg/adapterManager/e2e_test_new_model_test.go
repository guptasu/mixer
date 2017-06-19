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

package adapterManager

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"testing"
	"time"

	google_rpc "github.com/googleapis/googleapis/google/rpc"
	//adapter "istio.io/mixer/adapter"
	"istio.io/mixer/adapter/noop"
	"istio.io/mixer/adapter/noop2"
	pkgAdapter "istio.io/mixer/pkg/adapter"
	"istio.io/mixer/pkg/aspect"
	"istio.io/mixer/pkg/attribute"
	"istio.io/mixer/pkg/config"
	"istio.io/mixer/pkg/expr"
	"istio.io/mixer/pkg/pool"
)

const (
	globalCnfg = `
subject: namespace:ns
revision: "2022"
adapters:
  - name: no-op
    kind: metrics
    impl: no-op

manifests:
  - name: istio-proxy
    revision: "1"
    attributes:
      source.name:
        value_type: STRING
      target.name:
        value_type: STRING
      response.code:
        value_type: INT64
      api.name:
        value_type: STRING
      api.method:
        value_type: STRING

handlers:
  - name: noop2
    adapter: noop2
`

	srvcCnfg = `
subject: namespace:ns
constructors:
  - instanceName: MyMetricConstructor
    templateName: myMetricTypeReqCount
    params:
      value: response.code
      dimensions:
        source: source.name
        target: source.name

action_rules:
- selector: target.service == "*"
  actions:
  - handler: mygRPCAdapter
    instances:
    - MyMetricConstructor
`
)

func createCnfgs() (declarativeSrvcCnfg *os.File, declaredGlobalCnfg *os.File) {
	srvcCnfgFile, _ := ioutil.TempFile("", "managerDispatchBenchmarkTest")
	globalCnfgFile, _ := ioutil.TempFile("", "managerDispatchBenchmarkTest")

	_, _ = globalCnfgFile.Write([]byte(globalCnfg))
	_ = globalCnfgFile.Close()

	var srvcCnfgBuffer bytes.Buffer
	srvcCnfgBuffer.WriteString(srvcCnfg)

	_, _ = srvcCnfgFile.Write([]byte(srvcCnfgBuffer.String()))
	_ = srvcCnfgFile.Close()

	return srvcCnfgFile, globalCnfgFile
}

func testEnd2EndMixer(t *testing.T, declarativeSrvcCnfgFilePath string, declaredGlobalCnfgFilePath string) {
	apiPoolSize := 1024
	adapterPoolSize := 1024
	identityAttribute := "target.service"
	identityDomainAttribute := "svc.cluster.local"
	loopDelay := time.Second * 5
	singleThreadedGoRoutinePool := false

	gp := pool.NewGoroutinePool(apiPoolSize, singleThreadedGoRoutinePool)
	gp.AddWorkers(apiPoolSize)
	gp.AddWorkers(apiPoolSize)
	defer gp.Close()

	adapterGP := pool.NewGoroutinePool(adapterPoolSize, singleThreadedGoRoutinePool)
	adapterGP.AddWorkers(adapterPoolSize)
	defer adapterGP.Close()

	eval, err := expr.NewCEXLEvaluator(expr.DefaultCacheSize)
	if err != nil {
		t.Errorf("Failed to create expression evaluator: %v", err)
	}
	store, err := config.NewCompatFSStore(declaredGlobalCnfgFilePath, declarativeSrvcCnfgFilePath)
	if err != nil {
		t.Errorf("NewCompatFSStore failed: %v", err)
		return
	}

	adapterMgr := NewManager([]pkgAdapter.RegisterFn{
		noop.Register,
	},
		// TODO use Inventory2 here, but getting log_dir flag error
		[]pkgAdapter.RegisterFn2{
			noop2.Register,
		}, aspect.Inventory(), eval, gp, adapterGP)

	cnfgMgr := config.NewManager(eval, adapterMgr.AspectValidatorFinder, adapterMgr.BuilderValidatorFinder,
		adapterMgr.HandlerFinder, adapterMgr.SupportedKinds, store,
		loopDelay,
		identityAttribute, identityDomainAttribute)
	cnfgMgr.Register(adapterMgr)
	cnfgMgr.Start()

	requestBag := attribute.GetMutableBag(nil)
	requestBag.Set(identityAttribute, identityDomainAttribute)
	requestBag.Set("response.code", 200)
	requestBag.Set("source.name", "mysourcename")
	configs, err := adapterMgr.loadConfigs(requestBag, adapterMgr.reportKindSet, false, false)
	if err != nil {
		t.Errorf("adapterMgr.loadConfigs failed: %v", err)
		return
	}

	var r google_rpc.Status

	r = adapterMgr.dispatchReport(context.Background(), configs, requestBag, attribute.GetMutableBag(nil))

	if r.Code != 0 {
		t.Errorf("dispatchReport benchmark test returned status code %d; expected 0", r.Code)
	}

}

func TestEnd2EndMixer(t *testing.T) {
	sc, gsc := createCnfgs()
	testEnd2EndMixer(t, sc.Name(), gsc.Name())
	_ = os.Remove(sc.Name())
	_ = os.Remove(gsc.Name())
}
