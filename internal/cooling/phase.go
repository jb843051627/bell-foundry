package cooling

import "github.com/jb843051627/bell-foundry/internal/model"

// PhaseResult 是冷却相变分析结果。
type PhaseResult struct {
	Minute       float64
	Crossed      bool
	Plateau      bool
	LiquidusSeen bool
	SolidusSeen  bool
}

// DetectPhase 根据固相线穿越和局部斜率平台识别相变时刻。
func DetectPhase(curve *model.CoolingCurve) PhaseResult {
	result := PhaseResult{}
	if curve == nil || len(curve.Samples) == 0 {
		return result
	}
	for _, sample := range curve.Samples {
		if sample.TempC <= curve.LiquidusC {
			result.LiquidusSeen = true
		}
		if !result.Crossed && sample.TempC <= curve.SolidusC {
			result.Crossed = true
			result.SolidusSeen = true
			result.Minute = sample.Minute
		}
	}
	rates := Rates(curve)
	for _, rate := range rates {
		if rate >= 0 && rate < 3.0 {
			result.Plateau = true
			break
		}
	}
	return result
}

// SolidFraction 是温度在液相线到固相线区间内的线性固相比例估计。
func SolidFraction(curve *model.CoolingCurve, temperature float64) float64 {
	if curve == nil || curve.LiquidusC <= curve.SolidusC {
		return 0
	}
	if temperature >= curve.LiquidusC {
		return 0
	}
	if temperature <= curve.SolidusC {
		return 1
	}
	return (curve.LiquidusC - temperature) / (curve.LiquidusC - curve.SolidusC)
}

// IsComplete 判断曲线是否具备相变前后足够采样点。
func IsComplete(curve *model.CoolingCurve) bool {
	if curve == nil || len(curve.Samples) < 4 {
		return false
	}
	phase := DetectPhase(curve)
	return phase.LiquidusSeen && phase.SolidusSeen
}
