# 中文深度逐行解读（自动生成）

- 仓库: go_shop_second
- 起始引用: lesson24
- 结束引用: lesson25
- 生成时间: 2026-04-06 18:31:45 +08:00
- 解析范围: 仅 *.go，自动排除 *.pb.go 与 *.gen.go
- 解析模式: 祖先链模式（按提交）

## 提交 1: [fb48e4b] wait for

### 文件: internal/common/client/grpc.go

~~~go
   1: package client
   2: 
   3: import (
   4: 	"context"
   5: 	"errors"
   6: 	"net"
   7: 	"time"
   8: 
   9: 	"github.com/ghost-yu/go_shop_second/common/discovery"
  10: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
  11: 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
  12: 	"github.com/sirupsen/logrus"
  13: 	"github.com/spf13/viper"
  14: 	"google.golang.org/grpc"
  15: 	"google.golang.org/grpc/credentials/insecure"
  16: )
  17: 
  18: func NewStockGRPCClient(ctx context.Context) (client stockpb.StockServiceClient, close func() error, err error) {
  19: 	if !WaitForStockGRPCClient(viper.GetDuration("dial-grpc-timeout") * time.Second) {
  20: 		return nil, nil, errors.New("stock grpc not available")
  21: 	}
  22: 	grpcAddr, err := discovery.GetServiceAddr(ctx, viper.GetString("stock.service-name"))
  23: 	if err != nil {
  24: 		return nil, func() error { return nil }, err
  25: 	}
  26: 	if grpcAddr == "" {
  27: 		logrus.Warn("empty grpc addr for stock grpc")
  28: 	}
  29: 	opts := grpcDialOpts(grpcAddr)
  30: 	conn, err := grpc.NewClient(grpcAddr, opts...)
  31: 	if err != nil {
  32: 		return nil, func() error { return nil }, err
  33: 	}
  34: 	return stockpb.NewStockServiceClient(conn), conn.Close, nil
  35: }
  36: 
  37: func NewOrderGRPCClient(ctx context.Context) (client orderpb.OrderServiceClient, close func() error, err error) {
  38: 	if !WaitForOrderGRPCClient(viper.GetDuration("dial-grpc-timeout") * time.Second) {
  39: 		return nil, nil, errors.New("order grpc not available")
  40: 	}
  41: 	grpcAddr, err := discovery.GetServiceAddr(ctx, viper.GetString("order.service-name"))
  42: 	if err != nil {
  43: 		return nil, func() error { return nil }, err
  44: 	}
  45: 	if grpcAddr == "" {
  46: 		logrus.Warn("empty grpc addr for order grpc")
  47: 	}
  48: 	opts := grpcDialOpts(grpcAddr)
  49: 
  50: 	conn, err := grpc.NewClient(grpcAddr, opts...)
  51: 	if err != nil {
  52: 		return nil, func() error { return nil }, err
  53: 	}
  54: 	return orderpb.NewOrderServiceClient(conn), conn.Close, nil
  55: }
  56: 
  57: func grpcDialOpts(_ string) []grpc.DialOption {
  58: 	return []grpc.DialOption{
  59: 		grpc.WithTransportCredentials(insecure.NewCredentials()),
  60: 	}
  61: }
  62: 
  63: func WaitForOrderGRPCClient(timeout time.Duration) bool {
  64: 	logrus.Infof("waiting for order grpc client, timeout: %v seconds", timeout.Seconds())
  65: 	return waitFor(viper.GetString("order.grpc-addr"), timeout)
  66: }
  67: 
  68: func WaitForStockGRPCClient(timeout time.Duration) bool {
  69: 	logrus.Infof("waiting for stock grpc client, timeout: %v seconds", timeout.Seconds())
  70: 	return waitFor(viper.GetString("stock.grpc-addr"), timeout)
  71: }
  72: 
  73: func waitFor(addr string, timeout time.Duration) bool {
  74: 	portAvailable := make(chan struct{})
  75: 	timeoutCh := time.After(timeout)
  76: 
  77: 	go func() {
  78: 		for {
  79: 			select {
  80: 			case <-timeoutCh:
  81: 				return
  82: 			default:
  83: 				// continue
  84: 			}
  85: 			_, err := net.Dial("tcp", addr)
  86: 			if err == nil {
  87: 				close(portAvailable)
  88: 				return
  89: 			}
  90: 			time.Sleep(200 * time.Millisecond)
  91: 		}
  92: 	}()
  93: 
  94: 	select {
  95: 	case <-portAvailable:
  96: 		return true
  97: 	case <-timeoutCh:
  98: 		return false
  99: 	}
 100: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 8 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 9 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 10 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 11 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 12 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 13 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 14 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 15 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 16 | 语法块结束：关闭 import 或参数列表。 |
| 17 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 18 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 19 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 20 | 返回语句：输出当前结果并结束执行路径。 |
| 21 | 代码块结束：收束当前函数、分支或类型定义。 |
| 22 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 23 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 24 | 返回语句：输出当前结果并结束执行路径。 |
| 25 | 代码块结束：收束当前函数、分支或类型定义。 |
| 26 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 27 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 28 | 代码块结束：收束当前函数、分支或类型定义。 |
| 29 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 30 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 31 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 32 | 返回语句：输出当前结果并结束执行路径。 |
| 33 | 代码块结束：收束当前函数、分支或类型定义。 |
| 34 | 返回语句：输出当前结果并结束执行路径。 |
| 35 | 代码块结束：收束当前函数、分支或类型定义。 |
| 36 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 37 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 38 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 39 | 返回语句：输出当前结果并结束执行路径。 |
| 40 | 代码块结束：收束当前函数、分支或类型定义。 |
| 41 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 42 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 43 | 返回语句：输出当前结果并结束执行路径。 |
| 44 | 代码块结束：收束当前函数、分支或类型定义。 |
| 45 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 46 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 47 | 代码块结束：收束当前函数、分支或类型定义。 |
| 48 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 49 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 50 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 51 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 52 | 返回语句：输出当前结果并结束执行路径。 |
| 53 | 代码块结束：收束当前函数、分支或类型定义。 |
| 54 | 返回语句：输出当前结果并结束执行路径。 |
| 55 | 代码块结束：收束当前函数、分支或类型定义。 |
| 56 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 57 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 58 | 返回语句：输出当前结果并结束执行路径。 |
| 59 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 60 | 代码块结束：收束当前函数、分支或类型定义。 |
| 61 | 代码块结束：收束当前函数、分支或类型定义。 |
| 62 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 63 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 64 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 65 | 返回语句：输出当前结果并结束执行路径。 |
| 66 | 代码块结束：收束当前函数、分支或类型定义。 |
| 67 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 68 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 69 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 70 | 返回语句：输出当前结果并结束执行路径。 |
| 71 | 代码块结束：收束当前函数、分支或类型定义。 |
| 72 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 73 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 74 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 75 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 76 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 77 | goroutine 启动：引入并发执行，需关注生命周期管理。 |
| 78 | 循环结构：遍历数据或重复执行，需关注终止条件与副作用。 |
| 79 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 80 | 分支标签：定义 switch 的命中条件。 |
| 81 | 返回语句：输出当前结果并结束执行路径。 |
| 82 | 默认分支：兜底处理，防止未知输入静默失败。 |
| 83 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 84 | 代码块结束：收束当前函数、分支或类型定义。 |
| 85 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 86 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 87 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 88 | 返回语句：输出当前结果并结束执行路径。 |
| 89 | 代码块结束：收束当前函数、分支或类型定义。 |
| 90 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 91 | 代码块结束：收束当前函数、分支或类型定义。 |
| 92 | 代码块结束：收束当前函数、分支或类型定义。 |
| 93 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 94 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 95 | 分支标签：定义 switch 的命中条件。 |
| 96 | 返回语句：输出当前结果并结束执行路径。 |
| 97 | 分支标签：定义 switch 的命中条件。 |
| 98 | 返回语句：输出当前结果并结束执行路径。 |
| 99 | 代码块结束：收束当前函数、分支或类型定义。 |
| 100 | 代码块结束：收束当前函数、分支或类型定义。 |

## 提交 2: [f95add9] otel

### 文件: internal/common/client/grpc.go

~~~go
   1: package client
   2: 
   3: import (
   4: 	"context"
   5: 	"errors"
   6: 	"net"
   7: 	"time"
   8: 
   9: 	"github.com/ghost-yu/go_shop_second/common/discovery"
  10: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
  11: 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
  12: 	"github.com/sirupsen/logrus"
  13: 	"github.com/spf13/viper"
  14: 	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
  15: 	"google.golang.org/grpc"
  16: 	"google.golang.org/grpc/credentials/insecure"
  17: )
  18: 
  19: func NewStockGRPCClient(ctx context.Context) (client stockpb.StockServiceClient, close func() error, err error) {
  20: 	if !WaitForStockGRPCClient(viper.GetDuration("dial-grpc-timeout") * time.Second) {
  21: 		return nil, nil, errors.New("stock grpc not available")
  22: 	}
  23: 	grpcAddr, err := discovery.GetServiceAddr(ctx, viper.GetString("stock.service-name"))
  24: 	if err != nil {
  25: 		return nil, func() error { return nil }, err
  26: 	}
  27: 	if grpcAddr == "" {
  28: 		logrus.Warn("empty grpc addr for stock grpc")
  29: 	}
  30: 	opts := grpcDialOpts(grpcAddr)
  31: 	conn, err := grpc.NewClient(grpcAddr, opts...)
  32: 	if err != nil {
  33: 		return nil, func() error { return nil }, err
  34: 	}
  35: 	return stockpb.NewStockServiceClient(conn), conn.Close, nil
  36: }
  37: 
  38: func NewOrderGRPCClient(ctx context.Context) (client orderpb.OrderServiceClient, close func() error, err error) {
  39: 	if !WaitForOrderGRPCClient(viper.GetDuration("dial-grpc-timeout") * time.Second) {
  40: 		return nil, nil, errors.New("order grpc not available")
  41: 	}
  42: 	grpcAddr, err := discovery.GetServiceAddr(ctx, viper.GetString("order.service-name"))
  43: 	if err != nil {
  44: 		return nil, func() error { return nil }, err
  45: 	}
  46: 	if grpcAddr == "" {
  47: 		logrus.Warn("empty grpc addr for order grpc")
  48: 	}
  49: 	opts := grpcDialOpts(grpcAddr)
  50: 
  51: 	conn, err := grpc.NewClient(grpcAddr, opts...)
  52: 	if err != nil {
  53: 		return nil, func() error { return nil }, err
  54: 	}
  55: 	return orderpb.NewOrderServiceClient(conn), conn.Close, nil
  56: }
  57: 
  58: func grpcDialOpts(_ string) []grpc.DialOption {
  59: 	return []grpc.DialOption{
  60: 		grpc.WithTransportCredentials(insecure.NewCredentials()),
  61: 		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
  62: 	}
  63: }
  64: 
  65: func WaitForOrderGRPCClient(timeout time.Duration) bool {
  66: 	logrus.Infof("waiting for order grpc client, timeout: %v seconds", timeout.Seconds())
  67: 	return waitFor(viper.GetString("order.grpc-addr"), timeout)
  68: }
  69: 
  70: func WaitForStockGRPCClient(timeout time.Duration) bool {
  71: 	logrus.Infof("waiting for stock grpc client, timeout: %v seconds", timeout.Seconds())
  72: 	return waitFor(viper.GetString("stock.grpc-addr"), timeout)
  73: }
  74: 
  75: func waitFor(addr string, timeout time.Duration) bool {
  76: 	portAvailable := make(chan struct{})
  77: 	timeoutCh := time.After(timeout)
  78: 
  79: 	go func() {
  80: 		for {
  81: 			select {
  82: 			case <-timeoutCh:
  83: 				return
  84: 			default:
  85: 				// continue
  86: 			}
  87: 			_, err := net.Dial("tcp", addr)
  88: 			if err == nil {
  89: 				close(portAvailable)
  90: 				return
  91: 			}
  92: 			time.Sleep(200 * time.Millisecond)
  93: 		}
  94: 	}()
  95: 
  96: 	select {
  97: 	case <-portAvailable:
  98: 		return true
  99: 	case <-timeoutCh:
 100: 		return false
 101: 	}
 102: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 8 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 9 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 10 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 11 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 12 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 13 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 14 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 15 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 16 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 17 | 语法块结束：关闭 import 或参数列表。 |
| 18 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 19 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 20 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 21 | 返回语句：输出当前结果并结束执行路径。 |
| 22 | 代码块结束：收束当前函数、分支或类型定义。 |
| 23 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 24 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 25 | 返回语句：输出当前结果并结束执行路径。 |
| 26 | 代码块结束：收束当前函数、分支或类型定义。 |
| 27 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 28 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 29 | 代码块结束：收束当前函数、分支或类型定义。 |
| 30 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 31 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 32 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 33 | 返回语句：输出当前结果并结束执行路径。 |
| 34 | 代码块结束：收束当前函数、分支或类型定义。 |
| 35 | 返回语句：输出当前结果并结束执行路径。 |
| 36 | 代码块结束：收束当前函数、分支或类型定义。 |
| 37 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 38 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 39 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 40 | 返回语句：输出当前结果并结束执行路径。 |
| 41 | 代码块结束：收束当前函数、分支或类型定义。 |
| 42 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 43 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 44 | 返回语句：输出当前结果并结束执行路径。 |
| 45 | 代码块结束：收束当前函数、分支或类型定义。 |
| 46 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 47 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 48 | 代码块结束：收束当前函数、分支或类型定义。 |
| 49 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 50 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 51 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 52 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 53 | 返回语句：输出当前结果并结束执行路径。 |
| 54 | 代码块结束：收束当前函数、分支或类型定义。 |
| 55 | 返回语句：输出当前结果并结束执行路径。 |
| 56 | 代码块结束：收束当前函数、分支或类型定义。 |
| 57 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 58 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 59 | 返回语句：输出当前结果并结束执行路径。 |
| 60 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 61 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 62 | 代码块结束：收束当前函数、分支或类型定义。 |
| 63 | 代码块结束：收束当前函数、分支或类型定义。 |
| 64 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 65 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 66 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 67 | 返回语句：输出当前结果并结束执行路径。 |
| 68 | 代码块结束：收束当前函数、分支或类型定义。 |
| 69 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 70 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 71 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 72 | 返回语句：输出当前结果并结束执行路径。 |
| 73 | 代码块结束：收束当前函数、分支或类型定义。 |
| 74 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 75 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 76 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 77 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 78 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 79 | goroutine 启动：引入并发执行，需关注生命周期管理。 |
| 80 | 循环结构：遍历数据或重复执行，需关注终止条件与副作用。 |
| 81 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 82 | 分支标签：定义 switch 的命中条件。 |
| 83 | 返回语句：输出当前结果并结束执行路径。 |
| 84 | 默认分支：兜底处理，防止未知输入静默失败。 |
| 85 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 86 | 代码块结束：收束当前函数、分支或类型定义。 |
| 87 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 88 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 89 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 90 | 返回语句：输出当前结果并结束执行路径。 |
| 91 | 代码块结束：收束当前函数、分支或类型定义。 |
| 92 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 93 | 代码块结束：收束当前函数、分支或类型定义。 |
| 94 | 代码块结束：收束当前函数、分支或类型定义。 |
| 95 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 96 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 97 | 分支标签：定义 switch 的命中条件。 |
| 98 | 返回语句：输出当前结果并结束执行路径。 |
| 99 | 分支标签：定义 switch 的命中条件。 |
| 100 | 返回语句：输出当前结果并结束执行路径。 |
| 101 | 代码块结束：收束当前函数、分支或类型定义。 |
| 102 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/common/server/gprc.go

~~~go
   1: package server
   2: 
   3: import (
   4: 	"net"
   5: 
   6: 	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
   7: 	grpc_tags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
   8: 	"github.com/sirupsen/logrus"
   9: 	"github.com/spf13/viper"
  10: 	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
  11: 	"google.golang.org/grpc"
  12: )
  13: 
  14: func init() {
  15: 	logger := logrus.New()
  16: 	logger.SetLevel(logrus.WarnLevel)
  17: 	grpc_logrus.ReplaceGrpcLogger(logrus.NewEntry(logger))
  18: }
  19: 
  20: func RunGRPCServer(serviceName string, registerServer func(server *grpc.Server)) {
  21: 	addr := viper.Sub(serviceName).GetString("grpc-addr")
  22: 	if addr == "" {
  23: 		// TODO: Warning log
  24: 		addr = viper.GetString("fallback-grpc-addr")
  25: 	}
  26: 	RunGRPCServerOnAddr(addr, registerServer)
  27: }
  28: 
  29: func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server)) {
  30: 	logrusEntry := logrus.NewEntry(logrus.StandardLogger())
  31: 	grpcServer := grpc.NewServer(
  32: 		grpc.StatsHandler(otelgrpc.NewServerHandler()),
  33: 		grpc.ChainUnaryInterceptor(
  34: 			grpc_tags.UnaryServerInterceptor(grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor)),
  35: 			grpc_logrus.UnaryServerInterceptor(logrusEntry),
  36: 		),
  37: 		grpc.ChainStreamInterceptor(
  38: 			grpc_tags.StreamServerInterceptor(grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor)),
  39: 			grpc_logrus.StreamServerInterceptor(logrusEntry),
  40: 		),
  41: 	)
  42: 	registerServer(grpcServer)
  43: 
  44: 	listen, err := net.Listen("tcp", addr)
  45: 	if err != nil {
  46: 		logrus.Panic(err)
  47: 	}
  48: 	logrus.Infof("Starting gRPC server, Listening: %s", addr)
  49: 	if err := grpcServer.Serve(listen); err != nil {
  50: 		logrus.Panic(err)
  51: 	}
  52: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 6 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 7 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 8 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 9 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 10 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 11 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 12 | 语法块结束：关闭 import 或参数列表。 |
