# 中文深度逐行解读（自动生成）

- 仓库: go_shop_second
- 起始引用: lesson10
- 结束引用: lesson11
- 生成时间: 2026-04-06 18:30:24 +08:00
- 解析范围: 仅 *.go，自动排除 *.pb.go 与 *.gen.go
- 解析模式: 祖先链模式（按提交）

## 提交 1: [8ee5f39] stock query and grpc

### 文件: internal/order/adapters/grpc/stock_grpc.go

~~~go
   1: package grpc
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   7: 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
   8: 	"github.com/sirupsen/logrus"
   9: )
  10: 
  11: type StockGRPC struct {
  12: 	client stockpb.StockServiceClient
  13: }
  14: 
  15: func NewStockGRPC(client stockpb.StockServiceClient) *StockGRPC {
  16: 	return &StockGRPC{client: client}
  17: }
  18: 
  19: func (s StockGRPC) CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error) {
  20: 	resp, err := s.client.CheckIfItemsInStock(ctx, &stockpb.CheckIfItemsInStockRequest{Items: items})
  21: 	logrus.Info("stock_grpc response", resp)
  22: 	return resp, err
  23: }
  24: 
  25: func (s StockGRPC) GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error) {
  26: 	resp, err := s.client.GetItems(ctx, &stockpb.GetItemsRequest{ItemIDs: itemIDs})
  27: 	if err != nil {
  28: 		return nil, err
  29: 	}
  30: 	return resp.Items, nil
  31: }
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
| 9 | 语法块结束：关闭 import 或参数列表。 |
| 10 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 11 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 12 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 13 | 代码块结束：收束当前函数、分支或类型定义。 |
| 14 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 15 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 16 | 返回语句：输出当前结果并结束执行路径。 |
| 17 | 代码块结束：收束当前函数、分支或类型定义。 |
| 18 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 19 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 20 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 21 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 22 | 返回语句：输出当前结果并结束执行路径。 |
| 23 | 代码块结束：收束当前函数、分支或类型定义。 |
| 24 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 25 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 26 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 27 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 28 | 返回语句：输出当前结果并结束执行路径。 |
| 29 | 代码块结束：收束当前函数、分支或类型定义。 |
| 30 | 返回语句：输出当前结果并结束执行路径。 |
| 31 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/order/app/command/create_order.go

~~~go
   1: package command
   2: 
   3: import (
   4: 	"context"
   5: 	"errors"
   6: 
   7: 	"github.com/ghost-yu/go_shop_second/common/decorator"
   8: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   9: 	"github.com/ghost-yu/go_shop_second/order/app/query"
  10: 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
  11: 	"github.com/sirupsen/logrus"
  12: )
  13: 
  14: type CreateOrder struct {
  15: 	CustomerID string
  16: 	Items      []*orderpb.ItemWithQuantity
  17: }
  18: 
  19: type CreateOrderResult struct {
  20: 	OrderID string
  21: }
  22: 
  23: type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult]
  24: 
  25: type createOrderHandler struct {
  26: 	orderRepo domain.Repository
  27: 	stockGRPC query.StockService
  28: }
  29: 
  30: func NewCreateOrderHandler(
  31: 	orderRepo domain.Repository,
  32: 	stockGRPC query.StockService,
  33: 	logger *logrus.Entry,
  34: 	metricClient decorator.MetricsClient,
  35: ) CreateOrderHandler {
  36: 	if orderRepo == nil {
  37: 		panic("nil orderRepo")
  38: 	}
  39: 	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
  40: 		createOrderHandler{orderRepo: orderRepo, stockGRPC: stockGRPC},
  41: 		logger,
  42: 		metricClient,
  43: 	)
  44: }
  45: 
  46: func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
  47: 	validItems, err := c.validate(ctx, cmd.Items)
  48: 	if err != nil {
  49: 		return nil, err
  50: 	}
  51: 	o, err := c.orderRepo.Create(ctx, &domain.Order{
  52: 		CustomerID: cmd.CustomerID,
  53: 		Items:      validItems,
  54: 	})
  55: 	if err != nil {
  56: 		return nil, err
  57: 	}
  58: 	return &CreateOrderResult{OrderID: o.ID}, nil
  59: }
  60: 
  61: func (c createOrderHandler) validate(ctx context.Context, items []*orderpb.ItemWithQuantity) ([]*orderpb.Item, error) {
  62: 	if len(items) == 0 {
  63: 		return nil, errors.New("must have at least one item")
  64: 	}
  65: 	items = packItems(items)
  66: 	resp, err := c.stockGRPC.CheckIfItemsInStock(ctx, items)
  67: 	if err != nil {
  68: 		return nil, err
  69: 	}
  70: 	return resp.Items, nil
  71: }
  72: 
  73: func packItems(items []*orderpb.ItemWithQuantity) []*orderpb.ItemWithQuantity {
  74: 	merged := make(map[string]int32)
  75: 	for _, item := range items {
  76: 		merged[item.ID] += item.Quantity
  77: 	}
  78: 	var res []*orderpb.ItemWithQuantity
  79: 	for id, quantity := range merged {
  80: 		res = append(res, &orderpb.ItemWithQuantity{
  81: 			ID:       id,
  82: 			Quantity: quantity,
  83: 		})
  84: 	}
  85: 	return res
  86: }
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
| 10 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 11 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 12 | 语法块结束：关闭 import 或参数列表。 |
| 13 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 14 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 15 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 16 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 17 | 代码块结束：收束当前函数、分支或类型定义。 |
| 18 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 19 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 20 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 21 | 代码块结束：收束当前函数、分支或类型定义。 |
| 22 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 23 | 类型定义：建立语义模型，影响方法与边界设计。 |
| 24 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 25 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 26 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 27 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 28 | 代码块结束：收束当前函数、分支或类型定义。 |
| 29 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 30 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 31 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 32 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 33 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 34 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 35 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 36 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 37 | panic：不可恢复错误路径，常用于关键初始化失败。 |
| 38 | 代码块结束：收束当前函数、分支或类型定义。 |
| 39 | 返回语句：输出当前结果并结束执行路径。 |
| 40 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 41 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 42 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 43 | 语法块结束：关闭 import 或参数列表。 |
| 44 | 代码块结束：收束当前函数、分支或类型定义。 |
| 45 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 46 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 47 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 48 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 49 | 返回语句：输出当前结果并结束执行路径。 |
| 50 | 代码块结束：收束当前函数、分支或类型定义。 |
| 51 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 52 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 53 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 54 | 代码块结束：收束当前函数、分支或类型定义。 |
| 55 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 56 | 返回语句：输出当前结果并结束执行路径。 |
| 57 | 代码块结束：收束当前函数、分支或类型定义。 |
| 58 | 返回语句：输出当前结果并结束执行路径。 |
| 59 | 代码块结束：收束当前函数、分支或类型定义。 |
| 60 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 61 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 62 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 63 | 返回语句：输出当前结果并结束执行路径。 |
| 64 | 代码块结束：收束当前函数、分支或类型定义。 |
| 65 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 66 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 67 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 68 | 返回语句：输出当前结果并结束执行路径。 |
| 69 | 代码块结束：收束当前函数、分支或类型定义。 |
| 70 | 返回语句：输出当前结果并结束执行路径。 |
| 71 | 代码块结束：收束当前函数、分支或类型定义。 |
| 72 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 73 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 74 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 75 | 循环结构：遍历数据或重复执行，需关注终止条件与副作用。 |
| 76 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 77 | 代码块结束：收束当前函数、分支或类型定义。 |
| 78 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 79 | 循环结构：遍历数据或重复执行，需关注终止条件与副作用。 |
| 80 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 81 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 82 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 83 | 代码块结束：收束当前函数、分支或类型定义。 |
| 84 | 代码块结束：收束当前函数、分支或类型定义。 |
| 85 | 返回语句：输出当前结果并结束执行路径。 |
| 86 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/order/app/query/service.go

~~~go
   1: package query
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   7: 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
   8: )
   9: 
  10: type StockService interface {
  11: 	CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error)
  12: 	GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error)
  13: }
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
| 8 | 语法块结束：关闭 import 或参数列表。 |
| 9 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 10 | 接口定义：声明能力契约，用于解耦与可替换实现。 |
| 11 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 12 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 13 | 代码块结束：收束当前函数、分支或类型定义。 |

## 提交 2: [6d9bcf7] docs: 添加 Lesson5-11 关键代码解读（Go 小白版）

该提交没有可解析的非生成 Go 文件变更。

## 提交 3: [e2e4b04] docs: 扩展 Lesson5-11 为文件级复现手册

该提交没有可解析的非生成 Go 文件变更。

## 提交 4: [2d6db2e] docs: gprc.go & http.go 补全完整代码+逐段详细注释

该提交没有可解析的非生成 Go 文件变更。

## 提交 5: [853d6f0] docs: add learning comments for backend internship study

### 文件: internal/common/decorator/command.go

