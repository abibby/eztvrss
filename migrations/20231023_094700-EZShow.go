package migrations

import (
	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/database/schema"
)

func init() {
	migrations.Add(&migrate.Migration{
		Name: "20231023_094700-EZShow",
		Up: schema.Create("ez_shows", func(table *schema.Blueprint) {
			table.Int("id").Primary()
			table.String("slug")
		}),
		Down: schema.DropIfExists("ez_shows"),
	})
}
