package template

import (
	"github.com/golang/protobuf/proto"

	pb "istio.io/api/mixer/v1/config/descriptor"
)

type (
	// Repository defines all the helper functions to access the generated template specific types and fields.
	Repository interface {
		GetConstructorDefaultConfig(template string) (proto.Message, bool)
		GetTypeInferFn(template string) (InferTypeFn, bool)
	}
	// TypeEvalFn evaluates an expression and returns the ValueType for the expression.
	TypeEvalFn func(string) (pb.ValueType, error)
	// InferTypeFn does Type inference from the Constructor.params proto message.
	InferTypeFn  func(interface{}, TypeEvalFn) (proto.Message, error)
	templateRepo struct{}
)

func (t templateRepo) GetConstructorDefaultConfig(template string) (proto.Message, bool) {
	if templateConstructorParamMap != nil {
		if v, ok := templateConstructorParamMap[template]; ok {
			return proto.Clone(v), true
		}
	}
	return nil, false
}

func (t templateRepo) GetTypeInferFn(template string) (InferTypeFn, bool) {
	if templateConstructorParamMap != nil {
		if v, ok := templateInferTypeFnMapping[template]; ok {
			return v, true
		}
	}
	return nil, false
}

// NewTemplateRepository returns an implementation of Repository
func NewTemplateRepository() Repository {
	return templateRepo{}
}
