pushd $GOPATH/src/istio.io > /dev/null;

## RUN PROTOC ON ALL THE TEMPLATES
protoc mixer/configs/mixer/*.proto -I=. -I=api --go_out=Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:.

protoc mixer/configs/templates/generated_interfaces_for_adapters/list/ListTemplate.gen.proto -I=. -I=api --go_out=Mgoogle/protobuf/duration.proto=github.com/golang/protobuf/ptypes/duration,Mmixer/v1/config/descriptor/value_type.proto=istio.io/api/mixer/v1/config/descriptor,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor,Mgoogle/protobuf/struct.proto=github.com/golang/protobuf/ptypes/struct,Mmixer/configs/mixer/TemplateExtensions.proto=istio.io/mixer/configs/mixer:.
protoc mixer/configs/templates/generated_interfaces_for_adapters/quota/QuotaTemplate.gen.proto -I=. -I=api --go_out=Mgoogle/protobuf/duration.proto=github.com/golang/protobuf/ptypes/duration,Mmixer/v1/config/descriptor/value_type.proto=istio.io/api/mixer/v1/config/descriptor,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor,Mgoogle/protobuf/struct.proto=github.com/golang/protobuf/ptypes/struct,Mmixer/configs/mixer/TemplateExtensions.proto=istio.io/mixer/configs/mixer:.
protoc mixer/configs/templates/generated_interfaces_for_adapters/metric/MetricTemplate.gen.proto -I=. -I=api --go_out=Mgoogle/protobuf/duration.proto=github.com/golang/protobuf/ptypes/duration,Mmixer/v1/config/descriptor/value_type.proto=istio.io/api/mixer/v1/config/descriptor,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor,Mgoogle/protobuf/struct.proto=github.com/golang/protobuf/ptypes/struct,Mmixer/configs/mixer/TemplateExtensions.proto=istio.io/mixer/configs/mixer:.
protoc mixer/configs/templates/generated_interfaces_for_adapters/log/LogTemplate.gen.proto -I=. -I=api --go_out=Mgoogle/protobuf/duration.proto=github.com/golang/protobuf/ptypes/duration,Mmixer/v1/config/descriptor/value_type.proto=istio.io/api/mixer/v1/config/descriptor,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor,Mgoogle/protobuf/struct.proto=github.com/golang/protobuf/ptypes/struct,Mmixer/configs/mixer/TemplateExtensions.proto=istio.io/mixer/configs/mixer:.

## SOME MAGIC TO MAKE COMPILER HAPPY
sed -i \
  -e 's/ValueType_VALUE_TYPE_UNSPECIFIED/VALUE_TYPE_UNSPECIFIED/g' \
  mixer/configs/templates/generated_interfaces_for_adapters/metric/MetricTemplate.gen.pb.go;

sed -i \
  -e 's/ValueType_VALUE_TYPE_UNSPECIFIED/VALUE_TYPE_UNSPECIFIED/g' \
  mixer/configs/templates/generated_interfaces_for_adapters/list/ListTemplate.gen.pb.go;

## BUILD INDIVIDUAL GENERATED PROCESSORS

DIRS="list quota metric log"

for pkgdir in ${DIRS}; do
    pushd mixer/configs/templates/generated_interfaces_for_adapters/${pkgdir} > /dev/null; \
    go build
    popd > /dev/null;
done

popd > /dev/null;

echo All generated protos and interface Go code builds.. Yay

