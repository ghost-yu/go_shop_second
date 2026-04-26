# Lesson Pair Diff Report

- FromBranch: lesson59
- ToBranch: lesson60

## Short Summary

~~~text
 29 files changed, 507 insertions(+), 234 deletions(-)
~~~

## File Stats

~~~text
 internal/common/go.mod                             | 15 +++--
 internal/common/go.sum                             | 68 ++++++++++++++++++----
 internal/common/handler/redis/client.go            | 27 ++++-----
 internal/common/logging/grpc.go                    | 29 +++++++++
 internal/common/logging/logrus.go                  | 56 +++++++++++++++++-
 internal/common/logging/mysql.go                   | 12 ++--
 internal/common/logging/when.go                    | 59 +++++++++++++++++++
 internal/common/middleware/request.go              | 17 +++---
 internal/common/server/gprc.go                     |  2 +
 .../kitchen/infrastructure/consumer/consumer.go    | 34 ++++++-----
 internal/order/adapters/grpc/stock_grpc.go         | 16 +++--
 internal/order/adapters/order_inmem_repository.go  | 52 ++++++++++-------
 internal/order/adapters/order_mongo_repository.go  | 44 ++++++--------
 internal/order/app/command/create_order.go         |  4 ++
 internal/order/app/command/update_order.go         |  9 ++-
 internal/order/go.mod                              | 16 +++--
 internal/order/go.sum                              | 40 +++++++++----
 internal/order/infrastructure/consumer/consumer.go | 18 +++---
 internal/order/ports/grpc.go                       |  2 -
 internal/payment/adapters/order_grpc.go            |  7 ---
 internal/payment/app/command/create_payment.go     |  7 ++-
 internal/payment/go.mod                            | 17 +++---
 internal/payment/go.sum                            | 38 ++++++++----
 internal/payment/http.go                           | 38 +++++++-----
 .../payment/infrastructure/consumer/consumer.go    | 20 ++++---
 .../stock/app/query/check_if_items_in_stock.go     |  3 +-
 internal/stock/go.mod                              | 14 +++--
 internal/stock/go.sum                              | 36 ++++++++----
 internal/stock/infrastructure/persistent/mysql.go  | 41 ++++++-------
 29 files changed, 507 insertions(+), 234 deletions(-)
~~~

## Commit Comparison

~~~text
> 74d970e 全链路可观测性建设“
> 02c1271 Merge pull request #1 from Nicknamezz00/lesson59
~~~

## Changed Files

~~~text
internal/common/go.mod
internal/common/go.sum
internal/common/handler/redis/client.go
internal/common/logging/grpc.go
internal/common/logging/logrus.go
internal/common/logging/mysql.go
internal/common/logging/when.go
internal/common/middleware/request.go
internal/common/server/gprc.go
internal/kitchen/infrastructure/consumer/consumer.go
internal/order/adapters/grpc/stock_grpc.go
internal/order/adapters/order_inmem_repository.go
internal/order/adapters/order_mongo_repository.go
internal/order/app/command/create_order.go
internal/order/app/command/update_order.go
internal/order/go.mod
internal/order/go.sum
internal/order/infrastructure/consumer/consumer.go
internal/order/ports/grpc.go
internal/payment/adapters/order_grpc.go
internal/payment/app/command/create_payment.go
internal/payment/go.mod
internal/payment/go.sum
internal/payment/http.go
internal/payment/infrastructure/consumer/consumer.go
internal/stock/app/query/check_if_items_in_stock.go
internal/stock/go.mod
internal/stock/go.sum
internal/stock/infrastructure/persistent/mysql.go
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
internal/common/handler/redis/client.go
internal/common/logging/grpc.go
internal/common/logging/logrus.go
internal/common/logging/mysql.go
internal/common/logging/when.go
internal/common/middleware/request.go
internal/common/server/gprc.go
internal/kitchen/infrastructure/consumer/consumer.go
internal/order/adapters/grpc/stock_grpc.go
internal/order/adapters/order_inmem_repository.go
internal/order/adapters/order_mongo_repository.go
internal/order/app/command/create_order.go
internal/order/app/command/update_order.go
internal/order/infrastructure/consumer/consumer.go
internal/order/ports/grpc.go
internal/payment/adapters/order_grpc.go
internal/payment/app/command/create_payment.go
internal/payment/http.go
internal/payment/infrastructure/consumer/consumer.go
internal/stock/app/query/check_if_items_in_stock.go
internal/stock/infrastructure/persistent/mysql.go
~~~

## Full Diff

