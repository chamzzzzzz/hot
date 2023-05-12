package cninfo

import (
	"strings"
	"time"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "cninfo"
	ProxySwitch = false
	URL         = "http://www.cninfo.com.cn/new/index/getHotAnnouces"
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
	var body []body
	if err := httputil.Request("GET", URL, nil, "json", &body, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	for _, data := range body {
		title := strings.TrimSpace(data.ShortName)
		summary := strings.TrimSpace(data.Name)
		url := "http://static.cninfo.com.cn/" + strings.TrimSpace(data.URL)
		date := time.UnixMilli(data.AnnouncementTime)
		board.Append(&hot.Hot{Title: title, Summary: summary, URL: url, PublishDate: date})
	}
	return board, nil
}

type body struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	URL              string `json:"url"`
	ShortName        string `json:"shortName"`
	AnnouncementTime int64  `json:"announcementTime"`
}
