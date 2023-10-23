package migrations

import (
	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/database/schema"
)

func init() {
	migrations.Add(&migrate.Migration{
		Name: "20231023_094711-Show",
		Up: schema.Table("shows", func(table *schema.Blueprint) {
			table.Int("ez_show_id")
			// table.ForeignKey("ez_show_id", "ez_shows", "id")
		}),
		Down: schema.Table("shows", func(table *schema.Blueprint) {
			table.DropColumn("ez_show_id")
		}),
	})
}
