package decorator

import (
	"context"

	"github.com/sirupsen/logrus"
)

// CommandHandler 和 QueryHandler 对称，只是语义上代表“会改状态”的操作。
type CommandHandler[C, R any] interface {
	Handle(ctx context.Context, cmd C) (R, error)
}

// ApplyCommandDecorators 让所有 command 统一拥有日志和指标能力。
// 这样新增一个写操作时，只要实现业务 handler，不用重复写埋点代码。
func ApplyCommandDecorators[C, R any](handler CommandHandler[C, R], logger *logrus.Entry, metricsClient MetricsClient) CommandHandler[C, R] {
	return queryLoggingDecorator[C, R]{
		logger: logger,
		base: queryMetricsDecorator[C, R]{
			base:   handler,
			client: metricsClient,
		},
	}
}
