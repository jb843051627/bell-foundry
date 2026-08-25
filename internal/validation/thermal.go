package validation

import (
	"fmt"
	"math"

	"github.com/jb843051627/bell-foundry/internal/model"
)

// ThermalPoint 是一个经过基本范围检查的热过程点。
type ThermalPoint struct {
	Minute       float64 `json:"minute"`
	TemperatureC float64 `json:"temperature_c"`
}

// ValidateThermalPoint 校验热过程点的物理范围。
func ValidateThermalPoint(point ThermalPoint) error {
	if point.Minute < 0 || math.IsNaN(point.Minute) {
		return FieldError{Field: "minute", Message: "must be non-negative"}
	}
	if point.TemperatureC < -50 || point.TemperatureC > 1800 || math.IsNaN(point.TemperatureC) {
		return FieldError{Field: "temperature_c", Message: "outside sensor range"}
	}
	return nil
}

// ValidateCoolingSequence 校验时间严格递增，并返回详细问题。
func ValidateCoolingSequence(points []ThermalPoint) error {
	if len(points) < 2 {
		return fmt.Errorf("at least two thermal points required")
	}
	for i, point := range points {
		if err := ValidateThermalPoint(point); err != nil {
			return err
		}
		if i > 0 && point.Minute <= points[i-1].Minute {
			return fmt.Errorf("minute at %d is not increasing", i)
		}
	}
	return nil
}

// CoolingSlope 返回首尾温度斜率。
func CoolingSlope(points []ThermalPoint) (float64, error) {
	if err := ValidateCoolingSequence(points); err != nil {
		return 0, err
	}
	span := points[len(points)-1].Minute - points[0].Minute
	if span <= 0 {
		return 0, fmt.Errorf("zero time span")
	}
	return (points[len(points)-1].TemperatureC - points[0].TemperatureC) / span, nil
}

// TemperatureBand 判断所有点是否位于工艺温区。
func TemperatureBand(points []ThermalPoint, lower, upper float64) ([]int, error) {
	if lower >= upper {
		return nil, fmt.Errorf("invalid thermal band")
	}
	bad := make([]int, 0)
	for i, point := range points {
		if point.TemperatureC < lower || point.TemperatureC > upper {
			bad = append(bad, i)
		}
	}
	return bad, nil
}

// CurveFromPoints 将验证后的点转换为领域曲线。
func CurveFromPoints(pourID string, liquidus, solidus float64, points []ThermalPoint) (*model.CoolingCurve, error) {
	if liquidus <= solidus {
		return nil, fmt.Errorf("liquidus must exceed solidus")
	}
	if err := ValidateCoolingSequence(points); err != nil {
		return nil, err
	}
	curve := &model.CoolingCurve{PourID: pourID, LiquidusC: liquidus, SolidusC: solidus, Status: model.CurveCollecting}
	for _, point := range points {
		curve.AddSample(model.Sample{Minute: point.Minute, TempC: point.TemperatureC})
	}
	return curve, nil
}
