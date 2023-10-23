package jobs

import (
	"context"
	"errors"

	"github.com/abibby/eztvrss/app/events"
	"github.com/abibby/eztvrss/app/models"
	"github.com/abibby/eztvrss/database"
	"github.com/abibby/eztvrss/services/eztv"
	"github.com/abibby/salusa/clog"
	"github.com/abibby/salusa/database/model"
	"github.com/jmoiron/sqlx"
)

func FetchRSS(ctx context.Context, e *events.FetchRSSEvent) error {
	return database.Tx(ctx, func(tx *sqlx.Tx) error {
		clog.Use(ctx).Info("fetch rss")

		channel, err := eztv.RSS()
		if err != nil {
			return err
		}
		for _, i := range channel.Items {
			ep, err := models.EpisodeQuery().Where("guid", "=", i.GUID).First(tx)
			if err != nil {
				return err
			}
			if ep != nil {
				continue
			}

			s, err := models.FetchOrCreateShow(tx, i.Title)
			if errors.Is(err, models.ErrFailedToParse) {
				clog.Use(ctx).Warn("failed to parse", "title", i.Title, "error", err)
				continue
			} else if err != nil {
				return err
			}
			if i.Enclosure == nil {
				i.Enclosure = &eztv.Enclosure{}
			}
			ep = &models.Episode{
				ShowID: s.ID,

				Category:      i.Category,
				Link:          i.Link,
				GUID:          i.GUID,
				PubDate:       i.PubDate,
				Title:         i.Title,
				ContentLength: i.ContentLength,
				InfoHash:      i.InfoHash,
				MagnetURI:     i.MagnetURI.Value,
				Seeds:         i.Seeds,
				Peers:         i.Peers,
				Verified:      i.Verified,
				FileName:      i.FileName,
				TorrentURI:    i.Enclosure.URL,
			}

			err = model.Save(tx, ep)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
