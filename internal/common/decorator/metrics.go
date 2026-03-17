package decorator

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// MetricsClient 是一个很小的抽象，方便先接假实现，后续再切到真实监控系统。
type MetricsClient interface {
	Inc(key string, value int)
}

// queryMetricsDecorator 统计耗时、成功、失败次数。
// 它和日志装饰器可以自由组合，因为两者都只依赖同一个 handler 接口。
type queryMetricsDecorator[C, R any] struct {
	base   QueryHandler[C, R]
	client MetricsClient
}

func (q queryMetricsDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	// 用 defer 统一收尾，确保无论成功还是失败都能上报指标。
	start := time.Now()
	actionName := strings.ToLower(generateActionName(cmd))
	defer func() {
		end := time.Since(start)
		q.client.Inc(fmt.Sprintf("querys.%s.duration", actionName), int(end.Seconds()))
		if err == nil {
			q.client.Inc(fmt.Sprintf("querys.%s.success", actionName), 1)
		} else {
			q.client.Inc(fmt.Sprintf("querys.%s.failure", actionName), 1)
		}
	}()
	return q.base.Handle(ctx, cmd)
}
