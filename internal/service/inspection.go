package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/bell-foundry/internal/cooling"
	"github.com/jb843051627/bell-foundry/internal/model"
	"github.com/jb843051627/bell-foundry/internal/tuning"
)

// CreateInspection 创建一条检验记录。
func (l *Lab) CreateInspection(ctx context.Context, inspection model.Inspection) (*model.Inspection, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	if inspection.Stage == "" || inspection.Inspector == "" {
		return nil, model.ErrInvalidField("inspection")
	}
	if inspection.ID == "" {
		inspection.ID = model.NewID("inspection")
	}
	if inspection.CreatedAt.IsZero() {
		inspection.CreatedAt = l.now()
	}
	if inspection.Verdict == "" {
		inspection.Verdict = model.VerdictReview
	}
	if err := l.save(ctx, "inspection", inspection.ID, inspection); err != nil {
		return nil, err
	}
	_ = l.event(inspection.ID, "inspection.created")
	return &inspection, nil
}

// GetInspection 获取检验记录。
func (l *Lab) GetInspection(ctx context.Context, id string) (*model.Inspection, error) {
	var item model.Inspection
	if err := l.load(ctx, "inspection", id, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

// ListInspections 列出检验记录。
func (l *Lab) ListInspections(ctx context.Context, stage string) ([]model.Inspection, error) {
	result := make([]model.Inspection, 0)
	err := l.list(ctx, "inspection", func(raw []byte) error {
		var item model.Inspection
		if err := decode(raw, &item); err != nil {
			return err
		}
		if stage == "" || item.Stage == stage {
			result = append(result, item)
		}
		return nil
	})
	return result, err
}

// InspectCooling 根据曲线诊断结果创建检验结论。
func (l *Lab) InspectCooling(ctx context.Context, pourID, inspector string) (*model.Inspection, error) {
	curve, err := l.GetCurve(ctx, pourID)
	if err != nil {
		return nil, err
	}
	report := cooling.Diagnose(curve, 1)
	inspection := model.Inspection{BellID: "", Inspector: inspector, Stage: "cooling", Measurements: map[string]float64{"max_rate": report.MaximumRate, "noise_rms": cooling.NoiseRMS(curve)}, Verdict: model.VerdictPass}
	for _, warning := range report.Warnings {
		inspection.AddFinding(warning)
		inspection.Verdict = model.VerdictReview
	}
	return l.CreateInspection(ctx, inspection)
}

// InspectBell 根据五分音计算检验结论。
func (l *Lab) InspectBell(ctx context.Context, bellID, inspector string) (*model.Inspection, error) {
	bell, err := l.GetBell(ctx, bellID)
	if err != nil {
		return nil, err
	}
	if !tuning.Complete(bell.Partials) {
		return nil, model.ErrNotFullyMeasured
	}
	result := tuning.Evaluate(tuning.ExpectedFrequencies(bell.NominalFreqHz), bell.Partials, 8, 35)
	inspection := model.Inspection{BellID: bellID, Inspector: inspector, Stage: "tuning", Measurements: map[string]float64{"worst_cents": result.Worst, "mean_cents": result.Mean}, Verdict: model.VerdictPass}
	if result.Status == model.TuningNeedsRetune {
		inspection.Verdict = model.VerdictReview
		inspection.AddFinding("partial deviation requires retune")
	}
	if result.Status == model.TuningRecast {
		inspection.Verdict = model.VerdictHold
		inspection.AddFinding(fmt.Sprintf("%s exceeds recast limit", result.WorstName))
	}
	return l.CreateInspection(ctx, inspection)
}
