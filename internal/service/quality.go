package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/bell-foundry/internal/model"
	"github.com/jb843051627/bell-foundry/internal/notify"
)

// FileDefect 记录缺陷，并同步铸型的开放缺陷计数。
func (l *Lab) FileDefect(ctx context.Context, defect model.Defect) (*model.Defect, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	if defect.ID == "" {
		defect.ID = model.NewID(model.PrefixDefect)
	}
	if defect.Kind == "" || defect.Severity == "" {
		return nil, model.ErrInvalidField("defect")
	}
	if defect.Status == "" {
		defect.Status = model.DefectOpen
	}
	if defect.FiledAt.IsZero() {
		defect.FiledAt = l.now()
	}
	if err := l.save(ctx, "defect", defect.ID, defect); err != nil {
		return nil, err
	}
	if defect.MoldID != "" && defect.Status == model.DefectOpen {
		if mold, err := l.GetMold(ctx, defect.MoldID); err == nil {
			mold.OpenDefects++
			_ = l.save(ctx, "mold", mold.ID, *mold)
		}
	}
	if defect.Severity == model.DefectCritical {
		_, _ = l.RaiseAlert(ctx, "defect:"+defect.ID, notify.EscalationLevel(defect.Severity), defect.Description)
	}
	return &defect, nil
}

// ResolveDefect 关闭或豁免缺陷。
func (l *Lab) ResolveDefect(ctx context.Context, id, status string) (*model.Defect, error) {
	defect := model.Defect{}
	if err := l.load(ctx, "defect", id, &defect); err != nil {
		return nil, err
	}
	if defect.Status != model.DefectOpen {
		return nil, model.ErrBadTransition
	}
	if status != model.DefectResolved && status != model.DefectWaived {
		return nil, model.ErrInvalidField("status")
	}
	defect.Status = status
	now := l.now()
	defect.ResolvedAt = &now
	if err := l.save(ctx, "defect", id, defect); err != nil {
		return nil, err
	}
	if defect.MoldID != "" {
		if mold, err := l.GetMold(ctx, defect.MoldID); err == nil && mold.OpenDefects > 0 {
			mold.OpenDefects--
			_ = l.save(ctx, "mold", mold.ID, *mold)
		}
	}
	return &defect, nil
}

// ListOpenDefects 列出未解决缺陷。
func (l *Lab) ListOpenDefects(ctx context.Context) ([]model.Defect, error) {
	result := make([]model.Defect, 0)
	err := l.list(ctx, "defect", func(raw []byte) error {
		var item model.Defect
		if err := decode(raw, &item); err != nil {
			return err
		}
		if item.Status == model.DefectOpen {
			result = append(result, item)
		}
		return nil
	})
	return result, err
}

// RaiseAlert 创建去重后的告警并投递到通知出口。
func (l *Lab) RaiseAlert(ctx context.Context, subject, level, message string) (*model.Alert, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	level = notify.NormalizeLevel(level)
	key := notify.FormatKey(subject, level, message)
	if !l.dedup.Allow(key, l.now()) {
		return nil, fmt.Errorf("alert suppressed: %s", subject)
	}
	alert := model.Alert{ID: model.NewID(model.PrefixAlert), Subject: subject, Level: level, Message: message, CreatedAt: l.now()}
	if err := l.save(ctx, "alert", alert.ID, alert); err != nil {
		return nil, err
	}
	if err := l.sink.Send(ctx, notify.Message{Subject: subject, Level: level, Body: message}); err != nil {
		return nil, err
	}
	return &alert, nil
}

// ListAlerts 列出告警，可选择只看未确认。
func (l *Lab) ListAlerts(ctx context.Context, onlyOpen bool) ([]model.Alert, error) {
	result := make([]model.Alert, 0)
	err := l.list(ctx, "alert", func(raw []byte) error {
		var item model.Alert
		if err := decode(raw, &item); err != nil {
			return err
		}
		if !onlyOpen || !item.Acknowledged {
			result = append(result, item)
		}
		return nil
	})
	return result, err
}

// AcknowledgeAlert 确认告警。
func (l *Lab) AcknowledgeAlert(ctx context.Context, id string) (*model.Alert, error) {
	var alert model.Alert
	if err := l.load(ctx, "alert", id, &alert); err != nil {
		return nil, err
	}
	alert.Acknowledged = true
	if err := l.save(ctx, "alert", id, alert); err != nil {
		return nil, err
	}
	return &alert, nil
}
