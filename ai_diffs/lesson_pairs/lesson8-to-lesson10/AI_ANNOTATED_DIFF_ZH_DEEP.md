# `lesson8 -> lesson10` 带批注差异（手写版）

说明先说清楚：
- 标准 `diff` 语法本身没有“注释位”。
- 所以下面我采用“保留原始改动行 + 紧跟中文批注”的写法。
- 你可以把这些批注理解成我直接写在差异旁边的 code review 注释。
- 这一组里真正值得学的 Go 代码只集中在 3 个文件：`create_order.go`、`service.go`、`stock_grpc.go`。
- `go.mod`、`go.sum` 我会保留原始结论，但不逐行展开，不然会淹没主线。

## 这组差异先说人话

如果上一组 `lesson7 -> lesson8` 还是“先把订单服务和库存服务连起来”，那这一组就是：

`create_order.go` 开始从“调 RPC 的演示代码”走向“真正有校验、有整理输入、有清晰返回值的业务代码”。

这组改动最重要的不是代码变长了，而是这三件事：

1. 把“校验商品是否合法”从 `Handle()` 主流程里拆到了 `validate()`。
2. 不再只问“库存够不够”，而是直接拿库存服务返回的规范化商品列表。
3. 把重复商品先合并，再送去库存服务校验，避免同一个商品被重复传多次。

---

## 文件 1：`internal/order/adapters/grpc/stock_grpc.go`

```text
diff --git a/internal/order/adapters/grpc/stock_grpc.go b/internal/order/adapters/grpc/stock_grpc.go
index dfe61e6..c2397d6 100644
--- a/internal/order/adapters/grpc/stock_grpc.go
+++ b/internal/order/adapters/grpc/stock_grpc.go
@@ -16,10 +16,10 @@ func NewStockGRPC(client stockpb.StockServiceClient) *StockGRPC {
 	return &StockGRPC{client: client}
 }
 
-func (s StockGRPC) CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) error {
+func (s StockGRPC) CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error) {
+// 讲解：这里是这组差异的第一关键点。
+// 旧版只返回 error，调用方只能知道“库存检查有没有报错”。
+// 新版返回 (*stockpb.CheckIfItemsInStockResponse, error)，说明调用方现在不仅关心“有没有错”，
+// 还关心“库存服务到底返回了什么数据”。
+// 这就是接口设计升级：从只报状态，变成返回结果 + 错误。
+
 	resp, err := s.client.CheckIfItemsInStock(ctx, &stockpb.CheckIfItemsInStockRequest{Items: items})
 	logrus.Info("stock_grpc response", resp)
-	return err
+	return resp, err
+	// 讲解：旧版把 resp 白白丢掉了，新版把它往上层传。
+	// 这样 `create_order.go` 就可以直接用库存服务返回的 `resp.Items`。
+	// 这一步很重要，因为“库存服务返回的数据”会比“客户端自己传上来的数据”更适合作为下单时的可信来源。
 }
```

这里你要真正理解的不是“函数多了一个返回值”，而是职责变化：

- 旧版：库存 RPC 只是一个布尔式校验器，告诉你“行/不行”。
- 新版：库存 RPC 开始兼任“校验 + 返回规范化商品数据”的角色。

为什么这么做：

因为订单服务最终写入仓储时，不应该完全信任客户端自己传来的商品信息。更稳妥的做法是：
- 客户端只负责告诉你“我要买什么、多少件”
- 库存服务负责告诉你“这些商品在我这里是什么样子、是否可下单”
- 订单服务拿库存服务返回的结果入库

容易错的点：

1. `*stockpb.CheckIfItemsInStockResponse` 是 protobuf 生成类型，不是你手写的普通 struct。
2. gRPC 里 `resp != nil` 不代表没有错误，判断成败还是要先看 `err`。
3. 这里仍然保留了 `logrus.Info(...)`，说明当前还是偏调试态，不是最终的结构化日志方案。

---

## 文件 2：`internal/order/app/query/service.go`

```text
diff --git a/internal/order/app/query/service.go b/internal/order/app/query/service.go
index 3f419a9..2e3e4f4 100644
--- a/internal/order/app/query/service.go
+++ b/internal/order/app/query/service.go
@@ -4,9 +4,10 @@ import (
 	"context"
 
 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
+	// 讲解：这里新引入 stockpb，不是为了“多 import 一个包”，
+	// 而是因为接口返回值现在要显式暴露库存服务的响应类型。
 )
 
 type StockService interface {
-	CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) error
+	CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error)
+	// 讲解：接口签名必须跟具体实现一起升级。
+	// 否则 `create_order.go` 虽然想拿 resp.Items，编译器也不会允许。
+	// 这就是 Go 接口设计的好处：一旦接口变了，所有实现和调用点都会被编译器强制检查。
 	GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error)
 }
```

这里虽然只有几行，但它非常关键，因为它体现了 Go 里“面向接口编程”的真正效果。

你要这样理解：

