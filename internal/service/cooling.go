package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/bell-foundry/internal/alloy"
	curvecalc "github.com/jb843051627/bell-foundry/internal/cooling"
	"github.com/jb843051627/bell-foundry/internal/model"
)

// GetCurve 获取浇注对应的冷却曲线。
func (l *Lab) GetCurve(ctx context.Context, pourID string) (*model.CoolingCurve, error) {
	var curve model.CoolingCurve
	if err := l.load(ctx, "curve", pourID, &curve); err != nil {
		return nil, err
	}
	return &curve, nil
}

// AddCoolingSample 追加一个传感器采样点。
func (l *Lab) AddCoolingSample(ctx context.Context, pourID string, minute, temperature float64) (*model.CoolingCurve, error) {
	curve, err := l.GetCurve(ctx, pourID)
	if err != nil {
		return nil, err
	}
	if curve.Status != model.CurveCollecting {
		return nil, model.ErrBadTransition
	}
	if !curvecalc.Append(curve, model.Sample{Minute: minute, TempC: temperature}) {
		return nil, model.ErrInvalidField("sample")
	}
	if err := l.save(ctx, "curve", pourID, *curve); err != nil {
		return nil, err
	}
	return curve, nil
}

// AnalyzeCooling 分析相变点和最大降温速率。
func (l *Lab) AnalyzeCooling(ctx context.Context, pourID string) (*model.CoolingCurve, error) {
	curve, err := l.GetCurve(ctx, pourID)
	if err != nil {
		return nil, err
	}
	if len(curve.Samples) < 4 {
		return nil, model.ErrCurveIncomplete
	}
	phase := curvecalc.DetectPhase(curve)
	if !phase.Crossed {
		return nil, fmt.Errorf("curve %s has no solidus crossing: %w", pourID, model.ErrCurveIncomplete)
	}
	curve.PhaseMinute = phase.Minute
	curve.MaxRateCPM = curvecalc.MaxRate(curve)
	_, _, curve.TooFast = curvecalc.EvaluateRate(curve, 4, 18)
	curve.Status = model.CurveAnalyzed
	if err := l.save(ctx, "curve", pourID, *curve); err != nil {
		return nil, err
	}
	if curve.TooFast {
		_, _ = l.RaiseAlert(ctx, "curve:"+pourID, "warn", "cooling rate exceeds process band")
	}
	return curve, nil
}

// FinalizeBell 在冷却分析完成后生成成品钟记录。
func (l *Lab) FinalizeBell(ctx context.Context, pourID string, profile string, diameter, mass float64) (*model.Bell, error) {
	curve, err := l.GetCurve(ctx, pourID)
	if err != nil {
		return nil, err
	}
	if curve.Status != model.CurveAnalyzed || !curve.BelowSolidus() {
		return nil, model.ErrPreconditionFailed
	}
	pour, err := l.GetPour(ctx, pourID)
	if err != nil {
		return nil, err
	}
	if mass <= 0 {
		if p, ok := alloy.LookupProfile(profile); ok {
			if estimate, valid := alloy.EstimateMass(p, diameter); valid {
				mass = estimate
			}
		}
	}
	if diameter <= 0 || mass <= 0 {
		return nil, model.ErrInvalidField("diameter_or_mass")
	}
	bell := model.Bell{ID: model.NewID(model.PrefixBell), PourID: pourID, MoldID: pour.MoldID, MassKg: mass, DiameterMm: diameter, NominalFreqHz: alloy.EstimateNominalHz(diameter), TuningStatus: model.TuningUnmeasured, CastAt: l.now()}
	curve.Status = model.CurveFinalized
	if err := l.save(ctx, "curve", pourID, *curve); err != nil {
		return nil, err
	}
	if err := l.save(ctx, "bell", bell.ID, bell); err != nil {
		return nil, err
	}
	_ = l.event(bell.ID, "bell.created")
	return &bell, nil
}

// CoolingAdvice 将曲线分析结果转为工艺建议。
func (l *Lab) CoolingAdvice(ctx context.Context, pourID string) (string, error) {
	curve, err := l.GetCurve(ctx, pourID)
	if err != nil {
		return "", err
	}
	if curve.TooFast {
		return "检查砂型保温层并降低冷却梯度", nil
	}
	if curve.MaxRateCPM < 1 {
		return "检查炉料温度与砂芯含水率", nil
	}
	return "冷却曲线在工艺窗口内", nil
}
