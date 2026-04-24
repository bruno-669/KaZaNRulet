// cmd/server/main.go
package main

import (
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"accounting-doc-processor/internal/handler"
	"accounting-doc-processor/internal/pkg/logger"
	"accounting-doc-processor/internal/service/datamanager"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapped, r)
		slog.Info("Request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrapped.statusCode,
		)
	})
}

func main() {
	debug := flag.Bool("debug", false, "enable debug mode")
	flag.Parse()

	cleanup, err := logger.InitLogger(*debug)
	if err != nil {
		panic("ERROR logger" + err.Error())
	}
	defer cleanup()

	slog.Info("Start program")

	// Инициализируем репозиторий документов и внедряем в обработчики
	repo := datamanager.NewMemoryDocRepo()
	handler.SetDocumentRepo(repo)

	// Парсим шаблоны
	handler.TmplUpload = template.Must(template.ParseFiles(
		"web/templates/base.html",
		"web/templates/upload_page.html",
	))
	handler.TmplVerify = template.Must(template.ParseFiles(
		"web/templates/base.html",
		"web/templates/verification_page.html",
	))
	handler.TmplResults = template.Must(template.ParseFiles(
		"web/templates/base.html",
		"web/templates/results_page.html",
	))
	handler.TmplPartial = template.Must(template.ParseFiles(
		"web/templates/file_list_partial.html",
	))

	mux := http.NewServeMux()

	// Page handlers
	mux.HandleFunc("/", handler.UploadPageHandler)
	mux.HandleFunc("GET /verify/", handler.VerificationPageHandler)
	mux.HandleFunc("PUT /verify/", handler.UpdateFieldHandler)
	mux.HandleFunc("POST /verify/", handler.ApproveHandler)
	mux.HandleFunc("/results/", handler.ResultsPageHandler)
	mux.HandleFunc("GET /export/", handler.ExportHandler)
	mux.HandleFunc("GET /files/", handler.FileHandler)
	// HTMX partials & API
	mux.HandleFunc("GET /partials/file-list", handler.FileListPartialHandler)
	mux.HandleFunc("POST /upload", handler.UploadHandler)

	// Static files
	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	server := loggingMiddleware(mux)

	slog.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", server); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
