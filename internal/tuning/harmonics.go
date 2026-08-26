package tuning

import "math"

// HarmonicProfile 描述一口钟的目标泛音比值。
type HarmonicProfile struct {
	Name   string    `json:"name"`
	Ratios []float64 `json:"ratios"`
	Weight []float64 `json:"weight"`
}

// ConcertProfile 是标准音乐会钟的默认配置。
var ConcertProfile = HarmonicProfile{Name: "concert", Ratios: []float64{0.5, 1, 1.2, 1.5, 2}, Weight: []float64{1, 2, 2, 1, 3}}

// RatioDeviation 返回实测相对名义分音的比值偏差百分比。
func RatioDeviation(measured []float64, profile HarmonicProfile) []float64 {
	if len(measured) < len(profile.Ratios) || len(measured) == 0 || measured[len(measured)-1] <= 0 {
		return nil
	}
	nominal := measured[len(measured)-1]
	out := make([]float64, len(profile.Ratios))
	for i, target := range profile.Ratios {
		actual := measured[i] / nominal
		out[i] = (actual - target) / target * 100
	}
	return out
}

// ProfileScore 将比例偏差转换为 0 到 100 的质量分。
func ProfileScore(measured []float64, profile HarmonicProfile) float64 {
	deviations := RatioDeviation(measured, profile)
	if len(deviations) == 0 {
		return 0
	}
	var weighted, total float64
	for i, value := range deviations {
		weight := 1.0
		if i < len(profile.Weight) && profile.Weight[i] > 0 {
			weight = profile.Weight[i]
		}
		weighted += math.Abs(value) * weight
		total += weight
	}
	if total == 0 {
		return 0
	}
	score := 100 - weighted
	if score < 0 {
		return 0
	}
	return score
}

// HarmonicMask 标记哪些分音超过指定比例误差。
func HarmonicMask(measured []float64, profile HarmonicProfile, limitPct float64) []bool {
	deviations := RatioDeviation(measured, profile)
	out := make([]bool, len(deviations))
	for i, value := range deviations {
		out[i] = math.Abs(value) > limitPct
	}
	return out
}

// NearestFrequency 返回最接近目标音高的频率。
func NearestFrequency(reference float64, candidates []float64) (float64, bool) {
	if reference <= 0 || len(candidates) == 0 {
		return 0, false
	}
	best, distance := 0.0, math.Inf(1)
	for _, candidate := range candidates {
		if candidate <= 0 {
			continue
		}
		d := AbsCents(CentsBetween(reference, candidate))
		if d < distance {
			best, distance = candidate, d
		}
	}
	return best, !math.IsInf(distance, 1)
}