| 13 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 14 | init 函数：包加载时自动执行，常用于初始化与注册。 |
| 15 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 16 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 17 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 18 | 代码块结束：收束当前函数、分支或类型定义。 |
| 19 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 20 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 21 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 22 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 23 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 24 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 25 | 代码块结束：收束当前函数、分支或类型定义。 |
| 26 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 27 | 代码块结束：收束当前函数、分支或类型定义。 |
| 28 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 29 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 30 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 31 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 32 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 33 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 34 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 35 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 36 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 37 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 38 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 39 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 40 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 41 | 语法块结束：关闭 import 或参数列表。 |
| 42 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 43 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 44 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 45 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 46 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 47 | 代码块结束：收束当前函数、分支或类型定义。 |
| 48 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 49 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 50 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 51 | 代码块结束：收束当前函数、分支或类型定义。 |
| 52 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/common/server/http.go

~~~go
   1: package server
   2: 
   3: import (
   4: 	"github.com/gin-gonic/gin"
   5: 	"github.com/spf13/viper"
   6: 	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
   7: )
   8: 
   9: func RunHTTPServer(serviceName string, wrapper func(router *gin.Engine)) {
  10: 	addr := viper.Sub(serviceName).GetString("http-addr")
  11: 	if addr == "" {
  12: 		panic("empty http address")
  13: 	}
  14: 	RunHTTPServerOnAddr(addr, wrapper)
  15: }
  16: 
  17: func RunHTTPServerOnAddr(addr string, wrapper func(router *gin.Engine)) {
  18: 	apiRouter := gin.New()
  19: 	setMiddlewares(apiRouter)
  20: 	wrapper(apiRouter)
  21: 	apiRouter.Group("/api")
  22: 	if err := apiRouter.Run(addr); err != nil {
  23: 		panic(err)
  24: 	}
  25: }
  26: 
  27: func setMiddlewares(r *gin.Engine) {
  28: 	r.Use(gin.Recovery())
  29: 	r.Use(otelgin.Middleware("default_server"))
  30: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 语法块结束：关闭 import 或参数列表。 |
| 8 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 9 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 10 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 11 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 12 | panic：不可恢复错误路径，常用于关键初始化失败。 |
| 13 | 代码块结束：收束当前函数、分支或类型定义。 |
| 14 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 15 | 代码块结束：收束当前函数、分支或类型定义。 |
| 16 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 17 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 18 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 19 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 20 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 21 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 22 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 23 | panic：不可恢复错误路径，常用于关键初始化失败。 |
| 24 | 代码块结束：收束当前函数、分支或类型定义。 |
| 25 | 代码块结束：收束当前函数、分支或类型定义。 |
| 26 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 27 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 28 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 29 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 30 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/common/tracing/jaeger.go

~~~go
   1: package tracing
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"go.opentelemetry.io/contrib/propagators/b3"
   7: 	"go.opentelemetry.io/otel"
   8: 	"go.opentelemetry.io/otel/exporters/jaeger"
   9: 	"go.opentelemetry.io/otel/propagation"
  10: 	"go.opentelemetry.io/otel/sdk/resource"
  11: 	sdktrace "go.opentelemetry.io/otel/sdk/trace"
  12: 	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
  13: 	"go.opentelemetry.io/otel/trace"
  14: )
  15: 
  16: var tracer = otel.Tracer("default_tracer")
  17: 
  18: func InitJaegerProvider(jaegerURL, serviceName string) (func(ctx context.Context) error, error) {
  19: 	if jaegerURL == "" {
  20: 		panic("empty jaeger url")
  21: 	}
  22: 	tracer = otel.Tracer(serviceName)
  23: 	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerURL)))
  24: 	if err != nil {
  25: 		return nil, err
  26: 	}
  27: 	tp := sdktrace.NewTracerProvider(
  28: 		sdktrace.WithBatcher(exp),
  29: 		sdktrace.WithResource(resource.NewSchemaless(
  30: 			semconv.ServiceNameKey.String(serviceName),
  31: 		)),
  32: 	)
  33: 	otel.SetTracerProvider(tp)
  34: 	b3Propagator := b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader))
  35: 	p := propagation.NewCompositeTextMapPropagator(
  36: 		propagation.TraceContext{}, propagation.Baggage{}, b3Propagator,
  37: 	)
  38: 	otel.SetTextMapPropagator(p)
  39: 	return tp.Shutdown, nil
  40: }
  41: 
  42: func Start(ctx context.Context, name string) (context.Context, trace.Span) {
  43: 	return tracer.Start(ctx, name)
  44: }
  45: 
  46: func TraceID(ctx context.Context) string {
  47: 	spanCtx := trace.SpanContextFromContext(ctx)
  48: 	return spanCtx.TraceID().String()
  49: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 8 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 9 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 10 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 11 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 12 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 13 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 14 | 语法块结束：关闭 import 或参数列表。 |
| 15 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 16 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 17 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 18 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 19 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 20 | panic：不可恢复错误路径，常用于关键初始化失败。 |
| 21 | 代码块结束：收束当前函数、分支或类型定义。 |
| 22 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 23 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 24 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 25 | 返回语句：输出当前结果并结束执行路径。 |
| 26 | 代码块结束：收束当前函数、分支或类型定义。 |
| 27 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 28 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 29 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 30 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 31 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 32 | 语法块结束：关闭 import 或参数列表。 |
| 33 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 34 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 35 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 36 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 37 | 语法块结束：关闭 import 或参数列表。 |
| 38 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 39 | 返回语句：输出当前结果并结束执行路径。 |
| 40 | 代码块结束：收束当前函数、分支或类型定义。 |
| 41 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 42 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 43 | 返回语句：输出当前结果并结束执行路径。 |
| 44 | 代码块结束：收束当前函数、分支或类型定义。 |
| 45 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 46 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 47 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 48 | 返回语句：输出当前结果并结束执行路径。 |
| 49 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/order/http.go

~~~go
   1: package main
   2: 
   3: import (
   4: 	"fmt"
   5: 	"net/http"
   6: 
   7: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   8: 	"github.com/ghost-yu/go_shop_second/common/tracing"
   9: 	"github.com/ghost-yu/go_shop_second/order/app"
  10: 	"github.com/ghost-yu/go_shop_second/order/app/command"
  11: 	"github.com/ghost-yu/go_shop_second/order/app/query"
  12: 	"github.com/gin-gonic/gin"
  13: )
  14: 
  15: type HTTPServer struct {
  16: 	app app.Application
  17: }
  18: 
  19: func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
  20: 	ctx, span := tracing.Start(c, "PostCustomerCustomerIDOrders")
  21: 	defer span.End()
  22: 
  23: 	var req orderpb.CreateOrderRequest
  24: 	if err := c.ShouldBindJSON(&req); err != nil {
  25: 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
  26: 		return
  27: 	}
  28: 	r, err := H.app.Commands.CreateOrder.Handle(ctx, command.CreateOrder{
  29: 		CustomerID: req.CustomerID,
  30: 		Items:      req.Items,
  31: 	})
  32: 	if err != nil {
  33: 		c.JSON(http.StatusOK, gin.H{"error": err})
  34: 		return
  35: 	}
  36: 	c.JSON(http.StatusOK, gin.H{
  37: 		"message":      "success",
  38: 		"trace_id":     tracing.TraceID(ctx),
  39: 		"customer_id":  req.CustomerID,
  40: 		"order_id":     r.OrderID,
  41: 		"redirect_url": fmt.Sprintf("http://localhost:8282/success?customerID=%s&orderID=%s", req.CustomerID, r.OrderID),
  42: 	})
  43: }
  44: 
  45: func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
  46: 	ctx, span := tracing.Start(c, "GetCustomerCustomerIDOrdersOrderID")
  47: 	defer span.End()
  48: 
  49: 	o, err := H.app.Queries.GetCustomerOrder.Handle(ctx, query.GetCustomerOrder{
  50: 		OrderID:    orderID,
  51: 		CustomerID: customerID,
  52: 	})
  53: 	if err != nil {
  54: 		c.JSON(http.StatusOK, gin.H{"error": err})
  55: 		return
  56: 	}
  57: 	c.JSON(http.StatusOK, gin.H{
  58: 		"message":  "success",
  59: 		"trace_id": tracing.TraceID(ctx),
  60: 		"data": gin.H{
  61: 			"Order": o,
  62: 		},
  63: 	})
  64: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 6 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 7 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 8 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 9 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 10 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 11 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 12 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 13 | 语法块结束：关闭 import 或参数列表。 |
