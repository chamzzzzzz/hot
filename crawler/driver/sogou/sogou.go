package sogou

import (
	"encoding/json"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "sogou"
	ProxySwitch = false
	URL         = "https://www.sogou.com/suggnew/hotwords"
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
	data, err := httputil.RequestData("GET", URL, nil, httputil.NewOption(c.Option, ProxySwitch))
	if err != nil {
		return nil, err
	}

	body := []string{}
	if err := json.Unmarshal(data[20:], &body); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	for _, word := range body {
		board.Append1(word)
	}
	return board, nil
}
