package template

import (
	"github.com/golang/protobuf/proto"

	sample_report "istio.io/mixer/pkg/template/sample/report"
)

var templateConstructorParamMap = map[string]proto.Message{
	sample_report.TemplateName: &sample_report.ConstructorParam{},
}
