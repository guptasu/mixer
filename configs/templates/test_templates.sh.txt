pushd $GOPATH/src/istio.io

## RUN PROTOC ON ALL THE TEMPLATES
protoc mixer/configs/mixer/*.proto -I=. -I=api --go_out=Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:.

protoc mixer/configs/templates/ListTemplate.proto -I=. -I=api --go_out=Mgoogle/protobuf/duration.proto=github.com/golang/protobuf/ptypes/duration,Mmixer/v1/config/descriptor/value_type.proto=istio.io/api/mixer/v1/config/descriptor,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor,Mgoogle/protobuf/struct.proto=github.com/golang/protobuf/ptypes/struct,Mmixer/configs/mixer/TemplateExtensions.proto=istio.io/mixer/configs/mixer:.
protoc mixer/configs/templates/QuotaTemplate.proto -I=. -I=api --go_out=Mgoogle/protobuf/duration.proto=github.com/golang/protobuf/ptypes/duration,Mmixer/v1/config/descriptor/value_type.proto=istio.io/api/mixer/v1/config/descriptor,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor,Mgoogle/protobuf/struct.proto=github.com/golang/protobuf/ptypes/struct,Mmixer/configs/mixer/TemplateExtensions.proto=istio.io/mixer/configs/mixer:.
protoc mixer/configs/templates/MetricTemplate.proto -I=. -I=api --go_out=Mgoogle/protobuf/duration.proto=github.com/golang/protobuf/ptypes/duration,Mmixer/v1/config/descriptor/value_type.proto=istio.io/api/mixer/v1/config/descriptor,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor,Mgoogle/protobuf/struct.proto=github.com/golang/protobuf/ptypes/struct,Mmixer/configs/mixer/TemplateExtensions.proto=istio.io/mixer/configs/mixer:.
protoc mixer/configs/templates/LogTemplate.proto -I=. -I=api --go_out=Mgoogle/protobuf/duration.proto=github.com/golang/protobuf/ptypes/duration,Mmixer/v1/config/descriptor/value_type.proto=istio.io/api/mixer/v1/config/descriptor,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor,Mgoogle/protobuf/struct.proto=github.com/golang/protobuf/ptypes/struct,Mmixer/configs/mixer/TemplateExtensions.proto=istio.io/mixer/configs/mixer:.

## COPY THE PROTOC OUTPUT TO SPECIFIC FOLDERS, SINCE TO BUILD GO CODE, NEED SAME PACKAGES IN ON DIR
mkdir mixer/configs/templates/list & mv mixer/configs/templates/ListTemplate.pb.go mixer/configs/templates/list
mkdir mixer/configs/templates/quota & mv mixer/configs/templates/QuotaTemplate.pb.go mixer/configs/templates/quota
mkdir mixer/configs/templates/metric & mv mixer/configs/templates/MetricTemplate.pb.go mixer/configs/templates/metric
mkdir mixer/configs/templates/log & mv mixer/configs/templates/LogTemplate.pb.go mixer/configs/templates/log

## SOME MAGIC TO MAKE COMPILER HAPPY
sed -i \
  -e 's/ValueType_VALUE_TYPE_UNSPECIFIED/VALUE_TYPE_UNSPECIFIED/g' \
  mixer/configs/templates/metric/MetricTemplate.pb.go;

## BUILD INDIVIDUAL GENERATED PROCESSORS

pushd mixer/configs/templates/generated_interfaces_for_adapters/list
go build
popd

pushd mixer/configs/templates/generated_interfaces_for_adapters/log
go build
popd

#pushd mixer/configs/templates/generated_interfaces_for_adapters/quota
#go build
#popd

pushd mixer/configs/templates/generated_interfaces_for_adapters/metric
go build
popd

popd
