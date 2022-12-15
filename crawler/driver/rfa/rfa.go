package rfa

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
)

const (
	Cantonese = "cantonese"
	English   = "english"
	Mandarin  = "mandarin"
)

const (
	DriverName  = "rfa"
	ProxySwitch = true
	URL         = "https://www.rfa.org/"
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
	case Cantonese:
		return c.cantonese()
	case English:
		return c.english()
	case Mandarin:
		return c.mandarin()
	default:
		return c.all()
	}
}

func (c *Crawler) all() (*hot.Board, error) {
	board := hot.NewBoard(c.Name())
	if b, err := c.cantonese(); err != nil {
		return nil, err
	} else {
		for _, hot := range b.Hots {
			board.Append0(hot)
		}
	}

	if b, err := c.english(); err != nil {
		return nil, err
	} else {
		for _, hot := range b.Hots {
			board.Append0(hot)
		}
	}

	if b, err := c.mandarin(); err != nil {
		return nil, err
	} else {
		for _, hot := range b.Hots {
			board.Append0(hot)
		}
	}
	return board, nil
}

func (c *Crawler) cantonese() (*hot.Board, error) {
	return c.withLanguage(Cantonese)
}

func (c *Crawler) english() (*hot.Board, error) {
	return c.withLanguage(English)
}

func (c *Crawler) mandarin() (*hot.Board, error) {
	return c.withLanguage(Mandarin)
}

func (c *Crawler) withLanguage(language string) (*hot.Board, error) {
	dom := &httputil.DOM{}
	if err := httputil.Request("GET", URL+language, nil, "dom", dom, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	for _, div := range dom.FindAllStrict("div", "class", "most_read_only first_most_bold") {
		for _, li := range div.FindAll("li") {
			a := li.Find("a")
			if a.Error != nil {
				return nil, a.Error
			}
			title := strings.TrimSpace(a.Text())
			url := strings.TrimSpace(a.Attrs()["href"])
			catalog := language
			board.Append4(title, "", url, catalog)
		}
	}
	return board, nil
}
