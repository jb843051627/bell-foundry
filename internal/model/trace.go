package model

import "time"

// TraceEvent 将一个实体在工艺链路中的变化记录为可追溯事件。
type TraceEvent struct {
	ID         string            `json:"id"`
	EntityType string            `json:"entity_type"`
	EntityID   string            `json:"entity_id"`
	Action     string            `json:"action"`
	Actor      string            `json:"actor"`
	Payload    map[string]string `json:"payload,omitempty"`
	At         time.Time         `json:"at"`
}

// TraceChain 是单个实体的事件链。
type TraceChain struct {
	EntityType string       `json:"entity_type"`
	EntityID   string       `json:"entity_id"`
	Events     []TraceEvent `json:"events"`
}

// Append 追加事件，并保证事件 ID 和时间存在。
func (c *TraceChain) Append(event TraceEvent) {
	if event.ID == "" {
		event.ID = NewID("trace")
	}
	if event.At.IsZero() {
		event.At = time.Now().UTC()
	}
	event.EntityType = c.EntityType
	event.EntityID = c.EntityID
	c.Events = append(c.Events, event)
}

// LastAction 返回最近动作。
func (c *TraceChain) LastAction() string {
	if len(c.Events) == 0 {
		return ""
	}
	return c.Events[len(c.Events)-1].Action
}

// HasAction 判断事件链是否出现指定动作。
func (c *TraceChain) HasAction(action string) bool {
	for _, event := range c.Events {
		if event.Action == action {
			return true
		}
	}
	return false
}