| 14 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 15 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 16 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 17 | 代码块结束：收束当前函数、分支或类型定义。 |
| 18 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 19 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 20 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 21 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 22 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 23 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 24 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 25 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 26 | 返回语句：输出当前结果并结束执行路径。 |
| 27 | 代码块结束：收束当前函数、分支或类型定义。 |
| 28 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 29 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 30 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 31 | 代码块结束：收束当前函数、分支或类型定义。 |
| 32 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 33 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 34 | 返回语句：输出当前结果并结束执行路径。 |
| 35 | 代码块结束：收束当前函数、分支或类型定义。 |
| 36 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 37 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 38 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 39 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 40 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 41 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 42 | 代码块结束：收束当前函数、分支或类型定义。 |
| 43 | 代码块结束：收束当前函数、分支或类型定义。 |
| 44 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 45 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 46 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 47 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 48 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 49 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 50 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 51 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 52 | 代码块结束：收束当前函数、分支或类型定义。 |
| 53 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 54 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 55 | 返回语句：输出当前结果并结束执行路径。 |
| 56 | 代码块结束：收束当前函数、分支或类型定义。 |
| 57 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 58 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 59 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 60 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 61 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 62 | 代码块结束：收束当前函数、分支或类型定义。 |
| 63 | 代码块结束：收束当前函数、分支或类型定义。 |
| 64 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/order/main.go

