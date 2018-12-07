package(default_visibility = ["//dfs-client-golang:__subpackages__"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = glob(
        ["*.go"],
        exclude = ["*_test.go"],
    ),
    deps = [
        "//dfs-client-golang/handler:go_default_library",
        "//dfs-client-golang/loadbalance/hash:go_default_library",
        "//dfs-client-golang/proto/discovery:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
    ],
)
