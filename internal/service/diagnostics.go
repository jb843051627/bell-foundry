package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/bell-foundry/internal/cooling"
	"github.com/jb843051627/bell-foundry/internal/model"
)

// DiagnosticSummary 是一条可直接用于告警和现场页面的诊断。
type DiagnosticSummary struct {
	EntityID string   `json:"entity_id"`
	Status   string   `json:"status"`
	Findings []string `json:"findings"`
	Advice   []string `json:"advice"`
}

// DiagnoseCooling 对浇注冷却链路进行跨层诊断。
func (l *Lab) DiagnoseCooling(ctx context.Context, pourID string) (*DiagnosticSummary, error) {
	curve, err := l.GetCurve(ctx, pourID)
	if err != nil {
		return nil, err
	}
	report := cooling.Diagnose(curve, 2)
	result := &DiagnosticSummary{EntityID: pourID, Status: "normal"}
	result.Findings = append(result.Findings, report.Warnings...)
	if report.SampleCount < 4 {
		result.Status = "incomplete"
		result.Advice = append(result.Advice, "补采相变前后的温度点")
	}
	if !report.Monotonic {
		result.Status = "review"
		result.Advice = append(result.Advice, "检查热电偶固定和时间戳顺序")
	}
	if report.MaximumRate > 18 {
		result.Status = "critical"
		result.Advice = append(result.Advice, "检查砂型保温层与浇注落差")
	}
	if len(result.Findings) == 0 {
		result.Findings = append(result.Findings, "curve is internally consistent")
	}
	return result, nil
}

// DiagnoseMold 对铸型是否可以合箱给出解释性结果。
func (l *Lab) DiagnoseMold(ctx context.Context, moldID string, thresholds model.ProcessThresholds) (*DiagnosticSummary, error) {
	mold, err := l.GetMold(ctx, moldID)
	if err != nil {
		return nil, err
	}
	result := &DiagnosticSummary{EntityID: moldID, Status: "ready"}
	if mold.Status != model.MoldDried {
		result.Status = "blocked"
		result.Findings = append(result.Findings, fmt.Sprintf("status is %s", mold.Status))
		result.Advice = append(result.Advice, "完成烘干并记录含水率")
	}
	if mold.DryingHours < thresholds.MinDryHours {
		result.Status = "blocked"
		result.Findings = append(result.Findings, "drying time below threshold")
	}
	if mold.MoisturePct > thresholds.MaxMoisturePct {
		result.Status = "blocked"
		result.Findings = append(result.Findings, "moisture above threshold")
	}
	if mold.OpenDefects > thresholds.MaxOpenDefects {
		result.Status = "blocked"
		result.Findings = append(result.Findings, "open defects block closing")
	}
	if result.Status == "ready" {
		result.Advice = append(result.Advice, "可以进入合箱")
	}
	return result, nil
}