~~~go
   1: package decorator
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/sirupsen/logrus"
   7: )
   8: 
   9: // CommandHandler 和 QueryHandler 对称，只是语义上代表“会改状态”的操作。
  10: type CommandHandler[C, R any] interface {
  11: 	Handle(ctx context.Context, cmd C) (R, error)
  12: }
  13: 
  14: // ApplyCommandDecorators 让所有 command 统一拥有日志和指标能力。
  15: // 这样新增一个写操作时，只要实现业务 handler，不用重复写埋点代码。
  16: func ApplyCommandDecorators[C, R any](handler CommandHandler[C, R], logger *logrus.Entry, metricsClient MetricsClient) CommandHandler[C, R] {
  17: 	return queryLoggingDecorator[C, R]{
  18: 		logger: logger,
  19: 		base: queryMetricsDecorator[C, R]{
  20: 			base:   handler,
  21: 			client: metricsClient,
  22: 		},
  23: 	}
  24: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 语法块结束：关闭 import 或参数列表。 |
| 8 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 9 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 10 | 类型定义：建立语义模型，影响方法与边界设计。 |
| 11 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 12 | 代码块结束：收束当前函数、分支或类型定义。 |
| 13 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 14 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 15 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 16 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 17 | 返回语句：输出当前结果并结束执行路径。 |
| 18 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 19 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 20 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 21 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 22 | 代码块结束：收束当前函数、分支或类型定义。 |
| 23 | 代码块结束：收束当前函数、分支或类型定义。 |
| 24 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/common/decorator/logging.go

~~~go
   1: package decorator
   2: 
   3: import (
   4: 	"context"
   5: 	"fmt"
   6: 	"strings"
   7: 
   8: 	"github.com/sirupsen/logrus"
   9: )
  10: 
  11: // queryLoggingDecorator 在真实 handler 外层补一圈日志。
  12: // “Decorator” 不是框架魔法，本质就是包一层结构体再转调 base。
  13: type queryLoggingDecorator[C, R any] struct {
  14: 	logger *logrus.Entry
  15: 	base   QueryHandler[C, R]
  16: }
  17: 
  18: func (q queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
  19: 	// generateActionName 会把类型名拿出来，作为统一的日志 action 名。
  20: 	logger := q.logger.WithFields(logrus.Fields{
  21: 		"query":      generateActionName(cmd),
  22: 		"query_body": fmt.Sprintf("%#v", cmd),
  23: 	})
  24: 	logger.Debug("Executing query")
  25: 	defer func() {
  26: 		if err == nil {
  27: 			logger.Info("Query execute successfully")
  28: 		} else {
  29: 			logger.Error("Failed to execute query", err)
  30: 		}
  31: 	}()
  32: 	return q.base.Handle(ctx, cmd)
  33: }
  34: 
  35: // generateActionName 直接复用请求类型名，避免每个 handler 手动写一次 action 常量。
  36: func generateActionName(cmd any) string {
  37: 	return strings.Split(fmt.Sprintf("%T", cmd), ".")[1]
  38: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 8 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 9 | 语法块结束：关闭 import 或参数列表。 |
| 10 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 11 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 12 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 13 | 类型定义：建立语义模型，影响方法与边界设计。 |
| 14 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 15 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 16 | 代码块结束：收束当前函数、分支或类型定义。 |
| 17 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 18 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 19 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 20 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 21 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 22 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 23 | 代码块结束：收束当前函数、分支或类型定义。 |
| 24 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 25 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 26 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 27 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 28 | 代码块结束：收束当前函数、分支或类型定义。 |
| 29 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 30 | 代码块结束：收束当前函数、分支或类型定义。 |
| 31 | 代码块结束：收束当前函数、分支或类型定义。 |
| 32 | 返回语句：输出当前结果并结束执行路径。 |
| 33 | 代码块结束：收束当前函数、分支或类型定义。 |
| 34 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 35 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 36 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 37 | 返回语句：输出当前结果并结束执行路径。 |
| 38 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/common/decorator/metrics.go

~~~go
   1: package decorator
   2: 
   3: import (
   4: 	"context"
   5: 	"fmt"
   6: 	"strings"
   7: 	"time"
   8: )
   9: 
  10: // MetricsClient 是一个很小的抽象，方便先接假实现，后续再切到真实监控系统。
  11: type MetricsClient interface {
  12: 	Inc(key string, value int)
  13: }
  14: 
  15: // queryMetricsDecorator 统计耗时、成功、失败次数。
  16: // 它和日志装饰器可以自由组合，因为两者都只依赖同一个 handler 接口。
  17: type queryMetricsDecorator[C, R any] struct {
  18: 	base   QueryHandler[C, R]
  19: 	client MetricsClient
  20: }
  21: 
  22: func (q queryMetricsDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
  23: 	// 用 defer 统一收尾，确保无论成功还是失败都能上报指标。
  24: 	start := time.Now()
  25: 	actionName := strings.ToLower(generateActionName(cmd))
  26: 	defer func() {
  27: 		end := time.Since(start)
  28: 		q.client.Inc(fmt.Sprintf("querys.%s.duration", actionName), int(end.Seconds()))
  29: 		if err == nil {
  30: 			q.client.Inc(fmt.Sprintf("querys.%s.success", actionName), 1)
  31: 		} else {
  32: 			q.client.Inc(fmt.Sprintf("querys.%s.failure", actionName), 1)
  33: 		}
  34: 	}()
  35: 	return q.base.Handle(ctx, cmd)
  36: }
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
| 8 | 语法块结束：关闭 import 或参数列表。 |
| 9 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 10 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 11 | 接口定义：声明能力契约，用于解耦与可替换实现。 |
| 12 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 13 | 代码块结束：收束当前函数、分支或类型定义。 |
| 14 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 15 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 16 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 17 | 类型定义：建立语义模型，影响方法与边界设计。 |
| 18 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 19 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 20 | 代码块结束：收束当前函数、分支或类型定义。 |
| 21 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 22 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 23 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 24 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 25 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 26 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 27 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 28 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 29 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 30 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 31 | 代码块结束：收束当前函数、分支或类型定义。 |
| 32 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 33 | 代码块结束：收束当前函数、分支或类型定义。 |
| 34 | 代码块结束：收束当前函数、分支或类型定义。 |
| 35 | 返回语句：输出当前结果并结束执行路径。 |
| 36 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/common/decorator/query.go

~~~go
   1: package decorator
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/sirupsen/logrus"
   7: )
   8: 
   9: // QueryHandler defines a generic type that receives a Query Q,
  10: // and returns a result R
  11: // 泛型接口的好处是：同一套日志/指标装饰逻辑可以复用在不同查询上。
  12: type QueryHandler[Q, R any] interface {
  13: 	Handle(ctx context.Context, query Q) (R, error)
  14: }
  15: 
  16: // ApplyQueryDecorators 按固定顺序把日志和指标能力包在真实 handler 外面。
  17: // 这就是装饰器模式：不改业务代码，也能统一加横切能力。
  18: func ApplyQueryDecorators[H, R any](handler QueryHandler[H, R], logger *logrus.Entry, metricsClient MetricsClient) QueryHandler[H, R] {
  19: 	return queryLoggingDecorator[H, R]{
  20: 		logger: logger,
  21: 		base: queryMetricsDecorator[H, R]{
  22: 			base:   handler,
  23: 			client: metricsClient,
  24: 		},
  25: 	}
  26: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 语法块结束：关闭 import 或参数列表。 |
| 8 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 9 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 10 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 11 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 12 | 类型定义：建立语义模型，影响方法与边界设计。 |
| 13 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 14 | 代码块结束：收束当前函数、分支或类型定义。 |
| 15 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 16 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 17 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 18 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 19 | 返回语句：输出当前结果并结束执行路径。 |
| 20 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 21 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 22 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 23 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 24 | 代码块结束：收束当前函数、分支或类型定义。 |
| 25 | 代码块结束：收束当前函数、分支或类型定义。 |
| 26 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/common/metrics/todo_metrics.go

~~~go
   1: package metrics
   2: 
   3: // TodoMetrics 是占位实现，用来先打通接口而不真正上报指标。
   4: // 初学时先理解“依赖抽象”比马上接 Prometheus 更重要。
   5: type TodoMetrics struct{}
   6: 
   7: func (t TodoMetrics) Inc(_ string, _ int) {
   8: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 4 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 5 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 6 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 7 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 8 | 代码块结束：收束当前函数、分支或类型定义。 |

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
  10: 	"google.golang.org/grpc"
  11: )
  12: 
  13: // init 会在包加载时执行，这里统一替换 gRPC 默认日志实现，
  14: // 让框架日志和业务日志都走 logrus，便于初学者排查问题时只看一种格式。
  15: func init() {
  16: 	logger := logrus.New()
  17: 	logger.SetLevel(logrus.WarnLevel)
  18: 	grpc_logrus.ReplaceGrpcLogger(logrus.NewEntry(logger))
  19: }
  20: 
  21: // RunGRPCServer 负责按服务名读取配置，再把真正的启动动作委托给 RunGRPCServerOnAddr。
  22: // 这样测试时可以直接传地址，避免每次都依赖 viper 配置。
  23: func RunGRPCServer(serviceName string, registerServer func(server *grpc.Server)) {
  24: 	addr := viper.Sub(serviceName).GetString("grpc-addr")
  25: 	if addr == "" {
  26: 		// TODO: Warning log
  27: 		addr = viper.GetString("fallback-grpc-addr")
  28: 	}
  29: 	RunGRPCServerOnAddr(addr, registerServer)
  30: }
  31: 
  32: // RunGRPCServerOnAddr 创建 gRPC server，并通过 registerServer 回调注册具体业务服务。
  33: // 这类“框架负责启动，业务负责注册”的回调模式，在 Go 后端里非常常见。
  34: func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server)) {
  35: 	logrusEntry := logrus.NewEntry(logrus.StandardLogger())
  36: 	grpcServer := grpc.NewServer(
  37: 		grpc.ChainUnaryInterceptor(
  38: 			// 一元拦截器链类似 HTTP 中间件，用来放日志、指标、鉴权这类横切逻辑。
  39: 			grpc_tags.UnaryServerInterceptor(grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor)),
  40: 			grpc_logrus.UnaryServerInterceptor(logrusEntry),
  41: 			//otelgrpc.UnaryServerInterceptor(),
  42: 			//srvMetrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
  43: 			//logging.UnaryServerInterceptor(interceptorLogger(rpcLogger), logging.WithFieldsFromContext(logTraceID)),
  44: 			//selector.UnaryServerInterceptor(auth.UnaryServerInterceptor(authFn), selector.MatchFunc(allButHealthZ)),
  45: 			//recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
  46: 		),
  47: 		grpc.ChainStreamInterceptor(
  48: 			// 流式 RPC 和一元 RPC 分开配置，是因为两类调用的拦截接口不同。
  49: 			grpc_tags.StreamServerInterceptor(grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor)),
  50: 			grpc_logrus.StreamServerInterceptor(logrusEntry),
  51: 			//otelgrpc.StreamServerInterceptor(),
  52: 			//srvMetrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
  53: 			//logging.StreamServerInterceptor(interceptorLogger(rpcLogger), logging.WithFieldsFromContext(logTraceID)),
  54: 			//selector.StreamServerInterceptor(auth.StreamServerInterceptor(authFn), selector.MatchFunc(allButHealthZ)),
  55: 			//recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
  56: 		),
  57: 	)
  58: 	// 业务方在这里把 OrderService/StockService 的实现挂到 grpcServer 上。
  59: 	registerServer(grpcServer)
  60: 
  61: 	// net.Listen 只负责占住端口；真正开始接收 gRPC 请求要等 Serve 调用后才发生。
  62: 	listen, err := net.Listen("tcp", addr)
  63: 	if err != nil {
  64: 		logrus.Panic(err)
  65: 	}
  66: 	logrus.Infof("Starting gRPC server, Listening: %s", addr)
  67: 	if err := grpcServer.Serve(listen); err != nil {
  68: 		logrus.Panic(err)
  69: 	}
  70: }
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
| 11 | 语法块结束：关闭 import 或参数列表。 |
| 12 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 13 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 14 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 15 | init 函数：包加载时自动执行，常用于初始化与注册。 |
| 16 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 17 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 18 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 19 | 代码块结束：收束当前函数、分支或类型定义。 |
| 20 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 21 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 22 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 23 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 24 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 25 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 26 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 27 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 28 | 代码块结束：收束当前函数、分支或类型定义。 |
| 29 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 30 | 代码块结束：收束当前函数、分支或类型定义。 |
| 31 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 32 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 33 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 34 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 35 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 36 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 37 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 38 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 39 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 40 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 41 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 42 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 43 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 44 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 45 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 46 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 47 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 48 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 49 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 50 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 51 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 52 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 53 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 54 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 55 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 56 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 57 | 语法块结束：关闭 import 或参数列表。 |
| 58 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 59 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 60 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 61 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 62 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 63 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 64 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 65 | 代码块结束：收束当前函数、分支或类型定义。 |
| 66 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 67 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 68 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 69 | 代码块结束：收束当前函数、分支或类型定义。 |
| 70 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/common/server/http.go

