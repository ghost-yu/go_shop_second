# Lesson Pair Diff Report

- FromBranch: lesson22
- ToBranch: lesson23

## Short Summary

~~~text
 14 files changed, 193 insertions(+), 26 deletions(-)
~~~

## File Stats

~~~text
 .golangci.yaml                                     | 109 +++++++++++++++++++++
 Makefile                                           |  10 +-
 internal/common/client/grpc.go                     |  15 +--
 internal/common/genproto/orderpb/order.pb.go       |   5 +-
 internal/common/genproto/orderpb/order_grpc.pb.go  |   1 +
 internal/common/genproto/stockpb/stock.pb.go       |   5 +-
 internal/common/genproto/stockpb/stock_grpc.pb.go  |   1 +
 internal/common/server/http.go                     |   2 +-
 internal/order/infrastructure/consumer/consumer.go |   4 +-
 internal/order/ports/grpc.go                       |   4 +-
 .../payment/infrastructure/consumer/consumer.go    |   4 +-
 internal/payment/service/application.go            |   2 +-
 .../stock/app/query/check_if_items_in_stock.go     |   6 +-
 scripts/lint.sh                                    |  51 ++++++++++
 14 files changed, 193 insertions(+), 26 deletions(-)
~~~

## Commit Comparison

~~~text
> 8311ea7 lint
~~~

## Changed Files

~~~text
.golangci.yaml
Makefile
internal/common/client/grpc.go
internal/common/genproto/orderpb/order.pb.go
internal/common/genproto/orderpb/order_grpc.pb.go
internal/common/genproto/stockpb/stock.pb.go
internal/common/genproto/stockpb/stock_grpc.pb.go
internal/common/server/http.go
internal/order/infrastructure/consumer/consumer.go
internal/order/ports/grpc.go
internal/payment/infrastructure/consumer/consumer.go
internal/payment/service/application.go
internal/stock/app/query/check_if_items_in_stock.go
scripts/lint.sh
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
.golangci.yaml
Makefile
internal/common/client/grpc.go
internal/common/server/http.go
internal/order/infrastructure/consumer/consumer.go
internal/order/ports/grpc.go
internal/payment/infrastructure/consumer/consumer.go
internal/payment/service/application.go
internal/stock/app/query/check_if_items_in_stock.go
scripts/lint.sh
~~~

## Full Diff

~~~diff
diff --git a/.golangci.yaml b/.golangci.yaml
new file mode 100644
index 0000000..9bd7f8f
--- /dev/null
+++ b/.golangci.yaml
@@ -0,0 +1,109 @@
+---
+run:
+  timeout: 30m
+  tests: false
+issues:
+  max-same-issues: 0
+  # Excluding configuration per-path, per-linter, per-text and per-source
+  exclude-rules:
+    # exclude ineffassing linter for generated files for conversion
+    - path: conversion\.go
+      linters: [ineffassign]
+  exclude-files:
+    - ^zz_generated.*
+linters:
+  disable-all: true
+  enable: # please keep this alphabetized
+    # Don't use soon to deprecated[1] linters that lead to false
+    # https://github.com/golangci/golangci-lint/issues/1841
+    # - deadcode
+    # - structcheck
+    # - varcheck
+    - goimports
+    - ineffassign
+    - nakedret
+    - revive
+    - staticcheck
+    - stylecheck
+    - unconvert # Remove unnecessary type conversions
+    - unparam
+    - unused
+linters-settings: # please keep this alphabetized
+  #  goimports:
+  #    local-prefixes:   github.com/Nicknamezz00/gorder-pay/internal
+  nakedret:
+    # Align with https://github.com/alexkohler/nakedret/blob/v1.0.2/cmd/nakedret/main.go#L10
+    max-func-lines: 5
+  revive:
+    ignore-generated-header: false
+    severity: error
+    confidence: 0.8
+    enable-all-rules: false
+    rules:
+      - name: blank-imports
+        severity: error
+        disabled: false
+      - name: context-as-argument
+        severity: error
+        disabled: false
+      - name: dot-imports
+        severity: error
+        disabled: false
+      - name: error-return
+        severity: error
+        disabled: false
+      - name: error-naming
+        severity: error
+        disabled: false
+      - name: if-return
+        severity: error
+        disabled: false
+      - name: increment-decrement
+        severity: error
+        disabled: false
+      - name: var-declaration
+        severity: error
+        disabled: false
+      - name: package-comments
+        severity: error
+        disabled: false
+      - name: range
+        severity: error
+        disabled: false
+      - name: receiver-naming
+        severity: error
+        disabled: false
+      - name: time-naming
+        severity: error
+        disabled: false
+      - name: indent-error-flow
+        severity: error
+        disabled: false
+      - name: errorf
+        severity: error
+        disabled: false
+      - name: context-keys-type
+        severity: error
+        disabled: false
+      - name: error-strings
+        severity: error
+        disabled: false
+      - name: var-naming
+        disabled: false
+        arguments:
+          # The following is the configuration for var-naming rule, the first element is the allow list and the second element is the deny list.
+          - [] # AllowList: leave it empty to use the default (empty, too). This means that we're not relaxing the rule in any way, i.e. elementId will raise a violation, it should be elementID, refer to the next line to see the list of denied initialisms.
+          - ["GRPC", "WAL"] # DenyList: Add GRPC and WAL to strict the rule not allowing instances like Wal or Grpc. The default values are located at commonInitialisms, refer to: https://github.com/mgechev/revive/blob/v1.3.7/lint/utils.go#L93-L133.
+      # TODO: enable the following rules
+      - name: exported
+        disabled: true
+      - name: unexported-return
+        disabled: true
+  staticcheck:
+    checks:
+      - all
+      - -SA1019 # TODO(fix) Using a deprecated function, variable, constant or field
+      - -SA2002 # TODO(fix) Called testing.T.FailNow or SkipNow in a goroutine, which isn’t allowed
+  stylecheck:
+    checks:
+      - ST1019 # Importing the same package multiple times.
\ No newline at end of file
diff --git a/Makefile b/Makefile
index 6faaf1e..c07c1ad 100644
--- a/Makefile
+++ b/Makefile
@@ -7,4 +7,12 @@ genproto:
 
 .PHONY: genopenapi
 genopenapi:
