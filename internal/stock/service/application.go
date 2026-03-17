package service

import (
	"context"

	"github.com/ghost-yu/go_shop_second/stock/app"
)

// NewApplication 当前先返回一个空的应用层门面。
// lesson 后续扩展库存能力时，repo 和 handler 的装配点也会放在这里。
func NewApplication(ctx context.Context) app.Application {
	return app.Application{}
}
