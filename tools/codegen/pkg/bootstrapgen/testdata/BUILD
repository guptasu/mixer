package(default_visibility = ["//tools/codegen/pkg/bootstrapgen:__pkg__"])

load("@org_pubref_rules_protobuf//protobuf:rules.bzl", "proto_compile")
load("//tools/codegen:generate.bzl", "mixer_proto_library")

mixer_proto_library(
    name = "list_template_library",
    testonly = True,
    protos = ["ListTemplate.proto"],
)

mixer_proto_library(
    name = "log_template_library",
    testonly = True,
    protos = ["LogTemplate.proto"],
)

mixer_proto_library(
    name = "metric_template_library",
    testonly = True,
    protos = ["MetricTemplate.proto"],
)

mixer_proto_library(
    name = "quota_template_library",
    testonly = True,
    protos = ["QuotaTemplate.proto"],
)

filegroup(
    name = "descriptors",
    testonly = True,
    srcs = [
        "list_template_library_proto.descriptor_set",
        "log_template_library_proto.descriptor_set",
        "metric_template_library_proto.descriptor_set",
        "quota_template_library_proto.descriptor_set",
    ],
)

filegroup(
    name = "golden_files",
    testonly = True,
    srcs = [
        "ListTemplate.golden.go",
        "LogTemplate.golden.go",
        "MetricTemplate.golden.go",
        "QuotaTemplate.golden.go",
        "AllTemplates.golden.go"
    ],
)

proto_compile(
    name = "error_template",
    testonly = True,
    args = [
        "--include_imports",
        "--include_source_info",
    ],
    imports = [
        "external/com_github_google_protobuf/src",
    ],
    inputs = [
        "//pkg/adapter/template:protos",
        "@com_github_google_protobuf//:well_known_protos",
    ],
    protos = ["ErrorTemplate.proto"],
)

proto_compile(
    name = "error_template_1",
    testonly = True,
    args = [
        "--include_imports",
        "--include_source_info",
    ],
    imports = [
        "external/com_github_google_protobuf/src",
    ],
    inputs = [
        "//pkg/adapter/template:protos",
        "@com_github_google_protobuf//:well_known_protos",
    ],
    protos = ["ErrorTemplate1.proto"],
)

proto_compile(
    name = "error_template_2",
    testonly = True,
    args = [
        "--include_imports",
        "--include_source_info",
    ],
    imports = [
        "external/com_github_google_protobuf/src",
    ],
    inputs = [
        "//pkg/adapter/template:protos",
        "@com_github_google_protobuf//:well_known_protos",
    ],
    protos = ["ErrorTemplate2.proto"],
)


filegroup(
    name = "error_test_files",
    testonly = True,
    srcs = [
        "ErrorTemplate.baseline",
        "error_template_1.descriptor_set",
        "error_template_2.descriptor_set",
    ],
)
