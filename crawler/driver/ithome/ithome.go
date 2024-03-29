package ithome

import (
	"strings"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	IT      = "it"
	Game    = "game"
	Unknown = "unknown"
)

const (
	DriverName  = "ithome"
	ProxySwitch = false
	URL         = "https://www.ithome.com/block/rank.html?d=game"
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
	for _, ulId := range catalogtoids(c.Option.Catalog) {
		ul, err := dom.Find("ul", "id", ulId)
		if err != nil {
			return nil, err
		}
		for _, a := range ul.QueryAll("a") {
			title := strings.TrimSpace(a.Text())
			url := strings.TrimSpace(a.Href())
			catalog := idtocatalog(ulId)
			board.Append(&hot.Hot{
				Title:   title,
				URL:     url,
				Catalog: catalog,
			})
		}
	}
	return board, nil
}

func idtocatalog(id string) string {
	switch id {
	case "d-1":
		return IT
	case "d-4":
		return Game
	default:
		return Unknown
	}
}

func catalogtoids(catalog string) []string {
	switch catalog {
	case "":
		return []string{"d-1", "d-4"}
	case IT:
		return []string{"d-1"}
	case Game:
		return []string{"d-4"}
	default:
		return nil
	}
}
