package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/bell-foundry/internal/model"
	"github.com/jb843051627/bell-foundry/internal/tuning"
)

// GetBell 获取成品钟。
func (l *Lab) GetBell(ctx context.Context, id string) (*model.Bell, error) {
	var bell model.Bell
	if err := l.load(ctx, "bell", id, &bell); err != nil {
		return nil, err
	}
	return &bell, nil
}

// RecordPartials 保存五分音实测频率。
func (l *Lab) RecordPartials(ctx context.Context, id string, measured []float64) (*model.Bell, error) {
	bell, err := l.GetBell(ctx, id)
	if err != nil {
		return nil, err
	}
	if !tuning.Complete(measured) {
		return nil, model.ErrInvalidField("partials")
	}
	bell.Partials = append([]float64(nil), measured...)
	bell.TuningStatus = model.TuningUnmeasured
	if err := l.save(ctx, "bell", id, *bell); err != nil {
		return nil, err
	}
	_ = l.event(id, "bell.partials.recorded")
	return bell, nil
}

// EvaluateBell 根据名义频率评估调音状态。
func (l *Lab) EvaluateBell(ctx context.Context, id string, tuneLimit, retuneLimit float64) (*tuning.Result, error) {
	bell, err := l.GetBell(ctx, id)
	if err != nil {
		return nil, err
	}
	if !tuning.Complete(bell.Partials) {
		return nil, model.ErrNotFullyMeasured
	}
	targets := tuning.ExpectedFrequencies(bell.NominalFreqHz)
	result := tuning.Evaluate(targets, bell.Partials, tuneLimit, retuneLimit)
	bell.DetuneCents = result.Mean
	bell.TuningStatus = result.Status
	if err := l.save(ctx, "bell", id, *bell); err != nil {
		return nil, err
	}
	if result.Status != model.TuningInTune {
		_, _ = l.RaiseAlert(ctx, "bell:"+id, "warn", fmt.Sprintf("tuning status %s (%s)", result.Status, result.WorstName))
	}
	return &result, nil
}

// ListBellsNeedingRetune 返回待返修或重铸的钟。
func (l *Lab) ListBellsNeedingRetune(ctx context.Context) ([]model.Bell, error) {
	result := make([]model.Bell, 0)
	err := l.list(ctx, "bell", func(raw []byte) error {
		var bell model.Bell
		if err := decode(raw, &bell); err != nil {
			return err
		}
		if bell.TuningStatus == model.TuningNeedsRetune || bell.TuningStatus == model.TuningRecast {
			result = append(result, bell)
		}
		return nil
	})
	return result, err
}

// TuningPlan 返回调音调整步骤。
func (l *Lab) TuningPlan(ctx context.Context, id string) ([]string, error) {
	result, err := l.EvaluateBell(ctx, id, 8, 35)
	if err != nil {
		return nil, err
	}
	return tuning.RetuneSteps(*result), nil
}
