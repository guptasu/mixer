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
	"github.com/gogo/protobuf/types"
	"io/ioutil"
	"istio.io/mixer/pkg/adapter"
	pkgAdapter "istio.io/mixer/pkg/adapter"
	"istio.io/mixer/pkg/aspect"
	"istio.io/mixer/pkg/attribute"
	"istio.io/mixer/pkg/config"
	"istio.io/mixer/pkg/expr"
	"istio.io/mixer/pkg/pool"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
	"github.com/googleapis/googleapis/google/rpc"
)

const (
	minimalGlobalCnfgForMetricAspect = `
subject: namespace:ns
revision: "2022"
adapters:
  - name: noopMetricKindAdapter
    kind: metrics
    impl: noopMetricKindAdapter

manifests:
  - name: istio-proxy
    revision: "1"
    attributes:
    - name: source.name
      value_type: STRING
    - name: target.name
      value_type: STRING
    - name: response.code
      value_type: INT64
    - name: api.name
      value_type: STRING
    - name: api.method
      value_type: STRING

metrics:
  - name: request_count
    kind: COUNTER
    value: INT64
    description: request count by source, target, service, and code
    labels:
      source: 1 # STRING
      target: 1 # STRING
      service: 1 # STRING
      method: 1 # STRING
      response_code: 2 # INT64
  - name: request_latency
    kind: COUNTER
    value: DURATION
    description: request latency by source, target, and service
    labels:
      source: 1 # STRING
      target: 1 # STRING
      service: 1 # STRING
      method: 1 # STRING
      response_code: 2 # INT64
`
	srvcCnfgConstInitialSection = `
subject: namespace:ns
rules:
- selector: true
  aspects:
`
	srvcCnfgSimpleAspectFromat = `
  - kind: metrics
    adapter: noopMetricKindAdapter
    params:
      metrics:
      - descriptorName: request_count
        value: response.code | 100
        labels:
          source: source.name | "one"
          target: target.name | "one"
          service: api.name | "one"
          method: api.method | "one"
          response_code: response.code | 111
`

	srvcCnfgComplexAspectFromat = `
  - kind: metrics
    adapter: noopMetricKindAdapter
    params:
      metrics:
      - descriptorName: request_count
        value: response.code | 100
        labels:
          source: source.name | target.name | source.name | target.name | "one"
          target: source.name | target.name | source.name | "one"
          service: target.name | api.name | target.name | "one"
          method: api.method | target.name | api.method | target.name | "one"
          response_code: response.code | response.code | response.code | response.code | response.code | response.code | response.code | response.code | 111
`
)

func registerNoopAdapter(r adapter.Registrar) {
	r.RegisterMetricsBuilder(&fakeNoopAdapter{adapter.NewDefaultBuilder("noopMetricKindAdapter", "Publishes metrics", &types.Empty{})})
}
type fakeNoopAdapter struct {
	adapter.DefaultBuilder
}
func (f *fakeNoopAdapter) NewMetricsAspect(env adapter.Env, cfg adapter.Config, metrics map[string]*adapter.MetricDefinition) (adapter.MetricsAspect, error) {
	return &fakeNoopAdapter{}, nil
}
func (p *fakeNoopAdapter) Record(vals []adapter.Value) error {
	return nil
}
func (*fakeNoopAdapter) Close() error { return nil }

