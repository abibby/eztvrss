package models

import (
	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/model"
)

//go:generate spice generate:migration
type EZShow struct {
	model.BaseModel

	ID   int    `db:"id,primary"`
	Slug string `db:"slug"`

	Shows *builder.HasMany[*Show]
}

func EZShowQuery() *builder.Builder[*EZShow] {
	return builder.From[*EZShow]()
}

func (s *EZShow) Table() string {
	return "ez_shows"
}
