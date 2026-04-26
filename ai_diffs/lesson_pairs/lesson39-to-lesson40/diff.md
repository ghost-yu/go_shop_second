# Lesson Pair Diff Report

- FromBranch: lesson39
- ToBranch: lesson40

## Short Summary

~~~text
 7 files changed, 112 insertions(+), 18 deletions(-)
~~~

## File Stats

~~~text
 internal/common/decorator/command.go  |  4 ++--
 internal/common/decorator/logging.go  | 26 ++++++++++++++++++++-
 internal/common/decorator/metrics.go  | 20 ++++++++++++++++
 internal/common/middleware/logger.go  | 20 +++++++---------
 internal/common/middleware/request.go | 44 +++++++++++++++++++++++++++++++++++
 internal/common/response.go           | 15 ++++++++----
 internal/common/server/http.go        |  1 +
 7 files changed, 112 insertions(+), 18 deletions(-)
~~~

## Commit Comparison

~~~text
> d080dd9 log update
~~~

## Changed Files

~~~text
internal/common/decorator/command.go
internal/common/decorator/logging.go
internal/common/decorator/metrics.go
internal/common/middleware/logger.go
internal/common/middleware/request.go
internal/common/response.go
internal/common/server/http.go
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
internal/common/decorator/command.go
internal/common/decorator/logging.go
internal/common/decorator/metrics.go
internal/common/middleware/logger.go
internal/common/middleware/request.go
internal/common/response.go
internal/common/server/http.go
~~~

## Full Diff

