package datamanager

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"sync"
	"time"

	"accounting-doc-processor/internal/model"
)

type MemoryDocRepo struct {
	mu   sync.Mutex
	docs map[int]model.Document
}

func NewMemoryDocRepo() *MemoryDocRepo {
	repo := &MemoryDocRepo{
		docs: make(map[int]model.Document),
	}
	repo.initTestData()
	return repo
}

// initTestData заполняет хранилище начальными тестовыми документами.
func (r *MemoryDocRepo) initTestData() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.docs[1] = model.Document{
		ID:       1,
		FileName: "contract_001.pdf",
		FileExt:  "pdf",
		Status:   "processing",
		Fields: []model.Field{
			{Name: "number", Label: "Номер документа", Value: "NK-2026/001", Type: model.FieldText, Valid: true},
			{Name: "date", Label: "Дата", Value: "2026-04-24", Type: model.FieldDate, Valid: true},
			{Name: "supplier", Label: "Поставщик", Value: "ООО Ромашка", Type: model.FieldText, Valid: true},
			{Name: "supplier_inn", Label: "ИНН поставщика", Value: "1234ABCD", Type: model.FieldInn, Valid: false, Error: "ИНН должен содержать 10 или 12 цифр"},
			{Name: "buyer", Label: "Покупатель", Value: "ООО Покупатель", Type: model.FieldText, Valid: true},
			{Name: "buyer_inn", Label: "ИНН покупателя", Value: "7707083893", Type: model.FieldInn, Valid: true},
			{Name: "item_name", Label: "Наименование товара", Value: "Канцтовары", Type: model.FieldText, Valid: true},
			{Name: "quantity", Label: "Количество", Value: "10", Type: model.FieldNumber, Valid: true},
			{Name: "price", Label: "Цена", Value: "150.00", Type: model.FieldAmount, Valid: true},
			{Name: "total_sum", Label: "Общая сумма", Value: "1500.00", Type: model.FieldAmount, Valid: true},
		},
	}
	r.docs[2] = model.Document{
		ID: 2, FileName: "invoice_042.pdf", FileExt: "pdf", Status: "ready",
		Fields: []model.Field{
			{Name: "number", Label: "Номер документа", Value: "SF-2026/042", Type: model.FieldText, Valid: true},
			{Name: "date", Label: "Дата", Value: "2026-04-20", Type: model.FieldDate, Valid: true},
			{Name: "supplier", Label: "Поставщик", Value: "ООО Поставщик", Type: model.FieldText, Valid: true},
			{Name: "supplier_inn", Label: "ИНН поставщика", Value: "7707083893", Type: model.FieldInn, Valid: true},
			{Name: "buyer", Label: "Покупатель", Value: "ООО Заказчик", Type: model.FieldText, Valid: true},
			{Name: "buyer_inn", Label: "ИНН покупателя", Value: "7707083893", Type: model.FieldInn, Valid: true},
			{Name: "item_name", Label: "Наименование товара", Value: "Бумага А4", Type: model.FieldText, Valid: true},
			{Name: "quantity", Label: "Количество", Value: "100", Type: model.FieldNumber, Valid: true},
			{Name: "price", Label: "Цена", Value: "2.50", Type: model.FieldAmount, Valid: true},
			{Name: "total_sum", Label: "Общая сумма", Value: "250.00", Type: model.FieldAmount, Valid: true},
		},
	}
	r.docs[3] = model.Document{
		ID: 3, FileName: "act_007.pdf", FileExt: "pdf", Status: "approved",
		Fields: []model.Field{
			{Name: "number", Label: "Номер документа", Value: "ACT-2026/007", Type: model.FieldText, Valid: true},
			{Name: "date", Label: "Дата", Value: "2026-04-15", Type: model.FieldDate, Valid: true},
			{Name: "supplier", Label: "Поставщик", Value: "ИП Иванов", Type: model.FieldText, Valid: true},
			{Name: "supplier_inn", Label: "ИНН поставщика", Value: "123456789012", Type: model.FieldInn, Valid: true},
			{Name: "buyer", Label: "Покупатель", Value: "ООО Клиент", Type: model.FieldText, Valid: true},
			{Name: "buyer_inn", Label: "ИНН покупателя", Value: "7707083893", Type: model.FieldInn, Valid: true},
			{Name: "item_name", Label: "Наименование товара", Value: "Услуги", Type: model.FieldText, Valid: true},
			{Name: "quantity", Label: "Количество", Value: "1", Type: model.FieldNumber, Valid: true},
			{Name: "price", Label: "Цена", Value: "5000.00", Type: model.FieldAmount, Valid: true},
			{Name: "total_sum", Label: "Общая сумма", Value: "5000.00", Type: model.FieldAmount, Valid: true},
		},
	}
	slog.Info("Initialised test documents in memory store", "count", len(r.docs))
}

