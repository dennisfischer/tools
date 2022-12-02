load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "index_builder_lib",
    srcs = ["index_builder.go"],
    importpath = "github.com/dennisfischer/tools/cmd",
    visibility = ["//visibility:private"],
    deps = ["//indexer"],
)

go_binary(
    name = "index_builder",
    embed = [":index_builder_lib"],
    visibility = ["//visibility:public"],
)

go_library(
    name = "duplicate_finder_lib",
    srcs = ["duplicate_finder.go"],
    importpath = "github.com/dennisfischer/tools/cmd",
    visibility = ["//visibility:private"],
    deps = ["//indexer"],
)

go_binary(
    name = "duplicate_finder",
    embed = [":duplicate_finder_lib"],
    visibility = ["//visibility:public"],
)

go_library(
    name = "duplicate_remover_lib",
    srcs = ["duplicate_remover.go"],
    importpath = "github.com/dennisfischer/tools/cmd",
    visibility = ["//visibility:private"],
    deps = [
        "//indexer",
        "//proto",
    ],
)

go_binary(
    name = "duplicate_remover",
    embed = [":duplicate_remover_lib"],
    visibility = ["//visibility:public"],
)

go_library(
    name = "empty_directory_remover_lib",
    srcs = ["empty_directory_remover.go"],
    importpath = "github.com/dennisfischer/tools/cmd",
    visibility = ["//visibility:private"],
)

go_binary(
    name = "empty_directory_remover",
    embed = [":empty_directory_remover_lib"],
    visibility = ["//visibility:public"],
)
