package main

import (
	"fmt"
	"text/template"
	"os"
)

type Generator struct {
	outputDirFullPath string
}

type Data struct {
	TypeMessageName string
}

func (g *Generator) generate(args []string) error {
	fmt.Println(args)
	fmt.Println("output dir is: ", g.outputDirFullPath)
	// TODO : Embed these files as resources.
	tmpl, err := template.New("TemplateToProto.tmpl.proto").ParseFiles("TemplateToProto.tmpl.proto")
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(os.Stdout, Data{TypeMessageName:"MetricTemplateParam"})
	if err != nil { panic(err) }

	return nil
}
