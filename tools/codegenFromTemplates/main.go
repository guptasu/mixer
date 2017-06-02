package main

import (
	"github.com/spf13/cobra"
	"os"
	"fmt"
	"path/filepath"
)

func withArgs(args []string, errorf func(format string, a ...interface{})) {
	var outputDir string

	rootCmd := cobra.Command{
		Use:   "codegenFromTemplates <File descriptor protobuf>...",
		Short: `
Tool that converts [Templates](http://TODO) defined in the file
descriptor sets and create protos and helper code for mixer's
service configuration model.

Example: codegenFromTemplates metricTemplateFileDescriptorSet.pb
NOTE: you will have to first generate metricTemplateFileDescriptorSet.pb using protoc
protoc -o metricTemplateFileDescriptorSet.pb metricTemplate.proto
`,

		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				errorf("Must specify at least one file descriptor set proto file.")
			}

			outputDirFullPath, err := filepath.Abs(outputDir)
			if err != nil {
				errorf("Invalid path %s. %v", outputDir, err)
			}
			generator := Generator{outputDirFullPath: outputDirFullPath}
			if err := generator.generate(args); err != nil {
				errorf("%v", err)
			}
		},
	}

	rootCmd.SetArgs(args)
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", "./generated", "Output " +
		"directory for generated code")
	if err := rootCmd.Execute(); err != nil {
		errorf("%v", err)
	}
}

func main() {
	withArgs(os.Args[1:],
		func(format string, a ...interface{}) {
			fmt.Fprintf(os.Stderr, format+"\n", a...)
			os.Exit(1)
		})
}
