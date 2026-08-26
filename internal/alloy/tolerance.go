package alloy

import "math"

// Tolerance 描述称重偏差的三级判定边界。
type Tolerance struct {
	AcceptPct float64
	ReviewPct float64
}

// DefaultTolerance 是工厂默认的称重容差。
var DefaultTolerance = Tolerance{AcceptPct: 0.8, ReviewPct: 2.5}

// ClassifyDeviation 将最大相对偏差分类为 accepted/review/rejected。
func ClassifyDeviation(deviation float64, t Tolerance) string {
	if deviation < 0 || math.IsNaN(deviation) {
		return "rejected"
	}
	if deviation <= t.AcceptPct {
		return "accepted"
	}
	if deviation <= t.ReviewPct {
		return "review"
	}
	return "rejected"
}

// NormalizeWeights 将组分重量归一化到指定总重，保留两位小数。
func NormalizeWeights(values map[string]float64, target float64) map[string]float64 {
	out := make(map[string]float64, len(values))
	if target <= 0 {
		return out
	}
	total := 0.0
	for _, value := range values {
		if value > 0 {
			total += value
		}
	}
	if total == 0 {
		return out
	}
	for name, value := range values {
		if value <= 0 {
			continue
		}
		out[name] = math.Round(value/total*target*100) / 100
	}
	return out
}

// BalanceError 返回归一化后的组分和与目标值之间的误差。
func BalanceError(values map[string]float64, target float64) float64 {
	sum := 0.0
	for _, value := range values {
		sum += value
	}
	return sum - target
}

// ComponentOrder 返回稳定的组分顺序，便于报告与人工核对。
func ComponentOrder(values map[string]float64) []string {
	preferred := []string{"copper", "tin", "impurity_allowance", "reclaim"}
	seen := make(map[string]bool, len(values))
	out := make([]string, 0, len(values))
	for _, name := range preferred {
		if _, ok := values[name]; ok {
			out = append(out, name)
			seen[name] = true
		}
	}
	for name := range values {
		if !seen[name] {
			out = append(out, name)
		}
	}
	return out
}