~~~go
   1: package server
   2: 
   3: import (
   4: 	"github.com/gin-gonic/gin"
   5: 	"github.com/spf13/viper"
   6: )
   7: 
   8: // RunHTTPServer 按服务名从配置中读取监听地址，再交给 RunHTTPServerOnAddr 真正启动。
   9: func RunHTTPServer(serviceName string, wrapper func(router *gin.Engine)) {
  10: 	addr := viper.Sub(serviceName).GetString("http-addr")
  11: 	if addr == "" {
  12: 		// TODO: Warning log
  13: 	}
  14: 	RunHTTPServerOnAddr(addr, wrapper)
  15: }
  16: 
  17: // RunHTTPServerOnAddr 创建 gin.Engine，并把路由注册动作交给 wrapper 回调。
  18: func RunHTTPServerOnAddr(addr string, wrapper func(router *gin.Engine)) {
  19: 	// gin.New 不带默认中间件，适合教学阶段观察“哪些能力是手动加进去的”。
  20: 	apiRouter := gin.New()
  21: 	// wrapper 负责把 OpenAPI 生成的 handler 绑定到具体路由。
  22: 	wrapper(apiRouter)
  23: 	// 这一行不会修改已有路由，只是创建了一个没有被接住的 RouterGroup。
  24: 	apiRouter.Group("/api")
  25: 	if err := apiRouter.Run(addr); err != nil {
  26: 		panic(err)
  27: 	}
  28: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 6 | 语法块结束：关闭 import 或参数列表。 |
| 7 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 8 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 9 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 10 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 11 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 12 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 13 | 代码块结束：收束当前函数、分支或类型定义。 |
| 14 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 15 | 代码块结束：收束当前函数、分支或类型定义。 |
| 16 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 17 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 18 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 19 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 20 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 21 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 22 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 23 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 24 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 25 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 26 | panic：不可恢复错误路径，常用于关键初始化失败。 |
| 27 | 代码块结束：收束当前函数、分支或类型定义。 |
| 28 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/order/adapters/grpc/stock_grpc.go

~~~go
   1: package grpc
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   7: 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
   8: 	"github.com/sirupsen/logrus"
   9: )
  10: 
  11: // StockGRPC 是 order 侧访问 stock 服务的远程适配器。
  12: // 它实现的是应用层定义的 StockService 接口，而不是把 proto client 直接暴露出去。
  13: type StockGRPC struct {
  14: 	client stockpb.StockServiceClient
  15: }
  16: 
  17: func NewStockGRPC(client stockpb.StockServiceClient) *StockGRPC {
  18: 	return &StockGRPC{client: client}
  19: }
  20: 
  21: func (s StockGRPC) CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error) {
  22: 	// 这里负责把应用层参数翻译成 gRPC 请求对象。
  23: 	resp, err := s.client.CheckIfItemsInStock(ctx, &stockpb.CheckIfItemsInStockRequest{Items: items})
  24: 	logrus.Info("stock_grpc response", resp)
  25: 	return resp, err
  26: }
  27: 
  28: func (s StockGRPC) GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error) {
  29: 	// 对调用方来说这里只是“查商品”，至于底层走 gRPC 还是别的协议都被屏蔽了。
  30: 	resp, err := s.client.GetItems(ctx, &stockpb.GetItemsRequest{ItemIDs: itemIDs})
  31: 	if err != nil {
  32: 		return nil, err
  33: 	}
  34: 	return resp.Items, nil
  35: }
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
| 9 | 语法块结束：关闭 import 或参数列表。 |
| 10 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 11 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 12 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 13 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 14 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 15 | 代码块结束：收束当前函数、分支或类型定义。 |
| 16 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 17 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 18 | 返回语句：输出当前结果并结束执行路径。 |
| 19 | 代码块结束：收束当前函数、分支或类型定义。 |
| 20 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 21 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 22 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 23 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 24 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 25 | 返回语句：输出当前结果并结束执行路径。 |
| 26 | 代码块结束：收束当前函数、分支或类型定义。 |
| 27 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 28 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 29 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 30 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 31 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 32 | 返回语句：输出当前结果并结束执行路径。 |
| 33 | 代码块结束：收束当前函数、分支或类型定义。 |
| 34 | 返回语句：输出当前结果并结束执行路径。 |
| 35 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/order/adapters/order_inmem_repository.go

