package service

import (
	"context"
	"strings"

	"github.com/jb843051627/bell-foundry/internal/model"
)

// AuditFilter 控制审计事件查询范围。
type AuditFilter struct {
	EntityType string
	EntityID   string
	Actor      string
	Kind       model.EventKind
}

// AppendAudit 保存一条跨层审计事件。
func (l *Lab) AppendAudit(ctx context.Context, event model.ProcessEvent) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	if event.ID == "" {
		event.ID = model.NewID("audit")
	}
	if event.At.IsZero() {
		event.At = l.now()
	}
	if event.EntityID == "" || event.EntityType == "" {
		return model.ErrInvalidField("audit_entity")
	}
	return l.save(ctx, "audit", event.ID, event)
}

// ListAudits 返回符合过滤器的审计事件。
func (l *Lab) ListAudits(ctx context.Context, filter AuditFilter) ([]model.ProcessEvent, error) {
	result := make([]model.ProcessEvent, 0)
	err := l.list(ctx, "audit", func(raw []byte) error {
		var event model.ProcessEvent
		if err := decode(raw, &event); err != nil {
			return err
		}
		if filter.EntityType != "" && event.EntityType != filter.EntityType {
			return nil
		}
		if filter.EntityID != "" && event.EntityID != filter.EntityID {
			return nil
		}
		if filter.Actor != "" && event.Actor != filter.Actor {
			return nil
		}
		if filter.Kind != "" && event.Kind != filter.Kind {
			return nil
		}
		result = append(result, event)
		return nil
	})
	return result, err
}

// SearchNotes 按关键词检索操作员备注。
func (l *Lab) SearchNotes(ctx context.Context, keyword string) ([]OperatorNote, error) {
	result := make([]OperatorNote, 0)
	keyword = strings.ToLower(strings.TrimSpace(keyword))
	err := l.list(ctx, "note", func(raw []byte) error {
		var note OperatorNote
		if err := decode(raw, &note); err != nil {
			return err
		}
		if keyword == "" || strings.Contains(strings.ToLower(note.Text), keyword) {
			result = append(result, note)
		}
		return nil
	})
	return result, err
}
