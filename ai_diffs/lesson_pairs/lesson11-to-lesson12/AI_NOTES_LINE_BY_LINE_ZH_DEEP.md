# 中文深度逐行解读（自动生成）

- 仓库: go_shop_second
- 起始引用: lesson11
- 结束引用: lesson12
- 生成时间: 2026-04-06 18:30:36 +08:00
- 解析范围: 仅 *.go，自动排除 *.pb.go 与 *.gen.go
- 解析模式: 分支差异模式（非祖先，按 diff）

## 分支差异说明

### 文件: internal/common/client/grpc.go

~~~diff
diff --git a/internal/common/client/grpc.go b/internal/common/client/grpc.go
new file mode 100644
index 0000000..f9f514a
--- /dev/null
+++ b/internal/common/client/grpc.go
@@ -0,0 +1,29 @@
+package client
+
+import (
+	"context"
+
+	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
+	"github.com/spf13/viper"
+	"google.golang.org/grpc"
+	"google.golang.org/grpc/credentials/insecure"
+)
+
+func NewStockGRPCClient(ctx context.Context) (client stockpb.StockServiceClient, close func() error, err error) {
+	grpcAddr := viper.GetString("stock.grpc-addr")
+	opts, err := grpcDialOpts(grpcAddr)
+	if err != nil {
+		return nil, func() error { return nil }, err
+	}
+	conn, err := grpc.NewClient(grpcAddr, opts...)
+	if err != nil {
+		return nil, func() error { return nil }, err
+	}
+	return stockpb.NewStockServiceClient(conn), conn.Close, nil
+}
+
+func grpcDialOpts(addr string) ([]grpc.DialOption, error) {
+	return []grpc.DialOption{
+		grpc.WithTransportCredentials(insecure.NewCredentials()),
+	}, nil
+}
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 新增 | $safeCode | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 语法块结束：关闭 import 或参数列表。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 新增 | $safeCode | 短变量声明：就地定义并初始化，收窄作用域。 |
| 新增 | $safeCode | 短变量声明：就地定义并初始化，收窄作用域。 |
| 新增 | $safeCode | 条件分支：进行校验、错误拦截或流程分叉。 |
| 新增 | $safeCode | 返回语句：输出当前结果并结束执行路径。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 短变量声明：就地定义并初始化，收窄作用域。 |
| 新增 | $safeCode | 条件分支：进行校验、错误拦截或流程分叉。 |
| 新增 | $safeCode | 返回语句：输出当前结果并结束执行路径。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 返回语句：输出当前结果并结束执行路径。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 新增 | $safeCode | 返回语句：输出当前结果并结束执行路径。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/common/config/viper.go

