// internal/service/validator/inn.go
package validator

func IsValidINN(inn string) bool {
	if len(inn) != 10 && len(inn) != 12 {
		return false
	}
	for _, c := range inn {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
