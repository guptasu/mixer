package template

import (
	"fmt"

	"github.com/golang/protobuf/proto"

	pb "istio.io/api/mixer/v1/config/descriptor"
	sample_report "istio.io/mixer/pkg/template/sample/report"
)

var (
	templateInfos = map[string]Info{
		sample_report.TemplateName: {
			InferTypeFn:    inferTypeForSampleReport,
			CnstrDefConfig: &sample_report.ConstructorParam{},
		},
	}
)

func inferTypeForSampleReport(cp interface{}, tEvalFn TypeEvalFn) (proto.Message, error) {
	cpb := &sample_report.ConstructorParam{}
	var err error
	var ok bool

	if cpb, ok = cp.(*sample_report.ConstructorParam); !ok {
		return nil, fmt.Errorf("Constructor param %v is not of type %T", cp, cpb)
	}

	var infrdType = &sample_report.Type{}
	if infrdType.Value, err = tEvalFn(cpb.Value); err != nil {
		return nil, err
	}

	infrdType.Dimensions = make(map[string]pb.ValueType)
	for k, v := range cpb.Dimensions {
		if infrdType.Dimensions[k], err = tEvalFn(v); err != nil {
			return nil, err
		}
	}

	return infrdType, nil
}
