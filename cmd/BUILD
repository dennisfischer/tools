load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "cmd_lib",
    srcs = ["main.go"],
    importpath = "github.com/dennisfischer/tools/cmd",
    visibility = ["//visibility:private"],
    deps = ["//proto"],
)

go_binary(
    name = "main",
    embed = [":cmd_lib"],
    visibility = ["//visibility:public"],
)
