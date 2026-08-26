package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/bell-foundry/internal/model"
)

// StartHeat 创建并启动炉次。
func (l *Lab) StartHeat(ctx context.Context, furnace int, specID, batchID string, chargeKg, targetC, windowC float64) (*model.Heat, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	if furnace <= 0 || chargeKg <= 0 || targetC <= 0 || windowC <= 0 {
		return nil, model.ErrInvalidField("heat_parameters")
	}
	if _, err := l.GetSpec(ctx, specID); err != nil {
		return nil, err
	}
	heat := model.Heat{ID: model.NewID(model.PrefixHeat), FurnaceNo: furnace, SpecID: specID, BatchID: batchID, ChargeKg: chargeKg, TargetTempC: targetC, WindowC: windowC, Status: model.HeatCharging, StartedAt: l.now()}
	if err := l.save(ctx, "heat", heat.ID, heat); err != nil {
		return nil, err
	}
	_ = l.event(heat.ID, "heat.started")
	return &heat, nil
}

// GetHeat 获取炉次。
func (l *Lab) GetHeat(ctx context.Context, id string) (*model.Heat, error) {
	var heat model.Heat
	if err := l.load(ctx, "heat", id, &heat); err != nil {
		return nil, err
	}
	return &heat, nil
}

// ListHeats 列出炉次。
func (l *Lab) ListHeats(ctx context.Context) ([]model.Heat, error) {
	result := make([]model.Heat, 0)
	err := l.list(ctx, "heat", func(raw []byte) error {
		var item model.Heat
		if err := decode(raw, &item); err != nil {
			return err
		}
		result = append(result, item)
		return nil
	})
	return result, err
}

// RecordTemperature 记录炉温并将炉次推进到 heating。
func (l *Lab) RecordTemperature(ctx context.Context, id string, temperature float64) (*model.Heat, error) {
	heat, err := l.GetHeat(ctx, id)
	if err != nil {
		return nil, err
	}
	if heat.Status != model.HeatCharging && heat.Status != model.HeatHeating {
		return nil, fmt.Errorf("heat %s: %w", id, model.ErrBadTransition)
	}
	if temperature <= 0 {
		return nil, model.ErrInvalidField("temperature")
	}
	heat.MeasuredTempC = temperature
	heat.Status = model.HeatHeating
	if err := l.save(ctx, "heat", id, *heat); err != nil {
		return nil, err
	}
	return heat, nil
}

// MarkHeatReady 检查温度窗口并推进炉次。
func (l *Lab) MarkHeatReady(ctx context.Context, id string, temperature float64) (*model.Heat, error) {
	heat, err := l.GetHeat(ctx, id)
	if err != nil {
		return nil, err
	}
	if heat.Status != model.HeatHeating && heat.Status != model.HeatCharging {
		return nil, model.ErrBadTransition
	}
	if !heat.TempInWindow(temperature) {
		return nil, fmt.Errorf("temperature %.2f outside window: %w", temperature, model.ErrPreconditionFailed)
	}
	now := l.now()
	heat.MeasuredTempC = temperature
	heat.Status = model.HeatReady
	heat.ReadyAt = &now
	if err := l.save(ctx, "heat", id, *heat); err != nil {
		return nil, err
	}
	_ = l.event(id, "heat.ready")
	return heat, nil
}

// AbortHeat 中止未浇注炉次。
func (l *Lab) AbortHeat(ctx context.Context, id string, reason string) (*model.Heat, error) {
	heat, err := l.GetHeat(ctx, id)
	if err != nil {
		return nil, err
	}
	if heat.Status == model.HeatPoured || heat.Status == model.HeatAborted {
		return nil, model.ErrBadTransition
	}
	heat.Status = model.HeatAborted
	if err := l.save(ctx, "heat", id, *heat); err != nil {
		return nil, err
	}
	_ = l.event(id, "heat.aborted:"+reason)
	return heat, nil
}

// HeatReady 是浇注服务使用的状态检查。
func (l *Lab) HeatReady(ctx context.Context, id string) (bool, error) {
	heat, err := l.GetHeat(ctx, id)
	if err != nil {
		return false, err
	}
	return heat.Status == model.HeatReady, nil
}
