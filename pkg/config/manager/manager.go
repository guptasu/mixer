// Copyright 2017 Istio Authors
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

package configManager

import (
	"crypto/sha1"
	"io/ioutil"
	"sync"
	"bytes"
	"time"
	"github.com/golang/glog"
	"istio.io/mixer/pkg/config"
	"istio.io/mixer/pkg/adapter"
	aconfig "istio.io/mixer/pkg/aspect/config"
	"istio.io/mixer/pkg/attribute"
	"istio.io/mixer/pkg/config/descriptor"
	pb "istio.io/mixer/pkg/config/proto"


	"istio.io/mixer/pkg/expr"
	"fmt"
)

// Resolver resolves configuration to a list of combined configs.
type Resolver interface {
	// Resolve resolves configuration to a list of combined configs.
	Resolve(bag attribute.Bag, aspectSet config.AspectSet) ([]*pb.Combined, error)
	// get JS string
	GetJS() string
}

// ChangeListener listens for config change notifications.
type ChangeListener interface {
	ConfigChange(cfg Resolver, df descriptor.Finder)
}

// Manager represents the config Manager.
// It is responsible for fetching and receiving configuration changes.
// It applies validated changes to the registered config change listeners.
// api.Handler listens for config changes.
type Manager struct {
	eval             expr.Evaluator
	aspectFinder     config.AspectValidatorFinder
	builderFinder    config.BuilderValidatorFinder
	descriptorFinder descriptor.Finder
	findAspects      config.AdapterToAspectMapper
	loopDelay        time.Duration
	globalConfig     string
	serviceConfig    string

	cl      []ChangeListener
	closing chan bool
	scSHA   [sha1.Size]byte
	gcSHA   [sha1.Size]byte

	sync.RWMutex
	lastError error
}

// NewManager returns a config.Manager.
// Eval validates and evaluates selectors.
// It is also used downstream for attribute mapping.
// AspectFinder finds aspect validator given aspect 'Kind'.
// BuilderFinder finds builder validator given builder 'Impl'.
// LoopDelay determines how often configuration is updated.
// The following fields will be eventually replaced by a
// repository location. At present we use GlobalConfig and ServiceConfig
// as command line input parameters.
// GlobalConfig specifies the location of Global Config.
// ServiceConfig specifies the location of Service config.
func NewManager(eval expr.Evaluator, aspectFinder config.AspectValidatorFinder, builderFinder config.BuilderValidatorFinder,
	findAspects config.AdapterToAspectMapper, globalConfig string, serviceConfig string, loopDelay time.Duration) *Manager {
	m := &Manager{
		eval:          eval,
		aspectFinder:  aspectFinder,
		builderFinder: builderFinder,
		findAspects:   findAspects,
		loopDelay:     loopDelay,
		globalConfig:  globalConfig,
		serviceConfig: serviceConfig,
		closing:       make(chan bool),
	}
	return m
}

// Register makes the ConfigManager aware of a ConfigChangeListener.
func (c *Manager) Register(cc ChangeListener) {
	c.cl = append(c.cl, cc)
}

func read(fname string) ([sha1.Size]byte, string, error) {
	var data []byte
	var err error
	if data, err = ioutil.ReadFile(fname); err != nil {
		return [sha1.Size]byte{}, "", err
	}
	return sha1.Sum(data), string(data[:]), nil
}

// fetch config and return runtime if a new one is available.
func (c *Manager) fetch() (*config.Runtime, descriptor.Finder, error) {
	var vd *config.Validated
	var cerr *adapter.ConfigErrors

	gcSHA, gc, err2 := read(c.globalConfig)
	if err2 != nil {
		return nil, nil, err2
	}

	scSHA, sc, err1 := read(c.serviceConfig)
	if err1 != nil {
		return nil, nil, err1
	}

	if gcSHA == c.gcSHA && scSHA == c.scSHA {
		return nil, nil, nil
	}

	v := config.NewValidator(c.aspectFinder, c.builderFinder, c.findAspects, true, c.eval)
	if vd, cerr = v.Validate(sc, gc); cerr != nil {
		return nil, nil, cerr
	}

	c.descriptorFinder = descriptor.NewFinder(v.GetValidatedGSC())

	c.gcSHA = gcSHA
	c.scSHA = scSHA
	rt := config.NewRuntime(vd, c.eval)

	rt.JSStr = getJS(vd, c)
	return rt, c.descriptorFinder, nil
}

