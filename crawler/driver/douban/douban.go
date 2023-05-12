package douban

import (
	"strings"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	Note  = "note"
	Movie = "movie"
)

const (
	DriverName  = "douban"
	ProxySwitch = false
	URL         = "https://movie.douban.com"
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
	case Note:
		return c.note()
	case Movie:
		return c.movie()
	default:
		return c.all()
	}
}

func (c *Crawler) all() (*hot.Board, error) {
	board := hot.NewBoard(c.Name())
	b1, err := c.note()
	if err != nil {
		return nil, err
	}
	b2, err := c.movie()
	if err != nil {
		return nil, err
	}
	for _, hot := range b1.Hots {
		board.Append(hot)
	}
	for _, hot := range b2.Hots {
		board.Append(hot)
	}
	return board, nil
}

func (c *Crawler) movie() (*hot.Board, error) {
	dom := &httputil.DOM{}
	if err := httputil.Request("GET", URL, nil, "dom", dom, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	div, err := dom.Find("div", "class", "screening-bd")
	if err != nil {
		return nil, err
	}
	for _, li := range div.QueryAll("li", "class", "title") {
		a, err := li.Find("a")
		if err != nil {
			return nil, err
		}
		title := strings.TrimSpace(a.Text())
		url := strings.TrimSpace(a.Href())
		board.Append(&hot.Hot{
			Title:   title,
			URL:     url,
			Catalog: Movie,
		})
	}
	return board, nil
}

func (c *Crawler) note() (*hot.Board, error) {
	dom := &httputil.DOM{}
	if err := httputil.Request("GET", "https://www.douban.com", nil, "dom", dom, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	div, err := dom.Find("div", "class", "notes")
	if err != nil {
		return nil, err
	}
	for _, a := range div.QueryAll("a") {
		title := strings.TrimSpace(a.Text())
		url := strings.TrimSpace(a.Href())
		board.Append(&hot.Hot{
			Title:   title,
			URL:     url,
			Catalog: Note,
		})
	}
	return board, nil
}
