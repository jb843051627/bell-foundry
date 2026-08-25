package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/bell-foundry/internal/model"
)

// OperatorNote 是对工艺实体的人工备注。
type OperatorNote struct {
	ID         string `json:"id"`
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
	Author     string `json:"author"`
	Text       string `json:"text"`
	At         string `json:"at"`
}

// RecordOperatorNote 保存一条操作员备注，便于后续轨迹复核。
func (l *Lab) RecordOperatorNote(ctx context.Context, entityType, entityID, author, text string) (*OperatorNote, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	if entityType == "" || entityID == "" || author == "" || text == "" {
		return nil, model.ErrInvalidField("operator_note")
	}
	note := OperatorNote{ID: model.NewID("note"), EntityType: entityType, EntityID: entityID, Author: author, Text: text, At: l.now().Format("2006-01-02T15:04:05.999Z07:00")}
	if err := l.save(ctx, "note", note.ID, note); err != nil {
		return nil, err
	}
	_ = l.event(entityID, "note.added")
	return &note, nil
}

// ReopenMold 将误标为 dried 的铸型退回 drying，要求重新记录含水率。
func (l *Lab) ReopenMold(ctx context.Context, id, reason string) (*model.Mold, error) {
	mold, err := l.GetMold(ctx, id)
	if err != nil {
		return nil, err
	}
	if mold.Status != model.MoldDried && mold.Status != model.MoldClosed {
		return nil, model.ErrBadTransition
	}
	mold.Status = model.MoldDrying
	mold.DryingHours = 0
	mold.MoisturePct = 100
	mold.UpdatedAt = l.now()
	if err := l.save(ctx, "mold", id, *mold); err != nil {
		return nil, err
	}
	_ = l.event(id, "mold.reopened:"+reason)
	return mold, nil
}

// ValidatePourReadiness 聚合浇注前的跨实体检查。
func (l *Lab) ValidatePourReadiness(ctx context.Context, heatID, moldID string) error {
	heat, err := l.GetHeat(ctx, heatID)
	if err != nil {
		return err
	}
	mold, err := l.GetMold(ctx, moldID)
	if err != nil {
		return err
	}
	if heat.Status != model.HeatReady {
		return fmt.Errorf("heat not ready: %w", model.ErrPreconditionFailed)
	}
	if mold.Status != model.MoldClosed {
		return fmt.Errorf("mold not closed: %w", model.ErrPreconditionFailed)
	}
	if mold.OpenDefects > 0 {
		return fmt.Errorf("mold has open defects: %w", model.ErrPreconditionFailed)
	}
	return nil
}

// BulkCloseEligibleMolds 尝试批量合箱，返回成功和失败的 ID。
func (l *Lab) BulkCloseEligibleMolds(ctx context.Context, ids []string) (closed []model.Mold, failed map[string]string) {
	failed = make(map[string]string)
	for _, id := range ids {
		item, err := l.CloseMold(ctx, id)
		if err != nil {
			failed[id] = err.Error()
			continue
		}
		closed = append(closed, *item)
	}
	return closed, failed
}

// DeleteDraftBatch 删除尚未称重的计划批次。
func (l *Lab) DeleteDraftBatch(ctx context.Context, id string) error {
	batch, err := l.GetBatch(ctx, id)
	if err != nil {
		return err
	}
	if batch.Status != "planned" {
		return model.ErrPreconditionFailed
	}
	if err := contextErr(ctx); err != nil {
		return err
	}
	if err := l.repo.Delete("batch", id); err != nil {
		return err
	}
	return l.event(id, "batch.deleted")
}

// FindBellsByStatus 按调音状态检索成品钟。
func (l *Lab) FindBellsByStatus(ctx context.Context, status string) ([]model.Bell, error) {
	result := make([]model.Bell, 0)
	err := l.list(ctx, "bell", func(raw []byte) error {
		var item model.Bell
		if err := decode(raw, &item); err != nil {
			return err
		}
		if item.TuningStatus == status {
			result = append(result, item)
		}
		return nil
	})
	return result, err
}