~~~go
   1: package adapters
   2: 
   3: import (
   4: 	"context"
   5: 	"strconv"
   6: 	"sync"
   7: 	"time"
   8: 
   9: 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
  10: 	"github.com/sirupsen/logrus"
  11: )
  12: 
  13: // MemoryOrderRepository 是 domain.Repository 的内存版实现。
  14: // 教学项目先用它跑通流程，后面换数据库时只需要替换这一层。
  15: type MemoryOrderRepository struct {
  16: 	lock  *sync.RWMutex
  17: 	store []*domain.Order
  18: }
  19: 
  20: func NewMemoryOrderRepository() *MemoryOrderRepository {
  21: 	// 这里放一条假数据，方便一开始就能演示“查询已有订单”的路径。
  22: 	s := make([]*domain.Order, 0)
  23: 	s = append(s, &domain.Order{
  24: 		ID:          "fake-ID",
  25: 		CustomerID:  "fake-customer-id",
  26: 		Status:      "fake-status",
  27: 		PaymentLink: "fake-payment-link",
  28: 		Items:       nil,
  29: 	})
  30: 	return &MemoryOrderRepository{
  31: 		lock:  &sync.RWMutex{},
  32: 		store: s,
  33: 	}
  34: }
  35: 
  36: func (m *MemoryOrderRepository) Create(_ context.Context, order *domain.Order) (*domain.Order, error) {
  37: 	// 写操作要加互斥锁，避免多个请求并发修改切片时产生数据竞争。
  38: 	m.lock.Lock()
  39: 	defer m.lock.Unlock()
  40: 	newOrder := &domain.Order{
  41: 		// 当前用 Unix 时间戳凑一个简单 ID，真实项目里通常会换成雪花算法或 UUID。
  42: 		ID:          strconv.FormatInt(time.Now().Unix(), 10),
  43: 		CustomerID:  order.CustomerID,
  44: 		Status:      order.Status,
  45: 		PaymentLink: order.PaymentLink,
  46: 		Items:       order.Items,
  47: 	}
  48: 	m.store = append(m.store, newOrder)
  49: 	logrus.WithFields(logrus.Fields{
  50: 		"input_order":        order,
  51: 		"store_after_create": m.store,
  52: 	}).Info("memory_order_repo_create")
  53: 	return newOrder, nil
  54: }
  55: 
  56: func (m *MemoryOrderRepository) Get(_ context.Context, id, customerID string) (*domain.Order, error) {
  57: 	for i, v := range m.store {
  58: 		logrus.Infof("m.store[%d] = %+v", i, v)
  59: 	}
  60: 	// 读操作使用读锁，允许多个查询并发进行。
  61: 	m.lock.RLock()
  62: 	defer m.lock.RUnlock()
  63: 	for _, o := range m.store {
  64: 		if o.ID == id && o.CustomerID == customerID {
  65: 			logrus.Infof("memory_order_repo_get||found||id=%s||customerID=%s||res=%+v", id, customerID, *o)
  66: 			return o, nil
  67: 		}
  68: 	}
  69: 	return nil, domain.NotFoundError{OrderID: id}
  70: }
  71: 
  72: func (m *MemoryOrderRepository) Update(ctx context.Context, order *domain.Order, updateFn func(context.Context, *domain.Order) (*domain.Order, error)) error {
  73: 	// UpdateFn 把“怎么改”交给上层，把“在哪里存”留给仓储层，是一种职责分离。
  74: 	m.lock.Lock()
  75: 	defer m.lock.Unlock()
  76: 	found := false
  77: 	for i, o := range m.store {
  78: 		if o.ID == order.ID && o.CustomerID == order.CustomerID {
  79: 			found = true
  80: 			updatedOrder, err := updateFn(ctx, o)
  81: 			if err != nil {
  82: 				return err
  83: 			}
  84: 			m.store[i] = updatedOrder
  85: 		}
  86: 	}
  87: 	if !found {
  88: 		return domain.NotFoundError{OrderID: order.ID}
  89: 	}
  90: 	return nil
  91: }
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
| 9 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 10 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 11 | 语法块结束：关闭 import 或参数列表。 |
| 12 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 13 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 14 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 15 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 16 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 17 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 18 | 代码块结束：收束当前函数、分支或类型定义。 |
| 19 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 20 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 21 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 22 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 23 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 24 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 25 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 26 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 27 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 28 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 29 | 代码块结束：收束当前函数、分支或类型定义。 |
| 30 | 返回语句：输出当前结果并结束执行路径。 |
| 31 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 32 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 33 | 代码块结束：收束当前函数、分支或类型定义。 |
| 34 | 代码块结束：收束当前函数、分支或类型定义。 |
| 35 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 36 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 37 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 38 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 39 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 40 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 41 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 42 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 43 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 44 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 45 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 46 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 47 | 代码块结束：收束当前函数、分支或类型定义。 |
| 48 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 49 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 50 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 51 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 52 | 代码块结束：收束当前函数、分支或类型定义。 |
| 53 | 返回语句：输出当前结果并结束执行路径。 |
| 54 | 代码块结束：收束当前函数、分支或类型定义。 |
| 55 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 56 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 57 | 循环结构：遍历数据或重复执行，需关注终止条件与副作用。 |
| 58 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 59 | 代码块结束：收束当前函数、分支或类型定义。 |
| 60 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 61 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 62 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 63 | 循环结构：遍历数据或重复执行，需关注终止条件与副作用。 |
| 64 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 65 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 66 | 返回语句：输出当前结果并结束执行路径。 |
| 67 | 代码块结束：收束当前函数、分支或类型定义。 |
| 68 | 代码块结束：收束当前函数、分支或类型定义。 |
| 69 | 返回语句：输出当前结果并结束执行路径。 |
| 70 | 代码块结束：收束当前函数、分支或类型定义。 |
| 71 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 72 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 73 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 74 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 75 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 76 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 77 | 循环结构：遍历数据或重复执行，需关注终止条件与副作用。 |
| 78 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 79 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 80 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 81 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 82 | 返回语句：输出当前结果并结束执行路径。 |
| 83 | 代码块结束：收束当前函数、分支或类型定义。 |
| 84 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 85 | 代码块结束：收束当前函数、分支或类型定义。 |
| 86 | 代码块结束：收束当前函数、分支或类型定义。 |
| 87 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 88 | 返回语句：输出当前结果并结束执行路径。 |
| 89 | 代码块结束：收束当前函数、分支或类型定义。 |
| 90 | 返回语句：输出当前结果并结束执行路径。 |
| 91 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/order/app/app.go

~~~go
   1: package app
   2: 
   3: import (
   4: 	"github.com/ghost-yu/go_shop_second/order/app/command"
   5: 	"github.com/ghost-yu/go_shop_second/order/app/query"
   6: )
   7: 
   8: // Application 是应用层门面，让上层只依赖一个总入口，而不用感知每个 handler 的构造细节。
   9: type Application struct {
  10: 	Commands Commands
  11: 	Queries  Queries
  12: }
  13: 
  14: // Commands 聚合“会改状态”的用例，典型如创建、更新订单。
  15: type Commands struct {
  16: 	CreateOrder command.CreateOrderHandler
  17: 	UpdateOrder command.UpdateOrderHandler
  18: }
  19: 
  20: // Queries 聚合“只读”的用例，便于和 Commands 做职责分离。
  21: type Queries struct {
  22: 	GetCustomerOrder query.GetCustomerOrderHandler
  23: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 6 | 语法块结束：关闭 import 或参数列表。 |
| 7 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 8 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 9 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 10 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 11 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 12 | 代码块结束：收束当前函数、分支或类型定义。 |
| 13 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 14 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 15 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 16 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 17 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 18 | 代码块结束：收束当前函数、分支或类型定义。 |
| 19 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 20 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 21 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 22 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 23 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/order/app/command/create_order.go

~~~go
   1: package command
   2: 
   3: import (
   4: 	"context"
   5: 	"errors"
   6: 
   7: 	"github.com/ghost-yu/go_shop_second/common/decorator"
   8: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   9: 	"github.com/ghost-yu/go_shop_second/order/app/query"
  10: 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
  11: 	"github.com/sirupsen/logrus"
  12: )
  13: 
  14: // CreateOrder 是创建订单用例的输入模型。
  15: type CreateOrder struct {
  16: 	CustomerID string
  17: 	Items      []*orderpb.ItemWithQuantity
  18: }
  19: 
  20: // CreateOrderResult 是返回给端口层的结果，避免直接暴露底层仓储对象。
  21: type CreateOrderResult struct {
  22: 	OrderID string
  23: }
  24: 
  25: type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult]
  26: 
  27: // createOrderHandler 只持有完成下单流程所需的两个依赖：订单仓储和库存服务。
  28: type createOrderHandler struct {
  29: 	orderRepo domain.Repository
  30: 	stockGRPC query.StockService
  31: }
  32: 
  33: func NewCreateOrderHandler(
  34: 	orderRepo domain.Repository,
  35: 	stockGRPC query.StockService,
  36: 	logger *logrus.Entry,
  37: 	metricClient decorator.MetricsClient,
  38: ) CreateOrderHandler {
  39: 	if orderRepo == nil {
  40: 		panic("nil orderRepo")
  41: 	}
  42: 	// 下单也是 command，所以同样套上统一的日志和指标装饰器。
  43: 	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
  44: 		createOrderHandler{orderRepo: orderRepo, stockGRPC: stockGRPC},
  45: 		logger,
  46: 		metricClient,
  47: 	)
  48: }
  49: 
  50: func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
  51: 	// 先校验库存，再真正创建订单，避免生成一张无法履约的脏订单。
  52: 	validItems, err := c.validate(ctx, cmd.Items)
  53: 	if err != nil {
  54: 		return nil, err
  55: 	}
  56: 	o, err := c.orderRepo.Create(ctx, &domain.Order{
  57: 		CustomerID: cmd.CustomerID,
  58: 		Items:      validItems,
  59: 	})
  60: 	if err != nil {
  61: 		return nil, err
  62: 	}
  63: 	return &CreateOrderResult{OrderID: o.ID}, nil
  64: }
  65: 
  66: func (c createOrderHandler) validate(ctx context.Context, items []*orderpb.ItemWithQuantity) ([]*orderpb.Item, error) {
  67: 	if len(items) == 0 {
  68: 		return nil, errors.New("must have at least one item")
  69: 	}
  70: 	// packItems 先合并重复商品，避免把同一个商品重复发给库存服务校验。
  71: 	items = packItems(items)
  72: 	resp, err := c.stockGRPC.CheckIfItemsInStock(ctx, items)
  73: 	if err != nil {
  74: 		return nil, err
  75: 	}
  76: 	return resp.Items, nil
  77: }
  78: 
  79: func packItems(items []*orderpb.ItemWithQuantity) []*orderpb.ItemWithQuantity {
  80: 	// merged 以商品 ID 为 key，把数量累加起来。
  81: 	merged := make(map[string]int32)
  82: 	for _, item := range items {
  83: 		merged[item.ID] += item.Quantity
  84: 	}
  85: 	var res []*orderpb.ItemWithQuantity
  86: 	for id, quantity := range merged {
  87: 		res = append(res, &orderpb.ItemWithQuantity{
  88: 			ID:       id,
  89: 			Quantity: quantity,
  90: 		})
  91: 	}
  92: 	return res
  93: }
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
| 10 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 11 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 12 | 语法块结束：关闭 import 或参数列表。 |
| 13 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 14 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 15 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 16 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 17 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 18 | 代码块结束：收束当前函数、分支或类型定义。 |
| 19 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 20 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 21 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 22 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 23 | 代码块结束：收束当前函数、分支或类型定义。 |
| 24 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 25 | 类型定义：建立语义模型，影响方法与边界设计。 |
| 26 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 27 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 28 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 29 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 30 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 31 | 代码块结束：收束当前函数、分支或类型定义。 |
| 32 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 33 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 34 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 35 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 36 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 37 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 38 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 39 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 40 | panic：不可恢复错误路径，常用于关键初始化失败。 |
| 41 | 代码块结束：收束当前函数、分支或类型定义。 |
| 42 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 43 | 返回语句：输出当前结果并结束执行路径。 |
| 44 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 45 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 46 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 47 | 语法块结束：关闭 import 或参数列表。 |
| 48 | 代码块结束：收束当前函数、分支或类型定义。 |
| 49 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 50 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 51 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 52 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 53 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 54 | 返回语句：输出当前结果并结束执行路径。 |
| 55 | 代码块结束：收束当前函数、分支或类型定义。 |
| 56 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 57 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 58 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 59 | 代码块结束：收束当前函数、分支或类型定义。 |
| 60 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 61 | 返回语句：输出当前结果并结束执行路径。 |
| 62 | 代码块结束：收束当前函数、分支或类型定义。 |
| 63 | 返回语句：输出当前结果并结束执行路径。 |
| 64 | 代码块结束：收束当前函数、分支或类型定义。 |
| 65 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 66 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 67 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 68 | 返回语句：输出当前结果并结束执行路径。 |
| 69 | 代码块结束：收束当前函数、分支或类型定义。 |
| 70 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 71 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 72 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 73 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 74 | 返回语句：输出当前结果并结束执行路径。 |
| 75 | 代码块结束：收束当前函数、分支或类型定义。 |
| 76 | 返回语句：输出当前结果并结束执行路径。 |
| 77 | 代码块结束：收束当前函数、分支或类型定义。 |
| 78 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 79 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 80 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 81 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 82 | 循环结构：遍历数据或重复执行，需关注终止条件与副作用。 |
| 83 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 84 | 代码块结束：收束当前函数、分支或类型定义。 |
| 85 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 86 | 循环结构：遍历数据或重复执行，需关注终止条件与副作用。 |
| 87 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 88 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 89 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 90 | 代码块结束：收束当前函数、分支或类型定义。 |
| 91 | 代码块结束：收束当前函数、分支或类型定义。 |
| 92 | 返回语句：输出当前结果并结束执行路径。 |
| 93 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/order/app/command/update_order.go

