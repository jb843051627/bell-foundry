package cooling

import (
	"math"

	"github.com/jb843051627/bell-foundry/internal/model"
)

// Forecast 是基于最近斜率的冷却预测。
type Forecast struct {
	TargetC      float64 `json:"target_c"`
	Minutes      float64 `json:"minutes"`
	Confidence   string  `json:"confidence"`
	Extrapolated bool    `json:"extrapolated"`
}

// PredictTo 根据最近两个采样点预测到目标温度的时间。
func PredictTo(curve *model.CoolingCurve, target float64) Forecast {
	result := Forecast{TargetC: target}
	if curve == nil || len(curve.Samples) < 2 {
		return result
	}
	last := curve.Samples[len(curve.Samples)-1]
	previous := curve.Samples[len(curve.Samples)-2]
	dt := last.Minute - previous.Minute
	if dt <= 0 {
		return result
	}
	slope := (last.TempC - previous.TempC) / dt
	if slope >= 0 {
		result.Confidence = "invalid"
		return result
	}
	result.Minutes = last.Minute + math.Max(0, (target-last.TempC)/slope)
	result.Extrapolated = last.TempC > target
	if len(curve.Samples) >= 5 {
		result.Confidence = "high"
	} else {
		result.Confidence = "low"
	}
	return result
}

// CoolingRisk 给出预测风险标签。
func CoolingRisk(f Forecast, maxMinute float64) string {
	if f.Confidence == "invalid" {
		return "sensor-stalled"
	}
	if f.Minutes == 0 {
		return "insufficient-data"
	}
	if maxMinute > 0 && f.Minutes > maxMinute {
		return "slow"
	}
	return "normal"
}
