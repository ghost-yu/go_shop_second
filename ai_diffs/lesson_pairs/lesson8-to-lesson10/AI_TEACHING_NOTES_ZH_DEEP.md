# `lesson8 -> lesson10` 独立讲义（伪注释阅读版）

这次我按你的要求改回独立讲义，不再用“带批注 diff”做主文档。

这份讲义的写法是：
- 保留清晰的 Markdown 排版
- 只挑真正值得学的代码块
- 解释方式尽量像“我在代码旁边给你写注释”
- 不把 `go.sum` 这种依赖噪音塞满正文

## 1. 先说这节课到底在修什么

如果用一句话概括：

`lesson8 -> lesson10` 这组改动，是把“演示版创建订单逻辑”推进成“像样一点的业务逻辑”。

上一组 `lesson7 -> lesson8` 的 `create_order.go` 有几个明显问题：

1. 业务逻辑全挤在 `Handle()` 里，主流程很乱。
2. 前一个 RPC 的 `err` 会被后一个调用覆盖。
3. 还留着硬编码 `"123"` 这种明显没收口的调试代码。
4. 请求里的重复商品没有合并。
5. 最终落库的商品数据，还是更偏向“相信客户端传来的内容”。

这次改动的核心，就是把这些问题一个个收口。

## 2. 完整差异先看哪里

这次我不把完整 diff 全文直接塞进讲义，因为阅读体验很差，尤其 `go.sum` 会把主线淹没。

你先对照这个原始差异文件：

- [diff.md](/g:/shi/go_shop_second/ai_diffs/lesson_pairs/lesson8-to-lesson10/diff.md)

然后正文里我只展开真正需要讲的 Go 代码。

## 3. 这组代码最正确的阅读顺序

不要上来就看 `go.mod`，也不要按文件名顺序扫。

建议顺序：

1. [create_order.go](/g:/shi/go_shop_second/internal/order/app/command/create_order.go)
2. [service.go](/g:/shi/go_shop_second/internal/order/app/query/service.go)
3. [stock_grpc.go](/g:/shi/go_shop_second/internal/order/adapters/grpc/stock_grpc.go)
4. [stock.proto](/g:/shi/go_shop_second/api/stockpb/stock.proto)

为什么这样看：

- `create_order.go` 是业务主入口，先看它最容易知道“这节到底在干嘛”。
- `service.go` 告诉你抽象接口怎么变了。
- `stock_grpc.go` 告诉你具体实现怎么跟着接口升级。
- `stock.proto` 则是最后补协议层背景，解释为什么这次能拿到 `resp.Items`。

## 4. 先建立一个大脑里的总图

你先把这条链背下来：

`HTTP 请求 -> CreateOrder.Handle() -> validate() -> packItems() -> stockGRPC.CheckIfItemsInStock() -> orderRepo.Create()`

这条链一旦顺了，这一节就不会乱。

它的意思其实很朴素：

1. 用户传来下单请求
2. 订单服务先整理商品数据
3. 再问库存服务这些商品是否可下单
4. 库存服务如果确认没问题，就把商品信息返回
5. 订单服务再把这些“校验过的商品数据”写进订单仓储

## 5. 先看主角：`create_order.go`

这节最重要的文件就是它。

### 5.1 先看最终代码

