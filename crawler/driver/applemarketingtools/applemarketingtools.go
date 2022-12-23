package applemarketingtools

import (
	"fmt"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"time"
)

const (
	CN = "cn"
	US = "us"
	GB = "gb"
	HK = "hk"
	MO = "mo"
	JP = "jp"
	KR = "kr"
	SG = "sg"
	IN = "in"
)

var (
	Catalogs = []string{CN, US, GB, HK, MO, JP, KR, SG, IN}
)

const (
	DriverName  = "applemarketingtools"
	ProxySwitch = false
	URL         = "https://rss.applemarketingtools.com/api/v2/%s.json"
)

type Driver struct {
}

func (driver *Driver) Open(option driver.Option) (driver.Crawler, error) {
	return &Crawler{Option: option}, nil
}

func init() {
	driver.Register(DriverName, &Driver{})
}

type Crawler struct {
	Option driver.Option
}

func (c *Crawler) Driver() driver.Driver {
	return &Driver{}
}

func (c *Crawler) Name() string {
	return DriverName
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	if c.Option.Catalog == "" {
		return c.all()
	} else {
		return c.rss(nil, c.Option.Catalog)
	}
}

func (c *Crawler) all() (*hot.Board, error) {
	board := hot.NewBoard(c.Name())
	for _, catalog := range Catalogs {
		if _, err := c.rss(board, catalog); err != nil {
			return nil, err
		}
	}
	return board, nil
}

func (c *Crawler) rss(board *hot.Board, catalog string) (*hot.Board, error) {
	if board == nil {
		board = hot.NewBoard(c.Name())
	}
	for _, feed := range feeds(catalog) {
		if _, err := c.feed(board, catalog, feed); err != nil {
			return nil, err
		}
	}
	return board, nil
}

func (c *Crawler) feed(board *hot.Board, catalog string, feed string) (*hot.Board, error) {
	body := &body{}
	if err := httputil.Request("GET", fmt.Sprintf(URL, feed), nil, "json", body, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		fmt.Println(feed)
		return nil, err
	}

	if board == nil {
		board = hot.NewBoard(c.Name())
	}
	for _, item := range body.Feed.Results {
		title := item.Name
		summary := item.ArtistName
		url := item.URL
		pubdate := time.Now()
		if item.ReleaseDate != "" {
			if date, err := time.Parse("2006-01-02", item.ReleaseDate); err != nil {
				if date, err := time.ParseInLocation("2006", item.ReleaseDate, time.Local); err != nil {
					return nil, err
				} else {
					pubdate = date
				}
			} else {
				pubdate = date
			}
		}
		board.Append5x1(title, summary, url, catalog, feed, pubdate)
	}
	return board, nil
}

func feeds(catalog string) []string {
	var v []string
	v = append(v, fmt.Sprintf("%s/apps/top-free/50/apps", catalog))
	v = append(v, fmt.Sprintf("%s/apps/top-paid/50/apps", catalog))
	v = append(v, fmt.Sprintf("%s/music/most-played/50/albums", catalog))
	v = append(v, fmt.Sprintf("%s/music/most-played/50/music-videos", catalog))
	v = append(v, fmt.Sprintf("%s/music/most-played/50/songs", catalog))
	v = append(v, fmt.Sprintf("%s/podcasts/top/50/podcasts", catalog))

	switch catalog {
	case CN:
	default:
		v = append(v, fmt.Sprintf("%s/books/top-free/50/books", catalog))
		v = append(v, fmt.Sprintf("%s/books/top-paid/50/books", catalog))
	}

	switch catalog {
	case US, GB, JP:
		v = append(v, fmt.Sprintf("%s/audio-books/top/50/audio-books", catalog))
	}
	return v
}

type body struct {
	Feed struct {
		Results []struct {
			ArtistName  string `json:"artistName"`
			Name        string `json:"name"`
			ReleaseDate string `json:"releaseDate"`
			URL         string `json:"url"`
		} `json:"results"`
	} `json:"feed"`
}