~~~diff
diff --git a/internal/common/decorator/command.go b/internal/common/decorator/command.go
index 53f8c37..b89552d 100644
--- a/internal/common/decorator/command.go
+++ b/internal/common/decorator/command.go
@@ -11,9 +11,9 @@ type CommandHandler[C, R any] interface {
 }
 
 func ApplyCommandDecorators[C, R any](handler CommandHandler[C, R], logger *logrus.Entry, metricsClient MetricsClient) CommandHandler[C, R] {
-	return queryLoggingDecorator[C, R]{
+	return commandLoggingDecorator[C, R]{
 		logger: logger,
-		base: queryMetricsDecorator[C, R]{
+		base: commandMetricsDecorator[C, R]{
 			base:   handler,
 			client: metricsClient,
 		},
diff --git a/internal/common/decorator/logging.go b/internal/common/decorator/logging.go
index a9a2325..dbe402c 100644
--- a/internal/common/decorator/logging.go
+++ b/internal/common/decorator/logging.go
@@ -2,6 +2,7 @@ package decorator
 
 import (
 	"context"
+	"encoding/json"
 	"fmt"
 	"strings"
 
@@ -14,9 +15,10 @@ type queryLoggingDecorator[C, R any] struct {
 }
 
 func (q queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
+	body, _ := json.Marshal(cmd)
 	logger := q.logger.WithFields(logrus.Fields{
 		"query":      generateActionName(cmd),
-		"query_body": fmt.Sprintf("%#v", cmd),
+		"query_body": string(body),
 	})
 	logger.Debug("Executing query")
 	defer func() {
@@ -29,6 +31,28 @@ func (q queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result
 	return q.base.Handle(ctx, cmd)
 }
 
+type commandLoggingDecorator[C, R any] struct {
+	logger *logrus.Entry
+	base   CommandHandler[C, R]
+}
+
+func (q commandLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
+	body, _ := json.Marshal(cmd)
+	logger := q.logger.WithFields(logrus.Fields{
+		"command":      generateActionName(cmd),
+		"command_body": string(body),
+	})
+	logger.Debug("Executing command")
+	defer func() {
+		if err == nil {
+			logger.Info("Command execute successfully")
+		} else {
+			logger.Error("Failed to execute command", err)
+		}
+	}()
+	return q.base.Handle(ctx, cmd)
+}
+
 func generateActionName(cmd any) string {
 	return strings.Split(fmt.Sprintf("%T", cmd), ".")[1]
 }
diff --git a/internal/common/decorator/metrics.go b/internal/common/decorator/metrics.go
index fba9a77..41db94d 100644
--- a/internal/common/decorator/metrics.go
+++ b/internal/common/decorator/metrics.go
@@ -30,3 +30,23 @@ func (q queryMetricsDecorator[C, R]) Handle(ctx context.Context, cmd C) (result
 	}()
 	return q.base.Handle(ctx, cmd)
 }
+
+type commandMetricsDecorator[C, R any] struct {
+	base   CommandHandler[C, R]
+	client MetricsClient
+}
+
+func (q commandMetricsDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
+	start := time.Now()
+	actionName := strings.ToLower(generateActionName(cmd))
+	defer func() {
+		end := time.Since(start)
+		q.client.Inc(fmt.Sprintf("command.%s.duration", actionName), int(end.Seconds()))
+		if err == nil {
+			q.client.Inc(fmt.Sprintf("command.%s.success", actionName), 1)
+		} else {
+			q.client.Inc(fmt.Sprintf("command.%s.failure", actionName), 1)
+		}
+	}()
+	return q.base.Handle(ctx, cmd)
+}
diff --git a/internal/common/middleware/logger.go b/internal/common/middleware/logger.go
index 1a7e1c5..731a61d 100644
--- a/internal/common/middleware/logger.go
+++ b/internal/common/middleware/logger.go
@@ -1,23 +1,21 @@
 package middleware
 
 import (
-	"time"
-
 	"github.com/gin-gonic/gin"
 	"github.com/sirupsen/logrus"
 )
 
 func StructuredLog(l *logrus.Entry) gin.HandlerFunc {
 	return func(c *gin.Context) {
-		t := time.Now()
+		//t := time.Now()
 		c.Next()
-		elapsed := time.Since(t)
-		l.WithFields(logrus.Fields{
-			"time_elapsed_ms": elapsed.Milliseconds(),
-			"request_uri":     c.Request.RequestURI,
-			"remote_addr":     c.RemoteIP(),
-			"client_ip":       c.ClientIP(),
-			"full_path":       c.FullPath(),
-		}).Info("request_out")
+		//elapsed := time.Since(t)
+		//l.WithFields(logrus.Fields{
+		//	"time_elapsed_ms": elapsed.Milliseconds(),
+		//	"request_uri":     c.Request.RequestURI,
+		//	"remote_addr":     c.RemoteIP(),
+		//	"client_ip":       c.ClientIP(),
+		//	"full_path":       c.FullPath(),
+		//}).Info("request_out")
 	}
 }
diff --git a/internal/common/middleware/request.go b/internal/common/middleware/request.go
new file mode 100644
index 0000000..3b90df0
--- /dev/null
+++ b/internal/common/middleware/request.go
@@ -0,0 +1,44 @@
+package middleware
+
+import (
+	"bytes"
+	"encoding/json"
+	"io"
+	"time"
+
+	"github.com/gin-gonic/gin"
+	"github.com/sirupsen/logrus"
+)
+
+func RequestLog(l *logrus.Entry) gin.HandlerFunc {
+	return func(c *gin.Context) {
+		requestIn(c, l)
+		defer requestOut(c, l)
+		c.Next()
+	}
+}
+
+func requestOut(c *gin.Context, l *logrus.Entry) {
+	response, _ := c.Get("response")
+	start, _ := c.Get("request_start")
+	startTime := start.(time.Time)
+	l.WithContext(c.Request.Context()).WithFields(logrus.Fields{
+		"proc_time_ms": time.Since(startTime).Milliseconds(),
+		"response":     response,
+	}).Info("__request_out")
+}
+
+func requestIn(c *gin.Context, l *logrus.Entry) {
+	c.Set("request_start", time.Now())
+	body := c.Request.Body
+	bodyBytes, _ := io.ReadAll(body)
+	c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
+	var compactJson bytes.Buffer
+	_ = json.Compact(&compactJson, bodyBytes)
+	l.WithContext(c.Request.Context()).WithFields(logrus.Fields{
+		"start": time.Now().Unix(),
+		"args":  compactJson.String(),
+		"from":  c.RemoteIP(),
+		"uri":   c.Request.RequestURI,
+	}).Info("__request_in")
+}
diff --git a/internal/common/response.go b/internal/common/response.go
index 7be1fe7..4abf56e 100644
--- a/internal/common/response.go
+++ b/internal/common/response.go
@@ -1,6 +1,7 @@
 package common
 
 import (
+	"encoding/json"
 	"net/http"
 
 	"github.com/ghost-yu/go_shop_second/common/tracing"
@@ -25,19 +26,25 @@ func (base *BaseResponse) Response(c *gin.Context, err error, data interface{})
 }
 
 func (base *BaseResponse) success(c *gin.Context, data interface{}) {
-	c.JSON(http.StatusOK, response{
+	r := response{
 		Errno:   0,
 		Message: "success",
 		Data:    data,
 		TraceID: tracing.TraceID(c.Request.Context()),
-	})
+	}
+	resp, _ := json.Marshal(r)
+	c.Set("response", string(resp))
+	c.JSON(http.StatusOK, r)
 }
 
 func (base *BaseResponse) error(c *gin.Context, err error) {
-	c.JSON(http.StatusOK, response{
+	r := response{
 		Errno:   2,
 		Message: err.Error(),
 		Data:    nil,
 		TraceID: tracing.TraceID(c.Request.Context()),
-	})
+	}
+	resp, _ := json.Marshal(r)
+	c.Set("response", string(resp))
+	c.JSON(http.StatusOK, r)
 }
diff --git a/internal/common/server/http.go b/internal/common/server/http.go
index 584f5bb..5f280f4 100644
--- a/internal/common/server/http.go
+++ b/internal/common/server/http.go
@@ -29,5 +29,6 @@ func RunHTTPServerOnAddr(addr string, wrapper func(router *gin.Engine)) {
 func setMiddlewares(r *gin.Engine) {
 	r.Use(middleware.StructuredLog(logrus.NewEntry(logrus.StandardLogger())))
 	r.Use(gin.Recovery())
+	r.Use(middleware.RequestLog(logrus.NewEntry(logrus.StandardLogger())))
 	r.Use(otelgin.Middleware("default_server"))
 }
~~~