```go
package command

import (
    "context"
    "errors"

    "github.com/ghost-yu/go_shop_second/common/decorator"
    "github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
    "github.com/ghost-yu/go_shop_second/order/app/query"
    domain "github.com/ghost-yu/go_shop_second/order/domain/order"
    "github.com/sirupsen/logrus"
)

type CreateOrder struct {
    CustomerID string
    Items      []*orderpb.ItemWithQuantity
}

type CreateOrderResult struct {
    OrderID string
}

type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult]

type createOrderHandler struct {
    orderRepo domain.Repository
    stockGRPC query.StockService
}

func NewCreateOrderHandler(
    orderRepo domain.Repository,
    stockGRPC query.StockService,
    logger *logrus.Entry,
    metricClient decorator.MetricsClient,
) CreateOrderHandler {
    if orderRepo == nil {
        panic("nil orderRepo")
    }
    return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
        createOrderHandler{orderRepo: orderRepo, stockGRPC: stockGRPC},
        logger,
        metricClient,
    )
}

func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
    validItems, err := c.validate(ctx, cmd.Items)
    if err != nil {
        return nil, err
    }
    o, err := c.orderRepo.Create(ctx, &domain.Order{
        CustomerID: cmd.CustomerID,
        Items:      validItems,
    })
    if err != nil {
        return nil, err
    }
    return &CreateOrderResult{OrderID: o.ID}, nil
}

func (c createOrderHandler) validate(ctx context.Context, items []*orderpb.ItemWithQuantity) ([]*orderpb.Item, error) {
    if len(items) == 0 {
        return nil, errors.New("must have at least one item")
    }
    items = packItems(items)
    resp, err := c.stockGRPC.CheckIfItemsInStock(ctx, items)
    if err != nil {
        return nil, err
    }
    return resp.Items, nil
}

func packItems(items []*orderpb.ItemWithQuantity) []*orderpb.ItemWithQuantity {
    merged := make(map[string]int32)
    for _, item := range items {
        merged[item.ID] += item.Quantity
    }
    var res []*orderpb.ItemWithQuantity
    for id, quantity := range merged {
        res = append(res, &orderpb.ItemWithQuantity{
            ID:       id,
            Quantity: quantity,
        })
    }
    return res
}
```

### 5.2 伪注释版阅读

下面我不再重复贴源码，而是按“边看边写注释”的方式带你读。

#### `import` 这几行到底在告诉你什么

```go
import (
    "context"
    "errors"
    ...
)
```

这里最该注意的是新加的 `errors`。

你可以把它理解成：

```go
// 标准库 errors 用来创建简单业务错误。
// 这里不需要复杂错误包装，先直接 errors.New("...") 就够了。
```

为什么这行重要？

因为它说明这次改动开始认真处理“非法输入”了，而不是默认请求一定合法。

对于 Go 新手，一个很常见的误区是：

- 只会处理技术错误，比如 RPC 报错、数据库报错
- 忘了处理业务错误，比如“空订单不能创建”

这次就开始补这个洞。

---

#### `Handle()` 现在终于像个主流程了

```go
func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
    validItems, err := c.validate(ctx, cmd.Items)
    if err != nil {
        return nil, err
    }
    o, err := c.orderRepo.Create(ctx, &domain.Order{
        CustomerID: cmd.CustomerID,
        Items:      validItems,
    })
    if err != nil {
        return nil, err
    }
    return &CreateOrderResult{OrderID: o.ID}, nil
}
```

把它翻译成中文注释，其实就是：

```go
// 第一步：先校验商品数据，确认请求不是乱来的。
validItems, err := c.validate(ctx, cmd.Items)
if err != nil {
    // 只要校验不过，立刻返回，不继续往下执行。
    return nil, err
}

// 第二步：校验通过后，再真正创建订单。
o, err := c.orderRepo.Create(ctx, &domain.Order{
    CustomerID: cmd.CustomerID,
    Items:      validItems,
})
if err != nil {
    // 仓储层失败，继续往上抛错误。
    return nil, err
}

// 第三步：创建成功，把订单 ID 返回给上层。
return &CreateOrderResult{OrderID: o.ID}, nil
```

你一定要看懂，这里最大的进步不是“多了个 `validate()` 函数”，而是：

`Handle()` 只剩主流程了。

也就是：
- 先校验
- 再创建
- 最后返回结果

这就是好读的业务代码该有的样子。

### 5.3 为什么这比上一版好很多

上一版的问题是，`Handle()` 里同时混着做了这些事：

1. 调库存 RPC
2. 打日志
3. 处理库存返回
4. 手工组装商品数据
5. 创建订单

