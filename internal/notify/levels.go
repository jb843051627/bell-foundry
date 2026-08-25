package notify

// EscalationLevel 将业务状态映射为通知级别。
func EscalationLevel(status string) string {
	switch status {
	case "rejected", "recast", "critical":
		return "critical"
	case "review", "needs_retune", "fast":
		return "warn"
	default:
		return "info"
	}
}

// ShouldNotify 判断状态是否值得创建告警。
func ShouldNotify(status string) bool {
	return EscalationLevel(status) != "info"
}

// NormalizeLevel 限制外部传入级别到允许集合。
func NormalizeLevel(level string) string {
	switch level {
	case "critical", "warn", "info":
		return level
	default:
		return "info"
	}
}
