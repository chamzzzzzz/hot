package ixigua

import (
	"fmt"
	"strings"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "ixigua"
	ProxySwitch = false
	URL         = "https://i.snssdk.com/video/app/hotspot/hot_board_list?type=1"
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
	for _, data := range body.Data.List {
		url := fmt.Sprintf("https://m.ixigua.com/xigua_hot_spot/detail/%s?hotspotid=%s&enter_type=hotboard", data.ObjectID, data.ObjectID)
		board.Append(&hot.Hot{Title: strings.TrimSpace(data.Title), URL: url})
	}
	return board, nil
}

type body struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    struct {
		List []struct {
			ObjectID string `json:"object_id,omitempty"`
			Title    string `json:"title,omitempty"`
		} `json:"list,omitempty"`
	} `json:"data,omitempty"`
}

func (body *body) NormalizedCode() int {
	return body.Code
}
