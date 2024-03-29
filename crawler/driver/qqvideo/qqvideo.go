package qqvideo

import (
	"strings"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	Search  = "search"
	TV      = "tv"
	Variety = "variety"
	Cartoon = "cartoon"
	Child   = "child"
	Movie   = "movie"
	Doco    = "doco"
	Games   = "games"
	Music   = "music"
	Unknown = "unknown"
)

var catalogs = []string{Search, TV, Variety, Cartoon, Child, Movie, Doco, Games, Music}

const (
	DriverName  = "qqvideo"
	ProxySwitch = false
	URL         = "https://v.qq.com/biu/ranks/?t=hotsearch"
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
	dom := &httputil.DOM{}
	if err := httputil.Request("GET", URL, nil, "dom", dom, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	for i, ol := range dom.QueryAll("ol", "class", "hotlist") {
		catalog := itocatalog(i)
		if c.Option.Catalog != "" && c.Option.Catalog != catalog {
			continue
		}
		for _, li := range ol.QueryAll("li", "class", "item item_odd item_1") {
			a, err := li.Find("a")
			if err != nil {
				return nil, err
			}
			title := strings.TrimSpace(a.Title())
			url := "https:" + strings.TrimSpace(a.Href())
			board.Append(&hot.Hot{
				Title:   title,
				URL:     url,
				Catalog: catalog,
			})
		}
	}
	return board, nil
}

func itocatalog(i int) string {
	if i >= 0 && i < len(catalogs) {
		return catalogs[i]
	} else {
		return Unknown
	}
}
