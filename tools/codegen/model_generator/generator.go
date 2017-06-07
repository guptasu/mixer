package model_generator

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

type FileDescriptorSetParser struct {
	typeNameToObject  map[string]Object // Key is a fully-qualified name in input syntax.
	Param             map[string]string // Command-line parameters.
	PackageImportPath string            // Go import path of the package we're generating code for
	ImportMap         map[string]string // Mapping from .proto file name to import path

	Pkg map[string]string // The names under which we import support packages

	packageName    string                     // What we're calling ourselves.
	allFiles       []*FileDescriptor          // All files in the tree
	allFilesByName map[string]*FileDescriptor // All files by filename.
	usedPackages   map[string]bool            // Names of packages used in current file.

}

type common struct {
	file *descriptor.FileDescriptorProto // File this object comes from.
}

func CreateFileDescriptorSetParser(fds *descriptor.FileDescriptorSet, importMap map[string]string) (*FileDescriptorSetParser, error) {
	parser := &FileDescriptorSetParser{ImportMap: importMap}
	parser.WrapTypes(fds)
	parser.BuildTypeNameMap()
	return parser, nil
}
func (g *FileDescriptorSetParser) fileByName(filename string) *FileDescriptor {
	return g.allFilesByName[filename]
}

func (g *FileDescriptorSetParser) GoType(message *descriptor.DescriptorProto, field *descriptor.FieldDescriptorProto) (typ string) {
	switch *field.Type {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		typ = "float64"
	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		typ = "float32"
	case descriptor.FieldDescriptorProto_TYPE_INT64:
		typ = "int64"
	case descriptor.FieldDescriptorProto_TYPE_UINT64:
		typ = "uint64"
	case descriptor.FieldDescriptorProto_TYPE_INT32:
		typ = "int32"
	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		typ = "uint32"
	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		typ = "uint64"
	case descriptor.FieldDescriptorProto_TYPE_FIXED32:
		typ = "uint32"
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		typ = "bool"
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		typ = "string"
	case descriptor.FieldDescriptorProto_TYPE_GROUP:
		//desc := g.ObjectNamed(field.GetTypeName())
		//typ = "*"+g.TypeName(desc), "group"
		// TODO : What needs to be done in this case? Is this allowed for templates
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		desc := g.ObjectNamed(field.GetTypeName())
		typ = "*" + g.TypeName(desc)

		if d, ok := desc.(*Descriptor); ok && d.GetOptions().GetMapEntry() {
			keyField, valField := d.Field[0], d.Field[1]
			keyType := g.GoType(d.DescriptorProto, keyField)
			valType := g.GoType(d.DescriptorProto, valField)

			keyType = strings.TrimPrefix(keyType, "*")
			switch *valField.Type {
			case descriptor.FieldDescriptorProto_TYPE_ENUM:
				valType = strings.TrimPrefix(valType, "*")

			case descriptor.FieldDescriptorProto_TYPE_MESSAGE:

			default:
				valType = strings.TrimPrefix(valType, "*")
			}

			typ = fmt.Sprintf("map[%s]%s", keyType, valType)
			return
		}
		//typ = field.GetTypeName()
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		typ = "[]byte"
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		//desc := g.ObjectNamed(field.GetTypeName())
		//typ = g.TypeName(desc)
		typ = field.GetTypeName()
	case descriptor.FieldDescriptorProto_TYPE_SFIXED32:
		typ = "int32"
	case descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		typ = "int64"
	case descriptor.FieldDescriptorProto_TYPE_SINT32:
		typ = "int32"
	case descriptor.FieldDescriptorProto_TYPE_SINT64:
		typ = "int64"
	default:
		// g.Fail("unknown type for", field.GetName())

	}
	if isRepeated(field) {
		typ = "[]" + typ
		//} else if message != nil && message.proto3() {
		//	return
	} else if field.OneofIndex != nil && message != nil {
		return
	} else if needsStar(*field.Type) {
		typ = "*" + typ
	}
	return
}

func (g *FileDescriptorSetParser) TypeName(obj Object) string {
	return g.DefaultPackageName(obj) + CamelCaseSlice(obj.TypeName())
}

func (g *FileDescriptorSetParser) DefaultPackageName(obj Object) string {
	// TODO if the protoc is not executed with --include_imports, this
	// is guaranteed to throw NPE.
	pkg := obj.PackageName()
	if pkg == g.packageName {
		return ""
	}
	return pkg + "."
}

type Object interface {
	PackageName() string // The name we use in our output (a_b_c), possibly renamed for uniqueness.
	TypeName() []string
	//File() *descriptor.FileDescriptorProto
}

func (c *common) PackageName() string {
	f := c.file
	return goPackageName(*f.Package)
}

// ObjectNamed, given a fully-qualified input type name as it appears in the input data,
// returns the descriptor for the message or enum with that name.
func (g *FileDescriptorSetParser) ObjectNamed(typeName string) Object {
	o, ok := g.typeNameToObject[typeName]
	if !ok {
		// TODO : g.Fail("can't find object with type", typeName)
	}

	return o
}

func (g *FileDescriptorSetParser) BuildTypeNameMap() {
	g.typeNameToObject = make(map[string]Object)
	for _, f := range g.allFiles {
		// The names in this loop are defined by the proto world, not us, so the
		// package name may be empty.  If so, the dotted package name of X will
		// be ".X"; otherwise it will be ".pkg.X".
		dottedPkg := "." + f.GetPackage()
		if dottedPkg != "." {
			dottedPkg += "."
		}
		for _, enum := range f.enum {
			name := dottedPkg + dottedSlice(enum.TypeName())
			g.typeNameToObject[name] = enum
		}
		for _, desc := range f.desc {
			name := dottedPkg + dottedSlice(desc.TypeName())
			g.typeNameToObject[name] = desc
		}
	}
}

func dottedSlice(elem []string) string { return strings.Join(elem, ".") }

func isRepeated(field *descriptor.FieldDescriptorProto) bool {
	return field.Label != nil && *field.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED
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

// Is c an ASCII lower-case letter?
func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func CamelCaseSlice(elem []string) string { return CamelCase(strings.Join(elem, "_")) }

func fileIsProto3(file *descriptor.FileDescriptorProto) bool {
	return file.GetSyntax() == "proto3"
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

func goPackageName(pkg string) string {
	return strings.Map(badToUnderscore, pkg)
}
