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

package model_generator

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"

	"github.com/gogo/protobuf/proto"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

func TestNoPackageName(t *testing.T) {
	test(t,
		"testdata/NoPackageName.proto",
		"testdata/NoPackageName.baseline")
}

func test(t *testing.T, inputTemplateProto string, expected string) {

	tmpOutDirContainer := "testdata/generated"
	outDir := path.Join(tmpOutDirContainer, t.Name())
	_, _ = filepath.Abs(outDir)
	err := os.RemoveAll(outDir)
	os.MkdirAll(outDir, os.ModePerm)

	outFDS := path.Join(outDir, "outFDS.pb")
	defer os.Remove(outFDS)
	err = generteFDSFileHacky(inputTemplateProto, outFDS)
	if err != nil {
		t.Logf("Unable to generate file descriptor set %v", err)
		t.Fail()
		return
	}

	fds, err := getFileDescSet(outFDS)
	if err != nil {
		t.Logf("Unable to parse file descriptor set file %v", err)
		t.Fail()
	}

	//outFilePath := path.Join(outDir, "Processor.go")

	parser, err := CreateFileDescriptorSetParser(fds, make(map[string]string))
	model, err := CreateModel(parser)
	fmt.Println(model, err)
	//diffCmd := exec.Command("diff", outFilePath, expected, "--ignore-all-space")
	//diffCmd.Stdout = os.Stdout
	//diffCmd.Stderr = os.Stderr
	//err = diffCmd.Run()
	//if err != nil {
	//	t.Logf("Diff failed: %+v. Expected output is located at %s", err, outFilePath)
	//	t.FailNow()
	//	return
	//}
	//
	//// if the test succeeded, clean up
	//os.RemoveAll(outDir)
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

// TODO: This is blocking the test to be enabled from Bazel.
func generteFDSFileHacky(protoFile string, outputFDSFile string) error {

	// HACK HACK. Depending on dir structure is super fragile.
	// Explore how to generate File Descriptor set in a better way.
	protocCmd := []string{
		path.Join("mixer/tools/codegen/model_generator", protoFile),
		"-o",
		fmt.Sprintf("%s", path.Join("mixer/tools/codegen/model_generator", outputFDSFile)),
		"-I=.",
		"-I=api",
		"--include_imports",
	}
	cmd := exec.Command("protoc", protocCmd...)
	dir := path.Join(os.Getenv("GOPATH"), "src/istio.io")
	cmd.Dir = dir
	cmd.Stderr = os.Stderr // For debugging
	err := cmd.Run()
	return err
}
