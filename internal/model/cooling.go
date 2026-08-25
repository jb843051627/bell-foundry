package model

// 冷却曲线状态。
const (
	CurveCollecting = "collecting"
	CurveAnalyzed   = "analyzed"
	CurveFinalized  = "finalized"
)

// Sample 是冷却曲线上的一个采样点：入模后第 Minute 分钟的表面温度。
type Sample struct {
	Minute float64 `json:"minute"`
	TempC  float64 `json:"temp_c"`
}

// CoolingCurve 是某次浇注的冷却曲线及分析结论。
type CoolingCurve struct {
	PourID      string   `json:"pour_id"`
	LiquidusC   float64  `json:"liquidus_c"`
	SolidusC    float64  `json:"solidus_c"`
	Samples     []Sample `json:"samples"`
	Status      string   `json:"status"`
	PhaseMinute float64  `json:"phase_minute"` // 相变点（固相线穿越）时刻
	MaxRateCPM  float64  `json:"max_rate_cpm"` // 最大降温速率 ℃/min
	TooFast     bool     `json:"too_fast"`
}

// BelowSolidus 判断最新采样是否已低于固相线。
func (c *CoolingCurve) BelowSolidus() bool {
	if c == nil || len(c.Samples) == 0 {
		return false
	}
	return c.Samples[len(c.Samples)-1].TempC <= c.SolidusC
}

// AddSample 追加采样点并保持按时间升序；重复分钟返回 false。
func (c *CoolingCurve) AddSample(s Sample) bool {
	for _, exist := range c.Samples {
		if exist.Minute == s.Minute {
			return false
		}
	}
	c.Samples = append(c.Samples, s)
	for i := len(c.Samples) - 1; i > 0; i-- {
		if c.Samples[i-1].Minute <= c.Samples[i].Minute {
			break
		}
		c.Samples[i-1], c.Samples[i] = c.Samples[i], c.Samples[i-1]
	}
	return true
}

// Span 覆盖的时间跨度（首末采样点之差）。
func (c *CoolingCurve) Span() float64 {
	if len(c.Samples) < 2 {
		return 0
	}
	last := c.Samples[len(c.Samples)-1]
	first := c.Samples[0]
	return last.Minute - first.Minute
}
