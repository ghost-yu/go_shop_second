package app

// Application 预留给 stock 服务聚合 Commands 和 Queries。
// 当前 lesson 里 stock 逻辑还比较薄，所以两个分组都是空壳结构。
type Application struct {
	Commands Commands
	Queries  Queries
}

// Commands 未来承载库存写操作，例如扣减库存。
type Commands struct{}

// Queries 未来承载库存读操作，例如按商品 ID 查询库存。
type Queries struct{}
