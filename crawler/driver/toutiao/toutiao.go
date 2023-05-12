package toutiao

import (
	"encoding/json"
	"fmt"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "toutiao"
	ProxySwitch = false
	URL         = "https://i-lq.snssdk.com/api/feed/hotboard_online/v1/?category=hotboard_online&count=50"
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
	if body.Message != "success" {
		return nil, fmt.Errorf("body message: %s", body.Message)
	}

	board := hot.NewBoard(c.Name())
	for _, data := range body.Data {
		content := &content{}
		if err := json.Unmarshal([]byte(data.Content), content); err != nil {
			return nil, err
		}
		board.Append(&hot.Hot{
			Title: content.RawData.Title,
		})
	}
	return board, nil
}

type body struct {
	Message string `json:"message"`
	Data    []struct {
		Content string `json:"content"`
		Code    string `json:"code"`
	} `json:"data"`
}

type content struct {
	RawData struct {
		Title string `json:"title"`
	} `json:"raw_data"`
}
