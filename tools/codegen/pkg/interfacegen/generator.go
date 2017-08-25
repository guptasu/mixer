// Copyright 2017 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package interfacegen

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"golang.org/x/tools/imports"

	tmpl "istio.io/mixer/tools/codegen/pkg/interfacegen/template"
	"istio.io/mixer/tools/codegen/pkg/modelgen"
)

// Generator generates Go interfaces for adapters to implement for a given Template.
type Generator struct {
	OutInterfacePath   string
	OAugmentedTmplPath string
	ImptMap            map[string]string
}

const (
	fullGoNameOfValueTypeEnum = "istio_mixer_v1_config_descriptor.ValueType"
)

func toProtoMap(k string, v string) string {
	return fmt.Sprintf("map<%s, %s>", k, v)
}

func stringify(protoType modelgen.TypeInfo) string {
	if protoType.IsMap {
		return toProtoMap(stringify(*protoType.MapKey), stringify(*protoType.MapValue))
	}
	return "string"
}

func containsValueType(ti modelgen.TypeInfo) bool {
	return ti.IsValueType || ti.IsMap && ti.MapValue.IsValueType
}

// Generate creates a Go interfaces for adapters to implement for a given Template.
func (g *Generator) Generate(fdsFile string) error {

	fds, err := getFileDescSet(fdsFile)
	if err != nil {
		return fmt.Errorf("cannot parse file '%s' as a FileDescriptorSetProto: %v", fdsFile, err)
	}

	parser, err := modelgen.CreateFileDescriptorSetParser(fds, g.ImptMap, "")
	if err != nil {
		return fmt.Errorf("cannot parse file '%s' as a FileDescriptorSetProto: %v", fdsFile, err)
	}

	model, err := modelgen.Create(parser)
	if err != nil {
		return err
	}

	interfaceFileContent, err := g.getInterfaceGoFileContent(model)
	if err != nil {
		return err
	}

	augProtoData, err := g.getAugmentedProtoContent(model)
	if err != nil {
		return err
	}

	// Everything succeeded, now write to the file.
	f1, err := os.Create(g.OutInterfacePath)
	if err != nil {
		return err
	}
	defer func() { _ = f1.Close() }() // nolint: gas

	if _, err = f1.Write(interfaceFileContent); err != nil { // nolint: gas
		_ = f1.Close()           // nolint: gas
		_ = os.Remove(f1.Name()) // nolint: gas
		return err
	}

	f2, err := os.Create(g.OAugmentedTmplPath)
	if err != nil {
		return err
	}
	defer func() { _ = f2.Close() }() // nolint: gas
	if _, err = f2.Write(augProtoData); err != nil {
		_ = f2.Close()           // nolint: gas
		_ = os.Remove(f2.Name()) // nolint: gas
		return err
	}

	return nil
}

const goFileImportFmt = "import \"%s\""

