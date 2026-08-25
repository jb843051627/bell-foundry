package cooling

import (
	"sort"

	"github.com/jb843051627/bell-foundry/internal/model"
)

// Segment 是温度曲线中的连续区间。
type Segment struct {
	StartMinute float64  `json:"start_minute"`
	EndMinute   float64  `json:"end_minute"`
	StartTempC  float64  `json:"start_temp_c"`
	EndTempC    float64  `json:"end_temp_c"`
	RateCpm     float64  `json:"rate_cpm"`
	Band        RateBand `json:"band"`
}

// BuildSegments 按相邻采样点构造区间。
func BuildSegments(curve *model.CoolingCurve) []Segment {
	if curve == nil || len(curve.Samples) < 2 {
		return nil
	}
	result := make([]Segment, 0, len(curve.Samples)-1)
	for i := 1; i < len(curve.Samples); i++ {
		before, after := curve.Samples[i-1], curve.Samples[i]
		dt := after.Minute - before.Minute
		if dt <= 0 {
			continue
		}
		rate := (before.TempC - after.TempC) / dt
		result = append(result, Segment{StartMinute: before.Minute, EndMinute: after.Minute, StartTempC: before.TempC, EndTempC: after.TempC, RateCpm: rate, Band: ClassifyRate(rate, 4, 18)})
	}
	return result
}

// FastestSegment 返回降温最快的区间。
func FastestSegment(segments []Segment) (Segment, bool) {
	if len(segments) == 0 {
		return Segment{}, false
	}
	copyOf := append([]Segment(nil), segments...)
	sort.SliceStable(copyOf, func(i, j int) bool { return copyOf[i].RateCpm > copyOf[j].RateCpm })
	return copyOf[0], true
}

// SegmentsInBand 筛选指定工艺级别区间。
func SegmentsInBand(segments []Segment, band RateBand) []Segment {
	result := make([]Segment, 0)
	for _, segment := range segments {
		if segment.Band == band {
			result = append(result, segment)
		}
	}
	return result
}

// RateHistogram 按工艺级别汇总区间数。
func RateHistogram(segments []Segment) map[RateBand]int {
	result := make(map[RateBand]int)
	for _, segment := range segments {
		result[segment.Band]++
	}
	return result
}
