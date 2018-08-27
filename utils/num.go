package utils

import (
	"fmt"
	"strconv"
)

// ParseFloat ... parse float64's string into float64
// "1.94" tobe 1.94
func ParseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

// Decimal ... 1.0903920 to 1.09
func Decimal(val float64) float64 {
	val, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", val), 64)
	return val
}
