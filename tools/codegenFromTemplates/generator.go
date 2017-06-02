package main

import (
	"os"
	"path"
	"text/template"
	"bytes"
)

var templates = template.Must(template.ParseFiles(
	// TODO : Embed these files as resources.
	"templates/TemplateToProto.tmpl.proto"))
var templateToOutFile = map[string]string {
	"TemplateToProto.tmpl.proto": "{{.TemplateShortName}}.Proto",
}

type Generator struct {
	outputDirFullPath string
}

func (g *Generator) createFile(genDirFullPath string, fileName string) (*os.File, error) {
	os.MkdirAll(genDirFullPath, os.ModePerm)
	file, err := os.Create(path.Join(genDirFullPath, fileName))
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (g *Generator) generate(args []string) error {
	models := parseTemplates(args)
	for _, model := range models {
		for k,v := range templateToOutFile {

			t, _ := template.New("tmp").Parse(v)
			buf := bytes.NewBufferString("")
			t.ExecuteTemplate(buf, "tmp", model)

			f, err := g.createFile(g.outputDirFullPath, buf.String())
			if err != nil {
				panic(err)
			}
			err = templates.ExecuteTemplate(f, k, model)
			if err != nil {
				panic(err)
			}
		}
	}

	return nil
}
