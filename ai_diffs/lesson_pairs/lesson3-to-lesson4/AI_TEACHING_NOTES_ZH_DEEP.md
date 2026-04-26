# `lesson3 -> lesson4` 独立讲义（增强细化版）

这一组一定要从“业务功能”里抽离出来看。

因为它的重点根本不是“下单能不能用”，而是：

`项目第一次建立了协议层、代码生成流程、统一配置读取方式。`

如果你是 Go 小白，这一组最容易误解成：

- 怎么突然多了这么多文件
- 为什么一堆 `.pb.go`、`.gen.go`
- 这些东西是不是都要我手写、手读

答案是：

`不用。`

这组真正要学的是：

1. 什么是协议文件
2. 为什么要先写协议再写实现
3. 为什么生成文件不是阅读重点
4. Go 项目里配置读取为什么要尽早统一

## 1. 这组差异到底在做什么

一句话概括：

`lesson3 -> lesson4` 把 order 服务的 gRPC 协议、HTTP 协议、代码生成脚本、统一配置读取这几件“地基级能力”补齐了。`

你先记住四个核心文件：

1. [order.proto](/g:/shi/go_shop_second/api/orderpb/order.proto)
2. [genproto.sh](/g:/shi/go_shop_second/scripts/genproto.sh)
3. [genopenapi.sh](/g:/shi/go_shop_second/scripts/genopenapi.sh)
4. [viper.go](/g:/shi/go_shop_second/internal/common/config/viper.go)

这四个文件比那些生成出来的 `.pb.go`、`.gen.go` 重要得多。

因为它们分别回答了 4 个问题：

1. order 服务准备暴露哪些能力
2. gRPC 代码怎么生成
3. HTTP/OpenAPI 代码怎么生成
4. 配置从哪里读、怎么读

## 2. 为什么这组会突然多出很多文件

因为从这一组开始，项目第一次真正进入“协议优先 + 生成优先”的工程方式。

也就是说，作者不是想手写：

- gRPC 的请求结构体
- gRPC 的 server/client 接口
- OpenAPI 的 handler 接口
- OpenAPI 的 client 类型

而是先定义：

- `proto`
- `openapi yaml`

再让工具自动生成 Go 代码。

这就是为什么你会突然看到：

- `order.pb.go`
- `order_grpc.pb.go`
- `openapi_api.gen.go`
- `openapi_types.gen.go`
- `openapi_client.gen.go`

这些文件的本质都是：

`协议定义的机械展开结果。`

所以这组阅读时必须遵守一个原则：

`先看协议和生成脚本，最后才看生成结果。`

## 3. 先看 `order.proto`

文件：[order.proto](/g:/shi/go_shop_second/api/orderpb/order.proto)

最开头：

```proto
syntax = "proto3";
package orderpb;

