# Lesson Pair Diff Report

- FromBranch: lesson64
- ToBranch: lesson65

## Short Summary

~~~text
 26 files changed, 201 insertions(+), 190 deletions(-)
~~~

## File Stats

~~~text
 .gitignore                                         |  4 +-
 internal/common/decorator/command.go               |  2 +-
 internal/common/decorator/logging.go               | 23 ++++----
 internal/common/decorator/query.go                 |  2 +-
 internal/common/entity/entity.go                   |  8 ++-
 internal/common/go.mod                             |  9 ++--
 internal/common/go.sum                             | 58 ++++----------------
 internal/common/logging/logrus.go                  | 61 +++++++++++++++++++---
 internal/kitchen/go.mod                            | 10 ++--
 internal/kitchen/go.sum                            | 36 ++++++-------
 internal/order/app/command/create_order.go         |  8 +--
 internal/order/app/command/update_order.go         |  6 +--
 internal/order/app/query/get_customer_order.go     |  6 +--
 internal/order/go.mod                              |  7 ++-
 internal/order/go.sum                              | 26 ++++-----
 internal/order/service/application.go              |  7 ++-
 internal/payment/app/command/create_payment.go     |  2 +-
 internal/payment/go.mod                            |  9 ++--
 internal/payment/go.sum                            | 34 ++++++------
 internal/payment/service/application.go            |  3 +-
 internal/stock/adapters/stock_mysql_repository.go  |  4 +-
 .../stock/app/query/check_if_items_in_stock.go     | 16 +++++-
 internal/stock/app/query/get_items.go              |  2 +-
 internal/stock/go.mod                              |  9 ++--
 internal/stock/go.sum                              | 34 ++++++------
 internal/stock/service/application.go              |  5 +-
 26 files changed, 201 insertions(+), 190 deletions(-)
~~~

## Commit Comparison

~~~text
> dcf7511 rotate log
~~~

## Changed Files

~~~text
.gitignore
internal/common/decorator/command.go
internal/common/decorator/logging.go
internal/common/decorator/query.go
internal/common/entity/entity.go
internal/common/go.mod
internal/common/go.sum
internal/common/logging/logrus.go
internal/kitchen/go.mod
internal/kitchen/go.sum
internal/order/app/command/create_order.go
internal/order/app/command/update_order.go
internal/order/app/query/get_customer_order.go
internal/order/go.mod
internal/order/go.sum
internal/order/service/application.go
internal/payment/app/command/create_payment.go
internal/payment/go.mod
internal/payment/go.sum
internal/payment/service/application.go
internal/stock/adapters/stock_mysql_repository.go
internal/stock/app/query/check_if_items_in_stock.go
internal/stock/app/query/get_items.go
internal/stock/go.mod
internal/stock/go.sum
internal/stock/service/application.go
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
.gitignore
internal/common/decorator/command.go
internal/common/decorator/logging.go
internal/common/decorator/query.go
internal/common/entity/entity.go
internal/common/logging/logrus.go
internal/order/app/command/create_order.go
internal/order/app/command/update_order.go
internal/order/app/query/get_customer_order.go
internal/order/service/application.go
internal/payment/app/command/create_payment.go
internal/payment/service/application.go
internal/stock/adapters/stock_mysql_repository.go
internal/stock/app/query/check_if_items_in_stock.go
internal/stock/app/query/get_items.go
internal/stock/service/application.go
~~~

## Full Diff

~~~diff
diff --git a/.gitignore b/.gitignore
index 3672a40..32f1e67 100644
--- a/.gitignore
+++ b/.gitignore
@@ -32,4 +32,6 @@ go.work.sum
 **/bin/
 
 data
-mysql_data
\ No newline at end of file
+mysql_data
+**/log
+*.log
\ No newline at end of file
diff --git a/internal/common/decorator/command.go b/internal/common/decorator/command.go
index b89552d..007ca0c 100644
--- a/internal/common/decorator/command.go
+++ b/internal/common/decorator/command.go
@@ -10,7 +10,7 @@ type CommandHandler[C, R any] interface {
 	Handle(ctx context.Context, cmd C) (R, error)
 }
 
