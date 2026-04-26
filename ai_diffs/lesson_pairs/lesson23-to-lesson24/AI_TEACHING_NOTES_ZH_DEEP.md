# `lesson23 -> lesson24` 独立讲义（详细注释版）

这一组只有两个文件变化，但它解决的是一个很真实、很常见、也很基础设施味的问题：

`服务启动顺序问题。`

前面课程走到这里时，项目已经有这些依赖关系了：

- `order` 依赖 `stock` 的 gRPC
- `payment` 依赖 `order` 的 gRPC

这意味着什么？

意味着你不能再把每个服务都想成“自己启动自己就行”。

因为当一个服务在启动过程中要去创建别的服务的 gRPC client 时，如果对方还没起来，就很容易：
- 发现不到地址
- 拨号失败
- 启动直接 panic

字幕 `25.TXT` 讲的就是这个问题：
- 比如 `payment` 比 `order` 先起
- 但 `payment` 初始化时又要先拿 `order` 的 gRPC client
- 而 `order` 还没把自己的地址注册到 Consul
- 那 `payment` 拿到的可能就是空地址，或者直接失败

所以这一组在做的事不是“优化一点点体验”，而是：

`给服务启动过程加上一层“等待依赖服务端口可用”的保护，避免因为启动时序导致整个服务直接死掉。`

这组你可以把它理解成：

`在真正引入更成熟的服务发现/健康检查/启动编排之前，先加一个实用的本地等待机制。`

## 1. 原始 diff 去哪里看

原始差异在这里：
[diff.md](/g:/shi/go_shop_second/ai_diffs/lesson_pairs/lesson23-to-lesson24/diff.md)

这次我继续严格按固定要求来：
- 每个文件先贴自己的 diff
- 再贴带中文注释的代码/关键代码
- 再讲旧代码做什么、新代码做什么、为什么这样改
- 并结合字幕 `25.TXT`

## 2. 正确阅读顺序

这一组虽然文件少，但建议按下面顺序读：

1. [internal/common/client/grpc.go](/g:/shi/go_shop_second/internal/common/client/grpc.go)
2. [internal/common/config/global.yaml](/g:/shi/go_shop_second/internal/common/config/global.yaml)

原因：
- `grpc.go` 是这组的核心逻辑，真正实现了等待端口可用
- `global.yaml` 解释这个等待超时从哪里配置

## 3. 总调用链

这组的“基础设施调用链”是：

```text
payment / order 服务启动
-> 启动过程中要 NewOrderGRPCClient / NewStockGRPCClient
-> 在真正 discovery + grpc.NewClient 之前
-> 先 WaitForOrderGRPCClient / WaitForStockGRPCClient
-> 内部循环用 net.Dial 检查目标 TCP 端口是否可连
-> 如果在 timeout 内可连，则继续创建 gRPC client
-> 如果超时还不可连，则返回 “xxx grpc not available”
-> 上层服务启动失败，但失败原因更明确，而且不会因为时序随机成功/随机失败
```

注意这条链真正解决的不是“服务发现”本身，而是：

`在依赖服务尚未准备好时，不要立刻往下执行后续 client 初始化。`

它做的是一个“准备就绪前等待”的动作。

## 4. 这组差异到底在干什么

一句话概括：

