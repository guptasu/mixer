package log

import "io"
import proto "istio.io/mixer/adapter/log"

// TODO: how do we expose a validation interface? Could the validation interface be autogened too so it's well typed?

// We're generating this, no need to use a proto.Message I think:
// we can generate code w/ strong typing.
type Config proto.LogParams

type Adapter interface {
  io.Closer

  NewProcessor(Config, []*Template) (Processor, error)
}

type Processor interface {
  io.Closer

  // TODO: I still think we need the arity of the error to match the arity of the input...
  Process([]*Instance) error
}

type Instance struct {
    Template *proto.LogEntryTemplate
    Name string
    // Note that for logging specifically, this is different than how it works today:
    // today the adapter never sees the template or the payload format, they only see
    // a constructed object with the template already evaluated. But I don't know how
    // we could automatically make that semantic change in our code gen.
    PayloadFormat proto.LogEntryTemplate_PayloadFormat
    LogTemplate string
    Dimensions map[string]interface{}
}