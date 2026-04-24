package aiextractor

import (
	"strings"

	"accounting-doc-processor/internal/model"
)

// LocalClient is a stub implementation that returns hardcoded fields
// based on the file name.
type LocalClient struct{}

func (c *LocalClient) Extract(filePath string) ([]model.Field, error) {
	if strings.Contains(filePath, "contract_001.pdf") {
		return []model.Field{
			{Name: "number", Label: "Номер документа", Value: "NK-2026/001", Type: model.FieldText},
			{Name: "date", Label: "Дата", Value: "2026-04-24", Type: model.FieldDate},
			{Name: "supplier", Label: "Поставщик", Value: "ООО Ромашка", Type: model.FieldText},
			{Name: "supplier_inn", Label: "ИНН поставщика", Value: "1234ABCD", Type: model.FieldInn},
			{Name: "buyer", Label: "Покупатель", Value: "ООО Покупатель", Type: model.FieldText},
			{Name: "buyer_inn", Label: "ИНН покупателя", Value: "7707083893", Type: model.FieldInn},
			{Name: "item_name", Label: "Наименование товара", Value: "Канцтовары", Type: model.FieldText},
			{Name: "quantity", Label: "Количество", Value: "10", Type: model.FieldNumber},
			{Name: "price", Label: "Цена", Value: "150.00", Type: model.FieldAmount},
			{Name: "total_sum", Label: "Общая сумма", Value: "1500.00", Type: model.FieldAmount},
		}, nil
	}
	if strings.Contains(filePath, "invoice_042.pdf") {
		return []model.Field{
			{Name: "number", Label: "Номер документа", Value: "SF-2026/042", Type: model.FieldText},
			{Name: "date", Label: "Дата", Value: "2026-04-20", Type: model.FieldDate},
			{Name: "supplier", Label: "Поставщик", Value: "ООО Поставщик", Type: model.FieldText},
			{Name: "supplier_inn", Label: "ИНН поставщика", Value: "7707083893", Type: model.FieldInn},
			{Name: "buyer", Label: "Покупатель", Value: "ООО Заказчик", Type: model.FieldText},
			{Name: "buyer_inn", Label: "ИНН покупателя", Value: "7707083893", Type: model.FieldInn},
			{Name: "item_name", Label: "Наименование товара", Value: "Бумага А4", Type: model.FieldText},
			{Name: "quantity", Label: "Количество", Value: "100", Type: model.FieldNumber},
			{Name: "price", Label: "Цена", Value: "2.50", Type: model.FieldAmount},
			{Name: "total_sum", Label: "Общая сумма", Value: "250.00", Type: model.FieldAmount},
		}, nil
	}
	// Default minimal set
	return []model.Field{
		{Name: "number", Label: "Номер документа", Value: "ДЕМО-2026/001", Type: model.FieldText},
		{Name: "date", Label: "Дата", Value: "2026-04-24", Type: model.FieldDate},
		{Name: "supplier", Label: "Поставщик", Value: "ООО Демо-Поставщик", Type: model.FieldText},
		{Name: "supplier_inn", Label: "ИНН поставщика", Value: "1234ABCD", Type: model.FieldInn}, // намеренная ошибка
		{Name: "buyer", Label: "Покупатель", Value: "ООО Демо-Покупатель", Type: model.FieldText},
		{Name: "buyer_inn", Label: "ИНН покупателя", Value: "7707083893", Type: model.FieldInn},
		{Name: "item_name", Label: "Наименование товара", Value: "Тестовый товар", Type: model.FieldText},
		{Name: "quantity", Label: "Количество", Value: "5", Type: model.FieldNumber},
		{Name: "price", Label: "Цена", Value: "100.00", Type: model.FieldAmount},
		{Name: "total_sum", Label: "Общая сумма", Value: "500.00", Type: model.FieldAmount},
	}, nil
}
