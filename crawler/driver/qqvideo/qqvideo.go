package qqvideo

import (
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"io/ioutil"
	"net/http"
	"strings"
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
	DriverName = "qqvideo"
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
	req, err := http.NewRequest("GET", "https://v.qq.com/biu/ranks/?t=hotsearch", nil)
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
	for i, ol := range dom.FindAllStrict("ol", "class", "hotlist") {
		catalog := itocatalog(i)
		if c.Option.Catalog != "" && c.Option.Catalog != catalog {
			continue
		}
		for _, li := range ol.FindAllStrict("li", "class", "item item_odd item_1") {
			a := li.Find("a")
			if a.Error != nil {
				return nil, a.Error
			}
			title := strings.TrimSpace(a.Attrs()["title"])
			url := "https:" + strings.TrimSpace(a.Attrs()["href"])
			board.Append4(title, "", url, catalog)
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
