# Commit Diff Report

- Repo: go_shop_second
- Sequence: 004 / 10
- Commit: 49bfa8eb40f86817e6bb4f6af9d9f628070e4660
- ShortCommit: 49bfa8e
- Parent: 0059266defb11c717d031d510839a40c36cbd0ee
- Subject: add app to servers
- Author: ghost-yu <hgfhgfhgfhgfhgfhgf@yeah.net>
- Date: 2024-10-14 00:41:57 +0800
- GeneratedAt: 2026-04-06 17:43:36 +08:00

## Short Summary

~~~text
 7 files changed, 100 insertions(+), 8 deletions(-)
~~~

## File Stats

~~~text
 internal/stock/.air.toml              | 52 +++++++++++++++++++++++++++++++++++
 internal/stock/app/app.go             | 10 +++++++
 internal/stock/go.mod                 |  5 ++--
 internal/stock/go.sum                 |  5 ++--
 internal/stock/main.go                | 19 ++++++++++++-
 internal/stock/ports/grpc.go          |  6 ++--
 internal/stock/service/application.go | 11 ++++++++
 7 files changed, 100 insertions(+), 8 deletions(-)
~~~

## Changed Files

~~~text
internal/stock/.air.toml
internal/stock/app/app.go
internal/stock/go.mod
internal/stock/go.sum
internal/stock/main.go
internal/stock/ports/grpc.go
internal/stock/service/application.go
~~~

## Focus Files (Excluded: go.mod / go.sum)

~~~text
internal/stock/.air.toml
internal/stock/app/app.go
internal/stock/main.go
internal/stock/ports/grpc.go
internal/stock/service/application.go
~~~

## Patch

