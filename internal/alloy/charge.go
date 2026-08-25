package alloy

import (
	"math"
	"sort"
)

// ChargeItem 是一项可称量的炉料。
type ChargeItem struct {
	Name      string  `json:"name"`
	TargetKg  float64 `json:"target_kg"`
	ActualKg  float64 `json:"actual_kg"`
	Tolerance float64 `json:"tolerance"`
	Deviation float64 `json:"deviation"`
	Status    string  `json:"status"`
}

// BuildCharge 将计划 map 转为稳定排序的称量清单。
func BuildCharge(planned map[string]float64, tolerance float64) []ChargeItem {
	if tolerance <= 0 {
		tolerance = DefaultTolerance.AcceptPct
	}
	items := make([]ChargeItem, 0, len(planned))
	for name, target := range planned {
		if target <= 0 {
			continue
		}
		items = append(items, ChargeItem{Name: name, TargetKg: target, Tolerance: tolerance, Status: "pending"})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return items
}

// ApplyActual 将实际重量填入称量清单并给出单项状态。
func ApplyActual(items []ChargeItem, actual map[string]float64) []ChargeItem {
	out := make([]ChargeItem, len(items))
	copy(out, items)
	for i := range out {
		value, ok := actual[out[i].Name]
		if !ok {
			out[i].Status = "missing"
			continue
		}
		out[i].ActualKg = value
		if out[i].TargetKg == 0 {
			out[i].Status = "invalid"
			continue
		}
		out[i].Deviation = math.Abs(value-out[i].TargetKg) / out[i].TargetKg * 100
		if out[i].Deviation <= out[i].Tolerance {
			out[i].Status = "accepted"
		} else {
			out[i].Status = "review"
		}
	}
	return out
}

// ChargeStatus 汇总称量清单状态。
func ChargeStatus(items []ChargeItem) string {
	if len(items) == 0 {
		return "empty"
	}
	for _, item := range items {
		if item.Status == "missing" || item.Status == "invalid" {
			return "incomplete"
		}
		if item.Status != "accepted" {
			return "review"
		}
	}
	return "accepted"
}

// TotalCharge 汇总目标或实际重量。
func TotalCharge(items []ChargeItem, actual bool) float64 {
	total := 0.0
	for _, item := range items {
		if actual {
			total += item.ActualKg
		} else {
			total += item.TargetKg
		}
	}
	return math.Round(total*100) / 100
}
