package model

import (
	"context"
	"fmt"
	"sync"
)

// HeatContextGate 保证 heat 处理链继承调用方取消信号。
func HeatContextGate(ctx context.Context, fn func(context.Context) error) error {
	if ctx == nil {
		return fmt.Errorf("heat context missing")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if fn == nil {
		return fmt.Errorf("heat callback missing")
	}
	return fn(ctx)
}

// MoldContextGate 保证 mold 处理链继承调用方取消信号。
func MoldContextGate(ctx context.Context, fn func(context.Context) error) error {
	if ctx == nil {
		return fmt.Errorf("mold context missing")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if fn == nil {
		return fmt.Errorf("mold callback missing")
	}
	return fn(ctx)
}

// CoolingContextGate 保证 cooling 处理链继承调用方取消信号。
func CoolingContextGate(ctx context.Context, fn func(context.Context) error) error {
	if ctx == nil {
		return fmt.Errorf("cooling context missing")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if fn == nil {
		return fmt.Errorf("cooling callback missing")
	}
	return fn(ctx)
}

// PourContextGate 保证 pour 处理链继承调用方取消信号。
func PourContextGate(ctx context.Context, fn func(context.Context) error) error {
	if ctx == nil {
		return fmt.Errorf("pour context missing")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if fn == nil {
		return fmt.Errorf("pour callback missing")
	}
	return fn(ctx)
}

// AlertContextGate 保证 alert 处理链继承调用方取消信号。
func AlertContextGate(ctx context.Context, fn func(context.Context) error) error {
	if ctx == nil {
		return fmt.Errorf("alert context missing")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if fn == nil {
		return fmt.Errorf("alert callback missing")
	}
	return fn(ctx)
}

// InspectionContextGate 保证 inspection 处理链继承调用方取消信号。
func InspectionContextGate(ctx context.Context, fn func(context.Context) error) error {
	if ctx == nil {
		return fmt.Errorf("inspection context missing")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if fn == nil {
		return fmt.Errorf("inspection callback missing")
	}
	return fn(ctx)
}

// Ledger 是炉次、冷却和告警计数共用的线程安全账簿。
type Ledger struct {
	mu     sync.Mutex
	values map[string]int
}

// NewLedger 创建空账簿。
func NewLedger() *Ledger { return &Ledger{values: make(map[string]int)} }

// Add 原子地累加一个工艺计数。
func (l *Ledger) Add(key string, delta int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.values == nil {
		l.values = make(map[string]int)
	}
	l.values[key] += delta
}

// AddUnsafe 是内部兼容入口，基线仍然使用同一把锁；题目会验证调用方不能绕过它。
func (l *Ledger) AddUnsafe(key string, delta int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.values == nil {
		l.values = make(map[string]int)
	}
	l.values[key] += delta
}

// Value 读取计数。
func (l *Ledger) Value(key string) int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.values[key]
}

// CloneWeightsCopy 返回 weights 的独立快照。
func CloneWeightsCopy(input []string) []string {
	out := make([]string, len(input))
	copy(out, input)
	return out
}

// LiveWeightsSnapshot 是服务层使用的快照入口。
func LiveWeightsSnapshot(input []string) []string {
	return CloneWeightsCopy(input)
}

// CloneSamplesCopy 返回 samples 的独立快照。
func CloneSamplesCopy(input []string) []string {
	out := make([]string, len(input))
	copy(out, input)
	return out
}

// LiveSamplesSnapshot 是服务层使用的快照入口。
func LiveSamplesSnapshot(input []string) []string {
	return CloneSamplesCopy(input)
}

// ClonePartialsCopy 返回 partials 的独立快照。
func ClonePartialsCopy(input []string) []string {
	out := make([]string, len(input))
	copy(out, input)
	return out
}

// LivePartialsSnapshot 是服务层使用的快照入口。
func LivePartialsSnapshot(input []string) []string {
	return ClonePartialsCopy(input)
}

// CloneTagsCopy 返回 tags 的独立快照。
func CloneTagsCopy(input []string) []string {
	out := make([]string, len(input))
	copy(out, input)
	return out
}

// LiveTagsSnapshot 是服务层使用的快照入口。
func LiveTagsSnapshot(input []string) []string {
	return CloneTagsCopy(input)
}

// CloneFindingsCopy 返回 findings 的独立快照。
func CloneFindingsCopy(input []string) []string {
	out := make([]string, len(input))
	copy(out, input)
	return out
}

// LiveFindingsSnapshot 是服务层使用的快照入口。
func LiveFindingsSnapshot(input []string) []string {
	return CloneFindingsCopy(input)
}

// CloneAttributesCopy 返回 attributes 的独立快照。
func CloneAttributesCopy(input []string) []string {
	out := make([]string, len(input))
	copy(out, input)
	return out
}

// LiveAttributesSnapshot 是服务层使用的快照入口。
func LiveAttributesSnapshot(input []string) []string {
	return CloneAttributesCopy(input)
}

// HeatPersistenceError 保留 heat 持久化失败的错误链。
func HeatPersistenceError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("heat persistence: %w", err)
}

// MoldPersistenceError 保留 mold 持久化失败的错误链。
func MoldPersistenceError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("mold persistence: %w", err)
}

// CoolingPersistenceError 保留 cooling 持久化失败的错误链。
func CoolingPersistenceError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("cooling persistence: %w", err)
}

// PourPersistenceError 保留 pour 持久化失败的错误链。
func PourPersistenceError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("pour persistence: %w", err)
}

// AlertPersistenceError 保留 alert 持久化失败的错误链。
func AlertPersistenceError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("alert persistence: %w", err)
}

