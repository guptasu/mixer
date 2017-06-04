package main

import (
	"github.com/spf13/cobra"
	"os"
	"fmt"
	"path/filepath"
)

func withArgs(args []string, errorf func(format string, a ...interface{})) {
	var outFilePath string

	rootCmd := cobra.Command{
		Use:   "adapterInterfaceGen [OPTIONS] <File descriptor protobuf>",
		Short: `
Tool that parses a [Template](http://TODO), defined within file descriptor set, and generates go interface
for adapters to implement.

Example: adapterInterfaceGen metricTemplateFileDescriptorSet.pb
`,

		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				errorf("Must specify a file descriptor set protobuf file.")
			}
			if len(args) != 1 {
				errorf("Only one input file is allowed.")
			}
			outFileFullPath, err := filepath.Abs(outFilePath)
			if err != nil {
				errorf("Invalid path %s. %v", outFilePath, err)
			}
			generator := Generator{outFilePath: outFileFullPath}
			if err := generator.generate(args[0]); err != nil {
				errorf("%v", err)
			}
		},
	}

	rootCmd.SetArgs(args)
	rootCmd.PersistentFlags().StringVarP(&outFilePath, "output", "o", "./generated.go", "Output " +
		"location for generating the go file.")
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
