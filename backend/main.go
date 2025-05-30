package main

import (
	"log"
	"net/http"

	"aria/importer"
)

func main() {
	db := initDB()
	defer db.Close()

	importer.RunCSVImport(db)

	log.Println("Starting API server on :8080")
	http.HandleFunc("/api/expenses", ExpensesHandler(db))
	http.HandleFunc("/api/available-months", AvailableMonthsHandler(db))
    http.HandleFunc("/api/total", TotalHandler(db))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
