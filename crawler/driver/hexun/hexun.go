package hexun

import (
	"strings"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "hexun"
	ProxySwitch = false
	URL         = "https://m.hexun.com/api/getHotFocus"
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
	for _, data := range body.Data {
		board.Append(&hot.Hot{Title: strings.TrimSpace(data.Title), URL: strings.TrimSpace(data.URL)})
	}
	return board, nil
}

type body struct {
	Data []struct {
		Title string `json:"title"`
		URL   string `json:"newsUrl"`
	} `json:"datas"`
}
