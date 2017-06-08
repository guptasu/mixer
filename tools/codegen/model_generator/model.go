package model_generator

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
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

const FullNameOfExprMessage = "istio_mixer_v1_config_template.Expr"
const FullNameOfExprMessageWithPtr = "*" + FullNameOfExprMessage

// Creates a Model object
func CreateModel(parser *FileDescriptorSetParser) (Model, error) {
	model := &Model{}

	templateProto, err := getTmplFileDesc(parser.allFiles)
	if err != nil {
		return *model, err
	}

	// set the current package to the package of the templateProto.
	// This will make sure all references within the same file
	// do not get fully qualified.
	parser.packageName = goPackageName(templateProto.GetPackage())

	err = model.addTopLevelFields(templateProto)
	if err != nil {
		return *model, err
	}

	err = model.addInstanceFieldFromConstructor(parser, templateProto)
	if err != nil {
		return *model, err
	}

	err = model.getTypeNameForType(parser, templateProto)
	if err != nil {
		return *model, err
	}

	model.addImports(parser, templateProto)
	if err != nil {
		return *model, err
	}

	return *model, err
}

// Find the file that has the options TemplateVariety and TemplateName. There should only be one such file.
func getTmplFileDesc(fds []*FileDescriptor) (*FileDescriptor, error) {
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
				return nil, fmt.Errorf("Proto files %s and %s, both have"+
					" the options %s and %s. Only one proto file is allowed with those options",
					*fdp.Name, templateDescriptorProto.Name,
					tmplExtns.E_TemplateVariety.Name, tmplExtns.E_TemplateName.Name)

			}
		} else {
			return nil, fmt.Errorf("Proto files %s has only one of the "+
				"following two options %s and %s. Both options are required.",
				*fdp.Name, tmplExtns.E_TemplateVariety.Name, tmplExtns.E_TemplateName.Name)
		}
	}
	if templateDescriptorProto == nil {
		return nil, fmt.Errorf("There has to be one proto file that has both extensions %s and %s",
			tmplExtns.E_TemplateVariety.Name, tmplExtns.E_TemplateVariety.Name)
	}
	return templateDescriptorProto, nil
}

func (m *Model) addTopLevelFields(fdp *FileDescriptor) error {
	if fdp.Package == nil {
		return fmt.Errorf("package name missing on file %s", *fdp.Name)
	}
	m.PackageName = goPackageName(strings.TrimSpace(*fdp.Package))
	tmplName, err := proto.GetExtension(fdp.GetOptions(), tmplExtns.E_TemplateName)
	if err != nil {
		// This code should only get called for FileDescriptor that has this attribute. It is impossible to get to
		// this state.
		return fmt.Errorf("file option %s is required", tmplExtns.E_TemplateName.Name)
	}
	if name, ok := tmplName.(*string); !ok {
		// protoc should mandate the type. It is impossible to get to this state.
		return fmt.Errorf("%s should be of type string", tmplExtns.E_TemplateName.Name)
	} else {
		m.Name = *name
	}

	tmplVariety, err := proto.GetExtension(fdp.GetOptions(), tmplExtns.E_TemplateVariety)
	if err != nil {
		// This code should only get called for FileDescriptor that has this attribute. It is impossible to get to
		// this state.
		return fmt.Errorf("file option %s is required", tmplExtns.E_TemplateVariety.Name)
	}
	if tmplVariety == tmplExtns.TemplateVariety_TEMPLATE_VARIETY_CHECK {
		m.Check = true
		m.VarietyName = "Check"
	} else {
		m.Check = false
		m.VarietyName = "Report"
	}
	return nil
}

func (m *Model) getTypeNameForType(parser *FileDescriptorSetParser, fdp *FileDescriptor) error {
	typeDesc, err := getDescriptor(fdp, "Type")
	if err != nil {
		return err
	}
	m.TypeFullName = parser.TypeName(typeDesc)
	return nil
}