func getJS(vd *config.Validated, c *Manager) string {
	//vd.serviceConfig
	/*
	create methodStrs for each method (chk, report, quota)
	for each rule r:
	  create tmpMethodStrs for each method (chk, report, quota)
          for each aspect a:
	    generate JS k via Aspect Managers
	    add k to appropriate tmpMethodStrs.
	  for each non empty tmpMethodStrs {
	    insert : if (rule) {
	      tmpMethodStr
	    }
	  }
	  Do recurrs for each nested rules.
	*/
	var reportMethodStr string
	var injectionMethods string
	for _, aspectRule := range vd.GetValidatedSC().GetRules() {
		var tmpReportMethodStr string
		for _, aspect := range aspectRule.GetAspects() {
			if aspect.Kind == "metrics" {
				invocation, wrapperMtdToIngest := GetJSInvocationForMetricAspect(aspect.Params.(*aconfig.MetricsParams), aspect.Adapter)
				tmpReportMethodStr = tmpReportMethodStr + invocation
				// HACK. skipping duplication incertion of mtd.
				if len(injectionMethods) == 0 {
					injectionMethods = injectionMethods + wrapperMtdToIngest
				}
			}
		}
		if len(tmpReportMethodStr) > 0 {
			//TODO FIX FOR COMPLEXT OPTIONS
			var ifStatementStr string
			if len(aspectRule.Selector) > 0 {
				expressionStr := aspectRule.Selector
				ifStatementStr = getJSForExpression(expressionStr)

			} else {
				ifStatementStr = "true"
			}
			ruleIfBlocked := fmt.Sprintf(`
			if (%s) {
			  %s
			}
			`, ifStatementStr, tmpReportMethodStr)

			reportMethodStr = reportMethodStr + ruleIfBlocked
		}


	}
	allJSMethodFormat := `
		function report(propertyBag) {
		    %s
		}
		function check(propertyBag) {
		    // TODO
		}
		function quota(propertyBag) {
		    // TODO
		}
		`
	userJSAllCode := fmt.Sprintf(allJSMethodFormat, reportMethodStr)
	var userScript = userJSAllCode + "\n" + injectionMethods

	fmt.Println(userScript)

	return userScript
}
// fetchAndNotify fetches a new config and notifies listeners if something has changed
func (c *Manager) fetchAndNotify() error {
	rt, df, err := c.fetch()
	if err != nil {
		c.Lock()
		c.lastError = err
		c.Unlock()
		return err
	}
	if rt == nil {
		return nil
	}

	glog.Infof("Installing new config from %s sha=%x ", c.serviceConfig, c.scSHA)
	for _, cl := range c.cl {
		cl.ConfigChange(rt, df)
	}
	return nil
}

// LastError returns last error encountered by the manager while processing config.
func (c *Manager) LastError() (err error) {
	c.RLock()
	err = c.lastError
	c.RUnlock()
	return err
}

// Close stops the config manager go routine.
func (c *Manager) Close() { close(c.closing) }

func (c *Manager) loop() {
	ticker := time.NewTicker(c.loopDelay)
	defer ticker.Stop()
	done := false
	for !done {
		select {
		case <-ticker.C:
			err := c.fetchAndNotify()
			if err != nil {
				glog.Warning(err)
			}
		case <-c.closing:
			done = true
		}
	}
}

// Start watching for configuration changes and handle updates.
func (c *Manager) Start() {
	err := c.fetchAndNotify()
	// We make an attempt to synchronously fetch and notify configuration
	// If it is not successful, we will continue to watch for changes.
	go c.loop()
	if err != nil {
		glog.Warning("Unable to process config: ", err)
	}
}

///////////////////// THIS SHOULD BELONG TO METRIC ASPECT MANAGER /////////////////

