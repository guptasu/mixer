package template

import (
	"github.com/golang/protobuf/proto"
)

// Repository defines all the helper functions to access the generated template specific types and fields.
type Repository interface {
	GetConstructorDefaultConfig(template string) (proto.Message, bool)
}
type templateRepo struct{}

func (t templateRepo) GetConstructorDefaultConfig(template string) (proto.Message, bool) {
	if templateConstructorParamMap != nil {
		if v, ok := templateConstructorParamMap[template]; ok {
			return proto.Clone(v), true
		}
	}
	return nil, false
}

// NewTemplateRepository returns an implementation of Repository
func NewTemplateRepository() Repository {
	return templateRepo{}
}
