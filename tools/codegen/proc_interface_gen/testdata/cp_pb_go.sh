pushd ../../../../../
protoc mixer/tools/codegen/proc_interface_gen/testdata/QuotaTemplate.proto -I=. -I=api --go_out=Mgoogle/protobuf/duration.proto=github.com/golang/protobuf/ptypes/duration,Mmixer/v1/config/descriptor/value_type.proto=istio.io/api/mixer/v1/config/descriptor,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor,Mgoogle/protobuf/struct.proto=github.com/golang/protobuf/ptypes/struct,Mmixer/tools/codegen/template_extension/TemplateExtensions.proto=istio.io/mixer/tools/codegen/template_extension:.
mv mixer/tools/codegen/proc_interface_gen/testdata/QuotaTemplate.pb.go mixer/tools/codegen/proc_interface_gen/testdata/generated/TestQuotaTemplate

protoc mixer/tools/codegen/proc_interface_gen/testdata/MetricTemplate.proto -I=. -I=api --go_out=Mgoogle/protobuf/duration.proto=github.com/golang/protobuf/ptypes/duration,Mmixer/v1/config/descriptor/value_type.proto=istio.io/api/mixer/v1/config/descriptor,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor,Mgoogle/protobuf/struct.proto=github.com/golang/protobuf/ptypes/struct,Mmixer/tools/codegen/template_extension/TemplateExtensions.proto=istio.io/mixer/tools/codegen/template_extension:.
sed -i \
  -e 's/ValueType_VALUE_TYPE_UNSPECIFIED/VALUE_TYPE_UNSPECIFIED/g' \
  mixer/tools/codegen/proc_interface_gen/testdata/MetricTemplate.pb.go;
mv mixer/tools/codegen/proc_interface_gen/testdata/MetricTemplate.pb.go mixer/tools/codegen/proc_interface_gen/testdata/generated/TestMetricTemplate

protoc mixer/tools/codegen/proc_interface_gen/testdata/ListTemplate.proto -I=. -I=api --go_out=Mgoogle/protobuf/duration.proto=github.com/golang/protobuf/ptypes/duration,Mmixer/v1/config/descriptor/value_type.proto=istio.io/api/mixer/v1/config/descriptor,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor,Mgoogle/protobuf/struct.proto=github.com/golang/protobuf/ptypes/struct,Mmixer/tools/codegen/template_extension/TemplateExtensions.proto=istio.io/mixer/tools/codegen/template_extension:.
mv mixer/tools/codegen/proc_interface_gen/testdata/ListTemplate.pb.go mixer/tools/codegen/proc_interface_gen/testdata/generated/TestListTemplate

protoc mixer/tools/codegen/proc_interface_gen/testdata/LogTemplate.proto -I=. -I=api --go_out=Mgoogle/protobuf/duration.proto=github.com/golang/protobuf/ptypes/duration,Mmixer/v1/config/descriptor/value_type.proto=istio.io/api/mixer/v1/config/descriptor,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor,Mgoogle/protobuf/struct.proto=github.com/golang/protobuf/ptypes/struct,Mmixer/tools/codegen/template_extension/TemplateExtensions.proto=istio.io/mixer/tools/codegen/template_extension:.
mv mixer/tools/codegen/proc_interface_gen/testdata/LogTemplate.pb.go mixer/tools/codegen/proc_interface_gen/testdata/generated/TestLogTemplate
