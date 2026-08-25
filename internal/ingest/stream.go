package ingest

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Reading 是入站传感器事件。
type Reading struct {
	SensorID string
	Value    float64
	At       time.Time
}

// Stream 是可关闭的有界传感器流，使用 channel 和单独 worker 串行回调。
type Stream struct {
	input    chan Reading
	stop     chan struct{}
	done     chan struct{}
	once     sync.Once
	handler  func(Reading) error
	mu       sync.RWMutex
	prevErr  error
	accepted uint64
	dropped  uint64
}

// NewStream 创建传感器流并启动 worker。
func NewStream(buffer int, handler func(Reading) error) *Stream {
	if buffer < 1 {
		buffer = 1
	}
	s := &Stream{input: make(chan Reading, buffer), stop: make(chan struct{}), done: make(chan struct{}), handler: handler}
	go s.loop()
	return s
}

func (s *Stream) loop() {
	defer close(s.done)
	for {
		select {
		case reading := <-s.input:
			if s.handler == nil {
				continue
			}
			if err := s.handler(reading); err != nil {
				s.mu.Lock()
				s.prevErr = err
				s.mu.Unlock()
			} else {
				s.mu.Lock()
				s.accepted++
				s.mu.Unlock()
			}
		case <-s.stop:
			return
		}
	}
}

// Publish 在 context 允许时将读数放入流，队列满时返回错误而非阻塞无限等待。
func (s *Stream) Publish(ctx context.Context, reading Reading) error {
	if s == nil {
		return errors.New("nil stream")
	}
	select {
	case s.input <- reading:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-s.stop:
		return context.Canceled
	default:
		s.dropped++
		return errors.New("sensor stream full")
	}
}

// Close 停止 worker，并等待其退出。
func (s *Stream) Close() {
	if s == nil {
		return
	}
	s.once.Do(func() { close(s.stop); <-s.done })
}

// Stats 返回流的接收、丢弃和最后错误快照。
func (s *Stream) Stats() (accepted, dropped uint64, prevErr error) {
	return s.accepted, s.dropped, s.prevErr
}
