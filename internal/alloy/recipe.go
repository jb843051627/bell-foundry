package alloy

import (
	"fmt"
	"math"
	"sort"
)

// RecipeAudit 是配方计划的审计结果。
type RecipeAudit struct {
	TargetKg       float64         `json:"target_kg"`
	PlannedKg      float64         `json:"planned_kg"`
	BalanceErrorKg float64         `json:"balance_error_kg"`
	Components     []string        `json:"components"`
	Warnings       []string        `json:"warnings"`
	Checks         map[string]bool `json:"checks"`
}

// AuditRecipe 检查配料计划总重、组分和残余比例。
func AuditRecipe(planned map[string]float64, targetKg float64) RecipeAudit {
	audit := RecipeAudit{TargetKg: targetKg, Checks: make(map[string]bool)}
	for name, value := range planned {
		if value > 0 {
			audit.PlannedKg += value
			audit.Components = append(audit.Components, name)
		}
	}
	sort.Strings(audit.Components)
	audit.BalanceErrorKg = audit.PlannedKg - targetKg
	audit.Checks["positive_target"] = targetKg > 0
	audit.Checks["components_present"] = len(audit.Components) >= 2
	audit.Checks["balanced"] = math.Abs(audit.BalanceErrorKg) <= math.Max(0.05, targetKg*0.001)
	if !audit.Checks["positive_target"] {
		audit.Warnings = append(audit.Warnings, "target mass must be positive")
	}
	if !audit.Checks["components_present"] {
		audit.Warnings = append(audit.Warnings, "at least two components are expected")
	}
	if !audit.Checks["balanced"] {
		audit.Warnings = append(audit.Warnings, fmt.Sprintf("plan differs from target by %.2f kg", audit.BalanceErrorKg))
	}
	return audit
}

// IsReleaseable 判断配方是否可以送入称量环节。
func (a RecipeAudit) IsReleaseable() bool {
	for _, passed := range a.Checks {
		if !passed {
			return false
		}
	}
	return len(a.Checks) > 0
}

// DominantComponent 返回占比最大的组分。
func DominantComponent(values map[string]float64) (string, float64) {
	name, value := "", 0.0
	for key, current := range values {
		if current > value {
			name, value = key, current
		}
	}
	return name, value
}

// NormalizePlan 将计划重量按 targetKg 缩放。
func NormalizePlan(values map[string]float64, targetKg float64) map[string]float64 {
	result := make(map[string]float64)
	if targetKg <= 0 {
		return result
	}
	total := 0.0
	for _, value := range values {
		if value > 0 {
			total += value
		}
	}
	if total == 0 {
		return result
	}
	for name, value := range values {
		if value > 0 {
			result[name] = math.Round(value/total*targetKg*100) / 100
		}
	}
	return result
}
