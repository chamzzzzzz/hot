package qqnews

import (
	"fmt"
	"time"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "qqnews"
	ProxySwitch = false
	URL         = `https://r.inews.qq.com/gw/event/hot_ranking_list?offset=0&page_size=50`
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
	if len(body.Idlist) > 0 {
		for _, data := range body.Idlist[0].Newslist {
			if data.URL == "" {
				continue
			}
			date, err := time.ParseInLocation("2006-01-02 15:04:05", data.Time, time.Local)
			if err != nil {
				return nil, err
			}
			board.Append(&hot.Hot{Title: data.Title, URL: data.URL, PublishDate: date})
		}
	}
	return board, nil
}

type body struct {
	Ret    int `json:"ret"`
	Idlist []struct {
		Newslist []struct {
			Title string `json:"title"`
			URL   string `json:"url,omitempty"`
			Time  string `json:"time,omitempty"`
		} `json:"newslist"`
	} `json:"idlist"`
}
