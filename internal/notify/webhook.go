package notify

import (
	"context"
	"fmt"
	"sync"
)

// Message 是向值班人员投递的告警消息。
type Message struct {
	Subject string
	Level   string
	Body    string
}

// Sink 是通知出口抽象，便于测试和替换为 HTTP webhook。
type Sink interface {
	Send(context.Context, Message) error
}

// MemorySink 将最近消息保存在内存中，用于本地运行和验收探针。
type MemorySink struct {
	mu       sync.Mutex
	messages []Message
}

func (s *MemorySink) Send(ctx context.Context, message Message) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	s.messages = append(s.messages, message)
	s.mu.Unlock()
	return nil
}

func (s *MemorySink) Messages() []Message {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.messages
}

// FormatKey 生成稳定的去重键。
func FormatKey(subject, level, body string) string {
	return fmt.Sprintf("%s|%s|%s", subject, level, body)
}
