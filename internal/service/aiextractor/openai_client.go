package aiextractor
// internal/service/aiextractor/openai_client.go
package aiextractor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"accounting-doc-processor/internal/model"
)

type GigaChatExtractor struct {
	OCRURL      string
	GigaChatURL string
	Token       string
	httpClient  *http.Client
}

// Промпт из gigachat.py (адаптирован)
const promptTemplate = `Извлеки из текста бухгалтерского документа следующие поля в формате JSON:
{
  "номер_документа": "...",
  "дата": "ДД.ММ.ГГГГ",
  "поставщик_название": "...",
  "поставщик_инн": "...",
  "покупатель_название": "...",
  "покупатель_инн": "...",
  "товары": [{"наименование": "...", "количество": число, "цена": число, "сумма": число}],
  "сумма_итого": число
}
Если какого-то поля нет – не включай его.
Текст документа:
%s`

type gigaChatResponse struct {
	НомерДокумента   *string  `json:"номер_документа"`
	Дата            *string  `json:"дата"`
	ПоставщикНазв   *string  `json:"поставщик_название"`
	ПоставщикИНН    *string  `json:"поставщик_инн"`
	ПокупательНазв  *string  `json:"покупатель_название"`
	ПокупательИНН   *string  `json:"покупатель_инн"`
	Товары          []struct {
		Наименование string  `json:"наименование"`
		Количество   float64 `json:"количество"`
		Цена         float64 `json:"цена"`
		Сумма        float64 `json:"сумма"`
	} `json:"товары"`
	СуммаИтого *float64 `json:"сумма_итого"`
}

func (g *GigaChatExtractor) Extract(filePath string) ([]model.Field, error) {
	// 1. OCR
	rawText, err := g.ocrExtract(filePath)
	if err != nil {
		return nil, fmt.Errorf("ocr failed: %w", err)
	}

	// 2. GigaChat
	fields, err := g.gigaChatExtract(rawText)
	if err != nil {
		return nil, fmt.Errorf("gigachat failed: %w", err)
	}
	return fields, nil
}

func (g *GigaChatExtractor) ocrExtract(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(part, file); err != nil {
		return "", err
	}
	writer.Close()

	req, err := http.NewRequest("POST", g.OCRURL, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ocr returned status %d", resp.StatusCode)
	}

	var ocrResult struct {
		RawText string `json:"raw_text"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&ocrResult); err != nil {
		return "", err
	}
	return ocrResult.RawText, nil
}

func (g *GigaChatExtractor) gigaChatExtract(rawText string) ([]model.Field, error) {
	prompt := fmt.Sprintf(promptTemplate, rawText)
	payload := map[string]interface{}{
		"model": "GigaChat",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"temperature": 0,
	}
	payloadBytes, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", g.GigaChatURL, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+g.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gigachat returned status %d", resp.StatusCode)
	}

	var apiResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}
	if len(apiResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in gigachat response")
	}

	content := apiResp.Choices[0].Message.Content
	// Вырезаем JSON (может быть внутри markdown-блока)
	content = extractJSON(content)

	var gcFields gigaChatResponse
	if err := json.Unmarshal([]byte(content), &gcFields); err != nil {
		return nil, fmt.Errorf("failed to parse gigachat JSON: %w, content: %s", err, content)
	}

	return mapToFields(gcFields), nil
}

func extractJSON(s string) string {
	start := strings.Index(s, "{")
	if start == -1 {
		return s
	}
	end := strings.LastIndex(s, "}")
	if end == -1 || end <= start {
		return s
	}
	return s[start : end+1]
}

func mapToFields(gc gigaChatResponse) []model.Field {
	var fields []model.Field

	if gc.НомерДокумента != nil {
		fields = append(fields, model.Field{
			Name: "number", Label: "Номер документа",
			Value: *gc.НомерДокумента, Type: model.FieldText,
		})
	}
	if gc.Дата != nil {
		dateVal := normalizeDate(*gc.Дата)
		fields = append(fields, model.Field{
			Name: "date", Label: "Дата",
			Value: dateVal, Type: model.FieldDate,
		})
	}
	if gc.ПоставщикНазв != nil {
		fields = append(fields, model.Field{
			Name: "supplier", Label: "Поставщик",
			Value: *gc.ПоставщикНазв, Type: model.FieldText,
		})
	}
	if gc.ПоставщикИНН != nil {
		fields = append(fields, model.Field{
			Name: "supplier_inn", Label: "ИНН поставщика",
			Value: *gc.ПоставщикИНН, Type: model.FieldInn,
		})
	}
	if gc.ПокупательНазв != nil {
		fields = append(fields, model.Field{
			Name: "buyer", Label: "Покупатель",
			Value: *gc.ПокупательНазв, Type: model.FieldText,
		})
	}
	if gc.ПокупательИНН != nil {
		fields = append(fields, model.Field{
			Name: "buyer_inn", Label: "ИНН покупателя",
			Value: *gc.ПокупательИНН, Type: model.FieldInn,
		})
	}
	if len(gc.Товары) > 0 {
		var names, qty, price, sum []string
		for _, t := range gc.Товары {
			names = append(names, t.Наименование)
			qty = append(qty, fmt.Sprintf("%.2f", t.Количество))
			price = append(price, fmt.Sprintf("%.2f", t.Цена))
			sum = append(sum, fmt.Sprintf("%.2f", t.Сумма))
		}
		fields = append(fields, model.Field{
			Name: "item_name", Label: "Наименование товара",
			Value: strings.Join(names, "; "), Type: model.FieldText,
		})
		fields = append(fields, model.Field{
			Name: "quantity", Label: "Количество",
			Value: strings.Join(qty, "; "), Type: model.FieldNumber,
		})
		fields = append(fields, model.Field{
			Name: "price", Label: "Цена",
			Value: strings.Join(price, "; "), Type: model.FieldAmount,
		})
		fields = append(fields, model.Field{
			Name: "total_sum", Label: "Общая сумма",
			Value: strings.Join(sum, "; "), Type: model.FieldAmount,
		})
	} else if gc.СуммаИтого != nil {
		fields = append(fields, model.Field{
			Name: "total_sum", Label: "Общая сумма",
			Value: fmt.Sprintf("%.2f", *gc.СуммаИтого), Type: model.FieldAmount,
		})
	}
	return fields
}

func normalizeDate(d string) string {
	// Если дата в формате DD.MM.YYYY, конвертируем
	if strings.Count(d, ".") == 2 && len(d) == 10 {
		t, err := time.Parse("02.01.2006", d)
		if err == nil {
			return t.Format("2006-01-02")
		}
	}
	return d
}