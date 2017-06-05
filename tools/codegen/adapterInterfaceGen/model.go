package main

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"fmt"
)

type Model struct {
	Name string
}

func generateModel(fds *descriptor.FileDescriptorSet) (Model, error) {
	// TODO. Create a model for using the text tempaltes.
	model := Model{}
	for _, fdp := range fds.File {
		fmt.Println(fdp)
	}
	model.Name = "foo"
	return model, nil
}
