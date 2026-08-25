package model

import "time"

// 炉次状态：charging → heating → ready → poured / aborted。
const (
	HeatCharging = "charging"
	HeatHeating  = "heating"
	HeatReady    = "ready"
	HeatPoured   = "poured"
	HeatAborted  = "aborted"
)

// Heat 是一次熔炼作业，目标温度窗口以 TargetTempC ± WindowC 表达。
type Heat struct {
	ID            string     `json:"id"`
	FurnaceNo     int        `json:"furnace_no"`
	SpecID        string     `json:"spec_id"`
	BatchID       string     `json:"batch_id"`
	ChargeKg      float64    `json:"charge_kg"`
	TargetTempC   float64    `json:"target_temp_c"`
	WindowC       float64    `json:"window_c"`
	MeasuredTempC float64    `json:"measured_temp_c"`
	Status        string     `json:"status"`
	StartedAt     time.Time  `json:"started_at"`
	ReadyAt       *time.Time `json:"ready_at,omitempty"`
	PouredAt      *time.Time `json:"poured_at,omitempty"`
}

// TempInWindow 判断实测温度是否落在浇注窗口内。
func (h *Heat) TempInWindow(temp float64) bool {
	lower := h.TargetTempC - h.WindowC
	upper := h.TargetTempC + h.WindowC
	return temp >= lower && temp <= upper
}

// ElapsedMinutes 返回炉次从点火到当前的分钟数（未结束按当前时间计）。
func (h *Heat) ElapsedMinutes(now time.Time) float64 {
	end := now
	if h.PouredAt != nil {
		end = *h.PouredAt
	} else if h.ReadyAt != nil {
		end = *h.ReadyAt
	}
	d := end.Sub(h.StartedAt)
	if d < 0 {
		return 0
	}
	return d.Minutes()
}
