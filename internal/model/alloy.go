package model

import "time"

// AlloySpec 定义一种钟青铜合金配方，成分以质量百分比计。
type AlloySpec struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	CopperPct      float64   `json:"copper_pct"`
	TinPct         float64   `json:"tin_pct"`
	MinTinPct      float64   `json:"min_tin_pct"`
	MaxImpurityPct float64   `json:"max_impurity_pct"`
	LiquidusC      float64   `json:"liquidus_c"`
	SolidusC       float64   `json:"solidus_c"`
	CreatedAt      time.Time `json:"created_at"`
}

// Validate 校验配方的成分约束，返回第一条违规原因。
func (s *AlloySpec) Validate() error {
	if s.Name == "" {
		return ErrInvalidField("name")
	}
	if s.CopperPct <= 0 || s.TinPct <= 0 {
		return ErrInvalidField("composition")
	}
	if s.CopperPct+s.TinPct > 100.0 {
		return ErrCompositionOverflow
	}
	if s.TinPct < s.MinTinPct {
		return ErrTinBelowFloor
	}
	if s.LiquidusC <= s.SolidusC {
		return ErrPhasePointsInverted
	}
	return nil
}

// ImpurityPct 计算除铜锡外残余组分的百分比。
func (s *AlloySpec) ImpurityPct() float64 {
	rest := 100.0 - s.CopperPct - s.TinPct
	if rest < 0 {
		return 0
	}
	return rest
}

// MeetsBellGrade 判断配方是否满足钟青铜等级（锡下限与杂质上限）。
func (s *AlloySpec) MeetsBellGrade() bool {
	return s.TinPct >= s.MinTinPct && s.ImpurityPct() <= s.MaxImpurityPct
}

// AlloyBatch 是一次配料作业：按配方计划称重并记录实际值。
type AlloyBatch struct {
	ID           string             `json:"id"`
	SpecID       string             `json:"spec_id"`
	TargetKg     float64            `json:"target_kg"`
	Planned      map[string]float64 `json:"planned"`
	Weighed      map[string]float64 `json:"weighed,omitempty"`
	DeviationPct float64            `json:"deviation_pct"`
	Status       string             `json:"status"`
	CreatedAt    time.Time          `json:"created_at"`
}

// TotalPlanned 汇总计划组分重量。
func (b *AlloyBatch) TotalPlanned() float64 {
	total := 0.0
	for _, kg := range b.Planned {
		total += kg
	}
	return total
}

// TotalWeighed 汇总实际称重。
func (b *AlloyBatch) TotalWeighed() float64 {
	total := 0.0
	for _, kg := range b.Weighed {
		total += kg
	}
	return total
}

// ComponentMissing 返回实际称重中缺失的组分名列表。
func (b *AlloyBatch) ComponentMissing() []string {
	var missing []string
	for name := range b.Planned {
		if _, ok := b.Weighed[name]; !ok {
			missing = append(missing, name)
		}
	}
	return missing
}
