# logentry
--
    import "istio.io/mixer/template/logentry"


## Usage

```go
const TemplateName = "logentry"
```
Fully qualified name of this template

#### type Handler

```go
type Handler interface {
	adapter.Handler

	// HandleLogEntry is called by Mixer at request time to deliver instances to
	// to an adapter.
	HandleLogEntry(context.Context, []*Instance) error
}
```

Handler must be implemented by adapter code if it wants to process data
associated with the LogEntry template.

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

	// ConfigureLogEntryHandler is invoked by Mixer to pass all possible Types for instances that an adapter
	// may receive at runtime. Each type holds information about the shape of the instances.
	ConfigureLogEntryHandler(map[string]*Type) error
}
```

HandlerBuilder must be implemented by adapters if they want to process data
associated with the LogEntry template.

Mixer uses this interface to call into the adapter at configuration time to
configure it with adapter-specific configuration as well as all inferred types
the adapter is expected to handle.

#### type Instance

```go
type Instance struct {
	// Name of the instance as specified in configuration.
	Name string

	// Labels contain the payload for the log entry
	Labels map[string]interface{}

	// Severity indicates the log level for the log entry.
	Severity string
}
```

Instance is constructed by Mixer for the 'logentry.LogEntry' template.

LogEntry represents an individual entry within a log.
