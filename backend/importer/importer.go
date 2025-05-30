package importer

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
)

func RunCSVImport(db *sql.DB) {
	SeedCategories(db)
	files, _ := filepath.Glob("../data/*.csv")
	log.Printf("[INFO] Found %d files", len(files))

	for _, f := range files {
		transactions, err := preprocessCSV(f, db)
		if err != nil {
			log.Printf("[ERROR] Processing %s: %v", f, err)
			continue
		}
		log.Printf("[INFO] Processing %s with %d transactions", f, len(transactions))

		for _, row := range transactions {
			if !transactionExists(db, row.Fingerprint) {
				if err := insertTransaction(db, row); err != nil {
					log.Printf("[ERROR] Inserting txn: %v", err)
				}
			}
		}
	}

	rows, err := db.Query("SELECT c.name, SUM(t.amount)/100 as total FROM transactions t JOIN categories c ON c.id = t.category_id GROUP BY c.name;")
	if err != nil {
		log.Printf("[ERROR] Querying summary: %v", err)
		return
	}
	defer rows.Close()

	fmt.Println("Type\tCategory\tTotal")
	for rows.Next() {
        var cat string
        var total float64
        if err := rows.Scan(&cat, &total); err != nil {
            log.Printf("[ERROR] Scanning row: %v", err)
            continue
        }
		fmt.Printf("%s\t%.2f\n", cat, total)
	}
}