func (r *MemoryDocRepo) GetAll() []model.Document {
	r.mu.Lock()
	defer r.mu.Unlock()
	res := make([]model.Document, 0, len(r.docs))
	for _, d := range r.docs {
		res = append(res, d)
	}
	return res
}

func (r *MemoryDocRepo) Get(id int) (model.Document, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	doc, ok := r.docs[id]
	return doc, ok
}

func (r *MemoryDocRepo) UpdateField(id int, fieldName string, value string) (string, model.Document, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	doc, ok := r.docs[id]
	if !ok {
		return "", model.Document{}, fmt.Errorf("document %d not found", id)
	}
	found := false
	for i := range doc.Fields {
		if doc.Fields[i].Name == fieldName {
			doc.Fields[i].Value = value
			errMsg := validateField(doc.Fields[i].Type, value)
			if errMsg != "" {
				doc.Fields[i].Valid = false
				doc.Fields[i].Error = errMsg
			} else {
				doc.Fields[i].Valid = true
				doc.Fields[i].Error = ""
			}
			found = true
			break
		}
	}
	if !found {
		return fmt.Sprintf("field '%s' not found", fieldName), doc, nil
	}
	r.docs[id] = doc
	for _, f := range doc.Fields {
		if f.Name == fieldName {
			return f.Error, doc, nil
		}
	}
	return "", doc, nil
}

func (r *MemoryDocRepo) Approve(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	doc, ok := r.docs[id]
	if !ok {
		return fmt.Errorf("document %d not found", id)
	}
	doc.Status = "approved"
	r.docs[id] = doc
	slog.Info("Document approved", "id", id)
	return nil
}

func (r *MemoryDocRepo) Add(doc model.Document) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	newID := 1
	for id := range r.docs {
		if id >= newID {
			newID = id + 1
		}
	}
	doc.ID = newID
	if doc.Status == "" {
		doc.Status = "ready"
	}
	if doc.FileExt == "" && doc.FileName != "" {
		ext := filepath.Ext(doc.FileName)
		if ext != "" {
			doc.FileExt = ext[1:]
		} else {
			doc.FileExt = "unknown"
		}
	}
	r.docs[newID] = doc
	slog.Info("Document added", "id", newID, "status", doc.Status)
	return newID
}

func (r *MemoryDocRepo) Validate(id int) []string {
	doc, ok := r.Get(id)
	if !ok {
		return []string{"Документ не найден"}
	}
	var errs []string
	for _, f := range doc.Fields {
		if msg := validateField(f.Type, f.Value); msg != "" {
			errs = append(errs, fmt.Sprintf("%s: %s", f.Label, msg))
		}
	}
	return errs
}

func (r *MemoryDocRepo) Replace(id int, doc model.Document) error {
	if doc.ID != id {
		return fmt.Errorf("document ID mismatch")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.docs[id]; !ok {
		return fmt.Errorf("document %d not found", id)
	}
	r.docs[id] = doc
	return nil
}

// Вспомогательные функции валидации (оставлены внутри пакета)
func validateField(ft model.FieldType, value string) string {
	if value == "" {
		return ""
	}
	switch ft {
	case model.FieldInn:
		if !isValidINN(value) {
			return "ИНН должен содержать 10 или 12 цифр"
		}
	case model.FieldDate:
		if !isValidDate(value) {
			return "Дата должна быть в формате YYYY-MM-DD"
		}
	case model.FieldNumber, model.FieldAmount:
		f := parseFloat(value)
		if f < 0 {
			return "Значение не может быть отрицательным"
		}
		if _, err := fmt.Sscanf(value, "%f", &f); err != nil {
			return "Некорректное числовое значение"
		}
	}
	return ""
}

func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
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