- `create_order.go` 不依赖具体的 gRPC client。
- 它只依赖 `StockService` 这个抽象接口。
- 所以只要接口升级，调用方和实现方都会一起被编译器逼着对齐。

为什么这么做：

因为应用层不应该知道“底层是 gRPC、HTTP、mock 还是别的协议”。应用层只需要知道：“我需要一个库存服务，它能检查库存并返回商品信息。”

Go 小白容易忽略的点：

1. 接口一旦改签名，不是只有实现要改，所有调用方都得改。
2. 这不是坏事，反而是 Go 很强的地方：编译期就把不一致抓出来。
3. `stockpb` 这个包来自 protobuf 生成代码，所以返回类型通常比较长、比较具体，这在微服务项目里很常见。

---

## 文件 3：`internal/order/app/command/create_order.go`

这组差异的主角就是它。

```text
diff --git a/internal/order/app/command/create_order.go b/internal/order/app/command/create_order.go
index adb658f..1c72b04 100644
--- a/internal/order/app/command/create_order.go
+++ b/internal/order/app/command/create_order.go
@@ -2,6 +2,7 @@ package command
 
 import (
 	"context"
+	"errors"
+	// 讲解：这里引入的是 Go 标准库 `errors`，不是第三方库。
+	// 它最基础的用途就是 `errors.New("...")` 直接创建一个简单错误。
+	// 对新手来说，这种错误最容易上手。
 
 	"github.com/ghost-yu/go_shop_second/common/decorator"
 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
@@ -43,23 +44,43 @@ func NewCreateOrderHandler(
 }
 
 func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
-	// TODO: call stock grpc to get items.
-	err := c.stockGRPC.CheckIfItemsInStock(ctx, cmd.Items)
-	resp, err := c.stockGRPC.GetItems(ctx, []string{"123"})
-	logrus.Info("createOrderHandler||resp from stockGRPC.GetItems", resp)
-	var stockResponse []*orderpb.Item
-	for _, item := range cmd.Items {
-		stockResponse = append(stockResponse, &orderpb.Item{
-			ID:       item.ID,
-			Quantity: item.Quantity,
-		})
+	validItems, err := c.validate(ctx, cmd.Items)
+	if err != nil {
+		return nil, err
+	}
+	// 讲解：这一段是整组差异里最关键的改法。
+	// 旧版 `Handle()` 里直接混着做三件事：
+	// 1. 调库存 RPC
+	// 2. 打日志
+	// 3. 自己手工拼订单商品
+	// 结果就是主流程很乱，而且还有硬编码 "123" 这种明显还没收口的调试代码。
+	//
+	// 新版先抽出 `validate()`：
+	// - 如果校验失败，立刻返回
+	// - 如果校验成功，直接拿到 `validItems`
+	// 这样 `Handle()` 主流程就只剩一句话：校验通过后创建订单。
+	// 这就是“把脏活拆走，让主流程变干净”。
+	// 对 Go 新手来说，这种重构思路非常值得学。
+	// 不是所有逻辑都要塞进一个函数里。
+
 	o, err := c.orderRepo.Create(ctx, &domain.Order{
 		CustomerID: cmd.CustomerID,
-		Items:      stockResponse,
+		Items:      validItems,
+		// 讲解：这里也完成了“数据来源切换”。
+		// 旧版写入仓储的是 `stockResponse`，它其实是拿客户端传来的 `cmd.Items` 手工转出来的。
+		// 也就是说，旧版本质上还是在相信请求方自己传来的内容。
+		//
+		// 新版写入的是 `validItems`，它来自 `validate()`，而 `validate()` 又来自库存服务的响应。
+		// 这就更合理：真正落库的商品数据应该尽量来自系统内部可信服务，而不是直接来自外部请求。
 	})
 	if err != nil {
 		return nil, err
 	}
 	return &CreateOrderResult{OrderID: o.ID}, nil
 }
+
+func (c createOrderHandler) validate(ctx context.Context, items []*orderpb.ItemWithQuantity) ([]*orderpb.Item, error) {
+	if len(items) == 0 {
+		return nil, errors.New("must have at least one item")
+	}
+	// 讲解：这一步先挡住“空订单”。
+	// 为什么要先做这个检查？
+	// 因为“订单至少有一件商品”是很基础的业务约束。
+	// 这种约束应该尽量在业务入口就拦掉，不要等到更下游的 RPC 或仓储层再报错。
+	//
+	// `errors.New(...)` 是标准库最简单的错误创建方式，适合这种直白的业务校验。
+
+	items = packItems(items)
+	// 讲解：这里非常值得你注意。
+	// `packItems()` 的作用是把重复商品合并。
+	// 例如客户端传了：
+	// - itemA x1
+	// - itemA x2
+	// 合并后就变成：
+	// - itemA x3
+	//
+	// 为什么要这样做？
+	// 因为库存系统通常关心“每种商品总共要多少件”，
+	// 而不是你请求里把同一个商品拆成了几行。
+	// 先合并再校验，能让后续逻辑更稳定，也避免重复计算。
+
+	resp, err := c.stockGRPC.CheckIfItemsInStock(ctx, items)
+	if err != nil {
+		return nil, err
+	}
+	// 讲解：这里和上一版最大的区别，是终于做到了“拿到 err 立刻处理”。
+	// 上一版最大的坏味道之一，就是前一个 RPC 的错误被后一个调用覆盖了。
+	// 这一版结构清楚很多：
+	// - 调 RPC
+	// - 有错立刻返回
+	// - 没错再继续
+	// 这是 Go 里非常典型、也非常推荐的错误处理节奏。
+
+	return resp.Items, nil
+	// 讲解：这里返回的不是原始 `items`，而是 `resp.Items`。
+	// 说明现在校验函数除了“检查库存”以外，还承担了“从库存服务拿到最终商品信息”的职责。
+	// 这也就是为什么上面接口和 gRPC adapter 都要跟着升级返回值。
+}
+
+func packItems(items []*orderpb.ItemWithQuantity) []*orderpb.ItemWithQuantity {
+	merged := make(map[string]int32)
+	for _, item := range items {
+		merged[item.ID] += item.Quantity
+	}
+	// 讲解：这里用 `map[string]int32` 做聚合。
+	// key 是商品 ID，value 是累加后的数量。
+	// 写法很 Go：简单、直接、可读。
+	//
+	// 如果你是新手，这段一定要看懂：
+	// `merged[item.ID] += item.Quantity`
+	// 的意思就是“把同一个商品的数量不断累加”。
+
+	var res []*orderpb.ItemWithQuantity
+	for id, quantity := range merged {
+		res = append(res, &orderpb.ItemWithQuantity{
+			ID:       id,
+			Quantity: quantity,
+		})
+	}
+	// 讲解：这里再把 map 结果转回切片。
+	// 为什么不直接把 map 往下传？
+	// 因为下游接口 `CheckIfItemsInStock` 要的参数类型就是 `[]*orderpb.ItemWithQuantity`。
+	// 所以这里是在做“内部方便计算的数据结构” -> “外部接口要求的数据结构”的转换。
+
+	return res
+	// 讲解：这里有个你以后很容易踩的点。
+	// Go 的 map 遍历顺序是不保证稳定的，所以 `res` 的顺序可能不是固定的。
+	// 业务上如果你只关心“每种商品和数量”，这没问题。
+	// 但如果你后面写测试去严格比对切片顺序，就可能踩坑。
+}
```

