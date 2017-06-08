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
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"
)

func TestMetricTemplate(t *testing.T) {
	test(t,
		"testdata/MetricTemplate.proto",
		"testdata/MetricTemplateProcessorInterface.go.baseline")
}

func TestQuotaTemplate(t *testing.T) {
	test(t,
		"testdata/QuotaTemplate.proto",
		"testdata/QuotaTemplateProcessorInterface.go.baseline")
}

func TestLogTemplate(t *testing.T) {
	test(t,
		"testdata/LogTemplate.proto",
		"testdata/LogTemplateProcessorInterface.go.baseline")
}

func TestListTemplate(t *testing.T) {
	test(t,
		"testdata/ListTemplate.proto",
		"testdata/ListTemplateProcessorInterface.go.baseline")
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
	}

	outFilePath := path.Join(outDir, "Processor.go")
	generator := Generator{outFilePath: outFilePath, importMapping: map[string]string{
		"mixer/v1/config/descriptor/value_type.proto":                     "istio.io/api/mixer/v1/config/descriptor",
		"mixer/tools/codegen/template_extension/TemplateExtensions.proto": "istio.io/mixer/tools/codegen/template_extension",
		"google/protobuf/duration.proto":                                  "github.com/golang/protobuf/ptypes/duration",
	}}
	generator.Generate(outFDS)

	/*
	Below commented code is for testing if the generated code compiles correctly. Currently to test that, I have to
	run protoc separately copy the generated pb.go in the tmp output folder (doing it via a separate script),
	then uncomment the code and run the test. Need to find a cleaner automated way.
	*/

	// validate if the generated code builds
	// First copy all the
	//protocCmd := []string{
	//	"build",
	//}
	//cmd := exec.Command("go", protocCmd...)
	//cmd.Dir = absOutDir
	//cmd.Stderr = os.Stderr // For debugging
	//err = cmd.Run()
	//if err != nil {
	//	t.FailNow()
	//	return
	//}

	diffCmd := exec.Command("diff", outFilePath, expected, "--ignore-all-space")
	diffCmd.Stdout = os.Stdout
	diffCmd.Stderr = os.Stderr
	err = diffCmd.Run()
	if err != nil {
		t.Fatalf("Diff failed: %+v. Expected output is located at %s", err, outFilePath)
		return
	}

	// if the test succeeded, clean up
	os.RemoveAll(outDir)
}

// TODO: This is blocking the test to be enabled from Bazel.
func generteFDSFileHacky(protoFile string, outputFDSFile string) error {

	// HACK HACK. Depending on dir structure is super fragile.
	// Explore how to generate File Descriptor set in a better way.
	protocCmd := []string{
		path.Join("mixer/tools/codegen/proc_interface_gen", protoFile),
		"-o",
		fmt.Sprintf("%s", path.Join("mixer/tools/codegen/proc_interface_gen", outputFDSFile)),
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
