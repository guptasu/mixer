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
	"context"
	"istio.io/mixer/pkg/attribute"
	"istio.io/mixer/pkg/adapter"
)

// The adapter implementation fills this data and test can verify what was called.
// To use this variable across tests, every test should clean this variable value.
var globalActualHandlerCallInfoToValidate map[string]interface{}

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

func executeReport(
	t *testing.T,
	declarativeSrvcCnfgFilePath string,
	declaredGlobalCnfgFilePath string,
	configStore2URL string,
	requestBag attribute.Bag,
	adptInfos []adapter.InfoFn,
) error {
	dispatcher := getDispatcher(t, configStore2URL, declaredGlobalCnfgFilePath, declarativeSrvcCnfgFilePath, adptInfos)
	return dispatcher.Report(context.TODO(), requestBag)
}

func TestConfigFlow(t *testing.T) {
	sc, gsc := getCnfgs(srvConfig, globalCnfg)
	defer func() {
		_ = os.Remove(sc.Name())
		_ = os.Remove(gsc.Name())
	}() // nolint: gas

	globalActualHandlerCallInfoToValidate = make(map[string]interface{})

	requestBag := attribute.GetMutableBag(nil)
	requestBag.Set(configIdentityAttribute, identityDomainAttribute)
	adapters := []adapter.InfoFn{GetFakeHndlrBuilderInfo}

	executeReport(t, sc.Name(), gsc.Name(), "fs://"+path.Dir(sc.Name()), requestBag, adapters)

	// validate globalActualHandlerCallInfoToValidate
	if len(globalActualHandlerCallInfoToValidate) != 2 {
		t.Errorf("got call count %d\nwant %d", len(globalActualHandlerCallInfoToValidate), 2)
	}

	if globalActualHandlerCallInfoToValidate["ConfigureSampleReport"] == nil || globalActualHandlerCallInfoToValidate["Build"] == nil {
		t.Errorf("got call info as : %v. \nwant calls %s and %s to have been called", globalActualHandlerCallInfoToValidate, "ConfigureSample", "Build")
	}
}
