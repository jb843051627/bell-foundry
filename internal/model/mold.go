package model

import "time"

// 铸型状态机：assembled → drying → dried → closed → poured；任意状态可 scrapped。
const (
	MoldAssembled = "assembled"
	MoldDrying    = "drying"
	MoldDried     = "dried"
	MoldClosed    = "closed"
	MoldPoured    = "poured"
	MoldScrapped  = "scrapped"
)

// MoldTransitions 声明铸型合法状态迁移表。
var MoldTransitions = map[string][]string{
	MoldAssembled: {MoldDrying, MoldScrapped},
	MoldDrying:    {MoldDried, MoldScrapped},
	MoldDried:     {MoldClosed, MoldScrapped},
	MoldClosed:    {MoldPoured, MoldScrapped},
	MoldPoured:    {},
	MoldScrapped:  {},
}

// CanTransition 判断铸型是否允许从 from 迁移到 to。
func CanTransition(from, to string) bool {
	for _, next := range MoldTransitions[from] {
		if next == to {
			return true
		}
	}
	return false
}

// Mold 是一个钟的铸型（外模 + 砂芯合箱体）。
type Mold struct {
	ID            string     `json:"id"`
	ProfileCode   string     `json:"profile_code"`
	Status        string     `json:"status"`
	DryingStartAt *time.Time `json:"drying_start_at,omitempty"`
	DryingHours   float64    `json:"drying_hours"`
	MoisturePct   float64    `json:"moisture_pct"`
	OpenDefects   int        `json:"open_defects"`
	PourID        string     `json:"pour_id,omitempty"`
	ScrapReason   string     `json:"scrap_reason,omitempty"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// IsDryingComplete 判定烘干是否达标：时长不低于 minHours 且含水率不超过上限。
func (m *Mold) IsDryingComplete(minHours, maxMoisturePct float64) bool {
	if m.Status != MoldDrying && m.Status != MoldDried {
		return false
	}
	if m.DryingHours < minHours {
		return false
	}
	if m.MoisturePct > maxMoisturePct {
		return false
	}
	return true
}

// CanClose 合箱前置条件：已烘干且无未解决缺陷。
func (m *Mold) CanClose() bool {
	if m.Status != MoldDried {
		return false
	}
	return m.OpenDefects == 0
}
