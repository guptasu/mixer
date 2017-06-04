package main

import (
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"io/ioutil"
	"os"
	"text/template"
)

var templates = template.Must(template.ParseFiles(
	// TODO : Embed these files as resources.
	"text_templates/Interface.tmpl.go"))

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

	err = templates.ExecuteTemplate(f, "Interface.tmpl.go", model)
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

func generateModel(fileDescriptorSetPb *descriptor.FileDescriptorSet) (interface{}, error) {
	// TODO. Create a model for using the text tempaltes.
	return nil, nil
}
