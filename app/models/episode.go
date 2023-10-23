package models

import (
	"context"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/hooks"
	"github.com/abibby/salusa/database/model"
)

//go:generate spice generate:migration
type Episode struct {
	model.BaseModel

	ID int `db:"id,primary,autoincrement" xml:"-"`

	Category      string     `db:"category"             xml:"category"`
	Link          string     `db:"link"                 xml:"link"`
	GUID          string     `db:"guid"                 xml:"guid"`
	PubDate       string     `db:"pub_date"             xml:"pubDate"`
	Title         string     `db:"title"                xml:"title"`
	ContentLength int64      `db:"content_length"       xml:"torrent:contentLength"`
	InfoHash      string     `db:"info_hash"            xml:"torrent:infoHash"`
	MagnetURI     string     `db:"magnet_uri,type:text" xml:"-"`
	XMLMagnetURI  *MagnetURI `db:"-"                    xml:"torrent:magnetURI"`
	Seeds         int        `db:"seeds"                xml:"torrent:seeds"`
	Peers         int        `db:"peers"                xml:"torrent:peers"`
	Verified      int        `db:"verified"             xml:"torrent:verified"`
	FileName      string     `db:"filename"             xml:"torrent:fileName"`
	TorrentURI    string     `db:"torrent_uri"          xml:"-"`
	Enclosure     *Enclosure `db:"-"                    xml:"enclosure"`

	ShowID int `db:"show_id" xml:"-"`
	Show   *builder.BelongsTo[*Show]
}

type MagnetURI struct {
	Value string `xml:",cdata"`
}

type Enclosure struct {
	URL    string `xml:"url,attr"`
	Length int64  `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

var _ hooks.AfterLoader = (*Episode)(nil)
var _ hooks.BeforeSaver = (*Episode)(nil)

func EpisodeQuery() *builder.Builder[*Episode] {
	return builder.From[*Episode]()
}

func (e *Episode) AfterLoad(ctx context.Context, tx hooks.DB) error {
	e.XMLMagnetURI = &MagnetURI{Value: e.MagnetURI}
	e.Enclosure = &Enclosure{
		URL:    e.TorrentURI,
		Length: e.ContentLength,
		Type:   "application/x-bittorrent",
	}
	return nil
}
func (e *Episode) BeforeSave(ctx context.Context, tx hooks.DB) error {
	// e.RawMagnetURI = e.MagnetURI.Value
	// e.TorrentURI = e.Enclosure.URL
	return nil
}
