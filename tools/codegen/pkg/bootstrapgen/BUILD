package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "generator.go",
    ],
)

go_test(
    name = "small_tests",
    size = "small",
    srcs = [
        "generator_test.go",
    ],
    data = [
        "//tools/codegen/pkg/bootstrapgen/testdata:descriptors",
        "//tools/codegen/pkg/bootstrapgen/testdata:error_test_files",
        "//tools/codegen/pkg/bootstrapgen/testdata:golden_files",
    ],
    library = ":go_default_library",
)
