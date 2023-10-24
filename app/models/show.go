package models

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/model"
	"github.com/jmoiron/sqlx"
)

var ErrFailedToParse = errors.New("failed to parse")

//go:generate spice generate:migration
type Show struct {
	model.BaseModel

	ID   int    `db:"id,primary,autoincrement"`
	Name string `db:"name"`

	EZShowID int `db:"ez_show_id" xml:"-"`

	Episodes *builder.HasMany[*Episode]
}

func ShowQuery() *builder.Builder[*Show] {
	return builder.From[*Show]()
}

var showTitleRE = regexp.MustCompile(`^(.+) (S\d+E\d+|\d{4} \d{2} \d{2}) (.+)$`)

func ShowNameFromTitle(title string) (string, error) {
	matches := showTitleRE.FindStringSubmatch(title)
	if len(matches) == 0 {
		return "", fmt.Errorf("title %s: %w", title, ErrFailedToParse)
	}
	return matches[1], nil
}

func FetchOrCreateShow(tx *sqlx.Tx, title string) (*Show, error) {
	name, err := ShowNameFromTitle(title)
	if err != nil {
		return nil, err
	}

	s, err := ShowQuery().Where("name", "=", name).First(tx)
	if err != nil {
		return nil, err
	}
	if s != nil {
		return s, nil
	}
	s = &Show{
		Name: name,
	}
	err = model.Save(tx, s)
	if err != nil {
		return nil, err
	}
	return s, nil
}
