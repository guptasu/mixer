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

func TestGetTypeInferFn(t *testing.T) {
	for _, tst := range []struct {
		template   string
		expectedFn InferTypeFn
		present    bool
	}{
		{sample_report.TemplateName, inferTypeForSampleReport, true},
		{"unknown template", nil, false},
	} {
		t.Run(tst.template, func(t *testing.T) {
			tdf := templateRepo{}
			k, rpresent := tdf.GetTypeInferFn(tst.template)
			if !reflect.DeepEqual(reflect.TypeOf(k), reflect.TypeOf(tst.expectedFn)) || rpresent != tst.present {
				t.Errorf("GetTypeInferFn(%s) = %v,%v, want %v,%v", tst.template, reflect.TypeOf(k), rpresent,
					reflect.TypeOf(tst.expectedFn), tst.present)
			}
		})
	}
}