func getDescriptor(fdp *FileDescriptor, typeName string) (*Descriptor, error) {
	var cstrDesc *Descriptor = nil
	for _, desc := range fdp.desc {
		if *desc.Name == typeName {
			cstrDesc = desc
			break
		}
	}
	if cstrDesc == nil {
		return nil, fmt.Errorf("%s should have a message '%s'", *fdp.Name, typeName)
	}
	return cstrDesc, nil
}

// Build field information about the Constructor message.
func (m *Model) addInstanceFieldFromConstructor(parser *FileDescriptorSetParser, fdp *FileDescriptor) error {
	m.ConstructorFields = make([]FieldInfo, 0)
	cstrDesc, err := getDescriptor(fdp, "Constructor")
	if err != nil {
		return err
	}
	for _, fieldDesc := range cstrDesc.Field {
		fieldName := CamelCase(*fieldDesc.Name)
		typename := parser.GoType(cstrDesc.DescriptorProto, fieldDesc)
		// TODO : Can there be more than one expressions in a type for a field in Constructor ?
		typename = strings.Replace(typename, FullNameOfExprMessageWithPtr, "interface{}", 1)

		m.ConstructorFields = append(m.ConstructorFields, FieldInfo{Name: fieldName, Type: TypeInfo{Name: typename}})
	}
	return nil
}

func getUsedPackages(parser *FileDescriptorSetParser, desc *Descriptor) []string {
	result := make([]string, 0)
	for _, fieldDesc := range desc.Field {
		getUsedPackagesThisField(parser, fieldDesc, desc.PackageName(), &result)
	}

	return result
}

func getUsedPackagesThisField(parser *FileDescriptorSetParser, fieldDesc *descriptor.FieldDescriptorProto, parentPackage string, referencedPkgs *[]string) {
	if *fieldDesc.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE || *fieldDesc.Type == descriptor.FieldDescriptorProto_TYPE_ENUM {
		refDesc := parser.ObjectNamed(fieldDesc.GetTypeName())
		if d, ok := refDesc.(*Descriptor); ok {
			if fmt.Sprintf("%s.%s", d.PackageName(), d.GetName()) == FullNameOfExprMessage {
				return
			}
		}
		if d, ok := refDesc.(*Descriptor); ok && d.GetOptions().GetMapEntry() {
			keyField, valField := d.Field[0], d.Field[1]
			getUsedPackagesThisField(parser, keyField, parentPackage, referencedPkgs)
			getUsedPackagesThisField(parser, valField, parentPackage, referencedPkgs)
		} else {
			pname := goPackageName(parser.FileOf(refDesc.File()).GetPackage())
			if parentPackage != pname && !contains(*referencedPkgs, pname) {
				*referencedPkgs = append(*referencedPkgs, pname)
			}
		}
	}
}

func (m *Model) addImports(parser *FileDescriptorSetParser, fileDescriptor *FileDescriptor) {
	cstrDesc, _ := getDescriptor(fileDescriptor, "Constructor")
	if cstrDesc == nil {
		// should not happen because imports are supposed to be computed after all validation has occured and
		// all the types are loaded.
		return
	}

	usedPackages := getUsedPackages(parser, cstrDesc)
	m.Imports = make([]string, 0)
	for _, s := range fileDescriptor.Dependency {
		fd := parser.fileByName(s)
		// Do not import our own package.
		if fd.PackageName() == parser.packageName {
			continue
		}
		filename := fd.goFileName()
		// By default, import path is the dirname of the Go filename.
		importPath := path.Dir(filename)
		if substitution, ok := parser.ImportMap[s]; ok {
			importPath = substitution
		}

		pname := goPackageName(fd.GetPackage())
		if !contains(usedPackages, pname) {
			pname = "_"
		}
		m.Imports = append(m.Imports, pname+" "+strconv.Quote(importPath))
	}
}
