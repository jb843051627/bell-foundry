package model

import "time"

// 告警级别。
const (
	AlertInfo     = "info"
	AlertWarn     = "warn"
	AlertCritical = "critical"
)

// Alert 是质量告警事件。
type Alert struct {
	ID           string    `json:"id"`
	Subject      string    `json:"subject"`
	Level        string    `json:"level"`
	Message      string    `json:"message"`
	CreatedAt    time.Time `json:"created_at"`
	Acknowledged bool      `json:"acknowledged"`
}

// IsUrgent 告警是否需要立即处理（critical 或未确认的 warn）。
func (a *Alert) IsUrgent() bool {
	if a.Level == AlertCritical {
		return true
	}
	return a.Level == AlertWarn && !a.Acknowledged
}
