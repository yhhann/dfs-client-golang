package(default_visibility = ["//:__subpackages__"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = glob(
        ["*.go"],
        exclude = ["*_test.go"],
    ),
    deps = [
        "//dfs-client-golang/instrument:go_default_library",
        "//dfs-client-golang/proto/discovery:go_default_library",
        "//dfs-client-golang/proto/transfer:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/codes:go_default_library",
        "//third-party-go/vendor/gopkg.in/mgo.v2/bson:go_default_library",
    ],
)