~~~go
   1: package command
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/decorator"
   7: 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
   8: 	"github.com/sirupsen/logrus"
   9: )
  10: 
  11: // UpdateOrder 把“更新哪个订单”和“如何更新”一起传给 handler。
  12: type UpdateOrder struct {
  13: 	Order    *domain.Order
  14: 	UpdateFn func(context.Context, *domain.Order) (*domain.Order, error)
  15: }
  16: 
  17: type UpdateOrderHandler decorator.CommandHandler[UpdateOrder, interface{}]
  18: 
  19: // updateOrderHandler 负责流程控制，具体修改细节由 UpdateFn 注入。
  20: type updateOrderHandler struct {
  21: 	orderRepo domain.Repository
  22: 	//stockGRPC
  23: }
  24: 
  25: func NewUpdateOrderHandler(
  26: 	orderRepo domain.Repository,
  27: 	logger *logrus.Entry,
  28: 	metricClient decorator.MetricsClient,
  29: ) UpdateOrderHandler {
  30: 	if orderRepo == nil {
  31: 		panic("nil orderRepo")
  32: 	}
  33: 	return decorator.ApplyCommandDecorators[UpdateOrder, interface{}](
  34: 		updateOrderHandler{orderRepo: orderRepo},
  35: 		logger,
  36: 		metricClient,
  37: 	)
  38: }
  39: 
  40: func (c updateOrderHandler) Handle(ctx context.Context, cmd UpdateOrder) (interface{}, error) {
  41: 	// 给 nil UpdateFn 一个 no-op 默认值，避免直接调用时发生空指针问题。
  42: 	if cmd.UpdateFn == nil {
  43: 		logrus.Warnf("updateOrderHandler got nil UpdateFn, order=%#v", cmd.Order)
  44: 		cmd.UpdateFn = func(_ context.Context, order *domain.Order) (*domain.Order, error) { return order, nil }
  45: 	}
  46: 	if err := c.orderRepo.Update(ctx, cmd.Order, cmd.UpdateFn); err != nil {
  47: 		return nil, err
  48: 	}
  49: 	return nil, nil
  50: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 8 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 9 | 语法块结束：关闭 import 或参数列表。 |
| 10 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 11 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 12 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 13 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 14 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 15 | 代码块结束：收束当前函数、分支或类型定义。 |
| 16 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 17 | 类型定义：建立语义模型，影响方法与边界设计。 |
| 18 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 19 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 20 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 21 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 22 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 23 | 代码块结束：收束当前函数、分支或类型定义。 |
| 24 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 25 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 26 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 27 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 28 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 29 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 30 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 31 | panic：不可恢复错误路径，常用于关键初始化失败。 |
| 32 | 代码块结束：收束当前函数、分支或类型定义。 |
| 33 | 返回语句：输出当前结果并结束执行路径。 |
| 34 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 35 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 36 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 37 | 语法块结束：关闭 import 或参数列表。 |
| 38 | 代码块结束：收束当前函数、分支或类型定义。 |
| 39 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 40 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 41 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 42 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 43 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 44 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 45 | 代码块结束：收束当前函数、分支或类型定义。 |
| 46 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 47 | 返回语句：输出当前结果并结束执行路径。 |
| 48 | 代码块结束：收束当前函数、分支或类型定义。 |
| 49 | 返回语句：输出当前结果并结束执行路径。 |
| 50 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/order/app/query/get_customer_order.go

~~~go
   1: package query
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/decorator"
   7: 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
   8: 	"github.com/sirupsen/logrus"
   9: )
  10: 
  11: // GetCustomerOrder 是查询订单的输入对象。
  12: // 把参数收成一个结构体后，后续加字段不会破坏函数签名。
  13: type GetCustomerOrder struct {
  14: 	CustomerID string
  15: 	OrderID    string
  16: }
  17: 
  18: // GetCustomerOrderHandler 是一个已经套好泛型的查询处理器类型别名。
  19: type GetCustomerOrderHandler decorator.QueryHandler[GetCustomerOrder, *domain.Order]
  20: 
  21: // getCustomerOrderHandler 才是真正的业务实现，字段里只放它需要的依赖。
  22: type getCustomerOrderHandler struct {
  23: 	orderRepo domain.Repository
  24: }
  25: 
  26: func NewGetCustomerOrderHandler(
  27: 	orderRepo domain.Repository,
  28: 	logger *logrus.Entry,
  29: 	metricClient decorator.MetricsClient,
  30: ) GetCustomerOrderHandler {
  31: 	if orderRepo == nil {
  32: 		panic("nil orderRepo")
  33: 	}
  34: 	// 构造函数里统一加装饰器，这样调用方拿到的就是“带日志和指标能力”的 handler。
  35: 	return decorator.ApplyQueryDecorators[GetCustomerOrder, *domain.Order](
  36: 		getCustomerOrderHandler{orderRepo: orderRepo},
  37: 		logger,
  38: 		metricClient,
  39: 	)
  40: }
  41: 
  42: func (g getCustomerOrderHandler) Handle(ctx context.Context, query GetCustomerOrder) (*domain.Order, error) {
  43: 	// 查询用例本身很薄，只负责描述流程，把数据访问交给仓储接口。
  44: 	o, err := g.orderRepo.Get(ctx, query.OrderID, query.CustomerID)
  45: 	if err != nil {
  46: 		return nil, err
  47: 	}
  48: 	return o, nil
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
| 7 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 8 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 9 | 语法块结束：关闭 import 或参数列表。 |
| 10 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 11 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 12 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 13 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 14 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 15 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 16 | 代码块结束：收束当前函数、分支或类型定义。 |
| 17 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 18 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 19 | 类型定义：建立语义模型，影响方法与边界设计。 |
| 20 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 21 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 22 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 23 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 24 | 代码块结束：收束当前函数、分支或类型定义。 |
| 25 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 26 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 27 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 28 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 29 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 30 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 31 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 32 | panic：不可恢复错误路径，常用于关键初始化失败。 |
| 33 | 代码块结束：收束当前函数、分支或类型定义。 |
| 34 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 35 | 返回语句：输出当前结果并结束执行路径。 |
| 36 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 37 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 38 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 39 | 语法块结束：关闭 import 或参数列表。 |
| 40 | 代码块结束：收束当前函数、分支或类型定义。 |
| 41 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 42 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 43 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 44 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 45 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 46 | 返回语句：输出当前结果并结束执行路径。 |
| 47 | 代码块结束：收束当前函数、分支或类型定义。 |
| 48 | 返回语句：输出当前结果并结束执行路径。 |
| 49 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/order/app/query/service.go

~~~go
   1: package query
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   7: 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
   8: )
   9: 
  10: // StockService 是 order 应用层眼里的“库存能力端口”。
  11: // 它隔离了底层 gRPC 细节，让 create_order 用例只表达业务意图。
  12: type StockService interface {
  13: 	CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error)
  14: 	GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error)
  15: }
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
| 8 | 语法块结束：关闭 import 或参数列表。 |
| 9 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 10 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 11 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 12 | 接口定义：声明能力契约，用于解耦与可替换实现。 |
| 13 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 14 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 15 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/order/domain/order/order.go

