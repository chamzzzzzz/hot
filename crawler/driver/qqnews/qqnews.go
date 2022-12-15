package qqnews

import (
	"fmt"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"time"
)

const (
	DriverName  = "qqnews"
	ProxySwitch = false
	URL         = `https://i.news.qq.com/trpc.qqnews_web.kv_srv.kv_srv_http_proxy/list?sub_srv_id=24hours&srv_id=pc&offset=0&limit=20&strategy=1&ext={%22pool%22:[%22top%22],%22is_filter%22:7,%22check_type%22:true}`
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
	} else if body.Ret != 0 {
		return nil, fmt.Errorf("body ret: %d", body.Ret)
	}

	board := hot.NewBoard(c.Name())
	for _, data := range body.Data.List {
		date, err := time.ParseInLocation("2006-01-02 15:04:05", data.PublishTime, time.Local)
		if err != nil {
			return nil, err
		}
		board.Append3x1(data.Title, "", data.URL, date)
	}
	return board, nil
}

type body struct {
	Ret  int    `json:"ret"`
	Msg  string `json:"msg"`
	Data struct {
		List []struct {
			Title       string `json:"title"`
			URL         string `json:"url"`
			PublishTime string `json:"publish_time"`
		} `json:"list"`
	} `json:"data"`
}
