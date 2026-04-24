package datamanager

import "accounting-doc-processor/internal/model"

type DocumentRepository interface {
	GetAll() []model.Document
	Get(id int) (model.Document, bool)
	UpdateField(id int, fieldName string, value string) (string, model.Document, error)
	Approve(id int) error
	Add(doc model.Document) int
	Validate(id int) []string
	Replace(id int, doc model.Document) error
}