~~~go
   1: package main
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/broker"
   7: 	"github.com/ghost-yu/go_shop_second/common/config"
   8: 	"github.com/ghost-yu/go_shop_second/common/discovery"
   9: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
  10: 	"github.com/ghost-yu/go_shop_second/common/logging"
  11: 	"github.com/ghost-yu/go_shop_second/common/server"
  12: 	"github.com/ghost-yu/go_shop_second/common/tracing"
  13: 	"github.com/ghost-yu/go_shop_second/order/infrastructure/consumer"
  14: 	"github.com/ghost-yu/go_shop_second/order/ports"
  15: 	"github.com/ghost-yu/go_shop_second/order/service"
  16: 	"github.com/gin-gonic/gin"
  17: 	"github.com/sirupsen/logrus"
  18: 	"github.com/spf13/viper"
  19: 	"google.golang.org/grpc"
  20: )
  21: 
  22: func init() {
  23: 	logging.Init()
  24: 	if err := config.NewViperConfig(); err != nil {
  25: 		logrus.Fatal(err)
  26: 	}
  27: }
  28: 
  29: func main() {
  30: 	serviceName := viper.GetString("order.service-name")
  31: 
  32: 	ctx, cancel := context.WithCancel(context.Background())
  33: 	defer cancel()
  34: 
  35: 	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
  36: 	if err != nil {
  37: 		logrus.Fatal(err)
  38: 	}
  39: 	defer shutdown(ctx)
  40: 
  41: 	application, cleanup := service.NewApplication(ctx)
  42: 	defer cleanup()
  43: 
  44: 	deregisterFunc, err := discovery.RegisterToConsul(ctx, serviceName)
  45: 	if err != nil {
  46: 		logrus.Fatal(err)
  47: 	}
  48: 	defer func() {
  49: 		_ = deregisterFunc()
  50: 	}()
  51: 
  52: 	ch, closeCh := broker.Connect(
  53: 		viper.GetString("rabbitmq.user"),
  54: 		viper.GetString("rabbitmq.password"),
  55: 		viper.GetString("rabbitmq.host"),
  56: 		viper.GetString("rabbitmq.port"),
  57: 	)
  58: 	defer func() {
  59: 		_ = ch.Close()
  60: 		_ = closeCh()
  61: 	}()
  62: 	go consumer.NewConsumer(application).Listen(ch)
  63: 
  64: 	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
  65: 		svc := ports.NewGRPCServer(application)
  66: 		orderpb.RegisterOrderServiceServer(server, svc)
  67: 	})
  68: 
  69: 	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
  70: 		router.StaticFile("/success", "../../public/success.html")
  71: 		ports.RegisterHandlersWithOptions(router, HTTPServer{
  72: 			app: application,
  73: 		}, ports.GinServerOptions{
  74: 			BaseURL:      "/api",
  75: 			Middlewares:  nil,
  76: 			ErrorHandler: nil,
  77: 		})
  78: 	})
  79: 
  80: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 8 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 9 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 10 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 11 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 12 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 13 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 14 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 15 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 16 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 17 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 18 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 19 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 20 | 语法块结束：关闭 import 或参数列表。 |
