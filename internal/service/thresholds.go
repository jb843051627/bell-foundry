package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/bell-foundry/internal/model"
)

// ThresholdBook 是可持久化的工艺卡版本。
type ThresholdBook struct {
	ID      string                  `json:"id"`
	Name    string                  `json:"name"`
	Version int                     `json:"version"`
	Values  model.ProcessThresholds `json:"values"`
	Active  bool                    `json:"active"`
}

// SaveThresholds 保存一版工艺阈值。
func (l *Lab) SaveThresholds(ctx context.Context, book ThresholdBook) (*ThresholdBook, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	if book.Name == "" || !book.Values.Validate() {
		return nil, model.ErrInvalidField("threshold_book")
	}
	if book.ID == "" {
		book.ID = model.NewID("threshold")
	}
	if book.Version <= 0 {
		book.Version = 1
	}
	if err := l.save(ctx, "threshold", book.ID, book); err != nil {
		return nil, err
	}
	return &book, nil
}

// ActiveThresholds 返回当前生效工艺卡；没有配置时返回默认值。
func (l *Lab) ActiveThresholds(ctx context.Context) (model.ProcessThresholds, error) {
	var active model.ProcessThresholds
	err := l.list(ctx, "threshold", func(raw []byte) error {
		var book ThresholdBook
		if err := decode(raw, &book); err != nil {
			return err
		}
		if book.Active {
			active = book.Values
		}
		return nil
	})
	if err != nil {
		return active, err
	}
	if !active.Validate() {
		active = model.DefaultThresholds()
	}
	return active, nil
}

// ExplainThreshold 返回某项阈值的当前解释，供操作员页面使用。
func (l *Lab) ExplainThreshold(ctx context.Context, name string) (string, error) {
	values, err := l.ActiveThresholds(ctx)
	if err != nil {
		return "", err
	}
	switch name {
	case "dry_hours":
		return fmt.Sprintf("铸型至少烘干 %.1f 小时", values.MinDryHours), nil
	case "moisture":
		return fmt.Sprintf("合箱前含水率不超过 %.2f%%", values.MaxMoisturePct), nil
	case "cooling_rate":
		return fmt.Sprintf("最大冷却速率 %.1f ℃/min", values.MaxCoolingRate), nil
	case "tuning":
		return fmt.Sprintf("调音释放窗口 ±%.1f cents", values.TuneLimitCents), nil
	default:
		return "", model.ErrInvalidField("threshold_name")
	}
}
