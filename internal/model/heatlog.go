package model

import (
	"sort"
	"time"
)

// HeatLog 记录炉次中的连续温度和功率观测。
type HeatLog struct {
	HeatID   string           `json:"heat_id"`
	Readings []HeatLogReading `json:"readings"`
	Closed   bool             `json:"closed"`
}

// HeatLogReading 是炉次中的一条观测。
type HeatLogReading struct {
	At          time.Time `json:"at"`
	Temperature float64   `json:"temperature"`
	PowerPct    float64   `json:"power_pct"`
	Operator    string    `json:"operator"`
}

// Append 追加观测，并拒绝关闭后的写入。
func (l *HeatLog) Append(reading HeatLogReading) error {
	if l == nil || l.Closed {
		return ErrPreconditionFailed
	}
	if reading.At.IsZero() || reading.Temperature <= 0 || reading.PowerPct < 0 || reading.PowerPct > 100 {
		return ErrInvalidField("heat_reading")
	}
	l.Readings = append(l.Readings, reading)
	sort.SliceStable(l.Readings, func(i, j int) bool { return l.Readings[i].At.Before(l.Readings[j].At) })
	return nil
}

// Close 封存炉次日志。
func (l *HeatLog) Close() error {
	if l == nil || len(l.Readings) < 2 {
		return ErrPreconditionFailed
	}
	l.Closed = true
	return nil
}

// PeakTemperature 返回峰值温度及其位置。
func (l HeatLog) PeakTemperature() (float64, time.Time) {
	var peak float64
	var at time.Time
	for _, reading := range l.Readings {
		if reading.Temperature > peak {
			peak = reading.Temperature
			at = reading.At
		}
	}
	return peak, at
}

// AveragePower 返回平均功率百分比。
func (l HeatLog) AveragePower() float64 {
	if len(l.Readings) == 0 {
		return 0
	}
	var total float64
	for _, reading := range l.Readings {
		total += reading.PowerPct
	}
	return total / float64(len(l.Readings))
}

// TemperatureRise 返回首末温度变化。
func (l HeatLog) TemperatureRise() float64 {
	if len(l.Readings) < 2 {
		return 0
	}
	return l.Readings[len(l.Readings)-1].Temperature - l.Readings[0].Temperature
}

// Duration 返回日志覆盖时长。
func (l HeatLog) Duration() time.Duration {
	if len(l.Readings) < 2 {
		return 0
	}
	return l.Readings[len(l.Readings)-1].At.Sub(l.Readings[0].At)
}