| 21 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 22 | init 函数：包加载时自动执行，常用于初始化与注册。 |
| 23 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 24 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 25 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 26 | 代码块结束：收束当前函数、分支或类型定义。 |
| 27 | 代码块结束：收束当前函数、分支或类型定义。 |
| 28 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 29 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 30 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 31 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 32 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 33 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 34 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 35 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 36 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 37 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 38 | 代码块结束：收束当前函数、分支或类型定义。 |
| 39 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 40 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 41 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 42 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 43 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 44 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 45 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 46 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 47 | 代码块结束：收束当前函数、分支或类型定义。 |
| 48 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 49 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 50 | 代码块结束：收束当前函数、分支或类型定义。 |
| 51 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 52 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 53 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 54 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 55 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 56 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 57 | 语法块结束：关闭 import 或参数列表。 |
| 58 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 59 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 60 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 61 | 代码块结束：收束当前函数、分支或类型定义。 |
| 62 | goroutine 启动：引入并发执行，需关注生命周期管理。 |
| 63 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 64 | goroutine 启动：引入并发执行，需关注生命周期管理。 |
| 65 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 66 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 67 | 代码块结束：收束当前函数、分支或类型定义。 |
| 68 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 69 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 70 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 71 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 72 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 73 | 代码块结束：收束当前函数、分支或类型定义。 |
| 74 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 75 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 76 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 77 | 代码块结束：收束当前函数、分支或类型定义。 |
| 78 | 代码块结束：收束当前函数、分支或类型定义。 |
| 79 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 80 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/payment/main.go

