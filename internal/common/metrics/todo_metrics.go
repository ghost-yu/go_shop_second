package metrics

// TodoMetrics 是占位实现，用来先打通接口而不真正上报指标。
// 初学时先理解“依赖抽象”比马上接 Prometheus 更重要。
type TodoMetrics struct{}

func (t TodoMetrics) Inc(_ string, _ int) {
}