这就会导致函数越来越难读。

而这次的思路是：

- 主流程只保留最粗的业务步骤
- 细节拆到 `validate()` 和 `packItems()`

对 Go 小白来说，这比学语法还重要。

因为你以后写业务代码时，最容易犯的不是“不会写 `for range`”，而是“什么都堆到一个函数里”。

## 6. `validate()` 到底做了什么

这是这节最值得你吃透的函数。

```go
func (c createOrderHandler) validate(ctx context.Context, items []*orderpb.ItemWithQuantity) ([]*orderpb.Item, error) {
    if len(items) == 0 {
        return nil, errors.New("must have at least one item")
    }
    items = packItems(items)
    resp, err := c.stockGRPC.CheckIfItemsInStock(ctx, items)
    if err != nil {
        return nil, err
    }
    return resp.Items, nil
}
```

### 6.1 第一层意思：拦截空订单

```go
if len(items) == 0 {
    return nil, errors.New("must have at least one item")
}
```

你可以把它理解成：

```go
// 订单里一件商品都没有，那就不允许创建。
// 这是业务规则，不是技术错误。
```

为什么这里必须先拦？

因为“订单至少有一件商品”是最基础的业务约束。

如果你不在这里拦住，问题就会往后游走：
- 可能跑到库存服务才报错
- 可能跑到仓储层才报错
- 甚至可能静悄悄地创建出一张空订单

这些都不是好结果。

### 6.2 第二层意思：先整理输入，再问库存

```go
items = packItems(items)
```

这行非常关键。

它的意思不是“随便调用一个工具函数”，而是：

```go
// 客户端传来的 items 先不要直接拿去用。
// 先把重复商品合并，再交给库存服务。
```

为什么这样做？

假设客户端传来的是：

```text
itemA x1
itemA x2
itemB x3
```

如果你不合并，库存服务看到的是：
- itemA 两次
- itemB 一次

而库存真正关心的通常是：
- itemA 总共要 3 件
- itemB 总共要 3 件

所以先合并，再校验，会更合理。

### 6.3 第三层意思：现在终于正确处理 RPC 错误了

```go
resp, err := c.stockGRPC.CheckIfItemsInStock(ctx, items)
if err != nil {
    return nil, err
}
```

这里你要重点记住 Go 的一个基本节奏：

```go
result, err := someCall()
if err != nil {
    return nil, err
}
```

这在 Go 里非常常见，因为它清晰、直接、不会把错误拖来拖去。

上一版最大的坏味道之一，就是：
- 先调用一个 RPC
- 不立刻处理 `err`
- 然后又调第二个 RPC
- 把前面的错误覆盖掉

这次终于修正了这个问题。

对 Go 小白来说，这个经验非常重要：

`拿到 err，就尽快处理。`

别等，别攒，别覆盖。

### 6.4 第四层意思：返回的是 `resp.Items`，不是原始 `items`

```go
return resp.Items, nil
```

这行很容易被你一眼带过去，但它其实是这组差异的核心。

它真正表达的意思是：

```go
// 订单最终采用的商品数据，来自库存服务确认后的响应。
// 而不是继续相信外部请求自己传来的数据。
```

这是一次“数据可信来源”的升级。

为什么这很重要？

因为客户端传来的请求只能说是“意图”：
- 我想买什么
- 我想买多少

而库存服务返回的内容才更接近“系统内部确认过的真实数据”：
- 这个商品确实存在
- 这个商品库存足够
- 这个商品的信息是库存服务认可的版本

所以这次改动不是简单的“换个返回值”，而是在提高下单数据的可信度。

## 7. `packItems()` 是这节最适合新手练手的函数

```go
func packItems(items []*orderpb.ItemWithQuantity) []*orderpb.ItemWithQuantity {
    merged := make(map[string]int32)
    for _, item := range items {
        merged[item.ID] += item.Quantity
    }
    var res []*orderpb.ItemWithQuantity
    for id, quantity := range merged {
        res = append(res, &orderpb.ItemWithQuantity{
            ID:       id,
            Quantity: quantity,
        })
    }
    return res
}
```