func GetJSInvocationForMetricAspect(metricsParams *aconfig.MetricsParams, adapterName string) (string, string) {
	params := metricsParams
	//var allMetricsStr bytes.Buffer
	var metricsStr bytes.Buffer
	for _, metric := range params.Metrics {
		var labelStr bytes.Buffer
		labelStr.WriteString(fmt.Sprintf("%s: %s", "value", getJSForExpression(metric.Value)))
		labelLen := len(metric.Labels)
		if labelLen != 0 {
			labelStr.WriteString(",\n")
		}
		for key, value := range metric.Labels {
			labelStr.WriteString(fmt.Sprintf(`				    %s: %s`, key, getJSForExpression(value)))
			labelLen--
			if labelLen != 0 {
				labelStr.WriteString(",\n")
			}
		}
		metricsStr.WriteString(fmt.Sprintf(`
			      "%s": {
			        %s
			      },
			    `, metric.DescriptorName, labelStr.String()))

	}
	metricsStrBuilt := metricsStr.String()
	if metricsStrBuilt[len(metricsStrBuilt)-1] == ',' {
		metricsStrBuilt = metricsStrBuilt[0 : len(metricsStrBuilt)-1]
	}
	callStr := fmt.Sprintf(`
				%s({
	  			  %s
				})
	`, GetJSWrapperMethodNameForMetricAspect(adapterName), metricsStrBuilt)
	return callStr, GetJSWrapperMethodsToInjectForMetricAspect(adapterName, "metrics")
}

func GetJSWrapperMethodNameForMetricAspect(adapterName string) string {
	return "RecordTo" + adapterName
}

func GetJSWrapperMethodsToInjectForMetricAspect(adapterName string, kindName string) string {

	methodName := GetJSWrapperMethodNameForMetricAspect(adapterName)
	var embeddedMethodsInUserScriptFmt = `
        // such methods will be dynamically generated for each aspect in the user config.
        // For now assume there is only one aspect of kind metric for which we computed the call back
        // method name to be "ParticularAspectReport". Currenly the Aspect Manager is invoked for each
        // aspect, but eventually it will be invoked for check/report/quota call and calls to various
        // aspects will be done here in this file.
        function %s(val) {
          CallBackFromUserScript_go("%s", "%s", val)
        }
`
	return fmt.Sprintf(embeddedMethodsInUserScriptFmt, methodName, kindName, adapterName)
}

func getJSForExpression(expression string) string{
	ex, err := expr.Parse(expression)
	var out string
	if err != nil {
		glog.Warning("Unable to parse : %s. %v. Setting expression to false", expression, err)
		out = "false"
	} else {
		condition, _ := EvalJSExpession(ex, expr.FuncMap(), "propertyBag.Get")
		out = condition
	}
	return out
}

// Eval evaluates the expression given an attribute bag and a function map.
func EvalJSExpession(e *expr.Expression, fMap map[string]expr.FuncBase, getPropMtdName string) (string, error) {
	if e.Const != nil {
		return e.Const.StrValue, nil
	}
	if e.Var != nil {
		return fmt.Sprintf(getPropMtdName + "(\"%s\")[0]", e.Var.Name), nil
	}

	fn := fMap[e.Fn.Name]
	if fn == nil {
		return "", fmt.Errorf("unknown function: %s", e.Fn.Name)
	}
	// may panic
	if e.Fn.Name == "EQ" {
		leftStr, _ := EvalJSExpession(e.Fn.Args[0], fMap, getPropMtdName)
		rightStr, _ := EvalJSExpession(e.Fn.Args[1], fMap, getPropMtdName)
		return fmt.Sprintf("%s == %s", leftStr, rightStr), nil
	}
	if e.Fn.Name == "OR" {
		//(age < 18) ? "Too young":"Old enough"
		allArgs := e.Fn.Args
		if len(allArgs) > 0 {
			chkIfExists := fmt.Sprintf(getPropMtdName+"(\"%s\")[1]", allArgs[0].Var.Name)
			leftexp, _ := EvalJSExpession(e.Fn.Args[0], fMap, getPropMtdName)
			rightexp, _ := EvalJSExpession(e.Fn.Args[1], fMap, getPropMtdName)
			return fmt.Sprintf("%s ? %s : %s", chkIfExists, leftexp, rightexp), nil
		}
	}
	return "", nil
}