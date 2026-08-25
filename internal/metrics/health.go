package metrics

import (
	"sync/atomic"
	"time"
)

// HealthCollector 聚合服务运行指标。
type HealthCollector struct {
	requests atomic.Uint64
	failures atomic.Uint64
	started  time.Time
	latency  *Series
}

// NewHealthCollector 创建健康指标收集器。
func NewHealthCollector() *HealthCollector {
	return &HealthCollector{started: time.Now().UTC(), latency: NewSeries(256)}
}

// ObserveRequest 记录一次请求。
func (h *HealthCollector) ObserveRequest(duration time.Duration, failed bool) {
	h.requests.Store(h.requests.Load() + 1)
	if failed {
		h.failures.Store(h.failures.Load() + 1)
	}
	h.latency.Add(Point{At: time.Now().UTC(), Value: duration.Seconds()})
}

// Snapshot 返回健康指标。
func (h *HealthCollector) Snapshot() map[string]any {
	requests := h.requests.Load()
	failures := h.failures.Load()
	rate := 0.0
	if requests > 0 {
		rate = float64(failures) / float64(requests)
	}
	return map[string]any{"requests": requests, "failures": failures, "failure_rate": rate, "uptime_seconds": time.Since(h.started).Seconds(), "average_latency_seconds": h.latency.Average()}
}

// Reset 只重置计数，保留启动时间和延迟窗口。
func (h *HealthCollector) Reset() { h.requests = atomic.Uint64{}; h.failures = atomic.Uint64{} }
