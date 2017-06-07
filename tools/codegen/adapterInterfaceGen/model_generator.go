package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	multierror "github.com/hashicorp/go-multierror"
	tmplExtns "istio.io/mixer/tools/codegen/template_extension"
	"strings"
	"unicode"
	"path"
	"strconv"
)

type ModelGenerator struct {
	typeNameToObject map[string]Object          // Key is a fully-qualified name in input syntax.
	Param             map[string]string // Command-line parameters.
	PackageImportPath string            // Go import path of the package we're generating code for
	ImportMap         map[string]string // Mapping from .proto file name to import path

	Pkg map[string]string // The names under which we import support packages

	packageName      string                     // What we're calling ourselves.
	allFiles         []*FileDescriptor          // All files in the tree
	allFilesByName   map[string]*FileDescriptor // All files by filename.
	usedPackages     map[string]bool            // Names of packages used in current file.

	file             *FileDescriptor
}

type Descriptor struct {
	common
	*descriptor.DescriptorProto
	parent   *Descriptor            // The containing message, if any.
	nested   []*Descriptor          // Inner messages, if any.
	enums    []*EnumDescriptor      // Inner enums, if any.
	typename []string               // Cached typename vector.
	group    bool
}

type common struct {
	file *descriptor.FileDescriptorProto // File this object comes from.
}

type EnumDescriptor struct {
	common
	*descriptor.EnumDescriptorProto
	parent   *Descriptor // The containing message, if any.
	typename []string    // Cached typename vector.
}
type FileDescriptor struct {
	*descriptor.FileDescriptorProto
	desc []*Descriptor          // All the messages defined in this file.
	enum []*EnumDescriptor      // All the enums defined in this file.
	//imp  []*ImportedDescriptor  // All types defined in files publicly imported by this file.

	proto3 bool // whether to generate proto3 code for this file
}

func (g *ModelGenerator) WrapTypes(fds *descriptor.FileDescriptorSet) {
	g.allFiles = make([]*FileDescriptor, 0, len(fds.File))
	g.allFilesByName = make(map[string]*FileDescriptor, len(g.allFiles))
	for _, f := range fds.File {
		g.WrapFileDescriptor(f)

	}


	//for _, fd := range g.allFiles {
	//	g.allFiles = append(g.allFiles, fd)
	//	//fd.imp = wrapImported(fd.FileDescriptorProto, g)
	//}
	//g.genFiles = make([]*FileDescriptor, 0, len(g.Request.FileToGenerate))
	//for _, fileName := range g.Request.FileToGenerate {
	//	fd := g.allFilesByName[fileName]
	//	if fd == nil {
	//		g.Fail("could not find file named", fileName)
	//	}
	//	fd.index = len(g.genFiles)
	//	g.genFiles = append(g.genFiles, fd)
	//}
}

func (g *ModelGenerator) FileOf(fd *descriptor.FileDescriptorProto) *FileDescriptor {
	for _, file := range g.allFiles {
		if file.FileDescriptorProto == fd {
			return file
		}
	}
	// TODO: g.Fail("could not find file in table:", fd.GetName())
	return nil
}

func (g *ModelGenerator) WrapFileDescriptor(f *descriptor.FileDescriptorProto) {
	if _, ok := g.allFilesByName[f.GetName()]; !ok {
		// We must wrap the descriptors before we wrap the enums
		descs := wrapDescriptors(f)
		g.buildNestedDescriptors(descs)
		enums := wrapEnumDescriptors(f, descs)
		g.buildNestedEnums(descs, enums)
		fd := &FileDescriptor{
			FileDescriptorProto: f,
			desc:                descs,
			enum:                enums,
			proto3:              fileIsProto3(f),
		}
		g.allFiles = append(g.allFiles, fd)
		g.allFilesByName[f.GetName()] = fd
		//for _, deps := f.de
	}
}

//func wrapImported(file *descriptor.FileDescriptorProto, g *Generator) (sl []*ImportedDescriptor) {
//	for _, index := range file.PublicDependency {
//		df := g.fileByName(file.Dependency[index])
//		for _, d := range df.desc {
//			if d.GetOptions().GetMapEntry() {
//				continue
//			}
//			sl = append(sl, &ImportedDescriptor{common{file}, d})
//		}
//		for _, e := range df.enum {
//			sl = append(sl, &ImportedDescriptor{common{file}, e})
//		}
//		for _, ext := range df.ext {
//			sl = append(sl, &ImportedDescriptor{common{file}, ext})
//		}
//	}
//	return
//}

