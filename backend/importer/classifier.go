package importer

import (
	"crypto/sha256"
	"database/sql"
	"fmt"

	"log"
	"strings"

	"github.com/lib/pq"
)

type DBCategory struct {
	ID       int
	Name     string
	Keywords []string
}

func GetManualCategory(narrative string, availableCategories map[string]DBCategory) int {
	log.Printf("[INPUT] No category found for narrative: '%s'\n", narrative)
	var category string
	fmt.Scan(&category)
	for _, c := range availableCategories {
		if strings.EqualFold(c.Name, category) {
			log.Printf("[INPUT] Category '%s' selected with ID %d", c.Name, c.ID)
			return c.ID
		}
	}
	GetManualCategory(narrative, availableCategories)

	return 1
}

func GetCategoryID(narrative string, availableCategories map[string]DBCategory, paymentNr string) (id int, keyword string, err error) {
	narrativeLower := strings.ToLower(narrative)
	paymentNrLower := strings.ToLower(paymentNr)
	for _, c := range availableCategories {
		for _, kw := range c.Keywords {
			if strings.Contains(narrativeLower, kw) || strings.Contains(paymentNrLower, kw) {
				if c.ID == 1 {
					return 0, "", fmt.Errorf("Skipped category encountered")
				}
				return c.ID, kw, nil
			}
		}
	}

	categoryID := GetManualCategory(narrative, availableCategories)
	return categoryID, "", nil
}

func BuildKeywordIndex(cats map[int]DBCategory) (index map[string]int) {
	index = make(map[string]int)
	for id, c := range cats {
		for _, kw := range c.Keywords {
			index[strings.ToLower(kw)] = id
		}
	}
	return index
}

func LoadDBCategories(db *sql.DB) (map[string]DBCategory, error) {
	rows, err := db.Query(`SELECT id, name, keywords FROM categories`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make(map[string]DBCategory)
	for rows.Next() {
		var c DBCategory
		var keywords []string
		if err := rows.Scan(&c.ID, &c.Name, pq.Array(&keywords)); err != nil {
			return nil, err
		}
		c.Keywords = keywords
		categories[c.Name] = c
	}
	return categories, nil
}

func createFingerprint(date string, amount int, narrative string) string {
	raw := fmt.Sprintf("%s_%d_%s", date, amount, narrative[:min(50, len(narrative))])
	hash := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(raw))))
	return fmt.Sprintf("%x", hash)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
