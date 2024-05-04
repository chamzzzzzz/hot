package baike

import (
	"strings"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "baike"
	ProxySwitch = false
	URL         = "https://baike.baidu.com/cms/home/hotLemmas.json"
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
	for _, data := range body.Today {
		board.Append(&hot.Hot{Title: strings.TrimSpace(data.Title), URL: strings.TrimSpace(data.URL)})
	}
	return board, nil
}

type body struct {
	Today []struct {
		Title string `json:"title,omitempty"`
		URL   string `json:"url,omitempty"`
	} `json:"today,omitempty"`
}
