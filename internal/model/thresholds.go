package model

// ProcessThresholds 集中保存工艺阈值，避免各服务自行散落常量。
type ProcessThresholds struct {
	MinDryHours       float64 `json:"min_dry_hours"`
	MaxMoisturePct    float64 `json:"max_moisture_pct"`
	MinPourFlowSecond int     `json:"min_pour_flow_second"`
	MaxCoolingRate    float64 `json:"max_cooling_rate"`
	TuneLimitCents    float64 `json:"tune_limit_cents"`
	RetuneLimitCents  float64 `json:"retune_limit_cents"`
	MaxOpenDefects    int     `json:"max_open_defects"`
}

// DefaultThresholds 是默认工艺卡参数。
func DefaultThresholds() ProcessThresholds {
	return ProcessThresholds{MinDryHours: 18, MaxMoisturePct: 1.2, MinPourFlowSecond: 12, MaxCoolingRate: 18, TuneLimitCents: 8, RetuneLimitCents: 35, MaxOpenDefects: 0}
}

// Validate 检查阈值是否可用于一条工艺流程。
func (t ProcessThresholds) Validate() bool {
	return t.MinDryHours > 0 && t.MaxMoisturePct > 0 && t.MinPourFlowSecond > 0 && t.MaxCoolingRate > 0 && t.TuneLimitCents > 0 && t.RetuneLimitCents > t.TuneLimitCents && t.MaxOpenDefects >= 0
}

// Merge 用非零字段覆盖默认参数。
func (t ProcessThresholds) Merge(base ProcessThresholds) ProcessThresholds {
	if t.MinDryHours > 0 {
		base.MinDryHours = t.MinDryHours
	}
	if t.MaxMoisturePct > 0 {
		base.MaxMoisturePct = t.MaxMoisturePct
	}
	if t.MinPourFlowSecond > 0 {
		base.MinPourFlowSecond = t.MinPourFlowSecond
	}
	if t.MaxCoolingRate > 0 {
		base.MaxCoolingRate = t.MaxCoolingRate
	}
	if t.TuneLimitCents > 0 {
		base.TuneLimitCents = t.TuneLimitCents
	}
	if t.RetuneLimitCents > 0 {
		base.RetuneLimitCents = t.RetuneLimitCents
	}
	if t.MaxOpenDefects >= 0 {
		base.MaxOpenDefects = t.MaxOpenDefects
	}
	return base
}