option go_package = "github.com/ghost-yu/go_shop_second/internal/common/genproto/orderpb";
```

### 3.1 `syntax = "proto3"` 是什么

这一行的意思很简单：

`这份协议文件使用 protobuf 第 3 版语法。`

对你来说，现在先不用深究 proto2 / proto3 全部区别，
先知道它是：

- 一种定义消息结构的标准格式
- 也是 gRPC 最常见的协议描述方式

也就是说，这不是 Go 语法，
而是“接口契约语法”。

### 3.2 `package orderpb` 在做什么

```proto
package orderpb;
```

这是 proto 层自己的包名。

以后你在 Go 代码里会看到：

- `orderpb.Order`
- `orderpb.CreateOrderRequest`
- `orderpb.OrderServiceServer`

这些前缀就是从这里来的。

### 3.3 `go_package` 为什么不能乱写

```proto
option go_package = "github.com/ghost-yu/go_shop_second/internal/common/genproto/orderpb";
```

这一行是给 Go 生成器看的。

它在明确告诉生成器：

`你生成的 Go 代码属于哪个导入路径。`

如果这行配错，后面很容易出现：

- 生成文件包路径奇怪
- import 不对
- 多模块引用混乱
- 生成出来的代码能编译但不好用

所以这不是“可选美化项”，
它是生成代码是否可维护的重要设置。

## 4. `service OrderService` 到底在定义什么

代码：

```proto
service OrderService {
  rpc CreateOrder(CreateOrderRequest) returns (google.protobuf.Empty);
  rpc GetOrder(GetOrderRequest) returns (Order);
  rpc UpdateOrder(Order) returns (google.protobuf.Empty);
}
```

这一段非常关键。

它在定义：

`order 服务对外暴露哪些远程调用能力。`

也就是说，这不是“内部方法”，
而是“别的服务或者客户端可以远程调我的方法”。

逐个看：

### 4.1 `CreateOrder`

```proto
rpc CreateOrder(CreateOrderRequest) returns (google.protobuf.Empty);
```

含义是：

- 输入：创建订单请求
- 输出：空响应

为什么返回空响应？

这说明作者此时的设计是：

`创建成功与否最重要，暂时不把完整订单对象作为创建结果返回。`

这在 RPC 设计里很常见，
但它也有代价：

- 如果调用方创建后想拿订单详情，通常还要再查一次

### 4.2 `GetOrder`

```proto
rpc GetOrder(GetOrderRequest) returns (Order);
```

这个就很直接：

- 输入查询条件
- 输出完整订单对象

你要开始形成这种意识：

- `CreateOrder` 是写
- `GetOrder` 是读

虽然现在还没正式讲 CQRS，
但接口层已经有这个味道了。

### 4.3 `UpdateOrder`

```proto
rpc UpdateOrder(Order) returns (google.protobuf.Empty);
```

这是一种很朴素的更新接口写法：

- 直接把订单对象传进来
- 由服务端自己决定怎么更新

在项目早期这样写很常见，
因为实现快、改动少。

## 5. `message` 为什么要认真看

### 5.1 `CreateOrderRequest`

```proto
message CreateOrderRequest {
  string CustomerID = 1;
  repeated ItemWithQuantity Items = 2;
}
```

这在表达一个非常实际的业务意思：

`创建订单至少需要知道是谁下单，以及买了什么。`

这里最值得你记的点有两个：

1. `repeated` 表示列表
2. `= 1`、`= 2` 这种字段编号不是装饰，而是 protobuf 协议的一部分

字段编号以后不能随便改、不能随便复用，
这个意识现在就要建立。

### 5.2 `GetOrderRequest`

```proto
message GetOrderRequest {
  string OrderID = 1;
  string CustomerID = 2;
}
```

这里说明查询订单时，不只是用订单 ID，
还带上了 customerID。

这背后其实已经有一点业务意味：

`订单查询和客户身份是绑定的。`

### 5.3 `ItemWithQuantity` 和 `Item`

```proto
message ItemWithQuantity {
  string ID = 1;
  int32 Quantity = 2;
}
```

这个结构很像“创建订单时用户传来的最小商品信息”。

而：

```proto
message Item {
  string ID = 1;
  string Name = 2;
  int32 Quantity = 3;
  string PriceID = 4;
}
```

这个就更像“系统内部已经更丰富的商品信息”。

你要看出来：

- `ItemWithQuantity` 更偏输入
- `Item` 更偏系统处理后的完整对象

### 5.4 `Order`

```proto
message Order {
  string ID = 1;
  string CustomerID = 2;
  string Status = 3;
  repeated Item Items = 4;
}
```

这说明在 `lesson4` 时，作者对订单的理解还很早期：

- 订单号
- 客户
- 状态
- 商品列表

此时还没有：

- `PaymentLink`
- 支付时间
- 金额细节
- 更复杂的状态字段

这很正常，说明系统还在早期建模阶段。

## 6. 看 `genproto.sh`

文件：[genproto.sh](/g:/shi/go_shop_second/scripts/genproto.sh)

这个脚本的核心价值是：

`让 proto -> Go 代码生成变成可复用、可重复执行的工程流程。`

### 6.1 为什么不是手写 `order.pb.go`

因为 `.pb.go` 本质上是：

`协议文件的机械翻译结果。`

这类文件手写会带来三个问题：

1. 容易错
2. 难维护
3. 协议一变所有地方都要手改

所以必须交给生成器。

### 6.2 `find . -type f -name '*.proto'` 在做什么

它的含义是：

`扫描仓库里的 proto 文件。`

这说明这个脚本不是只针对一个文件写死的，
它有一点“面向多个 proto”的意图了。

### 6.3 为什么先清理输出目录

脚本里有清理逻辑：

```bash
run rm -rf $go_out
```

原因是生成文件如果不清理旧版本，
很容易出现：

- 旧文件残留
- 新旧版本混在一起
- import 或编译结果变脏

所以这里的思路是：

`每次生成都尽量从干净状态开始。`

### 6.4 `protoc` 命令真正生成了什么

```bash
--go_out=...
--go-grpc_out=...
```

意思是：

- `--go_out` 生成消息结构体代码
- `--go-grpc_out` 生成 gRPC server/client 接口代码

也就是：

- `order.pb.go`
- `order_grpc.pb.go`

这两类文件分别来自不同生成阶段。

## 7. 看 `genopenapi.sh`

文件：[genopenapi.sh](/g:/shi/go_shop_second/scripts/genopenapi.sh)

这个脚本回答的问题是：

`如果我要让 order 也走 HTTP 风格接口，怎么生成服务端和客户端代码？`

### 7.1 为什么指定 `gin-server`

```bash
GEN_SERVER=(
  "gin-server"
)
```

意思是：

`HTTP server 端代码生成目标框架是 Gin。`

也就是说，
项目不准备自己手搓一层随意的 HTTP handler 风格，
而是让 OpenAPI 生成器直接配合 Gin。

### 7.2 为什么同时生成 server 和 client

脚本里你会看到既生成：

- server 端 `openapi_api.gen.go`
- client 端 `openapi_client.gen.go`

这说明作者的思路不是只写文档，
而是：

`让 HTTP 协议也变成一套可以被代码消费的契约。`

这和 proto 的思路其实是一样的。

## 8. 看 `viper.go`

文件：[viper.go](/g:/shi/go_shop_second/internal/common/config/viper.go)

代码：

```go
func NewViperConfig() error {
    viper.SetConfigName("global")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("../common/config")
    viper.AutomaticEnv()
    return viper.ReadInConfig()
}
```

这组里它虽然短，
但意义其实很大。

### 8.1 为什么要做成公共函数

因为如果每个服务都自己写一套：

- 配置文件名
- 配置路径
- 是否读环境变量

后面一定会乱。

所以这里开始把“如何读配置”抽到 common 里。

### 8.2 `AutomaticEnv()` 为什么你要记住

```go
viper.AutomaticEnv()
```

这表示：

`除了 YAML 配置文件，环境变量也会参与配置覆盖。`

这对后面接：

- Stripe key
- 数据库密码
- Docker 环境
- CI/CD 环境

都非常重要。

也就是说，这一组虽然还没真正大量用环境变量，
但已经把能力准备好了。

## 9. 这组最容易混淆的点

1. `.pb.go`、`.gen.go` 是生成结果，不是手写业务代码
2. `proto` 和 `openapi` 都是协议，但服务对象不同
3. `go_package` 是给 Go 生成器看的，不是 proto 包名本身
4. 配置读取这时看起来简单，但它是后面所有服务启动的基础

## 10. 为什么这一组要这么做

因为 `lesson3` 只有“项目能起、配置能读”的最初地基，
但没有真正定义：

- order 服务到底暴露哪些能力
- HTTP 和 gRPC 的契约长什么样
- 这些契约以后怎么稳定地生成代码

所以这组必须先把“共同语言”建立起来。

微服务项目里，经常是：

`协议先行，业务后补。`

## 11. 这组你最该记住的结论

1. `lesson3 -> lesson4` 的核心是协议层和生成流程，不是业务层
2. `order.proto` 是这组最重要的阅读入口
3. `genproto.sh` 和 `genopenapi.sh` 让代码生成变成固定流程
4. `viper.go` 让配置读取开始统一
5. 大量生成文件不要当主阅读入口

## 12. 最推荐你的复读顺序

1. [order.proto](/g:/shi/go_shop_second/api/orderpb/order.proto)
2. [genproto.sh](/g:/shi/go_shop_second/scripts/genproto.sh)
3. [genopenapi.sh](/g:/shi/go_shop_second/scripts/genopenapi.sh)
4. [viper.go](/g:/shi/go_shop_second/internal/common/config/viper.go)
5. [Makefile](/g:/shi/go_shop_second/Makefile)

## 13. 一句话收尾

这组差异的本质是：

`项目第一次建立了“先定义协议，再生成代码，再往里填实现”的开发方式。`