package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"testing"
	"path/filepath"
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
	outDir :=  path.Join(tmpOutDirContainer, t.Name())
	_,_ = filepath.Abs(outDir)
	err := os.RemoveAll(outDir)
	os.MkdirAll(outDir, os.ModePerm)

	outFDS := path.Join(outDir, "outFDS.pb")
	defer os.Remove(outFDS)
	err = generteFDSFileHacky(inputTemplateProto, outFDS)
	if err != nil {
		t.Logf("Unable to generate file descriptor set %v", err)
	}

	outFilePath := path.Join(outDir, "Processor.go")
	generator := Generator{outFilePath: outFilePath, importMapping:map[string]string {
		"mixer/v1/config/descriptor/value_type.proto":"istio.io/api/mixer/v1/config/descriptor",
		"mixer/tools/codegen/template_extension/TemplateExtensions.proto":"istio.io/mixer/tools/codegen/template_extension",
		"google/protobuf/duration.proto":"github.com/golang/protobuf/ptypes/duration",
	}}
	generator.generate(outFDS)

	/*
	// validate if the generated code builds
	// First copy all the
	protocCmd := []string{
		"build",
	}
	cmd := exec.Command("go", protocCmd...)
	cmd.Dir = absOutDir
	cmd.Stderr = os.Stderr // For debugging
	err = cmd.Run()

	if err != nil {
		//t.Logf("go build failed in dir %s: %+v.", absOutDir, err)
		t.FailNow()
		return
	}
	*/

	diffCmd := exec.Command("diff", outFilePath, expected, "--ignore-all-space")
	diffCmd.Stdout = os.Stdout
	diffCmd.Stderr = os.Stderr
	err = diffCmd.Run()
	if err != nil {
		t.Logf("Diff failed: %+v. Expected output is located at %s", err, outFilePath)
		t.FailNow()
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
		path.Join("mixer/tools/codegen/procInterfaceGen", protoFile),
		"-o",
		fmt.Sprintf("%s", path.Join("mixer/tools/codegen/procInterfaceGen", outputFDSFile)),
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
