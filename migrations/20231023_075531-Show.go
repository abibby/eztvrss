package migrations

import (
	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/database/schema"
)

func init() {
	migrations.Add(&migrate.Migration{
		Name: "20231023_075531-Show",
		Up: schema.Create("shows", func(table *schema.Blueprint) {
			table.Int("id").Primary().AutoIncrement()
			table.String("name")
			table.String("slug")
		}),
		Down: schema.DropIfExists("shows"),
	})
}
