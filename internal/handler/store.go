package handler

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"sync"
	"time"

	"accounting-doc-processor/internal/model"
)

// Глобальное потокобезопасное хранилище документов
var docStore = struct {
	mu   sync.Mutex
	docs map[int]model.Document
}{
	docs: make(map[int]model.Document),
}

func init() {
	initTestData()
}

// initTestData добавляет стартовые документы с полными демонстрационными данными.
func initTestData() {
	docStore.mu.Lock()
	defer docStore.mu.Unlock()

	docStore.docs[1] = model.Document{
		ID:          1,
		FileName:    "contract_001.pdf",
		FileExt:     "pdf",
		Status:      "processing",
		Number:      "NK-2026/001",
		Date:        "2026-04-24",
		Supplier:    "ООО Ромашка",
		SupplierINN: "1234ABCD", // намеренно некорректный для демонстрации ошибки валидации
		Buyer:       "ООО Покупатель",
		BuyerINN:    "7707083893",
		ItemName:    "Канцтовары",
		Quantity:    10.0,
		Price:       150.00,
		TotalSum:    1500.00,
	}
	slog.Info("Added test document", "id", 1, "number", "NK-2026/001")

	docStore.docs[2] = model.Document{
		ID:          2,
		FileName:    "invoice_042.pdf",
		FileExt:     "pdf",
		Status:      "ready",
		Number:      "SF-2026/042",
		Date:        "2026-04-20",
		Supplier:    "ООО Поставщик",
		SupplierINN: "7707083893",
		Buyer:       "ООО Заказчик",
		BuyerINN:    "7707083893",
		ItemName:    "Бумага А4",
		Quantity:    100.0,
		Price:       2.50,
		TotalSum:    250.00,
	}
	slog.Info("Added test document", "id", 2, "number", "SF-2026/042")

	docStore.docs[3] = model.Document{
		ID:          3,
		FileName:    "act_007.pdf",
		FileExt:     "pdf",
		Status:      "approved",
		Number:      "ACT-2026/007",
		Date:        "2026-04-15",
		Supplier:    "ИП Иванов",
		SupplierINN: "123456789012",
		Buyer:       "ООО Клиент",
		BuyerINN:    "7707083893",
		ItemName:    "Услуги",
		Quantity:    1.0,
		Price:       5000.00,
		TotalSum:    5000.00,
	}
	slog.Info("Added test document", "id", 3, "number", "ACT-2026/007")

	slog.Info("Initialised test documents in store", "count", len(docStore.docs))
}

// GetAllDocuments возвращает срез всех документов (порядок не гарантирован)
func GetAllDocuments() []model.Document {
	docStore.mu.Lock()
	defer docStore.mu.Unlock()

	result := make([]model.Document, 0, len(docStore.docs))
	for _, d := range docStore.docs {
		result = append(result, d)
	}
	return result
}

// GetDocument возвращает документ по ID (и флаг существования)
func GetDocument(id int) (model.Document, bool) {
	docStore.mu.Lock()
	defer docStore.mu.Unlock()

	d, ok := docStore.docs[id]
	return d, ok
}

// UpdateDocumentField обновляет поле документа и возвращает ошибку валидации (пустая строка, если ОК)
// и обновлённый документ. Если документ не найден, возвращает ошибку.
func UpdateDocumentField(id int, fieldName string, value string) (string, model.Document, error) {
	docStore.mu.Lock()
	defer docStore.mu.Unlock()

	doc, ok := docStore.docs[id]
	if !ok {
		return "", model.Document{}, fmt.Errorf("document with id %d not found", id)
	}

	// Обновляем поле (только те, что есть в модели)
	switch fieldName {
	case "number":
		doc.Number = value
	case "date":
		doc.Date = value
	case "supplier":
		doc.Supplier = value
	case "supplier_inn":
		doc.SupplierINN = value
	case "buyer":
		doc.Buyer = value
	case "buyer_inn":
		doc.BuyerINN = value
	case "item_name":
		doc.ItemName = value
	case "quantity":
		doc.Quantity = parseFloat(value)
	case "price":
		doc.Price = parseFloat(value)
	case "total_sum":
		doc.TotalSum = parseFloat(value)
	default:
		return fmt.Sprintf("unknown field: %s", fieldName), doc, nil
	}

	docStore.docs[id] = doc

	// Простая валидация поля
	validationErr := validateField(fieldName, value)
	return validationErr, doc, nil
}

// ApproveDocument устанавливает статус "approved"
func ApproveDocument(id int) error {
	docStore.mu.Lock()
	defer docStore.mu.Unlock()

	doc, ok := docStore.docs[id]
	if !ok {
		return fmt.Errorf("document with id %d not found", id)
	}
	doc.Status = "approved"
	docStore.docs[id] = doc
	return nil
}

// AddDocument добавляет новый документ в хранилище и возвращает его ID
func AddDocument(doc model.Document) int {
	docStore.mu.Lock()
	defer docStore.mu.Unlock()

	// Генерируем простой ID: на 1 больше максимального существующего
	newID := 1
	for id := range docStore.docs {
		if id >= newID {
			newID = id + 1
		}
	}
	doc.ID = newID

	// Если FileExt не задан, определяем из имени файла
	if doc.FileExt == "" && doc.FileName != "" {
		ext := filepath.Ext(doc.FileName)
		if ext != "" {
			doc.FileExt = ext[1:] // убираем точку
		} else {
			doc.FileExt = "unknown"
		}
	}

	docStore.docs[newID] = doc
	return newID
}

// --- Вспомогательные функции ---

func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

func validateField(fieldName, value string) string {
	switch fieldName {
	case "supplier_inn", "buyer_inn":
		if !isValidINN(value) {
			return "ИНН должен содержать 10 или 12 цифр"
		}
	case "date":
		if !isValidDate(value) {
			return "Дата должна быть в формате YYYY-MM-DD"
		}
	case "quantity", "price", "total_sum":
		f := parseFloat(value)
		if f < 0 {
			return "Значение не может быть отрицательным"
		}
	}
	return ""
}

func isValidINN(inn string) bool {
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

func isValidDate(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

func ValidateDocument(id int) []string {
	doc, ok := GetDocument(id)
	if !ok {
		return []string{"Документ не найден"}
	}

	var errors []string

	// Поля для проверки: имя поля в модели, текущее значение, русское название
	checks := []struct {
		fieldName string
		value     string
		label     string
	}{
		{"supplier_inn", doc.SupplierINN, "ИНН поставщика"},
		{"buyer_inn", doc.BuyerINN, "ИНН покупателя"},
		{"date", doc.Date, "Дата"},
		{"quantity", fmt.Sprint(doc.Quantity), "Количество"},
		{"price", fmt.Sprint(doc.Price), "Цена"},
		{"total_sum", fmt.Sprint(doc.TotalSum), "Общая сумма"},
	}

	for _, c := range checks {
		if msg := validateField(c.fieldName, c.value); msg != "" {
			errors = append(errors, fmt.Sprintf("%s: %s", c.label, msg))
		}
	}

	return errors
}
