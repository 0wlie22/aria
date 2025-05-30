package importer

import (
	"database/sql"
	"log"
	"strings"
)

type Category struct {
	Name        string
	Type        string
	DisplayName string
	Keywords    []string
}

var defaultCategories = []Category{
	{"skip", "skip", "Skip", []string{"revolut", "revolt", "lv89parx0002056954007", "lv55parx0027855350002"}},
	{"salary", "income", "Salary", []string{"darba alga", "komand. d"}},
	{"scholarship", "income", "Scholarship", []string{"stipendija"}},
	{"other_income", "income", "Other Income", []string{}},

	{"food", "expense", "Food", []string{"rimi", "maxima", "lidl"}},
	{"wellness", "expense", "Wellness", []string{"drogas", "aptieka", "kiko", "fielmann"}},
	{"shopping", "expense", "Shopping", []string{"new yorker", "h&m", "lindex", "pepco", "reserved"}},
	{"gifts", "expense", "Gifts", []string{}},
	{"entertainment", "expense", "Entertainment", []string{}},
	{"transport", "expense", "Transport", []string{"citybee", "bolt", "narvesen", "tvm"}},
	{"phone", "expense", "Phone", []string{"tele2"}},
	{"hobby", "expense", "Hobby", []string{"nartiss.lv"}},
	{"sport", "expense", "Sport", []string{"trufit.eu", "nike", "decathlon", "baltic events", "myfitness", "sportdirect", "slidotava-akropole"}},

	{"citadele_investment", "investment", "Citadele Investment", []string{"lv51parx0027855351141", "lv98parx0027855350004"}},
	{"retirement_fund", "investment", "Retirement Fund", []string{"pensiju fond"}},
	{"roundups", "investment", "Roundups", []string{"pigr"}},
}

func SeedCategories(db *sql.DB) error {
	for _, c := range defaultCategories {
		_, err := db.Exec(`
			INSERT INTO categories (name, type, display_name, keywords)
			VALUES ($1, $2, $3, $4)
            ON CONFLICT (name) DO NOTHING
		`, c.Name, c.Type, c.DisplayName, "{"+strings.Join(c.Keywords, ",")+"}")
		if err != nil {
			log.Printf("[DB ERROR] Inserting category %s: %v", c.Name, err)
			return err
		}
	}

	log.Printf("[DB] Default categories seeded successfully")
	return nil
}
