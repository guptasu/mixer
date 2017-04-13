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
	"istio.io/mixer/adapter"
	"istio.io/mixer/pkg/aspect"
	"istio.io/mixer/pkg/attribute"
	"istio.io/mixer/pkg/config"
	"istio.io/mixer/pkg/expr"
	"istio.io/mixer/pkg/pool"
	"testing"
	"time"
	"strings"
	//"fmt"
	"strconv"
)

const globalCnfg = `
subject: namespace:ns
revision: "2022"
adapters:
  - name: default
    kind: quotas
    impl: memQuota
    params:
  - name: default
    impl: stdioLogger
    params:
      logStream: 0 # STDERR
  - name: prometheus
    kind: metrics
    impl: prometheus
    params:
  - name: default
    impl: denyChecker
attributes:
  # TODO: we really need to remove these, they're not part of the attribute vocab.
  - name: api.name
    value_type: 1 # STRING
  - name: api.method
    value_type: 1 # STRING

  - name: source.name
    value_type: 1 # STRING
  - name: target.name
    value_type: 1 # STRING
  - name: origin.ip
    value_type: 6 # IP_ADDRESS
  - name: origin.user
    value_type: 1 # STRING
  - name: request.time
    value_type: 5 # TIMESTAMP
  - name: request.method
    value_type: 1 # STRING
  - name: request.path
    value_type: 1 # STRING
  - name: request.scheme
    value_type: 1 # STRING
  - name: response.size
    value_type: 2 # INT64
  - name: response.code
    value_type: 2 # INT64
  - name: response.latency
    value_type: 10 # DURATION
metrics:
  - name: request_count
    kind: 2 # COUNTER
    value: 2 # INT64
    description: request count by source, target, service, and code
    labels:
    - name: source
      value_type: 1 # STRING
    - name: target
      value_type: 1 # STRING
    - name: service
      value_type: 1 # STRING
    - name: method
      value_type: 1 # STRING
    - name: response_code
      value_type: 2 # INT64
  - name: request_latency
    kind: 2 # COUNTER
    value: 10 # DURATION
    description: request latency by source, target, and service
    labels:
    - name: source
      value_type: 1 # STRING
    - name: target
      value_type: 1 # STRING
    - name: service
      value_type: 1 # STRING
    - name: method
      value_type: 1 # STRING
    - name: response_code
      value_type: 2 # INT64
quotas:
  - name: RequestCount
    max_amount: 5
    expiration:
      seconds: 1
logs:
  - name: accesslog.common
    display_name: Apache Common Log Format
    log_template: '{{or (.originIp) "-"}} - {{or (.sourceUser) "-"}} [{{or (.timestamp.Format "02/Jan/2006:15:04:05 -0700") "-"}}] "{{or (.method) "-"}} {{or (.url) "-"}} {{or (.protocol) "-"}}" {{or (.responseCode) "-"}} {{or (.responseSize) "-"}}'
    labels:
    - name: originIp
      value_type: 6 # IP_ADDRESS
    - name: sourceUser
      value_type: 1 # STRING
    - name: timestamp
      value_type: 5 # TIMESTAMP
    - name: method
      value_type: 1 # STRING
    - name: url
      value_type: 1 # STRING
    - name: protocol
      value_type: 1 # STRING
    - name: responseCode
      value_type: 2 # INT64
    - name: responseSize
      value_type: 2 # INT64
  - name: accesslog.combined
    display_name: Apache Combined Log Format
    log_template: '{{or (.originIp) "-"}} - {{or (.sourceUser) "-"}} [{{or (.timestamp.Format "02/Jan/2006:15:04:05 -0700") "-"}}] "{{or (.method) "-"}} {{or (.url) "-"}} {{or (.protocol) "-"}}" {{or (.responseCode) "-"}} {{or (.responseSize) "-"}} {{or (.referer) "-"}} {{or (.userAgent) "-"}}'
    labels:
    - name: originIp
      value_type: 6 # IP_ADDRESS
    - name: sourceUser
      value_type: 1 # STRING
    - name: timestamp
      value_type: 5 # TIMESTAMP
    - name: method
      value_type: 1 # STRING
    - name: url
      value_type: 1 # STRING
    - name: protocol
      value_type: 1 # STRING
    - name: responseCode
      value_type: 2 # INT64
    - name: responseSize
      value_type: 2 # INT64
    - name: referer
      value_type: 1 # STRING
    - name: userAgent
      value_type: 1 # STRING
`

