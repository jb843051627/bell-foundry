package model

import "time"

// 缺陷严重级别与状态。
const (
	DefectMinor    = "minor"
	DefectMajor    = "major"
	DefectCritical = "critical"

	DefectOpen     = "open"
	DefectResolved = "resolved"
	DefectWaived   = "waived"
)

// Defect 是检验环节归档的缺陷条目。
type Defect struct {
	ID          string     `json:"id"`
	BellID      string     `json:"bell_id,omitempty"`
	MoldID      string     `json:"mold_id,omitempty"`
	Kind        string     `json:"kind"`
	Severity    string     `json:"severity"`
	Stage       string     `json:"stage"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	FiledAt     time.Time  `json:"filed_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
}

// BlocksClosing critical 缺陷是否阻断合箱。
func (d *Defect) BlocksClosing() bool {
	return d.Severity == DefectCritical && d.Status == DefectOpen
}

// IsOpen 是否仍处于打开状态。
func (d *Defect) IsOpen() bool { return d.Status == DefectOpen }
