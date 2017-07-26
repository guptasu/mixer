package template

import (
	"github.com/gogo/protobuf/proto"

	"istio.io/mixer/pkg/template"

	"istio.io/mixer/template/list"

	"istio.io/mixer/template/metric"

	"istio.io/mixer/template/quota"

	"istio.io/mixer/template/log"
)

var (
	SupportedTmplInfo = map[string]template.Info{

		istio_mixer_template_list.TemplateName: {
			CtrCfg:   &istio_mixer_template_list.ConstructorParam{},
			BldrName: "istio.io/mixer/template/list.ListProcessorBuilder",
			SupportsTemplate: func(hndlrBuilder adptConfig.HandlerBuilder) bool {
				_, ok := hndlrBuilder.(istio_mixer_template_list.ListProcessorBuilder)
				return ok
			},
			InferType: func(cp proto.Message, tEvalFn template.TypeEvalFn) (proto.Message, error) {
				var err error
				cpb := cp.(*istio_mixer_template_list.ConstructorParam)
				infrdType := &istio_mixer_template_list.Type{}
				return infrdType, nil
			},
			ConfigureType: func(types map[string]proto.Message, builder *adptConfig.HandlerBuilder) error {

				castedBuilder := (*builder).(istio_mixer_adapter_sample_reportXXXXX.SampleProcessorBuilder)
				castedTypes := make(map[string]*istio_mixer_adapter_sample_reportXXXXX.Type)
				for k, v := range types {

					v1 := v.(*istio_mixer_adapter_sample_reportXXXXX.Type)
					castedTypes[k] = v1
				}
				return castedBuilder.ConfigureSample(castedTypes)
			},
		},

		istio_mixer_template_metric.TemplateName: {
			CtrCfg:   &istio_mixer_template_metric.ConstructorParam{},
			BldrName: "istio.io/mixer/template/metric.MetricProcessorBuilder",
			SupportsTemplate: func(hndlrBuilder adptConfig.HandlerBuilder) bool {
				_, ok := hndlrBuilder.(istio_mixer_template_metric.MetricProcessorBuilder)
				return ok
			},
			InferType: func(cp proto.Message, tEvalFn template.TypeEvalFn) (proto.Message, error) {
				var err error
				cpb := cp.(*istio_mixer_template_metric.ConstructorParam)
				infrdType := &istio_mixer_template_metric.Type{}
				return infrdType, nil
			},
			ConfigureType: func(types map[string]proto.Message, builder *adptConfig.HandlerBuilder) error {

				castedBuilder := (*builder).(istio_mixer_adapter_sample_reportXXXXX.SampleProcessorBuilder)
				castedTypes := make(map[string]*istio_mixer_adapter_sample_reportXXXXX.Type)
				for k, v := range types {

					v1 := v.(*istio_mixer_adapter_sample_reportXXXXX.Type)
					castedTypes[k] = v1
				}
				return castedBuilder.ConfigureSample(castedTypes)
			},
		},

		istio_mixer_template_quota.TemplateName: {
			CtrCfg:   &istio_mixer_template_quota.ConstructorParam{},
			BldrName: "istio.io/mixer/template/quota.QuotaProcessorBuilder",
			SupportsTemplate: func(hndlrBuilder adptConfig.HandlerBuilder) bool {
				_, ok := hndlrBuilder.(istio_mixer_template_quota.QuotaProcessorBuilder)
				return ok
			},
			InferType: func(cp proto.Message, tEvalFn template.TypeEvalFn) (proto.Message, error) {
				var err error
				cpb := cp.(*istio_mixer_template_quota.ConstructorParam)
				infrdType := &istio_mixer_template_quota.Type{}
				return infrdType, nil
			},
			ConfigureType: func(types map[string]proto.Message, builder *adptConfig.HandlerBuilder) error {

				castedBuilder := (*builder).(istio_mixer_adapter_sample_reportXXXXX.SampleProcessorBuilder)
				castedTypes := make(map[string]*istio_mixer_adapter_sample_reportXXXXX.Type)
				for k, v := range types {

					v1 := v.(*istio_mixer_adapter_sample_reportXXXXX.Type)
					castedTypes[k] = v1
				}
				return castedBuilder.ConfigureSample(castedTypes)
			},
		},

		istio_mixer_template_log.TemplateName: {
			CtrCfg:   &istio_mixer_template_log.ConstructorParam{},
			BldrName: "istio.io/mixer/template/log.LogProcessorBuilder",
			SupportsTemplate: func(hndlrBuilder adptConfig.HandlerBuilder) bool {
				_, ok := hndlrBuilder.(istio_mixer_template_log.LogProcessorBuilder)
				return ok
			},
			InferType: func(cp proto.Message, tEvalFn template.TypeEvalFn) (proto.Message, error) {
				var err error
				cpb := cp.(*istio_mixer_template_log.ConstructorParam)
				infrdType := &istio_mixer_template_log.Type{}
				return infrdType, nil
			},
			ConfigureType: func(types map[string]proto.Message, builder *adptConfig.HandlerBuilder) error {

				castedBuilder := (*builder).(istio_mixer_adapter_sample_reportXXXXX.SampleProcessorBuilder)
				castedTypes := make(map[string]*istio_mixer_adapter_sample_reportXXXXX.Type)
				for k, v := range types {

					v1 := v.(*istio_mixer_adapter_sample_reportXXXXX.Type)
					castedTypes[k] = v1
				}
				return castedBuilder.ConfigureSample(castedTypes)
			},
		},
	}
)