~~~diff
diff --git a/internal/stock/.air.toml b/internal/stock/.air.toml
new file mode 100644
index 0000000..92667de
--- /dev/null
+++ b/internal/stock/.air.toml
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
diff --git a/internal/stock/app/app.go b/internal/stock/app/app.go
new file mode 100644
index 0000000..42330cd
--- /dev/null
+++ b/internal/stock/app/app.go
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
diff --git a/internal/stock/go.mod b/internal/stock/go.mod
index 4602dec..83d4ef8 100644
--- a/internal/stock/go.mod
+++ b/internal/stock/go.mod
@@ -7,6 +7,7 @@ replace github.com/ghost-yu/go_shop_second/common => ../common
 require (
 	github.com/ghost-yu/go_shop_second/common v0.0.0-00010101000000-000000000000
 	github.com/spf13/viper v1.19.0
+	google.golang.org/grpc v1.62.1
 )
 
 require (
@@ -27,7 +28,6 @@ require (
 	github.com/hashicorp/hcl v1.0.0 // indirect
 	github.com/json-iterator/go v1.1.12 // indirect
 	github.com/klauspost/cpuid/v2 v2.2.7 // indirect
-	github.com/konsorten/go-windows-terminal-sequences v1.0.1 // indirect
 	github.com/leodido/go-urn v1.4.0 // indirect
 	github.com/magiconair/properties v1.8.7 // indirect
 	github.com/mattn/go-isatty v0.0.20 // indirect
@@ -37,7 +37,7 @@ require (
 	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
 	github.com/sagikazarmark/locafero v0.4.0 // indirect
 	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
-	github.com/sirupsen/logrus v1.4.2 // indirect
+	github.com/sirupsen/logrus v1.8.1 // indirect
 	github.com/sourcegraph/conc v0.3.0 // indirect
 	github.com/spf13/afero v1.11.0 // indirect
 	github.com/spf13/cast v1.6.0 // indirect
@@ -54,7 +54,6 @@ require (
 	golang.org/x/sys v0.20.0 // indirect
 	golang.org/x/text v0.15.0 // indirect
 	google.golang.org/genproto/googleapis/rpc v0.0.0-20240314234333-6e1732d8331c // indirect
-	google.golang.org/grpc v1.62.1 // indirect
 	google.golang.org/protobuf v1.34.1 // indirect
 	gopkg.in/ini.v1 v1.67.0 // indirect
 	gopkg.in/yaml.v3 v3.0.1 // indirect
diff --git a/internal/stock/go.sum b/internal/stock/go.sum
index 5cf423c..b380353 100644
--- a/internal/stock/go.sum
+++ b/internal/stock/go.sum
@@ -70,7 +70,6 @@ github.com/klauspost/cpuid/v2 v2.0.9/go.mod h1:FInQzS24/EEf25PyTYn52gqo7WaD8xa02
 github.com/klauspost/cpuid/v2 v2.2.7 h1:ZWSB3igEs+d0qvnxR/ZBzXVmxkgt8DdzP6m9pfuVLDM=
 github.com/klauspost/cpuid/v2 v2.2.7/go.mod h1:Lcz8mBdAVJIBVzewtcLocK12l3Y+JytZYpaMropDUws=
 github.com/knz/go-libedit v1.10.1/go.mod h1:MZTVkCWyz0oBc7JOWP3wNAzd002ZbM/5hgShxwh4x8M=
-github.com/konsorten/go-windows-terminal-sequences v1.0.1 h1:mweAR1A6xJ3oS2pRaGiHgQ4OO8tzTaLawm8vnODuwDk=
 github.com/konsorten/go-windows-terminal-sequences v1.0.1/go.mod h1:T0+1ngSBFLxvqU3pZ+m/2kptfBszLMUkC4ZK/EgS/cQ=
 github.com/kr/pretty v0.1.0/go.mod h1:dAy3ld7l9f0ibDNOQOHHMYYIIbhfbHSm3C4ZsoJORNo=
 github.com/kr/pretty v0.3.1 h1:flRD4NNwYAUpkphVc1HcthR4KEIFJ65n8Mw5qdRn3LE=
@@ -106,8 +105,9 @@ github.com/sagikazarmark/locafero v0.4.0 h1:HApY1R9zGo4DBgr7dqsTH/JJxLTTsOt7u6ke
 github.com/sagikazarmark/locafero v0.4.0/go.mod h1:Pe1W6UlPYUk/+wc/6KFhbORCfqzgYEpgQ3O5fPuL3H4=
 github.com/sagikazarmark/slog-shim v0.1.0 h1:diDBnUNK9N/354PgrxMywXnAwEr1QZcOr6gto+ugjYE=
 github.com/sagikazarmark/slog-shim v0.1.0/go.mod h1:SrcSrq8aKtyuqEI1uvTDTK1arOWRIczQRv+GVI1AkeQ=
-github.com/sirupsen/logrus v1.4.2 h1:SPIRibHv4MatM3XXNO2BJeFLZwZ2LvZgfQ5+UNI2im4=
 github.com/sirupsen/logrus v1.4.2/go.mod h1:tLMulIdttU9McNUspp0xgXVQah82FyeX6MwdIuYE2rE=
+github.com/sirupsen/logrus v1.8.1 h1:dJKuHgqk1NNQlqoA6BTlM1Wf9DOH3NBjQyu0h9+AZZE=
+github.com/sirupsen/logrus v1.8.1/go.mod h1:yWOB1SBYBC5VeMP7gHvWumXLIWorT60ONWic61uBYv0=
 github.com/sourcegraph/conc v0.3.0 h1:OQTbbt6P72L20UqAkXXuLOj79LfEanQ+YQFNpLA9ySo=
 github.com/sourcegraph/conc v0.3.0/go.mod h1:Sdozi7LEKbFPqYX2/J+iBAM6HpqSLTASQIKqDmF7Mt0=
 github.com/spf13/afero v1.11.0 h1:WJQKhtpdm3v2IzqG8VMqrr6Rf3UYpEF239Jy9wNepM8=
@@ -186,6 +186,7 @@ golang.org/x/sys v0.0.0-20180830151530-49385e6e1522/go.mod h1:STP8DvDyc/dI5b8T5h
 golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
 golang.org/x/sys v0.0.0-20190412213103-97732733099d/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20190422165155-953cdadca894/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20191026070338-33540a1f6037/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20200930185726-fdedc70b468f/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20211025201205-69cdffdb9359/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.5.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
diff --git a/internal/stock/main.go b/internal/stock/main.go
index 302ddd6..78ca718 100644
--- a/internal/stock/main.go
+++ b/internal/stock/main.go
@@ -1,21 +1,38 @@
 package main
 
 import (
+	"context"
+
+	"github.com/ghost-yu/go_shop_second/common/config"
 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
 	"github.com/ghost-yu/go_shop_second/common/server"
 	"github.com/ghost-yu/go_shop_second/stock/ports"
+	"github.com/ghost-yu/go_shop_second/stock/service"
+	"github.com/sirupsen/logrus"
 	"github.com/spf13/viper"
 	"google.golang.org/grpc"
 )
 
+func init() {
+	if err := config.NewViperConfig(); err != nil {
+		logrus.Fatal(err)
+	}
+}
+
 func main() {
 	serviceName := viper.GetString("stock.service-name")
 	serverType := viper.GetString("stock.server-to-run")
 
+	logrus.Info(serverType)
+
+	ctx, cancel := context.WithCancel(context.Background())
+	defer cancel()
+
+	application := service.NewApplication(ctx)
 	switch serverType {
 	case "grpc":
 		server.RunGRPCServer(serviceName, func(server *grpc.Server) {
-			svc := ports.NewGRPCServer()
+			svc := ports.NewGRPCServer(application)
 			stockpb.RegisterStockServiceServer(server, svc)
 		})
 	case "http":
diff --git a/internal/stock/ports/grpc.go b/internal/stock/ports/grpc.go
index 63863a2..fb41c1f 100644
--- a/internal/stock/ports/grpc.go
+++ b/internal/stock/ports/grpc.go
@@ -4,13 +4,15 @@ import (
 	context "context"
 
 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
+	"github.com/ghost-yu/go_shop_second/stock/app"
 )
 
 type GRPCServer struct {
+	app app.Application
 }
 
-func NewGRPCServer() *GRPCServer {
-	return &GRPCServer{}
+func NewGRPCServer(app app.Application) *GRPCServer {
+	return &GRPCServer{app: app}
 }
 
 func (G GRPCServer) GetItems(ctx context.Context, request *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {
diff --git a/internal/stock/service/application.go b/internal/stock/service/application.go
new file mode 100644
index 0000000..423b368
--- /dev/null
+++ b/internal/stock/service/application.go
@@ -0,0 +1,11 @@
+package service
+
+import (
+	"context"
+
+	"github.com/ghost-yu/go_shop_second/stock/app"
+)
+
+func NewApplication(ctx context.Context) app.Application {
+	return app.Application{}
+}
~~~
