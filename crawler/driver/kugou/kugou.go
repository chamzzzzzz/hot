package kugou

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	Surge    = "surge"
	Top500   = "top500"
	Douyin   = "douyin"
	Kuaishou = "kuaishou"
	DJ       = "dj"
	Mainland = "mainland"
	HK       = "hk"
	TW       = "tw"
)

var pages = map[string]string{
	Surge:    "6666",
	Top500:   "8888",
	Douyin:   "52144",
	Kuaishou: "52767",
	DJ:       "24971",
	Mainland: "31308",
	HK:       "31313",
	TW:       "54848",
}

const (
	DriverName = "kugou"
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
	case Surge, Top500, Douyin, Kuaishou, DJ, Mainland, HK, TW:
		return c.withCatalog(c.Option.Catalog, nil)
	default:
		return c.all()
	}
}

func (c *Crawler) all() (*hot.Board, error) {
	board := hot.NewBoard(c.Name())
	for _, catalog := range []string{Surge, Top500, Douyin, Kuaishou, DJ, Mainland, HK, TW} {
		if _, err := c.withCatalog(catalog, board); err != nil {
			return nil, err
		}
	}
	return board, nil
}

func (c *Crawler) withCatalog(catalog string, board *hot.Board) (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.kugou.com/yy/rank/home/1-%s.html?from=rank", pages[catalog]), nil)
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

	if board == nil {
		board = hot.NewBoard(c.Name())
	}
	div := dom.Find("div", "id", "rankWrap")
	if div.Error != nil {
		return nil, div.Error
	}
	for _, a := range div.FindAllStrict("a", "class", "pc_temp_songname") {
		title := strings.TrimSpace(a.Attrs()["title"])
		url := strings.TrimSpace(a.Attrs()["href"])
		board.AppendTitleURLCatalog(title, url, catalog)
	}
	return board, nil
}
