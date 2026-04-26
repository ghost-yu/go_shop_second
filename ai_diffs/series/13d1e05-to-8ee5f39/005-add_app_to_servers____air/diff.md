# Commit Diff Report

- Repo: go_shop_second
- Sequence: 005 / 10
- Commit: 49f4e56b544b754e388515fdd7f8b15bf1341e20
- ShortCommit: 49f4e56
- Parent: 49bfa8eb40f86817e6bb4f6af9d9f628070e4660
- Subject: add app to servers && air
- Author: ghost-yu <hgfhgfhgfhgfhgfhgf@yeah.net>
- Date: 2024-10-14 00:49:20 +0800
- GeneratedAt: 2026-04-06 17:43:36 +08:00

## Short Summary

~~~text
 15 files changed, 298 insertions(+), 21 deletions(-)
~~~

## File Stats

~~~text
 .gitignore                                   |  8 +++-
 api/orderpb/order.proto                      |  2 +-
 api/stockpb/stock.proto                      |  2 +-
 internal/common/genproto/orderpb/order.pb.go |  9 ++--
 internal/common/genproto/stockpb/stock.pb.go |  9 ++--
 internal/common/go.mod                       | 23 ++++++++++-
 internal/common/go.sum                       | 61 ++++++++++++++++++++++++++++
 internal/kitchen/.air.toml                   | 52 ++++++++++++++++++++++++
 internal/order/.air.toml                     | 52 ++++++++++++++++++++++++
 internal/order/app/app.go                    | 10 +++++
 internal/order/http.go                       |  5 ++-
 internal/order/main.go                       | 17 ++++++--
 internal/order/ports/grpc.go                 |  6 ++-
 internal/order/service/application.go        | 11 +++++
 internal/payment/.air.toml                   | 52 ++++++++++++++++++++++++
 15 files changed, 298 insertions(+), 21 deletions(-)
~~~

## Changed Files

~~~text
.gitignore
api/orderpb/order.proto
api/stockpb/stock.proto
internal/common/genproto/orderpb/order.pb.go
internal/common/genproto/stockpb/stock.pb.go
internal/common/go.mod
internal/common/go.sum
internal/kitchen/.air.toml
internal/order/.air.toml
internal/order/app/app.go
internal/order/http.go
internal/order/main.go
internal/order/ports/grpc.go
internal/order/service/application.go
internal/payment/.air.toml
~~~

## Focus Files (Excluded: go.mod / go.sum)

~~~text
.gitignore
api/orderpb/order.proto
api/stockpb/stock.proto
internal/common/genproto/orderpb/order.pb.go
internal/common/genproto/stockpb/stock.pb.go
internal/kitchen/.air.toml
internal/order/.air.toml
internal/order/app/app.go
internal/order/http.go
internal/order/main.go
internal/order/ports/grpc.go
internal/order/service/application.go
internal/payment/.air.toml
~~~

## Patch

~~~diff
diff --git a/.gitignore b/.gitignore
index 42e4918..3c67536 100644
--- a/.gitignore
+++ b/.gitignore
@@ -23,4 +23,10 @@ go.work.sum
 
 # env file
 .env
-.idea
\ No newline at end of file
+.idea
+.DS_Store
+
+.gobincache/
+
+**/tmp/
+**/bin/
\ No newline at end of file
diff --git a/api/orderpb/order.proto b/api/orderpb/order.proto
index 24b4e62..af9b2ba 100644
--- a/api/orderpb/order.proto
+++ b/api/orderpb/order.proto
@@ -1,7 +1,7 @@
 syntax = "proto3";
 package orderpb;
 
-option go_package = "github.com/ghost-yu/go_shop_second/internal/common/genproto/orderpb";
+option go_package = "github.com/ghost-yu/go_shop_second/common/genproto/orderpb";
 
 import "google/protobuf/empty.proto";
 
diff --git a/api/stockpb/stock.proto b/api/stockpb/stock.proto
index 0ef749a..4275b40 100644
--- a/api/stockpb/stock.proto
+++ b/api/stockpb/stock.proto
@@ -1,7 +1,7 @@
 syntax = "proto3";
 package stockpb;
 
-option go_package = "github.com/ghost-yu/go_shop_second/internal/common/genproto/stockpb";
+option go_package = "github.com/ghost-yu/go_shop_second/common/genproto/stockpb";
 
 import "orderpb/order.proto";
 
