package eztv

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/abibby/salusa/clog"
	"github.com/inhies/go-bytesize"
)

type ShowItem struct {
	Title         string
	Link          string
	PubDate       time.Time
	ContentLength int64
	MagnetURI     string
	TorrentURI    string
	Seeds         int
	FileName      string
}

func Show(ctx context.Context, id int, slug string) ([]*ShowItem, error) {
	uri := eztvURL(fmt.Sprintf("/shows/%d/%s/", id, slug))
	resp, err := http.Get(uri)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("http error: %s", resp.Status)
	}

	d, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	// table.forum_header_noborder
	episodeElements := d.Find("table.forum_header_noborder tr.forum_header_border")

	items := []*ShowItem{}

	torrentExtRE := regexp.MustCompile(`\.torrent$`)

	episodeElements.Each(func(i int, s *goquery.Selection) {
		children := s.Children()
		link := s.Find(".epinfo")
		seedsStr := s.Find(".forum_thread_post_end font").Text()
		seeds := 0
		if seedsStr != "" {
			seeds, err = strconv.Atoi(seedsStr)
			if err != nil {
				clog.Use(ctx).Warn("failed to parse seeds", "error", err)
				return
			}
		}
		lenStr := children.Eq(3).Text()
		length, err := bytesize.Parse(lenStr)
		if err != nil {
			return
		}
		torrentURI := s.Find(".download_1").AttrOr("href", "")

		items = append(items, &ShowItem{
			Title:         link.Text(),
			Link:          "https://eztv.re" + link.AttrOr("href", ""),
			PubDate:       time.Now(),
			MagnetURI:     s.Find(".magnet").AttrOr("href", ""),
			TorrentURI:    torrentURI,
			ContentLength: int64(length),
			Seeds:         seeds,
			FileName:      torrentExtRE.ReplaceAllString(path.Base(torrentURI), ""),
		})
	})

	return items, nil
}