~~~go
   1: package order
   2: 
   3: import "github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   4: 
   5: // Order 是订单领域对象，代表业务里真正被创建、查询、更新的核心实体。
   6: type Order struct {
   7: 	// ID 是订单主键，由仓储层创建时生成。
   8: 	ID string
   9: 	// CustomerID 表示订单属于哪个客户。
  10: 	CustomerID string
  11: 	// Status 预留给支付中、已支付等订单状态流转。
  12: 	Status string
  13: 	// PaymentLink 预留给支付服务返回的支付链接。
  14: 	PaymentLink string
  15: 	// Items 直接复用 proto 里的商品结构，减少服务边界两侧的对象转换成本。
  16: 	Items []*orderpb.Item
  17: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 单行导入：引入一个依赖包，常见于依赖较少的文件。 |
| 4 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 5 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 6 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 7 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 8 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 9 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 10 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 11 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 12 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 13 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 14 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 15 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 16 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 17 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/order/domain/order/repository.go

~~~go
   1: package order
   2: 
   3: import (
   4: 	"context"
   5: 	"fmt"
   6: )
   7: 
   8: // Repository 定义订单持久化能力，应用层只依赖这个接口，而不关心底层是内存还是数据库。
   9: type Repository interface {
  10: 	Create(context.Context, *Order) (*Order, error)
  11: 	Get(ctx context.Context, id, customerID string) (*Order, error)
  12: 	Update(
  13: 		ctx context.Context,
  14: 		o *Order,
  15: 		updateFn func(context.Context, *Order) (*Order, error),
  16: 	) error
  17: }
  18: 
  19: // NotFoundError 是带业务语义的错误类型，调用方可以明确知道“订单不存在”而不是普通异常。
  20: type NotFoundError struct {
  21: 	OrderID string
  22: }
  23: 
  24: func (e NotFoundError) Error() string {
  25: 	return fmt.Sprintf("order '%s' not found", e.OrderID)
  26: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 6 | 语法块结束：关闭 import 或参数列表。 |
| 7 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 8 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 9 | 接口定义：声明能力契约，用于解耦与可替换实现。 |
| 10 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 11 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 12 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 13 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 14 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 15 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 16 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 17 | 代码块结束：收束当前函数、分支或类型定义。 |
| 18 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 19 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 20 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 21 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 22 | 代码块结束：收束当前函数、分支或类型定义。 |
| 23 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 24 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 25 | 返回语句：输出当前结果并结束执行路径。 |
| 26 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/order/http.go

~~~go
   1: package main
   2: 
   3: import (
   4: 	"net/http"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   7: 	"github.com/ghost-yu/go_shop_second/order/app"
   8: 	"github.com/ghost-yu/go_shop_second/order/app/command"
   9: 	"github.com/ghost-yu/go_shop_second/order/app/query"
  10: 	"github.com/gin-gonic/gin"
  11: )
  12: 
  13: // HTTPServer 是 OpenAPI 生成接口在 order 服务里的具体实现。
  14: // 它本身不处理复杂业务，只负责把 HTTP 请求翻译成应用层调用。
  15: type HTTPServer struct {
  16: 	app app.Application
  17: }
  18: 
  19: func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
  20: 	var req orderpb.CreateOrderRequest
  21: 	// ShouldBindJSON 负责把请求体反序列化为 proto 请求对象。
  22: 	if err := c.ShouldBindJSON(&req); err != nil {
  23: 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
  24: 		return
  25: 	}
  26: 	// 进入应用层前，把 HTTP 层对象转换成 command 对象，避免业务逻辑依赖 gin。
  27: 	r, err := H.app.Commands.CreateOrder.Handle(c, command.CreateOrder{
  28: 		CustomerID: req.CustomerID,
  29: 		Items:      req.Items,
  30: 	})
  31: 	if err != nil {
  32: 		c.JSON(http.StatusOK, gin.H{"error": err})
  33: 		return
  34: 	}
  35: 	c.JSON(http.StatusOK, gin.H{
  36: 		"message":     "success",
  37: 		"customer_id": req.CustomerID,
  38: 		"order_id":    r.OrderID,
  39: 	})
  40: }
  41: 
  42: func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
  43: 	// 查询场景走 Queries 分组，体现读写分离的组织方式。
  44: 	o, err := H.app.Queries.GetCustomerOrder.Handle(c, query.GetCustomerOrder{
  45: 		OrderID:    orderID,
  46: 		CustomerID: customerID,
  47: 	})
  48: 	if err != nil {
  49: 		c.JSON(http.StatusOK, gin.H{"error": err})
  50: 		return
  51: 	}
  52: 	c.JSON(http.StatusOK, gin.H{"message": "success", "data": o})
  53: }
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
| 11 | 语法块结束：关闭 import 或参数列表。 |
| 12 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 13 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 14 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 15 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 16 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 17 | 代码块结束：收束当前函数、分支或类型定义。 |
| 18 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 19 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 20 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 21 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 22 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 23 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 24 | 返回语句：输出当前结果并结束执行路径。 |
| 25 | 代码块结束：收束当前函数、分支或类型定义。 |
| 26 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 27 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 28 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 29 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 30 | 代码块结束：收束当前函数、分支或类型定义。 |
| 31 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 32 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 33 | 返回语句：输出当前结果并结束执行路径。 |
| 34 | 代码块结束：收束当前函数、分支或类型定义。 |
| 35 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 36 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 37 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 38 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 39 | 代码块结束：收束当前函数、分支或类型定义。 |
| 40 | 代码块结束：收束当前函数、分支或类型定义。 |
| 41 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 42 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 43 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 44 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 45 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 46 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 47 | 代码块结束：收束当前函数、分支或类型定义。 |
| 48 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 49 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 50 | 返回语句：输出当前结果并结束执行路径。 |
| 51 | 代码块结束：收束当前函数、分支或类型定义。 |
| 52 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 53 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/order/main.go

~~~go
   1: package main
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/config"
   7: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   8: 	"github.com/ghost-yu/go_shop_second/common/server"
   9: 	"github.com/ghost-yu/go_shop_second/order/ports"
  10: 	"github.com/ghost-yu/go_shop_second/order/service"
  11: 	"github.com/gin-gonic/gin"
  12: 	"github.com/sirupsen/logrus"
  13: 	"github.com/spf13/viper"
  14: 	"google.golang.org/grpc"
  15: )
  16: 
  17: // init 先加载全局配置，这样 main 里启动服务时就能直接从 viper 读取地址和服务名。
  18: func init() {
  19: 	if err := config.NewViperConfig(); err != nil {
  20: 		logrus.Fatal(err)
  21: 	}
  22: }
  23: 
  24: func main() {
  25: 	// serviceName 对应 global.yaml 里的 order.service-name，后续 server 包会继续用它找端口配置。
  26: 	serviceName := viper.GetString("order.service-name")
  27: 
  28: 	// ctx 用来把进程级生命周期传给下游依赖，例如 gRPC 客户端连接。
  29: 	ctx, cancel := context.WithCancel(context.Background())
  30: 	defer cancel()
  31: 
  32: 	// NewApplication 是组合根：在这里把 repo、远程 client、handler 全部装配好。
  33: 	application, cleanup := service.NewApplication(ctx)
  34: 	defer cleanup()
  35: 
  36: 	// gRPC 和 HTTP 同时启动，说明 order 服务对内和对外用了两套协议暴露能力。
  37: 	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
  38: 		svc := ports.NewGRPCServer(application)
  39: 		orderpb.RegisterOrderServiceServer(server, svc)
  40: 	})
  41: 
  42: 	// HTTP 侧通过 OpenAPI 生成的 RegisterHandlersWithOptions 完成路由到实现的映射。
  43: 	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
  44: 		ports.RegisterHandlersWithOptions(router, HTTPServer{
  45: 			app: application,
  46: 		}, ports.GinServerOptions{
  47: 			BaseURL:      "/api",
  48: 			Middlewares:  nil,
  49: 			ErrorHandler: nil,
  50: 		})
  51: 	})
  52: 
  53: }
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
| 17 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 18 | init 函数：包加载时自动执行，常用于初始化与注册。 |
| 19 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 20 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 21 | 代码块结束：收束当前函数、分支或类型定义。 |
| 22 | 代码块结束：收束当前函数、分支或类型定义。 |
| 23 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 24 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 25 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 26 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 27 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 28 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 29 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 30 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 31 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 32 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 33 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 34 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 35 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 36 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 37 | goroutine 启动：引入并发执行，需关注生命周期管理。 |
| 38 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 39 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 40 | 代码块结束：收束当前函数、分支或类型定义。 |
| 41 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 42 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 43 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 44 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 45 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 46 | 代码块结束：收束当前函数、分支或类型定义。 |
| 47 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 48 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 49 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 50 | 代码块结束：收束当前函数、分支或类型定义。 |
| 51 | 代码块结束：收束当前函数、分支或类型定义。 |
| 52 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 53 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/order/ports/grpc.go

