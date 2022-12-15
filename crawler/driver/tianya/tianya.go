package tianya

import (
	"fmt"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
)

const (
	DriverName  = "tianya"
	ProxySwitch = false
	URL         = "https://bbs.tianya.cn/api?method=bbs.ice.getHotArticleList&params.pageSize=40&params.pageNum=1"
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
	if body.Code != "1" {
		return nil, fmt.Errorf("body code: %s", body.Code)
	}

	board := hot.NewBoard(c.Name())
	for _, row := range body.Data.Rows {
		title := strings.TrimSpace(row.Title)
		url := strings.TrimSpace(row.URL)
		board.AppendTitleURL(title, url)
	}
	return board, nil
}

type body struct {
	Success string `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Rows []struct {
			Title string `json:"title"`
			URL   string `json:"url"`
		} `json:"rows"`
	} `json:"data"`
}
