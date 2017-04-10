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
	"fmt"
	"io/ioutil"
	"istio.io/mixer/adapter"
	"istio.io/mixer/pkg/aspect"
	"istio.io/mixer/pkg/attribute"
	"istio.io/mixer/pkg/cnfgNormalizer"
	"istio.io/mixer/pkg/config"
	"istio.io/mixer/pkg/expr"
	"istio.io/mixer/pkg/pool"
	"istio.io/mixer/pkg/status"
	"testing"
	"time"
)

func TestActualReport(t *testing.T) {
	scYaml := `
subject: namespace:ns
rules:
        #- selector: service.name == “*”
        #- selector: service.name == "myservice"
- selector: true
  aspects:
  - name: prometheus_reporting_all_metrics
    kind: metrics
    adapter: prometheus
    params:
      metrics:
      - descriptorName: request_count
        # we want to increment this counter by 1 for each unique (source, target, service, method, response_code) tuple
        value: response.code | 100
        labels:
          source: source.name | "one"
          target: target.name | "one"
          service: api.name | "one"
          method: api.method | "one"
          response_code: response.code | 111
`
	globalCnfg := `
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

	tmpfile, _ := ioutil.TempFile("", "TestReportWithJS")
	fileSCName := tmpfile.Name()
	//defer func() { _ = os.Remove(gc) }()
	_, _ = tmpfile.Write([]byte(scYaml))
	_ = tmpfile.Close()

	tmpfile, _ = ioutil.TempFile("", "TestReportWithJS")
	fileGSCName := tmpfile.Name()
	//defer func() { _ = os.Remove(gc) }()
	_, _ = tmpfile.Write([]byte(globalCnfg))
	_ = tmpfile.Close()

	fmt.Println(fileSCName, fileGSCName)

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
		time.Second*time.Duration(1),
		"",
		cnfgNormalizer.NormalizedJavascriptConfigNormalizer{})
	fmt.Println(cnfgMgr)

	cnfgMgr.Register(adapterMgr)
	cnfgMgr.Start()

	requestBag := attribute.GetMutableBag(nil)
	responseBag := attribute.GetMutableBag(nil)

	out := adapterMgr.Report(context.Background(), requestBag, responseBag)

	if !status.IsOK(out) {
		t.Errorf("Report failed with %v", out)
	}
}
