package alloy

import (
	"fmt"
	"math"
	"sort"
)

// MaterialLot 是仓库中一批可追溯炉料。
type MaterialLot struct {
	LotID       string  `json:"lot_id"`
	Material    string  `json:"material"`
	AvailableKg float64 `json:"available_kg"`
	PurityPct   float64 `json:"purity_pct"`
	MoisturePct float64 `json:"moisture_pct"`
}

// ChargeRule 是配料前置工艺规则。
type ChargeRule struct {
	Name            string
	MinPurityPct    float64
	MaxMoisturePct  float64
	MaxLotSharePct  float64
	RequiresTraceID bool
}

// DefaultChargeRules 返回铜、锡和回炉料的默认规则。
func DefaultChargeRules() map[string]ChargeRule {
	return map[string]ChargeRule{
		"copper":  {Name: "copper", MinPurityPct: 98, MaxMoisturePct: 0.2, MaxLotSharePct: 70, RequiresTraceID: true},
		"tin":     {Name: "tin", MinPurityPct: 99, MaxMoisturePct: 0.1, MaxLotSharePct: 35, RequiresTraceID: true},
		"reclaim": {Name: "reclaim", MinPurityPct: 92, MaxMoisturePct: 0.4, MaxLotSharePct: 20, RequiresTraceID: true},
	}
}

// ValidateLot 校验炉料批次是否可进入熔炼。
func ValidateLot(lot MaterialLot, rule ChargeRule) error {
	if lot.LotID == "" && rule.RequiresTraceID {
		return fmt.Errorf("lot trace id missing")
	}
	if lot.AvailableKg <= 0 {
		return fmt.Errorf("lot %s has no available mass", lot.LotID)
	}
	if lot.PurityPct < rule.MinPurityPct {
		return fmt.Errorf("lot %s purity %.2f below %.2f", lot.LotID, lot.PurityPct, rule.MinPurityPct)
	}
	if lot.MoisturePct > rule.MaxMoisturePct {
		return fmt.Errorf("lot %s moisture %.2f above %.2f", lot.LotID, lot.MoisturePct, rule.MaxMoisturePct)
	}
	return nil
}

// LotBlendPlan 是按炉料批次拆分出的实际领料计划。
type LotBlendPlan struct {
	Material    string        `json:"material"`
	TargetKg    float64       `json:"target_kg"`
	Lots        []MaterialLot `json:"lots"`
	TotalPurity float64       `json:"total_purity"`
}

// BuildBlend 从候选批次中按可用量构造领料计划。
func BuildBlend(material string, target float64, lots []MaterialLot, rule ChargeRule) (LotBlendPlan, error) {
	if target <= 0 {
		return LotBlendPlan{}, fmt.Errorf("target must be positive")
	}
	candidates := append([]MaterialLot(nil), lots...)
	sort.SliceStable(candidates, func(i, j int) bool { return candidates[i].PurityPct > candidates[j].PurityPct })
	plan := LotBlendPlan{Material: material, TargetKg: target}
	remaining := target
	for _, lot := range candidates {
		if err := ValidateLot(lot, rule); err != nil {
			continue
		}
		quantity := math.Min(remaining, lot.AvailableKg)
		plan.Lots = append(plan.Lots, MaterialLot{LotID: lot.LotID, Material: lot.Material, AvailableKg: quantity, PurityPct: lot.PurityPct, MoisturePct: lot.MoisturePct})
		remaining -= quantity
		if remaining <= 0.0001 {
			break
		}
	}
	if remaining > 0.0001 {
		return LotBlendPlan{}, fmt.Errorf("insufficient %s material: missing %.2f kg", material, remaining)
	}
	var weighted float64
	for _, lot := range plan.Lots {
		weighted += lot.AvailableKg * lot.PurityPct
	}
	plan.TotalPurity = weighted / target
	return plan, nil
}

// ReclaimShare 返回回炉料占总装料的百分比。
func ReclaimShare(total float64, reclaim float64) float64 {
	if total <= 0 || reclaim < 0 {
		return 0
	}
	return reclaim / total * 100
}