~~~go
   1: package ports
   2: 
   3: import (
   4: 	context "context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   7: 	"github.com/ghost-yu/go_shop_second/order/app"
   8: 	"google.golang.org/protobuf/types/known/emptypb"
   9: )
  10: 
  11: // GRPCServer 是 orderpb.OrderServiceServer 的端口适配器实现。
  12: // 这一层的职责和 HTTPServer 一样，都是把外部协议请求转给应用层。
  13: type GRPCServer struct {
  14: 	app app.Application
  15: }
  16: 
  17: func NewGRPCServer(app app.Application) *GRPCServer {
  18: 	return &GRPCServer{app: app}
  19: }
  20: 
  21: func (G GRPCServer) CreateOrder(ctx context.Context, request *orderpb.CreateOrderRequest) (*emptypb.Empty, error) {
  22: 	// lesson 当前故意保留 TODO，方便后续逐步实现 gRPC 版本的下单流程。
  23: 	//TODO implement me
  24: 	panic("implement me")
  25: }
  26: 
  27: func (G GRPCServer) GetOrder(ctx context.Context, request *orderpb.GetOrderRequest) (*orderpb.Order, error) {
  28: 	// 这里最终会把 request 转为 query，再返回 proto 层定义的 orderpb.Order。
  29: 	//TODO implement me
  30: 	panic("implement me")
  31: }
  32: 
  33: func (G GRPCServer) UpdateOrder(ctx context.Context, order *orderpb.Order) (*emptypb.Empty, error) {
  34: 	// 更新接口通常会调用 Commands.UpdateOrder；当前阶段先保留占位。
  35: 	//TODO implement me
  36: 	panic("implement me")
  37: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 5 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 8 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 9 | 语法块结束：关闭 import 或参数列表。 |
| 10 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 11 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 12 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 13 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 14 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 15 | 代码块结束：收束当前函数、分支或类型定义。 |
| 16 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 17 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 18 | 返回语句：输出当前结果并结束执行路径。 |
| 19 | 代码块结束：收束当前函数、分支或类型定义。 |
| 20 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 21 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 22 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 23 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 24 | panic：不可恢复错误路径，常用于关键初始化失败。 |
| 25 | 代码块结束：收束当前函数、分支或类型定义。 |
| 26 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 27 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 28 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 29 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 30 | panic：不可恢复错误路径，常用于关键初始化失败。 |
| 31 | 代码块结束：收束当前函数、分支或类型定义。 |
| 32 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 33 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 34 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 35 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 36 | panic：不可恢复错误路径，常用于关键初始化失败。 |
| 37 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/order/service/application.go

~~~go
   1: package service
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	grpcClient "github.com/ghost-yu/go_shop_second/common/client"
   7: 	"github.com/ghost-yu/go_shop_second/common/metrics"
   8: 	"github.com/ghost-yu/go_shop_second/order/adapters"
   9: 	"github.com/ghost-yu/go_shop_second/order/adapters/grpc"
  10: 	"github.com/ghost-yu/go_shop_second/order/app"
  11: 	"github.com/ghost-yu/go_shop_second/order/app/command"
  12: 	"github.com/ghost-yu/go_shop_second/order/app/query"
  13: 	"github.com/sirupsen/logrus"
  14: )
  15: 
  16: // NewApplication 是 order 服务的组合根。
  17: // 它负责创建外部依赖、构造适配器，并返回一个已经装配好的应用层门面。
  18: func NewApplication(ctx context.Context) (app.Application, func()) {
  19: 	stockClient, closeStockClient, err := grpcClient.NewStockGRPCClient(ctx)
  20: 	if err != nil {
  21: 		panic(err)
  22: 	}
  23: 	// 这里把生成的 gRPC client 再包一层适配器，避免应用层直接依赖 proto 细节。
  24: 	stockGRPC := grpc.NewStockGRPC(stockClient)
  25: 	return newApplication(ctx, stockGRPC), func() {
  26: 		_ = closeStockClient()
  27: 	}
  28: }
  29: 
  30: func newApplication(_ context.Context, stockGRPC query.StockService) app.Application {
  31: 	// 组合根里统一决定“接口对应哪种实现”，后面要切数据库或 mock 时只改这里。
  32: 	orderRepo := adapters.NewMemoryOrderRepository()
  33: 	logger := logrus.NewEntry(logrus.StandardLogger())
  34: 	metricClient := metrics.TodoMetrics{}
  35: 	return app.Application{
  36: 		Commands: app.Commands{
  37: 			CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGRPC, logger, metricClient),
  38: 			UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger, metricClient),
  39: 		},
  40: 		Queries: app.Queries{
  41: 			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricClient),
  42: 		},
  43: 	}
  44: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 6 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 7 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 8 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 9 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 10 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 11 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 12 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 13 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 14 | 语法块结束：关闭 import 或参数列表。 |
| 15 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 16 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 17 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 18 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 19 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 20 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 21 | panic：不可恢复错误路径，常用于关键初始化失败。 |
| 22 | 代码块结束：收束当前函数、分支或类型定义。 |
| 23 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 24 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 25 | 返回语句：输出当前结果并结束执行路径。 |
| 26 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 27 | 代码块结束：收束当前函数、分支或类型定义。 |
| 28 | 代码块结束：收束当前函数、分支或类型定义。 |
| 29 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 30 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 31 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 32 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 33 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 34 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 35 | 返回语句：输出当前结果并结束执行路径。 |
| 36 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 37 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 38 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 39 | 代码块结束：收束当前函数、分支或类型定义。 |
| 40 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 41 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 42 | 代码块结束：收束当前函数、分支或类型定义。 |
| 43 | 代码块结束：收束当前函数、分支或类型定义。 |
| 44 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/stock/adapters/stock_inmem_repository.go

~~~go
   1: package adapters
   2: 
   3: import (
   4: 	"context"
   5: 	"sync"
   6: 
   7: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   8: 	domain "github.com/ghost-yu/go_shop_second/stock/domain/stock"
   9: )
  10: 
  11: // MemoryStockRepository 用 map 模拟库存数据源。
  12: // 相比切片，map 更适合按商品 ID 做快速查找。
  13: type MemoryStockRepository struct {
  14: 	lock  *sync.RWMutex
  15: 	store map[string]*orderpb.Item
  16: }
  17: 
  18: // stub 是教学阶段的假数据，用来模拟一个已经存在的库存商品。
  19: var stub = map[string]*orderpb.Item{
  20: 	"item_id": {
  21: 		ID:       "foo_item",
  22: 		Name:     "stub item",
  23: 		Quantity: 10000,
  24: 		PriceID:  "stub_item_price_id",
  25: 	},
  26: }
  27: 
  28: func NewMemoryStockRepository() *MemoryStockRepository {
  29: 	return &MemoryStockRepository{
  30: 		lock:  &sync.RWMutex{},
  31: 		store: stub,
  32: 	}
  33: }
  34: 
  35: func (m MemoryStockRepository) GetItems(ctx context.Context, ids []string) ([]*orderpb.Item, error) {
  36: 	// 这里只有读场景，所以拿读锁即可。
  37: 	m.lock.RLock()
  38: 	defer m.lock.RUnlock()
  39: 	var (
  40: 		res     []*orderpb.Item
  41: 		missing []string
  42: 	)
  43: 	for _, id := range ids {
  44: 		if item, exist := m.store[id]; exist {
  45: 			res = append(res, item)
  46: 		} else {
  47: 			missing = append(missing, id)
  48: 		}
  49: 	}
  50: 	if len(res) == len(ids) {
  51: 		return res, nil
  52: 	}
  53: 	// 返回“部分结果 + 缺失错误”可以帮助上层更明确地决定如何提示用户。
  54: 	return res, domain.NotFoundError{Missing: missing}
  55: }
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
| 8 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 9 | 语法块结束：关闭 import 或参数列表。 |
| 10 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 11 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 12 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 13 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 14 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 15 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 16 | 代码块结束：收束当前函数、分支或类型定义。 |
| 17 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 18 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 19 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 20 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 21 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 22 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 23 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 24 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 25 | 代码块结束：收束当前函数、分支或类型定义。 |
| 26 | 代码块结束：收束当前函数、分支或类型定义。 |
| 27 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 28 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 29 | 返回语句：输出当前结果并结束执行路径。 |
| 30 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 31 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 32 | 代码块结束：收束当前函数、分支或类型定义。 |
| 33 | 代码块结束：收束当前函数、分支或类型定义。 |
| 34 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 35 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 36 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 37 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 38 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 39 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 40 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 41 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 42 | 语法块结束：关闭 import 或参数列表。 |
| 43 | 循环结构：遍历数据或重复执行，需关注终止条件与副作用。 |
| 44 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 45 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 46 | 代码块结束：收束当前函数、分支或类型定义。 |
| 47 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 48 | 代码块结束：收束当前函数、分支或类型定义。 |
| 49 | 代码块结束：收束当前函数、分支或类型定义。 |
| 50 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 51 | 返回语句：输出当前结果并结束执行路径。 |
| 52 | 代码块结束：收束当前函数、分支或类型定义。 |
| 53 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 54 | 返回语句：输出当前结果并结束执行路径。 |
| 55 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/stock/app/app.go

