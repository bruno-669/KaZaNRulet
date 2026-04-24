package handler

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"accounting-doc-processor/internal/model"
	"accounting-doc-processor/internal/service/aiextractor"
	"accounting-doc-processor/internal/service/documentprocessor"
)

// FileListPartialHandler returns an HTML fragment with the document table.
func FileListPartialHandler(w http.ResponseWriter, r *http.Request) {
	documents := DocRepo.GetAll()

	data := struct {
		Documents []model.Document
	}{
		Documents: documents,
	}

	// Используем TmplPartial, который содержит только file_list_partial.html
	if err := TmplPartial.ExecuteTemplate(w, "file_list_partial.html", data); err != nil {
		slog.Error("Failed to execute file list partial template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// UploadHandler processes file upload and adds the document to the store.
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 32<<20)

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		slog.Error("Failed to parse multipart form", "error", err)
		http.Error(w, "File too large or bad request", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("documents")
	if err != nil {
		slog.Error("Failed to get uploaded file", "error", err)
		http.Error(w, "Missing file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	uploadDir := "web/uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		slog.Error("Failed to create upload directory", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fileName := filepath.Base(header.Filename)
	if fileName == "." || fileName == "/" {
		fileName = "uploaded_file"
	}
	destPath := filepath.Join(uploadDir, fileName)

	dst, err := os.Create(destPath)
	if err != nil {
		slog.Error("Failed to create destination file", "error", err, "path", destPath)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		slog.Error("Failed to save uploaded file", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Создаём запись документа
	doc := model.Document{
		FileName: fileName,
		Status:   "ready", // будет изменён в процессе обработки
	}
	newID := DocRepo.Add(doc)
	slog.Info("File saved, document record created", "id", newID, "filename", fileName)

	// Запускаем локальную обработку (заглушку AI)
	storedDoc, ok := DocRepo.Get(newID)
	if !ok {
		slog.Error("Failed to retrieve just created document", "id", newID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	client := &aiextractor.LocalClient{}
	processedDoc, err := documentprocessor.ProcessDocument(storedDoc, client, destPath)
	if err != nil {
		slog.Warn("Document processed with errors, keeping result", "id", newID, "error", err)
	}

	if err := DocRepo.Replace(newID, processedDoc); err != nil {
		slog.Error("Failed to update document after processing", "id", newID, "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Возвращаем обновлённый фрагмент таблицы
	FileListPartialHandler(w, r)
}

// FileHandler serves the document file by its ID.
func FileHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/files/")
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

	doc, ok := DocRepo.Get(id)
	if !ok || doc.FileName == "" {
		http.NotFound(w, r)
		return
	}

	uploadPath := filepath.Join("web", "uploads", doc.FileName)
	if _, err := os.Stat(uploadPath); err == nil {
		slog.Info("Serving file from uploads", "id", id, "path", uploadPath)
		http.ServeFile(w, r, uploadPath)
		return
	}

	testPath := filepath.Join("web", "testdocs", doc.FileName)
	if _, err := os.Stat(testPath); err == nil {
		slog.Info("Serving file from testdocs", "id", id, "path", testPath)
		http.ServeFile(w, r, testPath)
		return
	}

	slog.Warn("File not found in uploads or testdocs", "doc_id", id, "file_name", doc.FileName)
	http.NotFound(w, r)
}
