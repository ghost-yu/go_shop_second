# Lesson Pair Diff Report

- FromBranch: lesson24
- ToBranch: lesson25

## Short Summary

~~~text
 18 files changed, 395 insertions(+), 183 deletions(-)
~~~

## File Stats

~~~text
 docker-compose.yml                 |  13 ++++-
 internal/common/client/grpc.go     |   2 +
 internal/common/config/global.yaml |   3 ++
 internal/common/go.mod             |  47 ++++++++++-------
 internal/common/go.sum             | 105 +++++++++++++++++++++----------------
 internal/common/server/gprc.go     |  12 +----
 internal/common/server/http.go     |   7 +++
 internal/common/tracing/jaeger.go  |  49 +++++++++++++++++
 internal/order/go.mod              |  10 ++++
 internal/order/go.sum              |  21 ++++++++
 internal/order/http.go             |  15 ++++--
 internal/order/main.go             |   7 +++
 internal/payment/go.mod            |  46 ++++++++++------
 internal/payment/go.sum            | 105 +++++++++++++++++++++----------------
 internal/payment/main.go           |  10 +++-
 internal/stock/go.mod              |  37 ++++++++-----
 internal/stock/go.sum              |  80 +++++++++++++++++-----------
 internal/stock/main.go             |   9 +++-
 18 files changed, 395 insertions(+), 183 deletions(-)
~~~

## Commit Comparison

~~~text
> f95add9 otel
~~~

## Changed Files

~~~text
docker-compose.yml
internal/common/client/grpc.go
internal/common/config/global.yaml
internal/common/go.mod
internal/common/go.sum
internal/common/server/gprc.go
internal/common/server/http.go
internal/common/tracing/jaeger.go
internal/order/go.mod
internal/order/go.sum
internal/order/http.go
internal/order/main.go
internal/payment/go.mod
internal/payment/go.sum
internal/payment/main.go
internal/stock/go.mod
internal/stock/go.sum
internal/stock/main.go
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
docker-compose.yml
internal/common/client/grpc.go
internal/common/config/global.yaml
internal/common/server/gprc.go
internal/common/server/http.go
internal/common/tracing/jaeger.go
internal/order/http.go
internal/order/main.go
internal/payment/main.go
internal/stock/main.go
~~~

## Full Diff

~~~diff
diff --git a/docker-compose.yml b/docker-compose.yml
index 06beeb4..12f1e17 100644
--- a/docker-compose.yml
+++ b/docker-compose.yml
@@ -11,4 +11,15 @@ services:
     image: "rabbitmq:3-management"
     ports:
       - "15672:15672"
-      - "5672:5672"
\ No newline at end of file
+      - "5672:5672"
+
+  jaeger:
+    image: "jaegertracing/all-in-one:latest"
+    ports:
+      - "6831:6831"
+      - "16686:16686"
+      - "14268:14268"
+      - "4318:4318"
+      - "4317:4317"
+    environment:
+      COLLECTOR_OTLP_ENABLED: true
\ No newline at end of file
diff --git a/internal/common/client/grpc.go b/internal/common/client/grpc.go
index f95506b..7adc093 100644
--- a/internal/common/client/grpc.go
+++ b/internal/common/client/grpc.go
@@ -11,6 +11,7 @@ import (
 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
 	"github.com/sirupsen/logrus"
 	"github.com/spf13/viper"
+	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
 	"google.golang.org/grpc"
 	"google.golang.org/grpc/credentials/insecure"
 )
