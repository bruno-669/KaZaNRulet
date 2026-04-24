// internal/service/validator/numeric.go
package validator

import "fmt"

func IsValidNumber(value string) bool {
	if value == "" {
		return true
	}
	var f float64
	if _, err := fmt.Sscanf(value, "%f", &f); err != nil {
		return false
	}
	return f >= 0
}
