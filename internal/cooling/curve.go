package cooling

import (
	"sort"

	"github.com/jb843051627/bell-foundry/internal/model"
)

// Append 将采样点加入曲线，拒绝重复时间戳与负时间。
func Append(curve *model.CoolingCurve, sample model.Sample) bool {
	if curve == nil || sample.Minute < 0 {
		return false
	}
	for _, existing := range curve.Samples {
		if existing.Minute == sample.Minute {
			return false
		}
	}
	curve.Samples = append(curve.Samples, sample)
	sort.SliceStable(curve.Samples, func(i, j int) bool {
		return curve.Samples[i].Minute < curve.Samples[j].Minute
	})
	return true
}

// Rates 返回相邻采样点之间的降温速率，单位 ℃/min。
func Rates(curve *model.CoolingCurve) []float64 {
	if curve == nil || len(curve.Samples) < 2 {
		return nil
	}
	rates := make([]float64, 0, len(curve.Samples)-1)
	for i := 1; i < len(curve.Samples); i++ {
		dt := curve.Samples[i].Minute - curve.Samples[i-1].Minute
		if dt <= 0 {
			rates = append(rates, 0)
			continue
		}
		rates = append(rates, (curve.Samples[i-1].TempC-curve.Samples[i].TempC)/dt)
	}
	return rates
}

// MaxRate 返回最大绝对降温速率。
func MaxRate(curve *model.CoolingCurve) float64 {
	max := 0.0
	for _, rate := range Rates(curve) {
		if rate > max {
			max = rate
		}
	}
	return max
}

// FirstBelow 返回首次降到目标温度以下的时间。
func FirstBelow(curve *model.CoolingCurve, temperature float64) (float64, bool) {
	if curve == nil {
		return 0, false
	}
	for _, sample := range curve.Samples {
		if sample.TempC <= temperature {
			return sample.Minute, true
		}
	}
	return 0, false
}
