package ithome

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
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
		ul := dom.FindStrict("ul", "id", ulId)
		if ul.Error != nil {
			return nil, ul.Error
		}
		for _, a := range ul.FindAllStrict("a") {
			title := strings.TrimSpace(a.Text())
			url := strings.TrimSpace(a.Attrs()["href"])
			catalog := idtocatalog(ulId)
			board.Append4(title, "", url, catalog)
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
