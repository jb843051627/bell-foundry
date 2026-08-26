package tuning

import "sort"

// Interval 是调音结果中一个需要处理的音程区间。
type Interval struct {
	Name        string  `json:"name"`
	LowerHz     float64 `json:"lower_hz"`
	UpperHz     float64 `json:"upper_hz"`
	MeasuredHz  float64 `json:"measured_hz"`
	OffsetCents float64 `json:"offset_cents"`
}

// BuildIntervals 由实测五分音构造相邻区间。
func BuildIntervals(measured []float64) []Interval {
	if len(measured) < 2 {
		return nil
	}
	result := make([]Interval, 0, len(measured)-1)
	for i := 1; i < len(measured); i++ {
		if measured[i-1] <= 0 || measured[i] <= 0 {
			continue
		}
		result = append(result, Interval{Name: PartialName[i-1] + "-" + PartialName[i], LowerHz: measured[i-1], UpperHz: measured[i], MeasuredHz: (measured[i-1] + measured[i]) / 2, OffsetCents: CentsBetween(measured[i-1]*2, measured[i])})
	}
	return result
}

// SortByOffset 按绝对音分偏差从大到小排序。
func SortByOffset(intervals []Interval) []Interval {
	out := append([]Interval(nil), intervals...)
	sort.SliceStable(out, func(i, j int) bool { return AbsCents(out[i].OffsetCents) > AbsCents(out[j].OffsetCents) })
	return out
}

// LargestOffset 返回最大偏差音程。
func LargestOffset(intervals []Interval) (Interval, bool) {
	if len(intervals) == 0 {
		return Interval{}, false
	}
	sorted := SortByOffset(intervals)
	return sorted[0], true
}

// WithinBand 判断全部音程是否在允许带内。
func WithinBand(intervals []Interval, limit float64) bool {
	for _, interval := range intervals {
		if AbsCents(interval.OffsetCents) > limit {
			return false
		}
	}
	return len(intervals) > 0
}
