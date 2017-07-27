package template

import (
	"github.com/gogo/protobuf/proto"

	pb "istio.io/api/mixer/v1/config/descriptor"
	adptConfig "istio.io/mixer/pkg/adapter/config"
	"istio.io/mixer/pkg/template"

	"istio.io/mixer/template/sample/report"
)

var (
	SupportedTmplInfo = map[string]template.Info{

		istio_mixer_adapter_sample_report.TemplateName: {
			CtrCfg:   &istio_mixer_adapter_sample_report.ConstructorParam{},
			BldrName: "istio.io/mixer/template/sample/report.SampleProcessorBuilder",
			SupportsTemplate: func(hndlrBuilder adptConfig.HandlerBuilder) bool {
				_, ok := hndlrBuilder.(istio_mixer_adapter_sample_report.SampleProcessorBuilder)
				return ok
			},
			InferType: func(cp proto.Message, tEvalFn template.TypeEvalFn) (proto.Message, error) {
				var err error
				cpb := cp.(*istio_mixer_adapter_sample_report.ConstructorParam)
				infrdType := &istio_mixer_adapter_sample_report.Type{}

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
			},
			ConfigureType: func(types map[string]proto.Message, builder *adptConfig.HandlerBuilder) error {

				castedBuilder := (*builder).(istio_mixer_adapter_sample_report.SampleProcessorBuilder)
				castedTypes := make(map[string]*istio_mixer_adapter_sample_report.Type)
				for k, v := range types {

					v1 := v.(*istio_mixer_adapter_sample_report.Type)
					castedTypes[k] = v1
				}
				return castedBuilder.ConfigureSample(castedTypes)
			},
		},
	}
)