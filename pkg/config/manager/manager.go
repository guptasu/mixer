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
	"time"
	"github.com/golang/glog"
	"istio.io/mixer/pkg/config"
	"istio.io/mixer/pkg/adapter"
	"istio.io/mixer/pkg/attribute"
	"istio.io/mixer/pkg/config/descriptor"
	pb "istio.io/mixer/pkg/config/proto"
	"istio.io/mixer/pkg/aspect"
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

	rt.JSStr = getJS(vd)
	return rt, c.descriptorFinder, nil
}

func getJS(vd *config.Validated) string {
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
	//var reportMethodStr string
	//for _, aspectRule := range vd.serviceConfig.GetRules() {
		//var tmpReportMethodStr string
		//for _, aspect := range aspectRule.GetAspects() {
			//aspect.Kind
		//}
	//}
	var mgr aspect.Manager
	fmt.Println(mgr)
	js := `
		function report(propertyBag) {

RecordToprometheus({"request_count": {value: "1",target: "B",method: "DD",response_code: "200",service: "C",source: "A"}})
		}
		// more check and quota methods



        // such methods will be dynamically generated for each aspect in the user config.
        // For now assume there is only one aspect of kind metric for which we computed the call back
        // method name to be "ParticularAspectReport". Currenly the Aspect Manager is invoked for each
        // aspect, but eventually it will be invoked for check/report/quota call and calls to various
        // aspects will be done here in this file.
        function RecordToprometheus(val) {
          CallBackFromUserScript_go("metrics", "prometheus", val)
        }
	`
	return js
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
