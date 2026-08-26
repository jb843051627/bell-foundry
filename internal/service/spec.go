package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/bell-foundry/internal/alloy"
	"github.com/jb843051627/bell-foundry/internal/model"
)

// CreateSpec 保存并校验合金配方。
func (l *Lab) CreateSpec(ctx context.Context, spec model.AlloySpec) (*model.AlloySpec, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	if spec.ID == "" {
		spec.ID = model.NewID(model.PrefixSpec)
	}
	if spec.CreatedAt.IsZero() {
		spec.CreatedAt = l.now()
	}
	if err := spec.Validate(); err != nil {
		return nil, err
	}
	if err := l.save(ctx, "spec", spec.ID, spec); err != nil {
		return nil, err
	}
	_ = l.event(spec.ID, "spec.created")
	return &spec, nil
}

// GetSpec 获取配方。
func (l *Lab) GetSpec(ctx context.Context, id string) (*model.AlloySpec, error) {
	var spec model.AlloySpec
	if err := l.load(ctx, "spec", id, &spec); err != nil {
		return nil, err
	}
	return &spec, nil
}

// ListSpecs 列出全部配方。
func (l *Lab) ListSpecs(ctx context.Context) ([]model.AlloySpec, error) {
	result := make([]model.AlloySpec, 0)
	err := l.list(ctx, "spec", func(raw []byte) error {
		var spec model.AlloySpec
		if err := decode(raw, &spec); err != nil {
			return err
		}
		result = append(result, spec)
		return nil
	})
	return result, err
}

// PlanBatch 按配方计算一次配料计划。
func (l *Lab) PlanBatch(ctx context.Context, specID string, targetKg float64) (*model.AlloyBatch, error) {
	spec, err := l.GetSpec(ctx, specID)
	if err != nil {
		return nil, err
	}
	planned, err := l.calc.PlanComponents(spec, targetKg)
	if err != nil {
		return nil, err
	}
	batch := model.AlloyBatch{ID: model.NewID(model.PrefixBatch), SpecID: specID, TargetKg: targetKg, Planned: planned, Status: "planned", CreatedAt: l.now()}
	if err := l.save(ctx, "batch", batch.ID, batch); err != nil {
		return nil, err
	}
	_ = l.event(batch.ID, "batch.planned")
	return &batch, nil
}

// RecordWeighIn 保存实际称重并按容差更新批次状态。
func (l *Lab) RecordWeighIn(ctx context.Context, id string, weighed map[string]float64) (*model.AlloyBatch, error) {
	var batch model.AlloyBatch
	if err := l.load(ctx, "batch", id, &batch); err != nil {
		return nil, err
	}
	deviation, err := l.calc.DeviationPct(batch.Planned, weighed)
	if err != nil {
		return nil, err
	}
	batch.Weighed = cloneWeights(weighed)
	batch.DeviationPct = deviation
	batch.Status = alloy.ClassifyDeviation(deviation, alloy.DefaultTolerance)
	if err := l.save(ctx, "batch", batch.ID, batch); err != nil {
		return nil, err
	}
	_ = l.event(batch.ID, fmt.Sprintf("batch.weighed.%s", batch.Status))
	return &batch, nil
}

// GetBatch 获取配料批次。
func (l *Lab) GetBatch(ctx context.Context, id string) (*model.AlloyBatch, error) {
	var batch model.AlloyBatch
	if err := l.load(ctx, "batch", id, &batch); err != nil {
		return nil, err
	}
	return &batch, nil
}

// ListBatches 列出配料批次。
func (l *Lab) ListBatches(ctx context.Context) ([]model.AlloyBatch, error) {
	result := make([]model.AlloyBatch, 0)
	err := l.list(ctx, "batch", func(raw []byte) error {
		var item model.AlloyBatch
		if err := decode(raw, &item); err != nil {
			return err
		}
		result = append(result, item)
		return nil
	})
	return result, err
}

func cloneWeights(input map[string]float64) map[string]float64 {
	out := make(map[string]float64, len(input))
	for key, value := range input {
		out[key] = value
	}
	return out
}
