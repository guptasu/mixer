**Passing template specific types and adapter config to ****builder**

After builder object instantiation, Mixer configures the builder object by invoking various Template specific HandlerBuilder interface methods (example SetMetricTypes, SetQuotaTypes for 'metric' and 'quota' named Templates.) and passing a map of string-to-Type struct. The string key and the value Type represents the name of the instance as configured by the operator and the shape of the Instance object the adapter would receive during request time.

Given the above sample operator's handler configuration and 'metric' Template shows in above examples, the below examples shows the configuration-time call values:

![flow: example attr to types](./img/example%20attr%20to%20instance.svg)

During request time, every Instance object dispatch to the adapter has a 'Name' field. Adapter implementation should use the value of the Name field to lookup the shape description for the Instance object from the map of instance name(string)->Type that was passed during configuration time through the builder object.

Once Mixer has called into various template-specific Set****Types methods,

Mixer calls the SetAdapterConfig method on the builder, and once done then Mixer calls the Validate followed by the Build method. SetAdapterConfig gives the builder the adapter specific configuration, Validate allows builder to validate the operator configuration based on the the provided Template specific Types and the Adapter specific configuration.

**Instantiating ****handler**

Once builder is validated, Mixer calls its Build method, which returns a handler object which Mixer invokes during request processing. The handler instance constructed must implement all the Handler interfaces (runtime request serving interfaces) for all the templates the adapter has registered for. If the returned handler fails to implement the required interface for the adapter?s supported templates, Mixer reports an error and doesn?t serve runtime traffic to the particular handler.

*NOTE: *In the Build method, adapters must do all the bootstrapping work (example establishing connection with backend system, initializing cache and more) that they need to start receiving data during request time.

**Closing ****handler**

When a handler is no longer useful, Mixer calls it close method. In the Close method an adapter is expected to release all the allocated resources and close all remote connections to the backends if it has any.

### Request-time

During this time Mixer dispatches the instance objects to the adapter based on the routing rules that operator has configured. Mixer does this by invoking the Handle* functions on the handler object.

Given the above example operator's config (instance, action, handler configuration) and ['metric' Template](https://docs.google.com/document/d/1rKPt2Z2acy4pRwcvScPa-Na-EnU-6poGR-Rj_fZtaIc/edit#heading=h.ee6dn8otn4o0), the following examples shows the request-time Instance objects created for a given input attribute bag:

![flow: example attr o instance mapping](./img/example%20attr%20to%20instance.svg)

# Example

