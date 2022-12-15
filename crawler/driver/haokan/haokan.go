package haokan

import (
	"fmt"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "haokan"
	ProxySwitch = false
	URL         = "https://haokan.baidu.com/videoui/api/hotwords"
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
	if err := httputil.Request("GET", URL, nil, "json", &body, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	} else if body.Errno != 0 {
		return nil, fmt.Errorf("body errno: %d", body.Errno)
	}

	board := hot.NewBoard(c.Name())
	for _, hotword := range body.Data.Response.Hotwords {
		board.Append1(hotword)
	}
	return board, nil
}

type body struct {
	Data struct {
		Response struct {
			Hotwords []string `json:"hotwords"`
		} `json:"response"`
	} `json:"data"`
	Errno int    `json:"errno"`
	Error string `json:"error"`
}
