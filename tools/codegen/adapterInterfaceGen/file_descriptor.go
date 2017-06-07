package main

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"path"
	"strings"
)

type FileDescriptor struct {
	*descriptor.FileDescriptorProto
	desc []*Descriptor          // All the messages defined in this file.
	enum []*EnumDescriptor      // All the enums defined in this file.

	proto3 bool // whether to generate proto3 code for this file
}



func (g *ModelGenerator) WrapTypes(fds *descriptor.FileDescriptorSet) {
	g.allFiles = make([]*FileDescriptor, 0, len(fds.File))
	g.allFilesByName = make(map[string]*FileDescriptor, len(g.allFiles))
	for _, f := range fds.File {
		g.WrapFileDescriptor(f)

	}
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

func (d *FileDescriptor) PackageName() string { return PackageName(*d.FileDescriptorProto.Name) }