### 7.1 它到底在干嘛

把这段翻译成中文：

```go
// 先准备一个 map。
// key 是商品 ID。
// value 是这个商品累计后的总数量。
merged := make(map[string]int32)

// 遍历输入切片。
for _, item := range items {
    // 同一个商品出现多次，就把数量加起来。
    merged[item.ID] += item.Quantity
}

// 再把 map 转回切片。
var res []*orderpb.ItemWithQuantity
for id, quantity := range merged {
    res = append(res, &orderpb.ItemWithQuantity{
        ID:       id,
        Quantity: quantity,
    })
}

// 返回整理后的商品列表。
return res
```

### 7.2 这里为什么用 `map`

因为这个问题本质上是“按商品 ID 聚合”。

而 `map` 最适合做这种事情：
- 查找快
- 更新方便
- 写法直接

对于 Go 新手，这段很值得反复看，因为它就是非常典型的“先聚合，再输出”。

### 7.3 最容易忽略的坑

这里有一个非常经典的点：

`Go 的 map 遍历顺序不保证稳定。`

也就是说：

```go
for id, quantity := range merged { ... }
```

每次生成 `res` 时，里面元素的顺序不一定一样。

业务上如果你只关心“商品和数量对不对”，通常没问题。

但如果你后面写测试时去严格比较切片顺序，就可能踩坑。

这是很多 Go 新手第一次用 `map -> slice` 时最容易忽略的地方。

## 8. 接口为什么也必须一起改：`service.go`

现在来看这段：

```go
type StockService interface {
    CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error)
    GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error)
}
```

你可以把它理解成：

```go
// 订单应用层不直接依赖底层 gRPC client。
// 它只依赖一个抽象接口 StockService。
// 这个接口现在升级了：
// CheckIfItemsInStock 不再只告诉我“报没报错”，
// 而是直接把库存服务的响应也给我。
```

为什么这一步必须改？

因为 `create_order.go` 现在要写：

```go
resp, err := c.stockGRPC.CheckIfItemsInStock(...)
```

如果接口还是旧签名：

```go
CheckIfItemsInStock(...) error
```

那编译器根本不会让你这么写。

这恰好就是 Go 接口设计很强的一点：

- 你一旦升级接口
- 所有实现和调用点都会被编译器逼着一起对齐

这反而比动态语言更稳。

## 9. 具体实现怎么跟着接口一起改：`stock_grpc.go`

现在看实现：

```go
func (s StockGRPC) CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error) {
    resp, err := s.client.CheckIfItemsInStock(ctx, &stockpb.CheckIfItemsInStockRequest{Items: items})
    logrus.Info("stock_grpc response", resp)
    return resp, err
}
```

这段你可以按下面方式理解：

```go
// 订单服务把商品列表传进来。
// 这里负责把它包装成 gRPC 请求对象。
resp, err := s.client.CheckIfItemsInStock(ctx, &stockpb.CheckIfItemsInStockRequest{Items: items})

// 先简单打个日志，确认库存服务回了什么。
logrus.Info("stock_grpc response", resp)

// 再把响应和错误一起往上交。
return resp, err
```

### 9.1 它比上一版进步在哪里

上一版是：

```go
resp, err := s.client.CheckIfItemsInStock(...)
logrus.Info(...)
return err
```

也就是：
- 明明拿到了响应
- 却直接把响应丢掉了
- 上层只能知道“有没有错”

而这版终于把 `resp` 往上交了，所以 `create_order.go` 才能真正用 `resp.Items`。

### 9.2 这里涉及的库你要怎么理解

#### `stockpb`

这是 protobuf 生成代码里的包。

意思是：
- `CheckIfItemsInStockRequest` 不是你手写的 struct
- `CheckIfItemsInStockResponse` 也不是你随便定义的 struct
- 它们来自 `.proto` 文件定义

