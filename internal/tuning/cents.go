package tuning

import "math"

// CentsBetween 返回两个频率的音分差，第二个频率相对第一个频率。
func CentsBetween(reference, measured float64) float64 {
	if reference <= 0 || measured <= 0 {
		return 0
	}
	return 1200 * math.Log2(measured/reference)
}

// FrequencyAtCents 根据参考频率和音分偏差反推频率。
func FrequencyAtCents(reference, cents float64) float64 {
	if reference <= 0 {
		return 0
	}
	return reference * math.Pow(2, cents/1200)
}

// MeanCents 计算一组偏差的算术均值。
func MeanCents(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, value := range values {
		sum += value
	}
	return sum / float64(len(values))
}

// WeightedCents 计算加权音分均值。
func WeightedCents(values, weights []float64) float64 {
	if len(values) == 0 || len(values) != len(weights) {
		return 0
	}
	var sum, total float64
	for i, value := range values {
		if weights[i] <= 0 {
			continue
		}
		sum += value * weights[i]
		total += weights[i]
	}
	if total == 0 {
		return 0
	}
	return sum / total
}

// AbsCents 返回绝对音分偏差。
func AbsCents(value float64) float64 {
	if value < 0 {
		return -value
	}
	return value
}
