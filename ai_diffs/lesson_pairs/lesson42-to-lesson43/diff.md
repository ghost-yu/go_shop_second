# Lesson Pair Diff Report

- FromBranch: lesson42
- ToBranch: lesson43

## Short Summary

~~~text
 5 files changed, 128 insertions(+), 4 deletions(-)
~~~

## File Stats

~~~text
 docker-compose.yml        | 19 ++++++++++++-
 internal/kitchen/go.mod   |  8 ++++++
 internal/kitchen/go.sum   | 24 ++++++++++++++---
 internal/kitchen/prom.go  | 68 +++++++++++++++++++++++++++++++++++++++++++++++
 prometheus/prometheus.yml | 13 +++++++++
 5 files changed, 128 insertions(+), 4 deletions(-)
~~~

## Commit Comparison

~~~text
> 631aaeb prometheus & grafana
~~~

## Changed Files

~~~text
docker-compose.yml
internal/kitchen/go.mod
internal/kitchen/go.sum
internal/kitchen/prom.go
prometheus/prometheus.yml
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
docker-compose.yml
internal/kitchen/prom.go
prometheus/prometheus.yml
~~~

## Full Diff

~~~diff
diff --git a/docker-compose.yml b/docker-compose.yml
index ac2b22f..07e344f 100644
--- a/docker-compose.yml
+++ b/docker-compose.yml
@@ -42,4 +42,21 @@ services:
       ME_CONFIG_MONGODB_ADMINUSERNAME: root
       ME_CONFIG_MONGODB_ADMINPASSWORD: password
       ME_CONFIG_MONGODB_URL: mongodb://root:password@order-mongo:27017/
