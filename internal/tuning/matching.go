package tuning

import "sort"

// Candidate 是一个可作为名义音参考的频率候选。
type Candidate struct {
	Name      string  `json:"name"`
	Frequency float64 `json:"frequency"`
	Distance  float64 `json:"distance_cents"`
}

// RankCandidates 按与测得名义音的音分距离排序。
func RankCandidates(measured float64, references map[string]float64) []Candidate {
	result := make([]Candidate, 0, len(references))
	for name, reference := range references {
		if reference <= 0 {
			continue
		}
		result = append(result, Candidate{Name: name, Frequency: reference, Distance: AbsCents(CentsBetween(reference, measured))})
	}
	sort.SliceStable(result, func(i, j int) bool { return result[i].Distance < result[j].Distance })
	return result
}

// MatchProfile 返回距离在窗口内的所有候选。
func MatchProfile(measured float64, references map[string]float64, limitCents float64) []Candidate {
	all := RankCandidates(measured, references)
	result := make([]Candidate, 0, len(all))
	for _, candidate := range all {
		if candidate.Distance <= limitCents {
			result = append(result, candidate)
		}
	}
	return result
}

// Consensus 计算一组参考候选的中位频率。
func Consensus(candidates []Candidate) float64 {
	if len(candidates) == 0 {
		return 0
	}
	copyOf := append([]Candidate(nil), candidates...)
	sort.Slice(copyOf, func(i, j int) bool { return copyOf[i].Frequency < copyOf[j].Frequency })
	middle := len(copyOf) / 2
	if len(copyOf)%2 == 1 {
		return copyOf[middle].Frequency
	}
	return (copyOf[middle-1].Frequency + copyOf[middle].Frequency) / 2
}

// DriftReport 对比两次名义音测量。
type DriftReport struct {
	FirstHz   float64 `json:"first_hz"`
	SecondHz  float64 `json:"second_hz"`
	Cents     float64 `json:"cents"`
	Direction string  `json:"direction"`
}

// CompareDrift 生成频率漂移报告。
func CompareDrift(first, second float64) DriftReport {
	report := DriftReport{FirstHz: first, SecondHz: second, Cents: CentsBetween(first, second)}
	if report.Cents > 0.5 {
		report.Direction = "sharp"
	} else if report.Cents < -0.5 {
		report.Direction = "flat"
	} else {
		report.Direction = "stable"
	}
	return report
}
