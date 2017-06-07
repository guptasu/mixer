package main

import (
	"github.com/spf13/cobra"
	"os"
	"fmt"
	"path/filepath"
	"strings"
)

func withArgs(args []string, errorf func(format string, a ...interface{})) {
	var outFilePath string
	var mappings []string

	rootCmd := cobra.Command{
		Use:   "procInterfaceGen <File descriptor set protobuf>",
		Short: `
Tool that parses a [Template](http://TODO) and generates go interface for adapters to implement.

Example: procInterfaceGen metricTemplateFileDescriptorSet.pb -o MetricProcessor.go
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
			importMapping := make(map[string]string)
			for _, maps := range mappings {
				m := strings.Split(maps, ":")
				importMapping[m[0]] = m[1]
			}
			generator := Generator{outFilePath: outFileFullPath, importMapping: importMapping}
			if err := generator.generate(args[0]); err != nil {
				errorf("%v", err)
			}
		},
	}

	rootCmd.SetArgs(args)
	rootCmd.PersistentFlags().StringVarP(&outFilePath, "output", "o", "./generated.go", "Output " +
		"location for generating the go file.")

	rootCmd.PersistentFlags().StringArrayVarP(&mappings, "importmapping", "m", []string{}, "colon separated mapping of proto import to go package names." +
		" Example -m google/protobuf/descriptor.proto:github.com/golang/protobuf/protoc-gen-go/descriptor -m mixer/v1/config/descriptor/value_type.proto:istio.io/api/mixer/v1/config/descriptor")

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