~~~diff
diff --git a/internal/common/go.mod b/internal/common/go.mod
index f7ef6a8..6ec52f4 100644
--- a/internal/common/go.mod
+++ b/internal/common/go.mod
@@ -11,6 +11,7 @@ require (
 	github.com/redis/go-redis/v9 v9.7.0
 	github.com/sirupsen/logrus v1.8.1
 	github.com/spf13/viper v1.19.0
+	github.com/x-cray/logrus-prefixed-formatter v0.5.2
 	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0
 	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0
 	go.opentelemetry.io/contrib/propagators/b3 v1.31.0
@@ -19,7 +20,7 @@ require (
 	go.opentelemetry.io/otel/sdk v1.31.0
 	go.opentelemetry.io/otel/trace v1.31.0
 	google.golang.org/grpc v1.67.1
-	google.golang.org/protobuf v1.35.1
+	google.golang.org/protobuf v1.36.1
 )
 
 require (
@@ -58,10 +59,13 @@ require (
 	github.com/magiconair/properties v1.8.7 // indirect
 	github.com/mattn/go-colorable v0.1.13 // indirect
 	github.com/mattn/go-isatty v0.0.20 // indirect
+	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
 	github.com/mitchellh/go-homedir v1.1.0 // indirect
 	github.com/mitchellh/mapstructure v1.5.0 // indirect
 	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
 	github.com/modern-go/reflect2 v1.0.2 // indirect
+	github.com/onsi/ginkgo v1.16.5 // indirect
+	github.com/onsi/gomega v1.36.2 // indirect
 	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
 	github.com/sagikazarmark/locafero v0.4.0 // indirect
 	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
@@ -76,11 +80,12 @@ require (
 	go.uber.org/atomic v1.9.0 // indirect
 	go.uber.org/multierr v1.9.0 // indirect
 	golang.org/x/arch v0.11.0 // indirect
-	golang.org/x/crypto v0.28.0 // indirect
+	golang.org/x/crypto v0.31.0 // indirect
 	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
-	golang.org/x/net v0.30.0 // indirect
-	golang.org/x/sys v0.26.0 // indirect
-	golang.org/x/text v0.19.0 // indirect
+	golang.org/x/net v0.33.0 // indirect
+	golang.org/x/sys v0.28.0 // indirect
+	golang.org/x/term v0.27.0 // indirect
+	golang.org/x/text v0.21.0 // indirect
 	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
 	gopkg.in/ini.v1 v1.67.0 // indirect
 	gopkg.in/yaml.v3 v3.0.1 // indirect
diff --git a/internal/common/go.sum b/internal/common/go.sum
index 553c9f1..59fdfd2 100644
--- a/internal/common/go.sum
+++ b/internal/common/go.sum
@@ -58,6 +58,8 @@ github.com/fatih/color v1.14.1 h1:qfhVLaG5s+nCROl1zJsZRxFeYrHLqWroPOQ8BWiNb4w=
 github.com/fatih/color v1.14.1/go.mod h1:2oHN61fhTpgcxD3TSWCgKDiH1+x4OiDVVGH8WlgGZGg=
 github.com/frankban/quicktest v1.14.6 h1:7Xjx+VpznH+oBnejlPUj8oUpdxnVs4f8XU8WnHkI4W8=
 github.com/frankban/quicktest v1.14.6/go.mod h1:4ptaffx2x8+WTWXmUCuVU6aPUX1/Mz7zb5vbUoiM6w0=
+github.com/fsnotify/fsnotify v1.4.7/go.mod h1:jwhsz4b93w/PPRr/qN1Yymfu8t87LnFCMoQvtojpjFo=
+github.com/fsnotify/fsnotify v1.4.9/go.mod h1:znqG4EE+3YCdAaPaxE2ZRY/06pZUdp0tY4IgpuI1SZQ=
 github.com/fsnotify/fsnotify v1.7.0 h1:8JEhPFa5W2WU7YfeZzPNqzMP6Lwt7L2715Ggo0nosvA=
 github.com/fsnotify/fsnotify v1.7.0/go.mod h1:40Bi/Hjc2AVfZrqy+aj+yEI+/bRxZnMJyTJwOpGvigM=
 github.com/gabriel-vasile/mimetype v1.4.5 h1:J7wGKdGu33ocBOhGy0z653k/lFKLFDPJMG8Gql0kxn4=
@@ -86,6 +88,7 @@ github.com/go-playground/universal-translator v0.18.1/go.mod h1:xekY+UJKNuX9WP91
 github.com/go-playground/validator/v10 v10.22.1 h1:40JcKH+bBNGFczGuoBYgX4I6m/i27HYW8P9FDk5PbgA=
 github.com/go-playground/validator/v10 v10.22.1/go.mod h1:dbuPbCMFw/DrkbEynArYaCwl3amGuJotoKCe95atGMM=
 github.com/go-stack/stack v1.8.0/go.mod h1:v0f6uXyyMGvRgIKkXu+yp6POWl0qKG85gN/melR3HDY=
+github.com/go-task/slim-sprig v0.0.0-20210107165309-348f09dbbbc0/go.mod h1:fyg7847qk6SyHyPtNmDHnmrv/HOrqktSC+C9fM+CJOE=
 github.com/goccy/go-json v0.10.3 h1:KZ5WoDbxAIgm2HNbYckL0se1fHD6rz5j4ywS6ebzDqA=
 github.com/goccy/go-json v0.10.3/go.mod h1:oq7eo15ShAhp70Anwd5lgX2pLfOS3QCiwU/PULtXL6M=
 github.com/gogo/protobuf v1.1.1/go.mod h1:r8qH/GZQm5c6nD/R0oafs1akxWv10x8SbQlK7atdtwQ=
@@ -97,6 +100,12 @@ github.com/golang/protobuf v1.2.0/go.mod h1:6lQm79b+lXiMfvg/cZm0SGofjICqVBUtrP5y
 github.com/golang/protobuf v1.3.1/go.mod h1:6lQm79b+lXiMfvg/cZm0SGofjICqVBUtrP5yJMmIC1U=
 github.com/golang/protobuf v1.3.2/go.mod h1:6lQm79b+lXiMfvg/cZm0SGofjICqVBUtrP5yJMmIC1U=
 github.com/golang/protobuf v1.3.3/go.mod h1:vzj43D7+SQXF/4pzW/hwtAqwc6iTitCiVSaWz5lYuqw=
+github.com/golang/protobuf v1.4.0-rc.1/go.mod h1:ceaxUfeHdC40wWswd/P6IGgMaK3YpKi5j83Wpe3EHw8=
+github.com/golang/protobuf v1.4.0-rc.1.0.20200221234624-67d41d38c208/go.mod h1:xKAWHe0F5eneWXFV3EuXVDTCmh+JuBKY0li0aMyXATA=
+github.com/golang/protobuf v1.4.0-rc.2/go.mod h1:LlEzMj4AhA7rCAGe4KMBDvJI+AwstrUpVNzEA03Pprs=
+github.com/golang/protobuf v1.4.0-rc.4.0.20200313231945-b860323f09d0/go.mod h1:WU3c8KckQ9AFe+yFwt9sWVRKCVIyN9cPHBJSNnbL67w=
+github.com/golang/protobuf v1.4.0/go.mod h1:jodUvKwWbYaEsadDk5Fwe5c77LiNKVO9IDvqG2KuDX0=
+github.com/golang/protobuf v1.4.2/go.mod h1:oDoupMAO8OvCJWAcko0GGGIgR6R6ocIYbsSw735rRwI=
 github.com/golang/protobuf v1.5.0/go.mod h1:FsONVRAS9T7sI+LIUmWTfcYkHO4aIWwzhcaSAoJOfIk=
 github.com/golang/protobuf v1.5.3 h1:KhyjKVUg7Usr/dYsdSqoFveMYd5ko72D+zANwlG1mmg=
 github.com/golang/protobuf v1.5.3/go.mod h1:XVQd3VNwM+JqD3oG2Ue2ip4fOMUkwXdXDdiuN0vRsmY=
@@ -104,6 +113,7 @@ github.com/google/btree v0.0.0-20180813153112-4030bb1f1f0c/go.mod h1:lNA+9X1NB3Z
 github.com/google/btree v1.0.1 h1:gK4Kx5IaGY9CD5sPJ36FHiBJ6ZXl0kilRiiCj+jdYp4=
 github.com/google/btree v1.0.1/go.mod h1:xXMiIv4Fb/0kKde4SpL7qlzvu5cMJDRkFDxJfI9uaxA=
 github.com/google/go-cmp v0.2.0/go.mod h1:oXzfMopK8JAjlY9xF4vHSVASa0yLyX7SntLO5aqRK0M=
+github.com/google/go-cmp v0.3.0/go.mod h1:8QqcDgzrUqlUb/G2PQTWiueGozuR1884gddMywk6iLU=
 github.com/google/go-cmp v0.3.1/go.mod h1:8QqcDgzrUqlUb/G2PQTWiueGozuR1884gddMywk6iLU=
 github.com/google/go-cmp v0.4.0/go.mod h1:v8dTdLbMG2kIc/vJvl+f65V22dbkXbowE6jgT/gNBxE=
 github.com/google/go-cmp v0.5.5/go.mod h1:v8dTdLbMG2kIc/vJvl+f65V22dbkXbowE6jgT/gNBxE=
@@ -160,6 +170,7 @@ github.com/hashicorp/memberlist v0.5.0 h1:EtYPN8DpAURiapus508I4n9CzHs2W+8NZGbmmR
 github.com/hashicorp/memberlist v0.5.0/go.mod h1:yvyXLpo0QaGE59Y7hDTsTzDD25JYBZ4mHgHUZ8lrOI0=
 github.com/hashicorp/serf v0.10.1 h1:Z1H2J60yRKvfDYAOZLd2MU0ND4AH/WDz7xYHDWQsIPY=
 github.com/hashicorp/serf v0.10.1/go.mod h1:yL2t6BqATOLGc5HF7qbFkTfXoPIY0WZdWHfEvMqbG+4=
+github.com/hpcloud/tail v1.0.0/go.mod h1:ab1qPbhIpdTxEkNHXyeSf5vhxWSCs/tWer42PpOxQnU=
 github.com/json-iterator/go v1.1.6/go.mod h1:+SdeFBvtyEkXs7REEP0seUULqWtbJapLOCVDaaPEHmU=
 github.com/json-iterator/go v1.1.9/go.mod h1:KdQUCv79m/52Kvf8AW2vK1V8akMuk1QjK/uOdHXbAo4=
 github.com/json-iterator/go v1.1.12 h1:PV8peI4a0ysnczrg+LtxykD8LfKY9ML6u2jnxaEnrnM=
@@ -201,6 +212,8 @@ github.com/mattn/go-isatty v0.0.16/go.mod h1:kYGgaQfpe5nmfYZH+SKPsOc2e4SrIfOl2e/
 github.com/mattn/go-isatty v0.0.20 h1:xfD0iDuEKnDkl03q4limB+vH+GxLEtL/jb4xVJSWWEY=
 github.com/mattn/go-isatty v0.0.20/go.mod h1:W+V8PltTTMOvKvAeJH7IuucS94S2C6jfK/D7dTCTo3Y=
 github.com/matttproud/golang_protobuf_extensions v1.0.1/go.mod h1:D8He9yQNgCq6Z5Ld7szi9bcBfOoFv/3dc6xSMkL2PC0=
+github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d h1:5PJl274Y63IEHC+7izoQE9x6ikvDFZS2mDVS3drnohI=
+github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d/go.mod h1:01TrycV0kFyexm33Z7vhZRXopbI8J3TDReVlkTgMUxE=
 github.com/miekg/dns v1.1.26/go.mod h1:bPDLeHnStXmXAq1m/Ch/hvfNHr14JKNPMBo3VZKjuso=
 github.com/miekg/dns v1.1.41 h1:WMszZWJG0XmzbK9FEmzH2TVcqYzFesusSIB41b8KHxY=
 github.com/miekg/dns v1.1.41/go.mod h1:p6aan82bvRIyn+zDIv9xYNUpwa73JcSh9BKwknJysuI=
@@ -218,8 +231,19 @@ github.com/modern-go/reflect2 v1.0.1/go.mod h1:bx2lNnkwVCuqBIxFjflWJWanXIb3Rllmb
 github.com/modern-go/reflect2 v1.0.2 h1:xBagoLtFs94CBntxluKeaWgTMpvLxC4ur3nMaC9Gz0M=
 github.com/modern-go/reflect2 v1.0.2/go.mod h1:yWuevngMOJpCy52FWWMvUC8ws7m/LJsjYzDa0/r8luk=
 github.com/mwitkow/go-conntrack v0.0.0-20161129095857-cc309e4a2223/go.mod h1:qRWi+5nqEBWmkhHvq77mSJWrCKwh8bxhgT7d/eI7P4U=
+github.com/nxadm/tail v1.4.4/go.mod h1:kenIhsEOeOJmVchQTgglprH7qJGnHDVpk1VPCcaMI8A=
+github.com/nxadm/tail v1.4.8 h1:nPr65rt6Y5JFSKQO7qToXr7pePgD6Gwiw05lkbyAQTE=
+github.com/nxadm/tail v1.4.8/go.mod h1:+ncqLTQzXmGhMZNUePPaPqPvBxHAIsmXswZKocGu+AU=
 github.com/oapi-codegen/runtime v1.1.1 h1:EXLHh0DXIJnWhdRPN2w4MXAzFyE4CskzhNLUmtpMYro=
 github.com/oapi-codegen/runtime v1.1.1/go.mod h1:SK9X900oXmPWilYR5/WKPzt3Kqxn/uS/+lbpREv+eCg=
+github.com/onsi/ginkgo v1.6.0/go.mod h1:lLunBs/Ym6LB5Z9jYTR76FiuTmxDTDusOGeTQH+WWjE=
+github.com/onsi/ginkgo v1.12.1/go.mod h1:zj2OWP4+oCPe1qIXoGWkgMRwljMUYCdkwsT2108oapk=
+github.com/onsi/ginkgo v1.16.5 h1:8xi0RTUf59SOSfEtZMvwTvXYMzG4gV23XVHOZiXNtnE=
+github.com/onsi/ginkgo v1.16.5/go.mod h1:+E8gABHa3K6zRBolWtd+ROzc/U5bkGt0FwiG042wbpU=
+github.com/onsi/gomega v1.7.1/go.mod h1:XdKZgCCFLUoM/7CFJVPcG8C1xQ1AJ0vpAezJrB7JYyY=
+github.com/onsi/gomega v1.10.1/go.mod h1:iN09h71vgCQne3DLsj+A5owkum+a2tYe+TOCB1ybHNo=
+github.com/onsi/gomega v1.36.2 h1:koNYke6TVk6ZmnyHrCXba/T/MoLBXFjeC1PtvYgw0A8=
+github.com/onsi/gomega v1.36.2/go.mod h1:DdwyADRjrc825LhMEkD76cHR5+pUnjhUN8GlHlRPHzY=
 github.com/opentracing/opentracing-go v1.1.0/go.mod h1:UkNAQd3GIcIGf0SeVgPpRdFStlNbqXla1AfSYxPUl2o=
 github.com/pascaldekloe/goe v0.0.0-20180627143212-57f6aae5913c/go.mod h1:lzWF7FIEvWOWxwDKqyGYQf6ZUaNfKdP144TG7ZOy1lc=
 github.com/pascaldekloe/goe v0.1.0 h1:cBOtyMzM9HTpWjXfbbunk26uA6nG3a8n06Wieeh0MwY=
@@ -283,6 +307,7 @@ github.com/stretchr/objx v0.5.0/go.mod h1:Yh+to48EsGEfYuaHDzXPcE3xhTkx73EhmCGUpE
 github.com/stretchr/testify v1.2.2/go.mod h1:a8OnRcib4nhh0OaRAV+Yts87kKdq0PP7pXfy6kDkUVs=
 github.com/stretchr/testify v1.3.0/go.mod h1:M5WIy9Dh21IEIfnGCwXGc5bZfKNJtfHm1UVUgZn+9EI=
 github.com/stretchr/testify v1.4.0/go.mod h1:j7eGeouHqKxXV5pUuKE4zz7dFj8WfuZ+81PSLYec5m4=
+github.com/stretchr/testify v1.5.1/go.mod h1:5W2xD1RspED5o8YsWQXVCued0rvSQ+mT+I5cxcmMvtA=
 github.com/stretchr/testify v1.7.0/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
 github.com/stretchr/testify v1.7.1/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
 github.com/stretchr/testify v1.7.2/go.mod h1:R6va5+xMeoiuVRoj+gSkQ7d3FALtqAAGI1FQKckRals=
@@ -297,6 +322,8 @@ github.com/twitchyliquid64/golang-asm v0.15.1 h1:SU5vSMR7hnwNxj24w34ZyCi/FmDZTkS
 github.com/twitchyliquid64/golang-asm v0.15.1/go.mod h1:a1lVb/DtPvCB8fslRZhAngC2+aY1QWCk3Cedj/Gdt08=
 github.com/ugorji/go/codec v1.2.12 h1:9LC83zGrHhuUA9l16C9AHXAqEV/2wBQ4nkvumAE65EE=
 github.com/ugorji/go/codec v1.2.12/go.mod h1:UNopzCgEMSXjBc6AOMqYvWC1ktqTAfzJZUZgYf6w6lg=
+github.com/x-cray/logrus-prefixed-formatter v0.5.2 h1:00txxvfBM9muc0jiLIEAkAcIMJzfthRT6usrui8uGmg=
+github.com/x-cray/logrus-prefixed-formatter v0.5.2/go.mod h1:2duySbKsL6M18s5GU7VPsoEPHyzalCE06qoARUCeBBE=
 github.com/yuin/goldmark v1.1.27/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
 github.com/yuin/goldmark v1.2.1/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
 go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0 h1:0nTRpaCaILLdooXAQnfktlL6Zw1ECKEW9DZGH2byi2c=
@@ -332,8 +359,8 @@ golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2/go.mod h1:djNgcEr1/C05ACk
 golang.org/x/crypto v0.0.0-20190923035154-9ee001bba392/go.mod h1:/lpIB1dKB+9EgE3H3cr1v9wB50oz8l4C4h62xy7jSTY=
 golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550/go.mod h1:yigFU9vqHzYiE8UmvKecakEJjdnWj3jj499lnFckfCI=
 golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9/go.mod h1:LzIPMQfyMNhhGPhUkYOs5KpL4U8rLKemX1yGLhDgUto=
-golang.org/x/crypto v0.28.0 h1:GBDwsMXVQi34v5CCYUm2jkJvu4cbtru2U4TN2PSyQnw=
-golang.org/x/crypto v0.28.0/go.mod h1:rmgy+3RHxRZMyY0jjAJShp2zgEdOqj2AO7U0pYmeQ7U=
+golang.org/x/crypto v0.31.0 h1:ihbySMvVjLAeSH1IbfcRTkD/iNscyz8rGzjF/E5hV6U=
+golang.org/x/crypto v0.31.0/go.mod h1:kDsLvtWBEx7MV9tJOj9bnXsPbxwJQ6csT/x4KIN4Ssk=
 golang.org/x/exp v0.0.0-20190121172915-509febef88a4/go.mod h1:CJ0aWSM057203Lf6IL+f9T1iT9GByDxfZKAQTCR3kQA=
 golang.org/x/exp v0.0.0-20230905200255-921286631fa9 h1:GoHiUyI/Tp2nVkLI2mCxVkOjsbSXD66ic0XW0js0R9g=
 golang.org/x/exp v0.0.0-20230905200255-921286631fa9/go.mod h1:S2oDrQGGwySpoQPVqRShND87VCbxmc6bL1Yd2oYrm6k=
@@ -345,6 +372,7 @@ golang.org/x/mod v0.2.0/go.mod h1:s0Qsj1ACt9ePp/hMypM3fl4fZqREWJwdYDEqhRiZZUA=
 golang.org/x/mod v0.3.0/go.mod h1:s0Qsj1ACt9ePp/hMypM3fl4fZqREWJwdYDEqhRiZZUA=
 golang.org/x/net v0.0.0-20180724234803-3673e40ba225/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
 golang.org/x/net v0.0.0-20180826012351-8a410e7b638d/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
+golang.org/x/net v0.0.0-20180906233101-161cd47e91fd/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
 golang.org/x/net v0.0.0-20181114220301-adae6a3d119a/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
 golang.org/x/net v0.0.0-20190213061140-3a22650c66bd/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
 golang.org/x/net v0.0.0-20190311183353-d8887717615a/go.mod h1:t9HGtf8HONx5eT2rtn7q6eTqICYqUVnKs3thJo3Qplg=
@@ -353,11 +381,12 @@ golang.org/x/net v0.0.0-20190613194153-d28f0bde5980/go.mod h1:z5CRVTTTmAJ677TzLL
 golang.org/x/net v0.0.0-20190620200207-3b0461eec859/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
 golang.org/x/net v0.0.0-20190923162816-aa69164e4478/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
 golang.org/x/net v0.0.0-20200226121028-0de0cce0169b/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
+golang.org/x/net v0.0.0-20200520004742-59133d7f0dd7/go.mod h1:qpuaurCH72eLCgpAm/N6yyVIVM9cpaDIP3A8BGJEC5A=
 golang.org/x/net v0.0.0-20201021035429-f5854403a974/go.mod h1:sp8m0HH+o8qH0wwXwYZr8TS3Oi6o0r6Gce1SSxlDquU=
 golang.org/x/net v0.0.0-20210226172049-e18ecbb05110/go.mod h1:m0MpNAwzfU5UDzcl9v0D8zg8gWTRqZa9RBIspLL5mdg=
 golang.org/x/net v0.0.0-20210410081132-afb366fc7cd1/go.mod h1:9tjilg8BloeKEkVJvy7fQ90B1CfIiPueXVOjqfkSzI8=
-golang.org/x/net v0.30.0 h1:AcW1SDZMkb8IpzCdQUaIq2sP4sZ4zw+55h6ynffypl4=
-golang.org/x/net v0.30.0/go.mod h1:2wGyMJ5iFasEhkwi13ChkO/t1ECNC4X4eBKkVFyYFlU=
+golang.org/x/net v0.33.0 h1:74SYHlV8BIgHIFC/LrYkOGIwL19eTYXQ5wc6TBuO36I=
+golang.org/x/net v0.33.0/go.mod h1:HXLR5J+9DxmrqMwG9qjGCxZ+zKXxBru04zlTvWlWuN4=
 golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be/go.mod h1:N/0e6XlmueqKjAGxoOufVs8QHGRruUQn6yWY3a++T0U=
 golang.org/x/sync v0.0.0-20180314180146-1d60e4601c6f/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sync v0.0.0-20181108010431-42b317875d0f/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
@@ -369,19 +398,25 @@ golang.org/x/sync v0.0.0-20210220032951-036812b2e83c/go.mod h1:RxMgew5VJxzue5/jJ
 golang.org/x/sys v0.0.0-20180823144017-11551d06cbcc/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
 golang.org/x/sys v0.0.0-20180830151530-49385e6e1522/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
 golang.org/x/sys v0.0.0-20180905080454-ebe1bf3edb33/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
+golang.org/x/sys v0.0.0-20180909124046-d0be0721c37e/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
 golang.org/x/sys v0.0.0-20181116152217-5ac8a444bdc5/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
 golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
 golang.org/x/sys v0.0.0-20190222072716-a9d3bda3a223/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
 golang.org/x/sys v0.0.0-20190412213103-97732733099d/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20190422165155-953cdadca894/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20190904154756-749cb33beabd/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20190922100055-0a153f010e69/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20190924154521-2837fb4f24fe/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20191005200804-aed5e4c7ecf9/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20191026070338-33540a1f6037/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20191120155948-bd437916bb0e/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20200116001909-b77594299b42/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20200122134326-e047566fdf82/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20200223170610-d5e6a3e2c0ae/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20200323222414-85ca7c5b95cd/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20200930185726-fdedc70b468f/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20201119102817-f84b799fce68/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20210112080510-489259a85091/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20210303074136-134d130e1a04/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20210330210617-4fbd30eecc44/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
@@ -392,15 +427,17 @@ golang.org/x/sys v0.0.0-20220728004956-3c1f35247d10/go.mod h1:oPkhp1MJrh7nUepCBc
 golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.5.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.6.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
-golang.org/x/sys v0.26.0 h1:KHjCJyddX0LoSTb3J+vWpupP9p0oznkqVk/IfjymZbo=
-golang.org/x/sys v0.26.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
+golang.org/x/sys v0.28.0 h1:Fksou7UEQUWlKvIdsqzJmUmCX3cZuD2+P3XyyzwMhlA=
+golang.org/x/sys v0.28.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
 golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1/go.mod h1:bj7SfCRtBDWHUb9snDiAeCFNEtKQo2Wmx5Cou7ajbmo=
+golang.org/x/term v0.27.0 h1:WP60Sv1nlK1T6SupCHbXzSaN0b9wUmsPoRS9b61A23Q=
+golang.org/x/term v0.27.0/go.mod h1:iMsnZpn0cago0GOrHO2+Y7u7JPn5AylBrcoWkElMTSM=
 golang.org/x/text v0.3.0/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
 golang.org/x/text v0.3.2/go.mod h1:bEr9sfX3Q8Zfm5fL9x+3itogRgK3+ptLWKqgva+5dAk=
 golang.org/x/text v0.3.3/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
 golang.org/x/text v0.3.6/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
-golang.org/x/text v0.19.0 h1:kTxAhCbGbxhK0IwgSKiMO5awPoDQ0RpfiVYBfK860YM=
-golang.org/x/text v0.19.0/go.mod h1:BuEKDfySbSR4drPmRPG/7iBdf8hvFMuRexcpahXilzY=
+golang.org/x/text v0.21.0 h1:zyQAAkrwaneQ066sspRyJaG9VNi/YJ1NfzcGB3hZ/qo=
+golang.org/x/text v0.21.0/go.mod h1:4IBbMaMmOPCJ8SecivzSH54+73PCFmPWxNTLm+vZkEQ=
 golang.org/x/tools v0.0.0-20180917221912-90fa682c2a6e/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
 golang.org/x/tools v0.0.0-20190114222345-bf090417da8b/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
 golang.org/x/tools v0.0.0-20190226205152-f727befe758c/go.mod h1:9Yl7xja0Znq3iFh3HoIrodX9oNMXvdceNzlUR8zjMvY=
@@ -410,6 +447,7 @@ golang.org/x/tools v0.0.0-20190907020128-2ca718005c18/go.mod h1:b+2E5dAYhXwXZwtn
 golang.org/x/tools v0.0.0-20191108193012-7d206e10da11/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
 golang.org/x/tools v0.0.0-20191119224855-298f0cb1881e/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
 golang.org/x/tools v0.0.0-20200619180055-7c47624df98f/go.mod h1:EkVYQZoAsY45+roYkvgYkIh4xh/qjgUK9TdY2XT94GE=
+golang.org/x/tools v0.0.0-20201224043029-2b0845dc783e/go.mod h1:emZCQorbCU4vsT4fOWvOPXz4eW1wZW4PmDk9uLelYpA=
 golang.org/x/tools v0.0.0-20210106214847-113979e3529a/go.mod h1:emZCQorbCU4vsT4fOWvOPXz4eW1wZW4PmDk9uLelYpA=
 golang.org/x/xerrors v0.0.0-20190717185122-a985d3407aa7/go.mod h1:I/5z698sn9Ka8TeJc9MKroUUfqBBauWjQqLJ2OPfmY0=
 golang.org/x/xerrors v0.0.0-20191011141410-1b5146add898/go.mod h1:I/5z698sn9Ka8TeJc9MKroUUfqBBauWjQqLJ2OPfmY0=
@@ -429,22 +467,32 @@ google.golang.org/grpc v1.27.0/go.mod h1:qbnxyOmOxrQa7FizSgH+ReBfzJrCY1pSN7KXBS8
 google.golang.org/grpc v1.29.1/go.mod h1:itym6AZVZYACWQqET3MqgPpjcuV5QH3BxFS3IjizoKk=
 google.golang.org/grpc v1.67.1 h1:zWnc1Vrcno+lHZCOofnIMvycFcc0QRGIzm9dhnDX68E=
 google.golang.org/grpc v1.67.1/go.mod h1:1gLDyUQU7CTLJI90u3nXZ9ekeghjeM7pTDZlqFNg2AA=
+google.golang.org/protobuf v0.0.0-20200109180630-ec00e32a8dfd/go.mod h1:DFci5gLYBciE7Vtevhsrf46CRTquxDuWsQurQQe4oz8=
+google.golang.org/protobuf v0.0.0-20200221191635-4d8936d0db64/go.mod h1:kwYJMbMJ01Woi6D6+Kah6886xMZcty6N08ah7+eCXa0=
+google.golang.org/protobuf v0.0.0-20200228230310-ab0ca4ff8a60/go.mod h1:cfTl7dwQJ+fmap5saPgwCLgHXTUD7jkjRqWcaiX5VyM=
+google.golang.org/protobuf v1.20.1-0.20200309200217-e05f789c0967/go.mod h1:A+miEFZTKqfCUM6K7xSMQL9OKL/b6hQv+e19PK+JZNE=
+google.golang.org/protobuf v1.21.0/go.mod h1:47Nbq4nVaFHyn7ilMalzfO3qCViNmqZ2kzikPIcrTAo=
+google.golang.org/protobuf v1.23.0/go.mod h1:EGpADcykh3NcUnDUJcl1+ZksZNG86OlYog2l/sGQquU=
 google.golang.org/protobuf v1.26.0-rc.1/go.mod h1:jlhhOSvTdKEhbULTjvd4ARK9grFBp09yW+WbY/TyQbw=
 google.golang.org/protobuf v1.26.0/go.mod h1:9q0QmTI4eRPtz6boOQmLYwt+qCgq0jsYwAQnmE0givc=
-google.golang.org/protobuf v1.35.1 h1:m3LfL6/Ca+fqnjnlqQXNpFPABW1UD7mjh8KO2mKFytA=
-google.golang.org/protobuf v1.35.1/go.mod h1:9fA7Ob0pmnwhb644+1+CVWFRbNajQ6iRojtC/QF5bRE=
+google.golang.org/protobuf v1.36.1 h1:yBPeRvTftaleIgM3PZ/WBIZ7XM/eEYAaEyCwvyjq/gk=
+google.golang.org/protobuf v1.36.1/go.mod h1:9fA7Ob0pmnwhb644+1+CVWFRbNajQ6iRojtC/QF5bRE=
 gopkg.in/alecthomas/kingpin.v2 v2.2.6/go.mod h1:FMv+mEhP44yOT+4EoQTLFTRgOQ1FBLkstjWtayDeSgw=
 gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
 gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
 gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 h1:YR8cESwS4TdDjEe65xsg0ogRM/Nc3DYOhEAlW+xobZo=
 gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
+gopkg.in/fsnotify.v1 v1.4.7/go.mod h1:Tz8NjZHkW78fSQdbUxIjBTcgA1z1m8ZHf0WmKUhAMys=
 gopkg.in/ini.v1 v1.67.0 h1:Dgnx+6+nfE+IfzjUEISNeydPJh9AXNNsWbGP9KzCsOA=
 gopkg.in/ini.v1 v1.67.0/go.mod h1:pNLf8WUiyNEtQjuu5G5vTm06TEv9tsIgeAvK8hOrP4k=
+gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 h1:uRGJdciOHaEIrze2W8Q3AKkepLTh2hOroT7a+7czfdQ=
+gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7/go.mod h1:dt/ZhP58zS4L8KSrWDmTeBkI65Dw0HsyUHuEVlX15mw=
 gopkg.in/yaml.v2 v2.2.1/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.2/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.4/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.5/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.8/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
+gopkg.in/yaml.v2 v2.3.0/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
 gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
 gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
diff --git a/internal/common/handler/redis/client.go b/internal/common/handler/redis/client.go
index 99c04d3..11abef6 100644
--- a/internal/common/handler/redis/client.go
+++ b/internal/common/handler/redis/client.go
@@ -5,6 +5,7 @@ import (
 	"errors"
 	"time"
 
+	"github.com/ghost-yu/go_shop_second/common/logging"
 	"github.com/redis/go-redis/v9"
 	"github.com/sirupsen/logrus"
 )
@@ -13,16 +14,16 @@ func SetNX(ctx context.Context, client *redis.Client, key, value string, ttl tim
 	now := time.Now()
 	defer func() {
 		l := logrus.WithContext(ctx).WithFields(logrus.Fields{
-			"start": now,
-			"key":   key,
-			"value": value,
-			"err":   err,
-			"cost":  time.Since(now).Milliseconds(),
+			"start":       now,
+			"key":         key,
+			"value":       value,
+			logging.Error: err,
+			logging.Cost:  time.Since(now).Milliseconds(),
 		})
 		if err == nil {
-			l.Info("redis_setnx_success")
+			l.Info("_redis_setnx_success")
 		} else {
-			l.Warn("redis_setnx_error")
+			l.Warn("_redis_setnx_error")
 		}
 	}()
 	if client == nil {
@@ -36,15 +37,15 @@ func Del(ctx context.Context, client *redis.Client, key string) (err error) {
 	now := time.Now()
 	defer func() {
 		l := logrus.WithContext(ctx).WithFields(logrus.Fields{
-			"start": now,
-			"key":   key,
-			"err":   err,
-			"cost":  time.Since(now).Milliseconds(),
+			"start":       now,
+			"key":         key,
+			logging.Error: err,
+			logging.Cost:  time.Since(now).Milliseconds(),
 		})
 		if err == nil {
-			l.Info("redis_del_success")
+			l.Info("_redis_del_success")
 		} else {
-			l.Warn("redis_del_error")
+			l.Warn("_redis_del_error")
 		}
 	}()
 	if client == nil {
diff --git a/internal/common/logging/grpc.go b/internal/common/logging/grpc.go
new file mode 100644
index 0000000..70a37b8
--- /dev/null
+++ b/internal/common/logging/grpc.go
@@ -0,0 +1,29 @@
+package logging
+
+import (
+	"context"
+
+	"github.com/sirupsen/logrus"
+	"google.golang.org/grpc"
+	"google.golang.org/grpc/metadata"
+)
+
+func GRPCUnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
+	fields := logrus.Fields{
+		Args: req,
+	}
+	defer func() {
+		fields[Response] = resp
+		if err != nil {
+			fields[Error] = err.Error()
+			logf(ctx, logrus.ErrorLevel, fields, "%s", "_grpc_request_out")
+		}
+	}()
+	md, exist := metadata.FromIncomingContext(ctx)
+	if exist {
+		fields["grpc_metadata"] = md
+	}
+
+	logf(ctx, logrus.InfoLevel, fields, "%s", "_grpc_request_in")
+	return handler(ctx, req)
+}
diff --git a/internal/common/logging/logrus.go b/internal/common/logging/logrus.go
index 710ab8b..38eadf2 100644
--- a/internal/common/logging/logrus.go
+++ b/internal/common/logging/logrus.go
@@ -1,19 +1,28 @@
 package logging
 
 import (
+	"context"
 	"os"
 	"strconv"
+	"time"
 
+	"github.com/ghost-yu/go_shop_second/common/tracing"
 	"github.com/sirupsen/logrus"
+	prefixed "github.com/x-cray/logrus-prefixed-formatter"
 )
 
+// 要么用logging.Infof, Warnf...
+// 或者直接加hook，用 logrus.Infof...
+
 func Init() {
 	SetFormatter(logrus.StandardLogger())
 	logrus.SetLevel(logrus.DebugLevel)
+	logrus.AddHook(&traceHook{})
 }
 
 func SetFormatter(logger *logrus.Logger) {
 	logger.SetFormatter(&logrus.JSONFormatter{
+		TimestampFormat: time.RFC3339,
 		FieldMap: logrus.FieldMap{
 			logrus.FieldKeyLevel: "severity",
 			logrus.FieldKeyTime:  "time",
@@ -21,8 +30,49 @@ func SetFormatter(logger *logrus.Logger) {
 		},
 	})
 	if isLocal, _ := strconv.ParseBool(os.Getenv("LOCAL_ENV")); isLocal {
-		//logger.SetFormatter(&prefixed.TextFormatter{
-		//	ForceFormatting: false,
-		//})
+		logger.SetFormatter(&prefixed.TextFormatter{
+			ForceColors:     true,
+			ForceFormatting: true,
+			TimestampFormat: time.RFC3339,
+		})
+	}
+}
+
+func logf(ctx context.Context, level logrus.Level, fields logrus.Fields, format string, args ...any) {
+	logrus.WithContext(ctx).WithFields(fields).Logf(level, format, args...)
+}
+
+func InfofWithCost(ctx context.Context, fields logrus.Fields, start time.Time, format string, args ...any) {
+	fields[Cost] = time.Since(start).Milliseconds()
+	Infof(ctx, fields, format, args...)
+}
+
+func Infof(ctx context.Context, fields logrus.Fields, format string, args ...any) {
+	logrus.WithContext(ctx).WithFields(fields).Infof(format, args...)
+}
+
+func Errorf(ctx context.Context, fields logrus.Fields, format string, args ...any) {
+	logrus.WithContext(ctx).WithFields(fields).Errorf(format, args...)
+}
+
+func Warnf(ctx context.Context, fields logrus.Fields, format string, args ...any) {
+	logrus.WithContext(ctx).WithFields(fields).Warnf(format, args...)
+}
+
+func Panicf(ctx context.Context, fields logrus.Fields, format string, args ...any) {
+	logrus.WithContext(ctx).WithFields(fields).Panicf(format, args...)
+}
+
+type traceHook struct{}
+
+func (t traceHook) Levels() []logrus.Level {
+	return logrus.AllLevels
+}
+
+func (t traceHook) Fire(entry *logrus.Entry) error {
+	if entry.Context != nil {
+		entry.Data["trace"] = tracing.TraceID(entry.Context)
+		entry = entry.WithTime(time.Now())
 	}
+	return nil
 }
diff --git a/internal/common/logging/mysql.go b/internal/common/logging/mysql.go
index f2401bf..f9c4cfa 100644
--- a/internal/common/logging/mysql.go
+++ b/internal/common/logging/mysql.go
@@ -14,7 +14,7 @@ const (
 	Args     = "args"
 	Cost     = "cost_ms"
 	Response = "response"
-	Error    = "err"
+	Error    = "error"
 )
 
 type ArgFormatter interface {
@@ -24,7 +24,7 @@ type ArgFormatter interface {
 func WhenMySQL(ctx context.Context, method string, args ...any) (logrus.Fields, func(any, *error)) {
 	fields := logrus.Fields{
 		Method: method,
-		Args:   formatMySQLArgs(args),
+		Args:   formatArgs(args),
 	}
 	start := time.Now()
 	return fields, func(resp any, err *error) {
@@ -37,19 +37,19 @@ func WhenMySQL(ctx context.Context, method string, args ...any) (logrus.Fields,
 			fields[Error] = (*err).Error()
 		}
 
-		logrus.WithContext(ctx).WithFields(fields).Logf(level, "%s", msg)
+		logf(ctx, level, fields, "%s", msg)
 	}
 }
 
-func formatMySQLArgs(args []any) string {
+func formatArgs(args []any) string {
 	var item []string
 	for _, arg := range args {
-		item = append(item, formatMySQLArg(arg))
+		item = append(item, formatArg(arg))
 	}
 	return strings.Join(item, "||")
 }
 
-func formatMySQLArg(arg any) string {
+func formatArg(arg any) string {
 	var (
 		str string
 		err error
diff --git a/internal/common/logging/when.go b/internal/common/logging/when.go
new file mode 100644
index 0000000..f123008
--- /dev/null
+++ b/internal/common/logging/when.go
@@ -0,0 +1,59 @@
+package logging
+
+import (
+	"context"
+	"time"
+
+	"github.com/sirupsen/logrus"
+)
+
+func WhenCommandExecute(ctx context.Context, commandName string, cmd any, err error) {
+	fields := logrus.Fields{
+		"cmd": cmd,
+	}
+	if err == nil {
+		logf(ctx, logrus.InfoLevel, fields, "%s_command_success", commandName)
+	} else {
+		logf(ctx, logrus.ErrorLevel, fields, "%s_command_failed", commandName)
+	}
+}
+
+func WhenRequest(ctx context.Context, method string, args ...any) (logrus.Fields, func(any, *error)) {
+	fields := logrus.Fields{
+		Method: method,
+		Args:   formatArgs(args),
+	}
+	start := time.Now()
+	return fields, func(resp any, err *error) {
+		level, msg := logrus.InfoLevel, "_request_success"
+		fields[Cost] = time.Since(start).Milliseconds()
+		fields[Response] = resp
+
+		if err != nil && (*err != nil) {
+			level, msg = logrus.ErrorLevel, "_request_failed"
+			fields[Error] = (*err).Error()
+		}
+
+		logf(ctx, level, fields, "%s", msg)
+	}
+}
+
+func WhenEventPublish(ctx context.Context, method string, args ...any) (logrus.Fields, func(any, *error)) {
+	fields := logrus.Fields{
+		Method: method,
+		Args:   formatArgs(args),
+	}
+	start := time.Now()
+	return fields, func(resp any, err *error) {
+		level, msg := logrus.InfoLevel, "_mq_publish_success"
+		fields[Cost] = time.Since(start).Milliseconds()
+		fields[Response] = resp
+
+		if err != nil && (*err != nil) {
+			level, msg = logrus.ErrorLevel, "_mq_publish_failed"
+			fields[Error] = (*err).Error()
+		}
+
+		logf(ctx, level, fields, "%s", msg)
+	}
+}
diff --git a/internal/common/middleware/request.go b/internal/common/middleware/request.go
index 3b90df0..5bd5e3e 100644
--- a/internal/common/middleware/request.go
+++ b/internal/common/middleware/request.go
@@ -6,6 +6,7 @@ import (
 	"io"
 	"time"
 
+	"github.com/ghost-yu/go_shop_second/common/logging"
 	"github.com/gin-gonic/gin"
 	"github.com/sirupsen/logrus"
 )
@@ -23,9 +24,9 @@ func requestOut(c *gin.Context, l *logrus.Entry) {
 	start, _ := c.Get("request_start")
 	startTime := start.(time.Time)
 	l.WithContext(c.Request.Context()).WithFields(logrus.Fields{
-		"proc_time_ms": time.Since(startTime).Milliseconds(),
-		"response":     response,
-	}).Info("__request_out")
+		logging.Cost:     time.Since(startTime).Milliseconds(),
+		logging.Response: response,
+	}).Info("_request_out")
 }
 
 func requestIn(c *gin.Context, l *logrus.Entry) {
@@ -36,9 +37,9 @@ func requestIn(c *gin.Context, l *logrus.Entry) {
 	var compactJson bytes.Buffer
 	_ = json.Compact(&compactJson, bodyBytes)
 	l.WithContext(c.Request.Context()).WithFields(logrus.Fields{
-		"start": time.Now().Unix(),
-		"args":  compactJson.String(),
-		"from":  c.RemoteIP(),
-		"uri":   c.Request.RequestURI,
-	}).Info("__request_in")
+		"start":      time.Now().Unix(),
+		logging.Args: compactJson.String(),
+		"from":       c.RemoteIP(),
+		"uri":        c.Request.RequestURI,
+	}).Info("_request_in")
 }
diff --git a/internal/common/server/gprc.go b/internal/common/server/gprc.go
index 14f3132..e325975 100644
--- a/internal/common/server/gprc.go
+++ b/internal/common/server/gprc.go
@@ -3,6 +3,7 @@ package server
 import (
 	"net"
 
+	"github.com/ghost-yu/go_shop_second/common/logging"
 	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
 	grpc_tags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
 	"github.com/sirupsen/logrus"
@@ -33,6 +34,7 @@ func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server))
 		grpc.ChainUnaryInterceptor(
 			grpc_tags.UnaryServerInterceptor(grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor)),
 			grpc_logrus.UnaryServerInterceptor(logrusEntry),
+			logging.GRPCUnaryInterceptor,
 		),
 		grpc.ChainStreamInterceptor(
 			grpc_tags.StreamServerInterceptor(grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor)),
diff --git a/internal/kitchen/infrastructure/consumer/consumer.go b/internal/kitchen/infrastructure/consumer/consumer.go
index f17b444..9e17a9e 100644
--- a/internal/kitchen/infrastructure/consumer/consumer.go
+++ b/internal/kitchen/infrastructure/consumer/consumer.go
@@ -3,12 +3,13 @@ package consumer
 import (
 	"context"
 	"encoding/json"
-	"errors"
 	"fmt"
 	"time"
 
 	"github.com/ghost-yu/go_shop_second/common/broker"
 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+	"github.com/ghost-yu/go_shop_second/common/logging"
+	"github.com/pkg/errors"
 	amqp "github.com/rabbitmq/amqp091-go"
 	"github.com/sirupsen/logrus"
 	"go.opentelemetry.io/otel"
@@ -59,49 +60,50 @@ func (c *Consumer) Listen(ch *amqp.Channel) {
 }
 
 func (c *Consumer) handleMessage(ch *amqp.Channel, msg amqp.Delivery, q amqp.Queue) {
-	var err error
-	logrus.Infof("kitchen receive a message from %s, msg=%v", q.Name, string(msg.Body))
-	ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
 	tr := otel.Tracer("rabbitmq")
-	mqCtx, span := tr.Start(ctx, fmt.Sprintf("rabbitmq.%s.consume", q.Name))
+	ctx, span := tr.Start(broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers), fmt.Sprintf("rabbitmq.%s.consume", q.Name))
+	defer span.End()
+
+	var err error
 	defer func() {
-		span.End()
 		if err != nil {
+			logging.Warnf(ctx, nil, "consume failed||from=%s||msg=%+v||err=%v", q.Name, msg, err)
 			_ = msg.Nack(false, false)
 		} else {
+			logging.Infof(ctx, nil, "%s", "consume success")
 			_ = msg.Ack(false)
 		}
 	}()
 
 	o := &Order{}
-	if err := json.Unmarshal(msg.Body, o); err != nil {
-		logrus.Infof("failed to unmarshall msg to order, err=%v", err)
+	if err = json.Unmarshal(msg.Body, o); err != nil {
+		err = errors.Wrap(err, "error unmarshal msg.body into order")
 		return
 	}
 	if o.Status != "paid" {
 		err = errors.New("order not paid, cannot cook")
 		return
 	}
-	cook(o)
+	cook(ctx, o)
 	span.AddEvent(fmt.Sprintf("order_cook: %v", o))
-	if err := c.orderGRPC.UpdateOrder(mqCtx, &orderpb.Order{
+	if err = c.orderGRPC.UpdateOrder(ctx, &orderpb.Order{
 		ID:          o.ID,
 		CustomerID:  o.CustomerID,
 		Status:      "ready",
 		Items:       o.Items,
 		PaymentLink: o.PaymentLink,
 	}); err != nil {
-		if err = broker.HandleRetry(mqCtx, ch, &msg); err != nil {
-			logrus.Warnf("kitchen: error handling retry: err=%v", err)
+		logging.Errorf(ctx, nil, "error updating order||orderID=%s||err=%v", o.ID, err)
+		if err = broker.HandleRetry(ctx, ch, &msg); err != nil {
+			err = errors.Wrapf(err, "retry_error, error handling retry, messageID=%s||err=%v", msg.MessageId, err)
 		}
 		return
 	}
 	span.AddEvent("kitchen.order.finished.updated")
-	logrus.Info("consume success")
 }
 
-func cook(o *Order) {
-	logrus.Printf("cooking order: %s", o.ID)
+func cook(ctx context.Context, o *Order) {
+	logrus.WithContext(ctx).Printf("cooking order: %s", o.ID)
 	time.Sleep(5 * time.Second)
-	logrus.Printf("order %s done!", o.ID)
+	logrus.WithContext(ctx).Printf("order %s done!", o.ID)
 }
diff --git a/internal/order/adapters/grpc/stock_grpc.go b/internal/order/adapters/grpc/stock_grpc.go
index 549bdd1..958377f 100644
--- a/internal/order/adapters/grpc/stock_grpc.go
+++ b/internal/order/adapters/grpc/stock_grpc.go
@@ -6,7 +6,7 @@ import (
 
 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
-	"github.com/sirupsen/logrus"
+	"github.com/ghost-yu/go_shop_second/common/logging"
 )
 
 type StockGRPC struct {
@@ -17,16 +17,20 @@ func NewStockGRPC(client stockpb.StockServiceClient) *StockGRPC {
 	return &StockGRPC{client: client}
 }
 
-func (s StockGRPC) CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error) {
+func (s StockGRPC) CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (resp *stockpb.CheckIfItemsInStockResponse, err error) {
+	_, dLog := logging.WhenRequest(ctx, "StockGRPC.CheckIfItemsInStock", items)
+	defer dLog(resp, &err)
+
 	if items == nil {
 		return nil, errors.New("grpc items cannot be nil")
 	}
-	resp, err := s.client.CheckIfItemsInStock(ctx, &stockpb.CheckIfItemsInStockRequest{Items: items})
-	logrus.Info("stock_grpc response", resp)
-	return resp, err
+	return s.client.CheckIfItemsInStock(ctx, &stockpb.CheckIfItemsInStockRequest{Items: items})
 }
 
-func (s StockGRPC) GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error) {
+func (s StockGRPC) GetItems(ctx context.Context, itemIDs []string) (items []*orderpb.Item, err error) {
+	_, dLog := logging.WhenRequest(ctx, "StockGRPC.GetItems", items)
+	defer dLog(items, &err)
+
 	resp, err := s.client.GetItems(ctx, &stockpb.GetItemsRequest{ItemIDs: itemIDs})
 	if err != nil {
 		return nil, err
diff --git a/internal/order/adapters/order_inmem_repository.go b/internal/order/adapters/order_inmem_repository.go
index f910521..f50ab5f 100644
--- a/internal/order/adapters/order_inmem_repository.go
+++ b/internal/order/adapters/order_inmem_repository.go
@@ -6,8 +6,8 @@ import (
 	"sync"
 	"time"
 
+	"github.com/ghost-yu/go_shop_second/common/logging"
 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
-	"github.com/sirupsen/logrus"
 )
 
 type MemoryOrderRepository struct {
@@ -16,21 +16,25 @@ type MemoryOrderRepository struct {
 }
 
 func NewMemoryOrderRepository() *MemoryOrderRepository {
-	s := make([]*domain.Order, 0)
-	s = append(s, &domain.Order{
-		ID:          "fake-ID",
-		CustomerID:  "fake-customer-id",
-		Status:      "fake-status",
-		PaymentLink: "fake-payment-link",
-		Items:       nil,
-	})
+	s := []*domain.Order{
+		{
+			ID:          "fake-ID",
+			CustomerID:  "fake-customer-id",
+			Status:      "fake-status",
+			PaymentLink: "fake-payment-link",
+			Items:       nil,
+		},
+	}
 	return &MemoryOrderRepository{
 		lock:  &sync.RWMutex{},
 		store: s,
 	}
 }
 
-func (m *MemoryOrderRepository) Create(_ context.Context, order *domain.Order) (*domain.Order, error) {
+func (m *MemoryOrderRepository) Create(ctx context.Context, order *domain.Order) (created *domain.Order, err error) {
+	_, deferLog := logging.WhenRequest(ctx, "MemoryOrderRepository.Create", map[string]any{"order": order})
+	defer deferLog(created, &err)
+
 	m.lock.Lock()
 	defer m.lock.Unlock()
 	newOrder := &domain.Order{
@@ -40,30 +44,36 @@ func (m *MemoryOrderRepository) Create(_ context.Context, order *domain.Order) (
 		PaymentLink: order.PaymentLink,
 		Items:       order.Items,
 	}
-	m.store = append(m.store, newOrder)
-	logrus.WithFields(logrus.Fields{
-		"input_order":        order,
-		"store_after_create": m.store,
-	}).Info("memory_order_repo_create")
 	return newOrder, nil
 }
 
-func (m *MemoryOrderRepository) Get(_ context.Context, id, customerID string) (*domain.Order, error) {
-	for i, v := range m.store {
-		logrus.Infof("m.store[%d] = %+v", i, v)
-	}
+func (m *MemoryOrderRepository) Get(ctx context.Context, id, customerID string) (got *domain.Order, err error) {
+	_, deferLog := logging.WhenRequest(ctx, "MemoryOrderRepository.Get", map[string]any{
+		"id":         id,
+		"customerID": customerID,
+	})
+	defer deferLog(got, &err)
+
 	m.lock.RLock()
 	defer m.lock.RUnlock()
 	for _, o := range m.store {
 		if o.ID == id && o.CustomerID == customerID {
-			logrus.Infof("memory_order_repo_get||found||id=%s||customerID=%s||res=%+v", id, customerID, *o)
 			return o, nil
 		}
 	}
 	return nil, domain.NotFoundError{OrderID: id}
 }
 
-func (m *MemoryOrderRepository) Update(ctx context.Context, order *domain.Order, updateFn func(context.Context, *domain.Order) (*domain.Order, error)) error {
+func (m *MemoryOrderRepository) Update(
+	ctx context.Context,
+	order *domain.Order,
+	updateFn func(context.Context, *domain.Order) (*domain.Order, error),
+) (err error) {
+	_, deferLog := logging.WhenRequest(ctx, "MemoryOrderRepository.Update", map[string]any{
+		"order": order,
+	})
+	defer deferLog(nil, &err)
+
 	m.lock.Lock()
 	defer m.lock.Unlock()
 	found := false
diff --git a/internal/order/adapters/order_mongo_repository.go b/internal/order/adapters/order_mongo_repository.go
index ba202bd..fa7cf40 100644
--- a/internal/order/adapters/order_mongo_repository.go
+++ b/internal/order/adapters/order_mongo_repository.go
@@ -2,12 +2,11 @@ package adapters
 
 import (
 	"context"
-	"time"
 
 	_ "github.com/ghost-yu/go_shop_second/common/config"
+	"github.com/ghost-yu/go_shop_second/common/logging"
 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
 	"github.com/ghost-yu/go_shop_second/order/entity"
-	"github.com/sirupsen/logrus"
 	"github.com/spf13/viper"
 	"go.mongodb.org/mongo-driver/bson"
 	"go.mongodb.org/mongo-driver/bson/primitive"
@@ -41,7 +40,9 @@ type orderModel struct {
 }
 
 func (r *OrderRepositoryMongo) Create(ctx context.Context, order *domain.Order) (created *domain.Order, err error) {
-	defer r.logWithTag("create", err, order, created)
+	_, deferLog := logging.WhenRequest(ctx, "OrderRepositoryMongo.Create", map[string]any{"order": order})
+	defer deferLog(created, &err)
+
 	write := r.marshalToModel(order)
 	res, err := r.collection().InsertOne(ctx, write)
 	if err != nil {
@@ -52,23 +53,13 @@ func (r *OrderRepositoryMongo) Create(ctx context.Context, order *domain.Order)
 	return created, nil
 }
 
-func (r *OrderRepositoryMongo) logWithTag(tag string, err error, input *domain.Order, result interface{}) {
-	l := logrus.WithFields(logrus.Fields{
-		"tag":            "order_repository_mongo",
-		"input_order":    input,
-		"performed_time": time.Now().Unix(),
-		"err":            err,
-		"result":         result,
+func (r *OrderRepositoryMongo) Get(ctx context.Context, id, customerID string) (got *domain.Order, err error) {
+	_, deferLog := logging.WhenRequest(ctx, "OrderRepositoryMongo.Get", map[string]any{
+		"id":         id,
+		"customerID": customerID,
 	})
-	if err != nil {
-		l.Infof("%s_fail", tag)
-	} else {
-		l.Infof("%s_success", tag)
-	}
-}
+	defer deferLog(got, &err)
 
-func (r *OrderRepositoryMongo) Get(ctx context.Context, id, customerID string) (got *domain.Order, err error) {
-	defer r.logWithTag("get", err, nil, got)
 	read := &orderModel{}
 	mongoID, _ := primitive.ObjectIDFromHex(id)
 	cond := bson.M{"_id": mongoID}
@@ -86,12 +77,13 @@ func (r *OrderRepositoryMongo) Get(ctx context.Context, id, customerID string) (
 func (r *OrderRepositoryMongo) Update(
 	ctx context.Context,
 	order *domain.Order,
-	updateFn func(context.Context, *domain.Order,
-	) (*domain.Order, error)) (err error) {
-	defer r.logWithTag("update", err, order, nil)
-	if order == nil {
-		panic("got nil order")
-	}
+	updateFn func(context.Context, *domain.Order) (*domain.Order, error),
+) (err error) {
+	_, deferLog := logging.WhenRequest(ctx, "OrderRepositoryMongo.Update", map[string]any{
+		"order": order,
+	})
+	defer deferLog(nil, &err)
+
 	// 事务
 	session, err := r.db.StartSession()
 	if err != nil {
@@ -119,9 +111,8 @@ func (r *OrderRepositoryMongo) Update(
 	if err != nil {
 		return
 	}
-	logrus.Infof("update||oldOrder=%+v||updated=%+v", oldOrder, updated)
 	mongoID, _ := primitive.ObjectIDFromHex(oldOrder.ID)
-	res, err := r.collection().UpdateOne(
+	_, err = r.collection().UpdateOne(
 		ctx,
 		bson.M{"_id": mongoID, "customer_id": oldOrder.CustomerID},
 		bson.M{"$set": bson.M{
@@ -132,7 +123,6 @@ func (r *OrderRepositoryMongo) Update(
 	if err != nil {
 		return
 	}
-	r.logWithTag("finish_update", err, order, res)
 	return
 }
 
diff --git a/internal/order/app/command/create_order.go b/internal/order/app/command/create_order.go
index 998302b..83a8a90 100644
--- a/internal/order/app/command/create_order.go
+++ b/internal/order/app/command/create_order.go
@@ -7,6 +7,7 @@ import (
 
 	"github.com/ghost-yu/go_shop_second/common/broker"
 	"github.com/ghost-yu/go_shop_second/common/decorator"
+	"github.com/ghost-yu/go_shop_second/common/logging"
 	"github.com/ghost-yu/go_shop_second/order/app/query"
 	"github.com/ghost-yu/go_shop_second/order/convertor"
 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
@@ -63,6 +64,9 @@ func NewCreateOrderHandler(
 }
 
 func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
+	var err error
+	defer logging.WhenCommandExecute(ctx, "CreateOrderHandler", cmd, err)
+
 	q, err := c.channel.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
 	if err != nil {
 		return nil, err
diff --git a/internal/order/app/command/update_order.go b/internal/order/app/command/update_order.go
index f40716d..8849f02 100644
--- a/internal/order/app/command/update_order.go
+++ b/internal/order/app/command/update_order.go
@@ -4,6 +4,7 @@ import (
 	"context"
 
 	"github.com/ghost-yu/go_shop_second/common/decorator"
+	"github.com/ghost-yu/go_shop_second/common/logging"
 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
 	"github.com/sirupsen/logrus"
 )
@@ -36,11 +37,13 @@ func NewUpdateOrderHandler(
 }
 
 func (c updateOrderHandler) Handle(ctx context.Context, cmd UpdateOrder) (interface{}, error) {
+	var err error
+	defer logging.WhenCommandExecute(ctx, "UpdateOrderHandler", cmd, err)
+
 	if cmd.UpdateFn == nil {
-		logrus.Warnf("updateOrderHandler got nil UpdateFn, order=%#v", cmd.Order)
-		cmd.UpdateFn = func(_ context.Context, order *domain.Order) (*domain.Order, error) { return order, nil }
+		logrus.Panicf("UpdateOrderHandler got nil order, cmd=%+v", cmd)
 	}
-	if err := c.orderRepo.Update(ctx, cmd.Order, cmd.UpdateFn); err != nil {
+	if err = c.orderRepo.Update(ctx, cmd.Order, cmd.UpdateFn); err != nil {
 		return nil, err
 	}
 	return nil, nil
diff --git a/internal/order/go.mod b/internal/order/go.mod
index d780635..e4ad0ac 100644
--- a/internal/order/go.mod
+++ b/internal/order/go.mod
@@ -18,7 +18,7 @@ require (
 	go.mongodb.org/mongo-driver v1.17.1
 	go.opentelemetry.io/otel v1.31.0
 	google.golang.org/grpc v1.67.1
-	google.golang.org/protobuf v1.35.1
+	google.golang.org/protobuf v1.36.1
 )
 
 require (
@@ -59,11 +59,13 @@ require (
 	github.com/magiconair/properties v1.8.7 // indirect
 	github.com/mattn/go-colorable v0.1.13 // indirect
 	github.com/mattn/go-isatty v0.0.20 // indirect
+	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
 	github.com/mitchellh/go-homedir v1.1.0 // indirect
 	github.com/mitchellh/mapstructure v1.5.0 // indirect
 	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
 	github.com/modern-go/reflect2 v1.0.2 // indirect
 	github.com/montanaflynn/stats v0.7.1 // indirect
+	github.com/nxadm/tail v1.4.11 // indirect
 	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
 	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
 	github.com/sagikazarmark/locafero v0.6.0 // indirect
@@ -75,6 +77,7 @@ require (
 	github.com/subosito/gotenv v1.6.0 // indirect
 	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
 	github.com/ugorji/go/codec v1.2.12 // indirect
+	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
 	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
 	github.com/xdg-go/scram v1.1.2 // indirect
 	github.com/xdg-go/stringprep v1.0.4 // indirect
@@ -88,12 +91,13 @@ require (
 	go.opentelemetry.io/otel/trace v1.31.0 // indirect
 	go.uber.org/multierr v1.11.0 // indirect
 	golang.org/x/arch v0.11.0 // indirect
-	golang.org/x/crypto v0.28.0 // indirect
+	golang.org/x/crypto v0.31.0 // indirect
 	golang.org/x/exp v0.0.0-20241009180824-f66d83c29e7c // indirect
-	golang.org/x/net v0.30.0 // indirect
-	golang.org/x/sync v0.8.0 // indirect
-	golang.org/x/sys v0.26.0 // indirect
-	golang.org/x/text v0.19.0 // indirect
+	golang.org/x/net v0.33.0 // indirect
+	golang.org/x/sync v0.10.0 // indirect
+	golang.org/x/sys v0.28.0 // indirect
+	golang.org/x/term v0.27.0 // indirect
+	golang.org/x/text v0.21.0 // indirect
 	google.golang.org/genproto/googleapis/rpc v0.0.0-20241007155032-5fefd90f89a9 // indirect
 	gopkg.in/ini.v1 v1.67.0 // indirect
 	gopkg.in/yaml.v3 v3.0.1 // indirect
diff --git a/internal/order/go.sum b/internal/order/go.sum
index 89f4896..d255009 100644
--- a/internal/order/go.sum
+++ b/internal/order/go.sum
@@ -50,6 +50,7 @@ github.com/fatih/color v1.14.1 h1:qfhVLaG5s+nCROl1zJsZRxFeYrHLqWroPOQ8BWiNb4w=
 github.com/fatih/color v1.14.1/go.mod h1:2oHN61fhTpgcxD3TSWCgKDiH1+x4OiDVVGH8WlgGZGg=
 github.com/frankban/quicktest v1.14.6 h1:7Xjx+VpznH+oBnejlPUj8oUpdxnVs4f8XU8WnHkI4W8=
 github.com/frankban/quicktest v1.14.6/go.mod h1:4ptaffx2x8+WTWXmUCuVU6aPUX1/Mz7zb5vbUoiM6w0=
+github.com/fsnotify/fsnotify v1.6.0/go.mod h1:sl3t1tCWJFWoRz9R8WJCbQihKKwmorjAbSClcnxKAGw=
 github.com/fsnotify/fsnotify v1.7.0 h1:8JEhPFa5W2WU7YfeZzPNqzMP6Lwt7L2715Ggo0nosvA=
 github.com/fsnotify/fsnotify v1.7.0/go.mod h1:40Bi/Hjc2AVfZrqy+aj+yEI+/bRxZnMJyTJwOpGvigM=
 github.com/gabriel-vasile/mimetype v1.4.6 h1:3+PzJTKLkvgjeTbts6msPJt4DixhT4YtFNf1gtGe3zc=
@@ -195,6 +196,8 @@ github.com/mattn/go-isatty v0.0.16/go.mod h1:kYGgaQfpe5nmfYZH+SKPsOc2e4SrIfOl2e/
 github.com/mattn/go-isatty v0.0.20 h1:xfD0iDuEKnDkl03q4limB+vH+GxLEtL/jb4xVJSWWEY=
 github.com/mattn/go-isatty v0.0.20/go.mod h1:W+V8PltTTMOvKvAeJH7IuucS94S2C6jfK/D7dTCTo3Y=
 github.com/matttproud/golang_protobuf_extensions v1.0.1/go.mod h1:D8He9yQNgCq6Z5Ld7szi9bcBfOoFv/3dc6xSMkL2PC0=
+github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d h1:5PJl274Y63IEHC+7izoQE9x6ikvDFZS2mDVS3drnohI=
+github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d/go.mod h1:01TrycV0kFyexm33Z7vhZRXopbI8J3TDReVlkTgMUxE=
 github.com/miekg/dns v1.1.26/go.mod h1:bPDLeHnStXmXAq1m/Ch/hvfNHr14JKNPMBo3VZKjuso=
 github.com/miekg/dns v1.1.41 h1:WMszZWJG0XmzbK9FEmzH2TVcqYzFesusSIB41b8KHxY=
 github.com/miekg/dns v1.1.41/go.mod h1:p6aan82bvRIyn+zDIv9xYNUpwa73JcSh9BKwknJysuI=
@@ -214,8 +217,14 @@ github.com/modern-go/reflect2 v1.0.2/go.mod h1:yWuevngMOJpCy52FWWMvUC8ws7m/LJsjY
 github.com/montanaflynn/stats v0.7.1 h1:etflOAAHORrCC44V+aR6Ftzort912ZU+YLiSTuV8eaE=
 github.com/montanaflynn/stats v0.7.1/go.mod h1:etXPPgVO6n31NxCd9KQUMvCM+ve0ruNzt6R8Bnaayow=
 github.com/mwitkow/go-conntrack v0.0.0-20161129095857-cc309e4a2223/go.mod h1:qRWi+5nqEBWmkhHvq77mSJWrCKwh8bxhgT7d/eI7P4U=
+github.com/nxadm/tail v1.4.11 h1:8feyoE3OzPrcshW5/MJ4sGESc5cqmGkGCWlco4l0bqY=
+github.com/nxadm/tail v1.4.11/go.mod h1:OTaG3NK980DZzxbRq6lEuzgU+mug70nY11sMd4JXXHc=
 github.com/oapi-codegen/runtime v1.1.1 h1:EXLHh0DXIJnWhdRPN2w4MXAzFyE4CskzhNLUmtpMYro=
 github.com/oapi-codegen/runtime v1.1.1/go.mod h1:SK9X900oXmPWilYR5/WKPzt3Kqxn/uS/+lbpREv+eCg=
+github.com/onsi/ginkgo v1.16.5 h1:8xi0RTUf59SOSfEtZMvwTvXYMzG4gV23XVHOZiXNtnE=
+github.com/onsi/ginkgo v1.16.5/go.mod h1:+E8gABHa3K6zRBolWtd+ROzc/U5bkGt0FwiG042wbpU=
+github.com/onsi/gomega v1.36.2 h1:koNYke6TVk6ZmnyHrCXba/T/MoLBXFjeC1PtvYgw0A8=
+github.com/onsi/gomega v1.36.2/go.mod h1:DdwyADRjrc825LhMEkD76cHR5+pUnjhUN8GlHlRPHzY=
 github.com/opentracing/opentracing-go v1.1.0/go.mod h1:UkNAQd3GIcIGf0SeVgPpRdFStlNbqXla1AfSYxPUl2o=
 github.com/pascaldekloe/goe v0.0.0-20180627143212-57f6aae5913c/go.mod h1:lzWF7FIEvWOWxwDKqyGYQf6ZUaNfKdP144TG7ZOy1lc=
 github.com/pascaldekloe/goe v0.1.0 h1:cBOtyMzM9HTpWjXfbbunk26uA6nG3a8n06Wieeh0MwY=
@@ -294,6 +303,8 @@ github.com/twitchyliquid64/golang-asm v0.15.1 h1:SU5vSMR7hnwNxj24w34ZyCi/FmDZTkS
 github.com/twitchyliquid64/golang-asm v0.15.1/go.mod h1:a1lVb/DtPvCB8fslRZhAngC2+aY1QWCk3Cedj/Gdt08=
 github.com/ugorji/go/codec v1.2.12 h1:9LC83zGrHhuUA9l16C9AHXAqEV/2wBQ4nkvumAE65EE=
 github.com/ugorji/go/codec v1.2.12/go.mod h1:UNopzCgEMSXjBc6AOMqYvWC1ktqTAfzJZUZgYf6w6lg=
+github.com/x-cray/logrus-prefixed-formatter v0.5.2 h1:00txxvfBM9muc0jiLIEAkAcIMJzfthRT6usrui8uGmg=
+github.com/x-cray/logrus-prefixed-formatter v0.5.2/go.mod h1:2duySbKsL6M18s5GU7VPsoEPHyzalCE06qoARUCeBBE=
 github.com/xdg-go/pbkdf2 v1.0.0 h1:Su7DPu48wXMwC3bs7MCNG+z4FhcyEuz5dlvchbq0B0c=
 github.com/xdg-go/pbkdf2 v1.0.0/go.mod h1:jrpuAogTd400dnrH08LKmI/xc1MbPOebTwRqcT5RDeI=
 github.com/xdg-go/scram v1.1.2 h1:FHX5I5B4i4hKRVRBCFRxq1iQRej7WO3hhBuJf+UUySY=
@@ -339,8 +350,8 @@ golang.org/x/crypto v0.0.0-20190923035154-9ee001bba392/go.mod h1:/lpIB1dKB+9EgE3
 golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550/go.mod h1:yigFU9vqHzYiE8UmvKecakEJjdnWj3jj499lnFckfCI=
 golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9/go.mod h1:LzIPMQfyMNhhGPhUkYOs5KpL4U8rLKemX1yGLhDgUto=
 golang.org/x/crypto v0.0.0-20210921155107-089bfa567519/go.mod h1:GvvjBRRGRdwPK5ydBHafDWAxML/pGHZbMvKqRZ5+Abc=
-golang.org/x/crypto v0.28.0 h1:GBDwsMXVQi34v5CCYUm2jkJvu4cbtru2U4TN2PSyQnw=
-golang.org/x/crypto v0.28.0/go.mod h1:rmgy+3RHxRZMyY0jjAJShp2zgEdOqj2AO7U0pYmeQ7U=
+golang.org/x/crypto v0.31.0 h1:ihbySMvVjLAeSH1IbfcRTkD/iNscyz8rGzjF/E5hV6U=
+golang.org/x/crypto v0.31.0/go.mod h1:kDsLvtWBEx7MV9tJOj9bnXsPbxwJQ6csT/x4KIN4Ssk=
 golang.org/x/exp v0.0.0-20190121172915-509febef88a4/go.mod h1:CJ0aWSM057203Lf6IL+f9T1iT9GByDxfZKAQTCR3kQA=
 golang.org/x/exp v0.0.0-20241009180824-f66d83c29e7c h1:7dEasQXItcW1xKJ2+gg5VOiBnqWrJc+rq0DPKyvvdbY=
 golang.org/x/exp v0.0.0-20241009180824-f66d83c29e7c/go.mod h1:NQtJDoLvd6faHhE7m4T/1IY708gDefGGjR/iUW8yQQ8=
@@ -366,8 +377,8 @@ golang.org/x/net v0.0.0-20210226172049-e18ecbb05110/go.mod h1:m0MpNAwzfU5UDzcl9v
 golang.org/x/net v0.0.0-20210410081132-afb366fc7cd1/go.mod h1:9tjilg8BloeKEkVJvy7fQ90B1CfIiPueXVOjqfkSzI8=
 golang.org/x/net v0.0.0-20210520170846-37e1c6afe023/go.mod h1:9nx3DQGgdP8bBQD5qxJ1jj9UTztislL4KSBs9R2vV5Y=
 golang.org/x/net v0.0.0-20220722155237-a158d28d115b/go.mod h1:XRhObCWvk6IyKnWLug+ECip1KBveYUHfp+8e9klMJ9c=
-golang.org/x/net v0.30.0 h1:AcW1SDZMkb8IpzCdQUaIq2sP4sZ4zw+55h6ynffypl4=
-golang.org/x/net v0.30.0/go.mod h1:2wGyMJ5iFasEhkwi13ChkO/t1ECNC4X4eBKkVFyYFlU=
+golang.org/x/net v0.33.0 h1:74SYHlV8BIgHIFC/LrYkOGIwL19eTYXQ5wc6TBuO36I=
+golang.org/x/net v0.33.0/go.mod h1:HXLR5J+9DxmrqMwG9qjGCxZ+zKXxBru04zlTvWlWuN4=
 golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be/go.mod h1:N/0e6XlmueqKjAGxoOufVs8QHGRruUQn6yWY3a++T0U=
 golang.org/x/sync v0.0.0-20180314180146-1d60e4601c6f/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sync v0.0.0-20181108010431-42b317875d0f/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
@@ -377,8 +388,8 @@ golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e/go.mod h1:RxMgew5VJxzue5/jJ
 golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sync v0.0.0-20210220032951-036812b2e83c/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
-golang.org/x/sync v0.8.0 h1:3NFvSEYkUoMifnESzZl15y791HH1qU2xm6eCJU5ZPXQ=
-golang.org/x/sync v0.8.0/go.mod h1:Czt+wKu1gCyEFDUtn0jG5QVvpJ6rzVqr5aXyt9drQfk=
+golang.org/x/sync v0.10.0 h1:3NQrjDixjgGwUOCaF8w2+VYHv0Ve/vGYSbdkTa98gmQ=
+golang.org/x/sync v0.10.0/go.mod h1:Czt+wKu1gCyEFDUtn0jG5QVvpJ6rzVqr5aXyt9drQfk=
 golang.org/x/sys v0.0.0-20180823144017-11551d06cbcc/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
 golang.org/x/sys v0.0.0-20180830151530-49385e6e1522/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
 golang.org/x/sys v0.0.0-20180905080454-ebe1bf3edb33/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
@@ -408,20 +419,23 @@ golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8/go.mod h1:oPkhp1MJrh7nUepCBc
 golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20220728004956-3c1f35247d10/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
+golang.org/x/sys v0.0.0-20220908164124-27713097b956/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.5.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.6.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
-golang.org/x/sys v0.26.0 h1:KHjCJyddX0LoSTb3J+vWpupP9p0oznkqVk/IfjymZbo=
-golang.org/x/sys v0.26.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
+golang.org/x/sys v0.28.0 h1:Fksou7UEQUWlKvIdsqzJmUmCX3cZuD2+P3XyyzwMhlA=
+golang.org/x/sys v0.28.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
 golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1/go.mod h1:bj7SfCRtBDWHUb9snDiAeCFNEtKQo2Wmx5Cou7ajbmo=
 golang.org/x/term v0.0.0-20210927222741-03fcf44c2211/go.mod h1:jbD1KX2456YbFQfuXm/mYQcufACuNUgVhRMnK/tPxf8=
+golang.org/x/term v0.27.0 h1:WP60Sv1nlK1T6SupCHbXzSaN0b9wUmsPoRS9b61A23Q=
+golang.org/x/term v0.27.0/go.mod h1:iMsnZpn0cago0GOrHO2+Y7u7JPn5AylBrcoWkElMTSM=
 golang.org/x/text v0.3.0/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
 golang.org/x/text v0.3.2/go.mod h1:bEr9sfX3Q8Zfm5fL9x+3itogRgK3+ptLWKqgva+5dAk=
 golang.org/x/text v0.3.3/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
 golang.org/x/text v0.3.6/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
 golang.org/x/text v0.3.7/go.mod h1:u+2+/6zg+i71rQMx5EYifcz6MCKuco9NR6JIITiCfzQ=
 golang.org/x/text v0.3.8/go.mod h1:E6s5w1FMmriuDzIBO73fBruAKo1PCIq6d2Q6DHfQ8WQ=
-golang.org/x/text v0.19.0 h1:kTxAhCbGbxhK0IwgSKiMO5awPoDQ0RpfiVYBfK860YM=
-golang.org/x/text v0.19.0/go.mod h1:BuEKDfySbSR4drPmRPG/7iBdf8hvFMuRexcpahXilzY=
+golang.org/x/text v0.21.0 h1:zyQAAkrwaneQ066sspRyJaG9VNi/YJ1NfzcGB3hZ/qo=
+golang.org/x/text v0.21.0/go.mod h1:4IBbMaMmOPCJ8SecivzSH54+73PCFmPWxNTLm+vZkEQ=
 golang.org/x/tools v0.0.0-20180917221912-90fa682c2a6e/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
 golang.org/x/tools v0.0.0-20190114222345-bf090417da8b/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
 golang.org/x/tools v0.0.0-20190226205152-f727befe758c/go.mod h1:9Yl7xja0Znq3iFh3HoIrodX9oNMXvdceNzlUR8zjMvY=
@@ -451,8 +465,8 @@ google.golang.org/grpc v1.27.0/go.mod h1:qbnxyOmOxrQa7FizSgH+ReBfzJrCY1pSN7KXBS8
 google.golang.org/grpc v1.29.1/go.mod h1:itym6AZVZYACWQqET3MqgPpjcuV5QH3BxFS3IjizoKk=
 google.golang.org/grpc v1.67.1 h1:zWnc1Vrcno+lHZCOofnIMvycFcc0QRGIzm9dhnDX68E=
 google.golang.org/grpc v1.67.1/go.mod h1:1gLDyUQU7CTLJI90u3nXZ9ekeghjeM7pTDZlqFNg2AA=
-google.golang.org/protobuf v1.35.1 h1:m3LfL6/Ca+fqnjnlqQXNpFPABW1UD7mjh8KO2mKFytA=
-google.golang.org/protobuf v1.35.1/go.mod h1:9fA7Ob0pmnwhb644+1+CVWFRbNajQ6iRojtC/QF5bRE=
+google.golang.org/protobuf v1.36.1 h1:yBPeRvTftaleIgM3PZ/WBIZ7XM/eEYAaEyCwvyjq/gk=
+google.golang.org/protobuf v1.36.1/go.mod h1:9fA7Ob0pmnwhb644+1+CVWFRbNajQ6iRojtC/QF5bRE=
 gopkg.in/alecthomas/kingpin.v2 v2.2.6/go.mod h1:FMv+mEhP44yOT+4EoQTLFTRgOQ1FBLkstjWtayDeSgw=
 gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
 gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
@@ -460,6 +474,8 @@ gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 h1:YR8cESwS4TdDjEe65xsg0ogR
 gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
 gopkg.in/ini.v1 v1.67.0 h1:Dgnx+6+nfE+IfzjUEISNeydPJh9AXNNsWbGP9KzCsOA=
 gopkg.in/ini.v1 v1.67.0/go.mod h1:pNLf8WUiyNEtQjuu5G5vTm06TEv9tsIgeAvK8hOrP4k=
+gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 h1:uRGJdciOHaEIrze2W8Q3AKkepLTh2hOroT7a+7czfdQ=
+gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7/go.mod h1:dt/ZhP58zS4L8KSrWDmTeBkI65Dw0HsyUHuEVlX15mw=
 gopkg.in/yaml.v2 v2.2.1/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.2/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.4/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
diff --git a/internal/order/infrastructure/consumer/consumer.go b/internal/order/infrastructure/consumer/consumer.go
index 511e0e8..43cbb78 100644
--- a/internal/order/infrastructure/consumer/consumer.go
+++ b/internal/order/infrastructure/consumer/consumer.go
@@ -6,9 +6,11 @@ import (
 	"fmt"
 
 	"github.com/ghost-yu/go_shop_second/common/broker"
+	"github.com/ghost-yu/go_shop_second/common/logging"
 	"github.com/ghost-yu/go_shop_second/order/app"
 	"github.com/ghost-yu/go_shop_second/order/app/command"
 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
+	"github.com/pkg/errors"
 	amqp "github.com/rabbitmq/amqp091-go"
 	"github.com/sirupsen/logrus"
 	"go.opentelemetry.io/otel"
@@ -47,23 +49,24 @@ func (c *Consumer) Listen(ch *amqp.Channel) {
 }
 
 func (c *Consumer) handleMessage(ch *amqp.Channel, msg amqp.Delivery, q amqp.Queue) {
-	ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
-	t := otel.Tracer("rabbitmq")
-	_, span := t.Start(ctx, fmt.Sprintf("rabbitmq.%s.consume", q.Name))
+	tr := otel.Tracer("rabbitmq")
+	ctx, span := tr.Start(broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers), fmt.Sprintf("rabbitmq.%s.consume", q.Name))
 	defer span.End()
 
 	var err error
 	defer func() {
 		if err != nil {
+			logging.Warnf(ctx, nil, "consume failed||from=%s||msg=%+v||err=%v", q.Name, msg, err)
 			_ = msg.Nack(false, false)
 		} else {
+			logging.Infof(ctx, nil, "%s", "consume success")
 			_ = msg.Ack(false)
 		}
 	}()
 
 	o := &domain.Order{}
-	if err := json.Unmarshal(msg.Body, o); err != nil {
-		logrus.Infof("error unmarshal msg.body into domain.order, err = %v", err)
+	if err = json.Unmarshal(msg.Body, o); err != nil {
+		err = errors.Wrap(err, "error unmarshal msg.body into domain.order")
 		return
 	}
 	_, err = c.app.Commands.UpdateOrder.Handle(ctx, command.UpdateOrder{
@@ -76,13 +79,12 @@ func (c *Consumer) handleMessage(ch *amqp.Channel, msg amqp.Delivery, q amqp.Que
 		},
 	})
 	if err != nil {
-		logrus.Infof("error updating order, orderID = %s, err = %v", o.ID, err)
+		logging.Errorf(ctx, nil, "error updating order||orderID=%s||err=%v", o.ID, err)
 		if err = broker.HandleRetry(ctx, ch, &msg); err != nil {
-			logrus.Warnf("retry_error, error handling retry, messageID=%s, err=%v", msg.MessageId, err)
+			err = errors.Wrapf(err, "retry_error, error handling retry, messageID=%s||err=%v", msg.MessageId, err)
 		}
 		return
 	}
 
 	span.AddEvent("order.updated")
-	logrus.Info("order consume paid event success!")
 }
diff --git a/internal/order/ports/grpc.go b/internal/order/ports/grpc.go
index b366016..7863a7f 100644
--- a/internal/order/ports/grpc.go
+++ b/internal/order/ports/grpc.go
@@ -10,7 +10,6 @@ import (
 	"github.com/ghost-yu/go_shop_second/order/convertor"
 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
 	"github.com/golang/protobuf/ptypes/empty"
-	"github.com/sirupsen/logrus"
 	"google.golang.org/grpc/codes"
 	"google.golang.org/grpc/status"
 	"google.golang.org/protobuf/types/known/emptypb"
@@ -47,7 +46,6 @@ func (G GRPCServer) GetOrder(ctx context.Context, request *orderpb.GetOrderReque
 }
 
 func (G GRPCServer) UpdateOrder(ctx context.Context, request *orderpb.Order) (_ *emptypb.Empty, err error) {
-	logrus.Infof("order_grpc||request_in||request=%+v", request)
 	order, err := domain.NewOrder(
 		request.ID,
 		request.CustomerID,
diff --git a/internal/payment/adapters/order_grpc.go b/internal/payment/adapters/order_grpc.go
index 6890559..36e131e 100644
--- a/internal/payment/adapters/order_grpc.go
+++ b/internal/payment/adapters/order_grpc.go
@@ -5,7 +5,6 @@ import (
 
 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
 	"github.com/ghost-yu/go_shop_second/common/tracing"
-	"github.com/sirupsen/logrus"
 	"google.golang.org/grpc/status"
 )
 
@@ -18,12 +17,6 @@ func NewOrderGRPC(client orderpb.OrderServiceClient) *OrderGRPC {
 }
 
 func (o OrderGRPC) UpdateOrder(ctx context.Context, order *orderpb.Order) (err error) {
-	defer func() {
-		if err != nil {
-			logrus.Infof("payment_adapter||update_order,err=%v", err)
-		}
-	}()
-
 	ctx, span := tracing.Start(ctx, "order_grpc.update_order")
 	defer span.End()
 
diff --git a/internal/payment/app/command/create_payment.go b/internal/payment/app/command/create_payment.go
index cc7223e..d168300 100644
--- a/internal/payment/app/command/create_payment.go
+++ b/internal/payment/app/command/create_payment.go
@@ -5,10 +5,13 @@ import (
 
 	"github.com/ghost-yu/go_shop_second/common/decorator"
 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+	"github.com/ghost-yu/go_shop_second/common/logging"
 	"github.com/ghost-yu/go_shop_second/payment/domain"
 	"github.com/sirupsen/logrus"
 )
 
+// TODO: ACL 清理
+
 type CreatePayment struct {
 	Order *orderpb.Order
 }
@@ -21,11 +24,13 @@ type createPaymentHandler struct {
 }
 
 func (c createPaymentHandler) Handle(ctx context.Context, cmd CreatePayment) (string, error) {
+	var err error
+	defer logging.WhenCommandExecute(ctx, "CreatePaymentHandler", cmd, err)
+
 	link, err := c.processor.CreatePaymentLink(ctx, cmd.Order)
 	if err != nil {
 		return "", err
 	}
-	logrus.Infof("create payment link for order: %s success, payment link: %s", cmd.Order.ID, link)
 	newOrder := &orderpb.Order{
 		ID:          cmd.Order.ID,
 		CustomerID:  cmd.Order.CustomerID,
diff --git a/internal/payment/go.mod b/internal/payment/go.mod
index dff5f11..5d647ff 100644
--- a/internal/payment/go.mod
+++ b/internal/payment/go.mod
@@ -12,6 +12,7 @@ require (
 	github.com/spf13/viper v1.19.0
 	github.com/stripe/stripe-go/v79 v79.12.0
 	go.opentelemetry.io/otel v1.31.0
+	google.golang.org/grpc v1.67.1
 )
 
 require (
@@ -49,10 +50,12 @@ require (
 	github.com/magiconair/properties v1.8.7 // indirect
 	github.com/mattn/go-colorable v0.1.13 // indirect
 	github.com/mattn/go-isatty v0.0.20 // indirect
+	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
 	github.com/mitchellh/go-homedir v1.1.0 // indirect
 	github.com/mitchellh/mapstructure v1.5.0 // indirect
 	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
 	github.com/modern-go/reflect2 v1.0.2 // indirect
+	github.com/nxadm/tail v1.4.11 // indirect
 	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
 	github.com/sagikazarmark/locafero v0.4.0 // indirect
 	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
@@ -63,7 +66,7 @@ require (
 	github.com/subosito/gotenv v1.6.0 // indirect
 	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
 	github.com/ugorji/go/codec v1.2.12 // indirect
-	github.com/zenazn/goji v1.0.1 // indirect
+	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
 	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0 // indirect
 	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0 // indirect
 	go.opentelemetry.io/contrib/propagators/b3 v1.31.0 // indirect
@@ -74,14 +77,14 @@ require (
 	go.uber.org/atomic v1.9.0 // indirect
 	go.uber.org/multierr v1.9.0 // indirect
 	golang.org/x/arch v0.11.0 // indirect
-	golang.org/x/crypto v0.28.0 // indirect
+	golang.org/x/crypto v0.31.0 // indirect
 	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
-	golang.org/x/net v0.30.0 // indirect
-	golang.org/x/sys v0.26.0 // indirect
-	golang.org/x/text v0.19.0 // indirect
+	golang.org/x/net v0.33.0 // indirect
+	golang.org/x/sys v0.28.0 // indirect
+	golang.org/x/term v0.27.0 // indirect
+	golang.org/x/text v0.21.0 // indirect
 	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
-	google.golang.org/grpc v1.67.1 // indirect
-	google.golang.org/protobuf v1.35.1 // indirect
+	google.golang.org/protobuf v1.36.1 // indirect
 	gopkg.in/ini.v1 v1.67.0 // indirect
 	gopkg.in/yaml.v3 v3.0.1 // indirect
 )
diff --git a/internal/payment/go.sum b/internal/payment/go.sum
index 1685a68..d1f9f0d 100644
--- a/internal/payment/go.sum
+++ b/internal/payment/go.sum
@@ -46,6 +46,7 @@ github.com/fatih/color v1.14.1 h1:qfhVLaG5s+nCROl1zJsZRxFeYrHLqWroPOQ8BWiNb4w=
 github.com/fatih/color v1.14.1/go.mod h1:2oHN61fhTpgcxD3TSWCgKDiH1+x4OiDVVGH8WlgGZGg=
 github.com/frankban/quicktest v1.14.6 h1:7Xjx+VpznH+oBnejlPUj8oUpdxnVs4f8XU8WnHkI4W8=
 github.com/frankban/quicktest v1.14.6/go.mod h1:4ptaffx2x8+WTWXmUCuVU6aPUX1/Mz7zb5vbUoiM6w0=
+github.com/fsnotify/fsnotify v1.6.0/go.mod h1:sl3t1tCWJFWoRz9R8WJCbQihKKwmorjAbSClcnxKAGw=
 github.com/fsnotify/fsnotify v1.7.0 h1:8JEhPFa5W2WU7YfeZzPNqzMP6Lwt7L2715Ggo0nosvA=
 github.com/fsnotify/fsnotify v1.7.0/go.mod h1:40Bi/Hjc2AVfZrqy+aj+yEI+/bRxZnMJyTJwOpGvigM=
 github.com/gabriel-vasile/mimetype v1.4.5 h1:J7wGKdGu33ocBOhGy0z653k/lFKLFDPJMG8Gql0kxn4=
@@ -188,6 +189,8 @@ github.com/mattn/go-isatty v0.0.16/go.mod h1:kYGgaQfpe5nmfYZH+SKPsOc2e4SrIfOl2e/
 github.com/mattn/go-isatty v0.0.20 h1:xfD0iDuEKnDkl03q4limB+vH+GxLEtL/jb4xVJSWWEY=
 github.com/mattn/go-isatty v0.0.20/go.mod h1:W+V8PltTTMOvKvAeJH7IuucS94S2C6jfK/D7dTCTo3Y=
 github.com/matttproud/golang_protobuf_extensions v1.0.1/go.mod h1:D8He9yQNgCq6Z5Ld7szi9bcBfOoFv/3dc6xSMkL2PC0=
+github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d h1:5PJl274Y63IEHC+7izoQE9x6ikvDFZS2mDVS3drnohI=
+github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d/go.mod h1:01TrycV0kFyexm33Z7vhZRXopbI8J3TDReVlkTgMUxE=
 github.com/miekg/dns v1.1.26/go.mod h1:bPDLeHnStXmXAq1m/Ch/hvfNHr14JKNPMBo3VZKjuso=
 github.com/miekg/dns v1.1.41 h1:WMszZWJG0XmzbK9FEmzH2TVcqYzFesusSIB41b8KHxY=
 github.com/miekg/dns v1.1.41/go.mod h1:p6aan82bvRIyn+zDIv9xYNUpwa73JcSh9BKwknJysuI=
@@ -205,6 +208,12 @@ github.com/modern-go/reflect2 v1.0.1/go.mod h1:bx2lNnkwVCuqBIxFjflWJWanXIb3Rllmb
 github.com/modern-go/reflect2 v1.0.2 h1:xBagoLtFs94CBntxluKeaWgTMpvLxC4ur3nMaC9Gz0M=
 github.com/modern-go/reflect2 v1.0.2/go.mod h1:yWuevngMOJpCy52FWWMvUC8ws7m/LJsjYzDa0/r8luk=
 github.com/mwitkow/go-conntrack v0.0.0-20161129095857-cc309e4a2223/go.mod h1:qRWi+5nqEBWmkhHvq77mSJWrCKwh8bxhgT7d/eI7P4U=
+github.com/nxadm/tail v1.4.11 h1:8feyoE3OzPrcshW5/MJ4sGESc5cqmGkGCWlco4l0bqY=
+github.com/nxadm/tail v1.4.11/go.mod h1:OTaG3NK980DZzxbRq6lEuzgU+mug70nY11sMd4JXXHc=
+github.com/onsi/ginkgo v1.16.5 h1:8xi0RTUf59SOSfEtZMvwTvXYMzG4gV23XVHOZiXNtnE=
+github.com/onsi/ginkgo v1.16.5/go.mod h1:+E8gABHa3K6zRBolWtd+ROzc/U5bkGt0FwiG042wbpU=
+github.com/onsi/gomega v1.36.2 h1:koNYke6TVk6ZmnyHrCXba/T/MoLBXFjeC1PtvYgw0A8=
+github.com/onsi/gomega v1.36.2/go.mod h1:DdwyADRjrc825LhMEkD76cHR5+pUnjhUN8GlHlRPHzY=
 github.com/opentracing/opentracing-go v1.1.0/go.mod h1:UkNAQd3GIcIGf0SeVgPpRdFStlNbqXla1AfSYxPUl2o=
 github.com/pascaldekloe/goe v0.0.0-20180627143212-57f6aae5913c/go.mod h1:lzWF7FIEvWOWxwDKqyGYQf6ZUaNfKdP144TG7ZOy1lc=
 github.com/pascaldekloe/goe v0.1.0 h1:cBOtyMzM9HTpWjXfbbunk26uA6nG3a8n06Wieeh0MwY=
@@ -281,10 +290,10 @@ github.com/twitchyliquid64/golang-asm v0.15.1 h1:SU5vSMR7hnwNxj24w34ZyCi/FmDZTkS
 github.com/twitchyliquid64/golang-asm v0.15.1/go.mod h1:a1lVb/DtPvCB8fslRZhAngC2+aY1QWCk3Cedj/Gdt08=
 github.com/ugorji/go/codec v1.2.12 h1:9LC83zGrHhuUA9l16C9AHXAqEV/2wBQ4nkvumAE65EE=
 github.com/ugorji/go/codec v1.2.12/go.mod h1:UNopzCgEMSXjBc6AOMqYvWC1ktqTAfzJZUZgYf6w6lg=
+github.com/x-cray/logrus-prefixed-formatter v0.5.2 h1:00txxvfBM9muc0jiLIEAkAcIMJzfthRT6usrui8uGmg=
+github.com/x-cray/logrus-prefixed-formatter v0.5.2/go.mod h1:2duySbKsL6M18s5GU7VPsoEPHyzalCE06qoARUCeBBE=
 github.com/yuin/goldmark v1.1.27/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
 github.com/yuin/goldmark v1.2.1/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
-github.com/zenazn/goji v1.0.1 h1:4lbD8Mx2h7IvloP7r2C0D6ltZP6Ufip8Hn0wmSK5LR8=
-github.com/zenazn/goji v1.0.1/go.mod h1:7S9M489iMyHBNxwZnk9/EHS098H4/F6TATF2mIxtB1Q=
 go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0 h1:0nTRpaCaILLdooXAQnfktlL6Zw1ECKEW9DZGH2byi2c=
 go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0/go.mod h1:A7aFlp4WSLmeOnFRZwf2dMU+40THPc+rsr6KOwZLOcg=
 go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0 h1:4Pp6oUg3+e/6M4C0A/3kJ2VYa++dsWVTtGgLVj5xtHg=
@@ -318,8 +327,8 @@ golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2/go.mod h1:djNgcEr1/C05ACk
 golang.org/x/crypto v0.0.0-20190923035154-9ee001bba392/go.mod h1:/lpIB1dKB+9EgE3H3cr1v9wB50oz8l4C4h62xy7jSTY=
 golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550/go.mod h1:yigFU9vqHzYiE8UmvKecakEJjdnWj3jj499lnFckfCI=
 golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9/go.mod h1:LzIPMQfyMNhhGPhUkYOs5KpL4U8rLKemX1yGLhDgUto=
-golang.org/x/crypto v0.28.0 h1:GBDwsMXVQi34v5CCYUm2jkJvu4cbtru2U4TN2PSyQnw=
-golang.org/x/crypto v0.28.0/go.mod h1:rmgy+3RHxRZMyY0jjAJShp2zgEdOqj2AO7U0pYmeQ7U=
+golang.org/x/crypto v0.31.0 h1:ihbySMvVjLAeSH1IbfcRTkD/iNscyz8rGzjF/E5hV6U=
+golang.org/x/crypto v0.31.0/go.mod h1:kDsLvtWBEx7MV9tJOj9bnXsPbxwJQ6csT/x4KIN4Ssk=
 golang.org/x/exp v0.0.0-20190121172915-509febef88a4/go.mod h1:CJ0aWSM057203Lf6IL+f9T1iT9GByDxfZKAQTCR3kQA=
 golang.org/x/exp v0.0.0-20230905200255-921286631fa9 h1:GoHiUyI/Tp2nVkLI2mCxVkOjsbSXD66ic0XW0js0R9g=
 golang.org/x/exp v0.0.0-20230905200255-921286631fa9/go.mod h1:S2oDrQGGwySpoQPVqRShND87VCbxmc6bL1Yd2oYrm6k=
@@ -343,8 +352,8 @@ golang.org/x/net v0.0.0-20201021035429-f5854403a974/go.mod h1:sp8m0HH+o8qH0wwXwY
 golang.org/x/net v0.0.0-20210226172049-e18ecbb05110/go.mod h1:m0MpNAwzfU5UDzcl9v0D8zg8gWTRqZa9RBIspLL5mdg=
 golang.org/x/net v0.0.0-20210410081132-afb366fc7cd1/go.mod h1:9tjilg8BloeKEkVJvy7fQ90B1CfIiPueXVOjqfkSzI8=
 golang.org/x/net v0.0.0-20210520170846-37e1c6afe023/go.mod h1:9nx3DQGgdP8bBQD5qxJ1jj9UTztislL4KSBs9R2vV5Y=
-golang.org/x/net v0.30.0 h1:AcW1SDZMkb8IpzCdQUaIq2sP4sZ4zw+55h6ynffypl4=
-golang.org/x/net v0.30.0/go.mod h1:2wGyMJ5iFasEhkwi13ChkO/t1ECNC4X4eBKkVFyYFlU=
+golang.org/x/net v0.33.0 h1:74SYHlV8BIgHIFC/LrYkOGIwL19eTYXQ5wc6TBuO36I=
+golang.org/x/net v0.33.0/go.mod h1:HXLR5J+9DxmrqMwG9qjGCxZ+zKXxBru04zlTvWlWuN4=
 golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be/go.mod h1:N/0e6XlmueqKjAGxoOufVs8QHGRruUQn6yWY3a++T0U=
 golang.org/x/sync v0.0.0-20180314180146-1d60e4601c6f/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sync v0.0.0-20181108010431-42b317875d0f/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
@@ -378,17 +387,20 @@ golang.org/x/sys v0.0.0-20211025201205-69cdffdb9359/go.mod h1:oPkhp1MJrh7nUepCBc
 golang.org/x/sys v0.0.0-20220503163025-988cb79eb6c6/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20220728004956-3c1f35247d10/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
+golang.org/x/sys v0.0.0-20220908164124-27713097b956/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.5.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.6.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
-golang.org/x/sys v0.26.0 h1:KHjCJyddX0LoSTb3J+vWpupP9p0oznkqVk/IfjymZbo=
-golang.org/x/sys v0.26.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
+golang.org/x/sys v0.28.0 h1:Fksou7UEQUWlKvIdsqzJmUmCX3cZuD2+P3XyyzwMhlA=
+golang.org/x/sys v0.28.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
 golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1/go.mod h1:bj7SfCRtBDWHUb9snDiAeCFNEtKQo2Wmx5Cou7ajbmo=
+golang.org/x/term v0.27.0 h1:WP60Sv1nlK1T6SupCHbXzSaN0b9wUmsPoRS9b61A23Q=
+golang.org/x/term v0.27.0/go.mod h1:iMsnZpn0cago0GOrHO2+Y7u7JPn5AylBrcoWkElMTSM=
 golang.org/x/text v0.3.0/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
 golang.org/x/text v0.3.2/go.mod h1:bEr9sfX3Q8Zfm5fL9x+3itogRgK3+ptLWKqgva+5dAk=
 golang.org/x/text v0.3.3/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
 golang.org/x/text v0.3.6/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
-golang.org/x/text v0.19.0 h1:kTxAhCbGbxhK0IwgSKiMO5awPoDQ0RpfiVYBfK860YM=
-golang.org/x/text v0.19.0/go.mod h1:BuEKDfySbSR4drPmRPG/7iBdf8hvFMuRexcpahXilzY=
+golang.org/x/text v0.21.0 h1:zyQAAkrwaneQ066sspRyJaG9VNi/YJ1NfzcGB3hZ/qo=
+golang.org/x/text v0.21.0/go.mod h1:4IBbMaMmOPCJ8SecivzSH54+73PCFmPWxNTLm+vZkEQ=
 golang.org/x/tools v0.0.0-20180917221912-90fa682c2a6e/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
 golang.org/x/tools v0.0.0-20190114222345-bf090417da8b/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
 golang.org/x/tools v0.0.0-20190226205152-f727befe758c/go.mod h1:9Yl7xja0Znq3iFh3HoIrodX9oNMXvdceNzlUR8zjMvY=
@@ -419,8 +431,8 @@ google.golang.org/grpc v1.67.1 h1:zWnc1Vrcno+lHZCOofnIMvycFcc0QRGIzm9dhnDX68E=
 google.golang.org/grpc v1.67.1/go.mod h1:1gLDyUQU7CTLJI90u3nXZ9ekeghjeM7pTDZlqFNg2AA=
 google.golang.org/protobuf v1.26.0-rc.1/go.mod h1:jlhhOSvTdKEhbULTjvd4ARK9grFBp09yW+WbY/TyQbw=
 google.golang.org/protobuf v1.26.0/go.mod h1:9q0QmTI4eRPtz6boOQmLYwt+qCgq0jsYwAQnmE0givc=
-google.golang.org/protobuf v1.35.1 h1:m3LfL6/Ca+fqnjnlqQXNpFPABW1UD7mjh8KO2mKFytA=
-google.golang.org/protobuf v1.35.1/go.mod h1:9fA7Ob0pmnwhb644+1+CVWFRbNajQ6iRojtC/QF5bRE=
+google.golang.org/protobuf v1.36.1 h1:yBPeRvTftaleIgM3PZ/WBIZ7XM/eEYAaEyCwvyjq/gk=
+google.golang.org/protobuf v1.36.1/go.mod h1:9fA7Ob0pmnwhb644+1+CVWFRbNajQ6iRojtC/QF5bRE=
 gopkg.in/alecthomas/kingpin.v2 v2.2.6/go.mod h1:FMv+mEhP44yOT+4EoQTLFTRgOQ1FBLkstjWtayDeSgw=
 gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
 gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
@@ -428,6 +440,8 @@ gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 h1:YR8cESwS4TdDjEe65xsg0ogR
 gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
 gopkg.in/ini.v1 v1.67.0 h1:Dgnx+6+nfE+IfzjUEISNeydPJh9AXNNsWbGP9KzCsOA=
 gopkg.in/ini.v1 v1.67.0/go.mod h1:pNLf8WUiyNEtQjuu5G5vTm06TEv9tsIgeAvK8hOrP4k=
+gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 h1:uRGJdciOHaEIrze2W8Q3AKkepLTh2hOroT7a+7czfdQ=
+gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7/go.mod h1:dt/ZhP58zS4L8KSrWDmTeBkI65Dw0HsyUHuEVlX15mw=
 gopkg.in/yaml.v2 v2.2.1/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.2/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.4/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
diff --git a/internal/payment/http.go b/internal/payment/http.go
index 9671c55..97ade12 100644
--- a/internal/payment/http.go
+++ b/internal/payment/http.go
@@ -1,7 +1,6 @@
 package main
 
 import (
-	"context"
 	"encoding/json"
 	"fmt"
 	"io"
@@ -9,8 +8,10 @@ import (
 
 	"github.com/ghost-yu/go_shop_second/common/broker"
 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+	"github.com/ghost-yu/go_shop_second/common/logging"
 	"github.com/ghost-yu/go_shop_second/payment/domain"
 	"github.com/gin-gonic/gin"
+	"github.com/pkg/errors"
 	amqp "github.com/rabbitmq/amqp091-go"
 	"github.com/sirupsen/logrus"
 	"github.com/spf13/viper"
@@ -33,12 +34,21 @@ func (h *PaymentHandler) RegisterRoutes(c *gin.Engine) {
 }
 
 func (h *PaymentHandler) handleWebhook(c *gin.Context) {
-	logrus.Info("receive webhook from stripe")
+	logrus.WithContext(c.Request.Context()).Info("receive webhook from stripe")
+	var err error
+	defer func() {
+		if err != nil {
+			logging.Warnf(c.Request.Context(), nil, "handleWebhook err=%v", err)
+		} else {
+			logging.Infof(c.Request.Context(), nil, "%s", "handleWebhook success")
+		}
+	}()
+
 	const MaxBodyBytes = int64(65536)
 	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
 	payload, err := io.ReadAll(c.Request.Body)
 	if err != nil {
-		logrus.Infof("Error reading request body: %v\n", err)
+		err = errors.Wrap(err, "Error reading request body")
 		c.JSON(http.StatusServiceUnavailable, err.Error())
 		return
 	}
@@ -47,7 +57,7 @@ func (h *PaymentHandler) handleWebhook(c *gin.Context) {
 		viper.GetString("ENDPOINT_STRIPE_SECRET"))
 
 	if err != nil {
-		logrus.Infof("Error verifying webhook signature: %v\n", err)
+		err = errors.Wrap(err, "error verifying webhook signature")
 		c.JSON(http.StatusBadRequest, err.Error())
 		return
 	}
@@ -55,18 +65,13 @@ func (h *PaymentHandler) handleWebhook(c *gin.Context) {
 	switch event.Type {
 	case stripe.EventTypeCheckoutSessionCompleted:
 		var session stripe.CheckoutSession
-		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
-			logrus.Infof("error unmarshal event.data.raw into session, err = %v", err)
+		if err = json.Unmarshal(event.Data.Raw, &session); err != nil {
+			err = errors.Wrap(err, "error unmarshal event.data.raw into session")
 			c.JSON(http.StatusBadRequest, err.Error())
 			return
 		}
 
 		if session.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
-			logrus.Infof("payment for checkout session %v success!", session.ID)
-
-			ctx, cancel := context.WithCancel(context.TODO())
-			defer cancel()
-
 			var items []*orderpb.Item
 			_ = json.Unmarshal([]byte(session.Metadata["items"]), &items)
 
@@ -78,23 +83,24 @@ func (h *PaymentHandler) handleWebhook(c *gin.Context) {
 				Items:       items,
 			})
 			if err != nil {
-				logrus.Infof("error marshal domain.order, err = %v", err)
+				err = errors.Wrap(err, "error marshal domain.order")
 				c.JSON(http.StatusBadRequest, err.Error())
 				return
 			}
 
+			// TODO: mq logging
 			tr := otel.Tracer("rabbitmq")
-			mqCtx, span := tr.Start(ctx, fmt.Sprintf("rabbitmq.%s.publish", broker.EventOrderPaid))
+			ctx, span := tr.Start(c.Request.Context(), fmt.Sprintf("rabbitmq.%s.publish", broker.EventOrderPaid))
 			defer span.End()
 
-			headers := broker.InjectRabbitMQHeaders(mqCtx)
-			_ = h.channel.PublishWithContext(mqCtx, broker.EventOrderPaid, "", false, false, amqp.Publishing{
+			headers := broker.InjectRabbitMQHeaders(ctx)
+			_ = h.channel.PublishWithContext(ctx, broker.EventOrderPaid, "", false, false, amqp.Publishing{
 				ContentType:  "application/json",
 				DeliveryMode: amqp.Persistent,
 				Body:         marshalledOrder,
 				Headers:      headers,
 			})
-			logrus.Infof("message published to %s, body: %s", broker.EventOrderPaid, string(marshalledOrder))
+			logrus.WithContext(c).Infof("message published to %s, body: %s", broker.EventOrderPaid, string(marshalledOrder))
 		}
 	}
 	c.JSON(http.StatusOK, nil)
diff --git a/internal/payment/infrastructure/consumer/consumer.go b/internal/payment/infrastructure/consumer/consumer.go
index 611c5ae..60f9c19 100644
--- a/internal/payment/infrastructure/consumer/consumer.go
+++ b/internal/payment/infrastructure/consumer/consumer.go
@@ -7,8 +7,10 @@ import (
 
 	"github.com/ghost-yu/go_shop_second/common/broker"
 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+	"github.com/ghost-yu/go_shop_second/common/logging"
 	"github.com/ghost-yu/go_shop_second/payment/app"
 	"github.com/ghost-yu/go_shop_second/payment/app/command"
+	"github.com/pkg/errors"
 	amqp "github.com/rabbitmq/amqp091-go"
 	"github.com/sirupsen/logrus"
 	"go.opentelemetry.io/otel"
@@ -45,34 +47,34 @@ func (c *Consumer) Listen(ch *amqp.Channel) {
 }
 
 func (c *Consumer) handleMessage(ch *amqp.Channel, msg amqp.Delivery, q amqp.Queue) {
-	logrus.Infof("Payment receive a message from %s, msg=%v", q.Name, string(msg.Body))
-	ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
 	tr := otel.Tracer("rabbitmq")
-	_, span := tr.Start(ctx, fmt.Sprintf("rabbitmq.%s.consume", q.Name))
+	ctx, span := tr.Start(broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers), fmt.Sprintf("rabbitmq.%s.consume", q.Name))
 	defer span.End()
 
+	logging.Infof(ctx, nil, "Payment receive a message from %s, msg=%v", q.Name, string(msg.Body))
 	var err error
 	defer func() {
 		if err != nil {
+			logging.Warnf(ctx, nil, "consume failed||from=%s||msg=%+v||err=%v", q.Name, msg, err)
 			_ = msg.Nack(false, false)
 		} else {
+			logging.Infof(ctx, nil, "%s", "consume success")
 			_ = msg.Ack(false)
 		}
 	}()
 
 	o := &orderpb.Order{}
-	if err := json.Unmarshal(msg.Body, o); err != nil {
-		logrus.Infof("failed to unmarshall msg to order, err=%v", err)
+	if err = json.Unmarshal(msg.Body, o); err != nil {
+		err = errors.Wrap(err, "failed to unmarshall msg to order")
 		return
 	}
-	if _, err := c.app.Commands.CreatePayment.Handle(ctx, command.CreatePayment{Order: o}); err != nil {
-		logrus.Infof("failed to create payment, err=%v", err)
+	if _, err = c.app.Commands.CreatePayment.Handle(ctx, command.CreatePayment{Order: o}); err != nil {
+		err = errors.Wrap(err, "failed to create payment")
 		if err = broker.HandleRetry(ctx, ch, &msg); err != nil {
-			logrus.Warnf("retry_error, error handling retry, messageID=%s, err=%v", msg.MessageId, err)
+			err = errors.Wrapf(err, "retry_error, error handling retry, messageID=%s, err=%v", msg.MessageId, err)
 		}
 		return
 	}
 
 	span.AddEvent("payment.created")
-	logrus.Info("consume success")
 }
diff --git a/internal/stock/app/query/check_if_items_in_stock.go b/internal/stock/app/query/check_if_items_in_stock.go
index 2ddaa3a..c61aeed 100644
--- a/internal/stock/app/query/check_if_items_in_stock.go
+++ b/internal/stock/app/query/check_if_items_in_stock.go
@@ -7,6 +7,7 @@ import (
 
 	"github.com/ghost-yu/go_shop_second/common/decorator"
 	"github.com/ghost-yu/go_shop_second/common/handler/redis"
+	"github.com/ghost-yu/go_shop_second/common/logging"
 	domain "github.com/ghost-yu/go_shop_second/stock/domain/stock"
 	"github.com/ghost-yu/go_shop_second/stock/entity"
 	"github.com/ghost-yu/go_shop_second/stock/infrastructure/integration"
@@ -63,7 +64,7 @@ func (h checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfIte
 	}
 	defer func() {
 		if err := unlock(ctx, getLockKey(query)); err != nil {
-			logrus.Warnf("redis unlock fail, err=%v", err)
+			logging.Warnf(ctx, nil, "redis unlock fail, err=%v", err)
 		}
 	}()
 
diff --git a/internal/stock/go.mod b/internal/stock/go.mod
index e5a8122..6dc6ca8 100644
--- a/internal/stock/go.mod
+++ b/internal/stock/go.mod
@@ -58,10 +58,12 @@ require (
 	github.com/magiconair/properties v1.8.7 // indirect
 	github.com/mattn/go-colorable v0.1.13 // indirect
 	github.com/mattn/go-isatty v0.0.20 // indirect
+	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
 	github.com/mitchellh/go-homedir v1.1.0 // indirect
 	github.com/mitchellh/mapstructure v1.5.0 // indirect
 	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
 	github.com/modern-go/reflect2 v1.0.2 // indirect
+	github.com/nxadm/tail v1.4.11 // indirect
 	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
 	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
 	github.com/redis/go-redis/v9 v9.7.0 // indirect
@@ -74,6 +76,7 @@ require (
 	github.com/subosito/gotenv v1.6.0 // indirect
 	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
 	github.com/ugorji/go/codec v1.2.12 // indirect
+	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
 	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0 // indirect
 	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0 // indirect
 	go.opentelemetry.io/contrib/propagators/b3 v1.31.0 // indirect
@@ -85,13 +88,14 @@ require (
 	go.uber.org/atomic v1.9.0 // indirect
 	go.uber.org/multierr v1.9.0 // indirect
 	golang.org/x/arch v0.11.0 // indirect
-	golang.org/x/crypto v0.28.0 // indirect
+	golang.org/x/crypto v0.31.0 // indirect
 	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
-	golang.org/x/net v0.30.0 // indirect
-	golang.org/x/sys v0.26.0 // indirect
-	golang.org/x/text v0.19.0 // indirect
+	golang.org/x/net v0.33.0 // indirect
+	golang.org/x/sys v0.28.0 // indirect
+	golang.org/x/term v0.27.0 // indirect
+	golang.org/x/text v0.21.0 // indirect
 	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
-	google.golang.org/protobuf v1.35.1 // indirect
+	google.golang.org/protobuf v1.36.1 // indirect
 	gopkg.in/ini.v1 v1.67.0 // indirect
 	gopkg.in/yaml.v3 v3.0.1 // indirect
 )
diff --git a/internal/stock/go.sum b/internal/stock/go.sum
index b75a5d5..6e6ddad 100644
--- a/internal/stock/go.sum
+++ b/internal/stock/go.sum
@@ -54,6 +54,7 @@ github.com/fatih/color v1.14.1 h1:qfhVLaG5s+nCROl1zJsZRxFeYrHLqWroPOQ8BWiNb4w=
 github.com/fatih/color v1.14.1/go.mod h1:2oHN61fhTpgcxD3TSWCgKDiH1+x4OiDVVGH8WlgGZGg=
 github.com/frankban/quicktest v1.14.6 h1:7Xjx+VpznH+oBnejlPUj8oUpdxnVs4f8XU8WnHkI4W8=
 github.com/frankban/quicktest v1.14.6/go.mod h1:4ptaffx2x8+WTWXmUCuVU6aPUX1/Mz7zb5vbUoiM6w0=
+github.com/fsnotify/fsnotify v1.6.0/go.mod h1:sl3t1tCWJFWoRz9R8WJCbQihKKwmorjAbSClcnxKAGw=
 github.com/fsnotify/fsnotify v1.7.0 h1:8JEhPFa5W2WU7YfeZzPNqzMP6Lwt7L2715Ggo0nosvA=
 github.com/fsnotify/fsnotify v1.7.0/go.mod h1:40Bi/Hjc2AVfZrqy+aj+yEI+/bRxZnMJyTJwOpGvigM=
 github.com/gabriel-vasile/mimetype v1.4.5 h1:J7wGKdGu33ocBOhGy0z653k/lFKLFDPJMG8Gql0kxn4=
@@ -202,6 +203,8 @@ github.com/mattn/go-isatty v0.0.16/go.mod h1:kYGgaQfpe5nmfYZH+SKPsOc2e4SrIfOl2e/
 github.com/mattn/go-isatty v0.0.20 h1:xfD0iDuEKnDkl03q4limB+vH+GxLEtL/jb4xVJSWWEY=
 github.com/mattn/go-isatty v0.0.20/go.mod h1:W+V8PltTTMOvKvAeJH7IuucS94S2C6jfK/D7dTCTo3Y=
 github.com/matttproud/golang_protobuf_extensions v1.0.1/go.mod h1:D8He9yQNgCq6Z5Ld7szi9bcBfOoFv/3dc6xSMkL2PC0=
+github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d h1:5PJl274Y63IEHC+7izoQE9x6ikvDFZS2mDVS3drnohI=
+github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d/go.mod h1:01TrycV0kFyexm33Z7vhZRXopbI8J3TDReVlkTgMUxE=
 github.com/miekg/dns v1.1.26/go.mod h1:bPDLeHnStXmXAq1m/Ch/hvfNHr14JKNPMBo3VZKjuso=
 github.com/miekg/dns v1.1.41 h1:WMszZWJG0XmzbK9FEmzH2TVcqYzFesusSIB41b8KHxY=
 github.com/miekg/dns v1.1.41/go.mod h1:p6aan82bvRIyn+zDIv9xYNUpwa73JcSh9BKwknJysuI=
@@ -219,6 +222,12 @@ github.com/modern-go/reflect2 v1.0.1/go.mod h1:bx2lNnkwVCuqBIxFjflWJWanXIb3Rllmb
 github.com/modern-go/reflect2 v1.0.2 h1:xBagoLtFs94CBntxluKeaWgTMpvLxC4ur3nMaC9Gz0M=
 github.com/modern-go/reflect2 v1.0.2/go.mod h1:yWuevngMOJpCy52FWWMvUC8ws7m/LJsjYzDa0/r8luk=
 github.com/mwitkow/go-conntrack v0.0.0-20161129095857-cc309e4a2223/go.mod h1:qRWi+5nqEBWmkhHvq77mSJWrCKwh8bxhgT7d/eI7P4U=
+github.com/nxadm/tail v1.4.11 h1:8feyoE3OzPrcshW5/MJ4sGESc5cqmGkGCWlco4l0bqY=
+github.com/nxadm/tail v1.4.11/go.mod h1:OTaG3NK980DZzxbRq6lEuzgU+mug70nY11sMd4JXXHc=
+github.com/onsi/ginkgo v1.16.5 h1:8xi0RTUf59SOSfEtZMvwTvXYMzG4gV23XVHOZiXNtnE=
+github.com/onsi/ginkgo v1.16.5/go.mod h1:+E8gABHa3K6zRBolWtd+ROzc/U5bkGt0FwiG042wbpU=
+github.com/onsi/gomega v1.36.2 h1:koNYke6TVk6ZmnyHrCXba/T/MoLBXFjeC1PtvYgw0A8=
+github.com/onsi/gomega v1.36.2/go.mod h1:DdwyADRjrc825LhMEkD76cHR5+pUnjhUN8GlHlRPHzY=
 github.com/opentracing/opentracing-go v1.1.0/go.mod h1:UkNAQd3GIcIGf0SeVgPpRdFStlNbqXla1AfSYxPUl2o=
 github.com/pascaldekloe/goe v0.0.0-20180627143212-57f6aae5913c/go.mod h1:lzWF7FIEvWOWxwDKqyGYQf6ZUaNfKdP144TG7ZOy1lc=
 github.com/pascaldekloe/goe v0.1.0 h1:cBOtyMzM9HTpWjXfbbunk26uA6nG3a8n06Wieeh0MwY=
@@ -296,6 +305,8 @@ github.com/twitchyliquid64/golang-asm v0.15.1 h1:SU5vSMR7hnwNxj24w34ZyCi/FmDZTkS
 github.com/twitchyliquid64/golang-asm v0.15.1/go.mod h1:a1lVb/DtPvCB8fslRZhAngC2+aY1QWCk3Cedj/Gdt08=
 github.com/ugorji/go/codec v1.2.12 h1:9LC83zGrHhuUA9l16C9AHXAqEV/2wBQ4nkvumAE65EE=
 github.com/ugorji/go/codec v1.2.12/go.mod h1:UNopzCgEMSXjBc6AOMqYvWC1ktqTAfzJZUZgYf6w6lg=
+github.com/x-cray/logrus-prefixed-formatter v0.5.2 h1:00txxvfBM9muc0jiLIEAkAcIMJzfthRT6usrui8uGmg=
+github.com/x-cray/logrus-prefixed-formatter v0.5.2/go.mod h1:2duySbKsL6M18s5GU7VPsoEPHyzalCE06qoARUCeBBE=
 github.com/yuin/goldmark v1.1.27/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
 github.com/yuin/goldmark v1.2.1/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
 go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0 h1:0nTRpaCaILLdooXAQnfktlL6Zw1ECKEW9DZGH2byi2c=
@@ -329,8 +340,8 @@ golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2/go.mod h1:djNgcEr1/C05ACk
 golang.org/x/crypto v0.0.0-20190923035154-9ee001bba392/go.mod h1:/lpIB1dKB+9EgE3H3cr1v9wB50oz8l4C4h62xy7jSTY=
 golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550/go.mod h1:yigFU9vqHzYiE8UmvKecakEJjdnWj3jj499lnFckfCI=
 golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9/go.mod h1:LzIPMQfyMNhhGPhUkYOs5KpL4U8rLKemX1yGLhDgUto=
-golang.org/x/crypto v0.28.0 h1:GBDwsMXVQi34v5CCYUm2jkJvu4cbtru2U4TN2PSyQnw=
-golang.org/x/crypto v0.28.0/go.mod h1:rmgy+3RHxRZMyY0jjAJShp2zgEdOqj2AO7U0pYmeQ7U=
+golang.org/x/crypto v0.31.0 h1:ihbySMvVjLAeSH1IbfcRTkD/iNscyz8rGzjF/E5hV6U=
+golang.org/x/crypto v0.31.0/go.mod h1:kDsLvtWBEx7MV9tJOj9bnXsPbxwJQ6csT/x4KIN4Ssk=
 golang.org/x/exp v0.0.0-20190121172915-509febef88a4/go.mod h1:CJ0aWSM057203Lf6IL+f9T1iT9GByDxfZKAQTCR3kQA=
 golang.org/x/exp v0.0.0-20230905200255-921286631fa9 h1:GoHiUyI/Tp2nVkLI2mCxVkOjsbSXD66ic0XW0js0R9g=
 golang.org/x/exp v0.0.0-20230905200255-921286631fa9/go.mod h1:S2oDrQGGwySpoQPVqRShND87VCbxmc6bL1Yd2oYrm6k=
@@ -354,8 +365,8 @@ golang.org/x/net v0.0.0-20201021035429-f5854403a974/go.mod h1:sp8m0HH+o8qH0wwXwY
 golang.org/x/net v0.0.0-20210226172049-e18ecbb05110/go.mod h1:m0MpNAwzfU5UDzcl9v0D8zg8gWTRqZa9RBIspLL5mdg=
 golang.org/x/net v0.0.0-20210410081132-afb366fc7cd1/go.mod h1:9tjilg8BloeKEkVJvy7fQ90B1CfIiPueXVOjqfkSzI8=
 golang.org/x/net v0.0.0-20210520170846-37e1c6afe023/go.mod h1:9nx3DQGgdP8bBQD5qxJ1jj9UTztislL4KSBs9R2vV5Y=
-golang.org/x/net v0.30.0 h1:AcW1SDZMkb8IpzCdQUaIq2sP4sZ4zw+55h6ynffypl4=
-golang.org/x/net v0.30.0/go.mod h1:2wGyMJ5iFasEhkwi13ChkO/t1ECNC4X4eBKkVFyYFlU=
+golang.org/x/net v0.33.0 h1:74SYHlV8BIgHIFC/LrYkOGIwL19eTYXQ5wc6TBuO36I=
+golang.org/x/net v0.33.0/go.mod h1:HXLR5J+9DxmrqMwG9qjGCxZ+zKXxBru04zlTvWlWuN4=
 golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be/go.mod h1:N/0e6XlmueqKjAGxoOufVs8QHGRruUQn6yWY3a++T0U=
 golang.org/x/sync v0.0.0-20180314180146-1d60e4601c6f/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sync v0.0.0-20181108010431-42b317875d0f/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
@@ -389,17 +400,20 @@ golang.org/x/sys v0.0.0-20211025201205-69cdffdb9359/go.mod h1:oPkhp1MJrh7nUepCBc
 golang.org/x/sys v0.0.0-20220503163025-988cb79eb6c6/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20220728004956-3c1f35247d10/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
+golang.org/x/sys v0.0.0-20220908164124-27713097b956/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.5.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.6.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
-golang.org/x/sys v0.26.0 h1:KHjCJyddX0LoSTb3J+vWpupP9p0oznkqVk/IfjymZbo=
-golang.org/x/sys v0.26.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
+golang.org/x/sys v0.28.0 h1:Fksou7UEQUWlKvIdsqzJmUmCX3cZuD2+P3XyyzwMhlA=
+golang.org/x/sys v0.28.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
 golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1/go.mod h1:bj7SfCRtBDWHUb9snDiAeCFNEtKQo2Wmx5Cou7ajbmo=
+golang.org/x/term v0.27.0 h1:WP60Sv1nlK1T6SupCHbXzSaN0b9wUmsPoRS9b61A23Q=
+golang.org/x/term v0.27.0/go.mod h1:iMsnZpn0cago0GOrHO2+Y7u7JPn5AylBrcoWkElMTSM=
 golang.org/x/text v0.3.0/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
 golang.org/x/text v0.3.2/go.mod h1:bEr9sfX3Q8Zfm5fL9x+3itogRgK3+ptLWKqgva+5dAk=
 golang.org/x/text v0.3.3/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
 golang.org/x/text v0.3.6/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
-golang.org/x/text v0.19.0 h1:kTxAhCbGbxhK0IwgSKiMO5awPoDQ0RpfiVYBfK860YM=
-golang.org/x/text v0.19.0/go.mod h1:BuEKDfySbSR4drPmRPG/7iBdf8hvFMuRexcpahXilzY=
+golang.org/x/text v0.21.0 h1:zyQAAkrwaneQ066sspRyJaG9VNi/YJ1NfzcGB3hZ/qo=
+golang.org/x/text v0.21.0/go.mod h1:4IBbMaMmOPCJ8SecivzSH54+73PCFmPWxNTLm+vZkEQ=
 golang.org/x/tools v0.0.0-20180917221912-90fa682c2a6e/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
 golang.org/x/tools v0.0.0-20190114222345-bf090417da8b/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
 golang.org/x/tools v0.0.0-20190226205152-f727befe758c/go.mod h1:9Yl7xja0Znq3iFh3HoIrodX9oNMXvdceNzlUR8zjMvY=
@@ -430,8 +444,8 @@ google.golang.org/grpc v1.67.1 h1:zWnc1Vrcno+lHZCOofnIMvycFcc0QRGIzm9dhnDX68E=
 google.golang.org/grpc v1.67.1/go.mod h1:1gLDyUQU7CTLJI90u3nXZ9ekeghjeM7pTDZlqFNg2AA=
 google.golang.org/protobuf v1.26.0-rc.1/go.mod h1:jlhhOSvTdKEhbULTjvd4ARK9grFBp09yW+WbY/TyQbw=
 google.golang.org/protobuf v1.26.0/go.mod h1:9q0QmTI4eRPtz6boOQmLYwt+qCgq0jsYwAQnmE0givc=
-google.golang.org/protobuf v1.35.1 h1:m3LfL6/Ca+fqnjnlqQXNpFPABW1UD7mjh8KO2mKFytA=
-google.golang.org/protobuf v1.35.1/go.mod h1:9fA7Ob0pmnwhb644+1+CVWFRbNajQ6iRojtC/QF5bRE=
+google.golang.org/protobuf v1.36.1 h1:yBPeRvTftaleIgM3PZ/WBIZ7XM/eEYAaEyCwvyjq/gk=
+google.golang.org/protobuf v1.36.1/go.mod h1:9fA7Ob0pmnwhb644+1+CVWFRbNajQ6iRojtC/QF5bRE=
 gopkg.in/alecthomas/kingpin.v2 v2.2.6/go.mod h1:FMv+mEhP44yOT+4EoQTLFTRgOQ1FBLkstjWtayDeSgw=
 gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
 gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
@@ -439,6 +453,8 @@ gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 h1:YR8cESwS4TdDjEe65xsg0ogR
 gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
 gopkg.in/ini.v1 v1.67.0 h1:Dgnx+6+nfE+IfzjUEISNeydPJh9AXNNsWbGP9KzCsOA=
 gopkg.in/ini.v1 v1.67.0/go.mod h1:pNLf8WUiyNEtQjuu5G5vTm06TEv9tsIgeAvK8hOrP4k=
+gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 h1:uRGJdciOHaEIrze2W8Q3AKkepLTh2hOroT7a+7czfdQ=
+gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7/go.mod h1:dt/ZhP58zS4L8KSrWDmTeBkI65Dw0HsyUHuEVlX15mw=
 gopkg.in/yaml.v2 v2.2.1/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.2/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.4/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
diff --git a/internal/stock/infrastructure/persistent/mysql.go b/internal/stock/infrastructure/persistent/mysql.go
index dcd4ec6..ca97cb5 100644
--- a/internal/stock/infrastructure/persistent/mysql.go
+++ b/internal/stock/infrastructure/persistent/mysql.go
@@ -70,40 +70,41 @@ func (d MySQL) StartTransaction(fc func(tx *gorm.DB) error) error {
 	return d.db.Transaction(fc)
 }
 
-func (d MySQL) GetStockByID(ctx context.Context, query *builder.Stock) (*StockModel, error) {
+func (d MySQL) GetStockByID(ctx context.Context, query *builder.Stock) (result *StockModel, err error) {
 	_, deferLog := logging.WhenMySQL(ctx, "GetStockByID", query)
-	var result StockModel
-	tx := query.Fill(d.db.WithContext(ctx)).First(&result)
-	defer deferLog(result, &tx.Error)
-	if tx.Error != nil {
-		return nil, tx.Error
+	defer deferLog(result, &err)
+
+	err = query.Fill(d.db.WithContext(ctx)).First(&result).Error
+	if err != nil {
+		return nil, err
 	}
-	return &result, nil
+	return result, nil
 }
 
-func (d MySQL) BatchGetStockByID(ctx context.Context, query *builder.Stock) ([]StockModel, error) {
+func (d MySQL) BatchGetStockByID(ctx context.Context, query *builder.Stock) (result []StockModel, err error) {
 	_, deferLog := logging.WhenMySQL(ctx, "BatchGetStockByID", query)
-	var result []StockModel
-	tx := query.Fill(d.db.WithContext(ctx)).Find(&result)
-	defer deferLog(result, &tx.Error)
-	if tx.Error != nil {
-		return nil, tx.Error
+	defer deferLog(result, &err)
+
+	err = query.Fill(d.db.WithContext(ctx)).Find(&result).Error
+	if err != nil {
+		return nil, err
 	}
 	return result, nil
 }
 
-func (d MySQL) Update(ctx context.Context, tx *gorm.DB, cond *builder.Stock, update map[string]any) error {
-	_, deferLog := logging.WhenMySQL(ctx, "BatchUpdateStock", cond)
+func (d MySQL) Update(ctx context.Context, tx *gorm.DB, cond *builder.Stock, update map[string]any) (err error) {
 	var returning StockModel
+	_, deferLog := logging.WhenMySQL(ctx, "BatchUpdateStock", cond)
+	defer deferLog(returning, &err)
+
 	res := cond.Fill(d.UseTransaction(tx).WithContext(ctx).Model(&returning).Clauses(clause.Returning{})).Updates(update)
-	defer deferLog(returning, &res.Error)
 	return res.Error
 }
 
-func (d MySQL) Create(ctx context.Context, tx *gorm.DB, create *StockModel) error {
-	_, deferLog := logging.WhenMySQL(ctx, "Create", create)
+func (d MySQL) Create(ctx context.Context, tx *gorm.DB, create *StockModel) (err error) {
 	var returning StockModel
-	err := d.UseTransaction(tx).WithContext(ctx).Model(&returning).Clauses(clause.Returning{}).Create(create).Error
+	_, deferLog := logging.WhenMySQL(ctx, "Create", create)
 	defer deferLog(returning, &err)
-	return err
+
+	return d.UseTransaction(tx).WithContext(ctx).Model(&returning).Clauses(clause.Returning{}).Create(create).Error
 }
~~~
