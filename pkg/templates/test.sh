pushd ~/go/src/istio.io

protoc mixer/pkg/templates/*.proto -I=. -I=api --go_out=Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:.
protoc mixer/pkg/templates/metric/generated/config/*.proto -I=. -I=api --go_out=Mgoogle/protobuf/struct.proto=github.com/golang/protobuf/ptypes/struct,Mmixer/pkg/templates/TemplateExtensions.proto=istio.io/mixer/pkg/templates:.

sed -i \
  -e 's|mixer/v1/config/descriptor|istio.io/api/mixer/v1/config/descriptor|g' \
  -e 's/ValueType_VALUE_TYPE_UNSPECIFIED/VALUE_TYPE_UNSPECIFIED/g' \
  mixer/pkg/templates/metric/generated/config/MyNewMetric.pb.go;

pushd ~/go/src/istio.io/mixer/pkg/templates/metric/generated
go build
popd

pushd ~/go/src/istio.io/mixer/adapter/gRPCAdapter
go build
popd

popd


