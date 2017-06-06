package main

import (
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"os"
	"text/template"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// TODO : Add this to a file resource instead of a string here.
var processorInterfaceGenTemplate,_ = template.New("processorInterfaceGenTemplate").Parse(
`
// WARNING !! CURRENTLY THIS IS HARD CODED AND NOT AUTOGENERATED !!

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

// THIS FILE IS AUTOMATICALLY GENERATED.

package {{.PackageName}}

import "istio.io/mixer/configs/templates/metric"

type Instance struct {
  {{range .ConstructorFields}}
  {{.Name}} {{.Type.Name}}
  {{end}}
}

type {{.Name}}Processor interface {
  Configure{{.Name}}(types map[string]*{{.PackageName}}.Type) error
  {{if .Check -}}
    {{- .VarietyName}}{{.Name}}(instances map[string]*Instance) (bool, error)
  {{else -}}
    {{- .VarietyName}}{{.Name}}(instances map[string]*Instance) (error)
  {{end}}
}
`)

type Generator struct {
	outFilePath string
}

func (g *Generator) generate(fileDescriptorProtobufFile string) error {
	f, err := os.Create(g.outFilePath)
	if err != nil {
		return err
	}

	fileDescriptorSetPb, err := getFileDescSetPb(fileDescriptorProtobufFile)
	if err != nil {
		return err
	}

	model, err := generateModel(fileDescriptorSetPb)
	if err != nil {
		return err
	}

	err = processorInterfaceGenTemplate.Execute(f, model)
	return err
}

func getFileDescSetPb(path string) (*descriptor.FileDescriptorSet, error) {
	byts, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	fds := &descriptor.FileDescriptorSet{}
	err = proto.Unmarshal(byts, fds)

	return fds, err
}