-func ApplyCommandDecorators[C, R any](handler CommandHandler[C, R], logger *logrus.Entry, metricsClient MetricsClient) CommandHandler[C, R] {
+func ApplyCommandDecorators[C, R any](handler CommandHandler[C, R], logger *logrus.Logger, metricsClient MetricsClient) CommandHandler[C, R] {
 	return commandLoggingDecorator[C, R]{
 		logger: logger,
 		base: commandMetricsDecorator[C, R]{
diff --git a/internal/common/decorator/logging.go b/internal/common/decorator/logging.go
index 32f2f93..ca1838b 100644
--- a/internal/common/decorator/logging.go
+++ b/internal/common/decorator/logging.go
@@ -6,26 +6,26 @@ import (
 	"fmt"
 	"strings"
 
+	"github.com/ghost-yu/go_shop_second/common/logging"
 	"github.com/sirupsen/logrus"
 )
 
 type queryLoggingDecorator[C, R any] struct {
-	logger *logrus.Entry
+	logger *logrus.Logger
 	base   QueryHandler[C, R]
 }
 
 func (q queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
 	body, _ := json.Marshal(cmd)
-	logger := q.logger.WithFields(logrus.Fields{
+	fields := logrus.Fields{
 		"query":      generateActionName(cmd),
 		"query_body": string(body),
-	})
-	logger.Debug("Executing query")
+	}
 	defer func() {
 		if err == nil {
-			logger.Info("Query execute successfully")
+			logging.Infof(ctx, fields, "%s", "Query execute successfully")
 		} else {
-			logger.Error("Failed to execute query", err)
+			logging.Errorf(ctx, fields, "Failed to execute query, err=%v", err)
 		}
 	}()
 	result, err = q.base.Handle(ctx, cmd)
@@ -33,22 +33,21 @@ func (q queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result
 }
 
 type commandLoggingDecorator[C, R any] struct {
-	logger *logrus.Entry
+	logger *logrus.Logger
 	base   CommandHandler[C, R]
 }
 
 func (q commandLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
 	body, _ := json.Marshal(cmd)
-	logger := q.logger.WithFields(logrus.Fields{
+	fields := logrus.Fields{
 		"command":      generateActionName(cmd),
 		"command_body": string(body),
-	})
-	logger.Debug("Executing command")
+	}
 	defer func() {
 		if err == nil {
-			logger.Info("Command execute successfully")
+			logging.Infof(ctx, fields, "%s", "Query execute successfully")
 		} else {
-			logger.Error("Failed to execute command", err)
+			logging.Errorf(ctx, fields, "Failed to execute query, err=%v", err)
 		}
 	}()
 	result, err = q.base.Handle(ctx, cmd)
diff --git a/internal/common/decorator/query.go b/internal/common/decorator/query.go
index d3848fe..e3a71c3 100644
--- a/internal/common/decorator/query.go
+++ b/internal/common/decorator/query.go
@@ -12,7 +12,7 @@ type QueryHandler[Q, R any] interface {
 	Handle(ctx context.Context, query Q) (R, error)
 }
 
-func ApplyQueryDecorators[H, R any](handler QueryHandler[H, R], logger *logrus.Entry, metricsClient MetricsClient) QueryHandler[H, R] {
+func ApplyQueryDecorators[H, R any](handler QueryHandler[H, R], logger *logrus.Logger, metricsClient MetricsClient) QueryHandler[H, R] {
 	return queryLoggingDecorator[H, R]{
 		logger: logger,
 		base: queryMetricsDecorator[H, R]{
diff --git a/internal/common/entity/entity.go b/internal/common/entity/entity.go
index b0873a1..855b064 100644
--- a/internal/common/entity/entity.go
+++ b/internal/common/entity/entity.go
@@ -56,7 +56,13 @@ func (iq ItemWithQuantity) validate() error {
 	if iq.ID == "" {
 		invalidFields = append(invalidFields, "ID")
 	}
-	return errors.New(strings.Join(invalidFields, ","))
+	if iq.Quantity < 0 {
+		invalidFields = append(invalidFields, "Quantity")
+	}
+	if len(invalidFields) > 0 {
+		return errors.New(strings.Join(invalidFields, ","))
+	}
+	return nil
 }
 
 func NewItemWithQuantity(ID string, quantity int32) *ItemWithQuantity {
diff --git a/internal/common/go.mod b/internal/common/go.mod
index c208169..3b4c46d 100644
--- a/internal/common/go.mod
+++ b/internal/common/go.mod
@@ -6,13 +6,14 @@ require (
 	github.com/gin-gonic/gin v1.10.0
 	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
 	github.com/hashicorp/consul/api v1.28.2
+	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
 	github.com/oapi-codegen/runtime v1.1.1
 	github.com/pkg/errors v0.9.1
 	github.com/rabbitmq/amqp091-go v1.10.0
 	github.com/redis/go-redis/v9 v9.7.0
+	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
 	github.com/sirupsen/logrus v1.8.1
 	github.com/spf13/viper v1.19.0
-	github.com/x-cray/logrus-prefixed-formatter v0.5.2
 	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0
 	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0
 	go.opentelemetry.io/contrib/propagators/b3 v1.31.0
@@ -54,19 +55,18 @@ require (
 	github.com/hashicorp/golang-lru v0.5.4 // indirect
 	github.com/hashicorp/hcl v1.0.0 // indirect
 	github.com/hashicorp/serf v0.10.1 // indirect
+	github.com/jonboulle/clockwork v0.5.0 // indirect
 	github.com/json-iterator/go v1.1.12 // indirect
 	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
 	github.com/leodido/go-urn v1.4.0 // indirect
+	github.com/lestrrat-go/strftime v1.1.0 // indirect
 	github.com/magiconair/properties v1.8.7 // indirect
 	github.com/mattn/go-colorable v0.1.13 // indirect
 	github.com/mattn/go-isatty v0.0.20 // indirect
-	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
 	github.com/mitchellh/go-homedir v1.1.0 // indirect
 	github.com/mitchellh/mapstructure v1.5.0 // indirect
 	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
 	github.com/modern-go/reflect2 v1.0.2 // indirect
-	github.com/onsi/ginkgo v1.16.5 // indirect
-	github.com/onsi/gomega v1.36.2 // indirect
 	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
 	github.com/sagikazarmark/locafero v0.4.0 // indirect
 	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
@@ -85,7 +85,6 @@ require (
 	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
 	golang.org/x/net v0.33.0 // indirect
 	golang.org/x/sys v0.28.0 // indirect
-	golang.org/x/term v0.27.0 // indirect
 	golang.org/x/text v0.21.0 // indirect
 	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
 	gopkg.in/ini.v1 v1.67.0 // indirect
diff --git a/internal/common/go.sum b/internal/common/go.sum
index 59fdfd2..f953179 100644
--- a/internal/common/go.sum
+++ b/internal/common/go.sum
@@ -58,8 +58,6 @@ github.com/fatih/color v1.14.1 h1:qfhVLaG5s+nCROl1zJsZRxFeYrHLqWroPOQ8BWiNb4w=
 github.com/fatih/color v1.14.1/go.mod h1:2oHN61fhTpgcxD3TSWCgKDiH1+x4OiDVVGH8WlgGZGg=
 github.com/frankban/quicktest v1.14.6 h1:7Xjx+VpznH+oBnejlPUj8oUpdxnVs4f8XU8WnHkI4W8=
 github.com/frankban/quicktest v1.14.6/go.mod h1:4ptaffx2x8+WTWXmUCuVU6aPUX1/Mz7zb5vbUoiM6w0=
-github.com/fsnotify/fsnotify v1.4.7/go.mod h1:jwhsz4b93w/PPRr/qN1Yymfu8t87LnFCMoQvtojpjFo=
-github.com/fsnotify/fsnotify v1.4.9/go.mod h1:znqG4EE+3YCdAaPaxE2ZRY/06pZUdp0tY4IgpuI1SZQ=
 github.com/fsnotify/fsnotify v1.7.0 h1:8JEhPFa5W2WU7YfeZzPNqzMP6Lwt7L2715Ggo0nosvA=
 github.com/fsnotify/fsnotify v1.7.0/go.mod h1:40Bi/Hjc2AVfZrqy+aj+yEI+/bRxZnMJyTJwOpGvigM=
 github.com/gabriel-vasile/mimetype v1.4.5 h1:J7wGKdGu33ocBOhGy0z653k/lFKLFDPJMG8Gql0kxn4=
@@ -88,7 +86,6 @@ github.com/go-playground/universal-translator v0.18.1/go.mod h1:xekY+UJKNuX9WP91
 github.com/go-playground/validator/v10 v10.22.1 h1:40JcKH+bBNGFczGuoBYgX4I6m/i27HYW8P9FDk5PbgA=
 github.com/go-playground/validator/v10 v10.22.1/go.mod h1:dbuPbCMFw/DrkbEynArYaCwl3amGuJotoKCe95atGMM=
 github.com/go-stack/stack v1.8.0/go.mod h1:v0f6uXyyMGvRgIKkXu+yp6POWl0qKG85gN/melR3HDY=
-github.com/go-task/slim-sprig v0.0.0-20210107165309-348f09dbbbc0/go.mod h1:fyg7847qk6SyHyPtNmDHnmrv/HOrqktSC+C9fM+CJOE=
 github.com/goccy/go-json v0.10.3 h1:KZ5WoDbxAIgm2HNbYckL0se1fHD6rz5j4ywS6ebzDqA=
 github.com/goccy/go-json v0.10.3/go.mod h1:oq7eo15ShAhp70Anwd5lgX2pLfOS3QCiwU/PULtXL6M=
 github.com/gogo/protobuf v1.1.1/go.mod h1:r8qH/GZQm5c6nD/R0oafs1akxWv10x8SbQlK7atdtwQ=
@@ -100,12 +97,6 @@ github.com/golang/protobuf v1.2.0/go.mod h1:6lQm79b+lXiMfvg/cZm0SGofjICqVBUtrP5y
 github.com/golang/protobuf v1.3.1/go.mod h1:6lQm79b+lXiMfvg/cZm0SGofjICqVBUtrP5yJMmIC1U=
 github.com/golang/protobuf v1.3.2/go.mod h1:6lQm79b+lXiMfvg/cZm0SGofjICqVBUtrP5yJMmIC1U=
 github.com/golang/protobuf v1.3.3/go.mod h1:vzj43D7+SQXF/4pzW/hwtAqwc6iTitCiVSaWz5lYuqw=
-github.com/golang/protobuf v1.4.0-rc.1/go.mod h1:ceaxUfeHdC40wWswd/P6IGgMaK3YpKi5j83Wpe3EHw8=
-github.com/golang/protobuf v1.4.0-rc.1.0.20200221234624-67d41d38c208/go.mod h1:xKAWHe0F5eneWXFV3EuXVDTCmh+JuBKY0li0aMyXATA=
-github.com/golang/protobuf v1.4.0-rc.2/go.mod h1:LlEzMj4AhA7rCAGe4KMBDvJI+AwstrUpVNzEA03Pprs=
-github.com/golang/protobuf v1.4.0-rc.4.0.20200313231945-b860323f09d0/go.mod h1:WU3c8KckQ9AFe+yFwt9sWVRKCVIyN9cPHBJSNnbL67w=
-github.com/golang/protobuf v1.4.0/go.mod h1:jodUvKwWbYaEsadDk5Fwe5c77LiNKVO9IDvqG2KuDX0=
-github.com/golang/protobuf v1.4.2/go.mod h1:oDoupMAO8OvCJWAcko0GGGIgR6R6ocIYbsSw735rRwI=
 github.com/golang/protobuf v1.5.0/go.mod h1:FsONVRAS9T7sI+LIUmWTfcYkHO4aIWwzhcaSAoJOfIk=
 github.com/golang/protobuf v1.5.3 h1:KhyjKVUg7Usr/dYsdSqoFveMYd5ko72D+zANwlG1mmg=
 github.com/golang/protobuf v1.5.3/go.mod h1:XVQd3VNwM+JqD3oG2Ue2ip4fOMUkwXdXDdiuN0vRsmY=
@@ -113,7 +104,6 @@ github.com/google/btree v0.0.0-20180813153112-4030bb1f1f0c/go.mod h1:lNA+9X1NB3Z
 github.com/google/btree v1.0.1 h1:gK4Kx5IaGY9CD5sPJ36FHiBJ6ZXl0kilRiiCj+jdYp4=
 github.com/google/btree v1.0.1/go.mod h1:xXMiIv4Fb/0kKde4SpL7qlzvu5cMJDRkFDxJfI9uaxA=
 github.com/google/go-cmp v0.2.0/go.mod h1:oXzfMopK8JAjlY9xF4vHSVASa0yLyX7SntLO5aqRK0M=
-github.com/google/go-cmp v0.3.0/go.mod h1:8QqcDgzrUqlUb/G2PQTWiueGozuR1884gddMywk6iLU=
 github.com/google/go-cmp v0.3.1/go.mod h1:8QqcDgzrUqlUb/G2PQTWiueGozuR1884gddMywk6iLU=
 github.com/google/go-cmp v0.4.0/go.mod h1:v8dTdLbMG2kIc/vJvl+f65V22dbkXbowE6jgT/gNBxE=
 github.com/google/go-cmp v0.5.5/go.mod h1:v8dTdLbMG2kIc/vJvl+f65V22dbkXbowE6jgT/gNBxE=
@@ -170,7 +160,8 @@ github.com/hashicorp/memberlist v0.5.0 h1:EtYPN8DpAURiapus508I4n9CzHs2W+8NZGbmmR
 github.com/hashicorp/memberlist v0.5.0/go.mod h1:yvyXLpo0QaGE59Y7hDTsTzDD25JYBZ4mHgHUZ8lrOI0=
 github.com/hashicorp/serf v0.10.1 h1:Z1H2J60yRKvfDYAOZLd2MU0ND4AH/WDz7xYHDWQsIPY=
 github.com/hashicorp/serf v0.10.1/go.mod h1:yL2t6BqATOLGc5HF7qbFkTfXoPIY0WZdWHfEvMqbG+4=
-github.com/hpcloud/tail v1.0.0/go.mod h1:ab1qPbhIpdTxEkNHXyeSf5vhxWSCs/tWer42PpOxQnU=
+github.com/jonboulle/clockwork v0.5.0 h1:Hyh9A8u51kptdkR+cqRpT1EebBwTn1oK9YfGYbdFz6I=
+github.com/jonboulle/clockwork v0.5.0/go.mod h1:3mZlmanh0g2NDKO5TWZVJAfofYk64M7XN3SzBPjZF60=
 github.com/json-iterator/go v1.1.6/go.mod h1:+SdeFBvtyEkXs7REEP0seUULqWtbJapLOCVDaaPEHmU=
 github.com/json-iterator/go v1.1.9/go.mod h1:KdQUCv79m/52Kvf8AW2vK1V8akMuk1QjK/uOdHXbAo4=
 github.com/json-iterator/go v1.1.12 h1:PV8peI4a0ysnczrg+LtxykD8LfKY9ML6u2jnxaEnrnM=
@@ -194,6 +185,12 @@ github.com/kr/text v0.2.0 h1:5Nx0Ya0ZqY2ygV366QzturHI13Jq95ApcVaJBhpS+AY=
 github.com/kr/text v0.2.0/go.mod h1:eLer722TekiGuMkidMxC/pM04lWEeraHUUmBw8l2grE=
 github.com/leodido/go-urn v1.4.0 h1:WT9HwE9SGECu3lg4d/dIA+jxlljEa1/ffXKmRjqdmIQ=
 github.com/leodido/go-urn v1.4.0/go.mod h1:bvxc+MVxLKB4z00jd1z+Dvzr47oO32F/QSNjSBOlFxI=
+github.com/lestrrat-go/envload v0.0.0-20180220234015-a3eb8ddeffcc h1:RKf14vYWi2ttpEmkA4aQ3j4u9dStX2t4M8UM6qqNsG8=
+github.com/lestrrat-go/envload v0.0.0-20180220234015-a3eb8ddeffcc/go.mod h1:kopuH9ugFRkIXf3YoqHKyrJ9YfUFsckUU9S7B+XP+is=
+github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible h1:Y6sqxHMyB1D2YSzWkLibYKgg+SwmyFU9dF2hn6MdTj4=
+github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible/go.mod h1:ZQnN8lSECaebrkQytbHj4xNgtg8CR7RYXnPok8e0EHA=
+github.com/lestrrat-go/strftime v1.1.0 h1:gMESpZy44/4pXLO/m+sL0yBd1W6LjgjrrD4a68Gapyg=
+github.com/lestrrat-go/strftime v1.1.0/go.mod h1:uzeIB52CeUJenCo1syghlugshMysrqUT51HlxphXVeI=
 github.com/magiconair/properties v1.8.7 h1:IeQXZAiQcpL9mgcAe1Nu6cX9LLw6ExEHKjN0VQdvPDY=
 github.com/magiconair/properties v1.8.7/go.mod h1:Dhd985XPs7jluiymwWYZ0G4Z61jb3vdS329zhj2hYo0=
 github.com/mattn/go-colorable v0.0.9/go.mod h1:9vuHe8Xs5qXnSaW/c/ABM9alt+Vo+STaOChaDxuIBZU=
@@ -212,8 +209,6 @@ github.com/mattn/go-isatty v0.0.16/go.mod h1:kYGgaQfpe5nmfYZH+SKPsOc2e4SrIfOl2e/
 github.com/mattn/go-isatty v0.0.20 h1:xfD0iDuEKnDkl03q4limB+vH+GxLEtL/jb4xVJSWWEY=
 github.com/mattn/go-isatty v0.0.20/go.mod h1:W+V8PltTTMOvKvAeJH7IuucS94S2C6jfK/D7dTCTo3Y=
 github.com/matttproud/golang_protobuf_extensions v1.0.1/go.mod h1:D8He9yQNgCq6Z5Ld7szi9bcBfOoFv/3dc6xSMkL2PC0=
-github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d h1:5PJl274Y63IEHC+7izoQE9x6ikvDFZS2mDVS3drnohI=
-github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d/go.mod h1:01TrycV0kFyexm33Z7vhZRXopbI8J3TDReVlkTgMUxE=
 github.com/miekg/dns v1.1.26/go.mod h1:bPDLeHnStXmXAq1m/Ch/hvfNHr14JKNPMBo3VZKjuso=
 github.com/miekg/dns v1.1.41 h1:WMszZWJG0XmzbK9FEmzH2TVcqYzFesusSIB41b8KHxY=
 github.com/miekg/dns v1.1.41/go.mod h1:p6aan82bvRIyn+zDIv9xYNUpwa73JcSh9BKwknJysuI=
@@ -231,19 +226,8 @@ github.com/modern-go/reflect2 v1.0.1/go.mod h1:bx2lNnkwVCuqBIxFjflWJWanXIb3Rllmb
 github.com/modern-go/reflect2 v1.0.2 h1:xBagoLtFs94CBntxluKeaWgTMpvLxC4ur3nMaC9Gz0M=
 github.com/modern-go/reflect2 v1.0.2/go.mod h1:yWuevngMOJpCy52FWWMvUC8ws7m/LJsjYzDa0/r8luk=
 github.com/mwitkow/go-conntrack v0.0.0-20161129095857-cc309e4a2223/go.mod h1:qRWi+5nqEBWmkhHvq77mSJWrCKwh8bxhgT7d/eI7P4U=
-github.com/nxadm/tail v1.4.4/go.mod h1:kenIhsEOeOJmVchQTgglprH7qJGnHDVpk1VPCcaMI8A=
-github.com/nxadm/tail v1.4.8 h1:nPr65rt6Y5JFSKQO7qToXr7pePgD6Gwiw05lkbyAQTE=
-github.com/nxadm/tail v1.4.8/go.mod h1:+ncqLTQzXmGhMZNUePPaPqPvBxHAIsmXswZKocGu+AU=
 github.com/oapi-codegen/runtime v1.1.1 h1:EXLHh0DXIJnWhdRPN2w4MXAzFyE4CskzhNLUmtpMYro=
 github.com/oapi-codegen/runtime v1.1.1/go.mod h1:SK9X900oXmPWilYR5/WKPzt3Kqxn/uS/+lbpREv+eCg=
-github.com/onsi/ginkgo v1.6.0/go.mod h1:lLunBs/Ym6LB5Z9jYTR76FiuTmxDTDusOGeTQH+WWjE=
-github.com/onsi/ginkgo v1.12.1/go.mod h1:zj2OWP4+oCPe1qIXoGWkgMRwljMUYCdkwsT2108oapk=
-github.com/onsi/ginkgo v1.16.5 h1:8xi0RTUf59SOSfEtZMvwTvXYMzG4gV23XVHOZiXNtnE=
-github.com/onsi/ginkgo v1.16.5/go.mod h1:+E8gABHa3K6zRBolWtd+ROzc/U5bkGt0FwiG042wbpU=
-github.com/onsi/gomega v1.7.1/go.mod h1:XdKZgCCFLUoM/7CFJVPcG8C1xQ1AJ0vpAezJrB7JYyY=
-github.com/onsi/gomega v1.10.1/go.mod h1:iN09h71vgCQne3DLsj+A5owkum+a2tYe+TOCB1ybHNo=
-github.com/onsi/gomega v1.36.2 h1:koNYke6TVk6ZmnyHrCXba/T/MoLBXFjeC1PtvYgw0A8=
-github.com/onsi/gomega v1.36.2/go.mod h1:DdwyADRjrc825LhMEkD76cHR5+pUnjhUN8GlHlRPHzY=
 github.com/opentracing/opentracing-go v1.1.0/go.mod h1:UkNAQd3GIcIGf0SeVgPpRdFStlNbqXla1AfSYxPUl2o=
 github.com/pascaldekloe/goe v0.0.0-20180627143212-57f6aae5913c/go.mod h1:lzWF7FIEvWOWxwDKqyGYQf6ZUaNfKdP144TG7ZOy1lc=
 github.com/pascaldekloe/goe v0.1.0 h1:cBOtyMzM9HTpWjXfbbunk26uA6nG3a8n06Wieeh0MwY=
@@ -275,6 +259,8 @@ github.com/rabbitmq/amqp091-go v1.10.0 h1:STpn5XsHlHGcecLmMFCtg7mqq0RnD+zFr4uzuk
 github.com/rabbitmq/amqp091-go v1.10.0/go.mod h1:Hy4jKW5kQART1u+JkDTF9YYOQUHXqMuhrgxOEeS7G4o=
 github.com/redis/go-redis/v9 v9.7.0 h1:HhLSs+B6O021gwzl+locl0zEDnyNkxMtf/Z3NNBMa9E=
 github.com/redis/go-redis/v9 v9.7.0/go.mod h1:f6zhXITC7JUJIlPEiBOTXxJgPLdZcA93GewI7inzyWw=
+github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5 h1:mZHayPoR0lNmnHyvtYjDeq0zlVHn9K/ZXoy17ylucdo=
+github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5/go.mod h1:GEXHk5HgEKCvEIIrSpFI3ozzG5xOKA2DVlEX/gGnewM=
 github.com/rogpeppe/go-internal v1.9.0 h1:73kH8U+JUqXU8lRuOHeVHaa/SZPifC7BkcraZVejAe8=
 github.com/rogpeppe/go-internal v1.9.0/go.mod h1:WtVeX8xhTBvf0smdhujwtBcq4Qrzq/fJaraNFVN+nFs=
 github.com/ryanuber/columnize v0.0.0-20160712163229-9b3edd62028f/go.mod h1:sm1tb6uqfes/u+d4ooFouqFdy9/2g9QGwK3SQygK0Ts=
@@ -307,7 +293,6 @@ github.com/stretchr/objx v0.5.0/go.mod h1:Yh+to48EsGEfYuaHDzXPcE3xhTkx73EhmCGUpE
 github.com/stretchr/testify v1.2.2/go.mod h1:a8OnRcib4nhh0OaRAV+Yts87kKdq0PP7pXfy6kDkUVs=
 github.com/stretchr/testify v1.3.0/go.mod h1:M5WIy9Dh21IEIfnGCwXGc5bZfKNJtfHm1UVUgZn+9EI=
 github.com/stretchr/testify v1.4.0/go.mod h1:j7eGeouHqKxXV5pUuKE4zz7dFj8WfuZ+81PSLYec5m4=
-github.com/stretchr/testify v1.5.1/go.mod h1:5W2xD1RspED5o8YsWQXVCued0rvSQ+mT+I5cxcmMvtA=
 github.com/stretchr/testify v1.7.0/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
 github.com/stretchr/testify v1.7.1/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
 github.com/stretchr/testify v1.7.2/go.mod h1:R6va5+xMeoiuVRoj+gSkQ7d3FALtqAAGI1FQKckRals=
@@ -322,8 +307,6 @@ github.com/twitchyliquid64/golang-asm v0.15.1 h1:SU5vSMR7hnwNxj24w34ZyCi/FmDZTkS
 github.com/twitchyliquid64/golang-asm v0.15.1/go.mod h1:a1lVb/DtPvCB8fslRZhAngC2+aY1QWCk3Cedj/Gdt08=
 github.com/ugorji/go/codec v1.2.12 h1:9LC83zGrHhuUA9l16C9AHXAqEV/2wBQ4nkvumAE65EE=
 github.com/ugorji/go/codec v1.2.12/go.mod h1:UNopzCgEMSXjBc6AOMqYvWC1ktqTAfzJZUZgYf6w6lg=
-github.com/x-cray/logrus-prefixed-formatter v0.5.2 h1:00txxvfBM9muc0jiLIEAkAcIMJzfthRT6usrui8uGmg=
-github.com/x-cray/logrus-prefixed-formatter v0.5.2/go.mod h1:2duySbKsL6M18s5GU7VPsoEPHyzalCE06qoARUCeBBE=
 github.com/yuin/goldmark v1.1.27/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
 github.com/yuin/goldmark v1.2.1/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
 go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0 h1:0nTRpaCaILLdooXAQnfktlL6Zw1ECKEW9DZGH2byi2c=
@@ -372,7 +355,6 @@ golang.org/x/mod v0.2.0/go.mod h1:s0Qsj1ACt9ePp/hMypM3fl4fZqREWJwdYDEqhRiZZUA=
 golang.org/x/mod v0.3.0/go.mod h1:s0Qsj1ACt9ePp/hMypM3fl4fZqREWJwdYDEqhRiZZUA=
 golang.org/x/net v0.0.0-20180724234803-3673e40ba225/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
 golang.org/x/net v0.0.0-20180826012351-8a410e7b638d/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
-golang.org/x/net v0.0.0-20180906233101-161cd47e91fd/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
 golang.org/x/net v0.0.0-20181114220301-adae6a3d119a/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
 golang.org/x/net v0.0.0-20190213061140-3a22650c66bd/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
 golang.org/x/net v0.0.0-20190311183353-d8887717615a/go.mod h1:t9HGtf8HONx5eT2rtn7q6eTqICYqUVnKs3thJo3Qplg=
@@ -381,7 +363,6 @@ golang.org/x/net v0.0.0-20190613194153-d28f0bde5980/go.mod h1:z5CRVTTTmAJ677TzLL
 golang.org/x/net v0.0.0-20190620200207-3b0461eec859/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
 golang.org/x/net v0.0.0-20190923162816-aa69164e4478/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
 golang.org/x/net v0.0.0-20200226121028-0de0cce0169b/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
-golang.org/x/net v0.0.0-20200520004742-59133d7f0dd7/go.mod h1:qpuaurCH72eLCgpAm/N6yyVIVM9cpaDIP3A8BGJEC5A=
 golang.org/x/net v0.0.0-20201021035429-f5854403a974/go.mod h1:sp8m0HH+o8qH0wwXwYZr8TS3Oi6o0r6Gce1SSxlDquU=
 golang.org/x/net v0.0.0-20210226172049-e18ecbb05110/go.mod h1:m0MpNAwzfU5UDzcl9v0D8zg8gWTRqZa9RBIspLL5mdg=
 golang.org/x/net v0.0.0-20210410081132-afb366fc7cd1/go.mod h1:9tjilg8BloeKEkVJvy7fQ90B1CfIiPueXVOjqfkSzI8=
@@ -398,25 +379,19 @@ golang.org/x/sync v0.0.0-20210220032951-036812b2e83c/go.mod h1:RxMgew5VJxzue5/jJ
 golang.org/x/sys v0.0.0-20180823144017-11551d06cbcc/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
 golang.org/x/sys v0.0.0-20180830151530-49385e6e1522/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
 golang.org/x/sys v0.0.0-20180905080454-ebe1bf3edb33/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
-golang.org/x/sys v0.0.0-20180909124046-d0be0721c37e/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
 golang.org/x/sys v0.0.0-20181116152217-5ac8a444bdc5/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
 golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
 golang.org/x/sys v0.0.0-20190222072716-a9d3bda3a223/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
 golang.org/x/sys v0.0.0-20190412213103-97732733099d/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20190422165155-953cdadca894/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
-golang.org/x/sys v0.0.0-20190904154756-749cb33beabd/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20190922100055-0a153f010e69/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20190924154521-2837fb4f24fe/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
-golang.org/x/sys v0.0.0-20191005200804-aed5e4c7ecf9/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20191026070338-33540a1f6037/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
-golang.org/x/sys v0.0.0-20191120155948-bd437916bb0e/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20200116001909-b77594299b42/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20200122134326-e047566fdf82/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20200223170610-d5e6a3e2c0ae/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
-golang.org/x/sys v0.0.0-20200323222414-85ca7c5b95cd/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20200930185726-fdedc70b468f/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20201119102817-f84b799fce68/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
-golang.org/x/sys v0.0.0-20210112080510-489259a85091/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20210303074136-134d130e1a04/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20210330210617-4fbd30eecc44/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
@@ -430,8 +405,6 @@ golang.org/x/sys v0.6.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.28.0 h1:Fksou7UEQUWlKvIdsqzJmUmCX3cZuD2+P3XyyzwMhlA=
 golang.org/x/sys v0.28.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
 golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1/go.mod h1:bj7SfCRtBDWHUb9snDiAeCFNEtKQo2Wmx5Cou7ajbmo=
-golang.org/x/term v0.27.0 h1:WP60Sv1nlK1T6SupCHbXzSaN0b9wUmsPoRS9b61A23Q=
-golang.org/x/term v0.27.0/go.mod h1:iMsnZpn0cago0GOrHO2+Y7u7JPn5AylBrcoWkElMTSM=
 golang.org/x/text v0.3.0/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
 golang.org/x/text v0.3.2/go.mod h1:bEr9sfX3Q8Zfm5fL9x+3itogRgK3+ptLWKqgva+5dAk=
 golang.org/x/text v0.3.3/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
@@ -447,7 +420,6 @@ golang.org/x/tools v0.0.0-20190907020128-2ca718005c18/go.mod h1:b+2E5dAYhXwXZwtn
 golang.org/x/tools v0.0.0-20191108193012-7d206e10da11/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
 golang.org/x/tools v0.0.0-20191119224855-298f0cb1881e/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
 golang.org/x/tools v0.0.0-20200619180055-7c47624df98f/go.mod h1:EkVYQZoAsY45+roYkvgYkIh4xh/qjgUK9TdY2XT94GE=
-golang.org/x/tools v0.0.0-20201224043029-2b0845dc783e/go.mod h1:emZCQorbCU4vsT4fOWvOPXz4eW1wZW4PmDk9uLelYpA=
 golang.org/x/tools v0.0.0-20210106214847-113979e3529a/go.mod h1:emZCQorbCU4vsT4fOWvOPXz4eW1wZW4PmDk9uLelYpA=
 golang.org/x/xerrors v0.0.0-20190717185122-a985d3407aa7/go.mod h1:I/5z698sn9Ka8TeJc9MKroUUfqBBauWjQqLJ2OPfmY0=
 golang.org/x/xerrors v0.0.0-20191011141410-1b5146add898/go.mod h1:I/5z698sn9Ka8TeJc9MKroUUfqBBauWjQqLJ2OPfmY0=
@@ -467,12 +439,6 @@ google.golang.org/grpc v1.27.0/go.mod h1:qbnxyOmOxrQa7FizSgH+ReBfzJrCY1pSN7KXBS8
 google.golang.org/grpc v1.29.1/go.mod h1:itym6AZVZYACWQqET3MqgPpjcuV5QH3BxFS3IjizoKk=
 google.golang.org/grpc v1.67.1 h1:zWnc1Vrcno+lHZCOofnIMvycFcc0QRGIzm9dhnDX68E=
 google.golang.org/grpc v1.67.1/go.mod h1:1gLDyUQU7CTLJI90u3nXZ9ekeghjeM7pTDZlqFNg2AA=
-google.golang.org/protobuf v0.0.0-20200109180630-ec00e32a8dfd/go.mod h1:DFci5gLYBciE7Vtevhsrf46CRTquxDuWsQurQQe4oz8=
-google.golang.org/protobuf v0.0.0-20200221191635-4d8936d0db64/go.mod h1:kwYJMbMJ01Woi6D6+Kah6886xMZcty6N08ah7+eCXa0=
-google.golang.org/protobuf v0.0.0-20200228230310-ab0ca4ff8a60/go.mod h1:cfTl7dwQJ+fmap5saPgwCLgHXTUD7jkjRqWcaiX5VyM=
-google.golang.org/protobuf v1.20.1-0.20200309200217-e05f789c0967/go.mod h1:A+miEFZTKqfCUM6K7xSMQL9OKL/b6hQv+e19PK+JZNE=
-google.golang.org/protobuf v1.21.0/go.mod h1:47Nbq4nVaFHyn7ilMalzfO3qCViNmqZ2kzikPIcrTAo=
-google.golang.org/protobuf v1.23.0/go.mod h1:EGpADcykh3NcUnDUJcl1+ZksZNG86OlYog2l/sGQquU=
 google.golang.org/protobuf v1.26.0-rc.1/go.mod h1:jlhhOSvTdKEhbULTjvd4ARK9grFBp09yW+WbY/TyQbw=
 google.golang.org/protobuf v1.26.0/go.mod h1:9q0QmTI4eRPtz6boOQmLYwt+qCgq0jsYwAQnmE0givc=
 google.golang.org/protobuf v1.36.1 h1:yBPeRvTftaleIgM3PZ/WBIZ7XM/eEYAaEyCwvyjq/gk=
@@ -482,17 +448,13 @@ gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8
 gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
 gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 h1:YR8cESwS4TdDjEe65xsg0ogRM/Nc3DYOhEAlW+xobZo=
 gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
-gopkg.in/fsnotify.v1 v1.4.7/go.mod h1:Tz8NjZHkW78fSQdbUxIjBTcgA1z1m8ZHf0WmKUhAMys=
 gopkg.in/ini.v1 v1.67.0 h1:Dgnx+6+nfE+IfzjUEISNeydPJh9AXNNsWbGP9KzCsOA=
 gopkg.in/ini.v1 v1.67.0/go.mod h1:pNLf8WUiyNEtQjuu5G5vTm06TEv9tsIgeAvK8hOrP4k=
-gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 h1:uRGJdciOHaEIrze2W8Q3AKkepLTh2hOroT7a+7czfdQ=
-gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7/go.mod h1:dt/ZhP58zS4L8KSrWDmTeBkI65Dw0HsyUHuEVlX15mw=
 gopkg.in/yaml.v2 v2.2.1/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.2/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.4/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.5/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.8/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
-gopkg.in/yaml.v2 v2.3.0/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
 gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
 gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
diff --git a/internal/common/logging/logrus.go b/internal/common/logging/logrus.go
index 38eadf2..93a5b51 100644
--- a/internal/common/logging/logrus.go
+++ b/internal/common/logging/logrus.go
@@ -7,8 +7,9 @@ import (
 	"time"
 
 	"github.com/ghost-yu/go_shop_second/common/tracing"
+	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
+	"github.com/rifflock/lfshook"
 	"github.com/sirupsen/logrus"
-	prefixed "github.com/x-cray/logrus-prefixed-formatter"
 )
 
 // 要么用logging.Infof, Warnf...
@@ -17,9 +18,57 @@ import (
 func Init() {
 	SetFormatter(logrus.StandardLogger())
 	logrus.SetLevel(logrus.DebugLevel)
+	setOutput(logrus.StandardLogger())
 	logrus.AddHook(&traceHook{})
 }
 
+func setOutput(logger *logrus.Logger) {
+	var (
+		folder    = "./log/"
+		filePath  = "app.log"
+		errorPath = "errors.log"
+	)
+	if err := os.MkdirAll(folder, 0750); err != nil && !os.IsExist(err) {
+		panic(err)
+	}
+	file, err := os.OpenFile(folder+filePath, os.O_CREATE|os.O_RDWR, 0755)
+	if err != nil {
+		panic(err)
+	}
+	_, err = os.OpenFile(folder+errorPath, os.O_CREATE|os.O_RDWR, 0755)
+	if err != nil {
+		panic(err)
+	}
+	logger.SetOutput(file)
+
+	rotateInfo, err := rotatelogs.New(
+		folder+filePath+".%Y%m%d",
+		rotatelogs.WithLinkName("app.log"),
+		rotatelogs.WithMaxAge(7*24*time.Hour),
+		rotatelogs.WithRotationTime(1*time.Hour),
+	)
+	if err != nil {
+		panic(err)
+	}
+	rotateError, err := rotatelogs.New(
+		folder+errorPath+".%Y%m%d",
+		rotatelogs.WithLinkName("errors.log"),
+		rotatelogs.WithMaxAge(7*24*time.Hour),
+		rotatelogs.WithRotationTime(1*time.Hour),
+	)
+	rotationMap := lfshook.WriterMap{
+		logrus.DebugLevel: rotateInfo,
+		logrus.InfoLevel:  rotateInfo,
+		logrus.WarnLevel:  rotateError,
+		logrus.ErrorLevel: rotateError,
+		logrus.FatalLevel: rotateError,
+		logrus.PanicLevel: rotateError,
+	}
+	logrus.AddHook(lfshook.NewHook(rotationMap, &logrus.JSONFormatter{
+		TimestampFormat: time.DateTime,
+	}))
+}
+
 func SetFormatter(logger *logrus.Logger) {
 	logger.SetFormatter(&logrus.JSONFormatter{
 		TimestampFormat: time.RFC3339,
@@ -30,11 +79,11 @@ func SetFormatter(logger *logrus.Logger) {
 		},
 	})
 	if isLocal, _ := strconv.ParseBool(os.Getenv("LOCAL_ENV")); isLocal {
-		logger.SetFormatter(&prefixed.TextFormatter{
-			ForceColors:     true,
-			ForceFormatting: true,
-			TimestampFormat: time.RFC3339,
-		})
+		//logger.SetFormatter(&prefixed.TextFormatter{
+		//	ForceColors:     true,
+		//	ForceFormatting: true,
+		//	TimestampFormat: time.RFC3339,
+		//})
 	}
 }
 
diff --git a/internal/kitchen/go.mod b/internal/kitchen/go.mod
index 22b3e73..e13cc78 100644
--- a/internal/kitchen/go.mod
+++ b/internal/kitchen/go.mod
@@ -14,6 +14,7 @@ require (
 )
 
 require (
+	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
 	github.com/armon/go-metrics v0.4.1 // indirect
 	github.com/fatih/color v1.14.1 // indirect
 	github.com/fsnotify/fsnotify v1.7.0 // indirect
@@ -30,14 +31,16 @@ require (
 	github.com/hashicorp/golang-lru v0.5.4 // indirect
 	github.com/hashicorp/hcl v1.0.0 // indirect
 	github.com/hashicorp/serf v0.10.1 // indirect
+	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible // indirect
+	github.com/lestrrat-go/strftime v1.1.0 // indirect
 	github.com/magiconair/properties v1.8.7 // indirect
 	github.com/mattn/go-colorable v0.1.13 // indirect
 	github.com/mattn/go-isatty v0.0.20 // indirect
-	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
 	github.com/mitchellh/go-homedir v1.1.0 // indirect
 	github.com/mitchellh/mapstructure v1.5.0 // indirect
-	github.com/nxadm/tail v1.4.11 // indirect
+	github.com/oapi-codegen/runtime v1.1.1 // indirect
 	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
+	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5 // indirect
 	github.com/rogpeppe/go-internal v1.10.0 // indirect
 	github.com/sagikazarmark/locafero v0.4.0 // indirect
 	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
@@ -46,7 +49,6 @@ require (
 	github.com/spf13/cast v1.6.0 // indirect
 	github.com/spf13/pflag v1.0.5 // indirect
 	github.com/subosito/gotenv v1.6.0 // indirect
-	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
 	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0 // indirect
 	go.opentelemetry.io/contrib/propagators/b3 v1.31.0 // indirect
 	go.opentelemetry.io/otel/exporters/jaeger v1.17.0 // indirect
@@ -55,11 +57,9 @@ require (
 	go.opentelemetry.io/otel/trace v1.31.0 // indirect
 	go.uber.org/atomic v1.9.0 // indirect
 	go.uber.org/multierr v1.9.0 // indirect
-	golang.org/x/crypto v0.31.0 // indirect
 	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
 	golang.org/x/net v0.33.0 // indirect
 	golang.org/x/sys v0.28.0 // indirect
-	golang.org/x/term v0.27.0 // indirect
 	golang.org/x/text v0.21.0 // indirect
 	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
 	google.golang.org/grpc v1.67.1 // indirect
diff --git a/internal/kitchen/go.sum b/internal/kitchen/go.sum
index 1528c76..21e8c30 100644
--- a/internal/kitchen/go.sum
+++ b/internal/kitchen/go.sum
@@ -1,8 +1,11 @@
 github.com/DataDog/datadog-go v3.2.0+incompatible/go.mod h1:LButxg5PwREeZtORoXG3tL4fMGNddJ+vMq1mwgfaqoQ=
+github.com/RaveNoX/go-jsoncommentstrip v1.0.0/go.mod h1:78ihd09MekBnJnxpICcwzCMzGrKSKYe4AqU6PDYYpjk=
 github.com/alecthomas/template v0.0.0-20160405071501-a0175ee3bccc/go.mod h1:LOuyumcjzFXgccqObfd/Ljyb9UuFJ6TxHnclSeseNhc=
 github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751/go.mod h1:LOuyumcjzFXgccqObfd/Ljyb9UuFJ6TxHnclSeseNhc=
 github.com/alecthomas/units v0.0.0-20151022065526-2efee857e7cf/go.mod h1:ybxpYRFXyAe+OPACYpWeL0wqObRcbAqCMya13uyzqw0=
 github.com/alecthomas/units v0.0.0-20190717042225-c3de453c63f4/go.mod h1:ybxpYRFXyAe+OPACYpWeL0wqObRcbAqCMya13uyzqw0=
+github.com/apapsch/go-jsonmerge/v2 v2.0.0 h1:axGnT1gRIfimI7gJifB699GoE/oq+F2MU7Dml6nw9rQ=
+github.com/apapsch/go-jsonmerge/v2 v2.0.0/go.mod h1:lvDnEdqiQrp0O42VQGgmlKpxL1AP2+08jFMw88y4klk=
 github.com/armon/circbuf v0.0.0-20150827004946-bbbad097214e/go.mod h1:3U/XgcO3hCbHZ8TKRvWD2dDTCfh9M9ya+I9JpbB7O8o=
 github.com/armon/go-metrics v0.0.0-20180917152333-f0300d1749da/go.mod h1:Q73ZrmVTwzkszR9V5SSuryQ31EELlFMUz1kKyl939pY=
 github.com/armon/go-metrics v0.4.1 h1:hR91U9KYmb6bLBYLQjyM+3j+rcd/UhE+G78SFnF8gJA=
@@ -13,6 +16,7 @@ github.com/beorn7/perks v0.0.0-20180321164747-3a771d992973/go.mod h1:Dwedo/Wpr24
 github.com/beorn7/perks v1.0.0/go.mod h1:KWe93zE9D1o94FZ5RNwFwVgaQK1VOXiVxmqh+CedLV8=
 github.com/beorn7/perks v1.0.1/go.mod h1:G2ZrVWU2WbWT9wwq4/hrbKbnv/1ERSJQ0ibhJ6rlkpw=
 github.com/bgentry/speakeasy v0.1.0/go.mod h1:+zsyZBPWlz7T6j88CTgSN5bM796AkVf0kBD4zp0CCIs=
+github.com/bmatcuk/doublestar v1.1.1/go.mod h1:UD6OnuiIn0yFxxA2le/rnRU1G4RaI4UvFv1sNto9p6w=
 github.com/cespare/xxhash/v2 v2.1.1/go.mod h1:VGX0DQ3Q6kWi7AoAeZDth3/j3BFtOZR5XLFGgcrjCOs=
 github.com/circonus-labs/circonus-gometrics v2.3.1+incompatible/go.mod h1:nmEj6Dob7S7YxXgwXpfOuvO54S+tGdZdw9fuRZt25Ag=
 github.com/circonus-labs/circonusllhist v0.1.3/go.mod h1:kMXHVDlOchFAehlya5ePtbp5jckzBHf4XRpQvBOLI+I=
@@ -27,7 +31,6 @@ github.com/fatih/color v1.14.1 h1:qfhVLaG5s+nCROl1zJsZRxFeYrHLqWroPOQ8BWiNb4w=
 github.com/fatih/color v1.14.1/go.mod h1:2oHN61fhTpgcxD3TSWCgKDiH1+x4OiDVVGH8WlgGZGg=
 github.com/frankban/quicktest v1.14.6 h1:7Xjx+VpznH+oBnejlPUj8oUpdxnVs4f8XU8WnHkI4W8=
 github.com/frankban/quicktest v1.14.6/go.mod h1:4ptaffx2x8+WTWXmUCuVU6aPUX1/Mz7zb5vbUoiM6w0=
-github.com/fsnotify/fsnotify v1.6.0/go.mod h1:sl3t1tCWJFWoRz9R8WJCbQihKKwmorjAbSClcnxKAGw=
 github.com/fsnotify/fsnotify v1.7.0 h1:8JEhPFa5W2WU7YfeZzPNqzMP6Lwt7L2715Ggo0nosvA=
 github.com/fsnotify/fsnotify v1.7.0/go.mod h1:40Bi/Hjc2AVfZrqy+aj+yEI+/bRxZnMJyTJwOpGvigM=
 github.com/go-kit/kit v0.8.0/go.mod h1:xBxKIO96dXMWWy0MnWVtmwkA9/13aqxPnvrjFYMA2as=
@@ -100,8 +103,11 @@ github.com/hashicorp/memberlist v0.5.0 h1:EtYPN8DpAURiapus508I4n9CzHs2W+8NZGbmmR
 github.com/hashicorp/memberlist v0.5.0/go.mod h1:yvyXLpo0QaGE59Y7hDTsTzDD25JYBZ4mHgHUZ8lrOI0=
 github.com/hashicorp/serf v0.10.1 h1:Z1H2J60yRKvfDYAOZLd2MU0ND4AH/WDz7xYHDWQsIPY=
 github.com/hashicorp/serf v0.10.1/go.mod h1:yL2t6BqATOLGc5HF7qbFkTfXoPIY0WZdWHfEvMqbG+4=
+github.com/jonboulle/clockwork v0.5.0 h1:Hyh9A8u51kptdkR+cqRpT1EebBwTn1oK9YfGYbdFz6I=
+github.com/jonboulle/clockwork v0.5.0/go.mod h1:3mZlmanh0g2NDKO5TWZVJAfofYk64M7XN3SzBPjZF60=
 github.com/json-iterator/go v1.1.6/go.mod h1:+SdeFBvtyEkXs7REEP0seUULqWtbJapLOCVDaaPEHmU=
 github.com/json-iterator/go v1.1.9/go.mod h1:KdQUCv79m/52Kvf8AW2vK1V8akMuk1QjK/uOdHXbAo4=
+github.com/juju/gnuflag v0.0.0-20171113085948-2ce1bb71843d/go.mod h1:2PavIy+JPciBPrBUjwbNvtwB6RQlve+hkpll6QSNmOE=
 github.com/julienschmidt/httprouter v1.2.0/go.mod h1:SYymIcj16QtmaHHD7aYtjjsJG7VTCxuUUipMqKk8s4w=
 github.com/konsorten/go-windows-terminal-sequences v1.0.1/go.mod h1:T0+1ngSBFLxvqU3pZ+m/2kptfBszLMUkC4ZK/EgS/cQ=
 github.com/kr/logfmt v0.0.0-20140226030751-b84e30acd515/go.mod h1:+0opPa2QZZtGFBFZlji/RkVcI2GknAs/DXo4wKdlNEc=
@@ -113,6 +119,12 @@ github.com/kr/pty v1.1.1/go.mod h1:pFQYn66WHrOpPYNljwOMqo10TkYh1fy3cYio2l3bCsQ=
 github.com/kr/text v0.1.0/go.mod h1:4Jbv+DJW3UT/LiOwJeYQe1efqtUx/iVham/4vfdArNI=
 github.com/kr/text v0.2.0 h1:5Nx0Ya0ZqY2ygV366QzturHI13Jq95ApcVaJBhpS+AY=
 github.com/kr/text v0.2.0/go.mod h1:eLer722TekiGuMkidMxC/pM04lWEeraHUUmBw8l2grE=
+github.com/lestrrat-go/envload v0.0.0-20180220234015-a3eb8ddeffcc h1:RKf14vYWi2ttpEmkA4aQ3j4u9dStX2t4M8UM6qqNsG8=
+github.com/lestrrat-go/envload v0.0.0-20180220234015-a3eb8ddeffcc/go.mod h1:kopuH9ugFRkIXf3YoqHKyrJ9YfUFsckUU9S7B+XP+is=
+github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible h1:Y6sqxHMyB1D2YSzWkLibYKgg+SwmyFU9dF2hn6MdTj4=
+github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible/go.mod h1:ZQnN8lSECaebrkQytbHj4xNgtg8CR7RYXnPok8e0EHA=
+github.com/lestrrat-go/strftime v1.1.0 h1:gMESpZy44/4pXLO/m+sL0yBd1W6LjgjrrD4a68Gapyg=
+github.com/lestrrat-go/strftime v1.1.0/go.mod h1:uzeIB52CeUJenCo1syghlugshMysrqUT51HlxphXVeI=
 github.com/magiconair/properties v1.8.7 h1:IeQXZAiQcpL9mgcAe1Nu6cX9LLw6ExEHKjN0VQdvPDY=
 github.com/magiconair/properties v1.8.7/go.mod h1:Dhd985XPs7jluiymwWYZ0G4Z61jb3vdS329zhj2hYo0=
 github.com/mattn/go-colorable v0.0.9/go.mod h1:9vuHe8Xs5qXnSaW/c/ABM9alt+Vo+STaOChaDxuIBZU=
@@ -131,8 +143,6 @@ github.com/mattn/go-isatty v0.0.16/go.mod h1:kYGgaQfpe5nmfYZH+SKPsOc2e4SrIfOl2e/
 github.com/mattn/go-isatty v0.0.20 h1:xfD0iDuEKnDkl03q4limB+vH+GxLEtL/jb4xVJSWWEY=
 github.com/mattn/go-isatty v0.0.20/go.mod h1:W+V8PltTTMOvKvAeJH7IuucS94S2C6jfK/D7dTCTo3Y=
 github.com/matttproud/golang_protobuf_extensions v1.0.1/go.mod h1:D8He9yQNgCq6Z5Ld7szi9bcBfOoFv/3dc6xSMkL2PC0=
-github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d h1:5PJl274Y63IEHC+7izoQE9x6ikvDFZS2mDVS3drnohI=
-github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d/go.mod h1:01TrycV0kFyexm33Z7vhZRXopbI8J3TDReVlkTgMUxE=
 github.com/miekg/dns v1.1.26/go.mod h1:bPDLeHnStXmXAq1m/Ch/hvfNHr14JKNPMBo3VZKjuso=
 github.com/miekg/dns v1.1.41 h1:WMszZWJG0XmzbK9FEmzH2TVcqYzFesusSIB41b8KHxY=
 github.com/miekg/dns v1.1.41/go.mod h1:p6aan82bvRIyn+zDIv9xYNUpwa73JcSh9BKwknJysuI=
@@ -147,12 +157,8 @@ github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd/go.mod h1:6dJ
 github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742/go.mod h1:bx2lNnkwVCuqBIxFjflWJWanXIb3RllmbCylyMrvgv0=
 github.com/modern-go/reflect2 v1.0.1/go.mod h1:bx2lNnkwVCuqBIxFjflWJWanXIb3RllmbCylyMrvgv0=
 github.com/mwitkow/go-conntrack v0.0.0-20161129095857-cc309e4a2223/go.mod h1:qRWi+5nqEBWmkhHvq77mSJWrCKwh8bxhgT7d/eI7P4U=
-github.com/nxadm/tail v1.4.11 h1:8feyoE3OzPrcshW5/MJ4sGESc5cqmGkGCWlco4l0bqY=
-github.com/nxadm/tail v1.4.11/go.mod h1:OTaG3NK980DZzxbRq6lEuzgU+mug70nY11sMd4JXXHc=
-github.com/onsi/ginkgo v1.16.5 h1:8xi0RTUf59SOSfEtZMvwTvXYMzG4gV23XVHOZiXNtnE=
-github.com/onsi/ginkgo v1.16.5/go.mod h1:+E8gABHa3K6zRBolWtd+ROzc/U5bkGt0FwiG042wbpU=
-github.com/onsi/gomega v1.36.2 h1:koNYke6TVk6ZmnyHrCXba/T/MoLBXFjeC1PtvYgw0A8=
-github.com/onsi/gomega v1.36.2/go.mod h1:DdwyADRjrc825LhMEkD76cHR5+pUnjhUN8GlHlRPHzY=
+github.com/oapi-codegen/runtime v1.1.1 h1:EXLHh0DXIJnWhdRPN2w4MXAzFyE4CskzhNLUmtpMYro=
+github.com/oapi-codegen/runtime v1.1.1/go.mod h1:SK9X900oXmPWilYR5/WKPzt3Kqxn/uS/+lbpREv+eCg=
 github.com/pascaldekloe/goe v0.0.0-20180627143212-57f6aae5913c/go.mod h1:lzWF7FIEvWOWxwDKqyGYQf6ZUaNfKdP144TG7ZOy1lc=
 github.com/pascaldekloe/goe v0.1.0 h1:cBOtyMzM9HTpWjXfbbunk26uA6nG3a8n06Wieeh0MwY=
 github.com/pascaldekloe/goe v0.1.0/go.mod h1:lzWF7FIEvWOWxwDKqyGYQf6ZUaNfKdP144TG7ZOy1lc=
@@ -180,6 +186,8 @@ github.com/prometheus/procfs v0.0.2/go.mod h1:TjEm7ze935MbeOT/UhFTIMYKhuLP4wbCsT
 github.com/prometheus/procfs v0.0.8/go.mod h1:7Qr8sr6344vo1JqZ6HhLceV9o3AJ1Ff+GxbHq6oeK9A=
 github.com/rabbitmq/amqp091-go v1.10.0 h1:STpn5XsHlHGcecLmMFCtg7mqq0RnD+zFr4uzukfVhBw=
 github.com/rabbitmq/amqp091-go v1.10.0/go.mod h1:Hy4jKW5kQART1u+JkDTF9YYOQUHXqMuhrgxOEeS7G4o=
+github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5 h1:mZHayPoR0lNmnHyvtYjDeq0zlVHn9K/ZXoy17ylucdo=
+github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5/go.mod h1:GEXHk5HgEKCvEIIrSpFI3ozzG5xOKA2DVlEX/gGnewM=
 github.com/rogpeppe/go-internal v1.10.0 h1:TMyTOH3F/DB16zRVcYyreMH6GnZZrwQVAoYjRBZyWFQ=
 github.com/rogpeppe/go-internal v1.10.0/go.mod h1:UQnix2H7Ngw/k4C5ijL5+65zddjncjaFoBhdsK/akog=
 github.com/ryanuber/columnize v0.0.0-20160712163229-9b3edd62028f/go.mod h1:sm1tb6uqfes/u+d4ooFouqFdy9/2g9QGwK3SQygK0Ts=
@@ -203,6 +211,7 @@ github.com/spf13/pflag v1.0.5 h1:iy+VFUOCP1a+8yFto/drg2CJ5u0yRoB7fZw3DKv/JXA=
 github.com/spf13/pflag v1.0.5/go.mod h1:McXfInJRrz4CZXVZOBLb0bTZqETkiAhM9Iw0y3An2Bg=
 github.com/spf13/viper v1.19.0 h1:RWq5SEjt8o25SROyN3z2OrDB9l7RPd3lwTWU8EcEdcI=
 github.com/spf13/viper v1.19.0/go.mod h1:GQUN9bilAbhU/jgc1bKs99f/suXKeUMct8Adx5+Ntkg=
+github.com/spkg/bom v0.0.0-20160624110644-59b7046e48ad/go.mod h1:qLr4V1qq6nMqFKkMo8ZTx3f+BZEkzsRUY10Xsm2mwU0=
 github.com/stretchr/objx v0.1.0/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
 github.com/stretchr/objx v0.1.1/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
 github.com/stretchr/objx v0.5.0 h1:1zr/of2m5FGMsad5YfcqgdqdWrIhu+EBEJRhR1U7z/c=
@@ -217,8 +226,6 @@ github.com/stretchr/testify v1.9.0/go.mod h1:r2ic/lqez/lEtzL7wO/rwa5dbSLXVDPFyf8
 github.com/subosito/gotenv v1.6.0 h1:9NlTDc1FTs4qu0DDq7AEtTPNw6SVm7uBMsUCUjABIf8=
 github.com/subosito/gotenv v1.6.0/go.mod h1:Dk4QP5c2W3ibzajGcXpNraDfq2IrhjMIvMSWPKKo0FU=
 github.com/tv42/httpunix v0.0.0-20150427012821-b75d8614f926/go.mod h1:9ESjWnEqriFuLhtthL60Sar/7RFoluCcXsuvEwTV5KM=
-github.com/x-cray/logrus-prefixed-formatter v0.5.2 h1:00txxvfBM9muc0jiLIEAkAcIMJzfthRT6usrui8uGmg=
-github.com/x-cray/logrus-prefixed-formatter v0.5.2/go.mod h1:2duySbKsL6M18s5GU7VPsoEPHyzalCE06qoARUCeBBE=
 go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0 h1:4Pp6oUg3+e/6M4C0A/3kJ2VYa++dsWVTtGgLVj5xtHg=
 go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0/go.mod h1:Mjt1i1INqiaoZOMGR1RIUJN+i3ChKoFRqzrRQhlkbs0=
 go.opentelemetry.io/contrib/propagators/b3 v1.31.0 h1:PQPXYscmwbCp76QDvO4hMngF2j8Bx/OTV86laEl8uqo=
@@ -242,8 +249,6 @@ go.uber.org/multierr v1.9.0/go.mod h1:X2jQV1h+kxSjClGpnseKVIxpmcjrj7MNnI0bnlfKTV
 golang.org/x/crypto v0.0.0-20180904163835-0709b304e793/go.mod h1:6SG95UA2DQfeDnfUPMdvaQW0Q7yPrPDi9nlGo2tz2b4=
 golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2/go.mod h1:djNgcEr1/C05ACkg1iLfiJU5Ep61QUkGW8qpdssI0+w=
 golang.org/x/crypto v0.0.0-20190923035154-9ee001bba392/go.mod h1:/lpIB1dKB+9EgE3H3cr1v9wB50oz8l4C4h62xy7jSTY=
-golang.org/x/crypto v0.31.0 h1:ihbySMvVjLAeSH1IbfcRTkD/iNscyz8rGzjF/E5hV6U=
-golang.org/x/crypto v0.31.0/go.mod h1:kDsLvtWBEx7MV9tJOj9bnXsPbxwJQ6csT/x4KIN4Ssk=
 golang.org/x/exp v0.0.0-20230905200255-921286631fa9 h1:GoHiUyI/Tp2nVkLI2mCxVkOjsbSXD66ic0XW0js0R9g=
 golang.org/x/exp v0.0.0-20230905200255-921286631fa9/go.mod h1:S2oDrQGGwySpoQPVqRShND87VCbxmc6bL1Yd2oYrm6k=
 golang.org/x/net v0.0.0-20181114220301-adae6a3d119a/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
@@ -281,13 +286,10 @@ golang.org/x/sys v0.0.0-20220503163025-988cb79eb6c6/go.mod h1:oPkhp1MJrh7nUepCBc
 golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20220728004956-3c1f35247d10/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
-golang.org/x/sys v0.0.0-20220908164124-27713097b956/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.6.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.28.0 h1:Fksou7UEQUWlKvIdsqzJmUmCX3cZuD2+P3XyyzwMhlA=
 golang.org/x/sys v0.28.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
 golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1/go.mod h1:bj7SfCRtBDWHUb9snDiAeCFNEtKQo2Wmx5Cou7ajbmo=
-golang.org/x/term v0.27.0 h1:WP60Sv1nlK1T6SupCHbXzSaN0b9wUmsPoRS9b61A23Q=
-golang.org/x/term v0.27.0/go.mod h1:iMsnZpn0cago0GOrHO2+Y7u7JPn5AylBrcoWkElMTSM=
 golang.org/x/text v0.3.0/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
 golang.org/x/text v0.3.2/go.mod h1:bEr9sfX3Q8Zfm5fL9x+3itogRgK3+ptLWKqgva+5dAk=
 golang.org/x/text v0.3.3/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
@@ -311,8 +313,6 @@ gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c h1:Hei/4ADfdWqJk1ZMxUNpqntN
 gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c/go.mod h1:JHkPIbrfpd72SG/EVd6muEfDQjcINNoR0C8j2r3qZ4Q=
 gopkg.in/ini.v1 v1.67.0 h1:Dgnx+6+nfE+IfzjUEISNeydPJh9AXNNsWbGP9KzCsOA=
 gopkg.in/ini.v1 v1.67.0/go.mod h1:pNLf8WUiyNEtQjuu5G5vTm06TEv9tsIgeAvK8hOrP4k=
-gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 h1:uRGJdciOHaEIrze2W8Q3AKkepLTh2hOroT7a+7czfdQ=
-gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7/go.mod h1:dt/ZhP58zS4L8KSrWDmTeBkI65Dw0HsyUHuEVlX15mw=
 gopkg.in/yaml.v2 v2.2.1/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.2/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.4/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
diff --git a/internal/order/app/command/create_order.go b/internal/order/app/command/create_order.go
index ac0d9b7..df86d4a 100644
--- a/internal/order/app/command/create_order.go
+++ b/internal/order/app/command/create_order.go
@@ -35,13 +35,7 @@ type createOrderHandler struct {
 	channel   *amqp.Channel
 }
 
-func NewCreateOrderHandler(
-	orderRepo domain.Repository,
-	stockGRPC query.StockService,
-	channel *amqp.Channel,
-	logger *logrus.Entry,
-	metricClient decorator.MetricsClient,
-) CreateOrderHandler {
+func NewCreateOrderHandler(orderRepo domain.Repository, stockGRPC query.StockService, channel *amqp.Channel, logger *logrus.Logger, metricClient decorator.MetricsClient) CreateOrderHandler {
 	if orderRepo == nil {
 		panic("nil orderRepo")
 	}
diff --git a/internal/order/app/command/update_order.go b/internal/order/app/command/update_order.go
index 8849f02..c757589 100644
--- a/internal/order/app/command/update_order.go
+++ b/internal/order/app/command/update_order.go
@@ -21,11 +21,7 @@ type updateOrderHandler struct {
 	//stockGRPC
 }
 
-func NewUpdateOrderHandler(
-	orderRepo domain.Repository,
-	logger *logrus.Entry,
-	metricClient decorator.MetricsClient,
-) UpdateOrderHandler {
+func NewUpdateOrderHandler(orderRepo domain.Repository, logger *logrus.Logger, metricClient decorator.MetricsClient) UpdateOrderHandler {
 	if orderRepo == nil {
 		panic("nil orderRepo")
 	}
diff --git a/internal/order/app/query/get_customer_order.go b/internal/order/app/query/get_customer_order.go
index b107e01..1b4dfd5 100644
--- a/internal/order/app/query/get_customer_order.go
+++ b/internal/order/app/query/get_customer_order.go
@@ -19,11 +19,7 @@ type getCustomerOrderHandler struct {
 	orderRepo domain.Repository
 }
 
-func NewGetCustomerOrderHandler(
-	orderRepo domain.Repository,
-	logger *logrus.Entry,
-	metricClient decorator.MetricsClient,
-) GetCustomerOrderHandler {
+func NewGetCustomerOrderHandler(orderRepo domain.Repository, logger *logrus.Logger, metricClient decorator.MetricsClient) GetCustomerOrderHandler {
 	if orderRepo == nil {
 		panic("nil orderRepo")
 	}
diff --git a/internal/order/go.mod b/internal/order/go.mod
index e4ad0ac..26080d8 100644
--- a/internal/order/go.mod
+++ b/internal/order/go.mod
@@ -56,18 +56,19 @@ require (
 	github.com/klauspost/compress v1.17.2 // indirect
 	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
 	github.com/leodido/go-urn v1.4.0 // indirect
+	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible // indirect
+	github.com/lestrrat-go/strftime v1.1.0 // indirect
 	github.com/magiconair/properties v1.8.7 // indirect
 	github.com/mattn/go-colorable v0.1.13 // indirect
 	github.com/mattn/go-isatty v0.0.20 // indirect
-	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
 	github.com/mitchellh/go-homedir v1.1.0 // indirect
 	github.com/mitchellh/mapstructure v1.5.0 // indirect
 	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
 	github.com/modern-go/reflect2 v1.0.2 // indirect
 	github.com/montanaflynn/stats v0.7.1 // indirect
-	github.com/nxadm/tail v1.4.11 // indirect
 	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
 	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
+	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5 // indirect
 	github.com/sagikazarmark/locafero v0.6.0 // indirect
 	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
 	github.com/sourcegraph/conc v0.3.0 // indirect
@@ -77,7 +78,6 @@ require (
 	github.com/subosito/gotenv v1.6.0 // indirect
 	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
 	github.com/ugorji/go/codec v1.2.12 // indirect
-	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
 	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
 	github.com/xdg-go/scram v1.1.2 // indirect
 	github.com/xdg-go/stringprep v1.0.4 // indirect
@@ -96,7 +96,6 @@ require (
 	golang.org/x/net v0.33.0 // indirect
 	golang.org/x/sync v0.10.0 // indirect
 	golang.org/x/sys v0.28.0 // indirect
-	golang.org/x/term v0.27.0 // indirect
 	golang.org/x/text v0.21.0 // indirect
 	google.golang.org/genproto/googleapis/rpc v0.0.0-20241007155032-5fefd90f89a9 // indirect
 	gopkg.in/ini.v1 v1.67.0 // indirect
diff --git a/internal/order/go.sum b/internal/order/go.sum
index d255009..67b69b4 100644
--- a/internal/order/go.sum
+++ b/internal/order/go.sum
@@ -50,7 +50,6 @@ github.com/fatih/color v1.14.1 h1:qfhVLaG5s+nCROl1zJsZRxFeYrHLqWroPOQ8BWiNb4w=
 github.com/fatih/color v1.14.1/go.mod h1:2oHN61fhTpgcxD3TSWCgKDiH1+x4OiDVVGH8WlgGZGg=
 github.com/frankban/quicktest v1.14.6 h1:7Xjx+VpznH+oBnejlPUj8oUpdxnVs4f8XU8WnHkI4W8=
 github.com/frankban/quicktest v1.14.6/go.mod h1:4ptaffx2x8+WTWXmUCuVU6aPUX1/Mz7zb5vbUoiM6w0=
-github.com/fsnotify/fsnotify v1.6.0/go.mod h1:sl3t1tCWJFWoRz9R8WJCbQihKKwmorjAbSClcnxKAGw=
 github.com/fsnotify/fsnotify v1.7.0 h1:8JEhPFa5W2WU7YfeZzPNqzMP6Lwt7L2715Ggo0nosvA=
 github.com/fsnotify/fsnotify v1.7.0/go.mod h1:40Bi/Hjc2AVfZrqy+aj+yEI+/bRxZnMJyTJwOpGvigM=
 github.com/gabriel-vasile/mimetype v1.4.6 h1:3+PzJTKLkvgjeTbts6msPJt4DixhT4YtFNf1gtGe3zc=
@@ -153,6 +152,8 @@ github.com/hashicorp/memberlist v0.5.0 h1:EtYPN8DpAURiapus508I4n9CzHs2W+8NZGbmmR
 github.com/hashicorp/memberlist v0.5.0/go.mod h1:yvyXLpo0QaGE59Y7hDTsTzDD25JYBZ4mHgHUZ8lrOI0=
 github.com/hashicorp/serf v0.10.1 h1:Z1H2J60yRKvfDYAOZLd2MU0ND4AH/WDz7xYHDWQsIPY=
 github.com/hashicorp/serf v0.10.1/go.mod h1:yL2t6BqATOLGc5HF7qbFkTfXoPIY0WZdWHfEvMqbG+4=
+github.com/jonboulle/clockwork v0.5.0 h1:Hyh9A8u51kptdkR+cqRpT1EebBwTn1oK9YfGYbdFz6I=
+github.com/jonboulle/clockwork v0.5.0/go.mod h1:3mZlmanh0g2NDKO5TWZVJAfofYk64M7XN3SzBPjZF60=
 github.com/json-iterator/go v1.1.6/go.mod h1:+SdeFBvtyEkXs7REEP0seUULqWtbJapLOCVDaaPEHmU=
 github.com/json-iterator/go v1.1.9/go.mod h1:KdQUCv79m/52Kvf8AW2vK1V8akMuk1QjK/uOdHXbAo4=
 github.com/json-iterator/go v1.1.12 h1:PV8peI4a0ysnczrg+LtxykD8LfKY9ML6u2jnxaEnrnM=
@@ -178,6 +179,12 @@ github.com/kr/text v0.2.0 h1:5Nx0Ya0ZqY2ygV366QzturHI13Jq95ApcVaJBhpS+AY=
 github.com/kr/text v0.2.0/go.mod h1:eLer722TekiGuMkidMxC/pM04lWEeraHUUmBw8l2grE=
 github.com/leodido/go-urn v1.4.0 h1:WT9HwE9SGECu3lg4d/dIA+jxlljEa1/ffXKmRjqdmIQ=
 github.com/leodido/go-urn v1.4.0/go.mod h1:bvxc+MVxLKB4z00jd1z+Dvzr47oO32F/QSNjSBOlFxI=
+github.com/lestrrat-go/envload v0.0.0-20180220234015-a3eb8ddeffcc h1:RKf14vYWi2ttpEmkA4aQ3j4u9dStX2t4M8UM6qqNsG8=
+github.com/lestrrat-go/envload v0.0.0-20180220234015-a3eb8ddeffcc/go.mod h1:kopuH9ugFRkIXf3YoqHKyrJ9YfUFsckUU9S7B+XP+is=
+github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible h1:Y6sqxHMyB1D2YSzWkLibYKgg+SwmyFU9dF2hn6MdTj4=
+github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible/go.mod h1:ZQnN8lSECaebrkQytbHj4xNgtg8CR7RYXnPok8e0EHA=
+github.com/lestrrat-go/strftime v1.1.0 h1:gMESpZy44/4pXLO/m+sL0yBd1W6LjgjrrD4a68Gapyg=
+github.com/lestrrat-go/strftime v1.1.0/go.mod h1:uzeIB52CeUJenCo1syghlugshMysrqUT51HlxphXVeI=
 github.com/magiconair/properties v1.8.7 h1:IeQXZAiQcpL9mgcAe1Nu6cX9LLw6ExEHKjN0VQdvPDY=
 github.com/magiconair/properties v1.8.7/go.mod h1:Dhd985XPs7jluiymwWYZ0G4Z61jb3vdS329zhj2hYo0=
 github.com/mattn/go-colorable v0.0.9/go.mod h1:9vuHe8Xs5qXnSaW/c/ABM9alt+Vo+STaOChaDxuIBZU=
@@ -196,8 +203,6 @@ github.com/mattn/go-isatty v0.0.16/go.mod h1:kYGgaQfpe5nmfYZH+SKPsOc2e4SrIfOl2e/
 github.com/mattn/go-isatty v0.0.20 h1:xfD0iDuEKnDkl03q4limB+vH+GxLEtL/jb4xVJSWWEY=
 github.com/mattn/go-isatty v0.0.20/go.mod h1:W+V8PltTTMOvKvAeJH7IuucS94S2C6jfK/D7dTCTo3Y=
 github.com/matttproud/golang_protobuf_extensions v1.0.1/go.mod h1:D8He9yQNgCq6Z5Ld7szi9bcBfOoFv/3dc6xSMkL2PC0=
-github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d h1:5PJl274Y63IEHC+7izoQE9x6ikvDFZS2mDVS3drnohI=
-github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d/go.mod h1:01TrycV0kFyexm33Z7vhZRXopbI8J3TDReVlkTgMUxE=
 github.com/miekg/dns v1.1.26/go.mod h1:bPDLeHnStXmXAq1m/Ch/hvfNHr14JKNPMBo3VZKjuso=
 github.com/miekg/dns v1.1.41 h1:WMszZWJG0XmzbK9FEmzH2TVcqYzFesusSIB41b8KHxY=
 github.com/miekg/dns v1.1.41/go.mod h1:p6aan82bvRIyn+zDIv9xYNUpwa73JcSh9BKwknJysuI=
@@ -217,14 +222,8 @@ github.com/modern-go/reflect2 v1.0.2/go.mod h1:yWuevngMOJpCy52FWWMvUC8ws7m/LJsjY
 github.com/montanaflynn/stats v0.7.1 h1:etflOAAHORrCC44V+aR6Ftzort912ZU+YLiSTuV8eaE=
 github.com/montanaflynn/stats v0.7.1/go.mod h1:etXPPgVO6n31NxCd9KQUMvCM+ve0ruNzt6R8Bnaayow=
 github.com/mwitkow/go-conntrack v0.0.0-20161129095857-cc309e4a2223/go.mod h1:qRWi+5nqEBWmkhHvq77mSJWrCKwh8bxhgT7d/eI7P4U=
-github.com/nxadm/tail v1.4.11 h1:8feyoE3OzPrcshW5/MJ4sGESc5cqmGkGCWlco4l0bqY=
-github.com/nxadm/tail v1.4.11/go.mod h1:OTaG3NK980DZzxbRq6lEuzgU+mug70nY11sMd4JXXHc=
 github.com/oapi-codegen/runtime v1.1.1 h1:EXLHh0DXIJnWhdRPN2w4MXAzFyE4CskzhNLUmtpMYro=
 github.com/oapi-codegen/runtime v1.1.1/go.mod h1:SK9X900oXmPWilYR5/WKPzt3Kqxn/uS/+lbpREv+eCg=
-github.com/onsi/ginkgo v1.16.5 h1:8xi0RTUf59SOSfEtZMvwTvXYMzG4gV23XVHOZiXNtnE=
-github.com/onsi/ginkgo v1.16.5/go.mod h1:+E8gABHa3K6zRBolWtd+ROzc/U5bkGt0FwiG042wbpU=
-github.com/onsi/gomega v1.36.2 h1:koNYke6TVk6ZmnyHrCXba/T/MoLBXFjeC1PtvYgw0A8=
-github.com/onsi/gomega v1.36.2/go.mod h1:DdwyADRjrc825LhMEkD76cHR5+pUnjhUN8GlHlRPHzY=
 github.com/opentracing/opentracing-go v1.1.0/go.mod h1:UkNAQd3GIcIGf0SeVgPpRdFStlNbqXla1AfSYxPUl2o=
 github.com/pascaldekloe/goe v0.0.0-20180627143212-57f6aae5913c/go.mod h1:lzWF7FIEvWOWxwDKqyGYQf6ZUaNfKdP144TG7ZOy1lc=
 github.com/pascaldekloe/goe v0.1.0 h1:cBOtyMzM9HTpWjXfbbunk26uA6nG3a8n06Wieeh0MwY=
@@ -254,6 +253,8 @@ github.com/prometheus/procfs v0.0.2/go.mod h1:TjEm7ze935MbeOT/UhFTIMYKhuLP4wbCsT
 github.com/prometheus/procfs v0.0.8/go.mod h1:7Qr8sr6344vo1JqZ6HhLceV9o3AJ1Ff+GxbHq6oeK9A=
 github.com/rabbitmq/amqp091-go v1.10.0 h1:STpn5XsHlHGcecLmMFCtg7mqq0RnD+zFr4uzukfVhBw=
 github.com/rabbitmq/amqp091-go v1.10.0/go.mod h1:Hy4jKW5kQART1u+JkDTF9YYOQUHXqMuhrgxOEeS7G4o=
+github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5 h1:mZHayPoR0lNmnHyvtYjDeq0zlVHn9K/ZXoy17ylucdo=
+github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5/go.mod h1:GEXHk5HgEKCvEIIrSpFI3ozzG5xOKA2DVlEX/gGnewM=
 github.com/rogpeppe/go-internal v1.9.0 h1:73kH8U+JUqXU8lRuOHeVHaa/SZPifC7BkcraZVejAe8=
 github.com/rogpeppe/go-internal v1.9.0/go.mod h1:WtVeX8xhTBvf0smdhujwtBcq4Qrzq/fJaraNFVN+nFs=
 github.com/ryanuber/columnize v0.0.0-20160712163229-9b3edd62028f/go.mod h1:sm1tb6uqfes/u+d4ooFouqFdy9/2g9QGwK3SQygK0Ts=
@@ -303,8 +304,6 @@ github.com/twitchyliquid64/golang-asm v0.15.1 h1:SU5vSMR7hnwNxj24w34ZyCi/FmDZTkS
 github.com/twitchyliquid64/golang-asm v0.15.1/go.mod h1:a1lVb/DtPvCB8fslRZhAngC2+aY1QWCk3Cedj/Gdt08=
 github.com/ugorji/go/codec v1.2.12 h1:9LC83zGrHhuUA9l16C9AHXAqEV/2wBQ4nkvumAE65EE=
 github.com/ugorji/go/codec v1.2.12/go.mod h1:UNopzCgEMSXjBc6AOMqYvWC1ktqTAfzJZUZgYf6w6lg=
-github.com/x-cray/logrus-prefixed-formatter v0.5.2 h1:00txxvfBM9muc0jiLIEAkAcIMJzfthRT6usrui8uGmg=
-github.com/x-cray/logrus-prefixed-formatter v0.5.2/go.mod h1:2duySbKsL6M18s5GU7VPsoEPHyzalCE06qoARUCeBBE=
 github.com/xdg-go/pbkdf2 v1.0.0 h1:Su7DPu48wXMwC3bs7MCNG+z4FhcyEuz5dlvchbq0B0c=
 github.com/xdg-go/pbkdf2 v1.0.0/go.mod h1:jrpuAogTd400dnrH08LKmI/xc1MbPOebTwRqcT5RDeI=
 github.com/xdg-go/scram v1.1.2 h1:FHX5I5B4i4hKRVRBCFRxq1iQRej7WO3hhBuJf+UUySY=
@@ -419,15 +418,12 @@ golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8/go.mod h1:oPkhp1MJrh7nUepCBc
 golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20220728004956-3c1f35247d10/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
-golang.org/x/sys v0.0.0-20220908164124-27713097b956/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.5.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.6.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.28.0 h1:Fksou7UEQUWlKvIdsqzJmUmCX3cZuD2+P3XyyzwMhlA=
 golang.org/x/sys v0.28.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
 golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1/go.mod h1:bj7SfCRtBDWHUb9snDiAeCFNEtKQo2Wmx5Cou7ajbmo=
 golang.org/x/term v0.0.0-20210927222741-03fcf44c2211/go.mod h1:jbD1KX2456YbFQfuXm/mYQcufACuNUgVhRMnK/tPxf8=
-golang.org/x/term v0.27.0 h1:WP60Sv1nlK1T6SupCHbXzSaN0b9wUmsPoRS9b61A23Q=
-golang.org/x/term v0.27.0/go.mod h1:iMsnZpn0cago0GOrHO2+Y7u7JPn5AylBrcoWkElMTSM=
 golang.org/x/text v0.3.0/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
 golang.org/x/text v0.3.2/go.mod h1:bEr9sfX3Q8Zfm5fL9x+3itogRgK3+ptLWKqgva+5dAk=
 golang.org/x/text v0.3.3/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
@@ -474,8 +470,6 @@ gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 h1:YR8cESwS4TdDjEe65xsg0ogR
 gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
 gopkg.in/ini.v1 v1.67.0 h1:Dgnx+6+nfE+IfzjUEISNeydPJh9AXNNsWbGP9KzCsOA=
 gopkg.in/ini.v1 v1.67.0/go.mod h1:pNLf8WUiyNEtQjuu5G5vTm06TEv9tsIgeAvK8hOrP4k=
-gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 h1:uRGJdciOHaEIrze2W8Q3AKkepLTh2hOroT7a+7czfdQ=
-gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7/go.mod h1:dt/ZhP58zS4L8KSrWDmTeBkI65Dw0HsyUHuEVlX15mw=
 gopkg.in/yaml.v2 v2.2.1/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.2/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.4/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
diff --git a/internal/order/service/application.go b/internal/order/service/application.go
index 59bbe5a..2879b67 100644
--- a/internal/order/service/application.go
+++ b/internal/order/service/application.go
@@ -44,15 +44,14 @@ func newApplication(_ context.Context, stockGRPC query.StockService, ch *amqp.Ch
 	//orderRepo := adapters.NewMemoryOrderRepository()
 	mongoClient := newMongoClient()
 	orderRepo := adapters.NewOrderRepositoryMongo(mongoClient)
-	logger := logrus.NewEntry(logrus.StandardLogger())
 	metricClient := metrics.TodoMetrics{}
 	return app.Application{
 		Commands: app.Commands{
-			CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGRPC, ch, logger, metricClient),
-			UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger, metricClient),
+			CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGRPC, ch, logrus.StandardLogger(), metricClient),
+			UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logrus.StandardLogger(), metricClient),
 		},
 		Queries: app.Queries{
-			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricClient),
+			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logrus.StandardLogger(), metricClient),
 		},
 	}
 }
diff --git a/internal/payment/app/command/create_payment.go b/internal/payment/app/command/create_payment.go
index 84f629f..3a51acb 100644
--- a/internal/payment/app/command/create_payment.go
+++ b/internal/payment/app/command/create_payment.go
@@ -47,7 +47,7 @@ func (c createPaymentHandler) Handle(ctx context.Context, cmd CreatePayment) (st
 func NewCreatePaymentHandler(
 	processor domain.Processor,
 	orderGRPC OrderService,
-	logger *logrus.Entry,
+	logger *logrus.Logger,
 	metricClient decorator.MetricsClient,
 ) CreatePaymentHandler {
 	return decorator.ApplyCommandDecorators[CreatePayment, string](
diff --git a/internal/payment/go.mod b/internal/payment/go.mod
index 452ceaf..1bbedc8 100644
--- a/internal/payment/go.mod
+++ b/internal/payment/go.mod
@@ -17,6 +17,7 @@ require (
 )
 
 require (
+	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
 	github.com/armon/go-metrics v0.4.1 // indirect
 	github.com/bytedance/sonic v1.12.3 // indirect
 	github.com/bytedance/sonic/loader v0.2.0 // indirect
@@ -48,16 +49,18 @@ require (
 	github.com/json-iterator/go v1.1.12 // indirect
 	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
 	github.com/leodido/go-urn v1.4.0 // indirect
+	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible // indirect
+	github.com/lestrrat-go/strftime v1.1.0 // indirect
 	github.com/magiconair/properties v1.8.7 // indirect
 	github.com/mattn/go-colorable v0.1.13 // indirect
 	github.com/mattn/go-isatty v0.0.20 // indirect
-	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
 	github.com/mitchellh/go-homedir v1.1.0 // indirect
 	github.com/mitchellh/mapstructure v1.5.0 // indirect
 	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
 	github.com/modern-go/reflect2 v1.0.2 // indirect
-	github.com/nxadm/tail v1.4.11 // indirect
+	github.com/oapi-codegen/runtime v1.1.1 // indirect
 	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
+	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5 // indirect
 	github.com/sagikazarmark/locafero v0.4.0 // indirect
 	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
 	github.com/sourcegraph/conc v0.3.0 // indirect
@@ -67,7 +70,6 @@ require (
 	github.com/subosito/gotenv v1.6.0 // indirect
 	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
 	github.com/ugorji/go/codec v1.2.12 // indirect
-	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
 	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0 // indirect
 	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0 // indirect
 	go.opentelemetry.io/contrib/propagators/b3 v1.31.0 // indirect
@@ -82,7 +84,6 @@ require (
 	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
 	golang.org/x/net v0.33.0 // indirect
 	golang.org/x/sys v0.28.0 // indirect
-	golang.org/x/term v0.27.0 // indirect
 	golang.org/x/text v0.21.0 // indirect
 	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
 	google.golang.org/protobuf v1.36.1 // indirect
diff --git a/internal/payment/go.sum b/internal/payment/go.sum
index d1f9f0d..038e476 100644
--- a/internal/payment/go.sum
+++ b/internal/payment/go.sum
@@ -1,10 +1,13 @@
 cloud.google.com/go v0.26.0/go.mod h1:aQUYkXzVsufM+DwF1aE+0xfcU+56JwCaLick0ClmMTw=
 github.com/BurntSushi/toml v0.3.1/go.mod h1:xHWCNGjB5oqiDr8zfno3MHue2Ht5sIBksp03qcyfWMU=
 github.com/DataDog/datadog-go v3.2.0+incompatible/go.mod h1:LButxg5PwREeZtORoXG3tL4fMGNddJ+vMq1mwgfaqoQ=
+github.com/RaveNoX/go-jsoncommentstrip v1.0.0/go.mod h1:78ihd09MekBnJnxpICcwzCMzGrKSKYe4AqU6PDYYpjk=
 github.com/alecthomas/template v0.0.0-20160405071501-a0175ee3bccc/go.mod h1:LOuyumcjzFXgccqObfd/Ljyb9UuFJ6TxHnclSeseNhc=
 github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751/go.mod h1:LOuyumcjzFXgccqObfd/Ljyb9UuFJ6TxHnclSeseNhc=
 github.com/alecthomas/units v0.0.0-20151022065526-2efee857e7cf/go.mod h1:ybxpYRFXyAe+OPACYpWeL0wqObRcbAqCMya13uyzqw0=
 github.com/alecthomas/units v0.0.0-20190717042225-c3de453c63f4/go.mod h1:ybxpYRFXyAe+OPACYpWeL0wqObRcbAqCMya13uyzqw0=
+github.com/apapsch/go-jsonmerge/v2 v2.0.0 h1:axGnT1gRIfimI7gJifB699GoE/oq+F2MU7Dml6nw9rQ=
+github.com/apapsch/go-jsonmerge/v2 v2.0.0/go.mod h1:lvDnEdqiQrp0O42VQGgmlKpxL1AP2+08jFMw88y4klk=
 github.com/armon/circbuf v0.0.0-20150827004946-bbbad097214e/go.mod h1:3U/XgcO3hCbHZ8TKRvWD2dDTCfh9M9ya+I9JpbB7O8o=
 github.com/armon/go-metrics v0.0.0-20180917152333-f0300d1749da/go.mod h1:Q73ZrmVTwzkszR9V5SSuryQ31EELlFMUz1kKyl939pY=
 github.com/armon/go-metrics v0.4.1 h1:hR91U9KYmb6bLBYLQjyM+3j+rcd/UhE+G78SFnF8gJA=
@@ -16,6 +19,7 @@ github.com/beorn7/perks v0.0.0-20180321164747-3a771d992973/go.mod h1:Dwedo/Wpr24
 github.com/beorn7/perks v1.0.0/go.mod h1:KWe93zE9D1o94FZ5RNwFwVgaQK1VOXiVxmqh+CedLV8=
 github.com/beorn7/perks v1.0.1/go.mod h1:G2ZrVWU2WbWT9wwq4/hrbKbnv/1ERSJQ0ibhJ6rlkpw=
 github.com/bgentry/speakeasy v0.1.0/go.mod h1:+zsyZBPWlz7T6j88CTgSN5bM796AkVf0kBD4zp0CCIs=
+github.com/bmatcuk/doublestar v1.1.1/go.mod h1:UD6OnuiIn0yFxxA2le/rnRU1G4RaI4UvFv1sNto9p6w=
 github.com/bytedance/sonic v1.12.3 h1:W2MGa7RCU1QTeYRTPE3+88mVC0yXmsRQRChiyVocVjU=
 github.com/bytedance/sonic v1.12.3/go.mod h1:B8Gt/XvtZ3Fqj+iSKMypzymZxw/FVwgIGKzMzT9r/rk=
 github.com/bytedance/sonic/loader v0.1.1/go.mod h1:ncP89zfokxS5LZrJxl5z0UJcsk4M4yY2JpfqGeCtNLU=
@@ -46,7 +50,6 @@ github.com/fatih/color v1.14.1 h1:qfhVLaG5s+nCROl1zJsZRxFeYrHLqWroPOQ8BWiNb4w=
 github.com/fatih/color v1.14.1/go.mod h1:2oHN61fhTpgcxD3TSWCgKDiH1+x4OiDVVGH8WlgGZGg=
 github.com/frankban/quicktest v1.14.6 h1:7Xjx+VpznH+oBnejlPUj8oUpdxnVs4f8XU8WnHkI4W8=
 github.com/frankban/quicktest v1.14.6/go.mod h1:4ptaffx2x8+WTWXmUCuVU6aPUX1/Mz7zb5vbUoiM6w0=
-github.com/fsnotify/fsnotify v1.6.0/go.mod h1:sl3t1tCWJFWoRz9R8WJCbQihKKwmorjAbSClcnxKAGw=
 github.com/fsnotify/fsnotify v1.7.0 h1:8JEhPFa5W2WU7YfeZzPNqzMP6Lwt7L2715Ggo0nosvA=
 github.com/fsnotify/fsnotify v1.7.0/go.mod h1:40Bi/Hjc2AVfZrqy+aj+yEI+/bRxZnMJyTJwOpGvigM=
 github.com/gabriel-vasile/mimetype v1.4.5 h1:J7wGKdGu33ocBOhGy0z653k/lFKLFDPJMG8Gql0kxn4=
@@ -149,10 +152,13 @@ github.com/hashicorp/memberlist v0.5.0 h1:EtYPN8DpAURiapus508I4n9CzHs2W+8NZGbmmR
 github.com/hashicorp/memberlist v0.5.0/go.mod h1:yvyXLpo0QaGE59Y7hDTsTzDD25JYBZ4mHgHUZ8lrOI0=
 github.com/hashicorp/serf v0.10.1 h1:Z1H2J60yRKvfDYAOZLd2MU0ND4AH/WDz7xYHDWQsIPY=
 github.com/hashicorp/serf v0.10.1/go.mod h1:yL2t6BqATOLGc5HF7qbFkTfXoPIY0WZdWHfEvMqbG+4=
+github.com/jonboulle/clockwork v0.5.0 h1:Hyh9A8u51kptdkR+cqRpT1EebBwTn1oK9YfGYbdFz6I=
+github.com/jonboulle/clockwork v0.5.0/go.mod h1:3mZlmanh0g2NDKO5TWZVJAfofYk64M7XN3SzBPjZF60=
 github.com/json-iterator/go v1.1.6/go.mod h1:+SdeFBvtyEkXs7REEP0seUULqWtbJapLOCVDaaPEHmU=
 github.com/json-iterator/go v1.1.9/go.mod h1:KdQUCv79m/52Kvf8AW2vK1V8akMuk1QjK/uOdHXbAo4=
 github.com/json-iterator/go v1.1.12 h1:PV8peI4a0ysnczrg+LtxykD8LfKY9ML6u2jnxaEnrnM=
 github.com/json-iterator/go v1.1.12/go.mod h1:e30LSqwooZae/UwlEbR2852Gd8hjQvJoHmT4TnhNGBo=
+github.com/juju/gnuflag v0.0.0-20171113085948-2ce1bb71843d/go.mod h1:2PavIy+JPciBPrBUjwbNvtwB6RQlve+hkpll6QSNmOE=
 github.com/julienschmidt/httprouter v1.2.0/go.mod h1:SYymIcj16QtmaHHD7aYtjjsJG7VTCxuUUipMqKk8s4w=
 github.com/kisielk/errcheck v1.5.0/go.mod h1:pFxgyoBC7bSaBwPgfKdkLd5X25qrDl4LWUI2bnpBCr8=
 github.com/kisielk/gotool v1.0.0/go.mod h1:XhKaO+MFFWcvkIS/tQcRk01m1F5IRFswLeQ+oQHNcck=
@@ -171,6 +177,12 @@ github.com/kr/text v0.2.0 h1:5Nx0Ya0ZqY2ygV366QzturHI13Jq95ApcVaJBhpS+AY=
 github.com/kr/text v0.2.0/go.mod h1:eLer722TekiGuMkidMxC/pM04lWEeraHUUmBw8l2grE=
 github.com/leodido/go-urn v1.4.0 h1:WT9HwE9SGECu3lg4d/dIA+jxlljEa1/ffXKmRjqdmIQ=
 github.com/leodido/go-urn v1.4.0/go.mod h1:bvxc+MVxLKB4z00jd1z+Dvzr47oO32F/QSNjSBOlFxI=
+github.com/lestrrat-go/envload v0.0.0-20180220234015-a3eb8ddeffcc h1:RKf14vYWi2ttpEmkA4aQ3j4u9dStX2t4M8UM6qqNsG8=
+github.com/lestrrat-go/envload v0.0.0-20180220234015-a3eb8ddeffcc/go.mod h1:kopuH9ugFRkIXf3YoqHKyrJ9YfUFsckUU9S7B+XP+is=
+github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible h1:Y6sqxHMyB1D2YSzWkLibYKgg+SwmyFU9dF2hn6MdTj4=
+github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible/go.mod h1:ZQnN8lSECaebrkQytbHj4xNgtg8CR7RYXnPok8e0EHA=
+github.com/lestrrat-go/strftime v1.1.0 h1:gMESpZy44/4pXLO/m+sL0yBd1W6LjgjrrD4a68Gapyg=
+github.com/lestrrat-go/strftime v1.1.0/go.mod h1:uzeIB52CeUJenCo1syghlugshMysrqUT51HlxphXVeI=
 github.com/magiconair/properties v1.8.7 h1:IeQXZAiQcpL9mgcAe1Nu6cX9LLw6ExEHKjN0VQdvPDY=
 github.com/magiconair/properties v1.8.7/go.mod h1:Dhd985XPs7jluiymwWYZ0G4Z61jb3vdS329zhj2hYo0=
 github.com/mattn/go-colorable v0.0.9/go.mod h1:9vuHe8Xs5qXnSaW/c/ABM9alt+Vo+STaOChaDxuIBZU=
@@ -189,8 +201,6 @@ github.com/mattn/go-isatty v0.0.16/go.mod h1:kYGgaQfpe5nmfYZH+SKPsOc2e4SrIfOl2e/
 github.com/mattn/go-isatty v0.0.20 h1:xfD0iDuEKnDkl03q4limB+vH+GxLEtL/jb4xVJSWWEY=
 github.com/mattn/go-isatty v0.0.20/go.mod h1:W+V8PltTTMOvKvAeJH7IuucS94S2C6jfK/D7dTCTo3Y=
 github.com/matttproud/golang_protobuf_extensions v1.0.1/go.mod h1:D8He9yQNgCq6Z5Ld7szi9bcBfOoFv/3dc6xSMkL2PC0=
-github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d h1:5PJl274Y63IEHC+7izoQE9x6ikvDFZS2mDVS3drnohI=
-github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d/go.mod h1:01TrycV0kFyexm33Z7vhZRXopbI8J3TDReVlkTgMUxE=
 github.com/miekg/dns v1.1.26/go.mod h1:bPDLeHnStXmXAq1m/Ch/hvfNHr14JKNPMBo3VZKjuso=
 github.com/miekg/dns v1.1.41 h1:WMszZWJG0XmzbK9FEmzH2TVcqYzFesusSIB41b8KHxY=
 github.com/miekg/dns v1.1.41/go.mod h1:p6aan82bvRIyn+zDIv9xYNUpwa73JcSh9BKwknJysuI=
@@ -208,12 +218,8 @@ github.com/modern-go/reflect2 v1.0.1/go.mod h1:bx2lNnkwVCuqBIxFjflWJWanXIb3Rllmb
 github.com/modern-go/reflect2 v1.0.2 h1:xBagoLtFs94CBntxluKeaWgTMpvLxC4ur3nMaC9Gz0M=
 github.com/modern-go/reflect2 v1.0.2/go.mod h1:yWuevngMOJpCy52FWWMvUC8ws7m/LJsjYzDa0/r8luk=
 github.com/mwitkow/go-conntrack v0.0.0-20161129095857-cc309e4a2223/go.mod h1:qRWi+5nqEBWmkhHvq77mSJWrCKwh8bxhgT7d/eI7P4U=
-github.com/nxadm/tail v1.4.11 h1:8feyoE3OzPrcshW5/MJ4sGESc5cqmGkGCWlco4l0bqY=
-github.com/nxadm/tail v1.4.11/go.mod h1:OTaG3NK980DZzxbRq6lEuzgU+mug70nY11sMd4JXXHc=
-github.com/onsi/ginkgo v1.16.5 h1:8xi0RTUf59SOSfEtZMvwTvXYMzG4gV23XVHOZiXNtnE=
-github.com/onsi/ginkgo v1.16.5/go.mod h1:+E8gABHa3K6zRBolWtd+ROzc/U5bkGt0FwiG042wbpU=
-github.com/onsi/gomega v1.36.2 h1:koNYke6TVk6ZmnyHrCXba/T/MoLBXFjeC1PtvYgw0A8=
-github.com/onsi/gomega v1.36.2/go.mod h1:DdwyADRjrc825LhMEkD76cHR5+pUnjhUN8GlHlRPHzY=
+github.com/oapi-codegen/runtime v1.1.1 h1:EXLHh0DXIJnWhdRPN2w4MXAzFyE4CskzhNLUmtpMYro=
+github.com/oapi-codegen/runtime v1.1.1/go.mod h1:SK9X900oXmPWilYR5/WKPzt3Kqxn/uS/+lbpREv+eCg=
 github.com/opentracing/opentracing-go v1.1.0/go.mod h1:UkNAQd3GIcIGf0SeVgPpRdFStlNbqXla1AfSYxPUl2o=
 github.com/pascaldekloe/goe v0.0.0-20180627143212-57f6aae5913c/go.mod h1:lzWF7FIEvWOWxwDKqyGYQf6ZUaNfKdP144TG7ZOy1lc=
 github.com/pascaldekloe/goe v0.1.0 h1:cBOtyMzM9HTpWjXfbbunk26uA6nG3a8n06Wieeh0MwY=
@@ -243,6 +249,8 @@ github.com/prometheus/procfs v0.0.2/go.mod h1:TjEm7ze935MbeOT/UhFTIMYKhuLP4wbCsT
 github.com/prometheus/procfs v0.0.8/go.mod h1:7Qr8sr6344vo1JqZ6HhLceV9o3AJ1Ff+GxbHq6oeK9A=
 github.com/rabbitmq/amqp091-go v1.10.0 h1:STpn5XsHlHGcecLmMFCtg7mqq0RnD+zFr4uzukfVhBw=
 github.com/rabbitmq/amqp091-go v1.10.0/go.mod h1:Hy4jKW5kQART1u+JkDTF9YYOQUHXqMuhrgxOEeS7G4o=
+github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5 h1:mZHayPoR0lNmnHyvtYjDeq0zlVHn9K/ZXoy17ylucdo=
+github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5/go.mod h1:GEXHk5HgEKCvEIIrSpFI3ozzG5xOKA2DVlEX/gGnewM=
 github.com/rogpeppe/go-internal v1.9.0 h1:73kH8U+JUqXU8lRuOHeVHaa/SZPifC7BkcraZVejAe8=
 github.com/rogpeppe/go-internal v1.9.0/go.mod h1:WtVeX8xhTBvf0smdhujwtBcq4Qrzq/fJaraNFVN+nFs=
 github.com/ryanuber/columnize v0.0.0-20160712163229-9b3edd62028f/go.mod h1:sm1tb6uqfes/u+d4ooFouqFdy9/2g9QGwK3SQygK0Ts=
@@ -266,6 +274,7 @@ github.com/spf13/pflag v1.0.5 h1:iy+VFUOCP1a+8yFto/drg2CJ5u0yRoB7fZw3DKv/JXA=
 github.com/spf13/pflag v1.0.5/go.mod h1:McXfInJRrz4CZXVZOBLb0bTZqETkiAhM9Iw0y3An2Bg=
 github.com/spf13/viper v1.19.0 h1:RWq5SEjt8o25SROyN3z2OrDB9l7RPd3lwTWU8EcEdcI=
 github.com/spf13/viper v1.19.0/go.mod h1:GQUN9bilAbhU/jgc1bKs99f/suXKeUMct8Adx5+Ntkg=
+github.com/spkg/bom v0.0.0-20160624110644-59b7046e48ad/go.mod h1:qLr4V1qq6nMqFKkMo8ZTx3f+BZEkzsRUY10Xsm2mwU0=
 github.com/stretchr/objx v0.1.0/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
 github.com/stretchr/objx v0.1.1/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
 github.com/stretchr/objx v0.4.0/go.mod h1:YvHI0jy2hoMjB+UWwv71VJQ9isScKT/TqJzVSSt89Yw=
@@ -290,8 +299,6 @@ github.com/twitchyliquid64/golang-asm v0.15.1 h1:SU5vSMR7hnwNxj24w34ZyCi/FmDZTkS
 github.com/twitchyliquid64/golang-asm v0.15.1/go.mod h1:a1lVb/DtPvCB8fslRZhAngC2+aY1QWCk3Cedj/Gdt08=
 github.com/ugorji/go/codec v1.2.12 h1:9LC83zGrHhuUA9l16C9AHXAqEV/2wBQ4nkvumAE65EE=
 github.com/ugorji/go/codec v1.2.12/go.mod h1:UNopzCgEMSXjBc6AOMqYvWC1ktqTAfzJZUZgYf6w6lg=
-github.com/x-cray/logrus-prefixed-formatter v0.5.2 h1:00txxvfBM9muc0jiLIEAkAcIMJzfthRT6usrui8uGmg=
-github.com/x-cray/logrus-prefixed-formatter v0.5.2/go.mod h1:2duySbKsL6M18s5GU7VPsoEPHyzalCE06qoARUCeBBE=
 github.com/yuin/goldmark v1.1.27/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
 github.com/yuin/goldmark v1.2.1/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
 go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0 h1:0nTRpaCaILLdooXAQnfktlL6Zw1ECKEW9DZGH2byi2c=
@@ -387,14 +394,11 @@ golang.org/x/sys v0.0.0-20211025201205-69cdffdb9359/go.mod h1:oPkhp1MJrh7nUepCBc
 golang.org/x/sys v0.0.0-20220503163025-988cb79eb6c6/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20220728004956-3c1f35247d10/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
-golang.org/x/sys v0.0.0-20220908164124-27713097b956/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.5.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.6.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.28.0 h1:Fksou7UEQUWlKvIdsqzJmUmCX3cZuD2+P3XyyzwMhlA=
 golang.org/x/sys v0.28.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
 golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1/go.mod h1:bj7SfCRtBDWHUb9snDiAeCFNEtKQo2Wmx5Cou7ajbmo=
-golang.org/x/term v0.27.0 h1:WP60Sv1nlK1T6SupCHbXzSaN0b9wUmsPoRS9b61A23Q=
-golang.org/x/term v0.27.0/go.mod h1:iMsnZpn0cago0GOrHO2+Y7u7JPn5AylBrcoWkElMTSM=
 golang.org/x/text v0.3.0/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
 golang.org/x/text v0.3.2/go.mod h1:bEr9sfX3Q8Zfm5fL9x+3itogRgK3+ptLWKqgva+5dAk=
 golang.org/x/text v0.3.3/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
@@ -440,8 +444,6 @@ gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 h1:YR8cESwS4TdDjEe65xsg0ogR
 gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
 gopkg.in/ini.v1 v1.67.0 h1:Dgnx+6+nfE+IfzjUEISNeydPJh9AXNNsWbGP9KzCsOA=
 gopkg.in/ini.v1 v1.67.0/go.mod h1:pNLf8WUiyNEtQjuu5G5vTm06TEv9tsIgeAvK8hOrP4k=
-gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 h1:uRGJdciOHaEIrze2W8Q3AKkepLTh2hOroT7a+7czfdQ=
-gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7/go.mod h1:dt/ZhP58zS4L8KSrWDmTeBkI65Dw0HsyUHuEVlX15mw=
 gopkg.in/yaml.v2 v2.2.1/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.2/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.4/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
diff --git a/internal/payment/service/application.go b/internal/payment/service/application.go
index 0199c72..d944ffc 100644
--- a/internal/payment/service/application.go
+++ b/internal/payment/service/application.go
@@ -28,11 +28,10 @@ func NewApplication(ctx context.Context) (app.Application, func()) {
 }
 
 func newApplication(_ context.Context, orderGRPC command.OrderService, processor domain.Processor) app.Application {
-	logger := logrus.NewEntry(logrus.StandardLogger())
 	metricClient := metrics.TodoMetrics{}
 	return app.Application{
 		Commands: app.Commands{
-			CreatePayment: command.NewCreatePaymentHandler(processor, orderGRPC, logger, metricClient),
+			CreatePayment: command.NewCreatePaymentHandler(processor, orderGRPC, logrus.StandardLogger(), metricClient),
 		},
 	}
 }
diff --git a/internal/stock/adapters/stock_mysql_repository.go b/internal/stock/adapters/stock_mysql_repository.go
index bf5f54d..a9abaa8 100644
--- a/internal/stock/adapters/stock_mysql_repository.go
+++ b/internal/stock/adapters/stock_mysql_repository.go
@@ -4,10 +4,10 @@ import (
 	"context"
 
 	"github.com/ghost-yu/go_shop_second/common/entity"
+	"github.com/ghost-yu/go_shop_second/common/logging"
 	"github.com/ghost-yu/go_shop_second/stock/infrastructure/persistent"
 	"github.com/ghost-yu/go_shop_second/stock/infrastructure/persistent/builder"
 	"github.com/pkg/errors"
-	"github.com/sirupsen/logrus"
 	"gorm.io/gorm"
 )
 
@@ -51,7 +51,7 @@ func (m MySQLStockRepository) UpdateStock(
 	return m.db.StartTransaction(func(tx *gorm.DB) (err error) {
 		defer func() {
 			if err != nil {
-				logrus.Warnf("update stock transaction err=%v", err)
+				logging.Warnf(ctx, nil, "update stock transaction err=%v", err)
 			}
 		}()
 		err = m.updatePessimistic(ctx, tx, data, updateFn)
diff --git a/internal/stock/app/query/check_if_items_in_stock.go b/internal/stock/app/query/check_if_items_in_stock.go
index 4a5df16..6cec6b4 100644
--- a/internal/stock/app/query/check_if_items_in_stock.go
+++ b/internal/stock/app/query/check_if_items_in_stock.go
@@ -33,7 +33,7 @@ type checkIfItemsInStockHandler struct {
 func NewCheckIfItemsInStockHandler(
 	stockRepo domain.Repository,
 	stripeAPI *integration.StripeAPI,
-	logger *logrus.Entry,
+	logger *logrus.Logger,
 	metricClient decorator.MetricsClient,
 ) CheckIfItemsInStockHandler {
 	if stockRepo == nil {
@@ -67,8 +67,20 @@ func (h checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfIte
 			logging.Warnf(ctx, nil, "redis unlock fail, err=%v", err)
 		}
 	}()
-
+	var err error
 	var res []*entity.Item
+	defer func() {
+		f := logrus.Fields{
+			"query": query,
+			"res":   res,
+		}
+		if err != nil {
+			logging.Errorf(ctx, f, "checkIfItemsInStock err=%v", err)
+		} else {
+			logging.Infof(ctx, f, "%s", "checkIfItemsInStock success")
+		}
+	}()
+
 	for _, i := range query.Items {
 		priceID, err := h.stripeAPI.GetPriceByProductID(ctx, i.ID)
 		if err != nil || priceID == "" {
diff --git a/internal/stock/app/query/get_items.go b/internal/stock/app/query/get_items.go
index 67f9ae2..1173918 100644
--- a/internal/stock/app/query/get_items.go
+++ b/internal/stock/app/query/get_items.go
@@ -21,7 +21,7 @@ type getItemsHandler struct {
 
 func NewGetItemsHandler(
 	stockRepo domain.Repository,
-	logger *logrus.Entry,
+	logger *logrus.Logger,
 	metricClient decorator.MetricsClient,
 ) GetItemsHandler {
 	if stockRepo == nil {
diff --git a/internal/stock/go.mod b/internal/stock/go.mod
index 6dc6ca8..03273ac 100644
--- a/internal/stock/go.mod
+++ b/internal/stock/go.mod
@@ -17,6 +17,7 @@ require (
 )
 
 require (
+	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
 	github.com/armon/go-metrics v0.4.1 // indirect
 	github.com/bytedance/sonic v1.12.3 // indirect
 	github.com/bytedance/sonic/loader v0.2.0 // indirect
@@ -55,18 +56,20 @@ require (
 	github.com/json-iterator/go v1.1.12 // indirect
 	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
 	github.com/leodido/go-urn v1.4.0 // indirect
+	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible // indirect
+	github.com/lestrrat-go/strftime v1.1.0 // indirect
 	github.com/magiconair/properties v1.8.7 // indirect
 	github.com/mattn/go-colorable v0.1.13 // indirect
 	github.com/mattn/go-isatty v0.0.20 // indirect
-	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
 	github.com/mitchellh/go-homedir v1.1.0 // indirect
 	github.com/mitchellh/mapstructure v1.5.0 // indirect
 	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
 	github.com/modern-go/reflect2 v1.0.2 // indirect
-	github.com/nxadm/tail v1.4.11 // indirect
+	github.com/oapi-codegen/runtime v1.1.1 // indirect
 	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
 	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
 	github.com/redis/go-redis/v9 v9.7.0 // indirect
+	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5 // indirect
 	github.com/sagikazarmark/locafero v0.4.0 // indirect
 	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
 	github.com/sourcegraph/conc v0.3.0 // indirect
@@ -76,7 +79,6 @@ require (
 	github.com/subosito/gotenv v1.6.0 // indirect
 	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
 	github.com/ugorji/go/codec v1.2.12 // indirect
-	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
 	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0 // indirect
 	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0 // indirect
 	go.opentelemetry.io/contrib/propagators/b3 v1.31.0 // indirect
@@ -92,7 +94,6 @@ require (
 	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
 	golang.org/x/net v0.33.0 // indirect
 	golang.org/x/sys v0.28.0 // indirect
-	golang.org/x/term v0.27.0 // indirect
 	golang.org/x/text v0.21.0 // indirect
 	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
 	google.golang.org/protobuf v1.36.1 // indirect
diff --git a/internal/stock/go.sum b/internal/stock/go.sum
index 6e6ddad..3a14d50 100644
--- a/internal/stock/go.sum
+++ b/internal/stock/go.sum
@@ -1,10 +1,13 @@
 cloud.google.com/go v0.26.0/go.mod h1:aQUYkXzVsufM+DwF1aE+0xfcU+56JwCaLick0ClmMTw=
 github.com/BurntSushi/toml v0.3.1/go.mod h1:xHWCNGjB5oqiDr8zfno3MHue2Ht5sIBksp03qcyfWMU=
 github.com/DataDog/datadog-go v3.2.0+incompatible/go.mod h1:LButxg5PwREeZtORoXG3tL4fMGNddJ+vMq1mwgfaqoQ=
+github.com/RaveNoX/go-jsoncommentstrip v1.0.0/go.mod h1:78ihd09MekBnJnxpICcwzCMzGrKSKYe4AqU6PDYYpjk=
 github.com/alecthomas/template v0.0.0-20160405071501-a0175ee3bccc/go.mod h1:LOuyumcjzFXgccqObfd/Ljyb9UuFJ6TxHnclSeseNhc=
 github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751/go.mod h1:LOuyumcjzFXgccqObfd/Ljyb9UuFJ6TxHnclSeseNhc=
 github.com/alecthomas/units v0.0.0-20151022065526-2efee857e7cf/go.mod h1:ybxpYRFXyAe+OPACYpWeL0wqObRcbAqCMya13uyzqw0=
 github.com/alecthomas/units v0.0.0-20190717042225-c3de453c63f4/go.mod h1:ybxpYRFXyAe+OPACYpWeL0wqObRcbAqCMya13uyzqw0=
+github.com/apapsch/go-jsonmerge/v2 v2.0.0 h1:axGnT1gRIfimI7gJifB699GoE/oq+F2MU7Dml6nw9rQ=
+github.com/apapsch/go-jsonmerge/v2 v2.0.0/go.mod h1:lvDnEdqiQrp0O42VQGgmlKpxL1AP2+08jFMw88y4klk=
 github.com/armon/circbuf v0.0.0-20150827004946-bbbad097214e/go.mod h1:3U/XgcO3hCbHZ8TKRvWD2dDTCfh9M9ya+I9JpbB7O8o=
 github.com/armon/go-metrics v0.0.0-20180917152333-f0300d1749da/go.mod h1:Q73ZrmVTwzkszR9V5SSuryQ31EELlFMUz1kKyl939pY=
 github.com/armon/go-metrics v0.4.1 h1:hR91U9KYmb6bLBYLQjyM+3j+rcd/UhE+G78SFnF8gJA=
@@ -16,6 +19,7 @@ github.com/beorn7/perks v0.0.0-20180321164747-3a771d992973/go.mod h1:Dwedo/Wpr24
 github.com/beorn7/perks v1.0.0/go.mod h1:KWe93zE9D1o94FZ5RNwFwVgaQK1VOXiVxmqh+CedLV8=
 github.com/beorn7/perks v1.0.1/go.mod h1:G2ZrVWU2WbWT9wwq4/hrbKbnv/1ERSJQ0ibhJ6rlkpw=
 github.com/bgentry/speakeasy v0.1.0/go.mod h1:+zsyZBPWlz7T6j88CTgSN5bM796AkVf0kBD4zp0CCIs=
+github.com/bmatcuk/doublestar v1.1.1/go.mod h1:UD6OnuiIn0yFxxA2le/rnRU1G4RaI4UvFv1sNto9p6w=
 github.com/bsm/ginkgo/v2 v2.12.0 h1:Ny8MWAHyOepLGlLKYmXG4IEkioBysk6GpaRTLC8zwWs=
 github.com/bsm/ginkgo/v2 v2.12.0/go.mod h1:SwYbGRRDovPVboqFv0tPTcG1sN61LM1Z4ARdbAV9g4c=
 github.com/bsm/gomega v1.27.10 h1:yeMWxP2pV2fG3FgAODIY8EiRE3dy0aeFYt4l7wh6yKA=
@@ -54,7 +58,6 @@ github.com/fatih/color v1.14.1 h1:qfhVLaG5s+nCROl1zJsZRxFeYrHLqWroPOQ8BWiNb4w=
 github.com/fatih/color v1.14.1/go.mod h1:2oHN61fhTpgcxD3TSWCgKDiH1+x4OiDVVGH8WlgGZGg=
 github.com/frankban/quicktest v1.14.6 h1:7Xjx+VpznH+oBnejlPUj8oUpdxnVs4f8XU8WnHkI4W8=
 github.com/frankban/quicktest v1.14.6/go.mod h1:4ptaffx2x8+WTWXmUCuVU6aPUX1/Mz7zb5vbUoiM6w0=
-github.com/fsnotify/fsnotify v1.6.0/go.mod h1:sl3t1tCWJFWoRz9R8WJCbQihKKwmorjAbSClcnxKAGw=
 github.com/fsnotify/fsnotify v1.7.0 h1:8JEhPFa5W2WU7YfeZzPNqzMP6Lwt7L2715Ggo0nosvA=
 github.com/fsnotify/fsnotify v1.7.0/go.mod h1:40Bi/Hjc2AVfZrqy+aj+yEI+/bRxZnMJyTJwOpGvigM=
 github.com/gabriel-vasile/mimetype v1.4.5 h1:J7wGKdGu33ocBOhGy0z653k/lFKLFDPJMG8Gql0kxn4=
@@ -163,10 +166,13 @@ github.com/jinzhu/inflection v1.0.0 h1:K317FqzuhWc8YvSVlFMCCUb36O/S9MCKRDI7QkRKD
 github.com/jinzhu/inflection v1.0.0/go.mod h1:h+uFLlag+Qp1Va5pdKtLDYj+kHp5pxUVkryuEj+Srlc=
 github.com/jinzhu/now v1.1.5 h1:/o9tlHleP7gOFmsnYNz3RGnqzefHA47wQpKrrdTIwXQ=
 github.com/jinzhu/now v1.1.5/go.mod h1:d3SSVoowX0Lcu0IBviAWJpolVfI5UJVZZ7cO71lE/z8=
+github.com/jonboulle/clockwork v0.5.0 h1:Hyh9A8u51kptdkR+cqRpT1EebBwTn1oK9YfGYbdFz6I=
+github.com/jonboulle/clockwork v0.5.0/go.mod h1:3mZlmanh0g2NDKO5TWZVJAfofYk64M7XN3SzBPjZF60=
 github.com/json-iterator/go v1.1.6/go.mod h1:+SdeFBvtyEkXs7REEP0seUULqWtbJapLOCVDaaPEHmU=
 github.com/json-iterator/go v1.1.9/go.mod h1:KdQUCv79m/52Kvf8AW2vK1V8akMuk1QjK/uOdHXbAo4=
 github.com/json-iterator/go v1.1.12 h1:PV8peI4a0ysnczrg+LtxykD8LfKY9ML6u2jnxaEnrnM=
 github.com/json-iterator/go v1.1.12/go.mod h1:e30LSqwooZae/UwlEbR2852Gd8hjQvJoHmT4TnhNGBo=
+github.com/juju/gnuflag v0.0.0-20171113085948-2ce1bb71843d/go.mod h1:2PavIy+JPciBPrBUjwbNvtwB6RQlve+hkpll6QSNmOE=
 github.com/julienschmidt/httprouter v1.2.0/go.mod h1:SYymIcj16QtmaHHD7aYtjjsJG7VTCxuUUipMqKk8s4w=
 github.com/kisielk/errcheck v1.5.0/go.mod h1:pFxgyoBC7bSaBwPgfKdkLd5X25qrDl4LWUI2bnpBCr8=
 github.com/kisielk/gotool v1.0.0/go.mod h1:XhKaO+MFFWcvkIS/tQcRk01m1F5IRFswLeQ+oQHNcck=
@@ -185,6 +191,12 @@ github.com/kr/text v0.2.0 h1:5Nx0Ya0ZqY2ygV366QzturHI13Jq95ApcVaJBhpS+AY=
 github.com/kr/text v0.2.0/go.mod h1:eLer722TekiGuMkidMxC/pM04lWEeraHUUmBw8l2grE=
 github.com/leodido/go-urn v1.4.0 h1:WT9HwE9SGECu3lg4d/dIA+jxlljEa1/ffXKmRjqdmIQ=
 github.com/leodido/go-urn v1.4.0/go.mod h1:bvxc+MVxLKB4z00jd1z+Dvzr47oO32F/QSNjSBOlFxI=
+github.com/lestrrat-go/envload v0.0.0-20180220234015-a3eb8ddeffcc h1:RKf14vYWi2ttpEmkA4aQ3j4u9dStX2t4M8UM6qqNsG8=
+github.com/lestrrat-go/envload v0.0.0-20180220234015-a3eb8ddeffcc/go.mod h1:kopuH9ugFRkIXf3YoqHKyrJ9YfUFsckUU9S7B+XP+is=
+github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible h1:Y6sqxHMyB1D2YSzWkLibYKgg+SwmyFU9dF2hn6MdTj4=
+github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible/go.mod h1:ZQnN8lSECaebrkQytbHj4xNgtg8CR7RYXnPok8e0EHA=
+github.com/lestrrat-go/strftime v1.1.0 h1:gMESpZy44/4pXLO/m+sL0yBd1W6LjgjrrD4a68Gapyg=
+github.com/lestrrat-go/strftime v1.1.0/go.mod h1:uzeIB52CeUJenCo1syghlugshMysrqUT51HlxphXVeI=
 github.com/magiconair/properties v1.8.7 h1:IeQXZAiQcpL9mgcAe1Nu6cX9LLw6ExEHKjN0VQdvPDY=
 github.com/magiconair/properties v1.8.7/go.mod h1:Dhd985XPs7jluiymwWYZ0G4Z61jb3vdS329zhj2hYo0=
 github.com/mattn/go-colorable v0.0.9/go.mod h1:9vuHe8Xs5qXnSaW/c/ABM9alt+Vo+STaOChaDxuIBZU=
@@ -203,8 +215,6 @@ github.com/mattn/go-isatty v0.0.16/go.mod h1:kYGgaQfpe5nmfYZH+SKPsOc2e4SrIfOl2e/
 github.com/mattn/go-isatty v0.0.20 h1:xfD0iDuEKnDkl03q4limB+vH+GxLEtL/jb4xVJSWWEY=
 github.com/mattn/go-isatty v0.0.20/go.mod h1:W+V8PltTTMOvKvAeJH7IuucS94S2C6jfK/D7dTCTo3Y=
 github.com/matttproud/golang_protobuf_extensions v1.0.1/go.mod h1:D8He9yQNgCq6Z5Ld7szi9bcBfOoFv/3dc6xSMkL2PC0=
-github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d h1:5PJl274Y63IEHC+7izoQE9x6ikvDFZS2mDVS3drnohI=
-github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d/go.mod h1:01TrycV0kFyexm33Z7vhZRXopbI8J3TDReVlkTgMUxE=
 github.com/miekg/dns v1.1.26/go.mod h1:bPDLeHnStXmXAq1m/Ch/hvfNHr14JKNPMBo3VZKjuso=
 github.com/miekg/dns v1.1.41 h1:WMszZWJG0XmzbK9FEmzH2TVcqYzFesusSIB41b8KHxY=
 github.com/miekg/dns v1.1.41/go.mod h1:p6aan82bvRIyn+zDIv9xYNUpwa73JcSh9BKwknJysuI=
@@ -222,12 +232,8 @@ github.com/modern-go/reflect2 v1.0.1/go.mod h1:bx2lNnkwVCuqBIxFjflWJWanXIb3Rllmb
 github.com/modern-go/reflect2 v1.0.2 h1:xBagoLtFs94CBntxluKeaWgTMpvLxC4ur3nMaC9Gz0M=
 github.com/modern-go/reflect2 v1.0.2/go.mod h1:yWuevngMOJpCy52FWWMvUC8ws7m/LJsjYzDa0/r8luk=
 github.com/mwitkow/go-conntrack v0.0.0-20161129095857-cc309e4a2223/go.mod h1:qRWi+5nqEBWmkhHvq77mSJWrCKwh8bxhgT7d/eI7P4U=
-github.com/nxadm/tail v1.4.11 h1:8feyoE3OzPrcshW5/MJ4sGESc5cqmGkGCWlco4l0bqY=
-github.com/nxadm/tail v1.4.11/go.mod h1:OTaG3NK980DZzxbRq6lEuzgU+mug70nY11sMd4JXXHc=
-github.com/onsi/ginkgo v1.16.5 h1:8xi0RTUf59SOSfEtZMvwTvXYMzG4gV23XVHOZiXNtnE=
-github.com/onsi/ginkgo v1.16.5/go.mod h1:+E8gABHa3K6zRBolWtd+ROzc/U5bkGt0FwiG042wbpU=
-github.com/onsi/gomega v1.36.2 h1:koNYke6TVk6ZmnyHrCXba/T/MoLBXFjeC1PtvYgw0A8=
-github.com/onsi/gomega v1.36.2/go.mod h1:DdwyADRjrc825LhMEkD76cHR5+pUnjhUN8GlHlRPHzY=
+github.com/oapi-codegen/runtime v1.1.1 h1:EXLHh0DXIJnWhdRPN2w4MXAzFyE4CskzhNLUmtpMYro=
+github.com/oapi-codegen/runtime v1.1.1/go.mod h1:SK9X900oXmPWilYR5/WKPzt3Kqxn/uS/+lbpREv+eCg=
 github.com/opentracing/opentracing-go v1.1.0/go.mod h1:UkNAQd3GIcIGf0SeVgPpRdFStlNbqXla1AfSYxPUl2o=
 github.com/pascaldekloe/goe v0.0.0-20180627143212-57f6aae5913c/go.mod h1:lzWF7FIEvWOWxwDKqyGYQf6ZUaNfKdP144TG7ZOy1lc=
 github.com/pascaldekloe/goe v0.1.0 h1:cBOtyMzM9HTpWjXfbbunk26uA6nG3a8n06Wieeh0MwY=
@@ -257,6 +263,8 @@ github.com/prometheus/procfs v0.0.2/go.mod h1:TjEm7ze935MbeOT/UhFTIMYKhuLP4wbCsT
 github.com/prometheus/procfs v0.0.8/go.mod h1:7Qr8sr6344vo1JqZ6HhLceV9o3AJ1Ff+GxbHq6oeK9A=
 github.com/redis/go-redis/v9 v9.7.0 h1:HhLSs+B6O021gwzl+locl0zEDnyNkxMtf/Z3NNBMa9E=
 github.com/redis/go-redis/v9 v9.7.0/go.mod h1:f6zhXITC7JUJIlPEiBOTXxJgPLdZcA93GewI7inzyWw=
+github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5 h1:mZHayPoR0lNmnHyvtYjDeq0zlVHn9K/ZXoy17ylucdo=
+github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5/go.mod h1:GEXHk5HgEKCvEIIrSpFI3ozzG5xOKA2DVlEX/gGnewM=
 github.com/rogpeppe/go-internal v1.9.0 h1:73kH8U+JUqXU8lRuOHeVHaa/SZPifC7BkcraZVejAe8=
 github.com/rogpeppe/go-internal v1.9.0/go.mod h1:WtVeX8xhTBvf0smdhujwtBcq4Qrzq/fJaraNFVN+nFs=
 github.com/ryanuber/columnize v0.0.0-20160712163229-9b3edd62028f/go.mod h1:sm1tb6uqfes/u+d4ooFouqFdy9/2g9QGwK3SQygK0Ts=
@@ -280,6 +288,7 @@ github.com/spf13/pflag v1.0.5 h1:iy+VFUOCP1a+8yFto/drg2CJ5u0yRoB7fZw3DKv/JXA=
 github.com/spf13/pflag v1.0.5/go.mod h1:McXfInJRrz4CZXVZOBLb0bTZqETkiAhM9Iw0y3An2Bg=
 github.com/spf13/viper v1.19.0 h1:RWq5SEjt8o25SROyN3z2OrDB9l7RPd3lwTWU8EcEdcI=
 github.com/spf13/viper v1.19.0/go.mod h1:GQUN9bilAbhU/jgc1bKs99f/suXKeUMct8Adx5+Ntkg=
+github.com/spkg/bom v0.0.0-20160624110644-59b7046e48ad/go.mod h1:qLr4V1qq6nMqFKkMo8ZTx3f+BZEkzsRUY10Xsm2mwU0=
 github.com/stretchr/objx v0.1.0/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
 github.com/stretchr/objx v0.1.1/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
 github.com/stretchr/objx v0.4.0/go.mod h1:YvHI0jy2hoMjB+UWwv71VJQ9isScKT/TqJzVSSt89Yw=
@@ -305,8 +314,6 @@ github.com/twitchyliquid64/golang-asm v0.15.1 h1:SU5vSMR7hnwNxj24w34ZyCi/FmDZTkS
 github.com/twitchyliquid64/golang-asm v0.15.1/go.mod h1:a1lVb/DtPvCB8fslRZhAngC2+aY1QWCk3Cedj/Gdt08=
 github.com/ugorji/go/codec v1.2.12 h1:9LC83zGrHhuUA9l16C9AHXAqEV/2wBQ4nkvumAE65EE=
 github.com/ugorji/go/codec v1.2.12/go.mod h1:UNopzCgEMSXjBc6AOMqYvWC1ktqTAfzJZUZgYf6w6lg=
-github.com/x-cray/logrus-prefixed-formatter v0.5.2 h1:00txxvfBM9muc0jiLIEAkAcIMJzfthRT6usrui8uGmg=
-github.com/x-cray/logrus-prefixed-formatter v0.5.2/go.mod h1:2duySbKsL6M18s5GU7VPsoEPHyzalCE06qoARUCeBBE=
 github.com/yuin/goldmark v1.1.27/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
 github.com/yuin/goldmark v1.2.1/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
 go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0 h1:0nTRpaCaILLdooXAQnfktlL6Zw1ECKEW9DZGH2byi2c=
@@ -400,14 +407,11 @@ golang.org/x/sys v0.0.0-20211025201205-69cdffdb9359/go.mod h1:oPkhp1MJrh7nUepCBc
 golang.org/x/sys v0.0.0-20220503163025-988cb79eb6c6/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20220728004956-3c1f35247d10/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
-golang.org/x/sys v0.0.0-20220908164124-27713097b956/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.5.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.6.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.28.0 h1:Fksou7UEQUWlKvIdsqzJmUmCX3cZuD2+P3XyyzwMhlA=
 golang.org/x/sys v0.28.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
 golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1/go.mod h1:bj7SfCRtBDWHUb9snDiAeCFNEtKQo2Wmx5Cou7ajbmo=
-golang.org/x/term v0.27.0 h1:WP60Sv1nlK1T6SupCHbXzSaN0b9wUmsPoRS9b61A23Q=
-golang.org/x/term v0.27.0/go.mod h1:iMsnZpn0cago0GOrHO2+Y7u7JPn5AylBrcoWkElMTSM=
 golang.org/x/text v0.3.0/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
 golang.org/x/text v0.3.2/go.mod h1:bEr9sfX3Q8Zfm5fL9x+3itogRgK3+ptLWKqgva+5dAk=
 golang.org/x/text v0.3.3/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
@@ -453,8 +457,6 @@ gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 h1:YR8cESwS4TdDjEe65xsg0ogR
 gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
 gopkg.in/ini.v1 v1.67.0 h1:Dgnx+6+nfE+IfzjUEISNeydPJh9AXNNsWbGP9KzCsOA=
 gopkg.in/ini.v1 v1.67.0/go.mod h1:pNLf8WUiyNEtQjuu5G5vTm06TEv9tsIgeAvK8hOrP4k=
-gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 h1:uRGJdciOHaEIrze2W8Q3AKkepLTh2hOroT7a+7czfdQ=
-gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7/go.mod h1:dt/ZhP58zS4L8KSrWDmTeBkI65Dw0HsyUHuEVlX15mw=
 gopkg.in/yaml.v2 v2.2.1/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.2/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.4/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
diff --git a/internal/stock/service/application.go b/internal/stock/service/application.go
index 5feb83a..e0eace6 100644
--- a/internal/stock/service/application.go
+++ b/internal/stock/service/application.go
@@ -16,14 +16,13 @@ func NewApplication(_ context.Context) app.Application {
 	//stockRepo := adapters.NewMemoryStockRepository()
 	db := persistent.NewMySQL()
 	stockRepo := adapters.NewMySQLStockRepository(db)
-	logger := logrus.NewEntry(logrus.StandardLogger())
 	stripeAPI := integration.NewStripeAPI()
 	metricsClient := metrics.TodoMetrics{}
 	return app.Application{
 		Commands: app.Commands{},
 		Queries: app.Queries{
-			CheckIfItemsInStock: query.NewCheckIfItemsInStockHandler(stockRepo, stripeAPI, logger, metricsClient),
-			GetItems:            query.NewGetItemsHandler(stockRepo, logger, metricsClient),
+			CheckIfItemsInStock: query.NewCheckIfItemsInStockHandler(stockRepo, stripeAPI, logrus.StandardLogger(), metricsClient),
+			GetItems:            query.NewGetItemsHandler(stockRepo, logrus.StandardLogger(), metricsClient),
 		},
 	}
 }
~~~
