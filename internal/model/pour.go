package model

import "time"

// 浇注结果。
const (
	PourOK       = "ok"
	PourColdShut = "cold_shut"
	PourSplash   = "splash"
)

// PourRecord 记录一次浇注：哪一炉金属注入哪个铸型。
type PourRecord struct {
	ID          string    `json:"id"`
	HeatID      string    `json:"heat_id"`
	MoldID      string    `json:"mold_id"`
	BatchID     string    `json:"batch_id"`
	PourTempC   float64   `json:"pour_temp_c"`
	FlowSeconds int       `json:"flow_seconds"`
	Outcome     string    `json:"outcome"`
	At          time.Time `json:"at"`
}

// OutcomeFromTemp 依据浇注温度与炉次窗口给出初判结果：
// 低于窗口下限视为冷隔（cold_shut），过高视为飞溅风险，其余正常。
func OutcomeFromTemp(temp, lower, upper float64) string {
	switch {
	case temp < lower:
		return PourColdShut
	case temp > upper:
		return PourSplash
	default:
		return PourOK
	}
}
