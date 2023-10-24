package handlers

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/abibby/eztvrss/app/models"
	"github.com/abibby/eztvrss/services/eztv"
	"github.com/abibby/salusa/clog"
	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/set"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type ShowsRequest struct {
	ID      int    `url:"id"`
	Slug    string `url:"slug"`
	Query   string `query:"q"`
	Request *http.Request
	Ctx     context.Context
}
type ShowsResponse struct {
	RSS *Root
}

type Root struct {
	XMLName xml.Name `xml:"rss"`
	Channel *Channel `xml:"channel"`
}

type Channel struct {
	XMLName       xml.Name          `xml:"channel"`
	Title         string            `xml:"title"`
	Items         []*models.Episode `xml:"item"`
	LastBuildDate string            `xml:"lastBuildDate"`
}

var _ request.Responder = (*ShowsResponse)(nil)

func (sr *ShowsResponse) Respond(w http.ResponseWriter, r *http.Request) error {
	e := xml.NewEncoder(w)
	e.Indent("", "    ")
	return e.Encode(sr)
}

var fetchedShows = set.New[int]()

var Shows = request.Handler(func(r *ShowsRequest) (*ShowsResponse, error) {
	vars := mux.Vars(r.Request)
	r.ID, _ = strconv.Atoi(vars["id"])
	r.Slug = vars["id"]

	tx := request.UseTx(r.Request)

	if !fetchedShows.Has(r.ID) {
		err := fetchShows(r.Ctx, tx, r.ID, r.Slug)
		if err != nil {
			return nil, err
		}
		fetchedShows.Add(r.ID)
	}
	q := models.EpisodeQuery().WhereHas("Show", func(q *builder.SubBuilder) *builder.SubBuilder {
		return q.Where("ez_show_id", "=", r.ID)
	})
	if r.Query != "" {
		q = q.Where("title", "like", fmt.Sprintf("%%%s%%", r.Query))
	}

	episodes, err := q.Get(tx)
	if err != nil {
		return nil, err
	}

	return &ShowsResponse{
		RSS: &Root{
			Channel: &Channel{
				Title:         "TV Torrents RSS feed - EZTV",
				LastBuildDate: time.Now().Format(time.RFC1123Z),
				Items:         episodes,
			},
		},
	}, nil
})

func fetchShows(ctx context.Context, tx *sqlx.Tx, id int, slug string) error {
	ezEpisodes, err := eztv.Show(ctx, id, slug)
	if err != nil {
		return err
	}
	m := map[string]*models.Show{}
	for _, e := range ezEpisodes {
		name, err := models.ShowNameFromTitle(e.Title)
		if err != nil {
			clog.Use(ctx).Warn("failed to parse", "title", e.Title, "error", err)
			continue
		}

		if _, ok := m[name]; !ok {
			s, err := models.FetchOrCreateShow(tx, e.Title)
			if errors.Is(err, models.ErrFailedToParse) {
				clog.Use(ctx).Warn("failed to parse", "title", e.Title, "error", err)
				continue
			} else if err != nil {
				return err
			}

			if s.EZShowID != id {
				s.EZShowID = id
				err = model.Save(tx, s)
				if err != nil {
					return err
				}
			}
			m[name] = s
		}
	}
	for _, e := range ezEpisodes {
		ep, err := models.EpisodeQuery().Where("guid", "=", e.Link).First(tx)
		if err != nil {
			return err
		}
		if ep != nil {
			continue
		}

		u, err := url.Parse(e.MagnetURI)
		if err != nil {
			clog.Use(ctx).Warn("failed to parse magnet", "magnet", e.MagnetURI, "error", err)
			continue
		}
		name, err := models.ShowNameFromTitle(e.Title)
		if err != nil {
			continue
		}
		s := m[name]
		err = model.Save(tx, &models.Episode{
			Category:      "TV",
			Link:          e.Link,
			GUID:          e.Link,
			PubDate:       e.PubDate.Format(time.RFC1123Z),
			Title:         e.Title,
			ContentLength: e.ContentLength,
			InfoHash:      infoHashString(u.Query().Get("xt")),
			MagnetURI:     e.MagnetURI,
			Seeds:         e.Seeds,
			FileName:      e.FileName,
			TorrentURI:    e.TorrentURI,

			ShowID: s.ID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func infoHashString(xt string) string {
	hash, ok := strings.CutPrefix(xt, "urn:btih:")
	if ok {
		return hash
	}
	hash, ok = strings.CutPrefix(xt, "urn:btmh:")
	if ok {
		return hash
	}
	return ""
}
