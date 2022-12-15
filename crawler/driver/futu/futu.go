package futu

import (
	"fmt"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "futu"
	ProxySwitch = false
	URL         = "https://www.futunn.com/search-stock/hot-news"
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
	} else if body.Code != "0" {
		return nil, fmt.Errorf("body code: %s", body.Code)
	}

	board := hot.NewBoard(c.Name())
	for _, data := range body.Data {
		board.AppendTitleURL(data.Title, data.URL)
	}
	return board, nil
}

type body struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    []struct {
		PostID int    `json:"post_id"`
		Title  string `json:"title"`
		URL    string `json:"url"`
	} `json:"data"`
}
