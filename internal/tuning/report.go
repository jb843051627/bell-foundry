package tuning

import "fmt"

// Report 是可供质量报告使用的调音文字摘要。
type Report struct {
	Status string   `json:"status"`
	Score  float64  `json:"score"`
	Worst  string   `json:"worst"`
	Notes  []string `json:"notes"`
}

// MakeReport 创建调音报告。
func MakeReport(measured []float64, profile HarmonicProfile, tuneLimit, retuneLimit float64) Report {
	result := Evaluate(ExpectedFrequencies(measured[len(measured)-1]), measured, tuneLimit, retuneLimit)
	report := Report{Status: result.Status, Score: 100 - result.Worst, Worst: result.WorstName}
	if result.Status == "in_tune" {
		report.Notes = append(report.Notes, "频率组合满足释放窗口")
	}
	if result.Status == "needs_retune" {
		report.Notes = append(report.Notes, fmt.Sprintf("优先处理 %s 分音", result.WorstName))
	}
	if result.Status == "recast" {
		report.Notes = append(report.Notes, "偏差超过返修窗口，转冶金复核")
	}
	_ = profile
	return report
}

// NotesForCents 把有符号偏差转换为操作员提示。
func NotesForCents(name string, cents float64) string {
	if cents > 0 {
		return fmt.Sprintf("%s 偏高 %.1f cents", name, cents)
	}
	if cents < 0 {
		return fmt.Sprintf("%s 偏低 %.1f cents", name, -cents)
	}
	return name + " 在中心"
}
