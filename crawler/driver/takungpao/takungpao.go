package takungpao

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
	"time"
)

const (
	DriverName  = "takungpao"
	ProxySwitch = false
	URL         = "https://www.takungpao.com"
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
	dom := &httputil.DOM{}
	if err := httputil.Request("GET", URL, nil, "dom", dom, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	div := dom.Find("div", "class", "ranking")
	if div.Error != nil {
		return nil, div.Error
	}
	for _, dl := range div.FindAllStrict("dl") {
		a := dl.Find("a")
		if a.Error != nil {
			return nil, a.Error
		}
		dd := dl.Find("dd", "class", "time")
		if dd.Error != nil {
			return nil, dd.Error
		}

		title := strings.TrimSpace(a.Text())
		url := strings.TrimSpace(a.Attrs()["href"])
		date, err := time.ParseInLocation("2006-01-02 15:04:05", strings.TrimSpace(dd.Text()), time.Local)
		if err != nil {
			return nil, err
		}
		board.Append3x1(title, "", url, date)
	}
	return board, nil
}
