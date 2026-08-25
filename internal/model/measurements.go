package model

import (
	"math"
	"sort"
)

// TemperatureReading 是传感器在某个采样时刻的读数。
type TemperatureReading struct {
	Source   string  `json:"source"`
	Minute   float64 `json:"minute"`
	TempC    float64 `json:"temp_c"`
	Quality  string  `json:"quality"`
	Sequence int64   `json:"sequence"`
}

// MeasurementWindow 表示一个连续采集窗口。
type MeasurementWindow struct {
	ID       string               `json:"id"`
	Source   string               `json:"source"`
	Readings []TemperatureReading `json:"readings"`
	Closed   bool                 `json:"closed"`
}

// Add 追加传感器读数并按分钟排序。
func (w *MeasurementWindow) Add(reading TemperatureReading) error {
	if w == nil || w.Closed {
		return ErrPreconditionFailed
	}
	if reading.Source == "" || reading.Minute < 0 || math.IsNaN(reading.TempC) {
		return ErrInvalidField("reading")
	}
	for _, existing := range w.Readings {
		if existing.Sequence == reading.Sequence && reading.Sequence != 0 {
			return ErrAlreadyExists
		}
	}
	if reading.Quality == "" {
		reading.Quality = "good"
	}
	w.Readings = append(w.Readings, reading)
	sort.SliceStable(w.Readings, func(i, j int) bool { return w.Readings[i].Minute < w.Readings[j].Minute })
	return nil
}

// Close 关闭采集窗口。
func (w *MeasurementWindow) Close() error {
	if w == nil || len(w.Readings) == 0 {
		return ErrCurveIncomplete
	}
	w.Closed = true
	return nil
}

// Average 返回有效读数平均温度。
func (w *MeasurementWindow) Average() float64 {
	if w == nil {
		return 0
	}
	var sum float64
	var count int
	for _, reading := range w.Readings {
		if reading.Quality == "bad" {
			continue
		}
		sum += reading.TempC
		count++
	}
	if count == 0 {
		return 0
	}
	return sum / float64(count)
}

// Range 返回有效读数中的最低和最高温度。
func (w *MeasurementWindow) Range() (float64, float64) {
	if w == nil || len(w.Readings) == 0 {
		return 0, 0
	}
	min, max := math.Inf(1), math.Inf(-1)
	for _, reading := range w.Readings {
		if reading.Quality == "bad" {
			continue
		}
		if reading.TempC < min {
			min = reading.TempC
		}
		if reading.TempC > max {
			max = reading.TempC
		}
	}
	if math.IsInf(min, 1) {
		return 0, 0
	}
	return min, max
}

// Stable 判断连续读数是否在允许抖动范围内。
func (w *MeasurementWindow) Stable(maxDelta float64) bool {
	if w == nil || len(w.Readings) < 2 || maxDelta < 0 {
		return false
	}
	for i := 1; i < len(w.Readings); i++ {
		if math.Abs(w.Readings[i].TempC-w.Readings[i-1].TempC) > maxDelta {
			return false
		}
	}
	return true
}