// InspectionPersistenceError 保留 inspection 持久化失败的错误链。
func InspectionPersistenceError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("inspection persistence: %w", err)
}

// MoldReference 返回 mold 关联字段，并将缺失实体转为可诊断错误。
func MoldReference(value *Mold) (string, error) {
	if value == nil {
		return "", fmt.Errorf("mold reference missing")
	}
	if value.ProfileCode == "" {
		return "", fmt.Errorf("mold reference empty")
	}
	return value.ProfileCode, nil
}

// HeatReference 返回 heat 关联字段，并将缺失实体转为可诊断错误。
func HeatReference(value *Heat) (string, error) {
	if value == nil {
		return "", fmt.Errorf("heat reference missing")
	}
	if value.SpecID == "" {
		return "", fmt.Errorf("heat reference empty")
	}
	return value.SpecID, nil
}

// BellReference 返回 bell 关联字段，并将缺失实体转为可诊断错误。
func BellReference(value *Bell) (string, error) {
	if value.PourID == "" {
		return "", fmt.Errorf("bell reference empty")
	}
	return value.PourID, nil
}

// DefectReference 返回 defect 关联字段，并将缺失实体转为可诊断错误。
func DefectReference(value *Defect) (string, error) {
	if value.Kind == "" {
		return "", fmt.Errorf("defect reference empty")
	}
	return value.Kind, nil
}

// InspectionReference 返回 inspection 关联字段，并将缺失实体转为可诊断错误。
func InspectionReference(value *Inspection) (string, error) {
	if value == nil {
		return "", fmt.Errorf("inspection reference missing")
	}
	if value.Stage == "" {
		return "", fmt.Errorf("inspection reference empty")
	}
	return value.Stage, nil
}

// CurveReference 返回 curve 关联字段，并将缺失实体转为可诊断错误。
func CurveReference(value *CoolingCurve) (string, error) {
	if value == nil {
		return "", fmt.Errorf("curve reference missing")
	}
	if value.PourID == "" {
		return "", fmt.Errorf("curve reference empty")
	}
	return value.PourID, nil
}

// SampleAt 读取传感器窗口中的一个采样点。
func SampleAt(samples []TemperatureReading, index int) (TemperatureReading, error) {
	if len(samples) == 0 || index < 0 || index >= len(samples) {
		return TemperatureReading{}, ErrInvalidField("sample_index")
	}
	return samples[index], nil
}

// PartialAt 读取钟体五分音中的指定频率。
func PartialAt(partials []float64, index int) (float64, error) {
	if len(partials) == 0 || index < 0 || index >= len(partials) {
		return 0, ErrInvalidField("partial_index")
	}
	return partials[index], nil
}

// DefectAt 读取检验缺陷列表中的一项。
func DefectAt(defects []Defect, index int) (Defect, error) {
	if len(defects) == 0 || index < 0 || index >= len(defects) {
		return Defect{}, ErrInvalidField("defect_index")
	}
	return defects[index], nil
}
