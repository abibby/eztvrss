package eztv

import (
	"encoding/xml"
	"net/http"
)

type RSSRoot struct {
	XMLName xml.Name    `xml:"rss"`
	Channel *RSSChannel `xml:"channel"`
}

type RSSChannel struct {
	XMLName xml.Name `xml:"channel"`
	// Title         string   `xml:"title"`
	Items         []*RSSItem `xml:"item"`
	LastBuildDate string     `xml:"lastBuildDate"`
}

type RSSItem struct {
	XMLName xml.Name `xml:"item"`
	// Content  string   `xml:",innerxml"`
	Category string `xml:"category"`
	Link     string `xml:"link"`
	GUID     string `xml:"guid"`
	PubDate  string `xml:"pubDate"`
	Title    string `xml:"title"`

	// torrent: prefix
	ContentLength int64      `xml:"contentLength"`
	InfoHash      string     `xml:"infoHash"`
	MagnetURI     *MagnetURI `xml:"magnetURI"`
	Seeds         int        `xml:"seeds"`
	Peers         int        `xml:"peers"`
	Verified      int        `xml:"verified"`
	FileName      string     `xml:"fileName"`

	Enclosure *Enclosure `xml:"enclosure"`
}

type MagnetURI struct {
	// XMLName xml.Name `xml:"torrent:magnetURI"`
	Value string `xml:",cdata"`
}

type Enclosure struct {
	XMLName xml.Name `xml:"enclosure"`
	URL     string   `xml:"url,attr"`
	Length  int      `xml:"length,attr"`
	Type    string   `xml:"type,attr"`
}

func RSS() (*RSSChannel, error) {
	resp, err := http.Get(eztvURL("/ezrss.xml"))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r := &RSSRoot{}
	dec := xml.NewDecoder(resp.Body)
	dec.Strict = true
	err = dec.Decode(r)
	if err != nil {
		return nil, err
	}
	return r.Channel, nil
}