func generateYamlConfigs(srvcCnfgAspectFromat string, configRepeatCount int) (declarativeSrvcCnfgFilePath string, declaredGlobalCnfgFilePath string) {
	srvcCnfgFile, _ := ioutil.TempFile("", "TestReportWithJSServCnfg")
	declarativeSrvcCnfgFilePath = srvcCnfgFile.Name()

	globalCnfgFile, _ := ioutil.TempFile("", "TestReportWithJSGlobalcnfg")
	declaredGlobalCnfgFilePath = globalCnfgFile.Name()
	_, _ = globalCnfgFile.Write([]byte(minimalGlobalCnfgForMetricAspect))
	_ = globalCnfgFile.Close()

	var srvcCnfgBuffer bytes.Buffer
	srvcCnfgBuffer.WriteString(srvcCnfgConstInitialSection)
	for i := 0; i < configRepeatCount; i++ {
		srvcCnfgBuffer.WriteString(strings.Replace(srvcCnfgAspectFromat, "$s", strconv.FormatInt(int64(i), 10), 1))
	}
	_, _ = srvcCnfgFile.Write([]byte(srvcCnfgBuffer.String()))
	_ = srvcCnfgFile.Close()

	return declarativeSrvcCnfgFilePath, declaredGlobalCnfgFilePath
}

var rpcStatus google_rpc.Status
func benchmarkAdapterManagerDispatch(b *testing.B, declarativeSrvcCnfgFilePath string, declaredGlobalCnfgFilePath string) {
	apiPoolSize := 1
	adapterPoolSize := 1
	identityAttribute := "target.service"
	identityDomainAttribute := "svc.cluster.local"
	loopDelay := time.Second * time.Duration(1)

	gp := pool.NewGoroutinePool(apiPoolSize, true)
	gp.AddWorkers(apiPoolSize)
	gp.AddWorkers(apiPoolSize)
	defer gp.Close()

	adapterGP := pool.NewGoroutinePool(adapterPoolSize, true)
	adapterGP.AddWorkers(adapterPoolSize)
	defer adapterGP.Close()

	eval := expr.NewCEXLEvaluator()
	adapterMgr := NewManager([]pkgAdapter.RegisterFn{
		registerNoopAdapter,
	}, aspect.Inventory(), eval, gp, adapterGP)
	store, _ := config.NewCompatFSStore(declaredGlobalCnfgFilePath, declarativeSrvcCnfgFilePath)
	cnfgMgr := config.NewManager(eval, adapterMgr.AspectValidatorFinder, adapterMgr.BuilderValidatorFinder,
		adapterMgr.SupportedKinds, store,
		loopDelay,
		identityAttribute, identityDomainAttribute)

	cnfgMgr.Register(adapterMgr)
	cnfgMgr.Start()

	requestBag := attribute.GetMutableBag(nil)
	requestBag.Set(identityAttribute, identityDomainAttribute)
	configs, _ := adapterMgr.loadConfigs(requestBag, adapterMgr.reportKindSet, false, false)

	b.ResetTimer()
	var r google_rpc.Status
	for n := 0; n < b.N; n++ {
		r = adapterMgr.dispatchReport(context.Background(), configs, requestBag, attribute.GetMutableBag(nil))
	}
	rpcStatus = r
}

func BenchmarkOneSimpleAspect(b *testing.B) {
	sc, gsc := generateYamlConfigs(srvcCnfgSimpleAspectFromat, 1)
	defer os.Remove(sc)
	defer os.Remove(gsc)
	benchmarkAdapterManagerDispatch(b, sc, gsc)
}

func Benchmark50SimpleAspect(b *testing.B) {
	sc, gsc := generateYamlConfigs(srvcCnfgSimpleAspectFromat, 50)
	defer os.Remove(sc)
	defer os.Remove(gsc)
	benchmarkAdapterManagerDispatch(b, sc, gsc)
}

func BenchmarkOneComplexAspect(b *testing.B) {
	sc, gsc := generateYamlConfigs(srvcCnfgComplexAspectFromat, 1)
	defer os.Remove(sc)
	defer os.Remove(gsc)
	benchmarkAdapterManagerDispatch(b, sc, gsc)
}

func Benchmark50ComplexAspect(b *testing.B) {
	sc, gsc := generateYamlConfigs(srvcCnfgComplexAspectFromat, 50)
	defer os.Remove(sc)
	defer os.Remove(gsc)
	benchmarkAdapterManagerDispatch(b, sc, gsc)
}
