package main


type Model struct {
	FileName string

	TypeMessageName string

	TemplateShortName string
}

func parseTemplate(fdsPbFilePaths string) []*Model {
	return []*Model{
		{FileName: "SAMPLEGENERATED.proto",TypeMessageName: "MetricTemplateParam", TemplateShortName:"MyMetric"},
	}
}

func parseTemplates(fdsPbFilePaths []string) []*Model {
	models := make([]*Model, 0)
	for _, fdsPbFilePath := range fdsPbFilePaths {
		models = append(models, parseTemplate(fdsPbFilePath)...)
	}

	return models
}