const scYamlInitialSection = `
subject: namespace:ns
rules:
- selector: true
  aspects:
`
const scYamlSimpleOneAspectStrFromat = `
  - name: prometheus_reporting_all_metrics$s
    kind: metrics
    adapter: prometheus
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

// SOMETHING IS BREAKING WITH THIS HUGE LIST>>> NEED TO TEST MORE>..
//const scYamlOneLargeExprAspectStrFormat = `
//  - name: prometheus_reporting_all_metrics$s
//    kind: metrics
//    adapter: prometheus
//    params:
//      metrics:
//      - descriptorName: request_count
//        value: response.code | 100
//        labels:
//          source: source.name | target.name | source.name | target.name | source.name | target.name | source.name | target.name | source.name | target.name | source.name | target.name | source.name | target.name | source.name | target.name | "one"
//          target: target.name | source.name | target.name | source.name | target.name | source.name | target.name | source.name | target.name | source.name | target.name | source.name | target.name | source.name | target.name | source.name | "one"
//          service: api.name | target.name | api.name | target.name | api.name | target.name | api.name | target.name | api.name | target.name | api.name | target.name | api.name | target.name | api.name | target.name | api.name | target.name | "one"
//          method: api.method | target.name | api.method | target.name | api.method | target.name | api.method | target.name | api.method | target.name | api.method | target.name | api.method | target.name | api.method | target.name | api.method | target.name | "one"
//          response_code: response.code | response.code | response.code | response.code | response.code | response.code | response.code | response.code | response.code | response.code | response.code | response.code | response.code | response.code | 111
//`


const scYamlOneLargeExprAspectStrFormat = `
  - name: prometheus_reporting_all_metrics$s
    kind: metrics
    adapter: prometheus
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

func benchmarkDispatchSingleHugeAspect(b *testing.B, aspectStringFmt string, loopSize int64) {

	tmpfile, _ := ioutil.TempFile("", "TestReportWithJSServCnfg")
	fileSCName := tmpfile.Name()
	var scYaml string = scYamlInitialSection

	var i int64
	for i = 0; i < loopSize; i++ {
		scYaml = scYaml + strings.Replace(aspectStringFmt, "$s", strconv.FormatInt(i, 10), 1)
	}
	//fmt.Println(scYaml)
	_, _ = tmpfile.Write([]byte(scYaml))
	_ = tmpfile.Close()


	tmpfile, _ = ioutil.TempFile("", "TestReportWithJSGlobalcnfg")
	fileGSCName := tmpfile.Name()
	//defer func() { _ = os.Remove(gc) }()
	_, _ = tmpfile.Write([]byte(globalCnfg))
	_ = tmpfile.Close()


	apiPoolSize := 1
	adapterPoolSize := 1

	gp := pool.NewGoroutinePool(apiPoolSize, true)
	gp.AddWorkers(apiPoolSize)
	defer gp.Close()

	adapterGP := pool.NewGoroutinePool(adapterPoolSize, true)
	adapterGP.AddWorkers(adapterPoolSize)
	defer adapterGP.Close()

	eval := expr.NewCEXLEvaluator()
	adapterMgr := NewManager(adapter.Inventory(), aspect.Inventory(), eval, gp, adapterGP)
	cnfgMgr := config.NewManager(eval, adapterMgr.AspectValidatorFinder, adapterMgr.BuilderValidatorFinder,
		adapterMgr.SupportedKinds,
		fileGSCName,
		fileSCName,
		time.Second*time.Duration(1))

	cnfgMgr.Register(adapterMgr)
	cnfgMgr.Start()

	requestBag := attribute.GetMutableBag(nil)
	responseBag := attribute.GetMutableBag(nil)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = adapterMgr.Report(context.Background(), requestBag, responseBag)
		//if !status.IsOK(out) {
		//	t.Errorf("Report failed with %v", out)
		//}
	}
}

/*
Numbers with JS code:
BenchmarkDispatchSimpleOneAspect-12     	   10000	    157341 ns/op
BenchmarkDispatchSimple50Aspect-12      	     300	   5482677 ns/op
BenchmarkDispatchComplexOneAspect-12    	   10000	    176707 ns/op
BenchmarkDispatchComplex50Aspect-12     	     200	   6755903 ns/op
*/

/*
Numbers with existing code (no Javascript stuff):
BenchmarkDispatchSimpleOneAspect-12     	   10000	    169966 ns/op
BenchmarkDispatchSimple50Aspect-12      	     200	   6987477 ns/op
BenchmarkDispatchComplexOneAspect-12    	    5000	    212257 ns/op
BenchmarkDispatchComplex50Aspect-12     	     200	   8960063 ns/op
*/
func BenchmarkDispatchSimpleOneAspect(b *testing.B) {
	benchmarkDispatchSingleHugeAspect(b, scYamlSimpleOneAspectStrFromat, 1)
}

func BenchmarkDispatchSimple50Aspect(b *testing.B) {
	benchmarkDispatchSingleHugeAspect(b, scYamlSimpleOneAspectStrFromat, 50)
}

func BenchmarkDispatchComplexOneAspect(b *testing.B) {
	benchmarkDispatchSingleHugeAspect(b, scYamlOneLargeExprAspectStrFormat, 1)
}

func BenchmarkDispatchComplex50Aspect(b *testing.B) {
	benchmarkDispatchSingleHugeAspect(b, scYamlOneLargeExprAspectStrFormat, 50)
}