`lesson23 -> lesson24` 给 gRPC client 初始化加了一层端口等待机制，用来缓解服务启动顺序导致的依赖服务未就绪问题。`

这个问题为什么现实？

因为在多服务系统里，启动时经常会出现这样的情况：
- 程序 A 想连程序 B
- 但程序 B 还没监听端口
- 或者 B 虽然进程起来了，但还没把服务地址注册进 Consul

这时如果 A 直接继续初始化，就会出错。

所以这组不是业务功能升级，而是：

`给启动阶段补一个等待依赖的保护层。`

## 5. 关键文件详细讲解

### 5.1 [internal/common/client/grpc.go](/g:/shi/go_shop_second/internal/common/client/grpc.go)

这是这组的核心文件。

先看这个文件自己的 diff。

```diff
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
@@
 func NewStockGRPCClient(ctx context.Context) (client stockpb.StockServiceClient, close func() error, err error) {
+	if !WaitForStockGRPCClient(viper.GetDuration("dial-grpc-timeout") * time.Second) {
+		return nil, nil, errors.New("stock grpc not available")
+	}
 	grpcAddr, err := discovery.GetServiceAddr(ctx, viper.GetString("stock.service-name"))
@@
 func NewOrderGRPCClient(ctx context.Context) (client orderpb.OrderServiceClient, close func() error, err error) {
+	if !WaitForOrderGRPCClient(viper.GetDuration("dial-grpc-timeout") * time.Second) {
+		return nil, nil, errors.New("order grpc not available")
+	}
 	grpcAddr, err := discovery.GetServiceAddr(ctx, viper.GetString("order.service-name"))
@@
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
```

你第一次看可能会觉得：

`这不就是多了个 waitFor 吗？`

但这段代码其实把“服务依赖未就绪时怎么办”这件事，从“完全不处理”提升成了“显式处理”。

下面我把当前代码里相关部分贴出来，并直接写细注释。

```go
package client

import (
    "context"
    "errors"
    "net"
    "time"

    "github.com/ghost-yu/go_shop_second/common/discovery"
    "github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
    "github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
    "github.com/sirupsen/logrus"
    "github.com/spf13/viper"
    "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func NewStockGRPCClient(ctx context.Context) (client stockpb.StockServiceClient, close func() error, err error) {
    // 新增的第一层保护：先等 stock 的 gRPC 端口可用。
    // 如果在超时时间内端口始终没起来，直接返回错误，不继续往下做 discovery 和建连。
    if !WaitForStockGRPCClient(viper.GetDuration("dial-grpc-timeout") * time.Second) {
        return nil, nil, errors.New("stock grpc not available")
    }

    // 等到端口可用之后，再去 Consul 做服务发现。
    grpcAddr, err := discovery.GetServiceAddr(ctx, viper.GetString("stock.service-name"))
    if err != nil {
        return nil, func() error { return nil }, err
    }
    if grpcAddr == "" {
        logrus.Warn("empty grpc addr for stock grpc")
    }

    opts := grpcDialOpts(grpcAddr)
    conn, err := grpc.NewClient(grpcAddr, opts...)
    if err != nil {
        return nil, func() error { return nil }, err
    }
    return stockpb.NewStockServiceClient(conn), conn.Close, nil
}

func NewOrderGRPCClient(ctx context.Context) (client orderpb.OrderServiceClient, close func() error, err error) {
    // 对 order gRPC client 也做同样的等待。
    if !WaitForOrderGRPCClient(viper.GetDuration("dial-grpc-timeout") * time.Second) {
        return nil, nil, errors.New("order grpc not available")
    }

    grpcAddr, err := discovery.GetServiceAddr(ctx, viper.GetString("order.service-name"))
    if err != nil {
        return nil, func() error { return nil }, err
    }
    if grpcAddr == "" {
        logrus.Warn("empty grpc addr for order grpc")
    }

    opts := grpcDialOpts(grpcAddr)
    conn, err := grpc.NewClient(grpcAddr, opts...)
    if err != nil {
        return nil, func() error { return nil }, err
    }
    return orderpb.NewOrderServiceClient(conn), conn.Close, nil
}
```

### 这段代码到底在解决什么问题

旧代码的流程大致是：

1. 一启动就去 `discovery.GetServiceAddr(...)`
2. 然后立刻 `grpc.NewClient(...)`

问题在于：
- 依赖服务如果还没起来
- 或还没注册到 Consul
- 你拿到的地址可能为空或者根本不可连
- 后续初始化就直接出错

新代码做的就是把“真的去建 client 之前”，插入一个等待步骤。

你可以把它理解成：

`别急着建 gRPC client，先看看对面的端口至少活了没有。`

这不是最终完美方案，但它确实能缓解启动顺序问题。

---

继续看新增的等待函数：

```go
func WaitForOrderGRPCClient(timeout time.Duration) bool {
    logrus.Infof("waiting for order grpc client, timeout: %v seconds", timeout.Seconds())
    return waitFor(viper.GetString("order.grpc-addr"), timeout)
}

func WaitForStockGRPCClient(timeout time.Duration) bool {
    logrus.Infof("waiting for stock grpc client, timeout: %v seconds", timeout.Seconds())
    return waitFor(viper.GetString("stock.grpc-addr"), timeout)
}
```

这两个函数本身不复杂，它们只是：
- 打日志
- 从配置里拿端口地址
- 转给通用的 `waitFor(...)`

真正的核心在 `waitFor(...)`。

```go
func waitFor(addr string, timeout time.Duration) bool {
    // 用一个 channel 表示“端口已经可用”。
    portAvailable := make(chan struct{})

    // 另一个 channel 用来表示“总超时到点了”。
    timeoutCh := time.After(timeout)

    go func() {
        for {
            // 每轮循环先看看是不是已经超时。
            select {
            case <-timeoutCh:
                return
            default:
                // 没超时就继续
            }

            // 尝试用 TCP 去连这个地址。
            _, err := net.Dial("tcp", addr)
            if err == nil {
                // 只要 TCP 能连通，就说明这个端口已经开始监听。
                close(portAvailable)
                return
            }

            // 不能一直死命重试，否则 CPU 会白转，所以每 200ms 试一次。
            time.Sleep(200 * time.Millisecond)
        }
    }()

    // 主协程在这里等两个结果：
    // 要么端口可用，要么超时。
    select {
    case <-portAvailable:
        return true
    case <-timeoutCh:
        return false
    }
}
```

这段代码建议你慢慢看，因为它包含了几个基础设施里很常见的概念。

#### 概念 1：为什么要用 `net.Dial("tcp", addr)`

因为 gRPC 本质上也是跑在 TCP 连接上的。

所以这里没有先去做“业务级别”的请求，而只是做最基础的一件事：

`这个端口现在有没有开始监听？`

如果 `net.Dial("tcp", addr)` 成功，至少说明：
- 对方进程已经把这个 TCP 端口占起来了
- 不再是“完全没起来”的状态

这比直接盲目初始化 gRPC client 更稳一点。

#### 概念 2：为什么这里用 `channel + goroutine + select`

这是 Go 里非常典型的一种等待模式：

- goroutine 里不断探测
- channel 用来通知“探测成功了”
- 另一个 channel 表示“超时到了”
- `select` 用来等哪个先发生

你可以把它理解成：

`两件事在赛跑：`
- 端口先变可用
- 还是超时时间先到

谁先发生，就决定函数返回 `true` 还是 `false`。

#### 概念 3：为什么不能一直 while 死循环

如果你不 `Sleep(200 * time.Millisecond)`，那就是疯狂空转：
- 一秒钟重试成千上万次
- 白白吃 CPU

所以基础设施里这种“轮询等待”通常都会加间隔。

这里的 `200ms` 你可以理解成：
- 不追求毫秒级极致响应
- 但也不要等太久

是一个工程上比较朴素的折中值。

### 这里要特别指出它的局限

这段代码不是完美方案，它只是一个“先够用”的方案。

为什么？

因为它只检查：
- TCP 端口是不是开了

它没有检查：
- 服务是不是已经完全 ready
- gRPC handler 是不是都注册好了
- Consul 里是不是已经注册成功
- 这个服务是不是内部依赖也都准备齐了

也就是说：

`端口可连 != 服务真正完全可用`

这是你做基础设施时必须明白的一个非常重要的区别。

所以这组更准确地说，是：

`用“端口探活”缓解启动时序问题，而不是彻底解决服务就绪问题。`

### 这一文件你最该带走的结论

`grpc.go` 这组真正完成的是：把“依赖服务没起来就直接撞上去”改成了“先等端口可用，再继续初始化 gRPC client”。`

---

### 5.2 [internal/common/config/global.yaml](/g:/shi/go_shop_second/internal/common/config/global.yaml)

这个文件只改了一行，但这一行决定了等待机制的容忍时间。

先看 diff。

```diff
diff --git a/internal/common/config/global.yaml b/internal/common/config/global.yaml
index 7dd93c9..1bf2f6e 100644
--- a/internal/common/config/global.yaml
+++ b/internal/common/config/global.yaml
@@ -1,4 +1,5 @@
 fallback-grpc-addr: 127.0.0.1:3030
+dial-grpc-timeout: 10
```

当前代码里对应就是：

```yaml
dial-grpc-timeout: 10
```

这个值在代码里是怎么用的？

```go
viper.GetDuration("dial-grpc-timeout") * time.Second
```

所以这里的 `10` 表示：
- 先读成一个 duration 值
- 再乘上 `time.Second`
- 最终含义就是 10 秒

也就是说：

`如果 10 秒之内对方端口还没起来，就认为依赖服务不可用。`

#### 为什么把这个东西做成配置，而不是写死在代码里

因为“等多久算合理”是环境相关的。

比如：
- 本地开发机器慢一点，可能想等更久
- 某些环境服务起得快，可能几秒就够
- 以后如果接入容器编排，也许又会调整

所以这里把 timeout 提取成配置是对的。

这也是你需要建立的基础设施意识：

`凡是跟环境和部署节奏相关的参数，优先考虑做成配置。`

#### 这里补一个配置基础知识：为什么不直接写 `10s`

你可能会问：

`为什么不直接在 yaml 写成 10s？`

这取决于项目怎么取值。

当前代码写的是：

```go
viper.GetDuration("dial-grpc-timeout") * time.Second
```

这说明作者当前是把配置值当“秒数数字”在用。

如果以后想更直观，也可以改成：
- 配置直接写 `10s`
- 代码不再额外乘 `time.Second`

但当前这版还没走到那一步。

### 这一文件你最该带走的结论

`global.yaml` 这一行的意义不是“多了个配置项”，而是“服务启动等待策略第一次被参数化了”。`

## 6. 第三方库和基础设施解释

### 6.1 `net.Dial`

这是 Go 标准库里的网络拨号函数。

在这里它不是拿来做完整协议通信，而只是拿来探测：

`这个 TCP 地址现在能不能连上。`

所以你可以把它理解成最底层、最原始的连通性检查。

### 6.2 `time.After`

`time.After(timeout)` 会返回一个 channel。

过了指定时间后，这个 channel 就会收到一个时间值。

所以它很适合拿来做超时控制。

这也是 Go 并发里很经典的用法：
- 某个动作成功 channel
- 某个超时 channel
- `select` 谁先到就处理谁

### 6.3 `select`

你可以把 `select` 理解成：

`等待多个 channel 中的任意一个先准备好。`

这里它的作用是：
- 要么等到端口可用
- 要么等到超时

这比单纯 `sleep + if` 更适合做并发等待。

### 6.4 Consul 和端口等待的关系

这组很容易让你误会成：

`既然已经有 Consul，为什么还要手动 waitFor 端口？`

这是因为它们在解决的问题层次不完全一样。

- Consul 解决的是“服务地址发现”
- `waitFor` 解决的是“依赖服务是否至少已经开始监听端口”

而当前项目的真实问题是：
- 服务可能还没注册进 Consul
- 也可能虽然要注册，但在那之前端口还没准备好

所以作者先用一个简单的等待策略把问题缓一缓。

更成熟的做法通常会结合：
- 健康检查
- 注册中心状态
- 重试退避
- 容器编排里的 readiness probe

当前这组还没走到那一步。

## 7. 为什么这么设计

这组设计本质上是一个折中方案：

`不先引入复杂的服务编排，也不马上重构 discovery 流程，而是先在 gRPC client 初始化前加一层端口等待。`

这样做的好处：
- 改动小
- 立刻有效
- 非常适合本地多服务联调

坏处：
- 只解决“端口没开”的问题
- 不保证“服务完全 ready”
- 还是偏本地开发友好，而不是生产级就绪方案

所以这组是典型的教学项目式演进：

`先解决最疼的联调问题，再逐步收敛成更成熟的基础设施方案。`

## 8. 当前还不完美的地方

这一组虽然很实用，但你必须看到它的边界：

1. 只检查端口可连，不检查真正业务就绪
2. 如果服务端口起来了，但内部逻辑还没 ready，仍然可能失败
3. `WaitForOrderGRPCClient` / `WaitForStockGRPCClient` 现在还是分别写的，抽象层次不高
4. discovery 和 waitFor 的衔接还比较粗糙
5. 超时策略是固定轮询，不是更成熟的退避策略

所以这组更准确的定位是：

`启动依赖问题的第一版工程补丁。`

## 9. 这组最该带走的知识点

1. 多服务系统里，启动顺序和依赖就绪是一个真实问题，不是“运气不好”
2. gRPC client 初始化前先等依赖端口可用，是一种很常见的缓解方案
3. `channel + goroutine + select + time.After` 是 Go 里实现等待/超时控制的经典组合
4. 端口可连不等于服务完全 ready，这两个概念要分清
5. timeout 这类和环境相关的参数，最好配置化
6. 这一组解决的是“启动阶段别太脆弱”，不是“服务发现已经完美”

## 10. 一句话收住这组

`lesson23 -> lesson24` 的本质，不是“多了个 waitFor 函数”，而是“项目第一次正式处理服务启动依赖时序问题，让 gRPC client 在依赖服务未就绪时不再直接撞死”。`