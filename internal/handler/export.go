// internal/handler/export.go
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

// ExportHandler handles GET /export/{id}/{format}
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

	doc, ok := GetDocument(id)
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
	fmt.Fprintf(w, "ID,FileName,FileExt,Status,Number,Date,Supplier,SupplierINN,Buyer,BuyerINN,ItemName,Quantity,Price,TotalSum\n")
	fmt.Fprintf(w, "%d,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%.2f,%.2f,%.2f\n",
		doc.ID,
		escapeCSV(doc.FileName),
		escapeCSV(doc.FileExt),
		escapeCSV(doc.Status),
		escapeCSV(doc.Number),
		escapeCSV(doc.Date),
		escapeCSV(doc.Supplier),
		escapeCSV(doc.SupplierINN),
		escapeCSV(doc.Buyer),
		escapeCSV(doc.BuyerINN),
		escapeCSV(doc.ItemName),
		doc.Quantity,
		doc.Price,
		doc.TotalSum,
	)
}

func escapeCSV(s string) string {
	if strings.ContainsAny(s, ",\"\n") {
		return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
	}
	return s
}

func exportXML(w http.ResponseWriter, doc model.Document) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	type xmlDocument struct {
		XMLName     xml.Name `xml:"Document"`
		ID          int      `xml:"ID"`
		FileName    string   `xml:"FileName"`
		FileExt     string   `xml:"FileExt"`
		Status      string   `xml:"Status"`
		Number      string   `xml:"Number"`
		Date        string   `xml:"Date"`
		Supplier    string   `xml:"Supplier"`
		SupplierINN string   `xml:"SupplierINN"`
		Buyer       string   `xml:"Buyer"`
		BuyerINN    string   `xml:"BuyerINN"`
		ItemName    string   `xml:"ItemName"`
		Quantity    float64  `xml:"Quantity"`
		Price       float64  `xml:"Price"`
		TotalSum    float64  `xml:"TotalSum"`
	}
	x := xmlDocument{
		ID:          doc.ID,
		FileName:    doc.FileName,
		FileExt:     doc.FileExt,
		Status:      doc.Status,
		Number:      doc.Number,
		Date:        doc.Date,
		Supplier:    doc.Supplier,
		SupplierINN: doc.SupplierINN,
		Buyer:       doc.Buyer,
		BuyerINN:    doc.BuyerINN,
		ItemName:    doc.ItemName,
		Quantity:    doc.Quantity,
		Price:       doc.Price,
		TotalSum:    doc.TotalSum,
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
	fmt.Fprintf(w, "XLSX stub for document %s", doc.Number)
}
