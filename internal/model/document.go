package model

// FieldType определяет тип вводимых данных для поля.
type FieldType string

const (
	FieldText   FieldType = "text"   // произвольный текст
	FieldInn    FieldType = "inn"    // ИНН (10 или 12 цифр)
	FieldDate   FieldType = "date"   // дата в формате YYYY-MM-DD
	FieldNumber FieldType = "number" // целое или дробное число
	FieldAmount FieldType = "amount" // денежная сумма (аналогично number)
)

// Field представляет одно извлекаемое поле документа.
type Field struct {
	Name  string    `json:"name"`            // внутреннее имя (например "supplier_inn")
	Label string    `json:"label"`           // человекочитаемое название
	Value string    `json:"value"`           // текущее значение
	Type  FieldType `json:"type"`            // тип поля
	Valid bool      `json:"valid"`           // флаг валидности
	Error string    `json:"error,omitempty"` // текст ошибки валидации
}

// Document представляет бухгалтерский документ с динамическим набором полей.
type Document struct {
	ID       int     `json:"id"`
	FileName string  `json:"file_name"`
	FileExt  string  `json:"file_ext"` // расширение файла: pdf, jpg и т.д.
	Status   string  `json:"status"`   // processing, ready, approved
	Fields   []Field `json:"fields"`   // список распознанных полей
}
