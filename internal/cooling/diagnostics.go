package cooling

import (
	"math"

	"github.com/jb843051627/bell-foundry/internal/model"
)

// DiagnosticReport 描述冷却曲线的数据质量和工艺风险。
type DiagnosticReport struct {
	SampleCount    int       `json:"sample_count"`
	SpanMinutes    float64   `json:"span_minutes"`
	AverageRate    float64   `json:"average_rate"`
	MaximumRate    float64   `json:"maximum_rate"`
	Monotonic      bool      `json:"monotonic"`
	DuplicateTimes int       `json:"duplicate_times"`
	MissingGaps    []float64 `json:"missing_gaps"`
	Warnings       []string  `json:"warnings"`
}

// Diagnose 生成冷却曲线诊断报告。
func Diagnose(curve *model.CoolingCurve, expectedInterval float64) DiagnosticReport {
	report := DiagnosticReport{}
	if curve == nil {
		report.Warnings = append(report.Warnings, "curve is nil")
		return report
	}
	report.SampleCount = len(curve.Samples)
	report.SpanMinutes = curve.Span()
	report.Monotonic = IsMonotonic(curve)
	if !report.Monotonic {
		report.Warnings = append(report.Warnings, "temperature rises between samples")
	}
	rates := Rates(curve)
	for i, rate := range rates {
		if rate > report.MaximumRate {
			report.MaximumRate = rate
		}
		if i+1 < len(curve.Samples) {
			dt := curve.Samples[i+1].Minute - curve.Samples[i].Minute
			if expectedInterval > 0 && dt > expectedInterval*2 {
				report.MissingGaps = append(report.MissingGaps, dt)
			}
			if dt <= 0 {
				report.DuplicateTimes++
			}
		}
	}
	if report.SpanMinutes > 0 {
		report.AverageRate = (curve.Samples[0].TempC - curve.Samples[len(curve.Samples)-1].TempC) / report.SpanMinutes
	}
	if report.SampleCount < 4 {
		report.Warnings = append(report.Warnings, "too few samples for phase analysis")
	}
	if len(report.MissingGaps) > 0 {
		report.Warnings = append(report.Warnings, "sampling gaps exceed expected interval")
	}
	if report.MaximumRate > 18 {
		report.Warnings = append(report.Warnings, "cooling is too fast")
	}
	return report
}

// IsMonotonic 判断温度是否不升高（允许极小传感器噪声）。
func IsMonotonic(curve *model.CoolingCurve) bool {
	if curve == nil {
		return false
	}
	for i := 1; i < len(curve.Samples); i++ {
		if curve.Samples[i].TempC > curve.Samples[i-1].TempC+0.05 {
			return false
		}
	}
	return true
}

// NoiseRMS 计算相邻温度变化相对于平滑趋势的 RMS 噪声。
func NoiseRMS(curve *model.CoolingCurve) float64 {
	if curve == nil || len(curve.Samples) < 3 {
		return 0
	}
	var sum float64
	var count int
	for i := 1; i+1 < len(curve.Samples); i++ {
		predicted := (curve.Samples[i-1].TempC + curve.Samples[i+1].TempC) / 2
		delta := curve.Samples[i].TempC - predicted
		sum += delta * delta
		count++
	}
	if count == 0 {
		return 0
	}
	return math.Sqrt(sum / float64(count))
}

// ForecastBelow 通过末段平均速率预测降到目标温度的分钟数。
func ForecastBelow(curve *model.CoolingCurve, target float64) float64 {
	if curve == nil || len(curve.Samples) < 2 {
		return 0
	}
	last := curve.Samples[len(curve.Samples)-1]
	if last.TempC <= target {
		return last.Minute
	}
	rates := Rates(curve)
	var sum float64
	var count int
	for i := len(rates) - 1; i >= 0 && count < 3; i-- {
		if rates[i] > 0 {
			sum += rates[i]
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return last.Minute + (last.TempC-target)/(sum/float64(count))
}
