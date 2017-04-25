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
	"context"
	"io/ioutil"
	//"os"
	"bytes"
	pkgAdapter "istio.io/mixer/pkg/adapter"
	"istio.io/mixer/pkg/adapterManager/noopMetricKindAdapter"
	"istio.io/mixer/pkg/aspect"
	"istio.io/mixer/pkg/attribute"
	"istio.io/mixer/pkg/cnfgNormalizer"
	"istio.io/mixer/pkg/config"
	"istio.io/mixer/pkg/expr"
	"istio.io/mixer/pkg/pool"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	globalCnfgForMetricAspect = `
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
      value_type: 1 # STRING
    - name: target.name
      value_type: 1 # STRING
    - name: response.code
      value_type: 2 # INT64
    - name: api.name
      value_type: 1 # STRING
    - name: api.method
      value_type: 1 # STRING
metrics:
  - name: request_count
    kind: 2 # COUNTER
    value: 2 # INT64
    description: request count by source, target, service, and code
    labels:
      source: 1 # STRING
      target: 1 # STRING
      service: 1 # STRING
      method: 1 # STRING
      response_code: 2 # INT64
  - name: request_latency
    kind: 2 # COUNTER
    value: 10 # DURATION
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
	srvcCnfgYamlSimpleAspectStrFromat = `
  - name: prometheus_reporting_all_metrics$s
    kind: metrics
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

	srvcCnfgYamlComplexAspectStrFromat = `
  - name: prometheus_reporting_all_metrics$s
    kind: metrics
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

func benchmarkDispatchSingleHugeAspect(b *testing.B, aspectStringFmt string, loopSize int) {
	srvcCnfgFile, _ := ioutil.TempFile("", "TestReportWithJSServCnfg")
	srvcCnfgFilePath := srvcCnfgFile.Name()

	globalCnfgFile, _ := ioutil.TempFile("", "TestReportWithJSGlobalcnfg")
	globalCnfgFilePath := globalCnfgFile.Name()
	_, _ = globalCnfgFile.Write([]byte(globalCnfgForMetricAspect))
	_ = globalCnfgFile.Close()
	//defer os.Remove(globalCnfgFilePath)

	var srvcCnfgBuffer bytes.Buffer
	srvcCnfgBuffer.WriteString(srvcCnfgConstInitialSection)
	for i := 0; i < loopSize; i++ {
		srvcCnfgBuffer.WriteString(strings.Replace(aspectStringFmt, "$s", strconv.FormatInt(int64(i), 10), 1))
	}
	_, _ = srvcCnfgFile.Write([]byte(srvcCnfgBuffer.String()))
	_ = srvcCnfgFile.Close()
	//defer os.Remove(srvcCnfgFilePath)

	apiPoolSize := 1
	adapterPoolSize := 1
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
		noopMetricKindAdapter.Register,
	}, aspect.Inventory(), eval, gp, adapterGP)

	cnfgMgr := config.NewManager(eval, adapterMgr.AspectValidatorFinder, adapterMgr.BuilderValidatorFinder,
		adapterMgr.SupportedKinds,
		globalCnfgFilePath,
		srvcCnfgFilePath,
		loopDelay,
		"",
		cnfgNormalizer.NormalizedJavascriptConfigNormalizer{})

	cnfgMgr.Register(adapterMgr)
	cnfgMgr.Start()

	requestBag := attribute.GetMutableBag(nil)
	configs, _ := adapterMgr.loadConfigs(requestBag, adapterMgr.reportKindSet, false)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = adapterMgr.executeDispatch(context.Background(), configs, requestBag, attribute.GetMutableBag(nil))
	}
}

/* JS benchmark
BenchmarkDispatchSimpleOneAspect-12     	   20000	     74330 ns/op
BenchmarkDispatchSimple50Aspect-12      	     500	   2546216 ns/op
BenchmarkDispatchComplexOneAspect-12    	   20000	     92504 ns/op
BenchmarkDispatchComplex50Aspect-12     	     500	   3305152 ns/op
*/
/*
/* Without JS code
BenchmarkDispatchSimpleOneAspect-12     	   20000	     85643 ns/op
BenchmarkDispatchSimple50Aspect-12      	     300	   4180791 ns/op
BenchmarkDispatchComplexOneAspect-12    	   10000	    121990 ns/op
BenchmarkDispatchComplex50Aspect-12     	     200	   5866466 ns/op
*/

func BenchmarkDispatchSimpleOneAspect(b *testing.B) {
	benchmarkDispatchSingleHugeAspect(b, srvcCnfgYamlSimpleAspectStrFromat, 1)
}

func BenchmarkDispatchSimple50Aspect(b *testing.B) {
	benchmarkDispatchSingleHugeAspect(b, srvcCnfgYamlSimpleAspectStrFromat, 50)
}

func BenchmarkDispatchComplexOneAspect(b *testing.B) {
	benchmarkDispatchSingleHugeAspect(b, srvcCnfgYamlComplexAspectStrFromat, 1)
}

func BenchmarkDispatchComplex50Aspect(b *testing.B) {
	benchmarkDispatchSingleHugeAspect(b, srvcCnfgYamlComplexAspectStrFromat, 50)
}