~~~go
   1: package app
   2: 
   3: // Application 预留给 stock 服务聚合 Commands 和 Queries。
   4: // 当前 lesson 里 stock 逻辑还比较薄，所以两个分组都是空壳结构。
   5: type Application struct {
   6: 	Commands Commands
   7: 	Queries  Queries
   8: }
   9: 
  10: // Commands 未来承载库存写操作，例如扣减库存。
  11: type Commands struct{}
  12: 
  13: // Queries 未来承载库存读操作，例如按商品 ID 查询库存。
  14: type Queries struct{}
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 4 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 5 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 6 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 7 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 8 | 代码块结束：收束当前函数、分支或类型定义。 |
| 9 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 10 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 11 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 12 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 13 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 14 | 结构体定义：声明数据载体，承载状态或依赖。 |

### 文件: internal/stock/domain/stock/repository.go

~~~go
   1: package stock
   2: 
   3: import (
   4: 	"context"
   5: 	"fmt"
   6: 	"strings"
   7: 
   8: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   9: )
  10: 
  11: // Repository 定义库存服务需要提供的最小数据访问能力。
  12: // order 服务只要通过应用层依赖这个接口，就不需要知道库存数据放在哪里。
  13: type Repository interface {
  14: 	GetItems(ctx context.Context, ids []string) ([]*orderpb.Item, error)
  15: }
  16: 
  17: // NotFoundError 用来明确告诉调用方哪些商品在库存里不存在。
  18: type NotFoundError struct {
  19: 	Missing []string
  20: }
  21: 
  22: func (e NotFoundError) Error() string {
  23: 	return fmt.Sprintf("these items not found in stock: %s", strings.Join(e.Missing, ","))
  24: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 8 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 9 | 语法块结束：关闭 import 或参数列表。 |
| 10 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 11 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 12 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 13 | 接口定义：声明能力契约，用于解耦与可替换实现。 |
| 14 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 15 | 代码块结束：收束当前函数、分支或类型定义。 |
| 16 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 17 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 18 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 19 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 20 | 代码块结束：收束当前函数、分支或类型定义。 |
| 21 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 22 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 23 | 返回语句：输出当前结果并结束执行路径。 |
| 24 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/stock/main.go

~~~go
   1: package main
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/config"
   7: 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
   8: 	"github.com/ghost-yu/go_shop_second/common/server"
   9: 	"github.com/ghost-yu/go_shop_second/stock/ports"
  10: 	"github.com/ghost-yu/go_shop_second/stock/service"
  11: 	"github.com/sirupsen/logrus"
  12: 	"github.com/spf13/viper"
  13: 	"google.golang.org/grpc"
  14: )
  15: 
  16: // stock 服务和 order 服务一样，在进程入口先加载配置，避免后续各层自己读文件。
  17: func init() {
  18: 	if err := config.NewViperConfig(); err != nil {
  19: 		logrus.Fatal(err)
  20: 	}
  21: }
  22: 
  23: func main() {
  24: 	// serviceName 决定去哪个配置分组下找 grpc-addr 等参数。
  25: 	serviceName := viper.GetString("stock.service-name")
  26: 	// serverType 让同一个二进制可以选择跑哪种协议，当前 lesson 主要走 grpc 分支。
  27: 	serverType := viper.GetString("stock.server-to-run")
  28: 
  29: 	logrus.Info(serverType)
  30: 
  31: 	// 这里先组装应用层，再根据配置决定挂到哪个端口适配器上。
  32: 	ctx, cancel := context.WithCancel(context.Background())
  33: 	defer cancel()
  34: 
  35: 	application := service.NewApplication(ctx)
  36: 	switch serverType {
  37: 	case "grpc":
  38: 		// gRPC 分支把 ports.GRPCServer 注册为 StockService 的服务端实现。
  39: 		server.RunGRPCServer(serviceName, func(server *grpc.Server) {
  40: 			svc := ports.NewGRPCServer(application)
  41: 			stockpb.RegisterStockServiceServer(server, svc)
  42: 		})
  43: 	case "http":
  44: 		// 暂时不用
  45: 	default:
  46: 		panic("unexpected server type")
  47: 	}
  48: }
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
| 14 | 语法块结束：关闭 import 或参数列表。 |
| 15 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 16 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 17 | init 函数：包加载时自动执行，常用于初始化与注册。 |
| 18 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 19 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 20 | 代码块结束：收束当前函数、分支或类型定义。 |
| 21 | 代码块结束：收束当前函数、分支或类型定义。 |
| 22 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 23 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 24 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 25 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 26 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 27 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 28 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 29 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 30 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 31 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 32 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 33 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 34 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 35 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 36 | 多分支选择：按状态或类型分流执行路径。 |
| 37 | 分支标签：定义 switch 的命中条件。 |
| 38 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 39 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 40 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 41 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 42 | 代码块结束：收束当前函数、分支或类型定义。 |
| 43 | 分支标签：定义 switch 的命中条件。 |
| 44 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 45 | 默认分支：兜底处理，防止未知输入静默失败。 |
| 46 | panic：不可恢复错误路径，常用于关键初始化失败。 |
| 47 | 代码块结束：收束当前函数、分支或类型定义。 |
| 48 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/stock/ports/grpc.go

~~~go
   1: package ports
   2: 
   3: import (
   4: 	context "context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
   7: 	"github.com/ghost-yu/go_shop_second/stock/app"
   8: )
   9: 
  10: // GRPCServer 负责实现 stockpb.StockServiceServer 接口。
  11: // 外部请求进入后，应该在这里完成参数转换，再转给应用层处理。
  12: type GRPCServer struct {
  13: 	app app.Application
  14: }
  15: 
  16: func NewGRPCServer(app app.Application) *GRPCServer {
  17: 	return &GRPCServer{app: app}
  18: }
  19: 
  20: func (G GRPCServer) GetItems(ctx context.Context, request *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {
  21: 	// 这里后续会把商品 ID 列表转给库存查询用例。
  22: 	//TODO implement me
  23: 	panic("implement me")
  24: }
  25: 
  26: func (G GRPCServer) CheckIfItemsInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
  27: 	// 这一步对应 order 下单前的库存校验，是跨服务调用的关键入口。
  28: 	//TODO implement me
  29: 	panic("implement me")
  30: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 5 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 8 | 语法块结束：关闭 import 或参数列表。 |
| 9 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 10 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 11 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 12 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 13 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 14 | 代码块结束：收束当前函数、分支或类型定义。 |
| 15 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 16 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 17 | 返回语句：输出当前结果并结束执行路径。 |
| 18 | 代码块结束：收束当前函数、分支或类型定义。 |
| 19 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 20 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 21 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 22 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 23 | panic：不可恢复错误路径，常用于关键初始化失败。 |
| 24 | 代码块结束：收束当前函数、分支或类型定义。 |
| 25 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 26 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 27 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 28 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 29 | panic：不可恢复错误路径，常用于关键初始化失败。 |
| 30 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/stock/service/application.go

~~~go
   1: package service
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/stock/app"
   7: )
   8: 
   9: // NewApplication 当前先返回一个空的应用层门面。
  10: // lesson 后续扩展库存能力时，repo 和 handler 的装配点也会放在这里。
  11: func NewApplication(ctx context.Context) app.Application {
  12: 	return app.Application{}
  13: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 语法块结束：关闭 import 或参数列表。 |
| 8 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 9 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 10 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 11 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 12 | 返回语句：输出当前结果并结束执行路径。 |
| 13 | 代码块结束：收束当前函数、分支或类型定义。 |

## 提交 6: [3163c56] docs: add comments to viper config loader

### 文件: internal/common/config/viper.go

~~~go
   1: package config
   2: 
   3: import "github.com/spf13/viper"
   4: 
   5: // NewViperConfig 统一加载全局配置文件到内存。
   6: // 这个函数在 init() 里被调用，确保 main 启动前配置已经就位，
   7: // 后续各个函数可以直接用 viper.Get... 获取配置，而不用每次都读文件。
   8: func NewViperConfig() error {
   9: 	// SetConfigName 指定配置文件名（不包括扩展名）。
  10: 	// 这里是 "global"，对应 global.yaml。
  11: 	viper.SetConfigName("global")
  12: 
  13: 	// SetConfigType 告诉 viper 配置格式是什么。
  14: 	// viper 支持 json/yaml/toml/hcl 等多种格式。
  15: 	viper.SetConfigType("yaml")
  16: 
  17: 	// AddConfigPath 告诉 viper 去哪个目录找配置文件。
  18: 	// "../common/config" 是相对当前执行目录的路径，
  19: 	// 假设你从项目根目录启动 order 服务，viper 会去 internal/common/config 目录找 global.yaml。
  20: 	viper.AddConfigPath("../common/config")
  21: 
  22: 	// AutomaticEnv 让环境变量自动覆盖配置文件中的值。
  23: 	// 举例：如果环境里设置了 ORDER_GRPC_ADDR=":6666"，会覆盖 yaml 里的 order.grpc-addr。
  24: 	// 这样本地开发、测试、生产环境可以动态调整配置，而不改代码。
  25: 	viper.AutomaticEnv()
  26: 
  27: 	// ReadInConfig 真正去磁盘读配置文件，把内容解析到 viper 的内存结构里。
  28: 	// 如果文件不存在或格式有问题，这里会返回 error。
  29: 	return viper.ReadInConfig()
  30: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 单行导入：引入一个依赖包，常见于依赖较少的文件。 |
| 4 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 5 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 6 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 7 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 8 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 9 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 10 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 11 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 12 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 13 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 14 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 15 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 16 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 17 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 18 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 19 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 20 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 21 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 22 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 23 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 24 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 25 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 26 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 27 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 28 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 29 | 返回语句：输出当前结果并结束执行路径。 |
| 30 | 代码块结束：收束当前函数、分支或类型定义。 |


