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

	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

func TestNoPackageName(t *testing.T) {
	testError(t,
		"testdata/NoPackageName.proto",
		"package name missing")
}

func TestMissingTemplateNameExt(t *testing.T) {
	testError(t,
		"testdata/MissingTemplateNameExt.proto",
		"has only one of the following two options")
}

func TestMissingTemplateVarietyExt(t *testing.T) {
	testError(t,
		"testdata/MissingTemplateVarietyExt.proto",
		"has only one of the following two options")
}

func TestMissingBothRequriedExt(t *testing.T) {
	testError(t,
		"testdata/MissingBothRequiredExt.proto",
		"one proto file that has both extensions")
}

func testError(t *testing.T, inputTemplateProto string, expectedError string) {

	outDir := path.Join("testdata", t.Name())
	_, _ = filepath.Abs(outDir)
	err := os.RemoveAll(outDir)
	os.MkdirAll(outDir, os.ModePerm)
	defer os.RemoveAll(outDir)

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

	parser, err := CreateFileDescriptorSetParser(fds, map[string]string{})
	_, err = CreateModel(parser)

	if !strings.Contains(err.Error(), expectedError) {
		t.Fatalf("CreateModel(%s) = %v, \n wanted err that contains string `%v`", inputTemplateProto, err, fmt.Errorf(expectedError))
	}
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
