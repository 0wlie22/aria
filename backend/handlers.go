package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func setCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
}

func ExpensesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setCORS(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		month := r.URL.Query().Get("month")
		year := r.URL.Query().Get("year")
		log.Printf("Query expenses for month=%s year=%s", month, year)

		query := queryExpenses
		var args []any
		if month != "" && year != "" {
			query += " AND EXTRACT(MONTH FROM date) = $1 AND EXTRACT(YEAR FROM date) = $2"
			args = append(args, month, year)
		} else if year != "" {
			query += " AND EXTRACT(YEAR FROM date) = $1"
			args = append(args, year)
		}
		query += " GROUP BY category"

		rows, err := db.Query(query, args...)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		var expenses []Expense
		for rows.Next() {
			var e Expense
			if err := rows.Scan(&e.Category, &e.Total); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			expenses = append(expenses, e)
		}

		json.NewEncoder(w).Encode(expenses)
	}
}

func AvailableMonthsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setCORS(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		rows, err := db.Query(queryAvailableMonths)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		var months []MonthYear
		for rows.Next() {
			var m MonthYear
			if err := rows.Scan(&m.Year, &m.Month); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			months = append(months, m)
		}

		json.NewEncoder(w).Encode(months)
	}
}

func TotalHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setCORS(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		month := r.URL.Query().Get("month")
		year := r.URL.Query().Get("year")
		log.Printf("Query expenses for month=%s year=%s", month, year)

		rows, err := db.Query(queryTotalExpenseIncome)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		var totals []struct {
			Type  string  `json:"type"`
			Total float64 `json:"total"`
			Month int     `json:"month,omitempty"`
			Year  int     `json:"year,omitempty"`
		}

		for rows.Next() {
			var t struct {
				Type  string
				Total float64
				Month int
				Year  int
			}
			if err := rows.Scan(&t.Type, &t.Year, &t.Month, &t.Total); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			totals = append(totals, struct {
				Type  string  `json:"type"`
				Total float64 `json:"total"`
				Month int     `json:"month,omitempty"`
				Year  int     `json:"year,omitempty"`
			}{
				Type:  t.Type,
				Total: t.Total,
				Month: t.Month,
				Year:  t.Year,
			})
		}

		json.NewEncoder(w).Encode(totals)
	}
}