The following sample adapters just illustrate the basic skeleton of the adapter code and do not provide any functionality. They always return success. For examples of real world Adapters, see [implementation of built-in Adapters within Mixer framework](https://github.com/istio/mixer/tree/master/adapter).

* Sample no-op adapter that supports the above [sample 'metric' Template](https://docs.google.com/document/d/1rKPt2Z2acy4pRwcvScPa-Na-EnU-6poGR-Rj_fZtaIc/edit#heading=h.ee6dn8otn4o0)

<table>
  <tr>
    <td>type (
  builder struct{}
  handler struct{}
)

// ensure our types implement the requisite interfaces
var _ metric.HandlerBuilder = builder{}
var _ metric.Handler = handler{}

///////////////// Configuration Methods ///////////////

func (builder) Build(Context.Context, adapter.Env) (adapter.Handler, error) {
  return handler{}, nil
}
func (builder) SetAdapterConfig(adapter.Config)                      {}
func (builder) Validate() (*adapter.ConfigErrors)                 { return }

func (builder) SetMetricTypes(map[string]*metric.Type){
  ...
}

////////////////// Runtime Methods //////////////////////////

func (handler) HandleMetric(context.Context, []*metric.Instance) error {
  return nil
}

func (handler) Close() error { return nil }

////////////////// Bootstrap //////////////////////////
// GetInfo returns the Info for this adapter.

func GetInfo() adapter.BuilderInfo {
  return adapter.BuilderInfo{
     Name:        "istio.io/mixer/adapter/noop1",
     Description: "Does nothing",
     SupportedTemplates: []string{
        metric.TemplateName,
     },
     NewBuilder: func() adapter.HandlerBuilder { return builder{} },
     DefaultConfig:        &types.Empty{},
  }
}
</td>
  </tr>
</table>


* Sample no-op adapter that supports the above [sample 'listentry' Template](https://docs.google.com/document/d/1rKPt2Z2acy4pRwcvScPa-Na-EnU-6poGR-Rj_fZtaIc/edit#heading=h.qgv3mdgv1nfj).

<table>
  <tr>
    <td>type (
  builder struct{}
  handler struct{}
)

// ensure our types implement the requisite interfaces
var _ listentry.HandlerBuilder = builder{}
var _ listentry.Handler = handler{}

///////////////// Configuration Methods ///////////////

func (builder) Build(Context.Context, adapter.Env) (adapter.Handler, error) {
  return handler{}, nil
}
func (builder) SetAdapterConfig(adapter.Config)                      {}
func (builder) Validate() (*adapter.ConfigErrors)                 { return }

func (builder) SetListEntryTypes(map[string]*listentry.Type){
  ...
}


////////////////// Runtime Methods //////////////////////////

var checkResult = adapter.CheckResult{
  Status:        rpc.Status{Code: int32(rpc.OK)},
  ValidDuration: 1000000000 * time.Second,
  ValidUseCount: 1000000000,
}

func (handler) HandleListEntry(context.Context, *listentry.Instance) (adapter.CheckResult, error) {
  return checkResult, nil
}

func (handler) Close() error { return nil }

////////////////// Bootstrap //////////////////////////

// GetInfo returns the Info associated with this adapter implementation.
func GetInfo() adapter.BuilderInfo {
  return adapter.BuilderInfo{
     Name:        "istio.io/mixer/adapter/noop2",
     Description: "Does nothing",
     SupportedTemplates: []string{
        listentry.TemplateName,
     },
     NewBuilder: func() adapter.HandlerBuilder { return builder{} },
     DefaultConfig:        &types.Empty{},
  }
}
</td>
  </tr>
</table>


* Sample of a no-op adapter that supports the above [sample 'quota' Template](https://docs.google.com/document/d/1rKPt2Z2acy4pRwcvScPa-Na-EnU-6poGR-Rj_fZtaIc/edit#heading=h.67r0dd5r6jgw).

<table>
  <tr>
    <td>type (
  builder struct{}
  handler struct{}
)

// ensure our types implement the requisite interfaces
var _ quota.HandlerBuilder = builder{}
var _ quota.Handler = handler{}

///////////////// Configuration Methods ///////////////

func (builder) Build(Context.Context, adapter.Env) (adapter.Handler, error) {
  return handler{}, nil
}
func (builder) SetAdapterConfig(adapter.Config)                      {}
func (builder) Validate() (*adapter.ConfigErrors)                 { return }

func (builder) SetQuotaTypes(map[string]*quota.Type){
  ...
}

////////////////// Runtime Methods //////////////////////////

func (handler) HandleQuota(ctx context.Context, _ *quota.Instance, args adapter.QuotaRequestArgs) (adapter.QuotaResult2, error) {
  return adapter.QuotaResult2{
        ValidDuration: 1000000000 * time.Second,
        Amount:        args.QuotaAmount,
     },
     nil
}

func (handler) Close() error { return nil }

////////////////// Bootstrap //////////////////////////

// GetInfo returns the Info associated with this adapter implementation.
func GetInfo() adapter.BuilderInfo {
  return adapter.BuilderInfo{
     Name:        "istio.io/mixer/adapter/noop2",
     Description: "Does nothing",
     SupportedTemplates: []string{
        quota.TemplateName,
     },
     NewBuilder: func() adapter.HandlerBuilder { return builder{} },
     DefaultConfig:        &types.Empty{},
  }
}
</td>
  </tr>
</table>


The above section provided a complete example of a simple adapter. In the next sections we?ll look in more detail about how to build mixer with custom adapters that are are not shipped with the default Mixer build, and step-by-step guide to build a simple adapter.

# Summary diagrams

Below diagrams show the relationship between, Adapter, Mixer, Template and operator config. It also shows the flow of Mixer at boot time, how it interacts with Adapters and operator configuration. The diagrams also demonstrate how handler, rule and instance configuration of operator config is translated to calls into Adapters during config load time and request time.

First let's look into how Mixer, adapter, templates and operator configurations are related

![template, adapter and operator config relationship](./img/Template%2C%20Adapter%20and%20Operator%20config%20relationship.svg)

Now we have understood the relationship between various artifacts, let's look into what happens at the time of Mixer start, when operator configuration is loaded/changed and when request is received.

![flow: mixer start](./img/Mixer%20Start%20Flow.svg)

![flow: operator config change](./img/Operator%20config%20change.svg)

![flow: incoming api](./img/Request%20time%20.svg)

# Plug adapter into Mixer

For a new adapter to plug into Mixer, you will have to add your adapter's reference into Mixer's inventory [build file](https://github.com/istio/mixer/blob/master/adapter/BUILD)'s inventory_library rule. In the *deps *section add a reference to adapter's go_library build rule, and in the *packages* section add the short name and the go import path to the adapter package that implements the GetInfo function. These two changes will plug your custom adapter into Mixer.

# Testing

We provide a simple adapter test framework. The framework instantiates a in-proc Mixer gRPC server with a config store backed by local filesystem, and also a Mixer gRPC client in test process, which allows stepping through adapter code in test cases. The test framework is implemented in the [test/testenv](https://github.com/istio/mixer/blobhttps://github.com/istio/mixer/tree/master/test/testenvmaster/test/testenv/testenv_test.go) directory. A [sample](https://github.com/istio/mixer/blob/master/test/testenv/testenv_test.go) test is provided to show how to use this test framework to test a dummy adapter called denier. To setup the environment, adapter developer need author adapter config files. Sample adapter config can be found in [/testdata](https://github.com/istio/mixer/tree/master/testdata/config) directory.

# Do's and dont?s

* Adapters must use env.Logger for logging during execution. This logger understands about which adapter is running and routes the data to the place where the operator wants to see it.

* Adapters must use env.ScheduleWork or env.ScheduleDaemon in order to dispatch goroutines. This ensures all adapter goroutines are prevented from crashing Mixer as a whole by catching any panics they produce.

# Built-in Templates

Mixer ships with a set of built-in templates that are ready to use by adapters:

* [listentry](https://github.com/istio/mixer/tree/master/template/listentry)

* [logentry](https://github.com/istio/mixer/tree/master/template/logentry)

* [quota](https://github.com/istio/mixer/tree/master/template/quota)

* [metric](https://github.com/istio/mixer/tree/master/template/metric)

* [checknothing](https://github.com/istio/mixer/tree/master/template/checknothing)

* [reportnothing](https://github.com/istio/mixer/tree/master/template/reportnothing)

Using the above templates, the Mixer team has implemented a set of adapters that ships as part of the default Mixer binary. They are located at [istio/mixer/adapter](https://github.com/istio/mixer/tree/master/adapter). They are good examples for reference when implementing new adapters.

# Walkthrough of Adapter implementation (30 minutes)

Please refer to [Adapter Development Walkthrough](https://docs.google.com/document/d/1ZjGtmf27AQLxq7Au5lpI_P-YDdjDaqTpU2tNLdK3IMI/edit#)

