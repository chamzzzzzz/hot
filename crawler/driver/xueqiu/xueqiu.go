package xueqiu

import (
	"fmt"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
)

const (
	DriverName  = "xueqiu"
	ProxySwitch = false
	URL         = "https://xueqiu.com/query/v1/status/hots.json?count=10&page=1&scope=day&type=news"
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
	cookie string
}

func (c *Crawler) Driver() driver.Driver {
	return &Driver{}
}

func (c *Crawler) Name() string {
	return DriverName
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	if err := c.updatecookie(); err != nil {
		return nil, err
	}

	body := &body{}
	option := httputil.NewOption(c.Option, ProxySwitch)
	option.Header.Set("Cookie", c.cookie)
	if err := httputil.Request("GET", URL, nil, "json", &body, option); err != nil {
		return nil, err
	} else if body.Code != 200 {
		return nil, fmt.Errorf("body code: %d", body.Code)
	}

	board := hot.NewBoard(c.Name())
	for _, data := range body.Data {
		title := strings.TrimSpace(data.Title)
		url := "https://xueqiu.com" + strings.TrimSpace(data.Target)
		board.AppendTitleURL(title, url)
	}
	return board, nil
}

func (c *Crawler) updatecookie() error {
	return httputil.UpdateCookie("HEAD", "https://xueqiu.com", nil, httputil.NewOption(c.Option, ProxySwitch), &c.cookie)
}

type body struct {
	Code int `json:"code"`
	Data []struct {
		Target string `json:"target"`
		Text   string `json:"text"`
		Title  string `json:"title"`
	} `json:"data"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}
