package tuning

import "math"

// PartialName 是五分音名称。
var PartialName = []string{"hum", "prime", "tierce", "quint", "nominal"}

// DefaultRatios 是相对名义音的经验比值。
var DefaultRatios = []float64{0.500, 1.000, 1.200, 1.500, 2.000}

// ExpectedFrequencies 根据名义频率得到五分音目标频率。
func ExpectedFrequencies(nominal float64) []float64 {
	out := make([]float64, len(DefaultRatios))
	for i, ratio := range DefaultRatios {
		out[i] = nominal * ratio
	}
	return out
}

// PartialCents 返回各分音相对于目标频率的音分偏差。
func PartialCents(targets, measured []float64) []float64 {
	count := len(targets)
	if len(measured) < count {
		count = len(measured)
	}
	out := make([]float64, count)
	for i := 0; i < count; i++ {
		out[i] = CentsBetween(targets[i], measured[i])
	}
	return out
}

// Complete 判断是否有五个有效频率。
func Complete(measured []float64) bool {
	if len(measured) != len(PartialName) {
		return false
	}
	for _, value := range measured {
		if value <= 0 || math.IsNaN(value) || math.IsInf(value, 0) {
			return false
		}
	}
	return true
}

// RatioReport 返回各分音与名义音的实际比值。
func RatioReport(measured []float64) []float64 {
	if len(measured) < 5 || measured[4] <= 0 {
		return nil
	}
	out := make([]float64, 5)
	for i := range out {
		out[i] = measured[i] / measured[4]
	}
	return out
}
