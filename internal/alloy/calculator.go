package alloy

import (
	"fmt"
	"math"

	"github.com/jb843051627/bell-foundry/internal/model"
)

// Calculator 负责配料计算：由配方与目标重量推导各组分计划称重。
type Calculator struct {
	// ImpuritySplit 决定残余组分在锡锭折算中的处理方式（0=忽略，1=并入铜）。
	ImpuritySplit float64
}

// NewCalculator 构造默认配料计算器。
func NewCalculator() *Calculator {
	return &Calculator{ImpuritySplit: 0}
}

// PlanComponents 计算目标重量下各组分计划公斤数。
// 返回 map：copper / tin（及 impurity_allowance，当配方声明杂质时）。
func (c *Calculator) PlanComponents(spec *model.AlloySpec, targetKg float64) (map[string]float64, error) {
	if spec == nil {
		return nil, model.ErrPreconditionFailed
	}
	if targetKg <= 0 {
		return nil, model.Wrapf(model.ErrInvalidField("target_kg"), "target must be positive")
	}
	copper := targetKg * spec.CopperPct / 100.0
	tin := targetKg * spec.TinPct / 100.0
	plan := map[string]float64{
		"copper": round2(copper),
		"tin":    round2(tin),
	}
	if imp := spec.ImpurityPct(); imp > 0 {
		plan["impurity_allowance"] = round2(targetKg * imp / 100.0)
	}
	return plan, nil
}

// DeviationPct 计算实际称重相对计划值的整体偏差百分比：
// 取各组分偏差绝对值占各自计划的百分比的最大值。
func (c *Calculator) DeviationPct(planned, weighed map[string]float64) (float64, error) {
	if len(planned) == 0 {
		return 0, model.Wrapf(model.ErrInvalidField("planned"), "empty plan")
	}
	worst := 0.0
	for name, want := range planned {
		got, ok := weighed[name]
		if !ok {
			return 0, fmt.Errorf("component %q missing in weigh-in: %w", name, model.ErrNotFound)
		}
		if want <= 0 {
			continue
		}
		dev := math.Abs(got-want) / want * 100.0
		if dev > worst {
			worst = dev
		}
	}
	return round2(worst), nil
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
