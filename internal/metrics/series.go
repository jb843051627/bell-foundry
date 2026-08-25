package metrics

import (
	"sort"
	"sync"
	"time"
)

// Point 是一个带时间的指标点。
type Point struct {
	At    time.Time `json:"at"`
	Value float64   `json:"value"`
}

// Series 是线程安全的有限长度时间序列。
type Series struct {
	mu     sync.RWMutex
	limit  int
	values []Point
}

// NewSeries 创建时间序列。
func NewSeries(limit int) *Series {
	if limit < 1 {
		limit = 100
	}
	return &Series{limit: limit}
}

// Add 追加指标点并按时间排序。
func (s *Series) Add(point Point) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if point.At.IsZero() {
		point.At = time.Now().UTC()
	}
	s.values = append(s.values, point)
	sort.SliceStable(s.values, func(i, j int) bool { return s.values[i].At.Before(s.values[j].At) })
	if len(s.values) > s.limit {
		s.values = append([]Point(nil), s.values[len(s.values)-s.limit:]...)
	}
}

// Snapshot 返回序列拷贝。
func (s *Series) Snapshot() []Point {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]Point(nil), s.values...)
}

// Latest 返回最新点。
func (s *Series) Latest() (Point, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.values) == 0 {
		return Point{}, false
	}
	return s.values[len(s.values)-1], true
}

// Average 计算窗口平均值。
func (s *Series) Average() float64 {
	values := s.Snapshot()
	if len(values) == 0 {
		return 0
	}
	var sum float64
	for _, value := range values {
		sum += value.Value
	}
	return sum / float64(len(values))
}
