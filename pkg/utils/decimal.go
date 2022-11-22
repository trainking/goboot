package utils

import (
	"fmt"
	"strconv"
)

// decimalTemp 浮点数地模板
func decimalTemp(spec int) string {
	return "%." + strconv.Itoa(spec) + "f"
}

// Decimal 浮点数保留小数位数转换
func Decimal(value float64, spec int) float64 {
	valueStr := DecimalString(value, spec)
	value, _ = strconv.ParseFloat(valueStr, 64)
	return value
}

// DecimalString 小数位数转换成保留位数字符串
func DecimalString(value float64, spec int) string {
	return fmt.Sprintf(decimalTemp(spec), value)
}