~~~go
   1: package main
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/broker"
   7: 	"github.com/ghost-yu/go_shop_second/common/config"
   8: 	"github.com/ghost-yu/go_shop_second/common/logging"
   9: 	"github.com/ghost-yu/go_shop_second/common/server"
  10: 	"github.com/ghost-yu/go_shop_second/common/tracing"
  11: 	"github.com/ghost-yu/go_shop_second/payment/infrastructure/consumer"
  12: 	"github.com/ghost-yu/go_shop_second/payment/service"
  13: 	"github.com/sirupsen/logrus"
  14: 	"github.com/spf13/viper"
  15: )
  16: 
  17: func init() {
  18: 	logging.Init()
  19: 	if err := config.NewViperConfig(); err != nil {
  20: 		logrus.Fatal(err)
  21: 	}
  22: }
  23: 
  24: func main() {
  25: 	serviceName := viper.GetString("payment.service-name")
  26: 	ctx, cancel := context.WithCancel(context.Background())
  27: 	defer cancel()
  28: 
  29: 	serverType := viper.GetString("payment.server-to-run")
  30: 
  31: 	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
  32: 	if err != nil {
  33: 		logrus.Fatal(err)
  34: 	}
  35: 	defer shutdown(ctx)
  36: 
  37: 	application, cleanup := service.NewApplication(ctx)
  38: 	defer cleanup()
  39: 
  40: 	ch, closeCh := broker.Connect(
  41: 		viper.GetString("rabbitmq.user"),
  42: 		viper.GetString("rabbitmq.password"),
  43: 		viper.GetString("rabbitmq.host"),
  44: 		viper.GetString("rabbitmq.port"),
  45: 	)
  46: 	defer func() {
  47: 		_ = ch.Close()
  48: 		_ = closeCh()
  49: 	}()
  50: 
  51: 	go consumer.NewConsumer(application).Listen(ch)
  52: 
  53: 	paymentHandler := NewPaymentHandler(ch)
  54: 	switch serverType {
  55: 	case "http":
  56: 		server.RunHTTPServer(serviceName, paymentHandler.RegisterRoutes)
  57: 	case "grpc":
  58: 		logrus.Panic("unsupported server type: grpc")
  59: 	default:
  60: 		logrus.Panic("unreachable code")
  61: 	}
  62: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 8 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 9 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 10 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 11 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 12 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 13 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 14 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 15 | 语法块结束：关闭 import 或参数列表。 |
| 16 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 17 | init 函数：包加载时自动执行，常用于初始化与注册。 |
| 18 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 19 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 20 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 21 | 代码块结束：收束当前函数、分支或类型定义。 |
| 22 | 代码块结束：收束当前函数、分支或类型定义。 |
| 23 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 24 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 25 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 26 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 27 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 28 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 29 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 30 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 31 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 32 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 33 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 34 | 代码块结束：收束当前函数、分支或类型定义。 |
| 35 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 36 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 37 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 38 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 39 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 40 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 41 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 42 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 43 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 44 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 45 | 语法块结束：关闭 import 或参数列表。 |
| 46 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 47 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 48 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 49 | 代码块结束：收束当前函数、分支或类型定义。 |
| 50 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 51 | goroutine 启动：引入并发执行，需关注生命周期管理。 |
| 52 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 53 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 54 | 多分支选择：按状态或类型分流执行路径。 |
| 55 | 分支标签：定义 switch 的命中条件。 |
| 56 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 57 | 分支标签：定义 switch 的命中条件。 |
| 58 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 59 | 默认分支：兜底处理，防止未知输入静默失败。 |
| 60 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 61 | 代码块结束：收束当前函数、分支或类型定义。 |
| 62 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/stock/main.go

~~~go
   1: package main
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/config"
   7: 	"github.com/ghost-yu/go_shop_second/common/discovery"
   8: 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
   9: 	"github.com/ghost-yu/go_shop_second/common/logging"
  10: 	"github.com/ghost-yu/go_shop_second/common/server"
  11: 	"github.com/ghost-yu/go_shop_second/common/tracing"
  12: 	"github.com/ghost-yu/go_shop_second/stock/ports"
  13: 	"github.com/ghost-yu/go_shop_second/stock/service"
  14: 	"github.com/sirupsen/logrus"
  15: 	"github.com/spf13/viper"
  16: 	"google.golang.org/grpc"
  17: )
  18: 
  19: func init() {
  20: 	logging.Init()
  21: 	if err := config.NewViperConfig(); err != nil {
  22: 		logrus.Fatal(err)
  23: 	}
  24: }
  25: 
  26: func main() {
  27: 	serviceName := viper.GetString("stock.service-name")
  28: 	serverType := viper.GetString("stock.server-to-run")
  29: 
  30: 	ctx, cancel := context.WithCancel(context.Background())
  31: 	defer cancel()
  32: 
  33: 	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
  34: 	if err != nil {
  35: 		logrus.Fatal(err)
  36: 	}
  37: 	defer shutdown(ctx)
  38: 
  39: 	application := service.NewApplication(ctx)
  40: 
  41: 	deregisterFunc, err := discovery.RegisterToConsul(ctx, serviceName)
  42: 	if err != nil {
  43: 		logrus.Fatal(err)
  44: 	}
  45: 	defer func() {
  46: 		_ = deregisterFunc()
  47: 	}()
  48: 
  49: 	switch serverType {
  50: 	case "grpc":
  51: 		server.RunGRPCServer(serviceName, func(server *grpc.Server) {
  52: 			svc := ports.NewGRPCServer(application)
  53: 			stockpb.RegisterStockServiceServer(server, svc)
  54: 		})
  55: 	case "http":
  56: 		// 暂时不用
  57: 	default:
  58: 		panic("unexpected server type")
  59: 	}
  60: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 8 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 9 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 10 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 11 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 12 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 13 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 14 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 15 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 16 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 17 | 语法块结束：关闭 import 或参数列表。 |
| 18 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 19 | init 函数：包加载时自动执行，常用于初始化与注册。 |
| 20 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 21 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 22 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 23 | 代码块结束：收束当前函数、分支或类型定义。 |
| 24 | 代码块结束：收束当前函数、分支或类型定义。 |
| 25 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 26 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 27 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 28 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 29 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 30 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 31 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 32 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 33 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 34 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 35 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 36 | 代码块结束：收束当前函数、分支或类型定义。 |
| 37 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 38 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 39 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 40 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 41 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 42 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 43 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 44 | 代码块结束：收束当前函数、分支或类型定义。 |
| 45 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 46 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 47 | 代码块结束：收束当前函数、分支或类型定义。 |
| 48 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 49 | 多分支选择：按状态或类型分流执行路径。 |
| 50 | 分支标签：定义 switch 的命中条件。 |
| 51 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 52 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 53 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 54 | 代码块结束：收束当前函数、分支或类型定义。 |
| 55 | 分支标签：定义 switch 的命中条件。 |
| 56 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 57 | 默认分支：兜底处理，防止未知输入静默失败。 |
| 58 | panic：不可恢复错误路径，常用于关键初始化失败。 |
| 59 | 代码块结束：收束当前函数、分支或类型定义。 |
| 60 | 代码块结束：收束当前函数、分支或类型定义。 |