-      ME_CONFIG_BASICAUTH: false
\ No newline at end of file
+      ME_CONFIG_BASICAUTH: false
+
+  prometheus:
+    image: prom/prometheus
+    volumes:
+      - ./prometheus/:/etc/prometheus/
+    command:
+      - '--config.file=/etc/prometheus/prometheus.yml'
+      - '--storage.tsdb.path=/prometheus'
+      - '--web.console.libraries=/usr/share/prometheus/console_libraries'
+      - '--web.console.templates=/usr/share/prometheus/consoles'
+    ports:
+      - 9090:9090
+
+  grafana:
+    image: grafana/grafana
+    ports:
+      - 3000:3000
\ No newline at end of file
diff --git a/internal/kitchen/go.mod b/internal/kitchen/go.mod
index d68df76..2a29c39 100644
--- a/internal/kitchen/go.mod
+++ b/internal/kitchen/go.mod
@@ -6,6 +6,7 @@ replace github.com/ghost-yu/go_shop_second/common => ../common
 
 require (
 	github.com/ghost-yu/go_shop_second/common v0.0.0-00010101000000-000000000000
+	github.com/prometheus/client_golang v1.20.5
 	github.com/rabbitmq/amqp091-go v1.10.0
 	github.com/sirupsen/logrus v1.9.3
 	github.com/spf13/viper v1.19.0
@@ -14,6 +15,8 @@ require (
 
 require (
 	github.com/armon/go-metrics v0.4.1 // indirect
+	github.com/beorn7/perks v1.0.1 // indirect
+	github.com/cespare/xxhash/v2 v2.3.0 // indirect
 	github.com/fatih/color v1.14.1 // indirect
 	github.com/fsnotify/fsnotify v1.7.0 // indirect
 	github.com/go-logr/logr v1.4.2 // indirect
@@ -29,12 +32,17 @@ require (
 	github.com/hashicorp/golang-lru v0.5.4 // indirect
 	github.com/hashicorp/hcl v1.0.0 // indirect
 	github.com/hashicorp/serf v0.10.1 // indirect
+	github.com/klauspost/compress v1.17.9 // indirect
 	github.com/magiconair/properties v1.8.7 // indirect
 	github.com/mattn/go-colorable v0.1.13 // indirect
 	github.com/mattn/go-isatty v0.0.20 // indirect
 	github.com/mitchellh/go-homedir v1.1.0 // indirect
 	github.com/mitchellh/mapstructure v1.5.0 // indirect
+	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
 	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
+	github.com/prometheus/client_model v0.6.1 // indirect
+	github.com/prometheus/common v0.55.0 // indirect
+	github.com/prometheus/procfs v0.15.1 // indirect
 	github.com/sagikazarmark/locafero v0.4.0 // indirect
 	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
 	github.com/sourcegraph/conc v0.3.0 // indirect
diff --git a/internal/kitchen/go.sum b/internal/kitchen/go.sum
index b98b2e2..1ab1071 100644
--- a/internal/kitchen/go.sum
+++ b/internal/kitchen/go.sum
@@ -11,9 +11,12 @@ github.com/armon/go-radix v0.0.0-20180808171621-7fddfc383310/go.mod h1:ufUuZ+zHj
 github.com/armon/go-radix v1.0.0/go.mod h1:ufUuZ+zHj4x4TnLV4JWEpy2hxWSpsRywHrMgIH9cCH8=
 github.com/beorn7/perks v0.0.0-20180321164747-3a771d992973/go.mod h1:Dwedo/Wpr24TaqPxmxbtue+5NUziq4I4S80YR8gNf3Q=
 github.com/beorn7/perks v1.0.0/go.mod h1:KWe93zE9D1o94FZ5RNwFwVgaQK1VOXiVxmqh+CedLV8=
+github.com/beorn7/perks v1.0.1 h1:VlbKKnNfV8bJzeqoa4cOKqO6bYr3WgKZxO8Z16+hsOM=
 github.com/beorn7/perks v1.0.1/go.mod h1:G2ZrVWU2WbWT9wwq4/hrbKbnv/1ERSJQ0ibhJ6rlkpw=
 github.com/bgentry/speakeasy v0.1.0/go.mod h1:+zsyZBPWlz7T6j88CTgSN5bM796AkVf0kBD4zp0CCIs=
 github.com/cespare/xxhash/v2 v2.1.1/go.mod h1:VGX0DQ3Q6kWi7AoAeZDth3/j3BFtOZR5XLFGgcrjCOs=
+github.com/cespare/xxhash/v2 v2.3.0 h1:UL815xU9SqsFlibzuggzjXhog7bL6oX9BbNZnL2UFvs=
+github.com/cespare/xxhash/v2 v2.3.0/go.mod h1:VGX0DQ3Q6kWi7AoAeZDth3/j3BFtOZR5XLFGgcrjCOs=
 github.com/circonus-labs/circonus-gometrics v2.3.1+incompatible/go.mod h1:nmEj6Dob7S7YxXgwXpfOuvO54S+tGdZdw9fuRZt25Ag=
 github.com/circonus-labs/circonusllhist v0.1.3/go.mod h1:kMXHVDlOchFAehlya5ePtbp5jckzBHf4XRpQvBOLI+I=
 github.com/davecgh/go-spew v1.1.0/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
@@ -102,6 +105,8 @@ github.com/hashicorp/serf v0.10.1/go.mod h1:yL2t6BqATOLGc5HF7qbFkTfXoPIY0WZdWHfE
 github.com/json-iterator/go v1.1.6/go.mod h1:+SdeFBvtyEkXs7REEP0seUULqWtbJapLOCVDaaPEHmU=
 github.com/json-iterator/go v1.1.9/go.mod h1:KdQUCv79m/52Kvf8AW2vK1V8akMuk1QjK/uOdHXbAo4=
 github.com/julienschmidt/httprouter v1.2.0/go.mod h1:SYymIcj16QtmaHHD7aYtjjsJG7VTCxuUUipMqKk8s4w=
+github.com/klauspost/compress v1.17.9 h1:6KIumPrER1LHsvBVuDa0r5xaG0Es51mhhB9BQB2qeMA=
+github.com/klauspost/compress v1.17.9/go.mod h1:Di0epgTjJY877eYKx5yC51cX2A2Vl2ibi7bDH9ttBbw=
 github.com/konsorten/go-windows-terminal-sequences v1.0.1/go.mod h1:T0+1ngSBFLxvqU3pZ+m/2kptfBszLMUkC4ZK/EgS/cQ=
 github.com/kr/logfmt v0.0.0-20140226030751-b84e30acd515/go.mod h1:+0opPa2QZZtGFBFZlji/RkVcI2GknAs/DXo4wKdlNEc=
 github.com/kr/pretty v0.1.0/go.mod h1:dAy3ld7l9f0ibDNOQOHHMYYIIbhfbHSm3C4ZsoJORNo=
@@ -111,6 +116,8 @@ github.com/kr/pty v1.1.1/go.mod h1:pFQYn66WHrOpPYNljwOMqo10TkYh1fy3cYio2l3bCsQ=
 github.com/kr/text v0.1.0/go.mod h1:4Jbv+DJW3UT/LiOwJeYQe1efqtUx/iVham/4vfdArNI=
 github.com/kr/text v0.2.0 h1:5Nx0Ya0ZqY2ygV366QzturHI13Jq95ApcVaJBhpS+AY=
 github.com/kr/text v0.2.0/go.mod h1:eLer722TekiGuMkidMxC/pM04lWEeraHUUmBw8l2grE=
+github.com/kylelemons/godebug v1.1.0 h1:RPNrshWIDI6G2gRW9EHilWtl7Z6Sb1BR0xunSBf0SNc=
+github.com/kylelemons/godebug v1.1.0/go.mod h1:9/0rRGxNHcop5bhtWyNeEfOS8JIWk580+fNqagV/RAw=
 github.com/magiconair/properties v1.8.7 h1:IeQXZAiQcpL9mgcAe1Nu6cX9LLw6ExEHKjN0VQdvPDY=
 github.com/magiconair/properties v1.8.7/go.mod h1:Dhd985XPs7jluiymwWYZ0G4Z61jb3vdS329zhj2hYo0=
 github.com/mattn/go-colorable v0.0.9/go.mod h1:9vuHe8Xs5qXnSaW/c/ABM9alt+Vo+STaOChaDxuIBZU=
@@ -142,6 +149,8 @@ github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421/go.mod h1:6dJ
 github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd/go.mod h1:6dJC0mAP4ikYIbvyc7fijjWJddQyLn8Ig3JB5CqoB9Q=
 github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742/go.mod h1:bx2lNnkwVCuqBIxFjflWJWanXIb3RllmbCylyMrvgv0=
 github.com/modern-go/reflect2 v1.0.1/go.mod h1:bx2lNnkwVCuqBIxFjflWJWanXIb3RllmbCylyMrvgv0=
+github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 h1:C3w9PqII01/Oq1c1nUAm88MOHcQC9l5mIlSMApZMrHA=
+github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822/go.mod h1:+n7T8mK8HuQTcFwEeznm/DIxMOiR9yIdICNftLE1DvQ=
 github.com/mwitkow/go-conntrack v0.0.0-20161129095857-cc309e4a2223/go.mod h1:qRWi+5nqEBWmkhHvq77mSJWrCKwh8bxhgT7d/eI7P4U=
 github.com/pascaldekloe/goe v0.0.0-20180627143212-57f6aae5913c/go.mod h1:lzWF7FIEvWOWxwDKqyGYQf6ZUaNfKdP144TG7ZOy1lc=
 github.com/pascaldekloe/goe v0.1.0 h1:cBOtyMzM9HTpWjXfbbunk26uA6nG3a8n06Wieeh0MwY=
@@ -160,18 +169,26 @@ github.com/posener/complete v1.2.3/go.mod h1:WZIdtGGp+qx0sLrYKtIRAruyNpv6hFCicSg
 github.com/prometheus/client_golang v0.9.1/go.mod h1:7SWBe2y4D6OKWSNQJUaRYU/AaXPKyh/dDVn+NZz0KFw=
 github.com/prometheus/client_golang v1.0.0/go.mod h1:db9x61etRT2tGnBNRi70OPL5FsnadC4Ky3P0J6CfImo=
 github.com/prometheus/client_golang v1.4.0/go.mod h1:e9GMxYsXl05ICDXkRhurwBS4Q3OK1iX/F2sw+iXX5zU=
+github.com/prometheus/client_golang v1.20.5 h1:cxppBPuYhUnsO6yo/aoRol4L7q7UFfdm+bR9r+8l63Y=
+github.com/prometheus/client_golang v1.20.5/go.mod h1:PIEt8X02hGcP8JWbeHyeZ53Y/jReSnHgO035n//V5WE=
 github.com/prometheus/client_model v0.0.0-20180712105110-5c3871d89910/go.mod h1:MbSGuTsp3dbXC40dX6PRTWyKYBIrTGTE9sqQNg2J8bo=
 github.com/prometheus/client_model v0.0.0-20190129233127-fd36f4220a90/go.mod h1:xMI15A0UPsDsEKsMN9yxemIoYk6Tm2C1GtYGdfGttqA=
 github.com/prometheus/client_model v0.2.0/go.mod h1:xMI15A0UPsDsEKsMN9yxemIoYk6Tm2C1GtYGdfGttqA=
+github.com/prometheus/client_model v0.6.1 h1:ZKSh/rekM+n3CeS952MLRAdFwIKqeY8b62p8ais2e9E=
+github.com/prometheus/client_model v0.6.1/go.mod h1:OrxVMOVHjw3lKMa8+x6HeMGkHMQyHDk9E3jmP2AmGiY=
 github.com/prometheus/common v0.4.1/go.mod h1:TNfzLD0ON7rHzMJeJkieUDPYmFC7Snx/y86RQel1bk4=
 github.com/prometheus/common v0.9.1/go.mod h1:yhUN8i9wzaXS3w1O07YhxHEBxD+W35wd8bs7vj7HSQ4=
+github.com/prometheus/common v0.55.0 h1:KEi6DK7lXW/m7Ig5i47x0vRzuBsHuvJdi5ee6Y3G1dc=
+github.com/prometheus/common v0.55.0/go.mod h1:2SECS4xJG1kd8XF9IcM1gMX6510RAEL65zxzNImwdc8=
 github.com/prometheus/procfs v0.0.0-20181005140218-185b4288413d/go.mod h1:c3At6R/oaqEKCNdg8wHV1ftS6bRYblBhIjjI8uT2IGk=
 github.com/prometheus/procfs v0.0.2/go.mod h1:TjEm7ze935MbeOT/UhFTIMYKhuLP4wbCsTZCD3I8kEA=
 github.com/prometheus/procfs v0.0.8/go.mod h1:7Qr8sr6344vo1JqZ6HhLceV9o3AJ1Ff+GxbHq6oeK9A=
+github.com/prometheus/procfs v0.15.1 h1:YagwOFzUgYfKKHX6Dr+sHT7km/hxC76UB0learggepc=
+github.com/prometheus/procfs v0.15.1/go.mod h1:fB45yRUv8NstnjriLhBQLuOUt+WW4BsoGhij/e3PBqk=
 github.com/rabbitmq/amqp091-go v1.10.0 h1:STpn5XsHlHGcecLmMFCtg7mqq0RnD+zFr4uzukfVhBw=
 github.com/rabbitmq/amqp091-go v1.10.0/go.mod h1:Hy4jKW5kQART1u+JkDTF9YYOQUHXqMuhrgxOEeS7G4o=
-github.com/rogpeppe/go-internal v1.9.0 h1:73kH8U+JUqXU8lRuOHeVHaa/SZPifC7BkcraZVejAe8=
-github.com/rogpeppe/go-internal v1.9.0/go.mod h1:WtVeX8xhTBvf0smdhujwtBcq4Qrzq/fJaraNFVN+nFs=
+github.com/rogpeppe/go-internal v1.10.0 h1:TMyTOH3F/DB16zRVcYyreMH6GnZZrwQVAoYjRBZyWFQ=
+github.com/rogpeppe/go-internal v1.10.0/go.mod h1:UQnix2H7Ngw/k4C5ijL5+65zddjncjaFoBhdsK/akog=
 github.com/ryanuber/columnize v0.0.0-20160712163229-9b3edd62028f/go.mod h1:sm1tb6uqfes/u+d4ooFouqFdy9/2g9QGwK3SQygK0Ts=
 github.com/sagikazarmark/locafero v0.4.0 h1:HApY1R9zGo4DBgr7dqsTH/JJxLTTsOt7u6keLGt6kNQ=
 github.com/sagikazarmark/locafero v0.4.0/go.mod h1:Pe1W6UlPYUk/+wc/6KFhbORCfqzgYEpgQ3O5fPuL3H4=
@@ -289,8 +306,9 @@ google.golang.org/protobuf v1.35.1 h1:m3LfL6/Ca+fqnjnlqQXNpFPABW1UD7mjh8KO2mKFyt
 google.golang.org/protobuf v1.35.1/go.mod h1:9fA7Ob0pmnwhb644+1+CVWFRbNajQ6iRojtC/QF5bRE=
 gopkg.in/alecthomas/kingpin.v2 v2.2.6/go.mod h1:FMv+mEhP44yOT+4EoQTLFTRgOQ1FBLkstjWtayDeSgw=
 gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
-gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 h1:YR8cESwS4TdDjEe65xsg0ogRM/Nc3DYOhEAlW+xobZo=
 gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
+gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c h1:Hei/4ADfdWqJk1ZMxUNpqntNwaWcugrBjAiHlqqRiVk=
+gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c/go.mod h1:JHkPIbrfpd72SG/EVd6muEfDQjcINNoR0C8j2r3qZ4Q=
 gopkg.in/ini.v1 v1.67.0 h1:Dgnx+6+nfE+IfzjUEISNeydPJh9AXNNsWbGP9KzCsOA=
 gopkg.in/ini.v1 v1.67.0/go.mod h1:pNLf8WUiyNEtQjuu5G5vTm06TEv9tsIgeAvK8hOrP4k=
 gopkg.in/yaml.v2 v2.2.1/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
diff --git a/internal/kitchen/prom.go b/internal/kitchen/prom.go
new file mode 100644
index 0000000..bf07668
--- /dev/null
+++ b/internal/kitchen/prom.go
@@ -0,0 +1,68 @@
+package main
+
+import (
+	"bytes"
+	"encoding/json"
+	"log"
+	"math/rand"
+	"net/http"
+	"time"
+
+	"github.com/prometheus/client_golang/prometheus"
+	"github.com/prometheus/client_golang/prometheus/collectors"
+	"github.com/prometheus/client_golang/prometheus/promhttp"
+)
+
+const (
+	testAddr = "localhost:9123"
+)
+
+var httpStatusCodeCounter = prometheus.NewCounterVec(
+	prometheus.CounterOpts{
+		Name: "http_status_code_counter",
+		Help: "Count http status code",
+	},
+	[]string{"status_code"},
+)
+
+func main() {
+	go produceData()
+	reg := prometheus.NewRegistry()
+	prometheus.WrapRegistererWith(prometheus.Labels{"serviceName": "demo-service"}, reg).MustRegister(
+		collectors.NewGoCollector(),
+		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
+		httpStatusCodeCounter,
+	)
+	// localhost:9123/metrics
+	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
+	http.HandleFunc("/", sendMetricsHandler)
+	log.Fatal(http.ListenAndServe(testAddr, nil))
+}
+
+func sendMetricsHandler(w http.ResponseWriter, r *http.Request) {
+	var req request
+	defer func() {
+		httpStatusCodeCounter.WithLabelValues(req.StatusCode).Inc()
+		log.Printf("add 1 to %s", req.StatusCode)
+	}()
+	_ = json.NewDecoder(r.Body).Decode(&req)
+	log.Printf("receive req:%+v", req)
+	_, _ = w.Write([]byte(req.StatusCode))
+}
+
+type request struct {
+	StatusCode string
+}
+
+func produceData() {
+	codes := []string{"503", "404", "400", "200", "304", "500"}
+	for {
+		body, _ := json.Marshal(request{
+			StatusCode: codes[rand.Intn(len(codes))],
+		})
+		requestBody := bytes.NewBuffer(body)
+		http.Post("http://"+testAddr, "application/json", requestBody)
+		log.Printf("send request=%s to %s", requestBody.String(), testAddr)
+		time.Sleep(2 * time.Second)
+	}
+}
diff --git a/prometheus/prometheus.yml b/prometheus/prometheus.yml
new file mode 100644
index 0000000..a5f1952
--- /dev/null
+++ b/prometheus/prometheus.yml
@@ -0,0 +1,13 @@
+global:
+  scrape_interval:     15s # By default, scrape targets every 15 seconds.
+  evaluation_interval: 15s
+
+# A scrape configuration containing exactly one endpoint to scrape:
+# Here it's Prometheus itself.
+scrape_configs:
+  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
+  - job_name: 'demo-metrics'
+    metrics_path: /metrics
+    scrape_interval: 5s
+    static_configs:
+      - targets: ['host.docker.internal:9123']
\ No newline at end of file
~~~
