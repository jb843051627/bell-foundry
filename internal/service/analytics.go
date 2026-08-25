package service

import (
	"context"
	"sort"

	"github.com/jb843051627/bell-foundry/internal/cooling"
	"github.com/jb843051627/bell-foundry/internal/model"
)

// StatusCount 是状态聚合的一行。
type StatusCount struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

// CoolingRisk 是一条冷却风险摘要。
type CoolingRisk struct {
	PourID  string  `json:"pour_id"`
	MaxRate float64 `json:"max_rate"`
	TooFast bool    `json:"too_fast"`
	Samples int     `json:"samples"`
}

// QualitySnapshot 汇总质量链路当前状态。
type QualitySnapshot struct {
	OpenDefects     int `json:"open_defects"`
	CriticalDefects int `json:"critical_defects"`
	Unacknowledged  int `json:"unacknowledged"`
	RetuneBells     int `json:"retune_bells"`
}

// MoldStatuses 按铸型状态计数。
func (l *Lab) MoldStatuses(ctx context.Context) ([]StatusCount, error) {
	items, err := l.ListMolds(ctx)
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int)
	for _, item := range items {
		counts[item.Status]++
	}
	return sortedCounts(counts), nil
}

// HeatStatuses 按炉次状态计数。
func (l *Lab) HeatStatuses(ctx context.Context) ([]StatusCount, error) {
	items, err := l.ListHeats(ctx)
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int)
	for _, item := range items {
		counts[item.Status]++
	}
	return sortedCounts(counts), nil
}

// PourOutcomes 按浇注结果计数。
func (l *Lab) PourOutcomes(ctx context.Context) ([]StatusCount, error) {
	items, err := l.ListPours(ctx)
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int)
	for _, item := range items {
		counts[item.Outcome]++
	}
	return sortedCounts(counts), nil
}

// CoolingRisks 扫描所有曲线并重新计算风险。
func (l *Lab) CoolingRisks(ctx context.Context) ([]CoolingRisk, error) {
	result := make([]CoolingRisk, 0)
	err := l.list(ctx, "curve", func(raw []byte) error {
		var curve model.CoolingCurve
		if err := decode(raw, &curve); err != nil {
			return err
		}
		rate := cooling.MaxRate(&curve)
		result = append(result, CoolingRisk{PourID: curve.PourID, MaxRate: rate, TooFast: rate >= 18, Samples: len(curve.Samples)})
		return nil
	})
	sort.Slice(result, func(i, j int) bool { return result[i].MaxRate > result[j].MaxRate })
	return result, err
}

// TuningStatuses 按成品调音状态计数。
func (l *Lab) TuningStatuses(ctx context.Context) ([]StatusCount, error) {
	counts := make(map[string]int)
	err := l.list(ctx, "bell", func(raw []byte) error {
		var item model.Bell
		if err := decode(raw, &item); err != nil {
			return err
		}
		counts[item.TuningStatus]++
		return nil
	})
	return sortedCounts(counts), err
}

// Snapshot 生成质量摘要。
func (l *Lab) Snapshot(ctx context.Context) (*QualitySnapshot, error) {
	defects, err := l.ListOpenDefects(ctx)
	if err != nil {
		return nil, err
	}
	alerts, err := l.ListAlerts(ctx, true)
	if err != nil {
		return nil, err
	}
	bells, err := l.ListBellsNeedingRetune(ctx)
	if err != nil {
		return nil, err
	}
	snapshot := &QualitySnapshot{OpenDefects: len(defects), Unacknowledged: len(alerts), RetuneBells: len(bells)}
	for _, defect := range defects {
		if defect.Severity == model.DefectCritical {
			snapshot.CriticalDefects++
		}
	}
	return snapshot, nil
}

func sortedCounts(counts map[string]int) []StatusCount {
	result := make([]StatusCount, 0, len(counts))
	for status, count := range counts {
		result = append(result, StatusCount{Status: status, Count: count})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Count == result[j].Count {
			return result[i].Status < result[j].Status
		}
		return result[i].Count > result[j].Count
	})
	return result
}
