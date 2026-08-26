package model

import "time"

// SensorCalibration 记录温度传感器的校准参数和有效期。
type SensorCalibration struct {
	ID           string    `json:"id"`
	SensorID     string    `json:"sensor_id"`
	OffsetC      float64   `json:"offset_c"`
	Gain         float64   `json:"gain"`
	ReferenceC   float64   `json:"reference_c"`
	ToleranceC   float64   `json:"tolerance_c"`
	CalibratedAt time.Time `json:"calibrated_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	Technician   string    `json:"technician"`
	Status       string    `json:"status"`
}

const (
	CalibrationValid   = "valid"
	CalibrationExpired = "expired"
	CalibrationReview  = "review"
)

// Apply 将原始温度转换为校准后的温度。
func (c *SensorCalibration) Apply(raw float64) float64 { return raw*c.Gain + c.OffsetC }

// InTolerance 判断测量结果与参考温度的误差。
func (c *SensorCalibration) InTolerance(value float64) bool {
	limit := c.ToleranceC
	if limit <= 0 {
		limit = 2
	}
	delta := c.Apply(value) - c.ReferenceC
	return delta <= limit && delta >= -limit
}

// RefreshStatus 依据当前时间更新有效状态。
func (c *SensorCalibration) RefreshStatus(now time.Time) string {
	if c.ExpiresAt.IsZero() || now.Before(c.ExpiresAt) {
		c.Status = CalibrationValid
	} else {
		c.Status = CalibrationExpired
	}
	return c.Status
}
