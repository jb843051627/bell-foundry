package service

import (
	"context"
	"sort"

	"github.com/jb843051627/bell-foundry/internal/model"
)

// AppendTrace 追加实体追溯事件。
func (l *Lab) AppendTrace(ctx context.Context, event model.TraceEvent) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	if event.EntityType == "" || event.EntityID == "" || event.Action == "" {
		return model.ErrInvalidField("trace")
	}
	if event.ID == "" {
		event.ID = model.NewID("trace")
	}
	if event.At.IsZero() {
		event.At = l.now()
	}
	return l.save(ctx, "trace", event.ID, event)
}

// TraceEntity 获取某个实体的追溯链。
func (l *Lab) TraceEntity(ctx context.Context, entityType, entityID string) (*model.TraceChain, error) {
	chain := &model.TraceChain{EntityType: entityType, EntityID: entityID}
	err := l.list(ctx, "trace", func(raw []byte) error {
		var event model.TraceEvent
		if err := decode(raw, &event); err != nil {
			return err
		}
		if event.EntityType == entityType && event.EntityID == entityID {
			chain.Events = append(chain.Events, event)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.SliceStable(chain.Events, func(i, j int) bool { return chain.Events[i].At.Before(chain.Events[j].At) })
	return chain, nil
}

// TraceSummary 是追溯链的简化视图。
type TraceSummary struct {
	EntityType  string   `json:"entity_type"`
	EntityID    string   `json:"entity_id"`
	EventCount  int      `json:"event_count"`
	FirstAction string   `json:"first_action"`
	LastAction  string   `json:"last_action"`
	Actors      []string `json:"actors"`
}

// SummarizeTrace 将追溯链压缩为报告字段。
func SummarizeTrace(chain *model.TraceChain) TraceSummary {
	result := TraceSummary{}
	if chain == nil {
		return result
	}
	result.EntityType, result.EntityID, result.EventCount = chain.EntityType, chain.EntityID, len(chain.Events)
	seen := make(map[string]bool)
	for i, event := range chain.Events {
		if i == 0 {
			result.FirstAction = event.Action
		}
		result.LastAction = event.Action
		if event.Actor != "" && !seen[event.Actor] {
			result.Actors = append(result.Actors, event.Actor)
			seen[event.Actor] = true
		}
	}
	return result
}

// TraceWithNote 方便在业务操作后写入带备注的追溯事件。
func (l *Lab) TraceWithNote(ctx context.Context, entityType, entityID, action, actor, note string) error {
	return l.AppendTrace(ctx, model.TraceEvent{EntityType: entityType, EntityID: entityID, Action: action, Actor: actor, Payload: map[string]string{"note": note}, At: l.now()})
}