#### `logrus`

这里还是简单打印日志：

```go
logrus.Info("stock_grpc response", resp)
```

对入门来说足够看懂，但你要知道：

更成熟的写法通常会更结构化，比如：
- 把请求参数也打出来
- 把错误字段单独带上
- 用 `WithField` / `WithError` 而不是简单拼接参数

这说明当前代码还在“先把链路打通”的阶段。

## 10. 协议层背景：为什么这次能拿到 `resp.Items`

最后补一下 `stock.proto`，你就彻底不容易糊涂了。

```proto
message CheckIfItemsInStockResponse {
  int32 InStock = 1;
  repeated orderpb.Item Items = 2;
}
```

这段的意思是：

库存服务检查库存时，返回的不只是一个“是否可下单”的状态。
它还会返回 `Items`。

所以 `create_order.go` 现在写：

```go
return resp.Items, nil
```

不是拍脑袋乱写，而是协议层本来就给了这个字段。

这也解释了为什么这次改动是合理的：

- 协议本来支持返回 `Items`
- 只是上一版订单服务没把这个能力用起来
- 这次才把它真正接通

## 11. 这组改动为什么比上一版更像“真正业务代码”

你可以把这次演化总结成下面 4 句话：

1. 主流程更短了
2. 校验职责更清楚了
3. 输入数据先整理再处理了
4. 最终落库的数据来源更可信了

这几个点加起来，才是这节最该学的东西。

不是说它已经完美了，而是它明显开始从“演示代码”往“业务代码”靠近。

## 12. 这组代码还不完美，哪里还可以继续进步

我也直接告诉你这版还不够好的地方，这样你不会误以为它已经是标准答案。

### 12.1 `packItems()` 结果顺序不稳定

刚才说过，`map` 转回切片时顺序不固定。

如果后面业务或测试依赖顺序，就得手动排序。

### 12.2 `validate()` 现在只检查了“不能空订单”

它还没有检查：
- 商品 ID 是否为空
- 数量是否为负数
- 数量是否为 0

这些以后也可能继续补。

### 12.3 `CheckIfItemsInStock` 的返回值现在虽然更有用了，但 `InStock` 字段这版并没有显式使用

也就是说，当前逻辑主要利用了：
- `err`
- `resp.Items`

但 `resp.InStock` 这个字段没有被进一步解释或使用。

### 12.4 日志还是偏调试态

像：

```go
logrus.Info("stock_grpc response", resp)
```

更像联调期间的日志，而不是长期稳定的生产日志格式。

## 13. 你现在最应该带走的 Go 知识点

这组差异最适合你学这几个东西：

1. 早返回错误处理
2. 把大函数拆成主流程 + 小函数
3. `map` 做聚合
4. 接口升级后，调用方和实现方一起修改
5. gRPC/protobuf 返回值不是“只有 error”，还可能带业务数据
6. 客户端传来的数据不一定是最终可信数据

## 14. 如果我是带你复习这一组，我会怎么问你

你可以试着自己回答下面这些问题：

1. 为什么 `Handle()` 不再直接处理所有细节，而是先调用 `validate()`？
2. 为什么 `packItems()` 要放在调用库存服务之前？
3. 为什么最终写进订单仓储的是 `validItems`，不是原始 `cmd.Items`？
4. 为什么接口 `StockService` 必须跟着一起改签名？
5. `map -> slice` 这种写法有什么隐藏坑？

如果这 5 个问题你都能讲明白，这一组你就真的吃透了。

## 15. 最后一句总结

`lesson8 -> lesson10` 的本质，不是“多写了两个函数”，而是订单创建逻辑第一次开始有了像样的分层：

- `Handle()` 负责主流程
- `validate()` 负责业务校验
- `packItems()` 负责输入整理
- `stockGRPC` 负责拿库存服务的真实响应

这就是代码开始变得可维护的信号。
