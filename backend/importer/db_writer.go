package importer

import (
	"database/sql"
	"log"
)

func transactionExists(db *sql.DB, fingerprint string) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM transactions WHERE fingerprint = $1)", fingerprint).Scan(&exists)
	if err != nil {
		log.Printf("[ERROR] Checking fingerprint: %v", err)
		return false
	}
	return exists
}

func insertTransaction(db *sql.DB, row Transaction) error {
	log.Printf("[INFO] Inserting transaction: %s, CategoryID: %d, Date: %s, Amount: %d, Keyword: %s", row.Fingerprint, row.CategoryID, row.Date, row.Amount, row.Keyword)
	_, err := db.Exec(
		"INSERT INTO transactions (fingerprint, category_id, date, amount, keyword) VALUES ($1, $2, $3, $4, $5)",
		row.Fingerprint, row.CategoryID, row.Date, row.Amount, row.Keyword,
	)
	return err
}
