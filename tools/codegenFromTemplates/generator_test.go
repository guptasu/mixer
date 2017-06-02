package main

import (
	"testing"
	"os"
	"path/filepath"
	"io/ioutil"
	"os/exec"
)

func TestMetricTemplate(t *testing.T) {
	test(t,
		"testdata/MetricTemplate.proto",
		"testdata/MetricTemplate.proto.baseline")
}

func test(t *testing.T, input_file string, reference_file string) {

	outputDirFullPath := "testdata/generated"
	fullPathOutDir, _ := filepath.Abs(outputDirFullPath)
	err := os.RemoveAll(fullPathOutDir)
	if err != nil {
		panic (err)
	}
	generator := Generator{outputDirFullPath:fullPathOutDir}

	generator.generate([]string {input_file})

	flattenedOutput := filepath.Join(outputDirFullPath, "flattenedOutput")
	flattenDirContentToAFile(fullPathOutDir, flattenedOutput)

	//os.Remove(errors_file)
	err = exec.Command("diff", flattenedOutput, reference_file).Run()
	if err != nil {
		t.Logf("Diff failed: %+v", err)
		t.FailNow()
	} else {
		// if the test succeeded, clean up
		os.RemoveAll(fullPathOutDir)
	}
}

func flattenDirContentToAFile(dir string, outFilePath string) error {
	searchDir := dir

	fileList := []string{}
	filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			fileList = append(fileList, path)
		}
		return nil

	})

	outFile, err := os.Create(outFilePath)
	if err != nil {
		return err
	}
	for _, file := range fileList {
		s, _ := ioutil.ReadFile(file)
		outFile.WriteString("********* " + filepath.Base(file) + " *********\n")
		outFile.Write(s)
		outFile.WriteString("********* END - " + filepath.Base(file) + " *********\n\n")
	}
	outFile.Close()

	return nil
}
