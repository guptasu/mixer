package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	multierror "github.com/hashicorp/go-multierror"
	tmplExtns "istio.io/mixer/tools/codegen/template_extension"
	"strings"
	"unicode"
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

func validate(fds *descriptor.FileDescriptorSet) (Model, error) {
	result := &multierror.Error{}
	model := &Model{}
	model.Imports = make([]string, 0)

	templateProto := getTemplateProto(fds, result)
	if len(result.Errors) != 0 {
		return *model, result.ErrorOrNil()
	}

	addTopLevelFields(model, templateProto, result)
	addFieldsOfConstructor(model, templateProto, result)
	model.TypeFullName = "XXXXMyType"
	return *model, result.ErrorOrNil()
}

type ModelGenerator struct {
	typeNameToObject map[string]interface{}          // Key is a fully-qualified name in input syntax.
}

func generateModel(fds *descriptor.FileDescriptorSet) (Model, error) {
	// TODO. Create a model for using the text tempaltes.
	return validate(fds)
}

// badToUnderscore is the mapping function used to generate Go names from package names,
// which can be dotted in the input .proto file.  It replaces non-identifier characters such as
// dot or dash with underscore.
func badToUnderscore(r rune) rune {
	if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
		return r
	}
	return '_'
}

func PackageName(pkg string) string {
	return strings.Map(badToUnderscore, pkg)
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

func addFieldsOfConstructor(model *Model, fdp *descriptor.FileDescriptorProto, errors *multierror.Error) {
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
		typename, wiretype := GoType(cstrDesc, fieldDesc)

		model.ConstructorFields = append(model.ConstructorFields, FieldInfo{Name: fieldName, Type: getTypeInfo(fieldDesc)})
	}
}

func GoType(message *descriptor.DescriptorProto, field *descriptor.FieldDescriptorProto) (typ string, wire string) {
	switch *field.Type {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		typ, wire = "float64", "fixed64"
	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		typ, wire = "float32", "fixed32"
	case descriptor.FieldDescriptorProto_TYPE_INT64:
		typ, wire = "int64", "varint"
	case descriptor.FieldDescriptorProto_TYPE_UINT64:
		typ, wire = "uint64", "varint"
	case descriptor.FieldDescriptorProto_TYPE_INT32:
		typ, wire = "int32", "varint"
	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		typ, wire = "uint32", "varint"
	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		typ, wire = "uint64", "fixed64"
	case descriptor.FieldDescriptorProto_TYPE_FIXED32:
		typ, wire = "uint32", "fixed32"
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		typ, wire = "bool", "varint"
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		typ, wire = "string", "bytes"
	case descriptor.FieldDescriptorProto_TYPE_GROUP:
		//desc := g.ObjectNamed(field.GetTypeName())
		//typ, wire = "*"+g.TypeName(desc), "group"
		// TODO : What needs to be done in this case? Is this allowed for templates
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		desc := g.ObjectNamed(field.GetTypeName())
		typ, wire = "*"+g.TypeName(desc), "bytes"
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		typ, wire = "[]byte", "bytes"
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		desc := g.ObjectNamed(field.GetTypeName())
		typ, wire = g.TypeName(desc), "varint"
	case descriptor.FieldDescriptorProto_TYPE_SFIXED32:
		typ, wire = "int32", "fixed32"
	case descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		typ, wire = "int64", "fixed64"
	case descriptor.FieldDescriptorProto_TYPE_SINT32:
		typ, wire = "int32", "zigzag32"
	case descriptor.FieldDescriptorProto_TYPE_SINT64:
		typ, wire = "int64", "zigzag64"
	default:
		g.Fail("unknown type for", field.GetName())
	}
	if isRepeated(field) {
		typ = "[]" + typ
	} else if message != nil && message.proto3() {
		return
	} else if field.OneofIndex != nil && message != nil {
		return
	} else if needsStar(*field.Type) {
		typ = "*" + typ
	}
	return
}

func getTypeInfo(field *descriptor.FieldDescriptorProto) TypeInfo {
	return TypeInfo{Name:*field.TypeName}
}

func getTemplateProto(fds *descriptor.FileDescriptorSet, errors *multierror.Error) *descriptor.FileDescriptorProto {
	var templateDescriptorProto *descriptor.FileDescriptorProto = nil
	for _, fdp := range fds.File {
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

func CamelCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && isASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}
// Is c an ASCII lower-case letter?
func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// ObjectNamed, given a fully-qualified input type name as it appears in the input data,
// returns the descriptor for the message or enum with that name.
func (g *Generator) ObjectNamed(typeName string) Object {
	o, ok := g.typeNameToObject[typeName]
	if !ok {
		g.Fail("can't find object with type", typeName)
	}

	// If the file of this object isn't a direct dependency of the current file,
	// or in the current file, then this object has been publicly imported into
	// a dependency of the current file.
	// We should return the ImportedDescriptor object for it instead.
	direct := *o.File().Name == *g.file.Name
	if !direct {
		for _, dep := range g.file.Dependency {
			if *g.fileByName(dep).Name == *o.File().Name {
				direct = true
				break
			}
		}
	}
	if !direct {
		found := false
	Loop:
		for _, dep := range g.file.Dependency {
			df := g.fileByName(*g.fileByName(dep).Name)
			for _, td := range df.imp {
				if td.o == o {
					// Found it!
					o = td
					found = true
					break Loop
				}
			}
		}
	}

	return o
}

func (g *ModelGenerator) BuildTypeNameMap(fds descriptor.FileDescriptorSet) {
	g.typeNameToObject = make(map[string]Object)
	for _, f := range fds.File {
		// The names in this loop are defined by the proto world, not us, so the
		// package name may be empty.  If so, the dotted package name of X will
		// be ".X"; otherwise it will be ".pkg.X".
		dottedPkg := "." + f.GetPackage()
		if dottedPkg != "." {
			dottedPkg += "."
		}
		for _, enum := range f.EnumType {
			name := dottedPkg + *enum.Name
			g.typeNameToObject[name] = enum
		}
		for _, desc := range f.MessageType {
			name := dottedPkg + *desc.Name
			g.typeNameToObject[name] = desc
		}
	}
}

func needsStar(typ descriptor.FieldDescriptorProto_Type) bool {
	switch typ {
	case descriptor.FieldDescriptorProto_TYPE_GROUP:
		return false
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		return false
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		return false
	}
	return true
}

func TypeName(d descriptor.DescriptorProto) []string {
	n := 0
	for parent := d; parent != nil; parent = parent.parent {
		n++
	}
	s := make([]string, n, n)
	for parent := d; parent != nil; parent = parent.parent {
		n--
		s[n] = parent.GetName()
	}
	d.typename = s
	return s
}

func dottedSlice(elem []string) string { return strings.Join(elem, ".") }

func isRepeated(field *descriptor.FieldDescriptorProto) bool {
	return field.Label != nil && *field.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED
}
