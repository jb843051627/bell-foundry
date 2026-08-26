package model

import "fmt"

// ErrInvalidField 构造字段校验错误。
func ErrInvalidField(field string) error {
	return fmt.Errorf("model: invalid field %q", field)
}

// 领域级哨兵错误。
var (
	ErrCompositionOverflow = fmt.Errorf("model: copper+tin exceeds 100%%")
	ErrTinBelowFloor       = fmt.Errorf("model: tin below bell-grade floor")
	ErrPhasePointsInverted = fmt.Errorf("model: liquidus must exceed solidus")
	ErrNotFound            = fmt.Errorf("model: record not found")
	ErrAlreadyExists       = fmt.Errorf("model: record already exists")
	ErrBadTransition       = fmt.Errorf("model: illegal state transition")
	ErrPreconditionFailed  = fmt.Errorf("model: precondition failed")
	ErrNotFullyMeasured    = fmt.Errorf("model: partials not fully measured")
	ErrCurveIncomplete     = fmt.Errorf("model: cooling curve incomplete")

	// 落库/投递失败哨兵，供值班接口用 errors.Is 区分数据库失败与通知失败。
	ErrAlertPersistence      = fmt.Errorf("model: alert persistence failed")
	ErrAlertDelivery         = fmt.Errorf("model: alert delivery failed")
	ErrInspectionPersistence = fmt.Errorf("model: inspection persistence failed")
)

// Wrapf 在领域错误上附加上下文，保留 %w 错误链。
func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %v", fmt.Sprintf(format, args...), err)
}
