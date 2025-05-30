package importer

import (
	"database/sql"
	"encoding/csv"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Transaction struct {
	Date        string `json:"date"`
	Amount      int    `json:"amount"`
	CategoryID  int    `json:"category_id"`
	Fingerprint string `json:"fingerprint"`
	Keyword     string `json:"keyword,omitempty"`
}

func extractDateFromNarrative(narrative string) string {
	r := regexp.MustCompile(`\b(\d{2}/\d{2}/\d{4})\b`)
	if match := r.FindStringSubmatch(narrative); len(match) > 0 {
		return strings.ReplaceAll(match[1], "/", ".")
	}
	return ""
}

func preprocessCSV(filePath string, db *sql.DB) ([]Transaction, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := csv.NewReader(file)
	r.Comma = '|'
	r.FieldsPerRecord = -1
	r.LazyQuotes = true

	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	categories, err := LoadDBCategories(db)
	if err != nil {
		log.Printf("[ERROR] Loading categories: %v", err)
	}

	var transactions []Transaction
	processing := false

	for _, row := range records {
		if len(row) == 0 || strings.TrimSpace(row[0]) == "" {
			continue
		}

		first := strings.ToLower(strings.TrimSpace(row[0]))

		if first == "date" {
			processing = true
			continue
		}
		if !processing {
			continue
		}

		if strings.Contains(first, "closing balance") {
			log.Printf("[INFO] Reached end of transactions in %s", filePath)
			break
		}
		if len(row) < 6 {
			log.Printf("[WARN] Skipping row with insufficient columns: %+v", row)
			continue
		}

		narrative := strings.TrimSpace(row[2])
		amountStr := strings.TrimSpace(row[5])
		paymentNr := strings.TrimSpace(row[3])

		var date time.Time
		extracted := extractDateFromNarrative(narrative)
		if extracted != "" {
			date, err = time.Parse("02.01.2006", extracted)
		}
		if extracted == "" || err != nil {
			dateStr := strings.TrimSpace(row[0])
			date, err = time.Parse("02.01.2006", dateStr)
			if err != nil {
				log.Printf("[WARN] Invalid date for row: %+v", row)
				continue
			}
		}

		amountStr = strings.ReplaceAll(amountStr, ",", ".")
		amountFloat, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			log.Printf("[WARN] Skipping invalid amount '%s'", amountStr)
			continue
		}
		amount := int(amountFloat * 100)

		fingerprint := createFingerprint(date.Format("2006-01-02"), amount, narrative)
		if transactionExists(db, fingerprint) {
			continue
		}

		categoryID, keyword, err := GetCategoryID(narrative, categories, paymentNr)
		if categoryID == 0 {
			continue
		}
		if err != nil {
			log.Printf("[WARN] No valid category found for row: %+v, Error: %v", row, err)
		}

		transactions = append(transactions, Transaction{
			Date:        date.Format("2006-01-02"),
			Amount:      amount,
			CategoryID:  categoryID,
			Fingerprint: fingerprint,
			Keyword:     keyword,
		})
	}

	log.Printf("[INFO] Parsed %d transactions from %s", len(transactions), filePath)
	return transactions, nil
}
