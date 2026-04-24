package model

// Document represents an extracted accounting document with fields
// that can be displayed on verification and results pages.
type Document struct {
	ID          int     `json:"id"`
	FileName    string  `json:"file_name"`
	FileExt     string  `json:"file_ext"` // расширение файла: "pdf", "jpg", "png" и т.д.
	Status      string  `json:"status"`   // "processing", "ready", "approved"
	Number      string  `json:"number"`
	Date        string  `json:"date"`
	Supplier    string  `json:"supplier"`
	SupplierINN string  `json:"supplier_inn"`
	Buyer       string  `json:"buyer"`
	BuyerINN    string  `json:"buyer_inn"`
	ItemName    string  `json:"item_name"`
	Quantity    float64 `json:"quantity"`
	Price       float64 `json:"price"`
	TotalSum    float64 `json:"total_sum"`
}

type Field struct {
	Name    string
	Value   string
	IsValid bool
	Error   string
}
