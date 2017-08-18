# quota
--
    import "istio.io/mixer/template/quota"


## Usage

```go
const TemplateName = "quota"
```
Fully qualified name of this template

#### type Handler

```go
type Handler interface {
	adapter.Handler

	// HandleQuota is called by Mixer at request time to deliver instances to
	// to an adapter.
	HandleQuota(context.Context, *Instance, adapter.QuotaRequestArgs) (adapter.QuotaResult2, error)
}
```

Handler must be implemented by adapter code if it wants to process data
associated with the Quota template.

Mixer uses this interface to call into the adapter at request time in order to
dispatch created instances to the adapter. Adapters take the incoming instances
and do what they need to achieve their primary function.

The name of each instance can be used as a key into the Type map supplied to the
adapter at configuration time. These types provide descriptions of each specific
instances.

#### type HandlerBuilder

```go
type HandlerBuilder interface {
	adapter.HandlerBuilder

	// ConfigureQuotaHandler is invoked by Mixer to pass all possible Types for instances that an adapter
	// may receive at runtime. Each type holds information about the shape of the instances.
	ConfigureQuotaHandler(map[string]*Type) error
}
```

HandlerBuilder must be implemented by adapters if they want to process data
associated with the Quota template.

Mixer uses this interface to call into the adapter at configuration time to
configure it with adapter-specific configuration as well as all inferred types
the adapter is expected to handle.

#### type Instance

```go
type Instance struct {
	// Name of the instance as specified in configuration.
	Name string

	// The unique identity of the particular quota to manipulate.
	Dimensions map[string]interface{}
}
```

Instance is constructed by Mixer for the 'quota.Quota' template.
