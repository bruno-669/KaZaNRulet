package aiextractor

import "accounting-doc-processor/internal/model"

// ExtractorClient defines the interface for an AI field extraction service.
type ExtractorClient interface {
	Extract(filePath string) ([]model.Field, error)
}
