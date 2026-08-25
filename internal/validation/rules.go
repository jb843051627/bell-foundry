package validation

import (
	"fmt"
	"math"
	"strings"
)

// FieldError 包含可直接反馈给操作员的字段错误。
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e FieldError) Error() string { return fmt.Sprintf("%s: %s", e.Field, e.Message) }

// Required 检查字符串字段非空并去除首尾空白。
func Required(field, value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", FieldError{Field: field, Message: "required"}
	}
	return value, nil
}

// Positive 检查正数且不是 NaN/Inf。
func Positive(field string, value float64) error {
	if value <= 0 || math.IsNaN(value) || math.IsInf(value, 0) {
		return FieldError{Field: field, Message: "must be positive and finite"}
	}
	return nil
}

// Range 检查数值上下限。
func Range(field string, value, lower, upper float64) error {
	if value < lower || value > upper || math.IsNaN(value) {
		return FieldError{Field: field, Message: fmt.Sprintf("must be between %.2f and %.2f", lower, upper)}
	}
	return nil
}

// OneOf 检查枚举字段。
func OneOf(field, value string, allowed ...string) error {
	for _, item := range allowed {
		if value == item {
			return nil
		}
	}
	return FieldError{Field: field, Message: "unsupported value"}
}

// All 聚合多个字段错误。
func All(errs ...error) error {
	messages := make([]string, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			messages = append(messages, err.Error())
		}
	}
	if len(messages) == 0 {
		return nil
	}
	return fmt.Errorf("validation: %s", strings.Join(messages, "; "))
}

// Clamp 将数值限制在工艺范围内。
func Clamp(value, lower, upper float64) float64 {
	if value < lower {
		return lower
	}
	if value > upper {
		return upper
	}
	return value
}
