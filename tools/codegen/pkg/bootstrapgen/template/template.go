// Copyright 2016 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package template

import (
	"github.com/golang/protobuf/proto"
	pb "istio.io/api/mixer/v1/config/descriptor"
	adptConfig "istio.io/mixer/pkg/adapter/config"
	template "istio.io/mixer/pkg/template"

	istio_mixer_adapter_sample_reportXXXXX "istio.io/mixer/template/sample/report" // XXXXXXXXXXX
)

var (
	SupportedTmplInfo = map[string]template.Info{
		istio_mixer_adapter_sample_reportXXXXX.TemplateName: {
			InferTypeFn:     inferTypeForXXXX,
			CnstrDefConfig:  &istio_mixer_adapter_sample_reportXXXXX.ConstructorParam{},
			ConfigureTypeFn: configureTypeForXXXX,
		},
	}
)

/////////////////////// Start generated code for template XXXX ///////////////////////
func supportsXXXXBuilder(hndlrBuilder adptConfig.HandlerBuilder) bool {
	_, ok := hndlrBuilder.(istio_mixer_adapter_sample_reportXXXXX.SampleProcessorBuilder)
	return ok
}

func inferTypeForXXXX(cp proto.Message, tEvalFn template.TypeEvalFn) (proto.Message, error) {
	var err error

	// Mixer framework should have ensured the type safety.
	cpb := cp.(*istio_mixer_adapter_sample_reportXXXXX.ConstructorParam)

	infrdType := &istio_mixer_adapter_sample_reportXXXXX.Type{}


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

func configureTypeForXXXX(types map[string]proto.Message, builder *adptConfig.HandlerBuilder) error {

	// Mixer framework should have ensured the type safety.
	castedBuilder := (*builder).(istio_mixer_adapter_sample_reportXXXXX.SampleProcessorBuilder)

	castedTypes := make(map[string]*istio_mixer_adapter_sample_reportXXXXX.Type)
	for k, v := range types {
		// Mixer framework should have ensured the type safety.
		v1 := v.(*istio_mixer_adapter_sample_reportXXXXX.Type)
		castedTypes[k] = v1
	}

	return castedBuilder.ConfigureSample(castedTypes)
}

/////////////////////// End generated code for template XXXX ///////////////////////
