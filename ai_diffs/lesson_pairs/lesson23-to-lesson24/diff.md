# Lesson Pair Diff Report

- FromBranch: lesson23
- ToBranch: lesson24

## Short Summary

~~~text
 2 files changed, 49 insertions(+)
~~~

## File Stats

~~~text
 internal/common/client/grpc.go     | 48 ++++++++++++++++++++++++++++++++++++++
 internal/common/config/global.yaml |  1 +
 2 files changed, 49 insertions(+)
~~~

## Commit Comparison

~~~text
> fb48e4b wait for
~~~

## Changed Files

~~~text
internal/common/client/grpc.go
internal/common/config/global.yaml
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
internal/common/client/grpc.go
internal/common/config/global.yaml
~~~

## Full Diff

~~~diff
diff --git a/internal/common/client/grpc.go b/internal/common/client/grpc.go
index bbf0c76..f95506b 100644
--- a/internal/common/client/grpc.go
+++ b/internal/common/client/grpc.go
@@ -2,6 +2,9 @@ package client
 
 import (
 	"context"
+	"errors"
+	"net"
+	"time"
 
 	"github.com/ghost-yu/go_shop_second/common/discovery"
 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
@@ -13,6 +16,9 @@ import (
 )
 
 func NewStockGRPCClient(ctx context.Context) (client stockpb.StockServiceClient, close func() error, err error) {
+	if !WaitForStockGRPCClient(viper.GetDuration("dial-grpc-timeout") * time.Second) {
+		return nil, nil, errors.New("stock grpc not available")
+	}
 	grpcAddr, err := discovery.GetServiceAddr(ctx, viper.GetString("stock.service-name"))
 	if err != nil {
 		return nil, func() error { return nil }, err
@@ -29,6 +35,9 @@ func NewStockGRPCClient(ctx context.Context) (client stockpb.StockServiceClient,
 }
 
 func NewOrderGRPCClient(ctx context.Context) (client orderpb.OrderServiceClient, close func() error, err error) {
+	if !WaitForOrderGRPCClient(viper.GetDuration("dial-grpc-timeout") * time.Second) {
+		return nil, nil, errors.New("order grpc not available")
+	}
 	grpcAddr, err := discovery.GetServiceAddr(ctx, viper.GetString("order.service-name"))
 	if err != nil {
 		return nil, func() error { return nil }, err
@@ -50,3 +59,42 @@ func grpcDialOpts(_ string) []grpc.DialOption {
 		grpc.WithTransportCredentials(insecure.NewCredentials()),
 	}
 }
+
+func WaitForOrderGRPCClient(timeout time.Duration) bool {
+	logrus.Infof("waiting for order grpc client, timeout: %v seconds", timeout.Seconds())
+	return waitFor(viper.GetString("order.grpc-addr"), timeout)
+}
+
+func WaitForStockGRPCClient(timeout time.Duration) bool {
+	logrus.Infof("waiting for stock grpc client, timeout: %v seconds", timeout.Seconds())
+	return waitFor(viper.GetString("stock.grpc-addr"), timeout)
+}
+
+func waitFor(addr string, timeout time.Duration) bool {
+	portAvailable := make(chan struct{})
+	timeoutCh := time.After(timeout)
+
+	go func() {
+		for {
+			select {
+			case <-timeoutCh:
+				return
+			default:
+				// continue
+			}
+			_, err := net.Dial("tcp", addr)
+			if err == nil {
+				close(portAvailable)
+				return
+			}
+			time.Sleep(200 * time.Millisecond)
+		}
+	}()
+
+	select {
+	case <-portAvailable:
+		return true
+	case <-timeoutCh:
+		return false
+	}
+}
diff --git a/internal/common/config/global.yaml b/internal/common/config/global.yaml
index 7dd93c9..1bf2f6e 100644
--- a/internal/common/config/global.yaml
+++ b/internal/common/config/global.yaml
@@ -1,4 +1,5 @@
 fallback-grpc-addr: 127.0.0.1:3030
+dial-grpc-timeout: 10
 
 consul:
   addr: 127.0.0.1:8500
~~~
