package model

import "time"

// EventKind 是可追溯工艺事件类型。
type EventKind string

const (
	EventCreated    EventKind = "created"
	EventTransition EventKind = "transition"
	EventMeasured   EventKind = "measured"
	EventAccepted   EventKind = "accepted"
	EventRejected   EventKind = "rejected"
	EventEscalated  EventKind = "escalated"
)

// ProcessEvent 是跨实体审计事件，不承载业务对象本身。
type ProcessEvent struct {
	ID         string            `json:"id"`
	EntityType string            `json:"entity_type"`
	EntityID   string            `json:"entity_id"`
	Kind       EventKind         `json:"kind"`
	From       string            `json:"from,omitempty"`
	To         string            `json:"to,omitempty"`
	Actor      string            `json:"actor"`
	Note       string            `json:"note,omitempty"`
	At         time.Time         `json:"at"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

// NewProcessEvent 构造标准事件。
func NewProcessEvent(entityType, entityID string, kind EventKind, actor string, at time.Time) ProcessEvent {
	return ProcessEvent{ID: NewID("event"), EntityType: entityType, EntityID: entityID, Kind: kind, Actor: actor, At: at}
}

// Transition 设置状态迁移信息并返回自身。
func (e *ProcessEvent) Transition(from, to string) *ProcessEvent {
	e.From, e.To, e.Kind = from, to, EventTransition
	return e
}

// WithNote 添加操作员备注。
func (e *ProcessEvent) WithNote(note string) *ProcessEvent {
	e.Note = note
	return e
}

// IsStateChange 判断事件是否记录了状态变化。
func (e ProcessEvent) IsStateChange() bool { return e.Kind == EventTransition && e.From != e.To }

// AuditTrail 聚合一个实体的事件序列。
type AuditTrail struct {
	EntityID string         `json:"entity_id"`
	Events   []ProcessEvent `json:"events"`
}

// Append 追加事件并确保时间不倒退。
func (t *AuditTrail) Append(event ProcessEvent) error {
	if t == nil || event.EntityID == "" || event.EntityID != t.EntityID {
		return ErrInvalidField("event")
	}
	if len(t.Events) > 0 && event.At.Before(t.Events[len(t.Events)-1].At) {
		return ErrPreconditionFailed
	}
	t.Events = append(t.Events, event)
	return nil
}

// LastState 返回最后一次状态迁移后的状态。
func (t AuditTrail) LastState() string {
	for i := len(t.Events) - 1; i >= 0; i-- {
		if t.Events[i].IsStateChange() {
			return t.Events[i].To
		}
	}
	return ""
}
