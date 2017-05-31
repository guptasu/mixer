package main

import (
	"testing"
	//"os"
	"path/filepath"
)

func test_normal(t *testing.T, input_file string, reference_file string) {
	test_compiler(t, input_file, reference_file, false)
}
func test_compiler(t *testing.T, input_file string, reference_file string, expect_errors bool) {

	outputDirFullPath := "testdata/generated"

	//var output_dir string
	if expect_errors {
		// output_file = errors_file
		// TODO
	} else {
		//output_dir = outputDirFullPath
	}
	// remove any preexisting output files
	//os.Remove(output_dir)

	fullPathOutDir, _ := filepath.Abs(outputDirFullPath)
	generator := Generator{outputDirFullPath:fullPathOutDir}

	generator.generate([]string {input_file})


	//os.Remove(errors_file)
	//err = exec.Command("diff", output_file, reference_file).Run()
	//if err != nil {
	//	t.Logf("Diff failed: %+v", err)
	//	t.FailNow()
	//} else {
	//	// if the test succeeded, clean up
	//	os.Remove(text_file)
	//	os.Remove(errors_file)
	//}
}

func TestMetricTemplate(t *testing.T) {
	test_normal(t,
		"testdata/MetricTemplate.proto",
		"testdata/MetricTemplate.proto.baseline")
}
