package main

import (
	"testing"
	"os"
)

func test_normal(t *testing.T, input_file string, reference_file string) {
	test_compiler(t, input_file, reference_file, false)
}
func test_compiler(t *testing.T, input_file string, reference_file string, expect_errors bool) {
	text_file := strings.Replace(filepath.Base(input_file), filepath.Ext(input_file), ".text", 1)
	errors_file := strings.Replace(filepath.Base(input_file), filepath.Ext(input_file), ".errors", 1)
	// remove any preexisting output files
	os.Remove(text_file)
	os.Remove(errors_file)
	// run the compiler
	var err error
	var cmd = exec.Command(
		"gnostic",
		input_file,
		"--text-out=.",
		"--errors-out=.",
		"--resolve-refs")
	t.Log(cmd.Args)
	err = cmd.Run()
	if err != nil && !expect_errors {
		t.Logf("Compile failed: %+v", err)
		t.FailNow()
	}
	// verify the output against a reference
	var output_file string
	if expect_errors {
		output_file = errors_file
	} else {
		output_file = text_file
	}
	err = exec.Command("diff", output_file, reference_file).Run()
	if err != nil {
		t.Logf("Diff failed: %+v", err)
		t.FailNow()
	} else {
		// if the test succeeded, clean up
		os.Remove(text_file)
		os.Remove(errors_file)
	}
}

func TestMetricTemplate(t *testing.T) {
	test_normal(t,
		"testdata/MetricTemplate.proto",
		"testdata/MetricTemplate.proto.baseline")
}
