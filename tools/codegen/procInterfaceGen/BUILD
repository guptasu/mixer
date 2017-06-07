package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_prefix", "go_library", "go_test", "go_binary")



go_library(
    name = "go_default_library",
    srcs = [
        "generator.go"
    ],
    deps = [
        "@com_github_googleapis_googleapis//:google/rpc",
        "@com_github_hashicorp_go_multierror//:go_default_library",
        "@com_github_golang_protobuf//protoc-gen-go/descriptor:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
        "//tools/codegen/model_generator:go_default_library",
        "//tools/codegen/template_extension:go_default_library",
    ],
)

go_binary(
    name = "codegen",
    srcs = [
        "main.go",
    ],
    library = ":go_default_library",
    deps = [
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

# TODO: Undo this once I fiture out how to generate file descriptor set in the test file from a given proto file.
# Executing protoc from go code is tricky because of all the include dependencies.
#go_test(
#    name = "small_tests",
#    size = "small",
#    srcs = [
#        "generator_test.go",
#    ],
#    library = ":go_default_library",
#)
