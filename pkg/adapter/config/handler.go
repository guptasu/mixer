package config

import (
	"io"

	"github.com/golang/protobuf/proto"
)

// Handler represents default functionality every Adapter must implement.
type Handler interface {
	io.Closer

	// Name returns the official name of the aspects produced by this builder.
	Name() string
	// Description returns a user-friendly description of the aspects produced by this builder.
	Description() string
	DefaultConfig() proto.Message
	ValidateConfig(proto.Message) error
	Configure(proto.Message) error
}
