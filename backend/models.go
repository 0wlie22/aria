package main

type Expense struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
}
type MonthYear struct {
	Month int `json:"month"`
	Year  int `json:"year"`
}
