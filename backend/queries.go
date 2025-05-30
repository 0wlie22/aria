package main

const (
	queryExpenses = `
SELECT c.display_name AS category, COALESCE(SUM(t.amount), 0)/100.0 * -1 AS total
    FROM transactions t
    JOIN categories c ON t.category_id = c.id
    WHERE c.type = 'expense'
`

	queryAvailableMonths = `
    SELECT DISTINCT EXTRACT(YEAR FROM t.date) AS year, EXTRACT(MONTH FROM t.date) AS month
    FROM transactions t
    JOIN categories c ON t.category_id = c.id
    WHERE c.type = 'expense'
    ORDER BY year, month
`
    queryTotalExpenseIncome = `
    SELECT c.type, EXTRACT (YEAR FROM t.date) as year, EXTRACT(MONTH FROM t.date) AS month,
    COALESCE(SUM(t.amount), 0)/100 AS total
    FROM transactions t
    JOIN categories c ON t.category_id = c.id
    GROUP BY c.type, year, month
    ORDER BY year, month, c.type
`
)
