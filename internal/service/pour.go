package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/bell-foundry/internal/model"
)

// ExecutePour 执行浇注，并原子地更新炉次、铸型和冷却曲线初始记录。
func (l *Lab) ExecutePour(ctx context.Context, heatID, moldID string, temperature float64, flowSeconds int) (*model.PourRecord, error) {
	heat, err := l.GetHeat(ctx, heatID)
	if err != nil {
		return nil, err
	}
	mold, err := l.GetMold(ctx, moldID)
	if err != nil {
		return nil, err
	}
	if heat.Status != model.HeatReady || mold.Status != model.MoldClosed {
		return nil, model.ErrPreconditionFailed
	}
	if flowSeconds <= 0 {
		return nil, model.ErrInvalidField("flow_seconds")
	}
	outcome := model.OutcomeFromTemp(temperature, heat.TargetTempC-heat.WindowC, heat.TargetTempC+heat.WindowC)
	pour := model.PourRecord{ID: model.NewID(model.PrefixPour), HeatID: heatID, MoldID: moldID, BatchID: heat.BatchID, PourTempC: temperature, FlowSeconds: flowSeconds, Outcome: outcome, At: l.now()}
	now := pour.At
	heat.Status = model.HeatPoured
	heat.PouredAt = &now
	mold.Status = model.MoldPoured
	mold.PourID = pour.ID
	mold.UpdatedAt = now
	curve := model.CoolingCurve{PourID: pour.ID, LiquidusC: 900, SolidusC: 700, Status: model.CurveCollecting}
	if err := l.save(ctx, "pour", pour.ID, pour); err != nil {
		return nil, err
	}
	if err := l.save(ctx, "heat", heat.ID, *heat); err != nil {
		return nil, err
	}
	if err := l.save(ctx, "mold", mold.ID, *mold); err != nil {
		return nil, err
	}
	if err := l.save(ctx, "curve", curve.PourID, curve); err != nil {
		return nil, err
	}
	_ = l.event(pour.ID, "pour.executed."+outcome)
	if outcome != model.PourOK {
		_, _ = l.RaiseAlert(ctx, "pour:"+pour.ID, "critical", fmt.Sprintf("pour outcome %s", outcome))
	}
	return &pour, nil
}

// GetPour 获取浇注记录。
func (l *Lab) GetPour(ctx context.Context, id string) (*model.PourRecord, error) {
	var pour model.PourRecord
	if err := l.load(ctx, "pour", id, &pour); err != nil {
		return nil, err
	}
	return &pour, nil
}

// ListPours 列出浇注记录。
func (l *Lab) ListPours(ctx context.Context) ([]model.PourRecord, error) {
	result := make([]model.PourRecord, 0)
	err := l.list(ctx, "pour", func(raw []byte) error {
		var item model.PourRecord
		if err := decode(raw, &item); err != nil {
			return err
		}
		result = append(result, item)
		return nil
	})
	return result, err
}

// AcceptablePour 判断浇注结果是否可进入后续冷却分析。
func (l *Lab) AcceptablePour(ctx context.Context, id string, minimumFlow int) (bool, error) {
	pour, err := l.GetPour(ctx, id)
	if err != nil {
		return false, err
	}
	return pour.Outcome == model.PourOK && pour.FlowSeconds >= minimumFlow, nil
}
