// internal/handler/verification.go

package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

// UpdateFieldHandler handles PUT /verify/{id}/fields/{fieldName}
// It updates a single field and returns an error span or empty span.
func UpdateFieldHandler(w http.ResponseWriter, r *http.Request) {
	// Путь вида /verify/123/fields/supplier_inn
	path := strings.TrimPrefix(r.URL.Path, "/verify/")
	parts := strings.Split(path, "/")
	// Ожидаем: parts[0] = id, parts[1] = "fields", parts[2] = fieldName
	if len(parts) < 3 || parts[1] != "fields" {
		http.NotFound(w, r)
		return
	}
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fieldName := parts[2]
	value := r.FormValue("value")

	errMsg, _, err := DocRepo.UpdateField(id, fieldName, value)
	if err != nil {
		// Не удалось найти документ или другая внутренняя ошибка
		slog.Error("UpdateDocumentField failed", "id", id, "field", fieldName, "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Возвращаем фрагмент с сообщением об ошибке
	w.Header().Set("Content-Type", "text/html")
	if errMsg != "" {
		fmt.Fprintf(w, `<span class="error-text" id="err-%s">%s</span>`, fieldName, errMsg)
	} else {
		fmt.Fprintf(w, `<span id="err-%s"></span>`, fieldName)
	}
}

// ApproveHandler handles POST /verify/{id}/approve
// It approves the document and redirects to the results page via HX-Redirect.
func ApproveHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/verify/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 1 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}
	id, err := strconv.Atoi(parts[0])
	if err != nil || len(parts) < 2 || parts[1] != "approve" {
		http.NotFound(w, r)
		return
	}

	// Проверка полей документа
	validationErrors := DocRepo.Validate(id)
	if len(validationErrors) > 0 {
		slog.Warn("Validation failed before approve", "id", id, "errors", validationErrors)

		// Формируем HTML со списком ошибок
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `<div class="validation-errors"><ul>`)
		for _, e := range validationErrors {
			fmt.Fprintf(w, `<li>%s</li>`, e)
		}
		fmt.Fprint(w, `</ul></div>`)
		return
	}

	// Ошибок нет — утверждаем документ
	err = DocRepo.Approve(id)
	if err != nil {
		slog.Error("ApproveDocument failed", "id", id, "error", err)
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	slog.Info("Document approved", "id", id)
	w.Header().Set("HX-Redirect", fmt.Sprintf("/results/%d", id))
	w.WriteHeader(http.StatusOK)
}
