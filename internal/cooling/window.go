package cooling

import "github.com/jb843051627/bell-foundry/internal/model"

// Window 是带工艺温区的冷却检查器。
type Window struct {
	MinimumC float64
	MaximumC float64
	MaxRate  float64
}

// DefaultWindow 返回钟青铜默认冷却温区。
func DefaultWindow() Window { return Window{MinimumC: 650, MaximumC: 980, MaxRate: 18} }

// Contains 判断温度是否在温区内。
func (w Window) Contains(temperature float64) bool {
	return temperature >= w.MinimumC && temperature <= w.MaximumC
}

// Check 对整条曲线执行温区检查。
func (w Window) Check(curve *model.CoolingCurve) (bool, []string) {
	var findings []string
	if curve == nil || len(curve.Samples) < 2 {
		return false, []string{"not enough samples"}
	}
	for _, sample := range curve.Samples {
		if sample.TempC < w.MinimumC-250 || sample.TempC > w.MaximumC+250 {
			findings = append(findings, "sample outside sensor plausibility range")
			break
		}
	}
	if MaxRate(curve) > w.MaxRate {
		findings = append(findings, "cooling rate too fast")
	}
	return len(findings) == 0, findings
}

// Segment 返回目标温区内的连续采样段。
func (w Window) Segment(curve *model.CoolingCurve) []model.Sample {
	if curve == nil {
		return nil
	}
	out := make([]model.Sample, 0)
	for _, sample := range curve.Samples {
		if w.Contains(sample.TempC) {
			out = append(out, sample)
		}
	}
	return out
}

// Duration 返回曲线在温区内停留的时间。
func (w Window) Duration(curve *model.CoolingCurve) float64 {
	segment := w.Segment(curve)
	if len(segment) < 2 {
		return 0
	}
	return segment[len(segment)-1].Minute - segment[0].Minute
}
