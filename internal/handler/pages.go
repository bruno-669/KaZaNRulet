// internal/handler/pages.go
package handler

import (
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"accounting-doc-processor/internal/model"
)

// Раздельные шаблоны для каждой страницы, чтобы избежать конфликта блоков "content".
var (
	TmplUpload  *template.Template
	TmplVerify  *template.Template
	TmplResults *template.Template
	TmplPartial *template.Template
)

// UploadPageHandler renders the upload page with a list of example documents.
func UploadPageHandler(w http.ResponseWriter, r *http.Request) {
	documents := GetAllDocuments()

	pageData := struct {
		Documents []model.Document
	}{
		Documents: documents,
	}

	if err := TmplUpload.ExecuteTemplate(w, "base.html", pageData); err != nil {
		slog.Error("Failed to execute upload template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// VerificationPageHandler renders the verification page for a document.
func VerificationPageHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/verify/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 1 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}
	idStr := parts[0]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	doc, ok := GetDocument(id)
	if !ok {
		http.NotFound(w, r)
		return
	}

	slog.Info("Rendering verification page", "id", id)

	data := struct {
		Document model.Document
	}{
		Document: doc,
	}

	if err := TmplVerify.ExecuteTemplate(w, "base.html", data); err != nil {
		slog.Error("Failed to execute verification template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// ResultsPageHandler renders the results page for an approved document.
func ResultsPageHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/results/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 1 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	doc, ok := GetDocument(id)
	if !ok {
		http.NotFound(w, r)
		return
	}

	slog.Info("Rendering results page", "id", id)

	data := struct {
		Document model.Document
		Message  string
	}{
		Document: doc,
		Message:  "Документ успешно обработан и утверждён.",
	}

	if err := TmplResults.ExecuteTemplate(w, "base.html", data); err != nil {
		slog.Error("Failed to execute results template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
