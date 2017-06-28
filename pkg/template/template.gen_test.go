package template

import (
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"

	sample_report "istio.io/mixer/pkg/template/sample/report"
)

func TestGetConstructorDefaultConfig(t *testing.T) {
	for _, tst := range []struct {
		template      string
		expectedProto proto.Message
		present       bool
	}{
		{sample_report.TemplateName, &sample_report.ConstructorParam{}, true},
		{"unknown template", nil, false},
	} {
		t.Run(tst.template, func(t *testing.T) {
			tdf := templateRepo{}
			k, rpresent := tdf.GetConstructorDefaultConfig(tst.template)
			if !reflect.DeepEqual(k, tst.expectedProto) || rpresent != tst.present {
				t.Errorf("GetConstructorDefaultConfig(%s) = %v,%v, want %v,%v", tst.template, k, rpresent,
					tst.expectedProto, tst.present)
			}
		})
	}
}
