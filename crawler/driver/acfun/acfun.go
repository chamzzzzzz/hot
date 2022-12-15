package acfun

import (
	"fmt"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
	"time"
)

const (
	DriverName  = "acfun"
	ProxySwitch = false
	URL         = "https://www.acfun.cn/rest/pc-direct/rank/channel?channelId=&subChannelId=&rankLimit=30&rankPeriod=DAY"
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
	if body.Result != 0 {
		return nil, fmt.Errorf("body result: %d", body.Result)
	}

	board := hot.NewBoard(c.Name())
	for _, data := range body.RankList {
		title := strings.TrimSpace(data.ContentTitle)
		summary := strings.TrimSpace(data.ContentDesc)
		url := strings.TrimSpace(data.ShareURL)
		date := time.UnixMilli(data.CreateTimeMillis)
		board.Append3x1(title, summary, url, date)
	}
	return board, nil
}

type body struct {
	Result   int `json:"result"`
	RankList []struct {
		ContentID        int    `json:"contentId"`
		ContributeTime   int64  `json:"contributeTime"`
		ContentTitle     string `json:"contentTitle"`
		ContentDesc      string `json:"contentDesc"`
		CreateTimeMillis int64  `json:"createTimeMillis"`
		ShareURL         string `json:"shareUrl"`
		Title            string `json:"title"`
	} `json:"rankList"`
}
