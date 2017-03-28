package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "cnfgNormalizer.go",
        "jsFromExpression.go",
        "metricsAspectNormalizer.go",
    ],
    deps = [
        "//pkg/adapter:go_default_library",
        "//pkg/attribute:go_default_library",
        "//pkg/config/descriptor:go_default_library",
        "//pkg/config/proto:go_default_library",
        "//pkg/aspect/config:go_default_library",
        "//pkg/expr:go_default_library",
        "//pkg/aspect:go_default_library",
        "//pkg/config:go_default_library",
        "@com_github_robertkrimen_otto//:go_default_library",
    ],
)
