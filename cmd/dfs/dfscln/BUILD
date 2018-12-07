package(default_visibility = ["//dfs-client-test:__subpackages__"])

load("@io_bazel_rules_go//go:def.bzl", "go_binary")

go_binary(
  name = "dfscln",
  srcs = glob(["*.go"], exclude = ["*_test.go"]),
  deps = [
    "//dfs-client-golang/client:go_default_library",
    "//dfs-client-golang/handler:go_default_library",
    "//dfs-client-golang/proto/transfer:go_default_library",
    "//dfs-client-golang/proto/discovery:go_default_library",
    "//dfs-client-test/cmd:go_default_library",
    "//third-party-go/vendor/github.com/golang/glog:go_default_library",
  ],
  x_defs = {
    "main.buildTime": "2016/06/01_19:25",
  },
)
