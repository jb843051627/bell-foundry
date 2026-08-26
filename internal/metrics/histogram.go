package metrics

import (
	"sort"
	"sync"
)

// Histogram 维护一组数值样本并提供分位数。
type Histogram struct {
	mu     sync.RWMutex
	values []float64
	limit  int
}

// NewHistogram 创建有限容量直方图。
func NewHistogram(limit int) *Histogram {
	if limit < 1 {
		limit = 256
	}
	return &Histogram{limit: limit}
}

// Observe 记录样本。
func (h *Histogram) Observe(value float64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.values = append(h.values, value)
	if len(h.values) > h.limit {
		h.values = append([]float64(nil), h.values[len(h.values)-h.limit:]...)
	}
}

// Quantile 返回 [0,1] 区间内的经验分位数。
func (h *Histogram) Quantile(q float64) float64 {
	h.mu.RLock()
	values := append([]float64(nil), h.values...)
	h.mu.RUnlock()
	if len(values) == 0 {
		return 0
	}
	if q < 0 {
		q = 0
	}
	if q > 1 {
		q = 1
	}
	sort.Float64s(values)
	index := int(float64(len(values)-1) * q)
	return values[index]
}

// Count 返回样本数量。
func (h *Histogram) Count() int { h.mu.RLock(); defer h.mu.RUnlock(); return len(h.values) }
