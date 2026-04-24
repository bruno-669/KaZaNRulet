// internal/service/validator/validator.go
package validator

import (
	"accounting-doc-processor/internal/model"
)

func ValidateFields(fields []model.Field) []model.Field {
	for i := range fields {
		errMsg := ""
		switch fields[i].Type {
		case model.FieldInn:
			if !IsValidINN(fields[i].Value) {
				errMsg = "ИНН должен содержать 10 или 12 цифр"
			}
		case model.FieldDate:
			if !IsValidDate(fields[i].Value) {
				errMsg = "Дата должна быть в формате YYYY-MM-DD"
			}
		case model.FieldNumber, model.FieldAmount:
			if !IsValidNumber(fields[i].Value) {
				errMsg = "Некорректное числовое значение"
			}
		}
		if errMsg != "" {
			fields[i].Valid = false
			fields[i].Error = errMsg
		} else {
			fields[i].Valid = true
			fields[i].Error = ""
		}
	}
	return fields
}
