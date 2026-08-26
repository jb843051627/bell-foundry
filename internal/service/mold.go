package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/bell-foundry/internal/model"
)

// CreateMold 创建一个待烘干铸型。
func (l *Lab) CreateMold(ctx context.Context, profile string) (*model.Mold, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	if profile == "" {
		return nil, model.ErrInvalidField("profile_code")
	}
	mold := model.Mold{ID: model.NewID(model.PrefixMold), ProfileCode: profile, Status: model.MoldAssembled, UpdatedAt: l.now()}
	if err := l.save(ctx, "mold", mold.ID, mold); err != nil {
		return nil, err
	}
	_ = l.event(mold.ID, "mold.created")
	return &mold, nil
}

// GetMold 获取铸型。
func (l *Lab) GetMold(ctx context.Context, id string) (*model.Mold, error) {
	var mold model.Mold
	if err := l.load(ctx, "mold", id, &mold); err != nil {
		return nil, err
	}
	return &mold, nil
}

// ListMolds 列出铸型。
func (l *Lab) ListMolds(ctx context.Context) ([]model.Mold, error) {
	result := make([]model.Mold, 0)
	err := l.list(ctx, "mold", func(raw []byte) error {
		var item model.Mold
		if err := decode(raw, &item); err != nil {
			return err
		}
		result = append(result, item)
		return nil
	})
	return result, err
}

// StartDrying 将铸型送入烘房。
func (l *Lab) StartDrying(ctx context.Context, id string) (*model.Mold, error) {
	mold, err := l.GetMold(ctx, id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransition(mold.Status, model.MoldDrying) {
		return nil, fmt.Errorf("mold %s: %w", id, model.ErrBadTransition)
	}
	now := l.now()
	mold.Status = model.MoldDrying
	mold.DryingStartAt = &now
	mold.DryingHours = 0
	mold.UpdatedAt = now
	if err := l.save(ctx, "mold", id, *mold); err != nil {
		return nil, err
	}
	_ = l.event(id, "mold.drying.started")
	return mold, nil
}

// RecordMoisture 记录含水率和烘干时长。
func (l *Lab) RecordMoisture(ctx context.Context, id string, moisture, hours float64) (*model.Mold, error) {
	mold, err := l.GetMold(ctx, id)
	if err != nil {
		return nil, err
	}
	if mold.Status != model.MoldDrying {
		return nil, model.ErrPreconditionFailed
	}
	if moisture < 0 || moisture > 100 || hours < 0 {
		return nil, model.ErrInvalidField("moisture_or_hours")
	}
	mold.MoisturePct = moisture
	mold.DryingHours = hours
	mold.UpdatedAt = l.now()
	if err := l.save(ctx, "mold", id, *mold); err != nil {
		return nil, err
	}
	return mold, nil
}

// CompleteDrying 完成烘干检查。
func (l *Lab) CompleteDrying(ctx context.Context, id string, minHours, maxMoisture float64) (*model.Mold, error) {
	mold, err := l.GetMold(ctx, id)
	if err != nil {
		return nil, err
	}
	if !mold.IsDryingComplete(minHours, maxMoisture) {
		return nil, fmt.Errorf("mold %s: %w", id, model.ErrPreconditionFailed)
	}
	mold.Status = model.MoldDried
	mold.UpdatedAt = l.now()
	if err := l.save(ctx, "mold", id, *mold); err != nil {
		return nil, err
	}
	_ = l.event(id, "mold.drying.completed")
	return mold, nil
}

// CloseMold 合箱；存在 critical 未解决缺陷时拒绝。
func (l *Lab) CloseMold(ctx context.Context, id string) (*model.Mold, error) {
	mold, err := l.GetMold(ctx, id)
	if err != nil {
		return nil, err
	}
	if !mold.CanClose() {
		return nil, fmt.Errorf("mold %s: %w", id, model.ErrPreconditionFailed)
	}
	mold.Status = model.MoldClosed
	mold.UpdatedAt = l.now()
	if err := l.save(ctx, "mold", id, *mold); err != nil {
		return nil, err
	}
	_ = l.event(id, "mold.closed")
	return mold, nil
}

// ScrapMold 将铸型报废，允许从任意未浇注状态执行。
func (l *Lab) ScrapMold(ctx context.Context, id, reason string) (*model.Mold, error) {
	mold, err := l.GetMold(ctx, id)
	if err != nil {
		return nil, err
	}
	if mold.Status == model.MoldPoured || mold.Status == model.MoldScrapped {
		return nil, model.ErrBadTransition
	}
	mold.Status = model.MoldScrapped
	mold.ScrapReason = nonEmpty(reason, "operator decision")
	mold.UpdatedAt = l.now()
	if err := l.save(ctx, "mold", id, *mold); err != nil {
		return nil, err
	}
	_ = l.event(id, "mold.scrapped")
	return mold, nil
}

// DryingAge 返回铸型开始烘干以来的时长。
func (l *Lab) DryingAge(ctx context.Context, id string) (time.Duration, error) {
	mold, err := l.GetMold(ctx, id)
	if err != nil {
		return 0, err
	}
	if mold.DryingStartAt == nil {
		return 0, nil
	}
	return l.now().Sub(*mold.DryingStartAt), nil
}
