package decorator

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

// queryLoggingDecorator 在真实 handler 外层补一圈日志。
// “Decorator” 不是框架魔法，本质就是包一层结构体再转调 base。
type queryLoggingDecorator[C, R any] struct {
	logger *logrus.Entry
	base   QueryHandler[C, R]
}

func (q queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	// generateActionName 会把类型名拿出来，作为统一的日志 action 名。
	logger := q.logger.WithFields(logrus.Fields{
		"query":      generateActionName(cmd),
		"query_body": fmt.Sprintf("%#v", cmd),
	})
	logger.Debug("Executing query")
	defer func() {
		if err == nil {
			logger.Info("Query execute successfully")
		} else {
			logger.Error("Failed to execute query", err)
		}
	}()
	return q.base.Handle(ctx, cmd)
}

// generateActionName 直接复用请求类型名，避免每个 handler 手动写一次 action 常量。
func generateActionName(cmd any) string {
	return strings.Split(fmt.Sprintf("%T", cmd), ".")[1]
}
