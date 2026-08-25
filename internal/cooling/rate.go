package cooling

import (
	"math"

	"github.com/jb843051627/bell-foundry/internal/model"
)

// RateBand 是冷却速率的工艺分级。
type RateBand string

const (
	RateGentle  RateBand = "gentle"
	RateNominal RateBand = "nominal"
	RateFast    RateBand = "fast"
	RateUnknown RateBand = "unknown"
)

// ClassifyRate 根据最大速率分类；阈值来自工艺卡。
func ClassifyRate(rate, gentleLimit, fastLimit float64) RateBand {
	if rate < 0 || math.IsNaN(rate) {
		return RateUnknown
	}
	if rate <= gentleLimit {
		return RateGentle
	}
	if rate >= fastLimit {
		return RateFast
	}
	return RateNominal
}

// EvaluateRate 返回曲线速率、分级及是否需要告警。
func EvaluateRate(curve *model.CoolingCurve, gentleLimit, fastLimit float64) (float64, RateBand, bool) {
	rate := MaxRate(curve)
	band := ClassifyRate(rate, gentleLimit, fastLimit)
	return rate, band, band == RateFast || band == RateUnknown
}

// Interpolate 在相邻采样点之间线性插值温度。
func Interpolate(a, b model.Sample, minute float64) (float64, bool) {
	if b.Minute <= a.Minute || minute < a.Minute || minute > b.Minute {
		return 0, false
	}
	ratio := (minute - a.Minute) / (b.Minute - a.Minute)
	return a.TempC + (b.TempC-a.TempC)*ratio, true
}

// CoolingWindow 返回曲线从起始温度降到目标温度所需的估计分钟数。
func CoolingWindow(curve *model.CoolingCurve, target float64) float64 {
	if curve == nil || len(curve.Samples) < 2 {
		return 0
	}
	if minute, ok := FirstBelow(curve, target); ok {
		return minute
	}
	return curve.Span()
}
