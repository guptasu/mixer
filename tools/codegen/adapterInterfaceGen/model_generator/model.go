package model_generator

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	multierror "github.com/hashicorp/go-multierror"
	tmplExtns "istio.io/mixer/tools/codegen/template_extension"
	"strings"
)
type Model struct {
	// top level fields
	Name        string
	Check       bool
	PackageName string
	VarietyName string

	// types
	TypeFullName string

	// imports
	Imports []string

	ConstructorFields []FieldInfo
}

type FieldInfo struct {
	Name string
	Type TypeInfo
}

type TypeInfo struct {
	Name   string
	IsExpr bool

	IsMap     bool
}

const FullNameOfExprMessage = "*istio_mixer_v1_config_template.Expr"

func (g *ModelGenerator) ConstructModel(fds *descriptor.FileDescriptorSet) (Model, error) {
	result := &multierror.Error{}
	model := &Model{}
	model.Imports = make([]string, 0)

	templateProto := getTemplateProto(fds, result)
	g.file = g.FileOf(templateProto)
	if len(result.Errors) != 0 {
		return *model, result.ErrorOrNil()
	}

	addTopLevelFields(model, templateProto, result)
	g.addFieldsOfConstructor(model, templateProto, result)
	model.Imports = g.generateImports()
	g.getTypeNameForType(model, templateProto, result)
	return *model, result.ErrorOrNil()
}


func addTopLevelFields(model *Model, fdp *descriptor.FileDescriptorProto, errors *multierror.Error) {
	model.PackageName = PackageName(*fdp.Package)
	tmplName, _ := proto.GetExtension(fdp.GetOptions(), tmplExtns.E_TemplateName)
	if name, ok := tmplName.(*string); !ok {
		errors = multierror.Append(errors, fmt.Errorf("%s should be of type string", tmplExtns.E_TemplateName.Name))
	} else {
		model.Name = *name
	}

	tmplVariety, _ := proto.GetExtension(fdp.GetOptions(), tmplExtns.E_TemplateVariety)
	if tmplVariety == tmplExtns.TemplateVariety_TEMPLATE_VARIETY_CHECK {
		model.Check = true
		model.VarietyName = "Check"
	} else {
		model.Check = false
		model.VarietyName = "Report"
	}
}

func (g *ModelGenerator) getTypeNameForType(model *Model, fdp *descriptor.FileDescriptorProto, errors *multierror.Error) {
	var typeDesc *descriptor.DescriptorProto = nil
	for _, desc := range fdp.MessageType {
		if *desc.Name == "Type" {
			typeDesc = desc
			break
		}
	}
	if typeDesc == nil {
		errors = multierror.Append(errors, fmt.Errorf("%s should have a message 'Type'", fdp.Name))
	}

	model.TypeFullName = g.TypeName(newDescriptor(typeDesc, nil, fdp, 0))
}

func (g *ModelGenerator) addFieldsOfConstructor(model *Model, fdp *descriptor.FileDescriptorProto, errors *multierror.Error) {
	model.ConstructorFields = make([]FieldInfo, 0)
	var cstrDesc *descriptor.DescriptorProto = nil
	for _, desc := range fdp.MessageType {
		if *desc.Name == "Constructor" {
			cstrDesc = desc
			break
		}
	}
	if cstrDesc == nil {
		errors = multierror.Append(errors, fmt.Errorf("%s should have a message 'Constructor'", fdp.Name))
	}

	for _, fieldDesc := range cstrDesc.Field {

		fieldName := CamelCase(*fieldDesc.Name)
		typename := g.GoType(cstrDesc, fieldDesc)
		typename = strings.Replace(typename, FullNameOfExprMessage, "interface{}", 1)

		model.ConstructorFields = append(model.ConstructorFields, FieldInfo{Name: fieldName, Type: TypeInfo{Name:typename}})
	}
}

func getTemplateProto(fds *descriptor.FileDescriptorSet, errors *multierror.Error) *descriptor.FileDescriptorProto {
	var templateDescriptorProto *descriptor.FileDescriptorProto = nil

	erroneousFiles := []string {
		"mixer/v1/config/descriptor/value_type.proto",
		"mixer/tools/codegen/template_extension/TemplateExtensions.proto",
	}

	for _, fdp := range fds.File {
		// TODO : Temporary hack..
		// For some reason the below code is panicing for files that are specified in the list.
		if contains(erroneousFiles, *fdp.Name) {
			continue
		}
		if !proto.HasExtension(fdp.GetOptions(), tmplExtns.E_TemplateName) && !proto.HasExtension(fdp.GetOptions(), tmplExtns.E_TemplateVariety) {
			continue
		} else if proto.HasExtension(fdp.GetOptions(), tmplExtns.E_TemplateName) && proto.HasExtension(fdp.GetOptions(), tmplExtns.E_TemplateVariety) {
			if templateDescriptorProto == nil {
				templateDescriptorProto = fdp
			} else {
				errors = multierror.Append(errors, fmt.Errorf("Proto files %s and %s, both have"+
					" the options %s and %s. Only one proto file is allowed with those options",
					fdp.Name, templateDescriptorProto.Name,
					tmplExtns.E_TemplateVariety.Name, tmplExtns.E_TemplateName.Name))

			}
		} else {
			errors = multierror.Append(errors, fmt.Errorf("Proto files %s has only one of the "+
				"following two options %s and %s. Both options are required.",
				fdp.Name, tmplExtns.E_TemplateVariety.Name, tmplExtns.E_TemplateName.Name))
		}
	}
	return templateDescriptorProto
}