package handler

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"accounting-doc-processor/internal/model"
)

func ExportHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/export/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 {
		http.NotFound(w, r)
		return
	}
	idStr := parts[0]
	format := parts[1]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	doc, ok := DocRepo.Get(id)
	if !ok {
		http.NotFound(w, r)
		return
	}

	if doc.Status != "approved" {
		http.Error(w, "Document not approved yet", http.StatusForbidden)
		return
	}

	slog.Info("Exporting document", "id", id, "format", format)

	switch format {
	case "json":
		exportJSON(w, doc)
	case "csv":
		exportCSV(w, doc)
	case "xml":
		exportXML(w, doc)
	case "xlsx":
		exportXLSX(w, doc)
	default:
		http.Error(w, "Unsupported export format", http.StatusBadRequest)
	}
}

func exportJSON(w http.ResponseWriter, doc model.Document) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(doc); err != nil {
		slog.Error("Failed to encode JSON", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func exportCSV(w http.ResponseWriter, doc model.Document) {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	// Заголовки из названий полей
	headers := []string{"ID", "FileName", "FileExt", "Status"}
	values := []string{
		strconv.Itoa(doc.ID),
		escapeCSV(doc.FileName),
		escapeCSV(doc.FileExt),
		escapeCSV(doc.Status),
	}
	for _, f := range doc.Fields {
		headers = append(headers, f.Name)
		values = append(values, escapeCSV(f.Value))
	}
	fmt.Fprintf(w, "%s\n", strings.Join(headers, ","))
	fmt.Fprintf(w, "%s\n", strings.Join(values, ","))
}

func exportXML(w http.ResponseWriter, doc model.Document) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	type xmlField struct {
		Name  string `xml:"name,attr"`
		Value string `xml:",chardata"`
	}
	type xmlDocument struct {
		XMLName  xml.Name   `xml:"Document"`
		ID       int        `xml:"ID"`
		FileName string     `xml:"FileName"`
		FileExt  string     `xml:"FileExt"`
		Status   string     `xml:"Status"`
		Fields   []xmlField `xml:"Fields>Field"`
	}
	fields := make([]xmlField, len(doc.Fields))
	for i, f := range doc.Fields {
		fields[i] = xmlField{Name: f.Name, Value: f.Value}
	}
	x := xmlDocument{
		ID:       doc.ID,
		FileName: doc.FileName,
		FileExt:  doc.FileExt,
		Status:   doc.Status,
		Fields:   fields,
	}
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "  ")
	fmt.Fprint(w, xml.Header)
	if err := encoder.Encode(x); err != nil {
		slog.Error("Failed to encode XML", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func exportXLSX(w http.ResponseWriter, doc model.Document) {
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	numberLabel := "без номера"
	for _, f := range doc.Fields {
		if f.Name == "number" {
			numberLabel = f.Value
			break
		}
	}
	fmt.Fprintf(w, "XLSX stub for document %s", numberLabel)
}

func escapeCSV(s string) string {
	if strings.ContainsAny(s, ",\"\n") {
		return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
	}
	return s
}
