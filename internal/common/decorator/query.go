package decorator

import (
	"context"

	"github.com/sirupsen/logrus"
)

// QueryHandler defines a generic type that receives a Query Q,
// and returns a result R
// 泛型接口的好处是：同一套日志/指标装饰逻辑可以复用在不同查询上。
type QueryHandler[Q, R any] interface {
	Handle(ctx context.Context, query Q) (R, error)
}

// ApplyQueryDecorators 按固定顺序把日志和指标能力包在真实 handler 外面。
// 这就是装饰器模式：不改业务代码，也能统一加横切能力。
func ApplyQueryDecorators[H, R any](handler QueryHandler[H, R], logger *logrus.Entry, metricsClient MetricsClient) QueryHandler[H, R] {
	return queryLoggingDecorator[H, R]{
		logger: logger,
		base: queryMetricsDecorator[H, R]{
			base:   handler,
			client: metricsClient,
		},
	}
}
