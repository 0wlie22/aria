package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func initDB() *sql.DB {
	connStr := "postgres://finance_user:" + getPassword() + "@localhost:5432/aria?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("[DB ERROR] Cannot open DB:", err)
	}

	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    type TEXT CHECK (type IN ('income', 'expense', 'investment', 'skip')) NOT NULL,
    display_name TEXT NOT NULL,
    keywords TEXT[]
);
CREATE TABLE IF NOT EXISTS transactions (
    fingerprint TEXT PRIMARY KEY,
    category_id INT REFERENCES categories(id) ON DELETE SET NULL,
    date DATE NOT NULL,
    amount INTEGER NOT NULL,
    keyword TEXT
);
        `)
	if err != nil {
		log.Fatal("[DB ERROR] Failed to create table:", err)
	}

	log.Println("[DB] Connected and table ensured.")
	return db
}

func getPassword() string {
	pass := os.Getenv("PSQL_PASSWORD")
	if pass == "" {
		pass = "test"
	}
	return pass
}
