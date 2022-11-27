load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "indexer",
    srcs = ["indexer.go"],
    importpath = "github.com/dennisfischer/tools/indexer",
    visibility = ["//visibility:public"],
    deps = [
        "//proto",
        "@org_golang_google_protobuf//proto:go_default_library",
    ],
)
