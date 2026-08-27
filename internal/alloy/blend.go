package alloy

import "math"

// BlendPlan 描述回炉料与新料混合后的组成。
type BlendPlan struct {
	FreshKg   map[string]float64 `json:"fresh_kg"`
	ReclaimKg map[string]float64 `json:"reclaim_kg"`
	TotalKg   float64            `json:"total_kg"`
	Ratios    map[string]float64 `json:"ratios"`
}

// NewBlendPlan 合并新料与回炉料，忽略非正重量。
func NewBlendPlan(fresh, reclaim map[string]float64) BlendPlan {
	plan := BlendPlan{FreshKg: clone(fresh), ReclaimKg: clone(reclaim), Ratios: make(map[string]float64)}
	for _, value := range plan.FreshKg {
		if value > 0 {
			plan.TotalKg += value
		}
	}
	for _, value := range plan.ReclaimKg {
		if value > 0 {
			plan.TotalKg += value
		}
	}
	if plan.TotalKg > 0 {
		for name, value := range plan.FreshKg {
			if value > 0 {
				plan.Ratios["fresh:"+name] = value / plan.TotalKg
			}
		}
		for name, value := range plan.ReclaimKg {
			if value > 0 {
				plan.Ratios["reclaim:"+name] = value / plan.TotalKg
			}
		}
	}
	plan.TotalKg = math.Round(plan.TotalKg*100) / 100
	return plan
}

// ReclaimRatio 返回回炉料占总重量的比例。
func (p BlendPlan) ReclaimRatio() float64 {
	if p.TotalKg <= 0 {
		return 0
	}
	total := 0.0
	for _, value := range p.ReclaimKg {
		if value > 0 {
			total += value
		}
	}
	return total / p.TotalKg
}

// WithinReclaimLimit 判断回炉料比例是否在限制内。
func (p BlendPlan) WithinReclaimLimit(limit float64) bool { return p.ReclaimRatio() <= limit }

// ComponentMass 汇总同名的新料和回炉料重量。
func (p BlendPlan) ComponentMass(name string) float64 { return p.FreshKg[name] + p.ReclaimKg[name] }

func clone(input map[string]float64) map[string]float64 {
	out := make(map[string]float64, len(input))
	for key, value := range input {
		out[key] = value
	}
	return out
}
