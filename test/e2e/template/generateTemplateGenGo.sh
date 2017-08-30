MIXERPATH=$GOPATH/src/istio.io/mixer
pushd $MIXERPATH

bazel build test/e2e/template/check/... test/e2e/template/report/... test/e2e/template/quota/...
bazel build tools/...

bazel-bin/tools/codegen/cmd/mixgenbootstrap/mixgenbootstrap \
bazel-genfiles/test/e2e/template/check/go_default_library_proto.descriptor_set:istio.io/mixer/test/e2e/template/check \
bazel-genfiles/test/e2e/template/quota/go_default_library_proto.descriptor_set:istio.io/mixer/test/e2e/template/quota \
bazel-genfiles/test/e2e/template/report/go_default_library_proto.descriptor_set:istio.io/mixer/test/e2e/template/report \
-o $MIXERPATH/test/e2e/template/template.gen.go
