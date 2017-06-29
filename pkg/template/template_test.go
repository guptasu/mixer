package template

import (
	"reflect"
	"testing"

	sample_report "istio.io/mixer/pkg/template/sample/report"
)

func TestGetTemplateInfo(t *testing.T) {
	for _, tst := range []struct {
		template string
		expected TemplateInfo
		present  bool
	}{
		{sample_report.TemplateName, TemplateInfo{&sample_report.ConstructorParam{}, inferTypeForSampleReport}, true},
		{"unknown template", TemplateInfo{}, false},
	} {
		t.Run(tst.template, func(t *testing.T) {
			tdf := templateRepo{}
			k, rpresent := tdf.GetTemplateInfo(tst.template)
			if rpresent != tst.present ||
				!reflect.DeepEqual(k.CnstrDefConfig, tst.expected.CnstrDefConfig) ||
				!reflect.DeepEqual(reflect.TypeOf(k.InferTypeFn), reflect.TypeOf(tst.expected.InferTypeFn)) {
				t.Errorf("GetConstructorDefaultConfig(%s) = %v,%v, want %v,%v", tst.template, k, rpresent,
					tst.expected, tst.present)
			}
		})
	}
}
