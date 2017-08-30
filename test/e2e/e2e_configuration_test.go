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
	"testing"
	"context"
	"istio.io/mixer/pkg/attribute"
	"istio.io/mixer/pkg/adapter"
	e2eTmpl "istio.io/mixer/test/e2e/template"
	"istio.io/mixer/pkg/template"
)

// The adapter implementation fills this data and test can verify what was called.
// To use this variable across tests, every test should clean this variable value.
var globalActualHandlerCallInfoToValidate map[string]interface{}

const (

	globalCnfg = `
apiVersion: "config.istio.io/v1alpha2"
kind: attribute-manifest
metadata:
  name: istio-proxy
  namespace: default
spec:
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
---
`
	reportTestCnfg = `
apiVersion: "config.istio.io/v1alpha2"
kind: fakeHandler
metadata:
  name: fakeHandlerConfig
  namespace: istio-config-default

---

apiVersion: "config.istio.io/v1alpha2"
kind: report
metadata:
  name: reportInstance
  namespace: istio-config-default
spec:
  value: "2"
  dimensions:
    source: source.name
    target_ip: target.name

---

apiVersion: "config.istio.io/v1alpha2"
kind: mixer-rule
metadata:
  name: rule1
  namespace: istio-config-default
spec:
  selector: target.name == "*"
  actions:
  - handler: fakeHandlerConfig.fakeHandler.istio-config-default
    instances: [ reportInstance.report.istio-config-default ]

---
`
)

type adptr struct {
	name string
	handler adapter.Handler
	builder adapter.HandlerBuilder
}

func TestReport(t *testing.T) {

	tests := []struct {
		name      string
		oprtrCnfg string
		adapters []adapter.InfoFn
		templates map[string]template.Info
		validate func(t *testing.T, err error, callInfos map[string]interface{})
	}{
		{
			name:      "Report",
			oprtrCnfg: reportTestCnfg,
			adapters: []adapter.InfoFn{getFakeHndlrBldrInfoFn("fakeHandler")},
			templates: e2eTmpl.SupportedTmplInfo,
			validate: func(t *testing.T, err error, callInfo map[string]interface{}) {
				// validate globalActualHandlerCallInfoToValidate
				if len(callInfo) != 2 {
					t.Errorf("got call count %d\nwant %d", len(callInfo), 2)
				}

				if callInfo["ConfigureSampleReport"] == nil || callInfo["Build"] == nil {
					t.Errorf("got call info as : %v. \nwant calls %s and %s to have been called", callInfo, "ConfigureSample", "Build")
				}
			},
		},
	}
	for _, tt := range tests {
		configDir := getCnfgs(tt.oprtrCnfg, globalCnfg)
		defer func() {
			if !t.Failed() {
				_ = os.RemoveAll(configDir)
			} else {
				t.Logf("The configs are located at %s", configDir)
			}
		}() // nolint: gas

		globalActualHandlerCallInfoToValidate = make(map[string]interface{})

		requestBag := attribute.GetMutableBag(nil)
		requestBag.Set(configIdentityAttribute, identityDomainAttribute)

		dispatcher := getDispatcher(t, "fs://"+configDir, tt.adapters, tt.templates)
		err := dispatcher.Report(context.TODO(), requestBag)
		tt.validate(t, err, globalActualHandlerCallInfoToValidate)
	}
}