-	@./scripts/genopenapi.sh
\ No newline at end of file
+	@./scripts/genopenapi.sh
+
+.PHONY: fmt
+fmt:
+	goimports -l -w internal/
+
+.PHONY: lint
+lint:
+	@./scripts/lint.sh
\ No newline at end of file
diff --git a/internal/common/client/grpc.go b/internal/common/client/grpc.go
index b6e6bab..bbf0c76 100644
--- a/internal/common/client/grpc.go
+++ b/internal/common/client/grpc.go
@@ -20,10 +20,7 @@ func NewStockGRPCClient(ctx context.Context) (client stockpb.StockServiceClient,
 	if grpcAddr == "" {
 		logrus.Warn("empty grpc addr for stock grpc")
 	}
-	opts, err := grpcDialOpts(grpcAddr)
-	if err != nil {
-		return nil, func() error { return nil }, err
-	}
+	opts := grpcDialOpts(grpcAddr)
 	conn, err := grpc.NewClient(grpcAddr, opts...)
 	if err != nil {
 		return nil, func() error { return nil }, err
@@ -39,10 +36,8 @@ func NewOrderGRPCClient(ctx context.Context) (client orderpb.OrderServiceClient,
 	if grpcAddr == "" {
 		logrus.Warn("empty grpc addr for order grpc")
 	}
-	opts, err := grpcDialOpts(grpcAddr)
-	if err != nil {
-		return nil, func() error { return nil }, err
-	}
+	opts := grpcDialOpts(grpcAddr)
+
 	conn, err := grpc.NewClient(grpcAddr, opts...)
 	if err != nil {
 		return nil, func() error { return nil }, err
@@ -50,8 +45,8 @@ func NewOrderGRPCClient(ctx context.Context) (client orderpb.OrderServiceClient,
 	return orderpb.NewOrderServiceClient(conn), conn.Close, nil
 }
 
-func grpcDialOpts(addr string) ([]grpc.DialOption, error) {
+func grpcDialOpts(_ string) []grpc.DialOption {
 	return []grpc.DialOption{
 		grpc.WithTransportCredentials(insecure.NewCredentials()),
-	}, nil
+	}
 }
diff --git a/internal/common/genproto/orderpb/order.pb.go b/internal/common/genproto/orderpb/order.pb.go
index 745cee7..c2975bc 100644
--- a/internal/common/genproto/orderpb/order.pb.go
+++ b/internal/common/genproto/orderpb/order.pb.go
@@ -7,11 +7,12 @@
 package orderpb
 
 import (
+	reflect "reflect"
+	sync "sync"
+
 	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
 	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
 	emptypb "google.golang.org/protobuf/types/known/emptypb"
-	reflect "reflect"
-	sync "sync"
 )
 
 const (
diff --git a/internal/common/genproto/orderpb/order_grpc.pb.go b/internal/common/genproto/orderpb/order_grpc.pb.go
index f18a631..4b0b971 100644
--- a/internal/common/genproto/orderpb/order_grpc.pb.go
+++ b/internal/common/genproto/orderpb/order_grpc.pb.go
@@ -8,6 +8,7 @@ package orderpb
 
 import (
 	context "context"
+
 	grpc "google.golang.org/grpc"
 	codes "google.golang.org/grpc/codes"
 	status "google.golang.org/grpc/status"
diff --git a/internal/common/genproto/stockpb/stock.pb.go b/internal/common/genproto/stockpb/stock.pb.go
index 9566739..38878e4 100644
--- a/internal/common/genproto/stockpb/stock.pb.go
+++ b/internal/common/genproto/stockpb/stock.pb.go
@@ -7,11 +7,12 @@
 package stockpb
 
 import (
+	reflect "reflect"
+	sync "sync"
+
 	orderpb "github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
 	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
 	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
-	reflect "reflect"
-	sync "sync"
 )
 
 const (
diff --git a/internal/common/genproto/stockpb/stock_grpc.pb.go b/internal/common/genproto/stockpb/stock_grpc.pb.go
index 7b3a2a2..f71e2d9 100644
--- a/internal/common/genproto/stockpb/stock_grpc.pb.go
+++ b/internal/common/genproto/stockpb/stock_grpc.pb.go
@@ -8,6 +8,7 @@ package stockpb
 
 import (
 	context "context"
+
 	grpc "google.golang.org/grpc"
 	codes "google.golang.org/grpc/codes"
 	status "google.golang.org/grpc/status"
diff --git a/internal/common/server/http.go b/internal/common/server/http.go
index 7f39359..f22dce4 100644
--- a/internal/common/server/http.go
+++ b/internal/common/server/http.go
@@ -8,7 +8,7 @@ import (
 func RunHTTPServer(serviceName string, wrapper func(router *gin.Engine)) {
 	addr := viper.Sub(serviceName).GetString("http-addr")
 	if addr == "" {
-		// TODO: Warning log
+		panic("empty http address")
 	}
 	RunHTTPServerOnAddr(addr, wrapper)
 }
diff --git a/internal/order/infrastructure/consumer/consumer.go b/internal/order/infrastructure/consumer/consumer.go
index 7e2b8a6..e7ce252 100644
--- a/internal/order/infrastructure/consumer/consumer.go
+++ b/internal/order/infrastructure/consumer/consumer.go
@@ -38,13 +38,13 @@ func (c *Consumer) Listen(ch *amqp.Channel) {
 	var forever chan struct{}
 	go func() {
 		for msg := range msgs {
-			c.handleMessage(msg, q, ch)
+			c.handleMessage(msg)
 		}
 	}()
 	<-forever
 }
 
-func (c *Consumer) handleMessage(msg amqp.Delivery, q amqp.Queue, ch *amqp.Channel) {
+func (c *Consumer) handleMessage(msg amqp.Delivery) {
 	o := &domain.Order{}
 	if err := json.Unmarshal(msg.Body, o); err != nil {
 		logrus.Infof("error unmarshal msg.body into domain.order, err = %v", err)
diff --git a/internal/order/ports/grpc.go b/internal/order/ports/grpc.go
index 0e3805b..534ce55 100644
--- a/internal/order/ports/grpc.go
+++ b/internal/order/ports/grpc.go
@@ -50,7 +50,7 @@ func (G GRPCServer) UpdateOrder(ctx context.Context, request *orderpb.Order) (_
 	order, err := domain.NewOrder(request.ID, request.CustomerID, request.Status, request.PaymentLink, request.Items)
 	if err != nil {
 		err = status.Error(codes.Internal, err.Error())
-		return
+		return nil, err
 	}
 	_, err = G.app.Commands.UpdateOrder.Handle(ctx, command.UpdateOrder{
 		Order: order,
@@ -58,5 +58,5 @@ func (G GRPCServer) UpdateOrder(ctx context.Context, request *orderpb.Order) (_
 			return order, nil
 		},
 	})
-	return
+	return nil, err
 }
diff --git a/internal/payment/infrastructure/consumer/consumer.go b/internal/payment/infrastructure/consumer/consumer.go
index 213a00b..7f01e7a 100644
--- a/internal/payment/infrastructure/consumer/consumer.go
+++ b/internal/payment/infrastructure/consumer/consumer.go
@@ -36,13 +36,13 @@ func (c *Consumer) Listen(ch *amqp.Channel) {
 	var forever chan struct{}
 	go func() {
 		for msg := range msgs {
-			c.handleMessage(msg, q, ch)
+			c.handleMessage(msg, q)
 		}
 	}()
 	<-forever
 }
 
-func (c *Consumer) handleMessage(msg amqp.Delivery, q amqp.Queue, ch *amqp.Channel) {
+func (c *Consumer) handleMessage(msg amqp.Delivery, q amqp.Queue) {
 	logrus.Infof("Payment receive a message from %s, msg=%v", q.Name, string(msg.Body))
 
 	o := &orderpb.Order{}
diff --git a/internal/payment/service/application.go b/internal/payment/service/application.go
index 31b7d69..0199c72 100644
--- a/internal/payment/service/application.go
+++ b/internal/payment/service/application.go
@@ -27,7 +27,7 @@ func NewApplication(ctx context.Context) (app.Application, func()) {
 	}
 }
 
-func newApplication(ctx context.Context, orderGRPC command.OrderService, processor domain.Processor) app.Application {
+func newApplication(_ context.Context, orderGRPC command.OrderService, processor domain.Processor) app.Application {
 	logger := logrus.NewEntry(logrus.StandardLogger())
 	metricClient := metrics.TodoMetrics{}
 	return app.Application{
diff --git a/internal/stock/app/query/check_if_items_in_stock.go b/internal/stock/app/query/check_if_items_in_stock.go
index 4d07201..3d3114e 100644
--- a/internal/stock/app/query/check_if_items_in_stock.go
+++ b/internal/stock/app/query/check_if_items_in_stock.go
@@ -44,14 +44,14 @@ func (h checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfIte
 	var res []*orderpb.Item
 	for _, i := range query.Items {
 		// TODO: 改成从数据库 or stripe 获取
-		priceId, ok := stub[i.ID]
+		priceID, ok := stub[i.ID]
 		if !ok {
-			priceId = stub["1"]
+			priceID = stub["1"]
 		}
 		res = append(res, &orderpb.Item{
 			ID:       i.ID,
 			Quantity: i.Quantity,
-			PriceID:  priceId,
+			PriceID:  priceID,
 		})
 	}
 	return res, nil
diff --git a/scripts/lint.sh b/scripts/lint.sh
new file mode 100755
index 0000000..4a03ab4
--- /dev/null
+++ b/scripts/lint.sh
@@ -0,0 +1,51 @@
+#!/usr/bin/env bash
+
+set -euo pipefail
+
+source ./scripts/lib.sh
+
+function install_if_not_exist() {
+  TOOL_NAME=$1
+  INSTALL_URL=$2
+  if command -v $TOOL_NAME &> /dev/null
+  then
+    log_callout "$TOOL_NAME is already installed."
+  else
+    log_cmd "$TOOL_NAME is not installed. Installing..."
+    run go install "$INSTALL_URL"
+  fi
+}
+
+install_if_not_exist go-cleanarch github.com/roblaszczak/go-cleanarch@latest
+
+readonly LINT_VERSION="1.54.0"
+NEED_INSTALL=false
+if command -v golangci-lint >/dev/null 2>&1; then
+  # golangci-lint has version 1.54.0 built with go1.21.0 from c1d8c565 on 2023-08-09T11:50:00Z
+  CURRENT_VERSION=$(golangci-lint --version | awk '{print $4}' | sed 's/^v//')
+  log_callout "golangci-lint v$CURRENT_VERSION already installed."
+  if [ "$CURRENT_VERSION" == "$LINT_VERSION" ]; then
+    NEED_INSTALL=false
+  else
+    NEED_INSTALL=true
+  fi
+else
+  NEED_INSTALL=true
+fi
+
+if [ "$NEED_INSTALL" == true ]; then
+  run curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.0
+fi
+
+run go-cleanarch
+
+log_info "lint modules:"
+log_info "$(modules)"
+
+run goimports -w -l .
+
+while read -r module; do
+  run cd ./internal/"$module"
+  run golangci-lint run  --config "$ROOT_DIR/.golangci.yaml"
+  run cd -
+done < <(modules)
\ No newline at end of file
~~~
