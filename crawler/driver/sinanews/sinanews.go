package sinanews

import (
	"strings"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "sinanews"
	ProxySwitch = false
	URL         = "https://newsapp.sina.cn/api/hotlist"
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
	for _, data := range body.Data.HotList {
		board.Append(&hot.Hot{Title: strings.TrimSpace(data.Info.Title), URL: strings.TrimSpace(data.Base.Base.URL)})
	}
	return board, nil
}

type body struct {
	Status int    `json:"status,omitempty"`
	Msg    string `json:"msg,omitempty"`
	Data   struct {
		HotList []struct {
			Base struct {
				Base struct {
					URL string `json:"url,omitempty"`
				} `json:"base,omitempty"`
			} `json:"base,omitempty"`
			Info struct {
				Title string `json:"title,omitempty"`
			} `json:"info,omitempty"`
		} `json:"hotList,omitempty"`
	} `json:"data,omitempty"`
}

func (body *body) NormalizedCode() int {
	return body.Status
}
