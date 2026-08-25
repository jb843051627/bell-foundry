package model

import "time"

// Inspection 是成品检验记录，既可以关联钟体，也可以关联铸型或炉次。
type Inspection struct {
	ID           string             `json:"id"`
	BellID       string             `json:"bell_id,omitempty"`
	MoldID       string             `json:"mold_id,omitempty"`
	HeatID       string             `json:"heat_id,omitempty"`
	Inspector    string             `json:"inspector"`
	Stage        string             `json:"stage"`
	Measurements map[string]float64 `json:"measurements"`
	Findings     []string           `json:"findings"`
	Verdict      string             `json:"verdict"`
	CreatedAt    time.Time          `json:"created_at"`
}

const (
	VerdictPass   = "pass"
	VerdictReview = "review"
	VerdictHold   = "hold"
)

// IsPassing 判断检验是否可以进入下一道工艺。
func (i *Inspection) IsPassing() bool { return i.Verdict == VerdictPass }

// AddFinding 追加一条去重后的检查发现。
func (i *Inspection) AddFinding(finding string) {
	if finding == "" {
		return
	}
	for _, existing := range i.Findings {
		if existing == finding {
			return
		}
	}
	i.Findings = append(i.Findings, finding)
}

// SetMeasurement 记录测量值并拒绝非有限数。
func (i *Inspection) SetMeasurement(name string, value float64) bool {
	if name == "" || value != value || value > 1e12 || value < -1e12 {
		return false
	}
	if i.Measurements == nil {
		i.Measurements = make(map[string]float64)
	}
	i.Measurements[name] = value
	return true
}
