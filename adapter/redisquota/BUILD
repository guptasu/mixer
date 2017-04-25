package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "driver.go",
        "redisquota.go",
    ],
    deps = [
        "//adapter/memQuota/util:go_default_library",
        "//adapter/redisquota/config:go_default_library",
        "//pkg/adapter:go_default_library",
        "//pkg/pool:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_redis_client//pool:go_default_library",
        "@com_github_redis_client//redis:go_default_library",
    ],
)

