package migrations

import (
	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/database/schema"
)

func init() {
	migrations.Add(&migrate.Migration{
		Name: "20231023_080140-Episode",
		Up: schema.Create("episodes", func(table *schema.Blueprint) {
			table.Int("id").Primary().AutoIncrement()
			table.String("category")
			table.String("link")
			table.String("guid")
			table.String("pub_date")
			table.String("title")
			table.Int64("content_length")
			table.String("info_hash")
			table.Text("magnet_uri")
			table.Int("seeds")
			table.Int("peers")
			table.Int("verified")
			table.String("filename")
			table.String("torrent_uri")
			table.Int("show_id")
			table.ForeignKey("show_id", "shows", "id")
		}),
		Down: schema.DropIfExists("episodes"),
	})
}
