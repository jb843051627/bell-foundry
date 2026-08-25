package tuning

import "math"

// CalibrationPoint 是一次调音参考点测量。
type CalibrationPoint struct {
	ReferenceHz float64 `json:"reference_hz"`
	MeasuredHz  float64 `json:"measured_hz"`
	Cents       float64 `json:"cents"`
	Accepted    bool    `json:"accepted"`
}

// BuildCalibration 将参考频率和实测频率转换为校准点。
func BuildCalibration(references, measured []float64, limitCents float64) []CalibrationPoint {
	count := len(references)
	if len(measured) < count {
		count = len(measured)
	}
	out := make([]CalibrationPoint, 0, count)
	for i := 0; i < count; i++ {
		cents := CentsBetween(references[i], measured[i])
		out = append(out, CalibrationPoint{ReferenceHz: references[i], MeasuredHz: measured[i], Cents: cents, Accepted: AbsCents(cents) <= limitCents})
	}
	return out
}

// CalibrationError 返回校准点的均方音分误差。
func CalibrationError(points []CalibrationPoint) float64 {
	if len(points) == 0 {
		return 0
	}
	var sum float64
	for _, p := range points {
		sum += p.Cents * p.Cents
	}
	return math.Sqrt(sum / float64(len(points)))
}

// CalibrationStatus 汇总校准结果。
func CalibrationStatus(points []CalibrationPoint) string {
	if len(points) == 0 {
		return "empty"
	}
	for _, p := range points {
		if !p.Accepted {
			return "review"
		}
	}
	return "accepted"
}

// DriftCents 计算两次相同频率测量之间的漂移。
func DriftCents(first, second []float64) []float64 {
	count := len(first)
	if len(second) < count {
		count = len(second)
	}
	out := make([]float64, count)
	for i := 0; i < count; i++ {
		out[i] = CentsBetween(first[i], second[i])
	}
	return out
}
