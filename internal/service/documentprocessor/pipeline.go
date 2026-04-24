// internal/service/documentprocessor/pipeline.go
package documentprocessor

import (
	"log/slog"

	"accounting-doc-processor/internal/model"
	"accounting-doc-processor/internal/service/aiextractor"
	"accounting-doc-processor/internal/service/validator"
)

func ProcessDocument(doc model.Document, client aiextractor.ExtractorClient, filePath string) (model.Document, error) {
	slog.Info("AI processing started", "id", doc.ID)
	doc.Status = "processing"

	fields, err := client.Extract(filePath)
	if err != nil {
		doc.Status = "ready" // ready with empty fields on failure
		slog.Error("AI extraction failed", "id", doc.ID, "error", err)
		return doc, err
	}

	// Валидация всех извлечённых полей
	fields = validator.ValidateFields(fields)

	doc.Fields = fields
	doc.Status = "ready"
	slog.Info("AI processing finished", "id", doc.ID, "fields_count", len(fields))
	return doc, nil
}
