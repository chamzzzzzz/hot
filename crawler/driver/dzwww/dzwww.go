package dzwww

import (
	"strings"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	HotSearch = "hotsearch"
	HotNews   = "hotnews"
)

var act = map[string]string{
	HotSearch: "getPCHotSearch",
	HotNews:   "getPCHotNews",
}

const (
	DriverName  = "dzwww"
	ProxySwitch = false
	URL         = "https://w.dzwww.com/?act="
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
	switch c.Option.Catalog {
	case HotSearch, HotNews:
		return c.withCatalog(c.Option.Catalog, nil)
	default:
		return c.all()
	}
}

func (c *Crawler) all() (*hot.Board, error) {
	board := hot.NewBoard(c.Name())
	for _, catalog := range []string{HotSearch, HotNews} {
		if _, err := c.withCatalog(catalog, board); err != nil {
			return nil, err
		}
	}
	return board, nil
}

func (c *Crawler) withCatalog(catalog string, board *hot.Board) (*hot.Board, error) {
	var body []body
	option := httputil.NewOption(c.Option, ProxySwitch)
	option.Header.Set("Referer", "https://w.dzwww.com/")
	if err := httputil.Request("GET", URL+act[catalog], nil, "json", &body, option); err != nil {
		return nil, err
	}

	if board == nil {
		board = hot.NewBoard(c.Name())
	}
	for _, data := range body {
		url := data.RealURL
		if !strings.HasPrefix(url, "http") {
			url = "https:" + url
		}
		board.Append(&hot.Hot{Title: data.Title, URL: url, Catalog: catalog})
	}
	return board, nil
}

type body struct {
	Title   string `json:"title"`
	RealURL string `json:"real_url"`
}