diff --git a/internal/common/genproto/orderpb/order.pb.go b/internal/common/genproto/orderpb/order.pb.go
index f016711..7d1f1fe 100644
--- a/internal/common/genproto/orderpb/order.pb.go
+++ b/internal/common/genproto/orderpb/order.pb.go
@@ -375,12 +375,11 @@ var file_orderpb_order_proto_rawDesc = []byte{
 	0x65, 0x72, 0x12, 0x35, 0x0a, 0x0b, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4f, 0x72, 0x64, 0x65,
 	0x72, 0x12, 0x0e, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x70, 0x62, 0x2e, 0x4f, 0x72, 0x64, 0x65,
 	0x72, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
-	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x44, 0x5a, 0x42, 0x67, 0x69, 0x74,
+	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x3b, 0x5a, 0x39, 0x67, 0x69, 0x74,
 	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x4e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65,
-	0x7a, 0x7a, 0x30, 0x30, 0x2f, 0x67, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2d, 0x76, 0x32, 0x2f, 0x69,
-	0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x67,
-	0x65, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x70, 0x62, 0x62,
-	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
+	0x7a, 0x7a, 0x30, 0x30, 0x2f, 0x67, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2d, 0x76, 0x32, 0x2f, 0x63,
+	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x67, 0x65, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6f,
+	0x72, 0x64, 0x65, 0x72, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
 }
 
 var (
diff --git a/internal/common/genproto/stockpb/stock.pb.go b/internal/common/genproto/stockpb/stock.pb.go
index 09d4979..9566739 100644
--- a/internal/common/genproto/stockpb/stock.pb.go
+++ b/internal/common/genproto/stockpb/stock.pb.go
@@ -7,7 +7,7 @@
 package stockpb
 
 import (
-	orderpb "github.com/ghost-yu/go_shop_second/internal/common/genproto/orderpb"
+	orderpb "github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
 	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
 	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
 	reflect "reflect"
@@ -251,12 +251,11 @@ var file_stockpb_stock_proto_rawDesc = []byte{
 	0x66, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x49, 0x6e, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x52, 0x65, 0x71,
 	0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x70, 0x62, 0x2e, 0x43,
 	0x68, 0x65, 0x63, 0x6b, 0x49, 0x66, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x49, 0x6e, 0x53, 0x74, 0x6f,
-	0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x44, 0x5a, 0x42, 0x67, 0x69,
+	0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x3b, 0x5a, 0x39, 0x67, 0x69,
 	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x4e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d,
 	0x65, 0x7a, 0x7a, 0x30, 0x30, 0x2f, 0x67, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2d, 0x76, 0x32, 0x2f,
-	0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f,
-	0x67, 0x65, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x70, 0x62,
-	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
+	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x67, 0x65, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
+	0x73, 0x74, 0x6f, 0x63, 0x6b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
 }
 
 var (
diff --git a/internal/common/go.mod b/internal/common/go.mod
index 546aaa5..8355ff8 100644
--- a/internal/common/go.mod
+++ b/internal/common/go.mod
@@ -3,7 +3,10 @@ module github.com/ghost-yu/go_shop_second/common
 go 1.22.8
 
 require (
+	github.com/gin-gonic/gin v1.9.1
+	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
 	github.com/oapi-codegen/runtime v1.1.1
+	github.com/sirupsen/logrus v1.8.1
 	github.com/spf13/viper v1.19.0
 	google.golang.org/grpc v1.62.1
 	google.golang.org/protobuf v1.33.0
@@ -11,13 +14,27 @@ require (
 
 require (
 	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
+	github.com/bytedance/sonic v1.10.0-rc3 // indirect
+	github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d // indirect
+	github.com/chenzhuoyu/iasm v0.9.0 // indirect
 	github.com/fsnotify/fsnotify v1.7.0 // indirect
+	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
+	github.com/gin-contrib/sse v0.1.0 // indirect
+	github.com/go-playground/locales v0.14.1 // indirect
+	github.com/go-playground/universal-translator v0.18.1 // indirect
+	github.com/go-playground/validator/v10 v10.14.1 // indirect
+	github.com/goccy/go-json v0.10.2 // indirect
 	github.com/golang/protobuf v1.5.3 // indirect
 	github.com/google/uuid v1.6.0 // indirect
-	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
 	github.com/hashicorp/hcl v1.0.0 // indirect
+	github.com/json-iterator/go v1.1.12 // indirect
+	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
+	github.com/leodido/go-urn v1.2.4 // indirect
 	github.com/magiconair/properties v1.8.7 // indirect
+	github.com/mattn/go-isatty v0.0.20 // indirect
 	github.com/mitchellh/mapstructure v1.5.0 // indirect
+	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
+	github.com/modern-go/reflect2 v1.0.2 // indirect
 	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
 	github.com/sagikazarmark/locafero v0.4.0 // indirect
 	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
@@ -26,8 +43,12 @@ require (
 	github.com/spf13/cast v1.6.0 // indirect
 	github.com/spf13/pflag v1.0.5 // indirect
 	github.com/subosito/gotenv v1.6.0 // indirect
+	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
+	github.com/ugorji/go/codec v1.2.11 // indirect
 	go.uber.org/atomic v1.9.0 // indirect
 	go.uber.org/multierr v1.9.0 // indirect
+	golang.org/x/arch v0.4.0 // indirect
+	golang.org/x/crypto v0.21.0 // indirect
 	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
 	golang.org/x/net v0.23.0 // indirect
 	golang.org/x/sys v0.18.0 // indirect
diff --git a/internal/common/go.sum b/internal/common/go.sum
index 806458e..a404a00 100644
--- a/internal/common/go.sum
+++ b/internal/common/go.sum
@@ -5,7 +5,17 @@ github.com/apapsch/go-jsonmerge/v2 v2.0.0 h1:axGnT1gRIfimI7gJifB699GoE/oq+F2MU7D
 github.com/apapsch/go-jsonmerge/v2 v2.0.0/go.mod h1:lvDnEdqiQrp0O42VQGgmlKpxL1AP2+08jFMw88y4klk=
 github.com/benbjohnson/clock v1.1.0/go.mod h1:J11/hYXuz8f4ySSvYwY0FKfm+ezbsZBKZxNJlLklBHA=
 github.com/bmatcuk/doublestar v1.1.1/go.mod h1:UD6OnuiIn0yFxxA2le/rnRU1G4RaI4UvFv1sNto9p6w=
+github.com/bytedance/sonic v1.5.0/go.mod h1:ED5hyg4y6t3/9Ku1R6dU/4KyJ48DZ4jPhfY1O2AihPM=
+github.com/bytedance/sonic v1.10.0-rc/go.mod h1:ElCzW+ufi8qKqNW0FY314xriJhyJhuoJ3gFZdAHF7NM=
+github.com/bytedance/sonic v1.10.0-rc3 h1:uNSnscRapXTwUgTyOF0GVljYD08p9X/Lbr9MweSV3V0=
+github.com/bytedance/sonic v1.10.0-rc3/go.mod h1:iZcSUejdk5aukTND/Eu/ivjQuEL0Cu9/rf50Hi0u/g4=
 github.com/census-instrumentation/opencensus-proto v0.2.1/go.mod h1:f6KPmirojxKA12rnyqOA5BBL4O983OfeGPqjHWSTneU=
+github.com/chenzhuoyu/base64x v0.0.0-20211019084208-fb5309c8db06/go.mod h1:DH46F32mSOjUmXrMHnKwZdA8wcEefY7UVqBKYGjpdQY=
+github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311/go.mod h1:b583jCggY9gE99b6G5LEC39OIiVsWj+R97kbl5odCEk=
+github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d h1:77cEq6EriyTZ0g/qfRdp61a3Uu/AWrgIq2s0ClJV1g0=
+github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d/go.mod h1:8EPpVsBuRksnlj1mLy4AWzRNQYxauNi62uWcE3to6eA=
+github.com/chenzhuoyu/iasm v0.9.0 h1:9fhXjVzq5hUy2gkhhgHl95zG2cEAhw9OSGs8toWWAwo=
+github.com/chenzhuoyu/iasm v0.9.0/go.mod h1:Xjy2NpN3h7aUqeqM+woSuuvxmIe6+DDsiNLIrkAmYog=
 github.com/client9/misspell v0.3.4/go.mod h1:qj6jICC3Q7zFZvVWo7KLAzC3yx5G7kyvSDkc90ppPyw=
 github.com/cncf/udpa/go v0.0.0-20191209042840-269d4d468f6f/go.mod h1:M8M6+tZqaGXZJjfX53e64911xZQV5JYwmTeXPW+k8Sc=
 github.com/davecgh/go-spew v1.1.0/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
@@ -20,9 +30,26 @@ github.com/frankban/quicktest v1.14.6 h1:7Xjx+VpznH+oBnejlPUj8oUpdxnVs4f8XU8WnHk
 github.com/frankban/quicktest v1.14.6/go.mod h1:4ptaffx2x8+WTWXmUCuVU6aPUX1/Mz7zb5vbUoiM6w0=
 github.com/fsnotify/fsnotify v1.7.0 h1:8JEhPFa5W2WU7YfeZzPNqzMP6Lwt7L2715Ggo0nosvA=
 github.com/fsnotify/fsnotify v1.7.0/go.mod h1:40Bi/Hjc2AVfZrqy+aj+yEI+/bRxZnMJyTJwOpGvigM=
+github.com/gabriel-vasile/mimetype v1.4.2 h1:w5qFW6JKBz9Y393Y4q372O9A7cUSequkh1Q7OhCmWKU=
+github.com/gabriel-vasile/mimetype v1.4.2/go.mod h1:zApsH/mKG4w07erKIaJPFiX0Tsq9BFQgN3qGY5GnNgA=
+github.com/gin-contrib/sse v0.1.0 h1:Y/yl/+YNO8GZSjAhjMsSuLt29uWRFHdHYUb5lYOV9qE=
+github.com/gin-contrib/sse v0.1.0/go.mod h1:RHrZQHXnP2xjPF+u1gW/2HnVO7nvIa9PG3Gm+fLHvGI=
+github.com/gin-gonic/gin v1.9.1 h1:4idEAncQnU5cB7BeOkPtxjfCSye0AAm1R0RVIqJ+Jmg=
+github.com/gin-gonic/gin v1.9.1/go.mod h1:hPrL7YrpYKXt5YId3A/Tnip5kqbEAP+KLuI3SUcPTeU=
 github.com/go-kit/log v0.1.0/go.mod h1:zbhenjAZHb184qTLMA9ZjW7ThYL0H2mk7Q6pNt4vbaY=
 github.com/go-logfmt/logfmt v0.5.0/go.mod h1:wCYkCAKZfumFQihp8CzCvQ3paCTfi41vtzG1KdI/P7A=
+github.com/go-playground/assert/v2 v2.2.0 h1:JvknZsQTYeFEAhQwI4qEt9cyV5ONwRHC+lYKSsYSR8s=
+github.com/go-playground/assert/v2 v2.2.0/go.mod h1:VDjEfimB/XKnb+ZQfWdccd7VUvScMdVu0Titje2rxJ4=
+github.com/go-playground/locales v0.14.1 h1:EWaQ/wswjilfKLTECiXz7Rh+3BjFhfDFKv/oXslEjJA=
+github.com/go-playground/locales v0.14.1/go.mod h1:hxrqLVvrK65+Rwrd5Fc6F2O76J/NuW9t0sjnWqG1slY=
+github.com/go-playground/universal-translator v0.18.1 h1:Bcnm0ZwsGyWbCzImXv+pAJnYK9S473LQFuzCbDbfSFY=
+github.com/go-playground/universal-translator v0.18.1/go.mod h1:xekY+UJKNuX9WP91TpwSH2VMlDf28Uj24BCp08ZFTUY=
+github.com/go-playground/validator/v10 v10.14.1 h1:9c50NUPC30zyuKprjL3vNZ0m5oG+jU0zvx4AqHGnv4k=
+github.com/go-playground/validator/v10 v10.14.1/go.mod h1:9iXMNT7sEkjXb0I+enO7QXmzG6QCsPWY4zveKFVRSyU=
 github.com/go-stack/stack v1.8.0/go.mod h1:v0f6uXyyMGvRgIKkXu+yp6POWl0qKG85gN/melR3HDY=
+github.com/goccy/go-json v0.10.2 h1:CrxCmQqYDkv1z7lO7Wbh2HN93uovUHgrECaO5ZrCXAU=
+github.com/goccy/go-json v0.10.2/go.mod h1:6MelG93GURQebXPDq3khkgXZkazVtN9CRI+MGFi0w8I=
+github.com/gogo/protobuf v1.3.2 h1:Ov1cvc58UF3b5XjBnZv7+opcTcQFZebYjWzi34vdm4Q=
 github.com/gogo/protobuf v1.3.2/go.mod h1:P1XiOD3dCwIKUDQYPy72D8LYyHL2YPYrpS2s69NZV8Q=
 github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b/go.mod h1:SBH7ygxi8pfUlaOkMMuAQtPIUF8ecWP5IEl/CR7VP2Q=
 github.com/golang/mock v1.1.1/go.mod h1:oTYuIxOrZwtPieC+H1uAHpcLFnEyAGVDL/k47Jfbm0A=
@@ -36,15 +63,22 @@ github.com/google/go-cmp v0.2.0/go.mod h1:oXzfMopK8JAjlY9xF4vHSVASa0yLyX7SntLO5a
 github.com/google/go-cmp v0.5.5/go.mod h1:v8dTdLbMG2kIc/vJvl+f65V22dbkXbowE6jgT/gNBxE=
 github.com/google/go-cmp v0.6.0 h1:ofyhxvXcZhMsU5ulbFiLKl/XBFqE1GSq7atu8tAmTRI=
 github.com/google/go-cmp v0.6.0/go.mod h1:17dUlkBOakJ0+DkrSSNjCkIjxS6bF9zb3elmeNGIjoY=
+github.com/google/gofuzz v1.0.0/go.mod h1:dBl0BpW6vV/+mYPU4Po3pmUjxk6FQPldtuIdl/M65Eg=
 github.com/google/uuid v1.6.0 h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=
 github.com/google/uuid v1.6.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
 github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 h1:UH//fgunKIs4JdUbpDl1VZCDaL56wXCB/5+wF6uHfaI=
 github.com/grpc-ecosystem/go-grpc-middleware v1.4.0/go.mod h1:g5qyo/la0ALbONm6Vbp88Yd8NsDy6rZz+RcrMPxvld8=
 github.com/hashicorp/hcl v1.0.0 h1:0Anlzjpi4vEasTeNFn2mLJgTSwt0+6sfsiTG8qcWGx4=
 github.com/hashicorp/hcl v1.0.0/go.mod h1:E5yfLk+7swimpb2L/Alb/PJmXilQ/rhwaUYs4T20WEQ=
+github.com/json-iterator/go v1.1.12 h1:PV8peI4a0ysnczrg+LtxykD8LfKY9ML6u2jnxaEnrnM=
+github.com/json-iterator/go v1.1.12/go.mod h1:e30LSqwooZae/UwlEbR2852Gd8hjQvJoHmT4TnhNGBo=
 github.com/juju/gnuflag v0.0.0-20171113085948-2ce1bb71843d/go.mod h1:2PavIy+JPciBPrBUjwbNvtwB6RQlve+hkpll6QSNmOE=
 github.com/kisielk/errcheck v1.5.0/go.mod h1:pFxgyoBC7bSaBwPgfKdkLd5X25qrDl4LWUI2bnpBCr8=
 github.com/kisielk/gotool v1.0.0/go.mod h1:XhKaO+MFFWcvkIS/tQcRk01m1F5IRFswLeQ+oQHNcck=
+github.com/klauspost/cpuid/v2 v2.0.9/go.mod h1:FInQzS24/EEf25PyTYn52gqo7WaD8xa0213Md/qVLRg=
+github.com/klauspost/cpuid/v2 v2.2.5 h1:0E5MSMDEoAulmXNFquVs//DdoomxaoTY1kUhbc/qbZg=
+github.com/klauspost/cpuid/v2 v2.2.5/go.mod h1:Lcz8mBdAVJIBVzewtcLocK12l3Y+JytZYpaMropDUws=
+github.com/knz/go-libedit v1.10.1/go.mod h1:MZTVkCWyz0oBc7JOWP3wNAzd002ZbM/5hgShxwh4x8M=
 github.com/konsorten/go-windows-terminal-sequences v1.0.1/go.mod h1:T0+1ngSBFLxvqU3pZ+m/2kptfBszLMUkC4ZK/EgS/cQ=
 github.com/kr/pretty v0.1.0/go.mod h1:dAy3ld7l9f0ibDNOQOHHMYYIIbhfbHSm3C4ZsoJORNo=
 github.com/kr/pretty v0.3.1 h1:flRD4NNwYAUpkphVc1HcthR4KEIFJ65n8Mw5qdRn3LE=
@@ -53,10 +87,19 @@ github.com/kr/pty v1.1.1/go.mod h1:pFQYn66WHrOpPYNljwOMqo10TkYh1fy3cYio2l3bCsQ=
 github.com/kr/text v0.1.0/go.mod h1:4Jbv+DJW3UT/LiOwJeYQe1efqtUx/iVham/4vfdArNI=
 github.com/kr/text v0.2.0 h1:5Nx0Ya0ZqY2ygV366QzturHI13Jq95ApcVaJBhpS+AY=
 github.com/kr/text v0.2.0/go.mod h1:eLer722TekiGuMkidMxC/pM04lWEeraHUUmBw8l2grE=
+github.com/leodido/go-urn v1.2.4 h1:XlAE/cm/ms7TE/VMVoduSpNBoyc2dOxHs5MZSwAN63Q=
+github.com/leodido/go-urn v1.2.4/go.mod h1:7ZrI8mTSeBSHl/UaRyKQW1qZeMgak41ANeCNaVckg+4=
 github.com/magiconair/properties v1.8.7 h1:IeQXZAiQcpL9mgcAe1Nu6cX9LLw6ExEHKjN0VQdvPDY=
 github.com/magiconair/properties v1.8.7/go.mod h1:Dhd985XPs7jluiymwWYZ0G4Z61jb3vdS329zhj2hYo0=
+github.com/mattn/go-isatty v0.0.20 h1:xfD0iDuEKnDkl03q4limB+vH+GxLEtL/jb4xVJSWWEY=
+github.com/mattn/go-isatty v0.0.20/go.mod h1:W+V8PltTTMOvKvAeJH7IuucS94S2C6jfK/D7dTCTo3Y=
 github.com/mitchellh/mapstructure v1.5.0 h1:jeMsZIYE/09sWLaz43PL7Gy6RuMjD2eJVyuac5Z2hdY=
 github.com/mitchellh/mapstructure v1.5.0/go.mod h1:bFUtVrKA4DC2yAKiSyO/QUcy7e+RRV2QTWOzhPopBRo=
+github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421/go.mod h1:6dJC0mAP4ikYIbvyc7fijjWJddQyLn8Ig3JB5CqoB9Q=
+github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd h1:TRLaZ9cD/w8PVh93nsPXa1VrQ6jlwL5oN8l14QlcNfg=
+github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd/go.mod h1:6dJC0mAP4ikYIbvyc7fijjWJddQyLn8Ig3JB5CqoB9Q=
+github.com/modern-go/reflect2 v1.0.2 h1:xBagoLtFs94CBntxluKeaWgTMpvLxC4ur3nMaC9Gz0M=
+github.com/modern-go/reflect2 v1.0.2/go.mod h1:yWuevngMOJpCy52FWWMvUC8ws7m/LJsjYzDa0/r8luk=
 github.com/oapi-codegen/runtime v1.1.1 h1:EXLHh0DXIJnWhdRPN2w4MXAzFyE4CskzhNLUmtpMYro=
 github.com/oapi-codegen/runtime v1.1.1/go.mod h1:SK9X900oXmPWilYR5/WKPzt3Kqxn/uS/+lbpREv+eCg=
 github.com/opentracing/opentracing-go v1.1.0/go.mod h1:UkNAQd3GIcIGf0SeVgPpRdFStlNbqXla1AfSYxPUl2o=
@@ -74,6 +117,8 @@ github.com/sagikazarmark/locafero v0.4.0/go.mod h1:Pe1W6UlPYUk/+wc/6KFhbORCfqzgY
 github.com/sagikazarmark/slog-shim v0.1.0 h1:diDBnUNK9N/354PgrxMywXnAwEr1QZcOr6gto+ugjYE=
 github.com/sagikazarmark/slog-shim v0.1.0/go.mod h1:SrcSrq8aKtyuqEI1uvTDTK1arOWRIczQRv+GVI1AkeQ=
 github.com/sirupsen/logrus v1.4.2/go.mod h1:tLMulIdttU9McNUspp0xgXVQah82FyeX6MwdIuYE2rE=
+github.com/sirupsen/logrus v1.8.1 h1:dJKuHgqk1NNQlqoA6BTlM1Wf9DOH3NBjQyu0h9+AZZE=
+github.com/sirupsen/logrus v1.8.1/go.mod h1:yWOB1SBYBC5VeMP7gHvWumXLIWorT60ONWic61uBYv0=
 github.com/sourcegraph/conc v0.3.0 h1:OQTbbt6P72L20UqAkXXuLOj79LfEanQ+YQFNpLA9ySo=
 github.com/sourcegraph/conc v0.3.0/go.mod h1:Sdozi7LEKbFPqYX2/J+iBAM6HpqSLTASQIKqDmF7Mt0=
 github.com/spf13/afero v1.11.0 h1:WJQKhtpdm3v2IzqG8VMqrr6Rf3UYpEF239Jy9wNepM8=
@@ -96,11 +141,17 @@ github.com/stretchr/testify v1.4.0/go.mod h1:j7eGeouHqKxXV5pUuKE4zz7dFj8WfuZ+81P
 github.com/stretchr/testify v1.7.0/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
 github.com/stretchr/testify v1.7.1/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
 github.com/stretchr/testify v1.8.0/go.mod h1:yNjHg4UonilssWZ8iaSj1OCr/vHnekPRkoO+kdMU+MU=
+github.com/stretchr/testify v1.8.1/go.mod h1:w2LPCIKwWwSfY2zedu0+kehJoqGctiVI29o6fzry7u4=
+github.com/stretchr/testify v1.8.2/go.mod h1:w2LPCIKwWwSfY2zedu0+kehJoqGctiVI29o6fzry7u4=
 github.com/stretchr/testify v1.8.4/go.mod h1:sz/lmYIOXD/1dqDmKjjqLyZ2RngseejIcXlSw2iwfAo=
 github.com/stretchr/testify v1.9.0 h1:HtqpIVDClZ4nwg75+f6Lvsy/wHu+3BoSGCbBAcpTsTg=
 github.com/stretchr/testify v1.9.0/go.mod h1:r2ic/lqez/lEtzL7wO/rwa5dbSLXVDPFyf8C91i36aY=
 github.com/subosito/gotenv v1.6.0 h1:9NlTDc1FTs4qu0DDq7AEtTPNw6SVm7uBMsUCUjABIf8=
 github.com/subosito/gotenv v1.6.0/go.mod h1:Dk4QP5c2W3ibzajGcXpNraDfq2IrhjMIvMSWPKKo0FU=
+github.com/twitchyliquid64/golang-asm v0.15.1 h1:SU5vSMR7hnwNxj24w34ZyCi/FmDZTkS4MhqMhdFk5YI=
+github.com/twitchyliquid64/golang-asm v0.15.1/go.mod h1:a1lVb/DtPvCB8fslRZhAngC2+aY1QWCk3Cedj/Gdt08=
+github.com/ugorji/go/codec v1.2.11 h1:BMaWp1Bb6fHwEtbplGBGJ498wD+LKlNSl25MjdZY4dU=
+github.com/ugorji/go/codec v1.2.11/go.mod h1:UNopzCgEMSXjBc6AOMqYvWC1ktqTAfzJZUZgYf6w6lg=
 github.com/yuin/goldmark v1.1.27/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
 github.com/yuin/goldmark v1.2.1/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
 go.uber.org/atomic v1.7.0/go.mod h1:fEN4uk6kAWBTFdckzkM89CLk9XfWZrxpCo0nPH17wJc=
@@ -111,9 +162,14 @@ go.uber.org/multierr v1.6.0/go.mod h1:cdWPpRnG4AhwMwsgIHip0KRBQjJy5kYEpYjJxpXp9i
 go.uber.org/multierr v1.9.0 h1:7fIwc/ZtS0q++VgcfqFDxSBZVv/Xo49/SYnDFupUwlI=
 go.uber.org/multierr v1.9.0/go.mod h1:X2jQV1h+kxSjClGpnseKVIxpmcjrj7MNnI0bnlfKTVQ=
 go.uber.org/zap v1.18.1/go.mod h1:xg/QME4nWcxGxrpdeYfq7UvYrLh66cuVKdrbD1XF/NI=
+golang.org/x/arch v0.0.0-20210923205945-b76863e36670/go.mod h1:5om86z9Hs0C8fWVUuoMHwpExlXzs5Tkyp9hOrfG7pp8=
+golang.org/x/arch v0.4.0 h1:A8WCeEWhLwPBKNbFi5Wv5UTCBx5zzubnXDlMOFAzFMc=
+golang.org/x/arch v0.4.0/go.mod h1:5om86z9Hs0C8fWVUuoMHwpExlXzs5Tkyp9hOrfG7pp8=
 golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2/go.mod h1:djNgcEr1/C05ACkg1iLfiJU5Ep61QUkGW8qpdssI0+w=
 golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550/go.mod h1:yigFU9vqHzYiE8UmvKecakEJjdnWj3jj499lnFckfCI=
 golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9/go.mod h1:LzIPMQfyMNhhGPhUkYOs5KpL4U8rLKemX1yGLhDgUto=
+golang.org/x/crypto v0.21.0 h1:X31++rzVUdKhX5sWmSOFZxx8UW/ldWx55cbf08iNAMA=
+golang.org/x/crypto v0.21.0/go.mod h1:0BP7YvVV9gBbVKyeTG0Gyn+gZm94bibOW5BjDEYAOMs=
 golang.org/x/exp v0.0.0-20190121172915-509febef88a4/go.mod h1:CJ0aWSM057203Lf6IL+f9T1iT9GByDxfZKAQTCR3kQA=
 golang.org/x/exp v0.0.0-20230905200255-921286631fa9 h1:GoHiUyI/Tp2nVkLI2mCxVkOjsbSXD66ic0XW0js0R9g=
 golang.org/x/exp v0.0.0-20230905200255-921286631fa9/go.mod h1:S2oDrQGGwySpoQPVqRShND87VCbxmc6bL1Yd2oYrm6k=
@@ -143,8 +199,11 @@ golang.org/x/sys v0.0.0-20180830151530-49385e6e1522/go.mod h1:STP8DvDyc/dI5b8T5h
 golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
 golang.org/x/sys v0.0.0-20190412213103-97732733099d/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20190422165155-953cdadca894/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20191026070338-33540a1f6037/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20200930185726-fdedc70b468f/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20211025201205-69cdffdb9359/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
+golang.org/x/sys v0.5.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
+golang.org/x/sys v0.6.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.18.0 h1:DBdB3niSjOA/O0blCZBqDefyWNYveAYMNF1Wum0DYQ4=
 golang.org/x/sys v0.18.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
 golang.org/x/text v0.3.0/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
@@ -196,3 +255,5 @@ gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
 gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
 honnef.co/go/tools v0.0.0-20190102054323-c2f93a96b099/go.mod h1:rf3lG4BRIbNafJWhAfAdb/ePZxsR/4RtNHQocxwk9r4=
 honnef.co/go/tools v0.0.0-20190523083050-ea95bdfd59fc/go.mod h1:rf3lG4BRIbNafJWhAfAdb/ePZxsR/4RtNHQocxwk9r4=
+nullprogram.com/x/optparse v1.0.0/go.mod h1:KdyPE+Igbe0jQUrVfMqDMeJQIJZEuyV7pjYmp6pbG50=
+rsc.io/pdf v0.1.1/go.mod h1:n8OzWcQ6Sp37PL01nO98y4iUCRdTGarVfzxY20ICaU4=
diff --git a/internal/kitchen/.air.toml b/internal/kitchen/.air.toml
new file mode 100644
index 0000000..92667de
--- /dev/null
+++ b/internal/kitchen/.air.toml
@@ -0,0 +1,52 @@
+root = "."
+testdata_dir = "testdata"
+tmp_dir = "tmp"
+
+[build]
+  args_bin = []
+  bin = "./tmp/main"
+  cmd = "go build -o ./tmp/main ."
+  delay = 1000
+  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
+  exclude_file = []
+  exclude_regex = ["_test.go"]
+  exclude_unchanged = false
+  follow_symlink = false
+  full_bin = ""
+  include_dir = []
+  include_ext = ["go", "tpl", "tmpl", "html"]
+  include_file = []
+  kill_delay = "0s"
+  log = "build-errors.log"
+  poll = false
+  poll_interval = 0
+  post_cmd = []
+  pre_cmd = []
+  rerun = false
+  rerun_delay = 500
+  send_interrupt = true
+  stop_on_error = true
+
+[color]
+  app = ""
+  build = "yellow"
+  main = "magenta"
+  runner = "green"
+  watcher = "cyan"
+
+[log]
+  main_only = false
+  silent = false
+  time = false
+
+[misc]
+  clean_on_exit = false
+
+[proxy]
+  app_port = 0
+  enabled = false
+  proxy_port = 0
+
+[screen]
+  clear_on_rebuild = false
+  keep_scroll = true
diff --git a/internal/order/.air.toml b/internal/order/.air.toml
new file mode 100644
index 0000000..92667de
--- /dev/null
+++ b/internal/order/.air.toml
@@ -0,0 +1,52 @@
+root = "."
+testdata_dir = "testdata"
+tmp_dir = "tmp"
+
+[build]
+  args_bin = []
+  bin = "./tmp/main"
+  cmd = "go build -o ./tmp/main ."
+  delay = 1000
+  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
+  exclude_file = []
+  exclude_regex = ["_test.go"]
+  exclude_unchanged = false
+  follow_symlink = false
+  full_bin = ""
+  include_dir = []
+  include_ext = ["go", "tpl", "tmpl", "html"]
+  include_file = []
+  kill_delay = "0s"
+  log = "build-errors.log"
+  poll = false
+  poll_interval = 0
+  post_cmd = []
+  pre_cmd = []
+  rerun = false
+  rerun_delay = 500
+  send_interrupt = true
+  stop_on_error = true
+
+[color]
+  app = ""
+  build = "yellow"
+  main = "magenta"
+  runner = "green"
+  watcher = "cyan"
+
+[log]
+  main_only = false
+  silent = false
+  time = false
+
+[misc]
+  clean_on_exit = false
+
+[proxy]
+  app_port = 0
+  enabled = false
+  proxy_port = 0
+
+[screen]
+  clear_on_rebuild = false
+  keep_scroll = true
diff --git a/internal/order/app/app.go b/internal/order/app/app.go
new file mode 100644
index 0000000..42330cd
--- /dev/null
+++ b/internal/order/app/app.go
@@ -0,0 +1,10 @@
+package app
+
+type Application struct {
+	Commands Commands
+	Queries  Queries
+}
+
+type Commands struct{}
+
+type Queries struct{}
diff --git a/internal/order/http.go b/internal/order/http.go
index 4cfbaab..f18fa80 100644
--- a/internal/order/http.go
+++ b/internal/order/http.go
@@ -1,10 +1,13 @@
 package main
 
 import (
+	"github.com/ghost-yu/go_shop_second/order/app"
 	"github.com/gin-gonic/gin"
 )
 
-type HTTPServer struct{}
+type HTTPServer struct {
+	app app.Application
+}
 
 func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
 	//TODO implement me
diff --git a/internal/order/main.go b/internal/order/main.go
index 79bb274..6d2f2c5 100644
--- a/internal/order/main.go
+++ b/internal/order/main.go
@@ -1,33 +1,42 @@
 package main
 
 import (
-	"log"
+	"context"
 
 	"github.com/ghost-yu/go_shop_second/common/config"
 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
 	"github.com/ghost-yu/go_shop_second/common/server"
 	"github.com/ghost-yu/go_shop_second/order/ports"
+	"github.com/ghost-yu/go_shop_second/order/service"
 	"github.com/gin-gonic/gin"
+	"github.com/sirupsen/logrus"
 	"github.com/spf13/viper"
 	"google.golang.org/grpc"
 )
 
 func init() {
 	if err := config.NewViperConfig(); err != nil {
-		log.Fatal(err)
+		logrus.Fatal(err)
 	}
 }
 
 func main() {
 	serviceName := viper.GetString("order.service-name")
 
+	ctx, cancel := context.WithCancel(context.Background())
+	defer cancel()
+
+	application := service.NewApplication(ctx)
+
 	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
-		svc := ports.NewGRPCServer()
+		svc := ports.NewGRPCServer(application)
 		orderpb.RegisterOrderServiceServer(server, svc)
 	})
 
 	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
-		ports.RegisterHandlersWithOptions(router, HTTPServer{}, ports.GinServerOptions{
+		ports.RegisterHandlersWithOptions(router, HTTPServer{
+			app: application,
+		}, ports.GinServerOptions{
 			BaseURL:      "/api",
 			Middlewares:  nil,
 			ErrorHandler: nil,
diff --git a/internal/order/ports/grpc.go b/internal/order/ports/grpc.go
index b0b19db..e1e8621 100644
--- a/internal/order/ports/grpc.go
+++ b/internal/order/ports/grpc.go
@@ -4,14 +4,16 @@ import (
 	context "context"
 
 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+	"github.com/ghost-yu/go_shop_second/order/app"
 	"google.golang.org/protobuf/types/known/emptypb"
 )
 
 type GRPCServer struct {
+	app app.Application
 }
 
-func NewGRPCServer() *GRPCServer {
-	return &GRPCServer{}
+func NewGRPCServer(app app.Application) *GRPCServer {
+	return &GRPCServer{app: app}
 }
 
 func (G GRPCServer) CreateOrder(ctx context.Context, request *orderpb.CreateOrderRequest) (*emptypb.Empty, error) {
diff --git a/internal/order/service/application.go b/internal/order/service/application.go
new file mode 100644
index 0000000..122b22a
--- /dev/null
+++ b/internal/order/service/application.go
@@ -0,0 +1,11 @@
+package service
+
+import (
+	"context"
+
+	"github.com/ghost-yu/go_shop_second/order/app"
+)
+
+func NewApplication(ctx context.Context) app.Application {
+	return app.Application{}
+}
diff --git a/internal/payment/.air.toml b/internal/payment/.air.toml
new file mode 100644
index 0000000..92667de
--- /dev/null
+++ b/internal/payment/.air.toml
@@ -0,0 +1,52 @@
+root = "."
+testdata_dir = "testdata"
+tmp_dir = "tmp"
+
+[build]
+  args_bin = []
+  bin = "./tmp/main"
+  cmd = "go build -o ./tmp/main ."
+  delay = 1000
+  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
+  exclude_file = []
+  exclude_regex = ["_test.go"]
+  exclude_unchanged = false
+  follow_symlink = false
+  full_bin = ""
+  include_dir = []
+  include_ext = ["go", "tpl", "tmpl", "html"]
+  include_file = []
+  kill_delay = "0s"
+  log = "build-errors.log"
+  poll = false
+  poll_interval = 0
+  post_cmd = []
+  pre_cmd = []
+  rerun = false
+  rerun_delay = 500
+  send_interrupt = true
+  stop_on_error = true
+
+[color]
+  app = ""
+  build = "yellow"
+  main = "magenta"
+  runner = "green"
+  watcher = "cyan"
+
+[log]
+  main_only = false
+  silent = false
+  time = false
+
+[misc]
+  clean_on_exit = false
+
+[proxy]
+  app_port = 0
+  enabled = false
+  proxy_port = 0
+
+[screen]
+  clear_on_rebuild = false
+  keep_scroll = true
~~~
