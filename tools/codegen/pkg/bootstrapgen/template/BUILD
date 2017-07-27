package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "bootstrapTemplate.go",
    ],
    deps = [
        "//pkg/adapter/config:go_default_library",
        "//pkg/expr:go_default_library",
        "//pkg/template:go_default_library",
        "//template/sample/report:go_default_library",
        "@com_github_gogo_protobuf//jsonpb:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@com_github_istio_api//:mixer/v1/config/descriptor",
    ],
)
