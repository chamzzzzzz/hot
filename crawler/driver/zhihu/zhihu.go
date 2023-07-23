package zhihu

import (
	"strings"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "zhihu"
	ProxySwitch = false
	URL         = "https://www.zhihu.com/api/v4/search/top_search"
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
	body := &body{}
	if err := httputil.Request("GET", URL, nil, "json", body, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	for _, data := range body.TopSearch.Words {
		board.Append(&hot.Hot{Title: strings.TrimSpace(data.DisplayQuery)})
	}
	return board, nil
}

type body struct {
	TopSearch struct {
		Words []struct {
			Query        string `json:"query"`
			DisplayQuery string `json:"display_query"`
		} `json:"words"`
	} `json:"top_search"`
}
