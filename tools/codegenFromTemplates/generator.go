package main

import (
	"text/template"
	"os"
	"path"
)

type Generator struct {
	outputDirFullPath string
}

type Data struct {
	TypeMessageName string
}

func (g *Generator) generate(args []string) error {
	// TODO : Embed these files as resources.
	tmpl, err := template.New("TemplateToProto.tmpl.proto").ParseFiles("TemplateToProto.tmpl.proto")
	if err != nil {
		panic(err)
	}
	file,err := os.Create(path.Join(g.outputDirFullPath, "SAMPLEGENERATED.proto"))
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(file, Data{TypeMessageName:"MetricTemplateParam"})
	if err != nil { panic(err) }

	return nil
}
