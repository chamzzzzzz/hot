package bilibili

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
)

const (
	DriverName  = "bilibili"
	ProxySwitch = false
	URL         = "https://api.bilibili.com/x/web-interface/search/square?limit=10"
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
	for _, keyword := range body.Data.Trending.List {
		title := strings.TrimSpace(keyword.ShowName)
		board.Append1(title)
	}
	return board, nil
}

type body struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	TTL     int    `json:"ttl,omitempty"`
	Data    struct {
		Trending struct {
			Title   string `json:"title,omitempty"`
			Trackid string `json:"trackid,omitempty"`
			List    []struct {
				Keyword  string `json:"keyword,omitempty"`
				ShowName string `json:"show_name,omitempty"`
				Icon     string `json:"icon,omitempty"`
				URI      string `json:"uri,omitempty"`
				Goto     string `json:"goto,omitempty"`
			} `json:"list,omitempty"`
		} `json:"trending,omitempty"`
	} `json:"data,omitempty"`
}
