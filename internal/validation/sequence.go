package validation

import (
	"fmt"
	"sort"
)

// SequencePoint 是带序号的工艺事件点。
type SequencePoint struct {
	Sequence int64
	Minute   float64
	Label    string
}

// SortAndValidate 将事件按序号排序并检查重复。
func SortAndValidate(points []SequencePoint) ([]SequencePoint, error) {
	copyOf := append([]SequencePoint(nil), points...)
	sort.SliceStable(copyOf, func(i, j int) bool { return copyOf[i].Sequence < copyOf[j].Sequence })
	for i, point := range copyOf {
		if point.Sequence < 0 || point.Label == "" {
			return nil, fmt.Errorf("invalid sequence point %d", i)
		}
		if i > 0 && point.Sequence == copyOf[i-1].Sequence {
			return nil, fmt.Errorf("duplicate sequence %d", point.Sequence)
		}
		if i > 0 && point.Minute < copyOf[i-1].Minute {
			return nil, fmt.Errorf("time moves backwards at sequence %d", point.Sequence)
		}
	}
	return copyOf, nil
}

// MissingLabels 返回预期标签中未出现的项。
func MissingLabels(points []SequencePoint, expected []string) []string {
	seen := make(map[string]bool, len(points))
	for _, point := range points {
		seen[point.Label] = true
	}
	missing := make([]string, 0)
	for _, label := range expected {
		if !seen[label] {
			missing = append(missing, label)
		}
	}
	return missing
}

// Duration 返回排序后事件链覆盖时长。
func Duration(points []SequencePoint) float64 {
	if len(points) < 2 {
		return 0
	}
	sorted, err := SortAndValidate(points)
	if err != nil {
		return 0
	}
	return sorted[len(sorted)-1].Minute - sorted[0].Minute
}
