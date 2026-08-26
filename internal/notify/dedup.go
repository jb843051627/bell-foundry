package notify

import (
	"sync"
	"time"
)

// Deduper 在 TTL 内只允许同一 subject+message 发送一次。
type Deduper struct {
	mu    sync.Mutex
	items map[string]time.Time
	ttl   time.Duration
}

// NewDeduper 创建告警去重器。
func NewDeduper(ttl time.Duration) *Deduper {
	if ttl <= 0 {
		ttl = 10 * time.Minute
	}
	return &Deduper{items: make(map[string]time.Time), ttl: ttl}
}

// Allow 检查并登记一个键，返回是否允许发送。
// 多路传感器或并发的 HTTP 请求可能同时为同一告警窗口调用本方法，
// 因此整个"清理过期项 → 查重 → 登记"过程必须在锁内原子完成，
// 否则两个 goroutine 会同时通过查重检查并都返回 true，导致同一窗口重复推送。
func (d *Deduper) Allow(key string, now time.Time) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	for old, at := range d.items {
		if now.Sub(at) >= d.ttl {
			delete(d.items, old)
		}
	}
	if at, ok := d.items[key]; ok && now.Sub(at) < d.ttl {
		return false
	}
	d.items[key] = now
	return true
}

// Size 返回当前未过期键数量。
func (d *Deduper) Size() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.items)
}