func (g *Generator) getInterfaceGoFileContent(model *modelgen.Model) ([]byte, error) {
	importsStms := make([]string, 0)
	intfaceTmpl, err := template.New("ProcInterface").Funcs(
		template.FuncMap{
			"replaceGoValueTypeToInterface": func(typeInfo modelgen.TypeInfo) string {
				return strings.Replace(typeInfo.Name, fullGoNameOfValueTypeEnum, "interface{}", 1)
			},
			// The text/templates have code logic using which it decides the fields to be printed. Example
			// when printing 'Type' we skip fields that have static types. So, this callback method 'reportTypeUsed'
			// allows the template code to register which fields and types it actually printed. Based on what was actually
			// printed we can decide which imports should be added to the file. Therefore, import adding is a last step
			// after all fields and messages / structs are printed.
			// The template's responsibility is to have a placeholder for printing the imports $$imports$$ and
			// the generator will replace it with imports for fields that were actually printed in the generated file.
			"reportTypeUsed": func(ti modelgen.TypeInfo) string {
				if len(ti.ImportNames) > 0 {
					for _, i := range ti.ImportNames {
						imptStm := fmt.Sprintf(goFileImportFmt, i)
						if !contains(importsStms, imptStm) {
							importsStms = append(importsStms, imptStm)
						}
					}
				}
				// do nothing, just record the import so that we can add them later (only for the types that got printed)
				return ""
			},
		}).Parse(tmpl.InterfaceTemplate)
	if err != nil {
		return nil, fmt.Errorf("cannot load template: %v", err)
	}
	intfaceBuf := new(bytes.Buffer)
	err = intfaceTmpl.Execute(intfaceBuf, model)
	if err != nil {
		return nil, fmt.Errorf("cannot execute the template with the given data: %v", err)
	}

	str := strings.Replace(string(intfaceBuf.Bytes()), "$$additional_imports$$", strings.Join(importsStms, "\n"), 1)

	fmtd, err := format.Source([]byte(str))
	if err != nil {
		return nil, fmt.Errorf("could not format generated code: %v : %s", err, string(intfaceBuf.Bytes()))
	}

	imports.LocalPrefix = "istio.io"
	// OutFilePath provides context for import path. We rely on the supplied bytes for content.
	imptd, err := imports.Process(g.OutInterfacePath, fmtd, nil)
	if err != nil {
		return nil, fmt.Errorf("could not fix imports for generated code: %v", err)
	}

	return []byte(str), nil

	return imptd, nil
}

const protoFileImportFmt = "import \"%s\";"
const protoValueTypeImport = "mixer/v1/config/descriptor/value_type.proto"

func (g *Generator) getAugmentedProtoContent(model *modelgen.Model) ([]byte, error) {
	imports := make([]string, 0)

	augmentedTemplateTmpl, err := template.New("AugmentedTemplateTmpl").Funcs(
		template.FuncMap{
			"containsValueType": containsValueType,
			"stringify":         stringify,
			// The text/templates have code logic using which it decides the fields to be printed. Example
			// when printing 'Type' we skip fields that have static types. So, this callback method 'reportTypeUsed'
			// allows the template code to register which fields and types it actually printed. Based on what was actually
			// printed we can decide which imports should be added to the file. Therefore, import adding is a last step
			// after all fields and messages / structs are printed.
			// The template's responsibility is to have a placeholder for printing the imports $$imports$$ and
			// the generator will replace it with imports for fields that were actually printed in the generated file.
			"reportTypeUsed": func(ti modelgen.TypeInfo) string {
				if len(ti.ImportNames) > 0 {
					for _, i := range ti.ImportNames {
						imptStm := fmt.Sprintf(protoFileImportFmt, i)
						if !contains(imports, imptStm) {
							imports = append(imports, imptStm)
						}
					}
				} else if containsValueType(ti) {
					imptStm := fmt.Sprintf(protoFileImportFmt, protoValueTypeImport)
					if !contains(imports, imptStm) {
						imports = append(imports, imptStm)
					}
				}
				// do nothing, just record the import so that we can add them later (only for the types that got printed)
				return ""
			},
		},
	).Parse(tmpl.RevisedTemplateTmpl)
	if err != nil {
		return nil, fmt.Errorf("cannot load template: %v", err)
	}

	tmplBuf := new(bytes.Buffer)
	err = augmentedTemplateTmpl.Execute(tmplBuf, model)
	if err != nil {
		return nil, fmt.Errorf("cannot execute the template with the given data: %v", err)
	}

	str := strings.Replace(string(tmplBuf.Bytes()), "$$additional_imports$$", strings.Join(imports, "\n"), 1)
	return []byte(str), nil
}

func getFileDescSet(path string) (*descriptor.FileDescriptorSet, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	fds := &descriptor.FileDescriptorSet{}
	err = proto.Unmarshal(bytes, fds)
	if err != nil {
		return nil, err
	}

	return fds, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
