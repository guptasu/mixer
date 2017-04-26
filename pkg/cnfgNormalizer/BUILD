package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "cnfgToJavascriptNormalizer.go",
        "expressionsToTypeScriptHelper.go",
        "metricsAspectConfigToTypeScript.go",
        "attribsToTypeScript.go",
        "NormalizedJavascriptConfigWithOtto.go",
        "NormalizedJavascriptConfigWithV8.go",
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
        #"@com_github_augustoroman_v8//:go_default_library",
        "@com_github_istio_api//:mixer/v1/config/descriptor",
    ],
)
