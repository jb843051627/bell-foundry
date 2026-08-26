package service

import (
	"context"
	"time"
)

// DailySummary 是工艺日报的聚合结果。
type DailySummary struct {
	Date           string `json:"date"`
	Specs          int    `json:"specs"`
	Batches        int    `json:"batches"`
	Molds          int    `json:"molds"`
	Heats          int    `json:"heats"`
	Pours          int    `json:"pours"`
	Bells          int    `json:"bells"`
	OpenDefects    int    `json:"open_defects"`
	Unacknowledged int    `json:"unacknowledged_alerts"`
	NeedsRetune    int    `json:"needs_retune"`
}

// DailyReport 生成指定日期的当前聚合快照。
func (l *Lab) DailyReport(ctx context.Context, day time.Time) (*DailySummary, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
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
	specs, err := l.ListSpecs(ctx)
	if err != nil {
		return nil, err
	}
	batches, err := l.ListBatches(ctx)
	if err != nil {
		return nil, err
	}
	molds, err := l.ListMolds(ctx)
	if err != nil {
		return nil, err
	}
	heats, err := l.ListHeats(ctx)
	if err != nil {
		return nil, err
	}
	pours, err := l.ListPours(ctx)
	if err != nil {
		return nil, err
	}
	return &DailySummary{
		Date: day.Format("2006-01-02"), Specs: len(specs), Batches: len(batches), Molds: len(molds),
		Heats: len(heats), Pours: len(pours), Bells: l.count("bell"), OpenDefects: len(defects),
		Unacknowledged: len(alerts), NeedsRetune: len(bells),
	}, nil
}

// TodayReport 使用服务时钟生成当天报告。
func (l *Lab) TodayReport(ctx context.Context) (*DailySummary, error) {
	return l.DailyReport(ctx, l.now())
}

// EventCount 返回事件总数，供运维页面显示。
func (l *Lab) EventCount() int {
	return l.count("event")
}
