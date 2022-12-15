package nowcoder

import (
	"fmt"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "nowcoder"
	ProxySwitch = false
	URL         = "https://www.nowcoder.com/nccommon/search/hot-query"
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
	} else if body.Code != 0 {
		return nil, fmt.Errorf("body code: %d", body.Code)
	}

	board := hot.NewBoard(c.Name())
	for _, data := range body.Data {
		board.Append1(data.Query)
	}
	return board, nil
}

type body struct {
	Code int `json:"code"`
	Data []struct {
		HotValue int    `json:"hotValue"`
		Query    string `json:"query"`
		Rank     int    `json:"rank"`
	} `json:"data"`
	Msg string `json:"msg"`
}
