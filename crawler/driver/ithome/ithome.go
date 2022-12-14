package ithome

import (
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	IT      = "it"
	Game    = "game"
	Unknown = "unknown"
)

const (
	DriverName = "ithome"
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
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.ithome.com/block/rank.html?d=game", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	dom := soup.HTMLParse(string(data))
	if dom.Error != nil {
		return nil, dom.Error
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
