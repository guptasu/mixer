package template

import (
	"github.com/golang/protobuf/proto"

	pb "istio.io/api/mixer/v1/config/descriptor"
)

type (
	// Repository defines all the helper functions to access the generated template specific types and fields.
	Repository interface {
		GetTemplateInfo(template string) (Info, bool)
	}
	// TypeEvalFn evaluates an expression and returns the ValueType for the expression.
	TypeEvalFn func(string) (pb.ValueType, error)
	// InferTypeFn does Type inference from the Constructor.params proto message.
	InferTypeFn func(interface{}, TypeEvalFn) (proto.Message, error)
	// Info contains all the information related a template like
	// Default constructor params, type inference method etc.
	Info struct {
		CnstrDefConfig proto.Message
		InferTypeFn    InferTypeFn
	}
	// templateRepo implements Repository
	templateRepo struct{}
)

func (t templateRepo) GetTemplateInfo(template string) (Info, bool) {
	if v, ok := templateInfos[template]; ok {
		return v, true
	}
	return Info{}, false
}

// NewTemplateRepository returns an implementation of Repository
func NewTemplateRepository() Repository {
	return templateRepo{}
}
