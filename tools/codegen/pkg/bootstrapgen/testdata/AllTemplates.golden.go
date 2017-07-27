package template

import (
	"github.com/gogo/protobuf/proto"

	pb "istio.io/api/mixer/v1/config/descriptor"
	"istio.io/mixer/pkg/template"

	"istio.io/mixer/template/list"

	"istio.io/mixer/template/log"

	"istio.io/mixer/template/metric"

	"istio.io/mixer/template/quota"
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

				infrdType.CheckExpression = istio_mixer_v1_config_descriptor.STRING

				return infrdType, nil
			},
			ConfigureType: func(types map[string]proto.Message, builder *adptConfig.HandlerBuilder) error {

				castedBuilder := (*builder).(istio_mixer_template_list.ListProcessorBuilder)
				castedTypes := make(map[string]*istio_mixer_template_list.Type)
				for k, v := range types {

					v1 := v.(*istio_mixer_template_list.Type)
					castedTypes[k] = v1
				}
				return castedBuilder.ConfigureList(castedTypes)
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

				infrdType.Dimensions = make(map[string]pb.ValueType)
				for k, v := range cpb.Dimensions {
					if infrdType.Dimensions[k], err = tEvalFn(v); err != nil {
						return nil, err
					}
				}

				return infrdType, nil
			},
			ConfigureType: func(types map[string]proto.Message, builder *adptConfig.HandlerBuilder) error {

				castedBuilder := (*builder).(istio_mixer_template_log.LogProcessorBuilder)
				castedTypes := make(map[string]*istio_mixer_template_log.Type)
				for k, v := range types {

					v1 := v.(*istio_mixer_template_log.Type)
					castedTypes[k] = v1
				}
				return castedBuilder.ConfigureLog(castedTypes)
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

				castedBuilder := (*builder).(istio_mixer_template_metric.MetricProcessorBuilder)
				castedTypes := make(map[string]*istio_mixer_template_metric.Type)
				for k, v := range types {

					v1 := v.(*istio_mixer_template_metric.Type)
					castedTypes[k] = v1
				}
				return castedBuilder.ConfigureMetric(castedTypes)
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

				infrdType.Dimensions = make(map[string]pb.ValueType)
				for k, v := range cpb.Dimensions {
					if infrdType.Dimensions[k], err = tEvalFn(v); err != nil {
						return nil, err
					}
				}

				return infrdType, nil
			},
			ConfigureType: func(types map[string]proto.Message, builder *adptConfig.HandlerBuilder) error {

				castedBuilder := (*builder).(istio_mixer_template_quota.QuotaProcessorBuilder)
				castedTypes := make(map[string]*istio_mixer_template_quota.Type)
				for k, v := range types {

					v1 := v.(*istio_mixer_template_quota.Type)
					castedTypes[k] = v1
				}
				return castedBuilder.ConfigureQuota(castedTypes)
			},
		},
	}
)