@@ -57,6 +58,7 @@ func NewOrderGRPCClient(ctx context.Context) (client orderpb.OrderServiceClient,
 func grpcDialOpts(_ string) []grpc.DialOption {
 	return []grpc.DialOption{
 		grpc.WithTransportCredentials(insecure.NewCredentials()),
+		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
 	}
 }
 
diff --git a/internal/common/config/global.yaml b/internal/common/config/global.yaml
index 1bf2f6e..6b3d5eb 100644
--- a/internal/common/config/global.yaml
+++ b/internal/common/config/global.yaml
@@ -1,6 +1,9 @@
 fallback-grpc-addr: 127.0.0.1:3030
 dial-grpc-timeout: 10
 
+jaeger:
+  url: "http://127.0.0.1:14268/api/traces"
+
 consul:
   addr: 127.0.0.1:8500
 
diff --git a/internal/common/go.mod b/internal/common/go.mod
index af1aa77..c2021aa 100644
--- a/internal/common/go.mod
+++ b/internal/common/go.mod
@@ -3,13 +3,21 @@ module github.com/ghost-yu/go_shop_second/common
 go 1.22.8
 
 require (
-	github.com/gin-gonic/gin v1.9.1
+	github.com/gin-gonic/gin v1.10.0
 	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
 	github.com/hashicorp/consul/api v1.28.2
 	github.com/oapi-codegen/runtime v1.1.1
+	github.com/rabbitmq/amqp091-go v1.10.0
 	github.com/sirupsen/logrus v1.8.1
 	github.com/spf13/viper v1.19.0
 	github.com/x-cray/logrus-prefixed-formatter v0.5.2
+	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0
+	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0
+	go.opentelemetry.io/contrib/propagators/b3 v1.31.0
+	go.opentelemetry.io/otel v1.31.0
+	go.opentelemetry.io/otel/exporters/jaeger v1.17.0
+	go.opentelemetry.io/otel/sdk v1.31.0
+	go.opentelemetry.io/otel/trace v1.31.0
 	google.golang.org/grpc v1.67.1
 	google.golang.org/protobuf v1.35.1
 )
@@ -17,17 +25,20 @@ require (
 require (
 	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
 	github.com/armon/go-metrics v0.4.1 // indirect
-	github.com/bytedance/sonic v1.10.0-rc3 // indirect
-	github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d // indirect
-	github.com/chenzhuoyu/iasm v0.9.0 // indirect
+	github.com/bytedance/sonic v1.12.3 // indirect
+	github.com/bytedance/sonic/loader v0.2.0 // indirect
+	github.com/cloudwego/base64x v0.1.4 // indirect
+	github.com/cloudwego/iasm v0.2.0 // indirect
 	github.com/fatih/color v1.14.1 // indirect
 	github.com/fsnotify/fsnotify v1.7.0 // indirect
-	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
+	github.com/gabriel-vasile/mimetype v1.4.5 // indirect
 	github.com/gin-contrib/sse v0.1.0 // indirect
+	github.com/go-logr/logr v1.4.2 // indirect
+	github.com/go-logr/stdr v1.2.2 // indirect
 	github.com/go-playground/locales v0.14.1 // indirect
 	github.com/go-playground/universal-translator v0.18.1 // indirect
-	github.com/go-playground/validator/v10 v10.14.1 // indirect
-	github.com/goccy/go-json v0.10.2 // indirect
+	github.com/go-playground/validator/v10 v10.22.1 // indirect
+	github.com/goccy/go-json v0.10.3 // indirect
 	github.com/golang/protobuf v1.5.3 // indirect
 	github.com/google/uuid v1.6.0 // indirect
 	github.com/hashicorp/errwrap v1.1.0 // indirect
@@ -40,8 +51,8 @@ require (
 	github.com/hashicorp/hcl v1.0.0 // indirect
 	github.com/hashicorp/serf v0.10.1 // indirect
 	github.com/json-iterator/go v1.1.12 // indirect
-	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
-	github.com/leodido/go-urn v1.2.4 // indirect
+	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
+	github.com/leodido/go-urn v1.4.0 // indirect
 	github.com/magiconair/properties v1.8.7 // indirect
 	github.com/mattn/go-colorable v0.1.13 // indirect
 	github.com/mattn/go-isatty v0.0.20 // indirect
@@ -52,8 +63,7 @@ require (
 	github.com/modern-go/reflect2 v1.0.2 // indirect
 	github.com/onsi/ginkgo v1.16.5 // indirect
 	github.com/onsi/gomega v1.34.2 // indirect
-	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
-	github.com/rabbitmq/amqp091-go v1.10.0 // indirect
+	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
 	github.com/sagikazarmark/locafero v0.4.0 // indirect
 	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
 	github.com/sourcegraph/conc v0.3.0 // indirect
@@ -62,16 +72,17 @@ require (
 	github.com/spf13/pflag v1.0.5 // indirect
 	github.com/subosito/gotenv v1.6.0 // indirect
 	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
-	github.com/ugorji/go/codec v1.2.11 // indirect
+	github.com/ugorji/go/codec v1.2.12 // indirect
+	go.opentelemetry.io/otel/metric v1.31.0 // indirect
 	go.uber.org/atomic v1.9.0 // indirect
 	go.uber.org/multierr v1.9.0 // indirect
-	golang.org/x/arch v0.4.0 // indirect
-	golang.org/x/crypto v0.26.0 // indirect
+	golang.org/x/arch v0.11.0 // indirect
+	golang.org/x/crypto v0.28.0 // indirect
 	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
-	golang.org/x/net v0.28.0 // indirect
-	golang.org/x/sys v0.24.0 // indirect
-	golang.org/x/term v0.23.0 // indirect
-	golang.org/x/text v0.17.0 // indirect
+	golang.org/x/net v0.30.0 // indirect
+	golang.org/x/sys v0.26.0 // indirect
+	golang.org/x/term v0.25.0 // indirect
+	golang.org/x/text v0.19.0 // indirect
 	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
 	gopkg.in/ini.v1 v1.67.0 // indirect
 	gopkg.in/yaml.v3 v3.0.1 // indirect
diff --git a/internal/common/go.sum b/internal/common/go.sum
index 05ee4a2..f9de9f4 100644
--- a/internal/common/go.sum
+++ b/internal/common/go.sum
@@ -20,21 +20,20 @@ github.com/beorn7/perks v1.0.0/go.mod h1:KWe93zE9D1o94FZ5RNwFwVgaQK1VOXiVxmqh+Ce
 github.com/beorn7/perks v1.0.1/go.mod h1:G2ZrVWU2WbWT9wwq4/hrbKbnv/1ERSJQ0ibhJ6rlkpw=
 github.com/bgentry/speakeasy v0.1.0/go.mod h1:+zsyZBPWlz7T6j88CTgSN5bM796AkVf0kBD4zp0CCIs=
 github.com/bmatcuk/doublestar v1.1.1/go.mod h1:UD6OnuiIn0yFxxA2le/rnRU1G4RaI4UvFv1sNto9p6w=
-github.com/bytedance/sonic v1.5.0/go.mod h1:ED5hyg4y6t3/9Ku1R6dU/4KyJ48DZ4jPhfY1O2AihPM=
-github.com/bytedance/sonic v1.10.0-rc/go.mod h1:ElCzW+ufi8qKqNW0FY314xriJhyJhuoJ3gFZdAHF7NM=
-github.com/bytedance/sonic v1.10.0-rc3 h1:uNSnscRapXTwUgTyOF0GVljYD08p9X/Lbr9MweSV3V0=
-github.com/bytedance/sonic v1.10.0-rc3/go.mod h1:iZcSUejdk5aukTND/Eu/ivjQuEL0Cu9/rf50Hi0u/g4=
+github.com/bytedance/sonic v1.12.3 h1:W2MGa7RCU1QTeYRTPE3+88mVC0yXmsRQRChiyVocVjU=
+github.com/bytedance/sonic v1.12.3/go.mod h1:B8Gt/XvtZ3Fqj+iSKMypzymZxw/FVwgIGKzMzT9r/rk=
+github.com/bytedance/sonic/loader v0.1.1/go.mod h1:ncP89zfokxS5LZrJxl5z0UJcsk4M4yY2JpfqGeCtNLU=
+github.com/bytedance/sonic/loader v0.2.0 h1:zNprn+lsIP06C/IqCHs3gPQIvnvpKbbxyXQP1iU4kWM=
+github.com/bytedance/sonic/loader v0.2.0/go.mod h1:ncP89zfokxS5LZrJxl5z0UJcsk4M4yY2JpfqGeCtNLU=
 github.com/census-instrumentation/opencensus-proto v0.2.1/go.mod h1:f6KPmirojxKA12rnyqOA5BBL4O983OfeGPqjHWSTneU=
 github.com/cespare/xxhash/v2 v2.1.1/go.mod h1:VGX0DQ3Q6kWi7AoAeZDth3/j3BFtOZR5XLFGgcrjCOs=
-github.com/chenzhuoyu/base64x v0.0.0-20211019084208-fb5309c8db06/go.mod h1:DH46F32mSOjUmXrMHnKwZdA8wcEefY7UVqBKYGjpdQY=
-github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311/go.mod h1:b583jCggY9gE99b6G5LEC39OIiVsWj+R97kbl5odCEk=
-github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d h1:77cEq6EriyTZ0g/qfRdp61a3Uu/AWrgIq2s0ClJV1g0=
-github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d/go.mod h1:8EPpVsBuRksnlj1mLy4AWzRNQYxauNi62uWcE3to6eA=
-github.com/chenzhuoyu/iasm v0.9.0 h1:9fhXjVzq5hUy2gkhhgHl95zG2cEAhw9OSGs8toWWAwo=
-github.com/chenzhuoyu/iasm v0.9.0/go.mod h1:Xjy2NpN3h7aUqeqM+woSuuvxmIe6+DDsiNLIrkAmYog=
 github.com/circonus-labs/circonus-gometrics v2.3.1+incompatible/go.mod h1:nmEj6Dob7S7YxXgwXpfOuvO54S+tGdZdw9fuRZt25Ag=
 github.com/circonus-labs/circonusllhist v0.1.3/go.mod h1:kMXHVDlOchFAehlya5ePtbp5jckzBHf4XRpQvBOLI+I=
 github.com/client9/misspell v0.3.4/go.mod h1:qj6jICC3Q7zFZvVWo7KLAzC3yx5G7kyvSDkc90ppPyw=
+github.com/cloudwego/base64x v0.1.4 h1:jwCgWpFanWmN8xoIUHa2rtzmkd5J2plF/dnLS6Xd/0Y=
+github.com/cloudwego/base64x v0.1.4/go.mod h1:0zlkT4Wn5C6NdauXdJRhSKRlJvmclQ1hhJgA0rcu/8w=
+github.com/cloudwego/iasm v0.2.0 h1:1KNIy1I1H9hNNFEEH3DVnI4UujN+1zjpuk6gwHLTssg=
+github.com/cloudwego/iasm v0.2.0/go.mod h1:8rXZaNYT2n95jn+zTI1sDr+IgcD2GVs0nlbbQPiEFhY=
 github.com/cncf/udpa/go v0.0.0-20191209042840-269d4d468f6f/go.mod h1:M8M6+tZqaGXZJjfX53e64911xZQV5JYwmTeXPW+k8Sc=
 github.com/davecgh/go-spew v1.1.0/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
 github.com/davecgh/go-spew v1.1.1/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
@@ -55,30 +54,35 @@ github.com/fsnotify/fsnotify v1.4.7/go.mod h1:jwhsz4b93w/PPRr/qN1Yymfu8t87LnFCMo
 github.com/fsnotify/fsnotify v1.4.9/go.mod h1:znqG4EE+3YCdAaPaxE2ZRY/06pZUdp0tY4IgpuI1SZQ=
 github.com/fsnotify/fsnotify v1.7.0 h1:8JEhPFa5W2WU7YfeZzPNqzMP6Lwt7L2715Ggo0nosvA=
 github.com/fsnotify/fsnotify v1.7.0/go.mod h1:40Bi/Hjc2AVfZrqy+aj+yEI+/bRxZnMJyTJwOpGvigM=
-github.com/gabriel-vasile/mimetype v1.4.2 h1:w5qFW6JKBz9Y393Y4q372O9A7cUSequkh1Q7OhCmWKU=
-github.com/gabriel-vasile/mimetype v1.4.2/go.mod h1:zApsH/mKG4w07erKIaJPFiX0Tsq9BFQgN3qGY5GnNgA=
+github.com/gabriel-vasile/mimetype v1.4.5 h1:J7wGKdGu33ocBOhGy0z653k/lFKLFDPJMG8Gql0kxn4=
+github.com/gabriel-vasile/mimetype v1.4.5/go.mod h1:ibHel+/kbxn9x2407k1izTA1S81ku1z/DlgOW2QE0M4=
 github.com/gin-contrib/sse v0.1.0 h1:Y/yl/+YNO8GZSjAhjMsSuLt29uWRFHdHYUb5lYOV9qE=
 github.com/gin-contrib/sse v0.1.0/go.mod h1:RHrZQHXnP2xjPF+u1gW/2HnVO7nvIa9PG3Gm+fLHvGI=
-github.com/gin-gonic/gin v1.9.1 h1:4idEAncQnU5cB7BeOkPtxjfCSye0AAm1R0RVIqJ+Jmg=
-github.com/gin-gonic/gin v1.9.1/go.mod h1:hPrL7YrpYKXt5YId3A/Tnip5kqbEAP+KLuI3SUcPTeU=
+github.com/gin-gonic/gin v1.10.0 h1:nTuyha1TYqgedzytsKYqna+DfLos46nTv2ygFy86HFU=
+github.com/gin-gonic/gin v1.10.0/go.mod h1:4PMNQiOhvDRa013RKVbsiNwoyezlm2rm0uX/T7kzp5Y=
 github.com/go-kit/kit v0.8.0/go.mod h1:xBxKIO96dXMWWy0MnWVtmwkA9/13aqxPnvrjFYMA2as=
 github.com/go-kit/kit v0.9.0/go.mod h1:xBxKIO96dXMWWy0MnWVtmwkA9/13aqxPnvrjFYMA2as=
 github.com/go-kit/log v0.1.0/go.mod h1:zbhenjAZHb184qTLMA9ZjW7ThYL0H2mk7Q6pNt4vbaY=
 github.com/go-logfmt/logfmt v0.3.0/go.mod h1:Qt1PoO58o5twSAckw1HlFXLmHsOX5/0LbT9GBnD5lWE=
 github.com/go-logfmt/logfmt v0.4.0/go.mod h1:3RMwSq7FuexP4Kalkev3ejPJsZTpXXBr9+V4qmtdjCk=
 github.com/go-logfmt/logfmt v0.5.0/go.mod h1:wCYkCAKZfumFQihp8CzCvQ3paCTfi41vtzG1KdI/P7A=
+github.com/go-logr/logr v1.2.2/go.mod h1:jdQByPbusPIv2/zmleS9BjJVeZ6kBagPoEUsqbVz/1A=
+github.com/go-logr/logr v1.4.2 h1:6pFjapn8bFcIbiKo3XT4j/BhANplGihG6tvd+8rYgrY=
+github.com/go-logr/logr v1.4.2/go.mod h1:9T104GzyrTigFIr8wt5mBrctHMim0Nb2HLGrmQ40KvY=
+github.com/go-logr/stdr v1.2.2 h1:hSWxHoqTgW2S2qGc0LTAI563KZ5YKYRhT3MFKZMbjag=
+github.com/go-logr/stdr v1.2.2/go.mod h1:mMo/vtBO5dYbehREoey6XUKy/eSumjCCveDpRre4VKE=
 github.com/go-playground/assert/v2 v2.2.0 h1:JvknZsQTYeFEAhQwI4qEt9cyV5ONwRHC+lYKSsYSR8s=
 github.com/go-playground/assert/v2 v2.2.0/go.mod h1:VDjEfimB/XKnb+ZQfWdccd7VUvScMdVu0Titje2rxJ4=
 github.com/go-playground/locales v0.14.1 h1:EWaQ/wswjilfKLTECiXz7Rh+3BjFhfDFKv/oXslEjJA=
 github.com/go-playground/locales v0.14.1/go.mod h1:hxrqLVvrK65+Rwrd5Fc6F2O76J/NuW9t0sjnWqG1slY=
 github.com/go-playground/universal-translator v0.18.1 h1:Bcnm0ZwsGyWbCzImXv+pAJnYK9S473LQFuzCbDbfSFY=
 github.com/go-playground/universal-translator v0.18.1/go.mod h1:xekY+UJKNuX9WP91TpwSH2VMlDf28Uj24BCp08ZFTUY=
-github.com/go-playground/validator/v10 v10.14.1 h1:9c50NUPC30zyuKprjL3vNZ0m5oG+jU0zvx4AqHGnv4k=
-github.com/go-playground/validator/v10 v10.14.1/go.mod h1:9iXMNT7sEkjXb0I+enO7QXmzG6QCsPWY4zveKFVRSyU=
+github.com/go-playground/validator/v10 v10.22.1 h1:40JcKH+bBNGFczGuoBYgX4I6m/i27HYW8P9FDk5PbgA=
+github.com/go-playground/validator/v10 v10.22.1/go.mod h1:dbuPbCMFw/DrkbEynArYaCwl3amGuJotoKCe95atGMM=
 github.com/go-stack/stack v1.8.0/go.mod h1:v0f6uXyyMGvRgIKkXu+yp6POWl0qKG85gN/melR3HDY=
 github.com/go-task/slim-sprig v0.0.0-20210107165309-348f09dbbbc0/go.mod h1:fyg7847qk6SyHyPtNmDHnmrv/HOrqktSC+C9fM+CJOE=
-github.com/goccy/go-json v0.10.2 h1:CrxCmQqYDkv1z7lO7Wbh2HN93uovUHgrECaO5ZrCXAU=
-github.com/goccy/go-json v0.10.2/go.mod h1:6MelG93GURQebXPDq3khkgXZkazVtN9CRI+MGFi0w8I=
+github.com/goccy/go-json v0.10.3 h1:KZ5WoDbxAIgm2HNbYckL0se1fHD6rz5j4ywS6ebzDqA=
+github.com/goccy/go-json v0.10.3/go.mod h1:oq7eo15ShAhp70Anwd5lgX2pLfOS3QCiwU/PULtXL6M=
 github.com/gogo/protobuf v1.1.1/go.mod h1:r8qH/GZQm5c6nD/R0oafs1akxWv10x8SbQlK7atdtwQ=
 github.com/gogo/protobuf v1.3.2 h1:Ov1cvc58UF3b5XjBnZv7+opcTcQFZebYjWzi34vdm4Q=
 github.com/gogo/protobuf v1.3.2/go.mod h1:P1XiOD3dCwIKUDQYPy72D8LYyHL2YPYrpS2s69NZV8Q=
@@ -168,8 +172,8 @@ github.com/julienschmidt/httprouter v1.2.0/go.mod h1:SYymIcj16QtmaHHD7aYtjjsJG7V
 github.com/kisielk/errcheck v1.5.0/go.mod h1:pFxgyoBC7bSaBwPgfKdkLd5X25qrDl4LWUI2bnpBCr8=
 github.com/kisielk/gotool v1.0.0/go.mod h1:XhKaO+MFFWcvkIS/tQcRk01m1F5IRFswLeQ+oQHNcck=
 github.com/klauspost/cpuid/v2 v2.0.9/go.mod h1:FInQzS24/EEf25PyTYn52gqo7WaD8xa0213Md/qVLRg=
-github.com/klauspost/cpuid/v2 v2.2.5 h1:0E5MSMDEoAulmXNFquVs//DdoomxaoTY1kUhbc/qbZg=
-github.com/klauspost/cpuid/v2 v2.2.5/go.mod h1:Lcz8mBdAVJIBVzewtcLocK12l3Y+JytZYpaMropDUws=
+github.com/klauspost/cpuid/v2 v2.2.8 h1:+StwCXwm9PdpiEkPyzBXIy+M9KUb4ODm0Zarf1kS5BM=
+github.com/klauspost/cpuid/v2 v2.2.8/go.mod h1:Lcz8mBdAVJIBVzewtcLocK12l3Y+JytZYpaMropDUws=
 github.com/knz/go-libedit v1.10.1/go.mod h1:MZTVkCWyz0oBc7JOWP3wNAzd002ZbM/5hgShxwh4x8M=
 github.com/konsorten/go-windows-terminal-sequences v1.0.1/go.mod h1:T0+1ngSBFLxvqU3pZ+m/2kptfBszLMUkC4ZK/EgS/cQ=
 github.com/kr/logfmt v0.0.0-20140226030751-b84e30acd515/go.mod h1:+0opPa2QZZtGFBFZlji/RkVcI2GknAs/DXo4wKdlNEc=
@@ -180,8 +184,8 @@ github.com/kr/pty v1.1.1/go.mod h1:pFQYn66WHrOpPYNljwOMqo10TkYh1fy3cYio2l3bCsQ=
 github.com/kr/text v0.1.0/go.mod h1:4Jbv+DJW3UT/LiOwJeYQe1efqtUx/iVham/4vfdArNI=
 github.com/kr/text v0.2.0 h1:5Nx0Ya0ZqY2ygV366QzturHI13Jq95ApcVaJBhpS+AY=
 github.com/kr/text v0.2.0/go.mod h1:eLer722TekiGuMkidMxC/pM04lWEeraHUUmBw8l2grE=
-github.com/leodido/go-urn v1.2.4 h1:XlAE/cm/ms7TE/VMVoduSpNBoyc2dOxHs5MZSwAN63Q=
-github.com/leodido/go-urn v1.2.4/go.mod h1:7ZrI8mTSeBSHl/UaRyKQW1qZeMgak41ANeCNaVckg+4=
+github.com/leodido/go-urn v1.4.0 h1:WT9HwE9SGECu3lg4d/dIA+jxlljEa1/ffXKmRjqdmIQ=
+github.com/leodido/go-urn v1.4.0/go.mod h1:bvxc+MVxLKB4z00jd1z+Dvzr47oO32F/QSNjSBOlFxI=
 github.com/magiconair/properties v1.8.7 h1:IeQXZAiQcpL9mgcAe1Nu6cX9LLw6ExEHKjN0VQdvPDY=
 github.com/magiconair/properties v1.8.7/go.mod h1:Dhd985XPs7jluiymwWYZ0G4Z61jb3vdS329zhj2hYo0=
 github.com/mattn/go-colorable v0.0.9/go.mod h1:9vuHe8Xs5qXnSaW/c/ABM9alt+Vo+STaOChaDxuIBZU=
@@ -236,8 +240,8 @@ github.com/opentracing/opentracing-go v1.1.0/go.mod h1:UkNAQd3GIcIGf0SeVgPpRdFSt
 github.com/pascaldekloe/goe v0.0.0-20180627143212-57f6aae5913c/go.mod h1:lzWF7FIEvWOWxwDKqyGYQf6ZUaNfKdP144TG7ZOy1lc=
 github.com/pascaldekloe/goe v0.1.0 h1:cBOtyMzM9HTpWjXfbbunk26uA6nG3a8n06Wieeh0MwY=
 github.com/pascaldekloe/goe v0.1.0/go.mod h1:lzWF7FIEvWOWxwDKqyGYQf6ZUaNfKdP144TG7ZOy1lc=
-github.com/pelletier/go-toml/v2 v2.2.2 h1:aYUidT7k73Pcl9nb2gScu7NSrKCSHIDE89b3+6Wq+LM=
-github.com/pelletier/go-toml/v2 v2.2.2/go.mod h1:1t835xjRzz80PqgE6HHgN2JOsmgYu/h4qDAS4n929Rs=
+github.com/pelletier/go-toml/v2 v2.2.3 h1:YmeHyLY8mFWbdkNWwpr+qIL2bEqT0o95WSdkNHvL12M=
+github.com/pelletier/go-toml/v2 v2.2.3/go.mod h1:MfCQTFTvCcUyyvvwm1+G6H/jORL20Xlb6rzQu9GuUkc=
 github.com/pkg/errors v0.8.0/go.mod h1:bwawxfHBFNV+L2hUp1rHADufV3IMtnDRdf1r5NINEl0=
 github.com/pkg/errors v0.8.1/go.mod h1:bwawxfHBFNV+L2hUp1rHADufV3IMtnDRdf1r5NINEl0=
 github.com/pkg/errors v0.9.1 h1:FEBLx1zS214owpjy7qsBeixbURkuhQAwrK5UwLGTwt4=
@@ -288,9 +292,8 @@ github.com/spkg/bom v0.0.0-20160624110644-59b7046e48ad/go.mod h1:qLr4V1qq6nMqFKk
 github.com/stretchr/objx v0.1.0/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
 github.com/stretchr/objx v0.1.1/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
 github.com/stretchr/objx v0.4.0/go.mod h1:YvHI0jy2hoMjB+UWwv71VJQ9isScKT/TqJzVSSt89Yw=
+github.com/stretchr/objx v0.5.0 h1:1zr/of2m5FGMsad5YfcqgdqdWrIhu+EBEJRhR1U7z/c=
 github.com/stretchr/objx v0.5.0/go.mod h1:Yh+to48EsGEfYuaHDzXPcE3xhTkx73EhmCGUpEOglKo=
-github.com/stretchr/objx v0.5.2 h1:xuMeJ0Sdp5ZMRXx/aWO6RZxdr3beISkG5/G/aIRr3pY=
-github.com/stretchr/objx v0.5.2/go.mod h1:FRsXN1f5AsAjCGJKqEizvkpNtU+EGNCLh3NxZ/8L+MA=
 github.com/stretchr/testify v1.2.2/go.mod h1:a8OnRcib4nhh0OaRAV+Yts87kKdq0PP7pXfy6kDkUVs=
 github.com/stretchr/testify v1.3.0/go.mod h1:M5WIy9Dh21IEIfnGCwXGc5bZfKNJtfHm1UVUgZn+9EI=
 github.com/stretchr/testify v1.4.0/go.mod h1:j7eGeouHqKxXV5pUuKE4zz7dFj8WfuZ+81PSLYec5m4=
@@ -300,8 +303,6 @@ github.com/stretchr/testify v1.7.1/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/
 github.com/stretchr/testify v1.7.2/go.mod h1:R6va5+xMeoiuVRoj+gSkQ7d3FALtqAAGI1FQKckRals=
 github.com/stretchr/testify v1.8.0/go.mod h1:yNjHg4UonilssWZ8iaSj1OCr/vHnekPRkoO+kdMU+MU=
 github.com/stretchr/testify v1.8.1/go.mod h1:w2LPCIKwWwSfY2zedu0+kehJoqGctiVI29o6fzry7u4=
-github.com/stretchr/testify v1.8.2/go.mod h1:w2LPCIKwWwSfY2zedu0+kehJoqGctiVI29o6fzry7u4=
-github.com/stretchr/testify v1.8.4/go.mod h1:sz/lmYIOXD/1dqDmKjjqLyZ2RngseejIcXlSw2iwfAo=
 github.com/stretchr/testify v1.9.0 h1:HtqpIVDClZ4nwg75+f6Lvsy/wHu+3BoSGCbBAcpTsTg=
 github.com/stretchr/testify v1.9.0/go.mod h1:r2ic/lqez/lEtzL7wO/rwa5dbSLXVDPFyf8C91i36aY=
 github.com/subosito/gotenv v1.6.0 h1:9NlTDc1FTs4qu0DDq7AEtTPNw6SVm7uBMsUCUjABIf8=
@@ -309,30 +310,47 @@ github.com/subosito/gotenv v1.6.0/go.mod h1:Dk4QP5c2W3ibzajGcXpNraDfq2IrhjMIvMSW
 github.com/tv42/httpunix v0.0.0-20150427012821-b75d8614f926/go.mod h1:9ESjWnEqriFuLhtthL60Sar/7RFoluCcXsuvEwTV5KM=
 github.com/twitchyliquid64/golang-asm v0.15.1 h1:SU5vSMR7hnwNxj24w34ZyCi/FmDZTkS4MhqMhdFk5YI=
 github.com/twitchyliquid64/golang-asm v0.15.1/go.mod h1:a1lVb/DtPvCB8fslRZhAngC2+aY1QWCk3Cedj/Gdt08=
-github.com/ugorji/go/codec v1.2.11 h1:BMaWp1Bb6fHwEtbplGBGJ498wD+LKlNSl25MjdZY4dU=
-github.com/ugorji/go/codec v1.2.11/go.mod h1:UNopzCgEMSXjBc6AOMqYvWC1ktqTAfzJZUZgYf6w6lg=
+github.com/ugorji/go/codec v1.2.12 h1:9LC83zGrHhuUA9l16C9AHXAqEV/2wBQ4nkvumAE65EE=
+github.com/ugorji/go/codec v1.2.12/go.mod h1:UNopzCgEMSXjBc6AOMqYvWC1ktqTAfzJZUZgYf6w6lg=
 github.com/x-cray/logrus-prefixed-formatter v0.5.2 h1:00txxvfBM9muc0jiLIEAkAcIMJzfthRT6usrui8uGmg=
 github.com/x-cray/logrus-prefixed-formatter v0.5.2/go.mod h1:2duySbKsL6M18s5GU7VPsoEPHyzalCE06qoARUCeBBE=
 github.com/yuin/goldmark v1.1.27/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
 github.com/yuin/goldmark v1.2.1/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
+go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0 h1:0nTRpaCaILLdooXAQnfktlL6Zw1ECKEW9DZGH2byi2c=
+go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0/go.mod h1:A7aFlp4WSLmeOnFRZwf2dMU+40THPc+rsr6KOwZLOcg=
+go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0 h1:4Pp6oUg3+e/6M4C0A/3kJ2VYa++dsWVTtGgLVj5xtHg=
+go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0/go.mod h1:Mjt1i1INqiaoZOMGR1RIUJN+i3ChKoFRqzrRQhlkbs0=
+go.opentelemetry.io/contrib/propagators/b3 v1.31.0 h1:PQPXYscmwbCp76QDvO4hMngF2j8Bx/OTV86laEl8uqo=
+go.opentelemetry.io/contrib/propagators/b3 v1.31.0/go.mod h1:jbqfV8wDdqSDrAYxVpXQnpM0XFMq2FtDesblJ7blOwQ=
+go.opentelemetry.io/otel v1.31.0 h1:NsJcKPIW0D0H3NgzPDHmo0WW6SptzPdqg/L1zsIm2hY=
+go.opentelemetry.io/otel v1.31.0/go.mod h1:O0C14Yl9FgkjqcCZAsE053C13OaddMYr/hz6clDkEJE=
+go.opentelemetry.io/otel/exporters/jaeger v1.17.0 h1:D7UpUy2Xc2wsi1Ras6V40q806WM07rqoCWzXu7Sqy+4=
+go.opentelemetry.io/otel/exporters/jaeger v1.17.0/go.mod h1:nPCqOnEH9rNLKqH/+rrUjiMzHJdV1BlpKcTwRTyKkKI=
+go.opentelemetry.io/otel/metric v1.31.0 h1:FSErL0ATQAmYHUIzSezZibnyVlft1ybhy4ozRPcF2fE=
+go.opentelemetry.io/otel/metric v1.31.0/go.mod h1:C3dEloVbLuYoX41KpmAhOqNriGbA+qqH6PQ5E5mUfnY=
+go.opentelemetry.io/otel/sdk v1.31.0 h1:xLY3abVHYZ5HSfOg3l2E5LUj2Cwva5Y7yGxnSW9H5Gk=
+go.opentelemetry.io/otel/sdk v1.31.0/go.mod h1:TfRbMdhvxIIr/B2N2LQW2S5v9m3gOQ/08KsbbO5BPT0=
+go.opentelemetry.io/otel/trace v1.31.0 h1:ffjsj1aRouKewfr85U2aGagJ46+MvodynlQ1HYdmJys=
+go.opentelemetry.io/otel/trace v1.31.0/go.mod h1:TXZkRk7SM2ZQLtR6eoAWQFIHPvzQ06FJAsO1tJg480A=
 go.uber.org/atomic v1.7.0/go.mod h1:fEN4uk6kAWBTFdckzkM89CLk9XfWZrxpCo0nPH17wJc=
 go.uber.org/atomic v1.9.0 h1:ECmE8Bn/WFTYwEW/bpKD3M8VtR/zQVbavAoalC1PYyE=
 go.uber.org/atomic v1.9.0/go.mod h1:fEN4uk6kAWBTFdckzkM89CLk9XfWZrxpCo0nPH17wJc=
 go.uber.org/goleak v1.1.10/go.mod h1:8a7PlsEVH3e/a/GLqe5IIrQx6GzcnRmZEufDUTk4A7A=
+go.uber.org/goleak v1.3.0 h1:2K3zAYmnTNqV73imy9J1T3WC+gmCePx2hEGkimedGto=
+go.uber.org/goleak v1.3.0/go.mod h1:CoHD4mav9JJNrW/WLlf7HGZPjdw8EucARQHekz1X6bE=
 go.uber.org/multierr v1.6.0/go.mod h1:cdWPpRnG4AhwMwsgIHip0KRBQjJy5kYEpYjJxpXp9iU=
 go.uber.org/multierr v1.9.0 h1:7fIwc/ZtS0q++VgcfqFDxSBZVv/Xo49/SYnDFupUwlI=
 go.uber.org/multierr v1.9.0/go.mod h1:X2jQV1h+kxSjClGpnseKVIxpmcjrj7MNnI0bnlfKTVQ=
 go.uber.org/zap v1.18.1/go.mod h1:xg/QME4nWcxGxrpdeYfq7UvYrLh66cuVKdrbD1XF/NI=
-golang.org/x/arch v0.0.0-20210923205945-b76863e36670/go.mod h1:5om86z9Hs0C8fWVUuoMHwpExlXzs5Tkyp9hOrfG7pp8=
-golang.org/x/arch v0.4.0 h1:A8WCeEWhLwPBKNbFi5Wv5UTCBx5zzubnXDlMOFAzFMc=
-golang.org/x/arch v0.4.0/go.mod h1:5om86z9Hs0C8fWVUuoMHwpExlXzs5Tkyp9hOrfG7pp8=
+golang.org/x/arch v0.11.0 h1:KXV8WWKCXm6tRpLirl2szsO5j/oOODwZf4hATmGVNs4=
+golang.org/x/arch v0.11.0/go.mod h1:FEVrYAQjsQXMVJ1nsMoVVXPZg6p2JE2mx8psSWTDQys=
 golang.org/x/crypto v0.0.0-20180904163835-0709b304e793/go.mod h1:6SG95UA2DQfeDnfUPMdvaQW0Q7yPrPDi9nlGo2tz2b4=
 golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2/go.mod h1:djNgcEr1/C05ACkg1iLfiJU5Ep61QUkGW8qpdssI0+w=
 golang.org/x/crypto v0.0.0-20190923035154-9ee001bba392/go.mod h1:/lpIB1dKB+9EgE3H3cr1v9wB50oz8l4C4h62xy7jSTY=
 golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550/go.mod h1:yigFU9vqHzYiE8UmvKecakEJjdnWj3jj499lnFckfCI=
 golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9/go.mod h1:LzIPMQfyMNhhGPhUkYOs5KpL4U8rLKemX1yGLhDgUto=
-golang.org/x/crypto v0.26.0 h1:RrRspgV4mU+YwB4FYnuBoKsUapNIL5cohGAmSH3azsw=
-golang.org/x/crypto v0.26.0/go.mod h1:GY7jblb9wI+FOo5y8/S2oY4zWP07AkOJ4+jxCqdqn54=
+golang.org/x/crypto v0.28.0 h1:GBDwsMXVQi34v5CCYUm2jkJvu4cbtru2U4TN2PSyQnw=
+golang.org/x/crypto v0.28.0/go.mod h1:rmgy+3RHxRZMyY0jjAJShp2zgEdOqj2AO7U0pYmeQ7U=
 golang.org/x/exp v0.0.0-20190121172915-509febef88a4/go.mod h1:CJ0aWSM057203Lf6IL+f9T1iT9GByDxfZKAQTCR3kQA=
 golang.org/x/exp v0.0.0-20230905200255-921286631fa9 h1:GoHiUyI/Tp2nVkLI2mCxVkOjsbSXD66ic0XW0js0R9g=
 golang.org/x/exp v0.0.0-20230905200255-921286631fa9/go.mod h1:S2oDrQGGwySpoQPVqRShND87VCbxmc6bL1Yd2oYrm6k=
@@ -357,8 +375,8 @@ golang.org/x/net v0.0.0-20200520004742-59133d7f0dd7/go.mod h1:qpuaurCH72eLCgpAm/
 golang.org/x/net v0.0.0-20201021035429-f5854403a974/go.mod h1:sp8m0HH+o8qH0wwXwYZr8TS3Oi6o0r6Gce1SSxlDquU=
 golang.org/x/net v0.0.0-20210226172049-e18ecbb05110/go.mod h1:m0MpNAwzfU5UDzcl9v0D8zg8gWTRqZa9RBIspLL5mdg=
 golang.org/x/net v0.0.0-20210410081132-afb366fc7cd1/go.mod h1:9tjilg8BloeKEkVJvy7fQ90B1CfIiPueXVOjqfkSzI8=
-golang.org/x/net v0.28.0 h1:a9JDOJc5GMUJ0+UDqmLT86WiEy7iWyIhz8gz8E4e5hE=
-golang.org/x/net v0.28.0/go.mod h1:yqtgsTWOOnlGLG9GFRrK3++bGOUEkNBoHZc8MEDWPNg=
+golang.org/x/net v0.30.0 h1:AcW1SDZMkb8IpzCdQUaIq2sP4sZ4zw+55h6ynffypl4=
+golang.org/x/net v0.30.0/go.mod h1:2wGyMJ5iFasEhkwi13ChkO/t1ECNC4X4eBKkVFyYFlU=
 golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be/go.mod h1:N/0e6XlmueqKjAGxoOufVs8QHGRruUQn6yWY3a++T0U=
 golang.org/x/sync v0.0.0-20180314180146-1d60e4601c6f/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sync v0.0.0-20181108010431-42b317875d0f/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
@@ -399,17 +417,17 @@ golang.org/x/sys v0.0.0-20220728004956-3c1f35247d10/go.mod h1:oPkhp1MJrh7nUepCBc
 golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.5.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.6.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
-golang.org/x/sys v0.24.0 h1:Twjiwq9dn6R1fQcyiK+wQyHWfaz/BJB+YIpzU/Cv3Xg=
-golang.org/x/sys v0.24.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
+golang.org/x/sys v0.26.0 h1:KHjCJyddX0LoSTb3J+vWpupP9p0oznkqVk/IfjymZbo=
+golang.org/x/sys v0.26.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
 golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1/go.mod h1:bj7SfCRtBDWHUb9snDiAeCFNEtKQo2Wmx5Cou7ajbmo=
-golang.org/x/term v0.23.0 h1:F6D4vR+EHoL9/sWAWgAR1H2DcHr4PareCbAaCo1RpuU=
-golang.org/x/term v0.23.0/go.mod h1:DgV24QBUrK6jhZXl+20l6UWznPlwAHm1Q1mGHtydmSk=
+golang.org/x/term v0.25.0 h1:WtHI/ltw4NvSUig5KARz9h521QvRC8RmF/cuYqifU24=
+golang.org/x/term v0.25.0/go.mod h1:RPyXicDX+6vLxogjjRxjgD2TKtmAO6NZBsBRfrOLu7M=
 golang.org/x/text v0.3.0/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
 golang.org/x/text v0.3.2/go.mod h1:bEr9sfX3Q8Zfm5fL9x+3itogRgK3+ptLWKqgva+5dAk=
 golang.org/x/text v0.3.3/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
 golang.org/x/text v0.3.6/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
-golang.org/x/text v0.17.0 h1:XtiM5bkSOt+ewxlOE/aE/AKEHibwj/6gvWMl9Rsh0Qc=
-golang.org/x/text v0.17.0/go.mod h1:BuEKDfySbSR4drPmRPG/7iBdf8hvFMuRexcpahXilzY=
+golang.org/x/text v0.19.0 h1:kTxAhCbGbxhK0IwgSKiMO5awPoDQ0RpfiVYBfK860YM=
+golang.org/x/text v0.19.0/go.mod h1:BuEKDfySbSR4drPmRPG/7iBdf8hvFMuRexcpahXilzY=
 golang.org/x/tools v0.0.0-20180917221912-90fa682c2a6e/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
 golang.org/x/tools v0.0.0-20190114222345-bf090417da8b/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
 golang.org/x/tools v0.0.0-20190226205152-f727befe758c/go.mod h1:9Yl7xja0Znq3iFh3HoIrodX9oNMXvdceNzlUR8zjMvY=
@@ -472,4 +490,3 @@ gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
 honnef.co/go/tools v0.0.0-20190102054323-c2f93a96b099/go.mod h1:rf3lG4BRIbNafJWhAfAdb/ePZxsR/4RtNHQocxwk9r4=
 honnef.co/go/tools v0.0.0-20190523083050-ea95bdfd59fc/go.mod h1:rf3lG4BRIbNafJWhAfAdb/ePZxsR/4RtNHQocxwk9r4=
 nullprogram.com/x/optparse v1.0.0/go.mod h1:KdyPE+Igbe0jQUrVfMqDMeJQIJZEuyV7pjYmp6pbG50=
-rsc.io/pdf v0.1.1/go.mod h1:n8OzWcQ6Sp37PL01nO98y4iUCRdTGarVfzxY20ICaU4=
diff --git a/internal/common/server/gprc.go b/internal/common/server/gprc.go
index 817cfcb..14f3132 100644
--- a/internal/common/server/gprc.go
+++ b/internal/common/server/gprc.go
@@ -7,6 +7,7 @@ import (
 	grpc_tags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
 	"github.com/sirupsen/logrus"
 	"github.com/spf13/viper"
+	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
 	"google.golang.org/grpc"
 )
 
@@ -28,23 +29,14 @@ func RunGRPCServer(serviceName string, registerServer func(server *grpc.Server))
 func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server)) {
 	logrusEntry := logrus.NewEntry(logrus.StandardLogger())
 	grpcServer := grpc.NewServer(
+		grpc.StatsHandler(otelgrpc.NewServerHandler()),
 		grpc.ChainUnaryInterceptor(
 			grpc_tags.UnaryServerInterceptor(grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor)),
 			grpc_logrus.UnaryServerInterceptor(logrusEntry),
-			//otelgrpc.UnaryServerInterceptor(),
-			//srvMetrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
-			//logging.UnaryServerInterceptor(interceptorLogger(rpcLogger), logging.WithFieldsFromContext(logTraceID)),
-			//selector.UnaryServerInterceptor(auth.UnaryServerInterceptor(authFn), selector.MatchFunc(allButHealthZ)),
-			//recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
 		),
 		grpc.ChainStreamInterceptor(
 			grpc_tags.StreamServerInterceptor(grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor)),
 			grpc_logrus.StreamServerInterceptor(logrusEntry),
-			//otelgrpc.StreamServerInterceptor(),
-			//srvMetrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
-			//logging.StreamServerInterceptor(interceptorLogger(rpcLogger), logging.WithFieldsFromContext(logTraceID)),
-			//selector.StreamServerInterceptor(auth.StreamServerInterceptor(authFn), selector.MatchFunc(allButHealthZ)),
-			//recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
 		),
 	)
 	registerServer(grpcServer)
diff --git a/internal/common/server/http.go b/internal/common/server/http.go
index f22dce4..d1d93f5 100644
--- a/internal/common/server/http.go
+++ b/internal/common/server/http.go
@@ -3,6 +3,7 @@ package server
 import (
 	"github.com/gin-gonic/gin"
 	"github.com/spf13/viper"
+	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
 )
 
 func RunHTTPServer(serviceName string, wrapper func(router *gin.Engine)) {
@@ -15,9 +16,15 @@ func RunHTTPServer(serviceName string, wrapper func(router *gin.Engine)) {
 
 func RunHTTPServerOnAddr(addr string, wrapper func(router *gin.Engine)) {
 	apiRouter := gin.New()
+	setMiddlewares(apiRouter)
 	wrapper(apiRouter)
 	apiRouter.Group("/api")
 	if err := apiRouter.Run(addr); err != nil {
 		panic(err)
 	}
 }
+
+func setMiddlewares(r *gin.Engine) {
+	r.Use(gin.Recovery())
+	r.Use(otelgin.Middleware("default_server"))
+}
diff --git a/internal/common/tracing/jaeger.go b/internal/common/tracing/jaeger.go
new file mode 100644
index 0000000..3374da7
--- /dev/null
+++ b/internal/common/tracing/jaeger.go
@@ -0,0 +1,49 @@
+package tracing
+
+import (
+	"context"
+
+	"go.opentelemetry.io/contrib/propagators/b3"
+	"go.opentelemetry.io/otel"
+	"go.opentelemetry.io/otel/exporters/jaeger"
+	"go.opentelemetry.io/otel/propagation"
+	"go.opentelemetry.io/otel/sdk/resource"
+	sdktrace "go.opentelemetry.io/otel/sdk/trace"
+	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
+	"go.opentelemetry.io/otel/trace"
+)
+
+var tracer = otel.Tracer("default_tracer")
+
+func InitJaegerProvider(jaegerURL, serviceName string) (func(ctx context.Context) error, error) {
+	if jaegerURL == "" {
+		panic("empty jaeger url")
+	}
+	tracer = otel.Tracer(serviceName)
+	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerURL)))
+	if err != nil {
+		return nil, err
+	}
+	tp := sdktrace.NewTracerProvider(
+		sdktrace.WithBatcher(exp),
+		sdktrace.WithResource(resource.NewSchemaless(
+			semconv.ServiceNameKey.String(serviceName),
+		)),
+	)
+	otel.SetTracerProvider(tp)
+	b3Propagator := b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader))
+	p := propagation.NewCompositeTextMapPropagator(
+		propagation.TraceContext{}, propagation.Baggage{}, b3Propagator,
+	)
+	otel.SetTextMapPropagator(p)
+	return tp.Shutdown, nil
+}
+
+func Start(ctx context.Context, name string) (context.Context, trace.Span) {
+	return tracer.Start(ctx, name)
+}
+
+func TraceID(ctx context.Context) string {
+	spanCtx := trace.SpanContextFromContext(ctx)
+	return spanCtx.TraceID().String()
+}
diff --git a/internal/order/go.mod b/internal/order/go.mod
index 99f2d66..5022f4f 100644
--- a/internal/order/go.mod
+++ b/internal/order/go.mod
@@ -29,6 +29,8 @@ require (
 	github.com/fsnotify/fsnotify v1.7.0 // indirect
 	github.com/gabriel-vasile/mimetype v1.4.6 // indirect
 	github.com/gin-contrib/sse v0.1.0 // indirect
+	github.com/go-logr/logr v1.4.2 // indirect
+	github.com/go-logr/stdr v1.2.2 // indirect
 	github.com/go-playground/locales v0.14.1 // indirect
 	github.com/go-playground/universal-translator v0.18.1 // indirect
 	github.com/go-playground/validator/v10 v10.22.1 // indirect
@@ -68,6 +70,14 @@ require (
 	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
 	github.com/ugorji/go/codec v1.2.12 // indirect
 	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
+	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0 // indirect
+	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0 // indirect
+	go.opentelemetry.io/contrib/propagators/b3 v1.31.0 // indirect
+	go.opentelemetry.io/otel v1.31.0 // indirect
+	go.opentelemetry.io/otel/exporters/jaeger v1.17.0 // indirect
+	go.opentelemetry.io/otel/metric v1.31.0 // indirect
+	go.opentelemetry.io/otel/sdk v1.31.0 // indirect
+	go.opentelemetry.io/otel/trace v1.31.0 // indirect
 	go.uber.org/multierr v1.11.0 // indirect
 	golang.org/x/arch v0.11.0 // indirect
 	golang.org/x/crypto v0.28.0 // indirect
diff --git a/internal/order/go.sum b/internal/order/go.sum
index 9f6846a..09409af 100644
--- a/internal/order/go.sum
+++ b/internal/order/go.sum
@@ -65,6 +65,11 @@ github.com/go-kit/log v0.1.0/go.mod h1:zbhenjAZHb184qTLMA9ZjW7ThYL0H2mk7Q6pNt4vb
 github.com/go-logfmt/logfmt v0.3.0/go.mod h1:Qt1PoO58o5twSAckw1HlFXLmHsOX5/0LbT9GBnD5lWE=
 github.com/go-logfmt/logfmt v0.4.0/go.mod h1:3RMwSq7FuexP4Kalkev3ejPJsZTpXXBr9+V4qmtdjCk=
 github.com/go-logfmt/logfmt v0.5.0/go.mod h1:wCYkCAKZfumFQihp8CzCvQ3paCTfi41vtzG1KdI/P7A=
+github.com/go-logr/logr v1.2.2/go.mod h1:jdQByPbusPIv2/zmleS9BjJVeZ6kBagPoEUsqbVz/1A=
+github.com/go-logr/logr v1.4.2 h1:6pFjapn8bFcIbiKo3XT4j/BhANplGihG6tvd+8rYgrY=
+github.com/go-logr/logr v1.4.2/go.mod h1:9T104GzyrTigFIr8wt5mBrctHMim0Nb2HLGrmQ40KvY=
+github.com/go-logr/stdr v1.2.2 h1:hSWxHoqTgW2S2qGc0LTAI563KZ5YKYRhT3MFKZMbjag=
+github.com/go-logr/stdr v1.2.2/go.mod h1:mMo/vtBO5dYbehREoey6XUKy/eSumjCCveDpRre4VKE=
 github.com/go-playground/assert/v2 v2.2.0 h1:JvknZsQTYeFEAhQwI4qEt9cyV5ONwRHC+lYKSsYSR8s=
 github.com/go-playground/assert/v2 v2.2.0/go.mod h1:VDjEfimB/XKnb+ZQfWdccd7VUvScMdVu0Titje2rxJ4=
 github.com/go-playground/locales v0.14.1 h1:EWaQ/wswjilfKLTECiXz7Rh+3BjFhfDFKv/oXslEjJA=
@@ -295,6 +300,22 @@ github.com/x-cray/logrus-prefixed-formatter v0.5.2 h1:00txxvfBM9muc0jiLIEAkAcIMJ
 github.com/x-cray/logrus-prefixed-formatter v0.5.2/go.mod h1:2duySbKsL6M18s5GU7VPsoEPHyzalCE06qoARUCeBBE=
 github.com/yuin/goldmark v1.1.27/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
 github.com/yuin/goldmark v1.2.1/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
+go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0 h1:0nTRpaCaILLdooXAQnfktlL6Zw1ECKEW9DZGH2byi2c=
+go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0/go.mod h1:A7aFlp4WSLmeOnFRZwf2dMU+40THPc+rsr6KOwZLOcg=
+go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0 h1:4Pp6oUg3+e/6M4C0A/3kJ2VYa++dsWVTtGgLVj5xtHg=
+go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0/go.mod h1:Mjt1i1INqiaoZOMGR1RIUJN+i3ChKoFRqzrRQhlkbs0=
+go.opentelemetry.io/contrib/propagators/b3 v1.31.0 h1:PQPXYscmwbCp76QDvO4hMngF2j8Bx/OTV86laEl8uqo=
+go.opentelemetry.io/contrib/propagators/b3 v1.31.0/go.mod h1:jbqfV8wDdqSDrAYxVpXQnpM0XFMq2FtDesblJ7blOwQ=
+go.opentelemetry.io/otel v1.31.0 h1:NsJcKPIW0D0H3NgzPDHmo0WW6SptzPdqg/L1zsIm2hY=
+go.opentelemetry.io/otel v1.31.0/go.mod h1:O0C14Yl9FgkjqcCZAsE053C13OaddMYr/hz6clDkEJE=
+go.opentelemetry.io/otel/exporters/jaeger v1.17.0 h1:D7UpUy2Xc2wsi1Ras6V40q806WM07rqoCWzXu7Sqy+4=
+go.opentelemetry.io/otel/exporters/jaeger v1.17.0/go.mod h1:nPCqOnEH9rNLKqH/+rrUjiMzHJdV1BlpKcTwRTyKkKI=
+go.opentelemetry.io/otel/metric v1.31.0 h1:FSErL0ATQAmYHUIzSezZibnyVlft1ybhy4ozRPcF2fE=
+go.opentelemetry.io/otel/metric v1.31.0/go.mod h1:C3dEloVbLuYoX41KpmAhOqNriGbA+qqH6PQ5E5mUfnY=
+go.opentelemetry.io/otel/sdk v1.31.0 h1:xLY3abVHYZ5HSfOg3l2E5LUj2Cwva5Y7yGxnSW9H5Gk=
+go.opentelemetry.io/otel/sdk v1.31.0/go.mod h1:TfRbMdhvxIIr/B2N2LQW2S5v9m3gOQ/08KsbbO5BPT0=
+go.opentelemetry.io/otel/trace v1.31.0 h1:ffjsj1aRouKewfr85U2aGagJ46+MvodynlQ1HYdmJys=
+go.opentelemetry.io/otel/trace v1.31.0/go.mod h1:TXZkRk7SM2ZQLtR6eoAWQFIHPvzQ06FJAsO1tJg480A=
 go.uber.org/atomic v1.7.0/go.mod h1:fEN4uk6kAWBTFdckzkM89CLk9XfWZrxpCo0nPH17wJc=
 go.uber.org/goleak v1.1.10/go.mod h1:8a7PlsEVH3e/a/GLqe5IIrQx6GzcnRmZEufDUTk4A7A=
 go.uber.org/goleak v1.3.0 h1:2K3zAYmnTNqV73imy9J1T3WC+gmCePx2hEGkimedGto=
diff --git a/internal/order/http.go b/internal/order/http.go
index 2073a4f..413cdc1 100644
--- a/internal/order/http.go
+++ b/internal/order/http.go
@@ -5,6 +5,7 @@ import (
 	"net/http"
 
 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+	"github.com/ghost-yu/go_shop_second/common/tracing"
 	"github.com/ghost-yu/go_shop_second/order/app"
 	"github.com/ghost-yu/go_shop_second/order/app/command"
 	"github.com/ghost-yu/go_shop_second/order/app/query"
@@ -16,12 +17,15 @@ type HTTPServer struct {
 }
 
 func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
+	ctx, span := tracing.Start(c, "PostCustomerCustomerIDOrders")
+	defer span.End()
+
 	var req orderpb.CreateOrderRequest
 	if err := c.ShouldBindJSON(&req); err != nil {
 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
 		return
 	}
-	r, err := H.app.Commands.CreateOrder.Handle(c, command.CreateOrder{
+	r, err := H.app.Commands.CreateOrder.Handle(ctx, command.CreateOrder{
 		CustomerID: req.CustomerID,
 		Items:      req.Items,
 	})
@@ -31,6 +35,7 @@ func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID stri
 	}
 	c.JSON(http.StatusOK, gin.H{
 		"message":      "success",
+		"trace_id":     tracing.TraceID(ctx),
 		"customer_id":  req.CustomerID,
 		"order_id":     r.OrderID,
 		"redirect_url": fmt.Sprintf("http://localhost:8282/success?customerID=%s&orderID=%s", req.CustomerID, r.OrderID),
@@ -38,7 +43,10 @@ func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID stri
 }
 
 func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
-	o, err := H.app.Queries.GetCustomerOrder.Handle(c, query.GetCustomerOrder{
+	ctx, span := tracing.Start(c, "GetCustomerCustomerIDOrdersOrderID")
+	defer span.End()
+
+	o, err := H.app.Queries.GetCustomerOrder.Handle(ctx, query.GetCustomerOrder{
 		OrderID:    orderID,
 		CustomerID: customerID,
 	})
@@ -47,7 +55,8 @@ func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerI
 		return
 	}
 	c.JSON(http.StatusOK, gin.H{
-		"message": "success",
+		"message":  "success",
+		"trace_id": tracing.TraceID(ctx),
 		"data": gin.H{
 			"Order": o,
 		},
diff --git a/internal/order/main.go b/internal/order/main.go
index 4453aaf..04fb795 100644
--- a/internal/order/main.go
+++ b/internal/order/main.go
@@ -9,6 +9,7 @@ import (
 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
 	"github.com/ghost-yu/go_shop_second/common/logging"
 	"github.com/ghost-yu/go_shop_second/common/server"
+	"github.com/ghost-yu/go_shop_second/common/tracing"
 	"github.com/ghost-yu/go_shop_second/order/infrastructure/consumer"
 	"github.com/ghost-yu/go_shop_second/order/ports"
 	"github.com/ghost-yu/go_shop_second/order/service"
@@ -31,6 +32,12 @@ func main() {
 	ctx, cancel := context.WithCancel(context.Background())
 	defer cancel()
 
+	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
+	if err != nil {
+		logrus.Fatal(err)
+	}
+	defer shutdown(ctx)
+
 	application, cleanup := service.NewApplication(ctx)
 	defer cleanup()
 
diff --git a/internal/payment/go.mod b/internal/payment/go.mod
index 7b26a13..db456a2 100644
--- a/internal/payment/go.mod
+++ b/internal/payment/go.mod
@@ -6,7 +6,7 @@ replace github.com/ghost-yu/go_shop_second/common => ../common
 
 require (
 	github.com/ghost-yu/go_shop_second/common v0.0.0-00010101000000-000000000000
-	github.com/gin-gonic/gin v1.9.1
+	github.com/gin-gonic/gin v1.10.0
 	github.com/rabbitmq/amqp091-go v1.10.0
 	github.com/sirupsen/logrus v1.8.1
 	github.com/spf13/viper v1.19.0
@@ -15,18 +15,22 @@ require (
 
 require (
 	github.com/armon/go-metrics v0.4.1 // indirect
-	github.com/bytedance/sonic v1.10.0-rc3 // indirect
-	github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d // indirect
-	github.com/chenzhuoyu/iasm v0.9.0 // indirect
+	github.com/bytedance/sonic v1.12.3 // indirect
+	github.com/bytedance/sonic/loader v0.2.0 // indirect
+	github.com/cloudwego/base64x v0.1.4 // indirect
+	github.com/cloudwego/iasm v0.2.0 // indirect
 	github.com/fatih/color v1.14.1 // indirect
 	github.com/fsnotify/fsnotify v1.7.0 // indirect
-	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
+	github.com/gabriel-vasile/mimetype v1.4.5 // indirect
 	github.com/gin-contrib/sse v0.1.0 // indirect
+	github.com/go-logr/logr v1.4.2 // indirect
+	github.com/go-logr/stdr v1.2.2 // indirect
 	github.com/go-playground/locales v0.14.1 // indirect
 	github.com/go-playground/universal-translator v0.18.1 // indirect
-	github.com/go-playground/validator/v10 v10.14.1 // indirect
-	github.com/goccy/go-json v0.10.2 // indirect
+	github.com/go-playground/validator/v10 v10.22.1 // indirect
+	github.com/goccy/go-json v0.10.3 // indirect
 	github.com/golang/protobuf v1.5.3 // indirect
+	github.com/google/uuid v1.6.0 // indirect
 	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
 	github.com/hashicorp/consul/api v1.28.2 // indirect
 	github.com/hashicorp/errwrap v1.1.0 // indirect
@@ -39,8 +43,8 @@ require (
 	github.com/hashicorp/hcl v1.0.0 // indirect
 	github.com/hashicorp/serf v0.10.1 // indirect
 	github.com/json-iterator/go v1.1.12 // indirect
-	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
-	github.com/leodido/go-urn v1.2.4 // indirect
+	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
+	github.com/leodido/go-urn v1.4.0 // indirect
 	github.com/magiconair/properties v1.8.7 // indirect
 	github.com/mattn/go-colorable v0.1.13 // indirect
 	github.com/mattn/go-isatty v0.0.20 // indirect
@@ -50,7 +54,7 @@ require (
 	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
 	github.com/modern-go/reflect2 v1.0.2 // indirect
 	github.com/nxadm/tail v1.4.11 // indirect
-	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
+	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
 	github.com/sagikazarmark/locafero v0.4.0 // indirect
 	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
 	github.com/sourcegraph/conc v0.3.0 // indirect
@@ -59,17 +63,25 @@ require (
 	github.com/spf13/pflag v1.0.5 // indirect
 	github.com/subosito/gotenv v1.6.0 // indirect
 	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
-	github.com/ugorji/go/codec v1.2.11 // indirect
+	github.com/ugorji/go/codec v1.2.12 // indirect
 	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
+	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0 // indirect
+	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0 // indirect
+	go.opentelemetry.io/contrib/propagators/b3 v1.31.0 // indirect
+	go.opentelemetry.io/otel v1.31.0 // indirect
+	go.opentelemetry.io/otel/exporters/jaeger v1.17.0 // indirect
+	go.opentelemetry.io/otel/metric v1.31.0 // indirect
+	go.opentelemetry.io/otel/sdk v1.31.0 // indirect
+	go.opentelemetry.io/otel/trace v1.31.0 // indirect
 	go.uber.org/atomic v1.9.0 // indirect
 	go.uber.org/multierr v1.9.0 // indirect
-	golang.org/x/arch v0.4.0 // indirect
-	golang.org/x/crypto v0.26.0 // indirect
+	golang.org/x/arch v0.11.0 // indirect
+	golang.org/x/crypto v0.28.0 // indirect
 	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
-	golang.org/x/net v0.28.0 // indirect
-	golang.org/x/sys v0.24.0 // indirect
-	golang.org/x/term v0.23.0 // indirect
-	golang.org/x/text v0.17.0 // indirect
+	golang.org/x/net v0.30.0 // indirect
+	golang.org/x/sys v0.26.0 // indirect
+	golang.org/x/term v0.25.0 // indirect
+	golang.org/x/text v0.19.0 // indirect
 	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
 	google.golang.org/grpc v1.67.1 // indirect
 	google.golang.org/protobuf v1.35.1 // indirect
diff --git a/internal/payment/go.sum b/internal/payment/go.sum
index 05d1402..8883881 100644
--- a/internal/payment/go.sum
+++ b/internal/payment/go.sum
@@ -16,21 +16,20 @@ github.com/beorn7/perks v0.0.0-20180321164747-3a771d992973/go.mod h1:Dwedo/Wpr24
 github.com/beorn7/perks v1.0.0/go.mod h1:KWe93zE9D1o94FZ5RNwFwVgaQK1VOXiVxmqh+CedLV8=
 github.com/beorn7/perks v1.0.1/go.mod h1:G2ZrVWU2WbWT9wwq4/hrbKbnv/1ERSJQ0ibhJ6rlkpw=
 github.com/bgentry/speakeasy v0.1.0/go.mod h1:+zsyZBPWlz7T6j88CTgSN5bM796AkVf0kBD4zp0CCIs=
-github.com/bytedance/sonic v1.5.0/go.mod h1:ED5hyg4y6t3/9Ku1R6dU/4KyJ48DZ4jPhfY1O2AihPM=
-github.com/bytedance/sonic v1.10.0-rc/go.mod h1:ElCzW+ufi8qKqNW0FY314xriJhyJhuoJ3gFZdAHF7NM=
-github.com/bytedance/sonic v1.10.0-rc3 h1:uNSnscRapXTwUgTyOF0GVljYD08p9X/Lbr9MweSV3V0=
-github.com/bytedance/sonic v1.10.0-rc3/go.mod h1:iZcSUejdk5aukTND/Eu/ivjQuEL0Cu9/rf50Hi0u/g4=
+github.com/bytedance/sonic v1.12.3 h1:W2MGa7RCU1QTeYRTPE3+88mVC0yXmsRQRChiyVocVjU=
+github.com/bytedance/sonic v1.12.3/go.mod h1:B8Gt/XvtZ3Fqj+iSKMypzymZxw/FVwgIGKzMzT9r/rk=
+github.com/bytedance/sonic/loader v0.1.1/go.mod h1:ncP89zfokxS5LZrJxl5z0UJcsk4M4yY2JpfqGeCtNLU=
+github.com/bytedance/sonic/loader v0.2.0 h1:zNprn+lsIP06C/IqCHs3gPQIvnvpKbbxyXQP1iU4kWM=
+github.com/bytedance/sonic/loader v0.2.0/go.mod h1:ncP89zfokxS5LZrJxl5z0UJcsk4M4yY2JpfqGeCtNLU=
 github.com/census-instrumentation/opencensus-proto v0.2.1/go.mod h1:f6KPmirojxKA12rnyqOA5BBL4O983OfeGPqjHWSTneU=
 github.com/cespare/xxhash/v2 v2.1.1/go.mod h1:VGX0DQ3Q6kWi7AoAeZDth3/j3BFtOZR5XLFGgcrjCOs=
-github.com/chenzhuoyu/base64x v0.0.0-20211019084208-fb5309c8db06/go.mod h1:DH46F32mSOjUmXrMHnKwZdA8wcEefY7UVqBKYGjpdQY=
-github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311/go.mod h1:b583jCggY9gE99b6G5LEC39OIiVsWj+R97kbl5odCEk=
-github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d h1:77cEq6EriyTZ0g/qfRdp61a3Uu/AWrgIq2s0ClJV1g0=
-github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d/go.mod h1:8EPpVsBuRksnlj1mLy4AWzRNQYxauNi62uWcE3to6eA=
-github.com/chenzhuoyu/iasm v0.9.0 h1:9fhXjVzq5hUy2gkhhgHl95zG2cEAhw9OSGs8toWWAwo=
-github.com/chenzhuoyu/iasm v0.9.0/go.mod h1:Xjy2NpN3h7aUqeqM+woSuuvxmIe6+DDsiNLIrkAmYog=
 github.com/circonus-labs/circonus-gometrics v2.3.1+incompatible/go.mod h1:nmEj6Dob7S7YxXgwXpfOuvO54S+tGdZdw9fuRZt25Ag=
 github.com/circonus-labs/circonusllhist v0.1.3/go.mod h1:kMXHVDlOchFAehlya5ePtbp5jckzBHf4XRpQvBOLI+I=
 github.com/client9/misspell v0.3.4/go.mod h1:qj6jICC3Q7zFZvVWo7KLAzC3yx5G7kyvSDkc90ppPyw=
+github.com/cloudwego/base64x v0.1.4 h1:jwCgWpFanWmN8xoIUHa2rtzmkd5J2plF/dnLS6Xd/0Y=
+github.com/cloudwego/base64x v0.1.4/go.mod h1:0zlkT4Wn5C6NdauXdJRhSKRlJvmclQ1hhJgA0rcu/8w=
+github.com/cloudwego/iasm v0.2.0 h1:1KNIy1I1H9hNNFEEH3DVnI4UujN+1zjpuk6gwHLTssg=
+github.com/cloudwego/iasm v0.2.0/go.mod h1:8rXZaNYT2n95jn+zTI1sDr+IgcD2GVs0nlbbQPiEFhY=
 github.com/cncf/udpa/go v0.0.0-20191209042840-269d4d468f6f/go.mod h1:M8M6+tZqaGXZJjfX53e64911xZQV5JYwmTeXPW+k8Sc=
 github.com/davecgh/go-spew v1.1.0/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
 github.com/davecgh/go-spew v1.1.1/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
@@ -50,29 +49,34 @@ github.com/frankban/quicktest v1.14.6/go.mod h1:4ptaffx2x8+WTWXmUCuVU6aPUX1/Mz7z
 github.com/fsnotify/fsnotify v1.6.0/go.mod h1:sl3t1tCWJFWoRz9R8WJCbQihKKwmorjAbSClcnxKAGw=
 github.com/fsnotify/fsnotify v1.7.0 h1:8JEhPFa5W2WU7YfeZzPNqzMP6Lwt7L2715Ggo0nosvA=
 github.com/fsnotify/fsnotify v1.7.0/go.mod h1:40Bi/Hjc2AVfZrqy+aj+yEI+/bRxZnMJyTJwOpGvigM=
-github.com/gabriel-vasile/mimetype v1.4.2 h1:w5qFW6JKBz9Y393Y4q372O9A7cUSequkh1Q7OhCmWKU=
-github.com/gabriel-vasile/mimetype v1.4.2/go.mod h1:zApsH/mKG4w07erKIaJPFiX0Tsq9BFQgN3qGY5GnNgA=
+github.com/gabriel-vasile/mimetype v1.4.5 h1:J7wGKdGu33ocBOhGy0z653k/lFKLFDPJMG8Gql0kxn4=
+github.com/gabriel-vasile/mimetype v1.4.5/go.mod h1:ibHel+/kbxn9x2407k1izTA1S81ku1z/DlgOW2QE0M4=
 github.com/gin-contrib/sse v0.1.0 h1:Y/yl/+YNO8GZSjAhjMsSuLt29uWRFHdHYUb5lYOV9qE=
 github.com/gin-contrib/sse v0.1.0/go.mod h1:RHrZQHXnP2xjPF+u1gW/2HnVO7nvIa9PG3Gm+fLHvGI=
-github.com/gin-gonic/gin v1.9.1 h1:4idEAncQnU5cB7BeOkPtxjfCSye0AAm1R0RVIqJ+Jmg=
-github.com/gin-gonic/gin v1.9.1/go.mod h1:hPrL7YrpYKXt5YId3A/Tnip5kqbEAP+KLuI3SUcPTeU=
+github.com/gin-gonic/gin v1.10.0 h1:nTuyha1TYqgedzytsKYqna+DfLos46nTv2ygFy86HFU=
+github.com/gin-gonic/gin v1.10.0/go.mod h1:4PMNQiOhvDRa013RKVbsiNwoyezlm2rm0uX/T7kzp5Y=
 github.com/go-kit/kit v0.8.0/go.mod h1:xBxKIO96dXMWWy0MnWVtmwkA9/13aqxPnvrjFYMA2as=
 github.com/go-kit/kit v0.9.0/go.mod h1:xBxKIO96dXMWWy0MnWVtmwkA9/13aqxPnvrjFYMA2as=
 github.com/go-kit/log v0.1.0/go.mod h1:zbhenjAZHb184qTLMA9ZjW7ThYL0H2mk7Q6pNt4vbaY=
 github.com/go-logfmt/logfmt v0.3.0/go.mod h1:Qt1PoO58o5twSAckw1HlFXLmHsOX5/0LbT9GBnD5lWE=
 github.com/go-logfmt/logfmt v0.4.0/go.mod h1:3RMwSq7FuexP4Kalkev3ejPJsZTpXXBr9+V4qmtdjCk=
 github.com/go-logfmt/logfmt v0.5.0/go.mod h1:wCYkCAKZfumFQihp8CzCvQ3paCTfi41vtzG1KdI/P7A=
+github.com/go-logr/logr v1.2.2/go.mod h1:jdQByPbusPIv2/zmleS9BjJVeZ6kBagPoEUsqbVz/1A=
+github.com/go-logr/logr v1.4.2 h1:6pFjapn8bFcIbiKo3XT4j/BhANplGihG6tvd+8rYgrY=
+github.com/go-logr/logr v1.4.2/go.mod h1:9T104GzyrTigFIr8wt5mBrctHMim0Nb2HLGrmQ40KvY=
+github.com/go-logr/stdr v1.2.2 h1:hSWxHoqTgW2S2qGc0LTAI563KZ5YKYRhT3MFKZMbjag=
+github.com/go-logr/stdr v1.2.2/go.mod h1:mMo/vtBO5dYbehREoey6XUKy/eSumjCCveDpRre4VKE=
 github.com/go-playground/assert/v2 v2.2.0 h1:JvknZsQTYeFEAhQwI4qEt9cyV5ONwRHC+lYKSsYSR8s=
 github.com/go-playground/assert/v2 v2.2.0/go.mod h1:VDjEfimB/XKnb+ZQfWdccd7VUvScMdVu0Titje2rxJ4=
 github.com/go-playground/locales v0.14.1 h1:EWaQ/wswjilfKLTECiXz7Rh+3BjFhfDFKv/oXslEjJA=
 github.com/go-playground/locales v0.14.1/go.mod h1:hxrqLVvrK65+Rwrd5Fc6F2O76J/NuW9t0sjnWqG1slY=
 github.com/go-playground/universal-translator v0.18.1 h1:Bcnm0ZwsGyWbCzImXv+pAJnYK9S473LQFuzCbDbfSFY=
 github.com/go-playground/universal-translator v0.18.1/go.mod h1:xekY+UJKNuX9WP91TpwSH2VMlDf28Uj24BCp08ZFTUY=
-github.com/go-playground/validator/v10 v10.14.1 h1:9c50NUPC30zyuKprjL3vNZ0m5oG+jU0zvx4AqHGnv4k=
-github.com/go-playground/validator/v10 v10.14.1/go.mod h1:9iXMNT7sEkjXb0I+enO7QXmzG6QCsPWY4zveKFVRSyU=
+github.com/go-playground/validator/v10 v10.22.1 h1:40JcKH+bBNGFczGuoBYgX4I6m/i27HYW8P9FDk5PbgA=
+github.com/go-playground/validator/v10 v10.22.1/go.mod h1:dbuPbCMFw/DrkbEynArYaCwl3amGuJotoKCe95atGMM=
 github.com/go-stack/stack v1.8.0/go.mod h1:v0f6uXyyMGvRgIKkXu+yp6POWl0qKG85gN/melR3HDY=
-github.com/goccy/go-json v0.10.2 h1:CrxCmQqYDkv1z7lO7Wbh2HN93uovUHgrECaO5ZrCXAU=
-github.com/goccy/go-json v0.10.2/go.mod h1:6MelG93GURQebXPDq3khkgXZkazVtN9CRI+MGFi0w8I=
+github.com/goccy/go-json v0.10.3 h1:KZ5WoDbxAIgm2HNbYckL0se1fHD6rz5j4ywS6ebzDqA=
+github.com/goccy/go-json v0.10.3/go.mod h1:oq7eo15ShAhp70Anwd5lgX2pLfOS3QCiwU/PULtXL6M=
 github.com/gogo/protobuf v1.1.1/go.mod h1:r8qH/GZQm5c6nD/R0oafs1akxWv10x8SbQlK7atdtwQ=
 github.com/gogo/protobuf v1.3.2 h1:Ov1cvc58UF3b5XjBnZv7+opcTcQFZebYjWzi34vdm4Q=
 github.com/gogo/protobuf v1.3.2/go.mod h1:P1XiOD3dCwIKUDQYPy72D8LYyHL2YPYrpS2s69NZV8Q=
@@ -95,6 +99,8 @@ github.com/google/go-cmp v0.5.5/go.mod h1:v8dTdLbMG2kIc/vJvl+f65V22dbkXbowE6jgT/
 github.com/google/go-cmp v0.6.0 h1:ofyhxvXcZhMsU5ulbFiLKl/XBFqE1GSq7atu8tAmTRI=
 github.com/google/go-cmp v0.6.0/go.mod h1:17dUlkBOakJ0+DkrSSNjCkIjxS6bF9zb3elmeNGIjoY=
 github.com/google/gofuzz v1.0.0/go.mod h1:dBl0BpW6vV/+mYPU4Po3pmUjxk6FQPldtuIdl/M65Eg=
+github.com/google/uuid v1.6.0 h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=
+github.com/google/uuid v1.6.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
 github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 h1:UH//fgunKIs4JdUbpDl1VZCDaL56wXCB/5+wF6uHfaI=
 github.com/grpc-ecosystem/go-grpc-middleware v1.4.0/go.mod h1:g5qyo/la0ALbONm6Vbp88Yd8NsDy6rZz+RcrMPxvld8=
 github.com/hashicorp/consul/api v1.28.2 h1:mXfkRHrpHN4YY3RqL09nXU1eHKLNiuAN4kHvDQ16k/8=
@@ -151,8 +157,8 @@ github.com/julienschmidt/httprouter v1.2.0/go.mod h1:SYymIcj16QtmaHHD7aYtjjsJG7V
 github.com/kisielk/errcheck v1.5.0/go.mod h1:pFxgyoBC7bSaBwPgfKdkLd5X25qrDl4LWUI2bnpBCr8=
 github.com/kisielk/gotool v1.0.0/go.mod h1:XhKaO+MFFWcvkIS/tQcRk01m1F5IRFswLeQ+oQHNcck=
 github.com/klauspost/cpuid/v2 v2.0.9/go.mod h1:FInQzS24/EEf25PyTYn52gqo7WaD8xa0213Md/qVLRg=
-github.com/klauspost/cpuid/v2 v2.2.5 h1:0E5MSMDEoAulmXNFquVs//DdoomxaoTY1kUhbc/qbZg=
-github.com/klauspost/cpuid/v2 v2.2.5/go.mod h1:Lcz8mBdAVJIBVzewtcLocK12l3Y+JytZYpaMropDUws=
+github.com/klauspost/cpuid/v2 v2.2.8 h1:+StwCXwm9PdpiEkPyzBXIy+M9KUb4ODm0Zarf1kS5BM=
+github.com/klauspost/cpuid/v2 v2.2.8/go.mod h1:Lcz8mBdAVJIBVzewtcLocK12l3Y+JytZYpaMropDUws=
 github.com/knz/go-libedit v1.10.1/go.mod h1:MZTVkCWyz0oBc7JOWP3wNAzd002ZbM/5hgShxwh4x8M=
 github.com/konsorten/go-windows-terminal-sequences v1.0.1/go.mod h1:T0+1ngSBFLxvqU3pZ+m/2kptfBszLMUkC4ZK/EgS/cQ=
 github.com/kr/logfmt v0.0.0-20140226030751-b84e30acd515/go.mod h1:+0opPa2QZZtGFBFZlji/RkVcI2GknAs/DXo4wKdlNEc=
@@ -163,8 +169,8 @@ github.com/kr/pty v1.1.1/go.mod h1:pFQYn66WHrOpPYNljwOMqo10TkYh1fy3cYio2l3bCsQ=
 github.com/kr/text v0.1.0/go.mod h1:4Jbv+DJW3UT/LiOwJeYQe1efqtUx/iVham/4vfdArNI=
 github.com/kr/text v0.2.0 h1:5Nx0Ya0ZqY2ygV366QzturHI13Jq95ApcVaJBhpS+AY=
 github.com/kr/text v0.2.0/go.mod h1:eLer722TekiGuMkidMxC/pM04lWEeraHUUmBw8l2grE=
-github.com/leodido/go-urn v1.2.4 h1:XlAE/cm/ms7TE/VMVoduSpNBoyc2dOxHs5MZSwAN63Q=
-github.com/leodido/go-urn v1.2.4/go.mod h1:7ZrI8mTSeBSHl/UaRyKQW1qZeMgak41ANeCNaVckg+4=
+github.com/leodido/go-urn v1.4.0 h1:WT9HwE9SGECu3lg4d/dIA+jxlljEa1/ffXKmRjqdmIQ=
+github.com/leodido/go-urn v1.4.0/go.mod h1:bvxc+MVxLKB4z00jd1z+Dvzr47oO32F/QSNjSBOlFxI=
 github.com/magiconair/properties v1.8.7 h1:IeQXZAiQcpL9mgcAe1Nu6cX9LLw6ExEHKjN0VQdvPDY=
 github.com/magiconair/properties v1.8.7/go.mod h1:Dhd985XPs7jluiymwWYZ0G4Z61jb3vdS329zhj2hYo0=
 github.com/mattn/go-colorable v0.0.9/go.mod h1:9vuHe8Xs5qXnSaW/c/ABM9alt+Vo+STaOChaDxuIBZU=
@@ -212,8 +218,8 @@ github.com/opentracing/opentracing-go v1.1.0/go.mod h1:UkNAQd3GIcIGf0SeVgPpRdFSt
 github.com/pascaldekloe/goe v0.0.0-20180627143212-57f6aae5913c/go.mod h1:lzWF7FIEvWOWxwDKqyGYQf6ZUaNfKdP144TG7ZOy1lc=
 github.com/pascaldekloe/goe v0.1.0 h1:cBOtyMzM9HTpWjXfbbunk26uA6nG3a8n06Wieeh0MwY=
 github.com/pascaldekloe/goe v0.1.0/go.mod h1:lzWF7FIEvWOWxwDKqyGYQf6ZUaNfKdP144TG7ZOy1lc=
-github.com/pelletier/go-toml/v2 v2.2.2 h1:aYUidT7k73Pcl9nb2gScu7NSrKCSHIDE89b3+6Wq+LM=
-github.com/pelletier/go-toml/v2 v2.2.2/go.mod h1:1t835xjRzz80PqgE6HHgN2JOsmgYu/h4qDAS4n929Rs=
+github.com/pelletier/go-toml/v2 v2.2.3 h1:YmeHyLY8mFWbdkNWwpr+qIL2bEqT0o95WSdkNHvL12M=
+github.com/pelletier/go-toml/v2 v2.2.3/go.mod h1:MfCQTFTvCcUyyvvwm1+G6H/jORL20Xlb6rzQu9GuUkc=
 github.com/pkg/errors v0.8.0/go.mod h1:bwawxfHBFNV+L2hUp1rHADufV3IMtnDRdf1r5NINEl0=
 github.com/pkg/errors v0.8.1/go.mod h1:bwawxfHBFNV+L2hUp1rHADufV3IMtnDRdf1r5NINEl0=
 github.com/pkg/errors v0.9.1 h1:FEBLx1zS214owpjy7qsBeixbURkuhQAwrK5UwLGTwt4=
@@ -263,9 +269,8 @@ github.com/spf13/viper v1.19.0/go.mod h1:GQUN9bilAbhU/jgc1bKs99f/suXKeUMct8Adx5+
 github.com/stretchr/objx v0.1.0/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
 github.com/stretchr/objx v0.1.1/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
 github.com/stretchr/objx v0.4.0/go.mod h1:YvHI0jy2hoMjB+UWwv71VJQ9isScKT/TqJzVSSt89Yw=
+github.com/stretchr/objx v0.5.0 h1:1zr/of2m5FGMsad5YfcqgdqdWrIhu+EBEJRhR1U7z/c=
 github.com/stretchr/objx v0.5.0/go.mod h1:Yh+to48EsGEfYuaHDzXPcE3xhTkx73EhmCGUpEOglKo=
-github.com/stretchr/objx v0.5.2 h1:xuMeJ0Sdp5ZMRXx/aWO6RZxdr3beISkG5/G/aIRr3pY=
-github.com/stretchr/objx v0.5.2/go.mod h1:FRsXN1f5AsAjCGJKqEizvkpNtU+EGNCLh3NxZ/8L+MA=
 github.com/stretchr/testify v1.2.2/go.mod h1:a8OnRcib4nhh0OaRAV+Yts87kKdq0PP7pXfy6kDkUVs=
 github.com/stretchr/testify v1.3.0/go.mod h1:M5WIy9Dh21IEIfnGCwXGc5bZfKNJtfHm1UVUgZn+9EI=
 github.com/stretchr/testify v1.4.0/go.mod h1:j7eGeouHqKxXV5pUuKE4zz7dFj8WfuZ+81PSLYec5m4=
@@ -274,8 +279,6 @@ github.com/stretchr/testify v1.7.1/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/
 github.com/stretchr/testify v1.7.2/go.mod h1:R6va5+xMeoiuVRoj+gSkQ7d3FALtqAAGI1FQKckRals=
 github.com/stretchr/testify v1.8.0/go.mod h1:yNjHg4UonilssWZ8iaSj1OCr/vHnekPRkoO+kdMU+MU=
 github.com/stretchr/testify v1.8.1/go.mod h1:w2LPCIKwWwSfY2zedu0+kehJoqGctiVI29o6fzry7u4=
-github.com/stretchr/testify v1.8.2/go.mod h1:w2LPCIKwWwSfY2zedu0+kehJoqGctiVI29o6fzry7u4=
-github.com/stretchr/testify v1.8.4/go.mod h1:sz/lmYIOXD/1dqDmKjjqLyZ2RngseejIcXlSw2iwfAo=
 github.com/stretchr/testify v1.9.0 h1:HtqpIVDClZ4nwg75+f6Lvsy/wHu+3BoSGCbBAcpTsTg=
 github.com/stretchr/testify v1.9.0/go.mod h1:r2ic/lqez/lEtzL7wO/rwa5dbSLXVDPFyf8C91i36aY=
 github.com/stripe/stripe-go/v79 v79.12.0 h1:HQs/kxNEB3gYA7FnkSFkp0kSOeez0fsmCWev6SxftYs=
@@ -285,12 +288,28 @@ github.com/subosito/gotenv v1.6.0/go.mod h1:Dk4QP5c2W3ibzajGcXpNraDfq2IrhjMIvMSW
 github.com/tv42/httpunix v0.0.0-20150427012821-b75d8614f926/go.mod h1:9ESjWnEqriFuLhtthL60Sar/7RFoluCcXsuvEwTV5KM=
 github.com/twitchyliquid64/golang-asm v0.15.1 h1:SU5vSMR7hnwNxj24w34ZyCi/FmDZTkS4MhqMhdFk5YI=
 github.com/twitchyliquid64/golang-asm v0.15.1/go.mod h1:a1lVb/DtPvCB8fslRZhAngC2+aY1QWCk3Cedj/Gdt08=
-github.com/ugorji/go/codec v1.2.11 h1:BMaWp1Bb6fHwEtbplGBGJ498wD+LKlNSl25MjdZY4dU=
-github.com/ugorji/go/codec v1.2.11/go.mod h1:UNopzCgEMSXjBc6AOMqYvWC1ktqTAfzJZUZgYf6w6lg=
+github.com/ugorji/go/codec v1.2.12 h1:9LC83zGrHhuUA9l16C9AHXAqEV/2wBQ4nkvumAE65EE=
+github.com/ugorji/go/codec v1.2.12/go.mod h1:UNopzCgEMSXjBc6AOMqYvWC1ktqTAfzJZUZgYf6w6lg=
 github.com/x-cray/logrus-prefixed-formatter v0.5.2 h1:00txxvfBM9muc0jiLIEAkAcIMJzfthRT6usrui8uGmg=
 github.com/x-cray/logrus-prefixed-formatter v0.5.2/go.mod h1:2duySbKsL6M18s5GU7VPsoEPHyzalCE06qoARUCeBBE=
 github.com/yuin/goldmark v1.1.27/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
 github.com/yuin/goldmark v1.2.1/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
+go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0 h1:0nTRpaCaILLdooXAQnfktlL6Zw1ECKEW9DZGH2byi2c=
+go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0/go.mod h1:A7aFlp4WSLmeOnFRZwf2dMU+40THPc+rsr6KOwZLOcg=
+go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0 h1:4Pp6oUg3+e/6M4C0A/3kJ2VYa++dsWVTtGgLVj5xtHg=
+go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0/go.mod h1:Mjt1i1INqiaoZOMGR1RIUJN+i3ChKoFRqzrRQhlkbs0=
+go.opentelemetry.io/contrib/propagators/b3 v1.31.0 h1:PQPXYscmwbCp76QDvO4hMngF2j8Bx/OTV86laEl8uqo=
+go.opentelemetry.io/contrib/propagators/b3 v1.31.0/go.mod h1:jbqfV8wDdqSDrAYxVpXQnpM0XFMq2FtDesblJ7blOwQ=
+go.opentelemetry.io/otel v1.31.0 h1:NsJcKPIW0D0H3NgzPDHmo0WW6SptzPdqg/L1zsIm2hY=
+go.opentelemetry.io/otel v1.31.0/go.mod h1:O0C14Yl9FgkjqcCZAsE053C13OaddMYr/hz6clDkEJE=
+go.opentelemetry.io/otel/exporters/jaeger v1.17.0 h1:D7UpUy2Xc2wsi1Ras6V40q806WM07rqoCWzXu7Sqy+4=
+go.opentelemetry.io/otel/exporters/jaeger v1.17.0/go.mod h1:nPCqOnEH9rNLKqH/+rrUjiMzHJdV1BlpKcTwRTyKkKI=
+go.opentelemetry.io/otel/metric v1.31.0 h1:FSErL0ATQAmYHUIzSezZibnyVlft1ybhy4ozRPcF2fE=
+go.opentelemetry.io/otel/metric v1.31.0/go.mod h1:C3dEloVbLuYoX41KpmAhOqNriGbA+qqH6PQ5E5mUfnY=
+go.opentelemetry.io/otel/sdk v1.31.0 h1:xLY3abVHYZ5HSfOg3l2E5LUj2Cwva5Y7yGxnSW9H5Gk=
+go.opentelemetry.io/otel/sdk v1.31.0/go.mod h1:TfRbMdhvxIIr/B2N2LQW2S5v9m3gOQ/08KsbbO5BPT0=
+go.opentelemetry.io/otel/trace v1.31.0 h1:ffjsj1aRouKewfr85U2aGagJ46+MvodynlQ1HYdmJys=
+go.opentelemetry.io/otel/trace v1.31.0/go.mod h1:TXZkRk7SM2ZQLtR6eoAWQFIHPvzQ06FJAsO1tJg480A=
 go.uber.org/atomic v1.7.0/go.mod h1:fEN4uk6kAWBTFdckzkM89CLk9XfWZrxpCo0nPH17wJc=
 go.uber.org/atomic v1.9.0 h1:ECmE8Bn/WFTYwEW/bpKD3M8VtR/zQVbavAoalC1PYyE=
 go.uber.org/atomic v1.9.0/go.mod h1:fEN4uk6kAWBTFdckzkM89CLk9XfWZrxpCo0nPH17wJc=
@@ -301,16 +320,15 @@ go.uber.org/multierr v1.6.0/go.mod h1:cdWPpRnG4AhwMwsgIHip0KRBQjJy5kYEpYjJxpXp9i
 go.uber.org/multierr v1.9.0 h1:7fIwc/ZtS0q++VgcfqFDxSBZVv/Xo49/SYnDFupUwlI=
 go.uber.org/multierr v1.9.0/go.mod h1:X2jQV1h+kxSjClGpnseKVIxpmcjrj7MNnI0bnlfKTVQ=
 go.uber.org/zap v1.18.1/go.mod h1:xg/QME4nWcxGxrpdeYfq7UvYrLh66cuVKdrbD1XF/NI=
-golang.org/x/arch v0.0.0-20210923205945-b76863e36670/go.mod h1:5om86z9Hs0C8fWVUuoMHwpExlXzs5Tkyp9hOrfG7pp8=
-golang.org/x/arch v0.4.0 h1:A8WCeEWhLwPBKNbFi5Wv5UTCBx5zzubnXDlMOFAzFMc=
-golang.org/x/arch v0.4.0/go.mod h1:5om86z9Hs0C8fWVUuoMHwpExlXzs5Tkyp9hOrfG7pp8=
+golang.org/x/arch v0.11.0 h1:KXV8WWKCXm6tRpLirl2szsO5j/oOODwZf4hATmGVNs4=
+golang.org/x/arch v0.11.0/go.mod h1:FEVrYAQjsQXMVJ1nsMoVVXPZg6p2JE2mx8psSWTDQys=
 golang.org/x/crypto v0.0.0-20180904163835-0709b304e793/go.mod h1:6SG95UA2DQfeDnfUPMdvaQW0Q7yPrPDi9nlGo2tz2b4=
 golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2/go.mod h1:djNgcEr1/C05ACkg1iLfiJU5Ep61QUkGW8qpdssI0+w=
 golang.org/x/crypto v0.0.0-20190923035154-9ee001bba392/go.mod h1:/lpIB1dKB+9EgE3H3cr1v9wB50oz8l4C4h62xy7jSTY=
 golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550/go.mod h1:yigFU9vqHzYiE8UmvKecakEJjdnWj3jj499lnFckfCI=
 golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9/go.mod h1:LzIPMQfyMNhhGPhUkYOs5KpL4U8rLKemX1yGLhDgUto=
-golang.org/x/crypto v0.26.0 h1:RrRspgV4mU+YwB4FYnuBoKsUapNIL5cohGAmSH3azsw=
-golang.org/x/crypto v0.26.0/go.mod h1:GY7jblb9wI+FOo5y8/S2oY4zWP07AkOJ4+jxCqdqn54=
+golang.org/x/crypto v0.28.0 h1:GBDwsMXVQi34v5CCYUm2jkJvu4cbtru2U4TN2PSyQnw=
+golang.org/x/crypto v0.28.0/go.mod h1:rmgy+3RHxRZMyY0jjAJShp2zgEdOqj2AO7U0pYmeQ7U=
 golang.org/x/exp v0.0.0-20190121172915-509febef88a4/go.mod h1:CJ0aWSM057203Lf6IL+f9T1iT9GByDxfZKAQTCR3kQA=
 golang.org/x/exp v0.0.0-20230905200255-921286631fa9 h1:GoHiUyI/Tp2nVkLI2mCxVkOjsbSXD66ic0XW0js0R9g=
 golang.org/x/exp v0.0.0-20230905200255-921286631fa9/go.mod h1:S2oDrQGGwySpoQPVqRShND87VCbxmc6bL1Yd2oYrm6k=
@@ -334,8 +352,8 @@ golang.org/x/net v0.0.0-20201021035429-f5854403a974/go.mod h1:sp8m0HH+o8qH0wwXwY
 golang.org/x/net v0.0.0-20210226172049-e18ecbb05110/go.mod h1:m0MpNAwzfU5UDzcl9v0D8zg8gWTRqZa9RBIspLL5mdg=
 golang.org/x/net v0.0.0-20210410081132-afb366fc7cd1/go.mod h1:9tjilg8BloeKEkVJvy7fQ90B1CfIiPueXVOjqfkSzI8=
 golang.org/x/net v0.0.0-20210520170846-37e1c6afe023/go.mod h1:9nx3DQGgdP8bBQD5qxJ1jj9UTztislL4KSBs9R2vV5Y=
-golang.org/x/net v0.28.0 h1:a9JDOJc5GMUJ0+UDqmLT86WiEy7iWyIhz8gz8E4e5hE=
-golang.org/x/net v0.28.0/go.mod h1:yqtgsTWOOnlGLG9GFRrK3++bGOUEkNBoHZc8MEDWPNg=
+golang.org/x/net v0.30.0 h1:AcW1SDZMkb8IpzCdQUaIq2sP4sZ4zw+55h6ynffypl4=
+golang.org/x/net v0.30.0/go.mod h1:2wGyMJ5iFasEhkwi13ChkO/t1ECNC4X4eBKkVFyYFlU=
 golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be/go.mod h1:N/0e6XlmueqKjAGxoOufVs8QHGRruUQn6yWY3a++T0U=
 golang.org/x/sync v0.0.0-20180314180146-1d60e4601c6f/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sync v0.0.0-20181108010431-42b317875d0f/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
@@ -372,17 +390,17 @@ golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab/go.mod h1:oPkhp1MJrh7nUepCBc
 golang.org/x/sys v0.0.0-20220908164124-27713097b956/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.5.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.6.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
-golang.org/x/sys v0.24.0 h1:Twjiwq9dn6R1fQcyiK+wQyHWfaz/BJB+YIpzU/Cv3Xg=
-golang.org/x/sys v0.24.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
+golang.org/x/sys v0.26.0 h1:KHjCJyddX0LoSTb3J+vWpupP9p0oznkqVk/IfjymZbo=
+golang.org/x/sys v0.26.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
 golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1/go.mod h1:bj7SfCRtBDWHUb9snDiAeCFNEtKQo2Wmx5Cou7ajbmo=
-golang.org/x/term v0.23.0 h1:F6D4vR+EHoL9/sWAWgAR1H2DcHr4PareCbAaCo1RpuU=
-golang.org/x/term v0.23.0/go.mod h1:DgV24QBUrK6jhZXl+20l6UWznPlwAHm1Q1mGHtydmSk=
+golang.org/x/term v0.25.0 h1:WtHI/ltw4NvSUig5KARz9h521QvRC8RmF/cuYqifU24=
+golang.org/x/term v0.25.0/go.mod h1:RPyXicDX+6vLxogjjRxjgD2TKtmAO6NZBsBRfrOLu7M=
 golang.org/x/text v0.3.0/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
 golang.org/x/text v0.3.2/go.mod h1:bEr9sfX3Q8Zfm5fL9x+3itogRgK3+ptLWKqgva+5dAk=
 golang.org/x/text v0.3.3/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
 golang.org/x/text v0.3.6/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
-golang.org/x/text v0.17.0 h1:XtiM5bkSOt+ewxlOE/aE/AKEHibwj/6gvWMl9Rsh0Qc=
-golang.org/x/text v0.17.0/go.mod h1:BuEKDfySbSR4drPmRPG/7iBdf8hvFMuRexcpahXilzY=
+golang.org/x/text v0.19.0 h1:kTxAhCbGbxhK0IwgSKiMO5awPoDQ0RpfiVYBfK860YM=
+golang.org/x/text v0.19.0/go.mod h1:BuEKDfySbSR4drPmRPG/7iBdf8hvFMuRexcpahXilzY=
 golang.org/x/tools v0.0.0-20180917221912-90fa682c2a6e/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
 golang.org/x/tools v0.0.0-20190114222345-bf090417da8b/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
 golang.org/x/tools v0.0.0-20190226205152-f727befe758c/go.mod h1:9Yl7xja0Znq3iFh3HoIrodX9oNMXvdceNzlUR8zjMvY=
@@ -436,4 +454,3 @@ gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
 honnef.co/go/tools v0.0.0-20190102054323-c2f93a96b099/go.mod h1:rf3lG4BRIbNafJWhAfAdb/ePZxsR/4RtNHQocxwk9r4=
 honnef.co/go/tools v0.0.0-20190523083050-ea95bdfd59fc/go.mod h1:rf3lG4BRIbNafJWhAfAdb/ePZxsR/4RtNHQocxwk9r4=
 nullprogram.com/x/optparse v1.0.0/go.mod h1:KdyPE+Igbe0jQUrVfMqDMeJQIJZEuyV7pjYmp6pbG50=
-rsc.io/pdf v0.1.1/go.mod h1:n8OzWcQ6Sp37PL01nO98y4iUCRdTGarVfzxY20ICaU4=
diff --git a/internal/payment/main.go b/internal/payment/main.go
index b80b31c..d0e4165 100644
--- a/internal/payment/main.go
+++ b/internal/payment/main.go
@@ -7,6 +7,7 @@ import (
 	"github.com/ghost-yu/go_shop_second/common/config"
 	"github.com/ghost-yu/go_shop_second/common/logging"
 	"github.com/ghost-yu/go_shop_second/common/server"
+	"github.com/ghost-yu/go_shop_second/common/tracing"
 	"github.com/ghost-yu/go_shop_second/payment/infrastructure/consumer"
 	"github.com/ghost-yu/go_shop_second/payment/service"
 	"github.com/sirupsen/logrus"
@@ -21,11 +22,18 @@ func init() {
 }
 
 func main() {
+	serviceName := viper.GetString("payment.service-name")
 	ctx, cancel := context.WithCancel(context.Background())
 	defer cancel()
 
 	serverType := viper.GetString("payment.server-to-run")
 
+	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
+	if err != nil {
+		logrus.Fatal(err)
+	}
+	defer shutdown(ctx)
+
 	application, cleanup := service.NewApplication(ctx)
 	defer cleanup()
 
@@ -45,7 +53,7 @@ func main() {
 	paymentHandler := NewPaymentHandler(ch)
 	switch serverType {
 	case "http":
-		server.RunHTTPServer(viper.GetString("payment.service-name"), paymentHandler.RegisterRoutes)
+		server.RunHTTPServer(serviceName, paymentHandler.RegisterRoutes)
 	case "grpc":
 		logrus.Panic("unsupported server type: grpc")
 	default:
diff --git a/internal/stock/go.mod b/internal/stock/go.mod
index a648d88..4083f63 100644
--- a/internal/stock/go.mod
+++ b/internal/stock/go.mod
@@ -13,20 +13,23 @@ require (
 
 require (
 	github.com/armon/go-metrics v0.4.1 // indirect
-	github.com/bytedance/sonic v1.11.6 // indirect
-	github.com/bytedance/sonic/loader v0.1.1 // indirect
+	github.com/bytedance/sonic v1.12.3 // indirect
+	github.com/bytedance/sonic/loader v0.2.0 // indirect
 	github.com/cloudwego/base64x v0.1.4 // indirect
 	github.com/cloudwego/iasm v0.2.0 // indirect
 	github.com/fatih/color v1.14.1 // indirect
 	github.com/fsnotify/fsnotify v1.7.0 // indirect
-	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
+	github.com/gabriel-vasile/mimetype v1.4.5 // indirect
 	github.com/gin-contrib/sse v0.1.0 // indirect
 	github.com/gin-gonic/gin v1.10.0 // indirect
+	github.com/go-logr/logr v1.4.2 // indirect
+	github.com/go-logr/stdr v1.2.2 // indirect
 	github.com/go-playground/locales v0.14.1 // indirect
 	github.com/go-playground/universal-translator v0.18.1 // indirect
-	github.com/go-playground/validator/v10 v10.20.0 // indirect
-	github.com/goccy/go-json v0.10.2 // indirect
+	github.com/go-playground/validator/v10 v10.22.1 // indirect
+	github.com/goccy/go-json v0.10.3 // indirect
 	github.com/golang/protobuf v1.5.3 // indirect
+	github.com/google/uuid v1.6.0 // indirect
 	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
 	github.com/hashicorp/consul/api v1.28.2 // indirect
 	github.com/hashicorp/errwrap v1.1.0 // indirect
@@ -39,7 +42,7 @@ require (
 	github.com/hashicorp/hcl v1.0.0 // indirect
 	github.com/hashicorp/serf v0.10.1 // indirect
 	github.com/json-iterator/go v1.1.12 // indirect
-	github.com/klauspost/cpuid/v2 v2.2.7 // indirect
+	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
 	github.com/leodido/go-urn v1.4.0 // indirect
 	github.com/magiconair/properties v1.8.7 // indirect
 	github.com/mattn/go-colorable v0.1.13 // indirect
@@ -50,7 +53,7 @@ require (
 	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
 	github.com/modern-go/reflect2 v1.0.2 // indirect
 	github.com/nxadm/tail v1.4.11 // indirect
-	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
+	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
 	github.com/sagikazarmark/locafero v0.4.0 // indirect
 	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
 	github.com/sourcegraph/conc v0.3.0 // indirect
@@ -61,15 +64,23 @@ require (
 	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
 	github.com/ugorji/go/codec v1.2.12 // indirect
 	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
+	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0 // indirect
+	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0 // indirect
+	go.opentelemetry.io/contrib/propagators/b3 v1.31.0 // indirect
+	go.opentelemetry.io/otel v1.31.0 // indirect
+	go.opentelemetry.io/otel/exporters/jaeger v1.17.0 // indirect
+	go.opentelemetry.io/otel/metric v1.31.0 // indirect
+	go.opentelemetry.io/otel/sdk v1.31.0 // indirect
+	go.opentelemetry.io/otel/trace v1.31.0 // indirect
 	go.uber.org/atomic v1.9.0 // indirect
 	go.uber.org/multierr v1.9.0 // indirect
-	golang.org/x/arch v0.8.0 // indirect
-	golang.org/x/crypto v0.26.0 // indirect
+	golang.org/x/arch v0.11.0 // indirect
+	golang.org/x/crypto v0.28.0 // indirect
 	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
-	golang.org/x/net v0.28.0 // indirect
-	golang.org/x/sys v0.24.0 // indirect
-	golang.org/x/term v0.23.0 // indirect
-	golang.org/x/text v0.17.0 // indirect
+	golang.org/x/net v0.30.0 // indirect
+	golang.org/x/sys v0.26.0 // indirect
+	golang.org/x/term v0.25.0 // indirect
+	golang.org/x/text v0.19.0 // indirect
 	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
 	google.golang.org/protobuf v1.35.1 // indirect
 	gopkg.in/ini.v1 v1.67.0 // indirect
diff --git a/internal/stock/go.sum b/internal/stock/go.sum
index f93f793..f649128 100644
--- a/internal/stock/go.sum
+++ b/internal/stock/go.sum
@@ -16,10 +16,11 @@ github.com/beorn7/perks v0.0.0-20180321164747-3a771d992973/go.mod h1:Dwedo/Wpr24
 github.com/beorn7/perks v1.0.0/go.mod h1:KWe93zE9D1o94FZ5RNwFwVgaQK1VOXiVxmqh+CedLV8=
 github.com/beorn7/perks v1.0.1/go.mod h1:G2ZrVWU2WbWT9wwq4/hrbKbnv/1ERSJQ0ibhJ6rlkpw=
 github.com/bgentry/speakeasy v0.1.0/go.mod h1:+zsyZBPWlz7T6j88CTgSN5bM796AkVf0kBD4zp0CCIs=
-github.com/bytedance/sonic v1.11.6 h1:oUp34TzMlL+OY1OUWxHqsdkgC/Zfc85zGqw9siXjrc0=
-github.com/bytedance/sonic v1.11.6/go.mod h1:LysEHSvpvDySVdC2f87zGWf6CIKJcAvqab1ZaiQtds4=
-github.com/bytedance/sonic/loader v0.1.1 h1:c+e5Pt1k/cy5wMveRDyk2X4B9hF4g7an8N3zCYjJFNM=
+github.com/bytedance/sonic v1.12.3 h1:W2MGa7RCU1QTeYRTPE3+88mVC0yXmsRQRChiyVocVjU=
+github.com/bytedance/sonic v1.12.3/go.mod h1:B8Gt/XvtZ3Fqj+iSKMypzymZxw/FVwgIGKzMzT9r/rk=
 github.com/bytedance/sonic/loader v0.1.1/go.mod h1:ncP89zfokxS5LZrJxl5z0UJcsk4M4yY2JpfqGeCtNLU=
+github.com/bytedance/sonic/loader v0.2.0 h1:zNprn+lsIP06C/IqCHs3gPQIvnvpKbbxyXQP1iU4kWM=
+github.com/bytedance/sonic/loader v0.2.0/go.mod h1:ncP89zfokxS5LZrJxl5z0UJcsk4M4yY2JpfqGeCtNLU=
 github.com/census-instrumentation/opencensus-proto v0.2.1/go.mod h1:f6KPmirojxKA12rnyqOA5BBL4O983OfeGPqjHWSTneU=
 github.com/cespare/xxhash/v2 v2.1.1/go.mod h1:VGX0DQ3Q6kWi7AoAeZDth3/j3BFtOZR5XLFGgcrjCOs=
 github.com/circonus-labs/circonus-gometrics v2.3.1+incompatible/go.mod h1:nmEj6Dob7S7YxXgwXpfOuvO54S+tGdZdw9fuRZt25Ag=
@@ -48,8 +49,8 @@ github.com/frankban/quicktest v1.14.6/go.mod h1:4ptaffx2x8+WTWXmUCuVU6aPUX1/Mz7z
 github.com/fsnotify/fsnotify v1.6.0/go.mod h1:sl3t1tCWJFWoRz9R8WJCbQihKKwmorjAbSClcnxKAGw=
 github.com/fsnotify/fsnotify v1.7.0 h1:8JEhPFa5W2WU7YfeZzPNqzMP6Lwt7L2715Ggo0nosvA=
 github.com/fsnotify/fsnotify v1.7.0/go.mod h1:40Bi/Hjc2AVfZrqy+aj+yEI+/bRxZnMJyTJwOpGvigM=
-github.com/gabriel-vasile/mimetype v1.4.3 h1:in2uUcidCuFcDKtdcBxlR0rJ1+fsokWf+uqxgUFjbI0=
-github.com/gabriel-vasile/mimetype v1.4.3/go.mod h1:d8uq/6HKRL6CGdk+aubisF/M5GcPfT7nKyLpA0lbSSk=
+github.com/gabriel-vasile/mimetype v1.4.5 h1:J7wGKdGu33ocBOhGy0z653k/lFKLFDPJMG8Gql0kxn4=
+github.com/gabriel-vasile/mimetype v1.4.5/go.mod h1:ibHel+/kbxn9x2407k1izTA1S81ku1z/DlgOW2QE0M4=
 github.com/gin-contrib/sse v0.1.0 h1:Y/yl/+YNO8GZSjAhjMsSuLt29uWRFHdHYUb5lYOV9qE=
 github.com/gin-contrib/sse v0.1.0/go.mod h1:RHrZQHXnP2xjPF+u1gW/2HnVO7nvIa9PG3Gm+fLHvGI=
 github.com/gin-gonic/gin v1.10.0 h1:nTuyha1TYqgedzytsKYqna+DfLos46nTv2ygFy86HFU=
@@ -60,17 +61,22 @@ github.com/go-kit/log v0.1.0/go.mod h1:zbhenjAZHb184qTLMA9ZjW7ThYL0H2mk7Q6pNt4vb
 github.com/go-logfmt/logfmt v0.3.0/go.mod h1:Qt1PoO58o5twSAckw1HlFXLmHsOX5/0LbT9GBnD5lWE=
 github.com/go-logfmt/logfmt v0.4.0/go.mod h1:3RMwSq7FuexP4Kalkev3ejPJsZTpXXBr9+V4qmtdjCk=
 github.com/go-logfmt/logfmt v0.5.0/go.mod h1:wCYkCAKZfumFQihp8CzCvQ3paCTfi41vtzG1KdI/P7A=
+github.com/go-logr/logr v1.2.2/go.mod h1:jdQByPbusPIv2/zmleS9BjJVeZ6kBagPoEUsqbVz/1A=
+github.com/go-logr/logr v1.4.2 h1:6pFjapn8bFcIbiKo3XT4j/BhANplGihG6tvd+8rYgrY=
+github.com/go-logr/logr v1.4.2/go.mod h1:9T104GzyrTigFIr8wt5mBrctHMim0Nb2HLGrmQ40KvY=
+github.com/go-logr/stdr v1.2.2 h1:hSWxHoqTgW2S2qGc0LTAI563KZ5YKYRhT3MFKZMbjag=
+github.com/go-logr/stdr v1.2.2/go.mod h1:mMo/vtBO5dYbehREoey6XUKy/eSumjCCveDpRre4VKE=
 github.com/go-playground/assert/v2 v2.2.0 h1:JvknZsQTYeFEAhQwI4qEt9cyV5ONwRHC+lYKSsYSR8s=
 github.com/go-playground/assert/v2 v2.2.0/go.mod h1:VDjEfimB/XKnb+ZQfWdccd7VUvScMdVu0Titje2rxJ4=
 github.com/go-playground/locales v0.14.1 h1:EWaQ/wswjilfKLTECiXz7Rh+3BjFhfDFKv/oXslEjJA=
 github.com/go-playground/locales v0.14.1/go.mod h1:hxrqLVvrK65+Rwrd5Fc6F2O76J/NuW9t0sjnWqG1slY=
 github.com/go-playground/universal-translator v0.18.1 h1:Bcnm0ZwsGyWbCzImXv+pAJnYK9S473LQFuzCbDbfSFY=
 github.com/go-playground/universal-translator v0.18.1/go.mod h1:xekY+UJKNuX9WP91TpwSH2VMlDf28Uj24BCp08ZFTUY=
-github.com/go-playground/validator/v10 v10.20.0 h1:K9ISHbSaI0lyB2eWMPJo+kOS/FBExVwjEviJTixqxL8=
-github.com/go-playground/validator/v10 v10.20.0/go.mod h1:dbuPbCMFw/DrkbEynArYaCwl3amGuJotoKCe95atGMM=
+github.com/go-playground/validator/v10 v10.22.1 h1:40JcKH+bBNGFczGuoBYgX4I6m/i27HYW8P9FDk5PbgA=
+github.com/go-playground/validator/v10 v10.22.1/go.mod h1:dbuPbCMFw/DrkbEynArYaCwl3amGuJotoKCe95atGMM=
 github.com/go-stack/stack v1.8.0/go.mod h1:v0f6uXyyMGvRgIKkXu+yp6POWl0qKG85gN/melR3HDY=
-github.com/goccy/go-json v0.10.2 h1:CrxCmQqYDkv1z7lO7Wbh2HN93uovUHgrECaO5ZrCXAU=
-github.com/goccy/go-json v0.10.2/go.mod h1:6MelG93GURQebXPDq3khkgXZkazVtN9CRI+MGFi0w8I=
+github.com/goccy/go-json v0.10.3 h1:KZ5WoDbxAIgm2HNbYckL0se1fHD6rz5j4ywS6ebzDqA=
+github.com/goccy/go-json v0.10.3/go.mod h1:oq7eo15ShAhp70Anwd5lgX2pLfOS3QCiwU/PULtXL6M=
 github.com/gogo/protobuf v1.1.1/go.mod h1:r8qH/GZQm5c6nD/R0oafs1akxWv10x8SbQlK7atdtwQ=
 github.com/gogo/protobuf v1.3.2 h1:Ov1cvc58UF3b5XjBnZv7+opcTcQFZebYjWzi34vdm4Q=
 github.com/gogo/protobuf v1.3.2/go.mod h1:P1XiOD3dCwIKUDQYPy72D8LYyHL2YPYrpS2s69NZV8Q=
@@ -93,6 +99,8 @@ github.com/google/go-cmp v0.5.5/go.mod h1:v8dTdLbMG2kIc/vJvl+f65V22dbkXbowE6jgT/
 github.com/google/go-cmp v0.6.0 h1:ofyhxvXcZhMsU5ulbFiLKl/XBFqE1GSq7atu8tAmTRI=
 github.com/google/go-cmp v0.6.0/go.mod h1:17dUlkBOakJ0+DkrSSNjCkIjxS6bF9zb3elmeNGIjoY=
 github.com/google/gofuzz v1.0.0/go.mod h1:dBl0BpW6vV/+mYPU4Po3pmUjxk6FQPldtuIdl/M65Eg=
+github.com/google/uuid v1.6.0 h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=
+github.com/google/uuid v1.6.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
 github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 h1:UH//fgunKIs4JdUbpDl1VZCDaL56wXCB/5+wF6uHfaI=
 github.com/grpc-ecosystem/go-grpc-middleware v1.4.0/go.mod h1:g5qyo/la0ALbONm6Vbp88Yd8NsDy6rZz+RcrMPxvld8=
 github.com/hashicorp/consul/api v1.28.2 h1:mXfkRHrpHN4YY3RqL09nXU1eHKLNiuAN4kHvDQ16k/8=
@@ -149,8 +157,8 @@ github.com/julienschmidt/httprouter v1.2.0/go.mod h1:SYymIcj16QtmaHHD7aYtjjsJG7V
 github.com/kisielk/errcheck v1.5.0/go.mod h1:pFxgyoBC7bSaBwPgfKdkLd5X25qrDl4LWUI2bnpBCr8=
 github.com/kisielk/gotool v1.0.0/go.mod h1:XhKaO+MFFWcvkIS/tQcRk01m1F5IRFswLeQ+oQHNcck=
 github.com/klauspost/cpuid/v2 v2.0.9/go.mod h1:FInQzS24/EEf25PyTYn52gqo7WaD8xa0213Md/qVLRg=
-github.com/klauspost/cpuid/v2 v2.2.7 h1:ZWSB3igEs+d0qvnxR/ZBzXVmxkgt8DdzP6m9pfuVLDM=
-github.com/klauspost/cpuid/v2 v2.2.7/go.mod h1:Lcz8mBdAVJIBVzewtcLocK12l3Y+JytZYpaMropDUws=
+github.com/klauspost/cpuid/v2 v2.2.8 h1:+StwCXwm9PdpiEkPyzBXIy+M9KUb4ODm0Zarf1kS5BM=
+github.com/klauspost/cpuid/v2 v2.2.8/go.mod h1:Lcz8mBdAVJIBVzewtcLocK12l3Y+JytZYpaMropDUws=
 github.com/knz/go-libedit v1.10.1/go.mod h1:MZTVkCWyz0oBc7JOWP3wNAzd002ZbM/5hgShxwh4x8M=
 github.com/konsorten/go-windows-terminal-sequences v1.0.1/go.mod h1:T0+1ngSBFLxvqU3pZ+m/2kptfBszLMUkC4ZK/EgS/cQ=
 github.com/kr/logfmt v0.0.0-20140226030751-b84e30acd515/go.mod h1:+0opPa2QZZtGFBFZlji/RkVcI2GknAs/DXo4wKdlNEc=
@@ -210,8 +218,8 @@ github.com/opentracing/opentracing-go v1.1.0/go.mod h1:UkNAQd3GIcIGf0SeVgPpRdFSt
 github.com/pascaldekloe/goe v0.0.0-20180627143212-57f6aae5913c/go.mod h1:lzWF7FIEvWOWxwDKqyGYQf6ZUaNfKdP144TG7ZOy1lc=
 github.com/pascaldekloe/goe v0.1.0 h1:cBOtyMzM9HTpWjXfbbunk26uA6nG3a8n06Wieeh0MwY=
 github.com/pascaldekloe/goe v0.1.0/go.mod h1:lzWF7FIEvWOWxwDKqyGYQf6ZUaNfKdP144TG7ZOy1lc=
-github.com/pelletier/go-toml/v2 v2.2.2 h1:aYUidT7k73Pcl9nb2gScu7NSrKCSHIDE89b3+6Wq+LM=
-github.com/pelletier/go-toml/v2 v2.2.2/go.mod h1:1t835xjRzz80PqgE6HHgN2JOsmgYu/h4qDAS4n929Rs=
+github.com/pelletier/go-toml/v2 v2.2.3 h1:YmeHyLY8mFWbdkNWwpr+qIL2bEqT0o95WSdkNHvL12M=
+github.com/pelletier/go-toml/v2 v2.2.3/go.mod h1:MfCQTFTvCcUyyvvwm1+G6H/jORL20Xlb6rzQu9GuUkc=
 github.com/pkg/errors v0.8.0/go.mod h1:bwawxfHBFNV+L2hUp1rHADufV3IMtnDRdf1r5NINEl0=
 github.com/pkg/errors v0.8.1/go.mod h1:bwawxfHBFNV+L2hUp1rHADufV3IMtnDRdf1r5NINEl0=
 github.com/pkg/errors v0.9.1 h1:FEBLx1zS214owpjy7qsBeixbURkuhQAwrK5UwLGTwt4=
@@ -259,9 +267,8 @@ github.com/spf13/viper v1.19.0/go.mod h1:GQUN9bilAbhU/jgc1bKs99f/suXKeUMct8Adx5+
 github.com/stretchr/objx v0.1.0/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
 github.com/stretchr/objx v0.1.1/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
 github.com/stretchr/objx v0.4.0/go.mod h1:YvHI0jy2hoMjB+UWwv71VJQ9isScKT/TqJzVSSt89Yw=
+github.com/stretchr/objx v0.5.0 h1:1zr/of2m5FGMsad5YfcqgdqdWrIhu+EBEJRhR1U7z/c=
 github.com/stretchr/objx v0.5.0/go.mod h1:Yh+to48EsGEfYuaHDzXPcE3xhTkx73EhmCGUpEOglKo=
-github.com/stretchr/objx v0.5.2 h1:xuMeJ0Sdp5ZMRXx/aWO6RZxdr3beISkG5/G/aIRr3pY=
-github.com/stretchr/objx v0.5.2/go.mod h1:FRsXN1f5AsAjCGJKqEizvkpNtU+EGNCLh3NxZ/8L+MA=
 github.com/stretchr/testify v1.2.2/go.mod h1:a8OnRcib4nhh0OaRAV+Yts87kKdq0PP7pXfy6kDkUVs=
 github.com/stretchr/testify v1.3.0/go.mod h1:M5WIy9Dh21IEIfnGCwXGc5bZfKNJtfHm1UVUgZn+9EI=
 github.com/stretchr/testify v1.4.0/go.mod h1:j7eGeouHqKxXV5pUuKE4zz7dFj8WfuZ+81PSLYec5m4=
@@ -270,7 +277,6 @@ github.com/stretchr/testify v1.7.1/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/
 github.com/stretchr/testify v1.7.2/go.mod h1:R6va5+xMeoiuVRoj+gSkQ7d3FALtqAAGI1FQKckRals=
 github.com/stretchr/testify v1.8.0/go.mod h1:yNjHg4UonilssWZ8iaSj1OCr/vHnekPRkoO+kdMU+MU=
 github.com/stretchr/testify v1.8.1/go.mod h1:w2LPCIKwWwSfY2zedu0+kehJoqGctiVI29o6fzry7u4=
-github.com/stretchr/testify v1.8.4/go.mod h1:sz/lmYIOXD/1dqDmKjjqLyZ2RngseejIcXlSw2iwfAo=
 github.com/stretchr/testify v1.9.0 h1:HtqpIVDClZ4nwg75+f6Lvsy/wHu+3BoSGCbBAcpTsTg=
 github.com/stretchr/testify v1.9.0/go.mod h1:r2ic/lqez/lEtzL7wO/rwa5dbSLXVDPFyf8C91i36aY=
 github.com/subosito/gotenv v1.6.0 h1:9NlTDc1FTs4qu0DDq7AEtTPNw6SVm7uBMsUCUjABIf8=
@@ -284,6 +290,22 @@ github.com/x-cray/logrus-prefixed-formatter v0.5.2 h1:00txxvfBM9muc0jiLIEAkAcIMJ
 github.com/x-cray/logrus-prefixed-formatter v0.5.2/go.mod h1:2duySbKsL6M18s5GU7VPsoEPHyzalCE06qoARUCeBBE=
 github.com/yuin/goldmark v1.1.27/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
 github.com/yuin/goldmark v1.2.1/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
+go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0 h1:0nTRpaCaILLdooXAQnfktlL6Zw1ECKEW9DZGH2byi2c=
+go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0/go.mod h1:A7aFlp4WSLmeOnFRZwf2dMU+40THPc+rsr6KOwZLOcg=
+go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0 h1:4Pp6oUg3+e/6M4C0A/3kJ2VYa++dsWVTtGgLVj5xtHg=
+go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0/go.mod h1:Mjt1i1INqiaoZOMGR1RIUJN+i3ChKoFRqzrRQhlkbs0=
+go.opentelemetry.io/contrib/propagators/b3 v1.31.0 h1:PQPXYscmwbCp76QDvO4hMngF2j8Bx/OTV86laEl8uqo=
+go.opentelemetry.io/contrib/propagators/b3 v1.31.0/go.mod h1:jbqfV8wDdqSDrAYxVpXQnpM0XFMq2FtDesblJ7blOwQ=
+go.opentelemetry.io/otel v1.31.0 h1:NsJcKPIW0D0H3NgzPDHmo0WW6SptzPdqg/L1zsIm2hY=
+go.opentelemetry.io/otel v1.31.0/go.mod h1:O0C14Yl9FgkjqcCZAsE053C13OaddMYr/hz6clDkEJE=
+go.opentelemetry.io/otel/exporters/jaeger v1.17.0 h1:D7UpUy2Xc2wsi1Ras6V40q806WM07rqoCWzXu7Sqy+4=
+go.opentelemetry.io/otel/exporters/jaeger v1.17.0/go.mod h1:nPCqOnEH9rNLKqH/+rrUjiMzHJdV1BlpKcTwRTyKkKI=
+go.opentelemetry.io/otel/metric v1.31.0 h1:FSErL0ATQAmYHUIzSezZibnyVlft1ybhy4ozRPcF2fE=
+go.opentelemetry.io/otel/metric v1.31.0/go.mod h1:C3dEloVbLuYoX41KpmAhOqNriGbA+qqH6PQ5E5mUfnY=
+go.opentelemetry.io/otel/sdk v1.31.0 h1:xLY3abVHYZ5HSfOg3l2E5LUj2Cwva5Y7yGxnSW9H5Gk=
+go.opentelemetry.io/otel/sdk v1.31.0/go.mod h1:TfRbMdhvxIIr/B2N2LQW2S5v9m3gOQ/08KsbbO5BPT0=
+go.opentelemetry.io/otel/trace v1.31.0 h1:ffjsj1aRouKewfr85U2aGagJ46+MvodynlQ1HYdmJys=
+go.opentelemetry.io/otel/trace v1.31.0/go.mod h1:TXZkRk7SM2ZQLtR6eoAWQFIHPvzQ06FJAsO1tJg480A=
 go.uber.org/atomic v1.7.0/go.mod h1:fEN4uk6kAWBTFdckzkM89CLk9XfWZrxpCo0nPH17wJc=
 go.uber.org/atomic v1.9.0 h1:ECmE8Bn/WFTYwEW/bpKD3M8VtR/zQVbavAoalC1PYyE=
 go.uber.org/atomic v1.9.0/go.mod h1:fEN4uk6kAWBTFdckzkM89CLk9XfWZrxpCo0nPH17wJc=
@@ -292,16 +314,15 @@ go.uber.org/multierr v1.6.0/go.mod h1:cdWPpRnG4AhwMwsgIHip0KRBQjJy5kYEpYjJxpXp9i
 go.uber.org/multierr v1.9.0 h1:7fIwc/ZtS0q++VgcfqFDxSBZVv/Xo49/SYnDFupUwlI=
 go.uber.org/multierr v1.9.0/go.mod h1:X2jQV1h+kxSjClGpnseKVIxpmcjrj7MNnI0bnlfKTVQ=
 go.uber.org/zap v1.18.1/go.mod h1:xg/QME4nWcxGxrpdeYfq7UvYrLh66cuVKdrbD1XF/NI=
-golang.org/x/arch v0.0.0-20210923205945-b76863e36670/go.mod h1:5om86z9Hs0C8fWVUuoMHwpExlXzs5Tkyp9hOrfG7pp8=
-golang.org/x/arch v0.8.0 h1:3wRIsP3pM4yUptoR96otTUOXI367OS0+c9eeRi9doIc=
-golang.org/x/arch v0.8.0/go.mod h1:FEVrYAQjsQXMVJ1nsMoVVXPZg6p2JE2mx8psSWTDQys=
+golang.org/x/arch v0.11.0 h1:KXV8WWKCXm6tRpLirl2szsO5j/oOODwZf4hATmGVNs4=
+golang.org/x/arch v0.11.0/go.mod h1:FEVrYAQjsQXMVJ1nsMoVVXPZg6p2JE2mx8psSWTDQys=
 golang.org/x/crypto v0.0.0-20180904163835-0709b304e793/go.mod h1:6SG95UA2DQfeDnfUPMdvaQW0Q7yPrPDi9nlGo2tz2b4=
 golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2/go.mod h1:djNgcEr1/C05ACkg1iLfiJU5Ep61QUkGW8qpdssI0+w=
 golang.org/x/crypto v0.0.0-20190923035154-9ee001bba392/go.mod h1:/lpIB1dKB+9EgE3H3cr1v9wB50oz8l4C4h62xy7jSTY=
 golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550/go.mod h1:yigFU9vqHzYiE8UmvKecakEJjdnWj3jj499lnFckfCI=
 golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9/go.mod h1:LzIPMQfyMNhhGPhUkYOs5KpL4U8rLKemX1yGLhDgUto=
-golang.org/x/crypto v0.26.0 h1:RrRspgV4mU+YwB4FYnuBoKsUapNIL5cohGAmSH3azsw=
-golang.org/x/crypto v0.26.0/go.mod h1:GY7jblb9wI+FOo5y8/S2oY4zWP07AkOJ4+jxCqdqn54=
+golang.org/x/crypto v0.28.0 h1:GBDwsMXVQi34v5CCYUm2jkJvu4cbtru2U4TN2PSyQnw=
+golang.org/x/crypto v0.28.0/go.mod h1:rmgy+3RHxRZMyY0jjAJShp2zgEdOqj2AO7U0pYmeQ7U=
 golang.org/x/exp v0.0.0-20190121172915-509febef88a4/go.mod h1:CJ0aWSM057203Lf6IL+f9T1iT9GByDxfZKAQTCR3kQA=
 golang.org/x/exp v0.0.0-20230905200255-921286631fa9 h1:GoHiUyI/Tp2nVkLI2mCxVkOjsbSXD66ic0XW0js0R9g=
 golang.org/x/exp v0.0.0-20230905200255-921286631fa9/go.mod h1:S2oDrQGGwySpoQPVqRShND87VCbxmc6bL1Yd2oYrm6k=
@@ -324,8 +345,8 @@ golang.org/x/net v0.0.0-20200226121028-0de0cce0169b/go.mod h1:z5CRVTTTmAJ677TzLL
 golang.org/x/net v0.0.0-20201021035429-f5854403a974/go.mod h1:sp8m0HH+o8qH0wwXwYZr8TS3Oi6o0r6Gce1SSxlDquU=
 golang.org/x/net v0.0.0-20210226172049-e18ecbb05110/go.mod h1:m0MpNAwzfU5UDzcl9v0D8zg8gWTRqZa9RBIspLL5mdg=
 golang.org/x/net v0.0.0-20210410081132-afb366fc7cd1/go.mod h1:9tjilg8BloeKEkVJvy7fQ90B1CfIiPueXVOjqfkSzI8=
-golang.org/x/net v0.28.0 h1:a9JDOJc5GMUJ0+UDqmLT86WiEy7iWyIhz8gz8E4e5hE=
-golang.org/x/net v0.28.0/go.mod h1:yqtgsTWOOnlGLG9GFRrK3++bGOUEkNBoHZc8MEDWPNg=
+golang.org/x/net v0.30.0 h1:AcW1SDZMkb8IpzCdQUaIq2sP4sZ4zw+55h6ynffypl4=
+golang.org/x/net v0.30.0/go.mod h1:2wGyMJ5iFasEhkwi13ChkO/t1ECNC4X4eBKkVFyYFlU=
 golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be/go.mod h1:N/0e6XlmueqKjAGxoOufVs8QHGRruUQn6yWY3a++T0U=
 golang.org/x/sync v0.0.0-20180314180146-1d60e4601c6f/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sync v0.0.0-20181108010431-42b317875d0f/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
@@ -361,17 +382,17 @@ golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab/go.mod h1:oPkhp1MJrh7nUepCBc
 golang.org/x/sys v0.0.0-20220908164124-27713097b956/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.5.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.6.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
-golang.org/x/sys v0.24.0 h1:Twjiwq9dn6R1fQcyiK+wQyHWfaz/BJB+YIpzU/Cv3Xg=
-golang.org/x/sys v0.24.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
+golang.org/x/sys v0.26.0 h1:KHjCJyddX0LoSTb3J+vWpupP9p0oznkqVk/IfjymZbo=
+golang.org/x/sys v0.26.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
 golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1/go.mod h1:bj7SfCRtBDWHUb9snDiAeCFNEtKQo2Wmx5Cou7ajbmo=
-golang.org/x/term v0.23.0 h1:F6D4vR+EHoL9/sWAWgAR1H2DcHr4PareCbAaCo1RpuU=
-golang.org/x/term v0.23.0/go.mod h1:DgV24QBUrK6jhZXl+20l6UWznPlwAHm1Q1mGHtydmSk=
+golang.org/x/term v0.25.0 h1:WtHI/ltw4NvSUig5KARz9h521QvRC8RmF/cuYqifU24=
+golang.org/x/term v0.25.0/go.mod h1:RPyXicDX+6vLxogjjRxjgD2TKtmAO6NZBsBRfrOLu7M=
 golang.org/x/text v0.3.0/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
 golang.org/x/text v0.3.2/go.mod h1:bEr9sfX3Q8Zfm5fL9x+3itogRgK3+ptLWKqgva+5dAk=
 golang.org/x/text v0.3.3/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
 golang.org/x/text v0.3.6/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
-golang.org/x/text v0.17.0 h1:XtiM5bkSOt+ewxlOE/aE/AKEHibwj/6gvWMl9Rsh0Qc=
-golang.org/x/text v0.17.0/go.mod h1:BuEKDfySbSR4drPmRPG/7iBdf8hvFMuRexcpahXilzY=
+golang.org/x/text v0.19.0 h1:kTxAhCbGbxhK0IwgSKiMO5awPoDQ0RpfiVYBfK860YM=
+golang.org/x/text v0.19.0/go.mod h1:BuEKDfySbSR4drPmRPG/7iBdf8hvFMuRexcpahXilzY=
 golang.org/x/tools v0.0.0-20180917221912-90fa682c2a6e/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
 golang.org/x/tools v0.0.0-20190114222345-bf090417da8b/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
 golang.org/x/tools v0.0.0-20190226205152-f727befe758c/go.mod h1:9Yl7xja0Znq3iFh3HoIrodX9oNMXvdceNzlUR8zjMvY=
@@ -425,4 +446,3 @@ gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
 honnef.co/go/tools v0.0.0-20190102054323-c2f93a96b099/go.mod h1:rf3lG4BRIbNafJWhAfAdb/ePZxsR/4RtNHQocxwk9r4=
 honnef.co/go/tools v0.0.0-20190523083050-ea95bdfd59fc/go.mod h1:rf3lG4BRIbNafJWhAfAdb/ePZxsR/4RtNHQocxwk9r4=
 nullprogram.com/x/optparse v1.0.0/go.mod h1:KdyPE+Igbe0jQUrVfMqDMeJQIJZEuyV7pjYmp6pbG50=
-rsc.io/pdf v0.1.1/go.mod h1:n8OzWcQ6Sp37PL01nO98y4iUCRdTGarVfzxY20ICaU4=
diff --git a/internal/stock/main.go b/internal/stock/main.go
index ffd9970..548e6f1 100644
--- a/internal/stock/main.go
+++ b/internal/stock/main.go
@@ -8,6 +8,7 @@ import (
 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
 	"github.com/ghost-yu/go_shop_second/common/logging"
 	"github.com/ghost-yu/go_shop_second/common/server"
+	"github.com/ghost-yu/go_shop_second/common/tracing"
 	"github.com/ghost-yu/go_shop_second/stock/ports"
 	"github.com/ghost-yu/go_shop_second/stock/service"
 	"github.com/sirupsen/logrus"
@@ -26,11 +27,15 @@ func main() {
 	serviceName := viper.GetString("stock.service-name")
 	serverType := viper.GetString("stock.server-to-run")
 
-	logrus.Info(serverType)
-
 	ctx, cancel := context.WithCancel(context.Background())
 	defer cancel()
 
+	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
+	if err != nil {
+		logrus.Fatal(err)
+	}
+	defer shutdown(ctx)
+
 	application := service.NewApplication(ctx)
 
 	deregisterFunc, err := discovery.RegisterToConsul(ctx, serviceName)
~~~
