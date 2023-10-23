package migrations

import (
	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/database/schema"
)

func init() {
	migrations.Add(&migrate.Migration{
		Name: "20231023_093912-Show",
		Up: schema.Table("shows", func(table *schema.Blueprint) {
			table.DropColumn("slug")
		}),
		Down: schema.Table("shows", func(table *schema.Blueprint) {
			table.String("slug")
		}),
	})
}
