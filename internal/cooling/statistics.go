package cooling

import (
	"math"

	"github.com/jb843051627/bell-foundry/internal/model"
)

// Statistics 是冷却曲线的统计摘要。
type Statistics struct {
	Count      int     `json:"count"`
	FirstTempC float64 `json:"first_temp_c"`
	LastTempC  float64 `json:"last_temp_c"`
	MinimumC   float64 `json:"minimum_c"`
	MaximumC   float64 `json:"maximum_c"`
	MeanC      float64 `json:"mean_c"`
	SlopeCpm   float64 `json:"slope_cpm"`
	Variance   float64 `json:"variance"`
}

// Summarize 计算曲线的温度统计。
func Summarize(curve *model.CoolingCurve) Statistics {
	result := Statistics{}
	if curve == nil || len(curve.Samples) == 0 {
		return result
	}
	result.Count = len(curve.Samples)
	result.FirstTempC = curve.Samples[0].TempC
	result.LastTempC = curve.Samples[len(curve.Samples)-1].TempC
	result.MinimumC = curve.Samples[0].TempC
	result.MaximumC = curve.Samples[0].TempC
	for _, sample := range curve.Samples {
		result.MeanC += sample.TempC
		if sample.TempC < result.MinimumC {
			result.MinimumC = sample.TempC
		}
		if sample.TempC > result.MaximumC {
			result.MaximumC = sample.TempC
		}
	}
	result.MeanC /= float64(result.Count)
	for _, sample := range curve.Samples {
		delta := sample.TempC - result.MeanC
		result.Variance += delta * delta
	}
	result.Variance /= float64(result.Count)
	if curve.Span() > 0 {
		result.SlopeCpm = (result.LastTempC - result.FirstTempC) / curve.Span()
	}
	return result
}

// StandardDeviation 返回温度标准差。
func (s Statistics) StandardDeviation() float64 { return math.Sqrt(s.Variance) }

// Uniform 判断采样是否在给定的统计波动内。
func (s Statistics) Uniform(maxStdDev float64) bool {
	return s.Count > 1 && maxStdDev >= 0 && s.StandardDeviation() <= maxStdDev
}

// CoolingDirection 判断曲线整体方向。
func (s Statistics) CoolingDirection() string {
	if s.SlopeCpm < -0.1 {
		return "cooling"
	}
	if s.SlopeCpm > 0.1 {
		return "warming"
	}
	return "stable"
}
