package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"testing"
)

func TestMetricTemplate(t *testing.T) {
	test(t,
		"testdata/MetricTemplate.proto",
		"testdata/MetricTemplateProcessorInterface.go.baseline")
}

func test(t *testing.T, inputTemplateProto string, expected string) {

	outDir := path.Join("testdata/generated", t.Name())
	err := os.RemoveAll(outDir)
	os.MkdirAll(outDir, os.ModePerm)

	outFDS := path.Join(outDir, "outFDS.pb")
	defer os.Remove(outFDS)
	err = generteFDSFileHacky(inputTemplateProto, outFDS)
	if err != nil {
		t.Logf("Unable to generate file descriptor set %v", err)
	}

	outFilePath := path.Join(outDir, "MetricTemplateProcessorInterface.go")
	generator := Generator{outFilePath: outFilePath}
	generator.generate(outFDS)

	err = exec.Command("diff", outFilePath, expected).Run()
	if err != nil {
		t.Logf("Diff failed: %+v. Expected output is located at %s", err, outFilePath)
		t.FailNow()
	} else {
		// if the test succeeded, clean up
		err := os.RemoveAll(outDir)
		os.Remove(outDir)
		if err != nil {
			panic (err)
		}
	}
}

// TODO: This is blocking the test to be enabled from Bazel.
func generteFDSFileHacky(protoFile string, outputFDSFile string) error {

	// HACK HACK. Depending on dir structure is super fragile.
	// Explore how to generate File Descriptor set in a better way.
	protocCmd := []string{
		path.Join("mixer/tools/codegen/adapterInterfaceGen", protoFile),
		"-o",
		fmt.Sprintf("%s", path.Join("mixer/tools/codegen/adapterInterfaceGen", outputFDSFile)),
		"-I=.",
		"-I=api",
	}
	cmd := exec.Command("protoc", protocCmd...)
	dir := path.Join(os.Getenv("GOPATH"), "src/istio.io")
	cmd.Dir = dir
	cmd.Stderr = os.Stderr // For debugging
	err := cmd.Run()
	return err
}