func (g *ModelGenerator) validate(fds *descriptor.FileDescriptorSet) (Model, error) {
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

func (g *ModelGenerator) buildNestedEnums(descs []*Descriptor, enums []*EnumDescriptor) {
	for _, desc := range descs {
		if len(desc.EnumType) != 0 {
			for _, enum := range enums {
				if enum.parent == desc {
					desc.enums = append(desc.enums, enum)
				}
			}
			if len(desc.enums) != len(desc.EnumType) {
				// TODO g.Fail("internal error: enum nesting failure for", desc.GetName())
			}
		}
	}
}

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
	//for _, fieldDesc := range cstrDesc.Field {

		//fieldName := CamelCase(*fieldDesc.Name)
		//typename := g.GoType(cstrDesc, fieldDesc)

		//model.ConstructorFields = append(model.ConstructorFields, FieldInfo{Name: fieldName, Type: TypeInfo{Name:typename}})
	//}
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

func wrapEnumDescriptors(file *descriptor.FileDescriptorProto, descs []*Descriptor) []*EnumDescriptor {
	sl := make([]*EnumDescriptor, 0, len(file.EnumType)+10)
	// Top-level enums.
	for i, enum := range file.EnumType {
		sl = append(sl, newEnumDescriptor(enum, nil, file, i))
	}
	// Enums within messages. Enums within embedded messages appear in the outer-most message.
	for _, nested := range descs {
		for i, enum := range nested.EnumType {
			sl = append(sl, newEnumDescriptor(enum, nested, file, i))
		}
	}
	return sl
}

func newEnumDescriptor(desc *descriptor.EnumDescriptorProto, parent *Descriptor, file *descriptor.FileDescriptorProto, index int) *EnumDescriptor {
	ed := &EnumDescriptor{
		common:          common{file},
		EnumDescriptorProto: desc,
		parent:              parent,
	}
	return ed
}

func (g *ModelGenerator) fileByName(filename string) *FileDescriptor {
	return g.allFilesByName[filename]
}

func (d *FileDescriptor) PackageName() string { return PackageName(*d.FileDescriptorProto.Name) }

func (g *ModelGenerator) generateImports() []string {
	imports := make([]string,0)
	for _, s := range g.file.Dependency {
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

		// We need to import all the dependencies, even if we don't reference them,
		// because other code and tools depend on having the full transitive closure
		// of protocol buffer types in the binary.
		pname := fd.PackageName()
		if _, ok := g.usedPackages[pname]; !ok {
			pname = "_"
		}
		imports = append(imports, pname + " " + strconv.Quote(importPath))
	}
	return imports
}

func (d *FileDescriptor) goFileName() string {
	name := *d.Name
	if ext := path.Ext(name); ext == ".proto" || ext == ".protodevel" {
		name = name[:len(name)-len(ext)]
	}
	name += ".pb.go"

	// Does the file have a "go_package" option?
	// If it does, it may override the filename.
	if impPath, _, ok := d.goPackageOption(); ok && impPath != "" {
		// Replace the existing dirname with the declared import path.
		_, name = path.Split(name)
		name = path.Join(impPath, name)
		return name
	}

	return name
}

// goPackageOption interprets the file's go_package option.
// If there is no go_package, it returns ("", "", false).
// If there's a simple name, it returns ("", pkg, true).
// If the option implies an import path, it returns (impPath, pkg, true).
func (d *FileDescriptor) goPackageOption() (impPath, pkg string, ok bool) {
	pkg = d.GetOptions().GetGoPackage()
	if pkg == "" {
		return
	}
	ok = true
	// The presence of a slash implies there's an import path.
	slash := strings.LastIndex(pkg, "/")
	if slash < 0 {
		return
	}
	impPath, pkg = pkg, pkg[slash+1:]
	// A semicolon-delimited suffix overrides the package name.
	sc := strings.IndexByte(impPath, ';')
	if sc < 0 {
		return
	}
	impPath, pkg = impPath[:sc], impPath[sc+1:]
	return
}


func (g *ModelGenerator) GoType(message *descriptor.DescriptorProto, field *descriptor.FieldDescriptorProto) (typ string) {
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
		 typ = "*"+ g.TypeName(desc)

		if d, ok := desc.(*Descriptor); ok && d.GetOptions().GetMapEntry() {
			keyField, valField := d.Field[0], d.Field[1]
			keyType:= g.GoType(d.DescriptorProto, keyField)
			valType:= g.GoType(d.DescriptorProto, valField)

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
func CamelCaseSlice(elem []string) string { return CamelCase(strings.Join(elem, "_")) }

func (g *ModelGenerator) TypeName(obj Object) string {
	return g.DefaultPackageName(obj) + CamelCaseSlice(obj.TypeName())
}

func (g *ModelGenerator) DefaultPackageName(obj Object) string {
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
	return PackageName(*f.Package)
}

// ObjectNamed, given a fully-qualified input type name as it appears in the input data,
// returns the descriptor for the message or enum with that name.
func (g *ModelGenerator) ObjectNamed(typeName string) Object {
	o, ok := g.typeNameToObject[typeName]
	if !ok {
		// TODO : g.Fail("can't find object with type", typeName)
	}

	// TODO
	/*
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
		if !found {
			log.Printf("protoc-gen-go: WARNING: failed finding publicly imported dependency for %v, used in %v", typeName, *g.file.Name)
		}
	}
	*/

	return o
}

func getTypeInfo(field *descriptor.FieldDescriptorProto) TypeInfo {
	return TypeInfo{Name:*field.TypeName}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
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
//func (g *Generator) ObjectNamed(typeName string) Object {
//	o, ok := g.typeNameToObject[typeName]
//	if !ok {
//		g.Fail("can't find object with type", typeName)
//	}
//
//	// If the file of this object isn't a direct dependency of the current file,
//	// or in the current file, then this object has been publicly imported into
//	// a dependency of the current file.
//	// We should return the ImportedDescriptor object for it instead.
//	direct := *o.File().Name == *g.file.Name
//	if !direct {
//		for _, dep := range g.file.Dependency {
//			if *g.fileByName(dep).Name == *o.File().Name {
//				direct = true
//				break
//			}
//		}
//	}
//	if !direct {
//		found := false
//	Loop:
//		for _, dep := range g.file.Dependency {
//			df := g.fileByName(*g.fileByName(dep).Name)
//			for _, td := range df.imp {
//				if td.o == o {
//					// Found it!
//					o = td
//					found = true
//					break Loop
//				}
//			}
//		}
//	}
//
//	return o
//}
//
func (g *ModelGenerator) BuildTypeNameMap() {
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

func (e *EnumDescriptor) TypeName() (s []string) {
	if e.typename != nil {
		return e.typename
	}
	name := e.GetName()
	if e.parent == nil {
		s = make([]string, 1)
	} else {
		pname := e.parent.TypeName()
		s = make([]string, len(pname)+1)
		copy(s, pname)
	}
	s[len(s)-1] = name
	e.typename = s
	return s
}

func (d *Descriptor) TypeName() []string {
	if d.typename != nil {
		return d.typename
	}
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

// Return a slice of all the Descriptors defined within this file
func wrapDescriptors(file *descriptor.FileDescriptorProto) []*Descriptor {
	sl := make([]*Descriptor, 0, len(file.MessageType)+10)
	for i, desc := range file.MessageType {
		sl = wrapThisDescriptor(sl, desc, nil, file, i)
	}
	return sl
}

// Wrap this Descriptor, recursively
func wrapThisDescriptor(sl []*Descriptor, desc *descriptor.DescriptorProto, parent *Descriptor, file *descriptor.FileDescriptorProto, index int) []*Descriptor {
	sl = append(sl, newDescriptor(desc, parent, file, index))
	me := sl[len(sl)-1]
	for i, nested := range desc.NestedType {
		sl = wrapThisDescriptor(sl, nested, me, file, i)
	}
	return sl
}

func (g *ModelGenerator) buildNestedDescriptors(descs []*Descriptor) {
	for _, desc := range descs {
		if len(desc.NestedType) != 0 {
			for _, nest := range descs {
				if nest.parent == desc {
					desc.nested = append(desc.nested, nest)
				}
			}
			if len(desc.nested) != len(desc.NestedType) {
				// TODO g.Fail("internal error: nesting failure for", desc.GetName())
			}
		}
	}
}

func newDescriptor(desc *descriptor.DescriptorProto, parent *Descriptor, file *descriptor.FileDescriptorProto, index int) *Descriptor {
	d := &Descriptor{
		common:          common{file},
		DescriptorProto: desc,
		parent:          parent,

	}

	// The only way to distinguish a group from a message is whether
	// the containing message has a TYPE_GROUP field that matches.
	if parent != nil {
		parts := d.TypeName()
		if file.Package != nil {
			parts = append([]string{*file.Package}, parts...)
		}
		exp := "." + strings.Join(parts, ".")
		for _, field := range parent.Field {
			if field.GetType() == descriptor.FieldDescriptorProto_TYPE_GROUP && field.GetTypeName() == exp {
				d.group = true
				break
			}
		}
	}

	return d
}

//
//func TypeName(d descriptor.DescriptorProto) []string {
//	n := 0
//	for parent := d; parent != nil; parent = parent.parent {
//		n++
//	}
//	s := make([]string, n, n)
//	for parent := d; parent != nil; parent = parent.parent {
//		n--
//		s[n] = parent.GetName()
//	}
//	d.typename = s
//	return s
//}
//
//func dottedSlice(elem []string) string { return strings.Join(elem, ".") }
//

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