这段差异讲完后，你应该把变化理解成下面这条业务链：

- 旧版：`Handle()` 自己把所有事情搅在一起做，逻辑还带明显调试痕迹。
- 新版：`Handle()` 只负责主流程，`validate()` 负责业务校验，`packItems()` 负责输入整理。

这其实就是很标准的重构思路：

- 主函数只保留主流程
- 把校验拆出去
- 把数据整理拆出去
- 错误立刻返回

对 Go 小白来说，这组差异非常值得反复看，因为它不是“多学一个语法点”，而是在学“怎么把代码写得越来越像真正的业务代码”。

---

## 这组里不重点展开的依赖噪音

### `go.mod`

```text
diff --git a/internal/order/go.mod b/internal/order/go.mod
@@
-	go.uber.org/atomic v1.11.0 // indirect
```

这里不是业务改动，只是依赖树收缩后的结果。你现在不用花太多时间。

### `go.sum`

`go.sum` 这一大段变化本质上也是依赖校准，不是这节课的学习重点。你先知道：
- 它记录依赖校验信息
- 变化通常来自依赖升级、降级、清理
- 不是你现在分析业务时最值得看的地方

---

## 这一组差异里，你最该带走的 6 个点

1. `Handle()` 主流程应该尽量短，只保留主干业务。
2. 校验逻辑可以拆成独立函数，比如 `validate()`。
3. 输入数据进入核心业务前，往往要先整理，比如 `packItems()` 合并重复商品。
4. Go 里拿到 `err` 后要尽快处理，别拖，别覆盖。
5. 接口一旦升级，调用方和实现方会一起被编译器逼着对齐，这是 Go 的优点。
6. `map` 聚合很常见，但要记住遍历顺序不稳定。

## 你现在最适合怎么读这组

建议顺序：

1. 先看 [create_order.go](/g:/shi/go_shop_second/internal/order/app/command/create_order.go)
2. 再看 [service.go](/g:/shi/go_shop_second/internal/order/app/query/service.go)
3. 最后看 [stock_grpc.go](/g:/shi/go_shop_second/internal/order/adapters/grpc/stock_grpc.go)

为什么这样看：
- `create_order.go` 是业务主角
- `service.go` 是抽象契约变化
- `stock_grpc.go` 是具体实现如何跟着接口一起升级

这次我已经按你说的方式改了：解释直接贴在差异里，不再让你来回跳讲义。
