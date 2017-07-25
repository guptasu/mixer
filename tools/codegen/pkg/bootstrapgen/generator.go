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

package bootstrapgen

import (
	"os"
	"text/template"
	tmplPkg "istio.io/mixer/tools/codegen/pkg/bootstrapgen/template"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"fmt"
	"io/ioutil"
	"github.com/gogo/protobuf/proto"
	"istio.io/mixer/tools/codegen/pkg/modelgen"
	"bytes"
	"golang.org/x/tools/imports"
	"go/format"
)

// Generator creates a Go file that will be build inside mixer framework. The generated file contains all the
// template specific code that mixer needs to add support for different passed in templates.
type Generator struct {
	OutFilePath   string
	ImportMapping map[string]string
}

// Generate creates a Go file that will be build inside mixer framework. The generated file contains all the
// template specific code that mixer needs to add support for different passed in templates.
func (g *Generator) Generate(fdsFiles []string) error {


	tmpl, err := template.New("ProcInterface").Parse(tmplPkg.InterfaceTemplate)
	if err != nil {
		return fmt.Errorf("cannot load template: %v", err)
	}

	models := make([]*modelgen.Model, 0)
	for _, fdsFile := range fdsFiles {
		fds, err := getFileDescSet(fdsFile)
		if err != nil {
			return fmt.Errorf("cannot parse file '%s' as a FileDescriptorSetProto. %v", fdsFile, err)
		}
		parser, err := modelgen.CreateFileDescriptorSetParser(fds, g.ImportMapping)
		if err != nil {
			return fmt.Errorf("cannot parse file '%s' as a FileDescriptorSetProto. %v", fdsFile, err)
		}

		model, err := modelgen.Create(parser)
		if err != nil {
			return err
		}
		models = append(models, model)
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, models)
	if err != nil {
		return fmt.Errorf("cannot execute the template with the given data: %v", err)
	}
	fmtd, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("could not format generated code: %v", err)
	}

	imports.LocalPrefix = "istio.io"
	// OutFilePath provides context for import path. We rely on the supplied bytes for content.
	imptd, err := imports.Process(g.OutFilePath, fmtd, nil)
	if err != nil {
		return fmt.Errorf("could not fix imports for generated code: %v", err)
	}

	f1, err := os.Create(g.OutFilePath)
	if err != nil {
		return err
	}
	defer func() { _ = f1.Close() }()
	if _, err = f1.Write(imptd); err != nil {
		_ = f1.Close()
		_ = os.Remove(f1.Name())
		return err
	}

	return nil
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
