package main

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"fmt"
	"github.com/golang/protobuf/proto"
	tmplExtns "istio.io/mixer/tools/codegen/template_extension"
	multierror "github.com/hashicorp/go-multierror"
)

type Model struct {
	Name string
	Check bool
	PackageName string
}

func validate(fds *descriptor.FileDescriptorSet) (Model, error) {
	result := &multierror.Error{}

	var templateDescriptorProto *descriptor.FileDescriptorProto = nil
	model := &Model{}
	for _, fdp := range fds.File {
		fmt.Println(proto.HasExtension(fdp.GetOptions(), tmplExtns.E_TemplateName))
		fmt.Println(proto.HasExtension(fdp.GetOptions(), tmplExtns.E_TemplateVariety))
		if !proto.HasExtension(fdp.GetOptions(), tmplExtns.E_TemplateName) && !proto.HasExtension(fdp.GetOptions(), tmplExtns.E_TemplateVariety) {
			continue
		} else if proto.HasExtension(fdp.GetOptions(), tmplExtns.E_TemplateName) && proto.HasExtension(fdp.GetOptions(), tmplExtns.E_TemplateVariety) {
			if templateDescriptorProto == nil {
				templateDescriptorProto = fdp
				model.PackageName = *fdp.Package
				tmplName, _ := proto.GetExtension(fdp.GetOptions(), tmplExtns.E_TemplateName)
				if name,ok := tmplName.(*string); !ok {
					result = multierror.Append(result, fmt.Errorf("%s should be of type string", tmplExtns.E_TemplateName.Name))
				} else {
					model.Name = *name
				}

				tmplVariety, _ := proto.GetExtension(fdp.GetOptions(), tmplExtns.E_TemplateVariety)
				model.Check = tmplVariety == tmplExtns.TemplateVariety_TEMPLATE_VARIETY_CHECK

			} else {
				result = multierror.Append(result, fmt.Errorf("Proto files %s and %s, both have" +
					" the options %s and %s. Only one proto file is allowed with those options",
					fdp.Name, templateDescriptorProto.Name,
					tmplExtns.E_TemplateVariety.Name, tmplExtns.E_TemplateName.Name))

			}
		} else {
			result = multierror.Append(result, fmt.Errorf("Proto files %s has only one of the " +
				"following two options %s and %s. Both options are required.",
				fdp.Name, tmplExtns.E_TemplateVariety.Name, tmplExtns.E_TemplateName.Name))
		}
	}

	if len(result.Errors) != 0 {
		return *model, result.ErrorOrNil()
	}


	return *model, result.ErrorOrNil()
}

func generateModel(fds *descriptor.FileDescriptorSet) (Model, error) {
	// TODO. Create a model for using the text tempaltes.
	return validate(fds)
}
