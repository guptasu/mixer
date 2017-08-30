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
	"istio.io/mixer/pkg/adapter"
	e2eTmpl "istio.io/mixer/test/e2e/template"
	"istio.io/mixer/pkg/template"
)


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

type testData struct {
	name          string
	oprtrCnfg     string
	adptBehaviors []AdptBehavior
	templates     map[string]template.Info
	attribs       map[string]interface{}
	validate      func(t *testing.T, err error, sypAdpts []*spyAdapter)
}

func TestReport(t *testing.T) {

	tests := []testData{
		{
			name:          "Report",
			oprtrCnfg:     reportTestCnfg,
			adptBehaviors: []AdptBehavior{AdptBehavior{name:"fakeHandler"}},
			templates:     e2eTmpl.SupportedTmplInfo,
			attribs:       map[string]interface{}{},
			validate: func(t *testing.T, err error, spyAdpts []*spyAdapter) {
				// validate globalActualHandlerCallInfoToValidate
				if len(spyAdpts[0].builderCallData.data) != 2 {
					t.Errorf("got call count %d\nwant %d", len(spyAdpts[0].builderCallData.data), 2)
				}

				if spyAdpts[0].builderCallData.data["ConfigureSampleReport"] == nil || spyAdpts[0].builderCallData.data["Build"] == nil {
					t.Errorf("got call info as : %v. \nwant calls %s and %s to have been called", spyAdpts[0].builderCallData.data, "ConfigureSample", "Build")
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

		adapterInfos, spyAdapters := cnstrAdapterInfos(tt.adptBehaviors)
		dispatcher := getDispatcher(t, "fs://"+configDir, adapterInfos, tt.templates)
		err := dispatcher.Report(context.TODO(), getAttrBag(tt.attribs))

		tt.validate(t, err, spyAdapters)
	}
}
