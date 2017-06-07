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

package main

import (
	"io/ioutil"
	"os"
	"text/template"

	"fmt"

	"path/filepath"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"istio.io/mixer/tools/codegen/model_generator"
)

type Generator struct {
	outFilePath   string
	importMapping map[string]string
}

func (g *Generator) Generate(fdsFile string) error {
	// This path works for bazel. TODO  Help pls !!
	tmplPathForBazel, _ := filepath.Abs("../../../../../../../../tools/codegen/procInterfaceGen/ProcInterface.tmpl")
	tmplPathForLocalIntellij, _ := filepath.Abs("../../../tools/codegen/procInterfaceGen/ProcInterface.tmpl")
	a, err := ioutil.ReadFile(tmplPathForBazel)
	if err != nil {
		a, err = ioutil.ReadFile(tmplPathForLocalIntellij)
	}

	if err != nil {
		panic(fmt.Errorf("cannot load template from path %s or %s", tmplPathForBazel, tmplPathForLocalIntellij))
	}

	tmpl, err := template.New("processorInterfaceGenTemplate").Parse(string(a))
	if err != nil {
		panic(err)
	}

	f, err := os.Create(g.outFilePath)
	if err != nil {
		return err
	}

	fds, err := getFileDescSet(fdsFile)
	if err != nil {
		return err
	}

	parser := &model_generator.FileDescriptorSetParser{ImportMap: g.importMapping}
	parser.WrapTypes(fds)
	parser.BuildTypeNameMap()
	model, err := parser.ConstructModel(fds)
	if err != nil {
		return err
	}

	return tmpl.Execute(f, model)
}

func getFileDescSet(path string) (*descriptor.FileDescriptorSet, error) {
	byts, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	fds := &descriptor.FileDescriptorSet{}
	err = proto.Unmarshal(byts, fds)

	return fds, err
}
