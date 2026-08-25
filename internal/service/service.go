package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jb843051627/bell-foundry/internal/alloy"
	"github.com/jb843051627/bell-foundry/internal/clock"
	"github.com/jb843051627/bell-foundry/internal/model"
	"github.com/jb843051627/bell-foundry/internal/notify"
	"github.com/jb843051627/bell-foundry/internal/store"
)

// Lab 是钟铸造工艺服务的组合根。
// 所有实体都通过 Store 持久化，服务方法在写入前检查 context。
type Lab struct {
	repo  *store.Store
	calc  *alloy.Calculator
	clock clock.Clock
	sink  notify.Sink
	dedup *notify.Deduper
	mu    sync.RWMutex
}

// NewLab 创建生产用服务组合。
func NewLab(repo *store.Store) *Lab {
	return NewLabWith(repo, clock.System{}, &notify.MemorySink{})
}

// NewLabWith 允许测试注入时钟和通知出口。
func NewLabWith(repo *store.Store, clk clock.Clock, sink notify.Sink) *Lab {
	if clk == nil {
		clk = clock.System{}
	}
	if sink == nil {
		sink = &notify.MemorySink{}
	}
	return &Lab{
		repo: repo, calc: alloy.NewCalculator(), clock: clk, sink: sink,
		dedup: notify.NewDeduper(10 * time.Minute),
	}
}

// Close 保留组合根生命周期接口；Store 由 main 负责关闭。
func (l *Lab) Close() error { return nil }

// Repository 返回底层仓储，仅供健康检查和证据探针使用。
func (l *Lab) Repository() *store.Store { return l.repo }

func contextErr(ctx context.Context) error {
	if ctx == nil {
		return nil
	}
	return ctx.Err()
}

func (l *Lab) now() time.Time { return l.clock.Now().UTC() }

func (l *Lab) save(ctx context.Context, kind, id string, value any) error {
	if err := contextErr(ctx); err != nil {
		return fmt.Errorf("save %s: %w", kind, err)
	}
	if l.repo == nil {
		return errors.New("service: nil repository")
	}
	return l.repo.Save(kind, id, value)
}

func (l *Lab) load(ctx context.Context, kind, id string, value any) error {
	if err := contextErr(ctx); err != nil {
		return fmt.Errorf("load %s: %w", kind, err)
	}
	if err := l.repo.Load(kind, id, value); err != nil {
		if err == sql.ErrNoRows {
			return model.Wrapf(model.ErrNotFound, "%s/%s", kind, id)
		}
		return err
	}
	return nil
}

func (l *Lab) list(ctx context.Context, kind string, decode func([]byte) error) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	return l.repo.List(kind, decode)
}

func (l *Lab) event(subject, action string) error {
	if l.repo == nil {
		return errors.New("service: nil repository")
	}
	return l.repo.Event(subject, action)
}

func (l *Lab) count(kind string) int {
	n, err := l.repo.Count(kind)
	if err != nil {
		return 0
	}
	return n
}

// Health 返回服务核心组件的可观测状态。
func (l *Lab) Health(ctx context.Context) map[string]any {
	result := map[string]any{"service": "bell-foundry", "time": l.now()}
	if err := contextErr(ctx); err != nil {
		result["context_error"] = err.Error()
		return result
	}
	if l.repo == nil {
		result["database"] = "missing"
		return result
	}
	if err := l.repo.Ping(); err != nil {
		result["database"] = err.Error()
	} else {
		result["database"] = "ok"
	}
	result["specs"] = l.count("spec")
	result["molds"] = l.count("mold")
	result["heats"] = l.count("heat")
	result["bells"] = l.count("bell")
	return result
}
