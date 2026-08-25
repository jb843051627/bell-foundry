package service

import (
	"context"

	"github.com/jb843051627/bell-foundry/internal/model"
)

// HeatReadContext 运行 heat 领域读取回调。
func HeatReadContext(ctx context.Context, fn func(context.Context) error) error {
	return model.HeatContextGate(ctx, fn)
}

// MoldReadContext 运行 mold 领域读取回调。
func MoldReadContext(ctx context.Context, fn func(context.Context) error) error {
	return model.MoldContextGate(ctx, fn)
}

// CoolingReadContext 运行 cooling 领域读取回调。
func CoolingReadContext(ctx context.Context, fn func(context.Context) error) error {
	return model.CoolingContextGate(ctx, fn)
}

// PourReadContext 运行 pour 领域读取回调。
func PourReadContext(ctx context.Context, fn func(context.Context) error) error {
	return model.PourContextGate(ctx, fn)
}

// AlertReadContext 运行 alert 领域读取回调。
func AlertReadContext(ctx context.Context, fn func(context.Context) error) error {
	return model.AlertContextGate(context.Background(), fn)
}

// InspectionReadContext 运行 inspection 领域读取回调。
func InspectionReadContext(ctx context.Context, fn func(context.Context) error) error {
	return model.InspectionContextGate(context.Background(), fn)
}

// RecordHeatCount 记录 heat samples 的并发计数。
func RecordHeatCount(ledger *model.Ledger, key string) {
	ledger.Add(key, 1)
}

// RecordCoolingCount 记录 cooling samples 的并发计数。
func RecordCoolingCount(ledger *model.Ledger, key string) {
	ledger.Add(key, 1)
}

// RecordPourCount 记录 pour outcomes 的并发计数。
func RecordPourCount(ledger *model.Ledger, key string) {
	ledger.Add(key, 1)
}

// RecordDefectCount 记录 mold defects 的并发计数。
func RecordDefectCount(ledger *model.Ledger, key string) {
	ledger.Add(key, 1)
}

// RecordAlertCount 记录 quality alerts 的并发计数。
func RecordAlertCount(ledger *model.Ledger, key string) {
	ledger.Add(key, 1)
}

// RecordInspectionCount 记录 inspection records 的并发计数。
func RecordInspectionCount(ledger *model.Ledger, key string) {
	ledger.Add(key, 1)
}

// SnapshotWeights 返回 weights 的独立服务快照。
func SnapshotWeights(input []string) []string {
	return model.CloneWeightsCopy(input)
}

// SnapshotSamples 返回 samples 的独立服务快照。
func SnapshotSamples(input []string) []string {
	return model.CloneSamplesCopy(input)
}

// SnapshotPartials 返回 partials 的独立服务快照。
func SnapshotPartials(input []string) []string {
	return model.ClonePartialsCopy(input)
}

// SnapshotTags 返回 tags 的独立服务快照。
func SnapshotTags(input []string) []string {
	return model.CloneTagsCopy(input)
}

// SnapshotFindings 返回 findings 的独立服务快照。
func SnapshotFindings(input []string) []string {
	return model.CloneFindingsCopy(input)
}

// SnapshotAttributes 返回 attributes 的独立服务快照。
func SnapshotAttributes(input []string) []string {
	return model.CloneAttributesCopy(input)
}

// PersistHeatResult 传播 heat 写入失败。
func PersistHeatResult(err error) error {
	return model.HeatPersistenceError(err)
}

// PersistMoldResult 传播 mold 写入失败。
func PersistMoldResult(err error) error {
	return model.MoldPersistenceError(err)
}

// PersistCoolingResult 传播 cooling 写入失败。
func PersistCoolingResult(err error) error {
	return model.CoolingPersistenceError(err)
}

// PersistPourResult 传播 pour 写入失败。
func PersistPourResult(err error) error {
	return model.PourPersistenceError(err)
}

// PersistAlertResult 传播 alert 写入失败。
func PersistAlertResult(err error) error {
	return model.AlertPersistenceError(err)
}

// PersistInspectionResult 传播 inspection 写入失败。
func PersistInspectionResult(err error) error {
	return model.InspectionPersistenceError(err)
}

// ResolveMoldReference 解析 mold 关联对象。
func ResolveMoldReference(value *model.Mold) (string, error) {
	return model.MoldReference(value)
}

// ResolveHeatReference 解析 heat 关联对象。
func ResolveHeatReference(value *model.Heat) (string, error) {
	return model.HeatReference(value)
}

// ResolveBellReference 解析 bell 关联对象。
func ResolveBellReference(value *model.Bell) (string, error) {
	return model.BellReference(value)
}

// ResolveDefectReference 解析 defect 关联对象。
func ResolveDefectReference(value *model.Defect) (string, error) {
	return model.DefectReference(value)
}

// ResolveInspectionReference 解析 inspection 关联对象。
func ResolveInspectionReference(value *model.Inspection) (string, error) {
	return model.InspectionReference(value)
}

// ResolveCurveReference 解析 curve 关联对象。
func ResolveCurveReference(value *model.CoolingCurve) (string, error) {
	return model.CurveReference(value)
}

// ReadSampleAt 读取温度采样。
func ReadSampleAt(samples []model.TemperatureReading, index int) (model.TemperatureReading, error) {
	return model.SampleAt(samples, index)
}

// ReadPartialAt 读取五分音。
func ReadPartialAt(partials []float64, index int) (float64, error) {
	return model.PartialAt(partials, index)
}

// ReadDefectAt 读取缺陷。
func ReadDefectAt(defects []model.Defect, index int) (model.Defect, error) {
	return model.DefectAt(defects, index)
}