~~~diff
diff --git a/internal/common/config/viper.go b/internal/common/config/viper.go
index 5c2fe30..103d697 100644
--- a/internal/common/config/viper.go
+++ b/internal/common/config/viper.go
@@ -5,3 +4,0 @@ import "github.com/spf13/viper"
-// NewViperConfig 统一加载全局配置文件到内存。
-// 这个函数在 init() 里被调用，确保 main 启动前配置已经就位，
-// 后续各个函数可以直接用 viper.Get... 获取配置，而不用每次都读文件。
@@ -9,2 +5,0 @@ func NewViperConfig() error {
-	// SetConfigName 指定配置文件名（不包括扩展名）。
-	// 这里是 "global"，对应 global.yaml。
@@ -12,3 +6,0 @@ func NewViperConfig() error {
-
-	// SetConfigType 告诉 viper 配置格式是什么。
-	// viper 支持 json/yaml/toml/hcl 等多种格式。
@@ -16,4 +7,0 @@ func NewViperConfig() error {
-
-	// AddConfigPath 告诉 viper 去哪个目录找配置文件。
-	// "../common/config" 是相对当前执行目录的路径，
-	// 假设你从项目根目录启动 order 服务，viper 会去 internal/common/config 目录找 global.yaml。
@@ -21,4 +8,0 @@ func NewViperConfig() error {
-
-	// AutomaticEnv 让环境变量自动覆盖配置文件中的值。
-	// 举例：如果环境里设置了 ORDER_GRPC_ADDR=":6666"，会覆盖 yaml 里的 order.grpc-addr。
-	// 这样本地开发、测试、生产环境可以动态调整配置，而不改代码。
@@ -26,3 +9,0 @@ func NewViperConfig() error {
-
-	// ReadInConfig 真正去磁盘读配置文件，把内容解析到 viper 的内存结构里。
-	// 如果文件不存在或格式有问题，这里会返回 error。
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |

### 文件: internal/common/decorator/command.go

~~~diff
diff --git a/internal/common/decorator/command.go b/internal/common/decorator/command.go
index e820637..53f8c37 100644
--- a/internal/common/decorator/command.go
+++ b/internal/common/decorator/command.go
@@ -9 +8,0 @@ import (
-// CommandHandler 和 QueryHandler 对称，只是语义上代表“会改状态”的操作。
@@ -14,2 +12,0 @@ type CommandHandler[C, R any] interface {
-// ApplyCommandDecorators 让所有 command 统一拥有日志和指标能力。
-// 这样新增一个写操作时，只要实现业务 handler，不用重复写埋点代码。
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |

### 文件: internal/common/decorator/logging.go

~~~diff
diff --git a/internal/common/decorator/logging.go b/internal/common/decorator/logging.go
index 2b158ae..a9a2325 100644
--- a/internal/common/decorator/logging.go
+++ b/internal/common/decorator/logging.go
@@ -11,2 +10,0 @@ import (
-// queryLoggingDecorator 在真实 handler 外层补一圈日志。
-// “Decorator” 不是框架魔法，本质就是包一层结构体再转调 base。
@@ -19 +16,0 @@ func (q queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result
-	// generateActionName 会把类型名拿出来，作为统一的日志 action 名。
@@ -35 +31,0 @@ func (q queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result
-// generateActionName 直接复用请求类型名，避免每个 handler 手动写一次 action 常量。
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |

### 文件: internal/common/decorator/metrics.go

~~~diff
diff --git a/internal/common/decorator/metrics.go b/internal/common/decorator/metrics.go
index ba04887..fba9a77 100644
--- a/internal/common/decorator/metrics.go
+++ b/internal/common/decorator/metrics.go
@@ -10 +9,0 @@ import (
-// MetricsClient 是一个很小的抽象，方便先接假实现，后续再切到真实监控系统。
@@ -15,2 +13,0 @@ type MetricsClient interface {
-// queryMetricsDecorator 统计耗时、成功、失败次数。
-// 它和日志装饰器可以自由组合，因为两者都只依赖同一个 handler 接口。
@@ -23 +19,0 @@ func (q queryMetricsDecorator[C, R]) Handle(ctx context.Context, cmd C) (result
-	// 用 defer 统一收尾，确保无论成功还是失败都能上报指标。
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |

### 文件: internal/common/decorator/query.go

~~~diff
diff --git a/internal/common/decorator/query.go b/internal/common/decorator/query.go
index b5e9cf0..d3848fe 100644
--- a/internal/common/decorator/query.go
+++ b/internal/common/decorator/query.go
@@ -11 +10,0 @@ import (
-// 泛型接口的好处是：同一套日志/指标装饰逻辑可以复用在不同查询上。
@@ -16,2 +14,0 @@ type QueryHandler[Q, R any] interface {
-// ApplyQueryDecorators 按固定顺序把日志和指标能力包在真实 handler 外面。
-// 这就是装饰器模式：不改业务代码，也能统一加横切能力。
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |

### 文件: internal/common/discovery/consul/consul.go

~~~diff
diff --git a/internal/common/discovery/consul/consul.go b/internal/common/discovery/consul/consul.go
new file mode 100644
index 0000000..b0ef77f
--- /dev/null
+++ b/internal/common/discovery/consul/consul.go
@@ -0,0 +1,88 @@
+package consul
+
+import (
+	"context"
+	"errors"
+	"fmt"
+	"strconv"
+	"strings"
+	"sync"
+
+	"github.com/hashicorp/consul/api"
+	"github.com/sirupsen/logrus"
+)
+
+type Registry struct {
+	client *api.Client
+}
+
+var (
+	consulClient *Registry
+	once         sync.Once
+	initErr      error
+)
+
+func New(consulAddr string) (*Registry, error) {
+	once.Do(func() {
+		config := api.DefaultConfig()
+		config.Address = consulAddr
+		client, err := api.NewClient(config)
+		if err != nil {
+			initErr = err
+			return
+		}
+		consulClient = &Registry{
+			client: client,
+		}
+	})
+	if initErr != nil {
+		return nil, initErr
+	}
+	return consulClient, nil
+}
+
+func (r *Registry) Register(_ context.Context, instanceID, serviceName, hostPort string) error {
+	parts := strings.Split(hostPort, ":")
+	if len(parts) != 2 {
+		return errors.New("invalid host:port format")
+	}
+	host := parts[0]
+	port, _ := strconv.Atoi(parts[1])
+	return r.client.Agent().ServiceRegister(&api.AgentServiceRegistration{
+		ID:      instanceID,
+		Address: host,
+		Port:    port,
+		Name:    serviceName,
+		Check: &api.AgentServiceCheck{
+			CheckID:                        instanceID,
+			TLSSkipVerify:                  false,
+			TTL:                            "5s",
+			Timeout:                        "5s",
+			DeregisterCriticalServiceAfter: "10s",
+		},
+	})
+}
+
+func (r *Registry) Deregister(_ context.Context, instanceID, serviceName string) error {
+	logrus.WithFields(logrus.Fields{
+		"instanceID":  instanceID,
+		"serviceName": serviceName,
+	}).Info("deregister from consul")
+	return r.client.Agent().CheckDeregister(instanceID)
+}
+
+func (r *Registry) Discover(ctx context.Context, serviceName string) ([]string, error) {
+	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)
+	if err != nil {
+		return nil, err
+	}
+	var ips []string
+	for _, e := range entries {
+		ips = append(ips, fmt.Sprintf("%s:%d", e.Service.Address, e.Service.Port))
+	}
+	return ips, nil
+}
+
+func (r *Registry) HealthCheck(instanceID, serviceName string) error {
+	return r.client.Agent().UpdateTTL(instanceID, "online", api.HealthPassing)
+}
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 新增 | $safeCode | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 语法块结束：关闭 import 或参数列表。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 结构体定义：声明数据载体，承载状态或依赖。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 语法块结束：关闭 import 或参数列表。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 短变量声明：就地定义并初始化，收窄作用域。 |
| 新增 | $safeCode | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 新增 | $safeCode | 短变量声明：就地定义并初始化，收窄作用域。 |
| 新增 | $safeCode | 条件分支：进行校验、错误拦截或流程分叉。 |
| 新增 | $safeCode | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 新增 | $safeCode | 返回语句：输出当前结果并结束执行路径。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 条件分支：进行校验、错误拦截或流程分叉。 |
| 新增 | $safeCode | 返回语句：输出当前结果并结束执行路径。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 返回语句：输出当前结果并结束执行路径。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 新增 | $safeCode | 短变量声明：就地定义并初始化，收窄作用域。 |
| 新增 | $safeCode | 条件分支：进行校验、错误拦截或流程分叉。 |
| 新增 | $safeCode | 返回语句：输出当前结果并结束执行路径。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 短变量声明：就地定义并初始化，收窄作用域。 |
| 新增 | $safeCode | 短变量声明：就地定义并初始化，收窄作用域。 |
| 新增 | $safeCode | 返回语句：输出当前结果并结束执行路径。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 返回语句：输出当前结果并结束执行路径。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 新增 | $safeCode | 短变量声明：就地定义并初始化，收窄作用域。 |
| 新增 | $safeCode | 条件分支：进行校验、错误拦截或流程分叉。 |
| 新增 | $safeCode | 返回语句：输出当前结果并结束执行路径。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 循环结构：遍历数据或重复执行，需关注终止条件与副作用。 |
| 新增 | $safeCode | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 返回语句：输出当前结果并结束执行路径。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 新增 | $safeCode | 返回语句：输出当前结果并结束执行路径。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/common/discovery/discovery.go

~~~diff
diff --git a/internal/common/discovery/discovery.go b/internal/common/discovery/discovery.go
new file mode 100644
index 0000000..7cf6f58
--- /dev/null
+++ b/internal/common/discovery/discovery.go
@@ -0,0 +1,20 @@
+package discovery
+
+import (
+	"context"
+	"fmt"
+	"math/rand"
+	"time"
+)
+
+type Registry interface {
+	Register(ctx context.Context, instanceID, serviceName, hostPort string) error
+	Deregister(ctx context.Context, instanceID, serviceName string) error
+	Discover(ctx context.Context, serviceName string) ([]string, error)
+	HealthCheck(instanceID, serviceName string) error
+}
+
+func GenerateInstanceID(serviceName string) string {
+	x := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
+	return fmt.Sprintf("%s-%d", serviceName, x)
+}
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 新增 | $safeCode | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 语法块结束：关闭 import 或参数列表。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 接口定义：声明能力契约，用于解耦与可替换实现。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 新增 | $safeCode | 短变量声明：就地定义并初始化，收窄作用域。 |
| 新增 | $safeCode | 返回语句：输出当前结果并结束执行路径。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/common/discovery/grpc.go

~~~diff
diff --git a/internal/common/discovery/grpc.go b/internal/common/discovery/grpc.go
new file mode 100644
index 0000000..cda522a
--- /dev/null
+++ b/internal/common/discovery/grpc.go
@@ -0,0 +1,37 @@
+package discovery
+
+import (
+	"context"
+	"time"
+
+	"github.com/ghost-yu/go_shop_second/common/discovery/consul"
+	"github.com/sirupsen/logrus"
+	"github.com/spf13/viper"
+)
+
+func RegisterToConsul(ctx context.Context, serviceName string) (func() error, error) {
+	registry, err := consul.New(viper.GetString("consul.addr"))
+	if err != nil {
+		return func() error { return nil }, err
+	}
+	instanceID := GenerateInstanceID(serviceName)
+	grpcAddr := viper.Sub(serviceName).GetString("grpc-addr")
+	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
+		return func() error { return nil }, err
+	}
+	go func() {
+		for {
+			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
+				logrus.Panicf("no heartbeat from %s to registry, err=%v", serviceName, err)
+			}
+			time.Sleep(1 * time.Second)
+		}
+	}()
+	logrus.WithFields(logrus.Fields{
+		"serviceName": serviceName,
+		"addr":        grpcAddr,
+	}).Info("registered to consul")
+	return func() error {
+		return registry.Deregister(ctx, instanceID, serviceName)
+	}, nil
+}
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 新增 | $safeCode | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 语法块结束：关闭 import 或参数列表。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 新增 | $safeCode | 短变量声明：就地定义并初始化，收窄作用域。 |
| 新增 | $safeCode | 条件分支：进行校验、错误拦截或流程分叉。 |
| 新增 | $safeCode | 返回语句：输出当前结果并结束执行路径。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 短变量声明：就地定义并初始化，收窄作用域。 |
| 新增 | $safeCode | 短变量声明：就地定义并初始化，收窄作用域。 |
| 新增 | $safeCode | 条件分支：进行校验、错误拦截或流程分叉。 |
| 新增 | $safeCode | 返回语句：输出当前结果并结束执行路径。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | goroutine 启动：引入并发执行，需关注生命周期管理。 |
| 新增 | $safeCode | 循环结构：遍历数据或重复执行，需关注终止条件与副作用。 |
| 新增 | $safeCode | 条件分支：进行校验、错误拦截或流程分叉。 |
| 新增 | $safeCode | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 返回语句：输出当前结果并结束执行路径。 |
| 新增 | $safeCode | 返回语句：输出当前结果并结束执行路径。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/common/metrics/todo_metrics.go

~~~diff
diff --git a/internal/common/metrics/todo_metrics.go b/internal/common/metrics/todo_metrics.go
index b379f64..a05f6b2 100644
--- a/internal/common/metrics/todo_metrics.go
+++ b/internal/common/metrics/todo_metrics.go
@@ -3,2 +2,0 @@ package metrics
-// TodoMetrics 是占位实现，用来先打通接口而不真正上报指标。
-// 初学时先理解“依赖抽象”比马上接 Prometheus 更重要。
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |

### 文件: internal/common/server/gprc.go

~~~diff
diff --git a/internal/common/server/gprc.go b/internal/common/server/gprc.go
index ba9f6e1..817cfcb 100644
--- a/internal/common/server/gprc.go
+++ b/internal/common/server/gprc.go
@@ -13,2 +12,0 @@ import (
-// init 会在包加载时执行，这里统一替换 gRPC 默认日志实现，
-// 让框架日志和业务日志都走 logrus，便于初学者排查问题时只看一种格式。
@@ -21,2 +18,0 @@ func init() {
-// RunGRPCServer 负责按服务名读取配置，再把真正的启动动作委托给 RunGRPCServerOnAddr。
-// 这样测试时可以直接传地址，避免每次都依赖 viper 配置。
@@ -32,2 +27,0 @@ func RunGRPCServer(serviceName string, registerServer func(server *grpc.Server))
-// RunGRPCServerOnAddr 创建 gRPC server，并通过 registerServer 回调注册具体业务服务。
-// 这类“框架负责启动，业务负责注册”的回调模式，在 Go 后端里非常常见。
@@ -38 +31,0 @@ func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server))
-			// 一元拦截器链类似 HTTP 中间件，用来放日志、指标、鉴权这类横切逻辑。
@@ -48 +40,0 @@ func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server))
-			// 流式 RPC 和一元 RPC 分开配置，是因为两类调用的拦截接口不同。
@@ -58 +49,0 @@ func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server))
-	// 业务方在这里把 OrderService/StockService 的实现挂到 grpcServer 上。
@@ -61 +51,0 @@ func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server))
-	// net.Listen 只负责占住端口；真正开始接收 gRPC 请求要等 Serve 调用后才发生。
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |

### 文件: internal/common/server/http.go

~~~diff
diff --git a/internal/common/server/http.go b/internal/common/server/http.go
index 1c561ce..7f39359 100644
--- a/internal/common/server/http.go
+++ b/internal/common/server/http.go
@@ -8 +7,0 @@ import (
-// RunHTTPServer 按服务名从配置中读取监听地址，再交给 RunHTTPServerOnAddr 真正启动。
@@ -17 +15,0 @@ func RunHTTPServer(serviceName string, wrapper func(router *gin.Engine)) {
-// RunHTTPServerOnAddr 创建 gin.Engine，并把路由注册动作交给 wrapper 回调。
@@ -19 +16,0 @@ func RunHTTPServerOnAddr(addr string, wrapper func(router *gin.Engine)) {
-	// gin.New 不带默认中间件，适合教学阶段观察“哪些能力是手动加进去的”。
@@ -21 +17,0 @@ func RunHTTPServerOnAddr(addr string, wrapper func(router *gin.Engine)) {
-	// wrapper 负责把 OpenAPI 生成的 handler 绑定到具体路由。
@@ -23 +18,0 @@ func RunHTTPServerOnAddr(addr string, wrapper func(router *gin.Engine)) {
-	// 这一行不会修改已有路由，只是创建了一个没有被接住的 RouterGroup。
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |

### 文件: internal/order/adapters/grpc/stock_grpc.go

~~~diff
diff --git a/internal/order/adapters/grpc/stock_grpc.go b/internal/order/adapters/grpc/stock_grpc.go
index f3acb2b..c2397d6 100644
--- a/internal/order/adapters/grpc/stock_grpc.go
+++ b/internal/order/adapters/grpc/stock_grpc.go
@@ -11,2 +10,0 @@ import (
-// StockGRPC 是 order 侧访问 stock 服务的远程适配器。
-// 它实现的是应用层定义的 StockService 接口，而不是把 proto client 直接暴露出去。
@@ -22 +19,0 @@ func (s StockGRPC) CheckIfItemsInStock(ctx context.Context, items []*orderpb.Ite
-	// 这里负责把应用层参数翻译成 gRPC 请求对象。
@@ -29 +25,0 @@ func (s StockGRPC) GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.I
-	// 对调用方来说这里只是“查商品”，至于底层走 gRPC 还是别的协议都被屏蔽了。
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |

### 文件: internal/order/adapters/order_inmem_repository.go

~~~diff
diff --git a/internal/order/adapters/order_inmem_repository.go b/internal/order/adapters/order_inmem_repository.go
index 94364ec..818b51f 100644
--- a/internal/order/adapters/order_inmem_repository.go
+++ b/internal/order/adapters/order_inmem_repository.go
@@ -13,2 +12,0 @@ import (
-// MemoryOrderRepository 是 domain.Repository 的内存版实现。
-// 教学项目先用它跑通流程，后面换数据库时只需要替换这一层。
@@ -21 +18,0 @@ func NewMemoryOrderRepository() *MemoryOrderRepository {
-	// 这里放一条假数据，方便一开始就能演示“查询已有订单”的路径。
@@ -37 +33,0 @@ func (m *MemoryOrderRepository) Create(_ context.Context, order *domain.Order) (
-	// 写操作要加互斥锁，避免多个请求并发修改切片时产生数据竞争。
@@ -41 +36,0 @@ func (m *MemoryOrderRepository) Create(_ context.Context, order *domain.Order) (
-		// 当前用 Unix 时间戳凑一个简单 ID，真实项目里通常会换成雪花算法或 UUID。
@@ -60 +54,0 @@ func (m *MemoryOrderRepository) Get(_ context.Context, id, customerID string) (*
-	// 读操作使用读锁，允许多个查询并发进行。
@@ -73 +66,0 @@ func (m *MemoryOrderRepository) Update(ctx context.Context, order *domain.Order,
-	// UpdateFn 把“怎么改”交给上层，把“在哪里存”留给仓储层，是一种职责分离。
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |

### 文件: internal/order/app/app.go

~~~diff
diff --git a/internal/order/app/app.go b/internal/order/app/app.go
index 5963e7c..b2e04ce 100644
--- a/internal/order/app/app.go
+++ b/internal/order/app/app.go
@@ -8 +7,0 @@ import (
-// Application 是应用层门面，让上层只依赖一个总入口，而不用感知每个 handler 的构造细节。
@@ -14 +12,0 @@ type Application struct {
-// Commands 聚合“会改状态”的用例，典型如创建、更新订单。
@@ -20 +17,0 @@ type Commands struct {
-// Queries 聚合“只读”的用例，便于和 Commands 做职责分离。
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |

### 文件: internal/order/app/command/create_order.go

~~~diff
diff --git a/internal/order/app/command/create_order.go b/internal/order/app/command/create_order.go
index 859e170..1c72b04 100644
--- a/internal/order/app/command/create_order.go
+++ b/internal/order/app/command/create_order.go
@@ -14 +13,0 @@ import (
-// CreateOrder 是创建订单用例的输入模型。
@@ -20 +18,0 @@ type CreateOrder struct {
-// CreateOrderResult 是返回给端口层的结果，避免直接暴露底层仓储对象。
@@ -27 +24,0 @@ type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult
-// createOrderHandler 只持有完成下单流程所需的两个依赖：订单仓储和库存服务。
@@ -42 +38,0 @@ func NewCreateOrderHandler(
-	// 下单也是 command，所以同样套上统一的日志和指标装饰器。
@@ -51 +46,0 @@ func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*Creat
-	// 先校验库存，再真正创建订单，避免生成一张无法履约的脏订单。
@@ -70 +64,0 @@ func (c createOrderHandler) validate(ctx context.Context, items []*orderpb.ItemW
-	// packItems 先合并重复商品，避免把同一个商品重复发给库存服务校验。
@@ -80 +73,0 @@ func packItems(items []*orderpb.ItemWithQuantity) []*orderpb.ItemWithQuantity {
-	// merged 以商品 ID 为 key，把数量累加起来。
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |

### 文件: internal/order/app/command/update_order.go

~~~diff
diff --git a/internal/order/app/command/update_order.go b/internal/order/app/command/update_order.go
index 3cce1e7..f40716d 100644
--- a/internal/order/app/command/update_order.go
+++ b/internal/order/app/command/update_order.go
@@ -11 +10,0 @@ import (
-// UpdateOrder 把“更新哪个订单”和“如何更新”一起传给 handler。
@@ -19 +17,0 @@ type UpdateOrderHandler decorator.CommandHandler[UpdateOrder, interface{}]
-// updateOrderHandler 负责流程控制，具体修改细节由 UpdateFn 注入。
@@ -41 +38,0 @@ func (c updateOrderHandler) Handle(ctx context.Context, cmd UpdateOrder) (interf
-	// 给 nil UpdateFn 一个 no-op 默认值，避免直接调用时发生空指针问题。
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |

### 文件: internal/order/app/query/get_customer_order.go

~~~diff
diff --git a/internal/order/app/query/get_customer_order.go b/internal/order/app/query/get_customer_order.go
index 19b7ffc..b107e01 100644
--- a/internal/order/app/query/get_customer_order.go
+++ b/internal/order/app/query/get_customer_order.go
@@ -11,2 +10,0 @@ import (
-// GetCustomerOrder 是查询订单的输入对象。
-// 把参数收成一个结构体后，后续加字段不会破坏函数签名。
@@ -18 +15,0 @@ type GetCustomerOrder struct {
-// GetCustomerOrderHandler 是一个已经套好泛型的查询处理器类型别名。
@@ -21 +17,0 @@ type GetCustomerOrderHandler decorator.QueryHandler[GetCustomerOrder, *domain.Or
-// getCustomerOrderHandler 才是真正的业务实现，字段里只放它需要的依赖。
@@ -34 +29,0 @@ func NewGetCustomerOrderHandler(
-	// 构造函数里统一加装饰器，这样调用方拿到的就是“带日志和指标能力”的 handler。
@@ -43 +37,0 @@ func (g getCustomerOrderHandler) Handle(ctx context.Context, query GetCustomerOr
-	// 查询用例本身很薄，只负责描述流程，把数据访问交给仓储接口。
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |

### 文件: internal/order/app/query/service.go

~~~diff
diff --git a/internal/order/app/query/service.go b/internal/order/app/query/service.go
index 0d33749..2e3e4f4 100644
--- a/internal/order/app/query/service.go
+++ b/internal/order/app/query/service.go
@@ -10,2 +9,0 @@ import (
-// StockService 是 order 应用层眼里的“库存能力端口”。
-// 它隔离了底层 gRPC 细节，让 create_order 用例只表达业务意图。
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |

### 文件: internal/order/domain/order/order.go

~~~diff
diff --git a/internal/order/domain/order/order.go b/internal/order/domain/order/order.go
index 073b8f0..6a9e94b 100644
--- a/internal/order/domain/order/order.go
+++ b/internal/order/domain/order/order.go
@@ -5 +4,0 @@ import "github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
-// Order 是订单领域对象，代表业务里真正被创建、查询、更新的核心实体。
@@ -7,7 +6,3 @@ type Order struct {
-	// ID 是订单主键，由仓储层创建时生成。
-	ID string
-	// CustomerID 表示订单属于哪个客户。
-	CustomerID string
-	// Status 预留给支付中、已支付等订单状态流转。
-	Status string
-	// PaymentLink 预留给支付服务返回的支付链接。
+	ID          string
+	CustomerID  string
+	Status      string
@@ -15,2 +10 @@ type Order struct {
-	// Items 直接复用 proto 里的商品结构，减少服务边界两侧的对象转换成本。
-	Items []*orderpb.Item
+	Items       []*orderpb.Item
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |

### 文件: internal/order/domain/order/repository.go

~~~diff
diff --git a/internal/order/domain/order/repository.go b/internal/order/domain/order/repository.go
index 685eddc..04b783f 100644
--- a/internal/order/domain/order/repository.go
+++ b/internal/order/domain/order/repository.go
@@ -8 +7,0 @@ import (
-// Repository 定义订单持久化能力，应用层只依赖这个接口，而不关心底层是内存还是数据库。
@@ -19 +17,0 @@ type Repository interface {
-// NotFoundError 是带业务语义的错误类型，调用方可以明确知道“订单不存在”而不是普通异常。
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |

### 文件: internal/order/http.go

~~~diff
diff --git a/internal/order/http.go b/internal/order/http.go
index 8b5d70b..b40adc7 100644
--- a/internal/order/http.go
+++ b/internal/order/http.go
@@ -13,2 +12,0 @@ import (
-// HTTPServer 是 OpenAPI 生成接口在 order 服务里的具体实现。
-// 它本身不处理复杂业务，只负责把 HTTP 请求翻译成应用层调用。
@@ -21 +18,0 @@ func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID stri
-	// ShouldBindJSON 负责把请求体反序列化为 proto 请求对象。
@@ -26 +22,0 @@ func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID stri
-	// 进入应用层前，把 HTTP 层对象转换成 command 对象，避免业务逻辑依赖 gin。
@@ -43 +38,0 @@ func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerI
-	// 查询场景走 Queries 分组，体现读写分离的组织方式。
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |

### 文件: internal/order/main.go

~~~diff
diff --git a/internal/order/main.go b/internal/order/main.go
index 9807588..cb08cc4 100644
--- a/internal/order/main.go
+++ b/internal/order/main.go
@@ -6,0 +7 @@ import (
+	"github.com/ghost-yu/go_shop_second/common/discovery"
@@ -17 +17,0 @@ import (
-// init 先加载全局配置，这样 main 里启动服务时就能直接从 viper 读取地址和服务名。
@@ -25 +24,0 @@ func main() {
-	// serviceName 对应 global.yaml 里的 order.service-name，后续 server 包会继续用它找端口配置。
@@ -28 +26,0 @@ func main() {
-	// ctx 用来把进程级生命周期传给下游依赖，例如 gRPC 客户端连接。
@@ -32 +29,0 @@ func main() {
-	// NewApplication 是组合根：在这里把 repo、远程 client、handler 全部装配好。
@@ -36 +33,8 @@ func main() {
-	// gRPC 和 HTTP 同时启动，说明 order 服务对内和对外用了两套协议暴露能力。
+	deregisterFunc, err := discovery.RegisterToConsul(ctx, serviceName)
+	if err != nil {
+		logrus.Fatal(err)
+	}
+	defer func() {
+		_ = deregisterFunc()
+	}()
+
@@ -42 +45,0 @@ func main() {
-	// HTTP 侧通过 OpenAPI 生成的 RegisterHandlersWithOptions 完成路由到实现的映射。
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 新增 | $safeCode | 短变量声明：就地定义并初始化，收窄作用域。 |
| 新增 | $safeCode | 条件分支：进行校验、错误拦截或流程分叉。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | defer：函数退出前执行，常用于资源释放与收尾。 |
| 新增 | $safeCode | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |

### 文件: internal/order/ports/grpc.go

~~~diff
diff --git a/internal/order/ports/grpc.go b/internal/order/ports/grpc.go
index 4841586..e1e8621 100644
--- a/internal/order/ports/grpc.go
+++ b/internal/order/ports/grpc.go
@@ -11,2 +10,0 @@ import (
-// GRPCServer 是 orderpb.OrderServiceServer 的端口适配器实现。
-// 这一层的职责和 HTTPServer 一样，都是把外部协议请求转给应用层。
@@ -22 +19,0 @@ func (G GRPCServer) CreateOrder(ctx context.Context, request *orderpb.CreateOrde
-	// lesson 当前故意保留 TODO，方便后续逐步实现 gRPC 版本的下单流程。
@@ -28 +24,0 @@ func (G GRPCServer) GetOrder(ctx context.Context, request *orderpb.GetOrderReque
-	// 这里最终会把 request 转为 query，再返回 proto 层定义的 orderpb.Order。
@@ -34 +29,0 @@ func (G GRPCServer) UpdateOrder(ctx context.Context, order *orderpb.Order) (*emp
-	// 更新接口通常会调用 Commands.UpdateOrder；当前阶段先保留占位。
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |

### 文件: internal/order/service/application.go

~~~diff
diff --git a/internal/order/service/application.go b/internal/order/service/application.go
index faf1f24..648e438 100644
--- a/internal/order/service/application.go
+++ b/internal/order/service/application.go
@@ -16,2 +15,0 @@ import (
-// NewApplication 是 order 服务的组合根。
-// 它负责创建外部依赖、构造适配器，并返回一个已经装配好的应用层门面。
@@ -23 +20,0 @@ func NewApplication(ctx context.Context) (app.Application, func()) {
-	// 这里把生成的 gRPC client 再包一层适配器，避免应用层直接依赖 proto 细节。
@@ -31 +27,0 @@ func newApplication(_ context.Context, stockGRPC query.StockService) app.Applica
-	// 组合根里统一决定“接口对应哪种实现”，后面要切数据库或 mock 时只改这里。
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |

### 文件: internal/stock/adapters/stock_inmem_repository.go

~~~diff
diff --git a/internal/stock/adapters/stock_inmem_repository.go b/internal/stock/adapters/stock_inmem_repository.go
index c38b115..f4aed23 100644
--- a/internal/stock/adapters/stock_inmem_repository.go
+++ b/internal/stock/adapters/stock_inmem_repository.go
@@ -11,2 +10,0 @@ import (
-// MemoryStockRepository 用 map 模拟库存数据源。
-// 相比切片，map 更适合按商品 ID 做快速查找。
@@ -18 +15,0 @@ type MemoryStockRepository struct {
-// stub 是教学阶段的假数据，用来模拟一个已经存在的库存商品。
@@ -25,0 +23,18 @@ var stub = map[string]*orderpb.Item{
+	"item1": {
+		ID:       "item1",
+		Name:     "stub item 1",
+		Quantity: 10000,
+		PriceID:  "stub_item1_price_id",
+	},
+	"item2": {
+		ID:       "item2",
+		Name:     "stub item 2",
+		Quantity: 10000,
+		PriceID:  "stub_item2_price_id",
+	},
+	"item3": {
+		ID:       "item3",
+		Name:     "stub item 3",
+		Quantity: 10000,
+		PriceID:  "stub_item3_price_id",
+	},
@@ -36 +50,0 @@ func (m MemoryStockRepository) GetItems(ctx context.Context, ids []string) ([]*o
-	// 这里只有读场景，所以拿读锁即可。
@@ -53 +66,0 @@ func (m MemoryStockRepository) GetItems(ctx context.Context, ids []string) ([]*o
-	// 返回“部分结果 + 缺失错误”可以帮助上层更明确地决定如何提示用户。
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |

### 文件: internal/stock/app/app.go

~~~diff
diff --git a/internal/stock/app/app.go b/internal/stock/app/app.go
index c159f54..94a18fb 100644
--- a/internal/stock/app/app.go
+++ b/internal/stock/app/app.go
@@ -3,2 +3,2 @@ package app
-// Application 预留给 stock 服务聚合 Commands 和 Queries。
-// 当前 lesson 里 stock 逻辑还比较薄，所以两个分组都是空壳结构。
+import "github.com/ghost-yu/go_shop_second/stock/app/query"
+
@@ -10 +9,0 @@ type Application struct {
-// Commands 未来承载库存写操作，例如扣减库存。
@@ -13,2 +12,4 @@ type Commands struct{}
-// Queries 未来承载库存读操作，例如按商品 ID 查询库存。
-type Queries struct{}
+type Queries struct {
+	CheckIfItemsInStock query.CheckIfItemsInStockHandler
+	GetItems            query.GetItemsHandler
+}
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 新增 | $safeCode | 单行导入：引入一个依赖包，常见于依赖较少的文件。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 结构体定义：声明数据载体，承载状态或依赖。 |
| 新增 | $safeCode | 结构体定义：声明数据载体，承载状态或依赖。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/stock/app/query/check_if_items_in_stock.go

~~~diff
diff --git a/internal/stock/app/query/check_if_items_in_stock.go b/internal/stock/app/query/check_if_items_in_stock.go
new file mode 100644
index 0000000..d1078f0
--- /dev/null
+++ b/internal/stock/app/query/check_if_items_in_stock.go
@@ -0,0 +1,46 @@
+package query
+
+import (
+	"context"
+
+	"github.com/ghost-yu/go_shop_second/common/decorator"
+	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+	domain "github.com/ghost-yu/go_shop_second/stock/domain/stock"
+	"github.com/sirupsen/logrus"
+)
+
+type CheckIfItemsInStock struct {
+	Items []*orderpb.ItemWithQuantity
+}
+
+type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, []*orderpb.Item]
+
+type checkIfItemsInStockHandler struct {
+	stockRepo domain.Repository
+}
+
+func NewCheckIfItemsInStockHandler(
+	stockRepo domain.Repository,
+	logger *logrus.Entry,
+	metricClient decorator.MetricsClient,
+) CheckIfItemsInStockHandler {
+	if stockRepo == nil {
+		panic("nil stockRepo")
+	}
+	return decorator.ApplyQueryDecorators[CheckIfItemsInStock, []*orderpb.Item](
+		checkIfItemsInStockHandler{stockRepo: stockRepo},
+		logger,
+		metricClient,
+	)
+}
+
+func (h checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*orderpb.Item, error) {
+	var res []*orderpb.Item
+	for _, i := range query.Items {
+		res = append(res, &orderpb.Item{
+			ID:       i.ID,
+			Quantity: i.Quantity,
+		})
+	}
+	return res, nil
+}
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 新增 | $safeCode | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 语法块结束：关闭 import 或参数列表。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 结构体定义：声明数据载体，承载状态或依赖。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 类型定义：建立语义模型，影响方法与边界设计。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 结构体定义：声明数据载体，承载状态或依赖。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 条件分支：进行校验、错误拦截或流程分叉。 |
| 新增 | $safeCode | panic：不可恢复错误路径，常用于关键初始化失败。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 返回语句：输出当前结果并结束执行路径。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 语法块结束：关闭 import 或参数列表。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 循环结构：遍历数据或重复执行，需关注终止条件与副作用。 |
| 新增 | $safeCode | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 返回语句：输出当前结果并结束执行路径。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/stock/app/query/get_items.go

~~~diff
diff --git a/internal/stock/app/query/get_items.go b/internal/stock/app/query/get_items.go
new file mode 100644
index 0000000..063a6b0
--- /dev/null
+++ b/internal/stock/app/query/get_items.go
@@ -0,0 +1,43 @@
+package query
+
+import (
+	"context"
+
+	"github.com/ghost-yu/go_shop_second/common/decorator"
+	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+	domain "github.com/ghost-yu/go_shop_second/stock/domain/stock"
+	"github.com/sirupsen/logrus"
+)
+
+type GetItems struct {
+	ItemIDs []string
+}
+
+type GetItemsHandler decorator.QueryHandler[GetItems, []*orderpb.Item]
+
+type getItemsHandler struct {
+	stockRepo domain.Repository
+}
+
+func NewGetItemsHandler(
+	stockRepo domain.Repository,
+	logger *logrus.Entry,
+	metricClient decorator.MetricsClient,
+) GetItemsHandler {
+	if stockRepo == nil {
+		panic("nil stockRepo")
+	}
+	return decorator.ApplyQueryDecorators[GetItems, []*orderpb.Item](
+		getItemsHandler{stockRepo: stockRepo},
+		logger,
+		metricClient,
+	)
+}
+
+func (g getItemsHandler) Handle(ctx context.Context, query GetItems) ([]*orderpb.Item, error) {
+	items, err := g.stockRepo.GetItems(ctx, query.ItemIDs)
+	if err != nil {
+		return nil, err
+	}
+	return items, nil
+}
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 新增 | $safeCode | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 语法块结束：关闭 import 或参数列表。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 结构体定义：声明数据载体，承载状态或依赖。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 类型定义：建立语义模型，影响方法与边界设计。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 结构体定义：声明数据载体，承载状态或依赖。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 条件分支：进行校验、错误拦截或流程分叉。 |
| 新增 | $safeCode | panic：不可恢复错误路径，常用于关键初始化失败。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 返回语句：输出当前结果并结束执行路径。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 语法块结束：关闭 import 或参数列表。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 新增 | $safeCode | 短变量声明：就地定义并初始化，收窄作用域。 |
| 新增 | $safeCode | 条件分支：进行校验、错误拦截或流程分叉。 |
| 新增 | $safeCode | 返回语句：输出当前结果并结束执行路径。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 返回语句：输出当前结果并结束执行路径。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/stock/domain/stock/repository.go

~~~diff
diff --git a/internal/stock/domain/stock/repository.go b/internal/stock/domain/stock/repository.go
index b0f1796..7a58e6a 100644
--- a/internal/stock/domain/stock/repository.go
+++ b/internal/stock/domain/stock/repository.go
@@ -11,2 +10,0 @@ import (
-// Repository 定义库存服务需要提供的最小数据访问能力。
-// order 服务只要通过应用层依赖这个接口，就不需要知道库存数据放在哪里。
@@ -17 +14,0 @@ type Repository interface {
-// NotFoundError 用来明确告诉调用方哪些商品在库存里不存在。
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |

### 文件: internal/stock/main.go

~~~diff
diff --git a/internal/stock/main.go b/internal/stock/main.go
index e1fe459..8990bf4 100644
--- a/internal/stock/main.go
+++ b/internal/stock/main.go
@@ -6,0 +7 @@ import (
+	"github.com/ghost-yu/go_shop_second/common/discovery"
@@ -16 +16,0 @@ import (
-// stock 服务和 order 服务一样，在进程入口先加载配置，避免后续各层自己读文件。
@@ -24 +23,0 @@ func main() {
-	// serviceName 决定去哪个配置分组下找 grpc-addr 等参数。
@@ -26 +24,0 @@ func main() {
-	// serverType 让同一个二进制可以选择跑哪种协议，当前 lesson 主要走 grpc 分支。
@@ -31 +28,0 @@ func main() {
-	// 这里先组装应用层，再根据配置决定挂到哪个端口适配器上。
@@ -35,0 +33,9 @@ func main() {
+
+	deregisterFunc, err := discovery.RegisterToConsul(ctx, serviceName)
+	if err != nil {
+		logrus.Fatal(err)
+	}
+	defer func() {
+		_ = deregisterFunc()
+	}()
+
@@ -38 +43,0 @@ func main() {
-		// gRPC 分支把 ports.GRPCServer 注册为 StockService 的服务端实现。
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 新增 | $safeCode | 短变量声明：就地定义并初始化，收窄作用域。 |
| 新增 | $safeCode | 条件分支：进行校验、错误拦截或流程分叉。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | defer：函数退出前执行，常用于资源释放与收尾。 |
| 新增 | $safeCode | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |

### 文件: internal/stock/ports/grpc.go

~~~diff
diff --git a/internal/stock/ports/grpc.go b/internal/stock/ports/grpc.go
index 72bc356..91e34e8 100644
--- a/internal/stock/ports/grpc.go
+++ b/internal/stock/ports/grpc.go
@@ -7,0 +8 @@ import (
+	"github.com/ghost-yu/go_shop_second/stock/app/query"
@@ -10,2 +10,0 @@ import (
-// GRPCServer 负责实现 stockpb.StockServiceServer 接口。
-// 外部请求进入后，应该在这里完成参数转换，再转给应用层处理。
@@ -21,3 +20,5 @@ func (G GRPCServer) GetItems(ctx context.Context, request *stockpb.GetItemsReque
-	// 这里后续会把商品 ID 列表转给库存查询用例。
-	//TODO implement me
-	panic("implement me")
+	items, err := G.app.Queries.GetItems.Handle(ctx, query.GetItems{ItemIDs: request.ItemIDs})
+	if err != nil {
+		return nil, err
+	}
+	return &stockpb.GetItemsResponse{Items: items}, nil
@@ -27,3 +28,8 @@ func (G GRPCServer) CheckIfItemsInStock(ctx context.Context, request *stockpb.Ch
-	// 这一步对应 order 下单前的库存校验，是跨服务调用的关键入口。
-	//TODO implement me
-	panic("implement me")
+	items, err := G.app.Queries.CheckIfItemsInStock.Handle(ctx, query.CheckIfItemsInStock{Items: request.Items})
+	if err != nil {
+		return nil, err
+	}
+	return &stockpb.CheckIfItemsInStockResponse{
+		InStock: 1,
+		Items:   items,
+	}, nil
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | panic：不可恢复错误路径，常用于关键初始化失败。 |
| 新增 | $safeCode | 短变量声明：就地定义并初始化，收窄作用域。 |
| 新增 | $safeCode | 条件分支：进行校验、错误拦截或流程分叉。 |
| 新增 | $safeCode | 返回语句：输出当前结果并结束执行路径。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 返回语句：输出当前结果并结束执行路径。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | panic：不可恢复错误路径，常用于关键初始化失败。 |
| 新增 | $safeCode | 短变量声明：就地定义并初始化，收窄作用域。 |
| 新增 | $safeCode | 条件分支：进行校验、错误拦截或流程分叉。 |
| 新增 | $safeCode | 返回语句：输出当前结果并结束执行路径。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 返回语句：输出当前结果并结束执行路径。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/stock/service/application.go

~~~diff
diff --git a/internal/stock/service/application.go b/internal/stock/service/application.go
index 23daf9a..0ccb6a8 100644
--- a/internal/stock/service/application.go
+++ b/internal/stock/service/application.go
@@ -5,0 +6,2 @@ import (
+	"github.com/ghost-yu/go_shop_second/common/metrics"
+	"github.com/ghost-yu/go_shop_second/stock/adapters"
@@ -6,0 +9,2 @@ import (
+	"github.com/ghost-yu/go_shop_second/stock/app/query"
+	"github.com/sirupsen/logrus"
@@ -9,4 +13,11 @@ import (
-// NewApplication 当前先返回一个空的应用层门面。
-// lesson 后续扩展库存能力时，repo 和 handler 的装配点也会放在这里。
-func NewApplication(ctx context.Context) app.Application {
-	return app.Application{}
+func NewApplication(_ context.Context) app.Application {
+	stockRepo := adapters.NewMemoryStockRepository()
+	logger := logrus.NewEntry(logrus.StandardLogger())
+	metricsClient := metrics.TodoMetrics{}
+	return app.Application{
+		Commands: app.Commands{},
+		Queries: app.Queries{
+			CheckIfItemsInStock: query.NewCheckIfItemsInStockHandler(stockRepo, logger, metricsClient),
+			GetItems:            query.NewGetItemsHandler(stockRepo, logger, metricsClient),
+		},
+	}
~~~

| 变更类型 | 代码行 | 中文深度解释 |
| --- | --- | --- |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 新增 | $safeCode | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 注释：解释意图、风险或待办，帮助理解设计。 |
| 删除 | $safeCode | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 删除 | $safeCode | 返回语句：输出当前结果并结束执行路径。 |
| 新增 | $safeCode | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 新增 | $safeCode | 短变量声明：就地定义并初始化，收窄作用域。 |
| 新增 | $safeCode | 短变量声明：就地定义并初始化，收窄作用域。 |
| 新增 | $safeCode | 短变量声明：就地定义并初始化，收窄作用域。 |
| 新增 | $safeCode | 返回语句：输出当前结果并结束执行路径。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |
| 新增 | $safeCode | 代码块结束：收束当前函数、分支或类型定义。 |


