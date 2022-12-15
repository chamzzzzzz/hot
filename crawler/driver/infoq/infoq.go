package infoq

import (
	"bytes"
	"fmt"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "infoq"
	ProxySwitch = false
	URL         = "https://www.infoq.cn/public/v1/article/getHotList"
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
	option := httputil.NewOption(c.Option, ProxySwitch)
	option.Header.Set("Origin", "https://www.infoq.cn")
	option.Header.Set("Content-Type", "application/json")
	if err := httputil.Request("POST", URL, bytes.NewReader([]byte(`{"score":null,"type":1,"size":30}`)), "json", body, option); err != nil {
		return nil, err
	} else if body.Code != 0 {
		return nil, fmt.Errorf("body code: %d", body.Code)
	}

	board := hot.NewBoard(c.Name())
	for _, data := range body.Data {
		url := "https://www.infoq.cn/article/" + data.UUID
		board.AppendTitleSummaryURL(data.ArticleTitle, data.ArticleSummary, url)
	}
	return board, nil
}

type body struct {
	Code int `json:"code"`
	Data []struct {
		ArticleSummary string `json:"article_summary"`
		ArticleTitle   string `json:"article_title"`
		UUID           string `json:"uuid"`
	} `json:"data"`
}
