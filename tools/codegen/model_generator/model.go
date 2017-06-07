package model_generator

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	multierror "github.com/hashicorp/go-multierror"
	tmplExtns "istio.io/mixer/tools/codegen/template_extension"
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

	// Constructor fields
	ConstructorFields []FieldInfo
}

type FieldInfo struct {
	Name string
	Type TypeInfo
}

type TypeInfo struct {
	Name string
}

const FullNameOfExprMessage = "*istio_mixer_v1_config_template.Expr"

// Creates a Model object
func CreateModel(parser *FileDescriptorSetParser) (Model, error) {
	model := &Model{}
	result := &multierror.Error{}

	templateProto := getTmplFileDesc(parser.allFiles, result)
	if len(result.Errors) != 0 {
		return *model, result.ErrorOrNil()
	}

	model.addTopLevelFields(templateProto, result)
	if len(result.Errors) != 0 {
		return *model, result.ErrorOrNil()
	}

	model.addFieldsOfConstructor(parser, templateProto, result)
	if len(result.Errors) != 0 {
		return *model, result.ErrorOrNil()
	}

	model.addImports(parser, templateProto)
	if len(result.Errors) != 0 {
		return *model, result.ErrorOrNil()
	}

	model.getTypeNameForType(parser, templateProto, result)
	if len(result.Errors) != 0 {
		return *model, result.ErrorOrNil()
	}

	return *model, result.ErrorOrNil()
}

// Find the file that has the options TemplateVariety and TemplateName. There should only be one such file.
func getTmplFileDesc(fds []*FileDescriptor, errors *multierror.Error) *FileDescriptor {
	var templateDescriptorProto *FileDescriptor = nil

	for _, fdp := range fds {
		if fdp.GetOptions() == nil {
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
					*fdp.Name, templateDescriptorProto.Name,
					tmplExtns.E_TemplateVariety.Name, tmplExtns.E_TemplateName.Name))

			}
		} else {
			errors = multierror.Append(errors, fmt.Errorf("Proto files %s has only one of the "+
				"following two options %s and %s. Both options are required.",
				*fdp.Name, tmplExtns.E_TemplateVariety.Name, tmplExtns.E_TemplateName.Name))
		}
	}
	if templateDescriptorProto == nil {
		errors = multierror.Append(errors, fmt.Errorf("There has to be one proto file that has both extensions %s and %s",
			tmplExtns.E_TemplateVariety.Name, tmplExtns.E_TemplateVariety.Name))
	}
	return templateDescriptorProto
}

func (m *Model) addTopLevelFields(fdp *FileDescriptor, errors *multierror.Error) {
	if fdp.Package == nil {
		errors = multierror.Append(errors, fmt.Errorf("package name missing on file %s", *fdp.Name))
		return
	}
	m.PackageName = goPackageName(strings.TrimSpace(*fdp.Package))
	tmplName, err := proto.GetExtension(fdp.GetOptions(), tmplExtns.E_TemplateName)
	if err != nil {
		errors = multierror.Append(errors, fmt.Errorf("file option %s is required", tmplExtns.E_TemplateName.Name))
		return
	}
	if name, ok := tmplName.(*string); !ok {
		errors = multierror.Append(errors, fmt.Errorf("%s should be of type string", tmplExtns.E_TemplateName.Name))
	} else {
		m.Name = *name
	}

	tmplVariety, err := proto.GetExtension(fdp.GetOptions(), tmplExtns.E_TemplateVariety)
	if err != nil {
		errors = multierror.Append(errors, fmt.Errorf("file option %s is required", tmplExtns.E_TemplateVariety.Name))
		return
	}
	if tmplVariety == tmplExtns.TemplateVariety_TEMPLATE_VARIETY_CHECK {
		m.Check = true
		m.VarietyName = "Check"
	} else {
		m.Check = false
		m.VarietyName = "Report"
	}
}

func (m *Model) getTypeNameForType(parser *FileDescriptorSetParser, fdp *FileDescriptor, errors *multierror.Error) {
	var typeDesc *Descriptor = nil
	for _, desc := range fdp.desc {
		if *desc.Name == "Type" {
			typeDesc = desc
			break
		}
	}
	if typeDesc == nil {
		errors = multierror.Append(errors, fmt.Errorf("%s should have a message 'Type'", *fdp.Name))
	}

	m.TypeFullName = parser.TypeName(typeDesc)
}

// Build field information about the Constructor message.
func (m *Model) addFieldsOfConstructor(parser *FileDescriptorSetParser, fdp *FileDescriptor, errors *multierror.Error) {
	m.ConstructorFields = make([]FieldInfo, 0)
	var cstrDesc *Descriptor = nil
	for _, desc := range fdp.desc {
		if *desc.Name == "Constructor" {
			cstrDesc = desc
			break
		}
	}
	if cstrDesc == nil {
		errors = multierror.Append(errors, fmt.Errorf("%s should have a message 'Constructor'", *fdp.Name))
	}

	for _, fieldDesc := range cstrDesc.Field {

		fieldName := CamelCase(*fieldDesc.Name)
		typename := parser.GoType(cstrDesc.DescriptorProto, fieldDesc)
		// TODO : Can there be more than one expressions in a type for a field in Constructor ?
		typename = strings.Replace(typename, FullNameOfExprMessage, "interface{}", 1)

		m.ConstructorFields = append(m.ConstructorFields, FieldInfo{Name: fieldName, Type: TypeInfo{Name: typename}})
	}
}

func (m *Model) addImports(g *FileDescriptorSetParser, fileDescriptor *FileDescriptor) {
	m.Imports = make([]string, 0)
	for _, s := range fileDescriptor.Dependency {
		fd := g.fileByName(s)
		// Do not import our own package.
		if fd.PackageName() == g.packageName {
			continue
		}
		filename := fd.goFileName()
		// By default, import path is the dirname of the Go filename.
		importPath := path.Dir(filename)
		if substitution, ok := g.ImportMap[s]; ok {
			importPath = substitution
		}

		pname := fd.PackageName()
		if _, ok := g.usedPackages[pname]; !ok {
			pname = "_"
		}
		m.Imports = append(m.Imports, pname+" "+strconv.Quote(importPath))
	}
}
