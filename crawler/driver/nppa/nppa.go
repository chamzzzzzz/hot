package nppa

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
	"time"
)

const (
	DriverName  = "nppa"
	ProxySwitch = false
	URL         = "https://www.nppa.gov.cn/nppa/channels/718.shtml"
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
	ul := dom.FindStrict("ul", "class", "m2nrul")
	if ul.Error != nil {
		return nil, ul.Error
	}
	for _, li := range ul.FindAllStrict("li") {
		a := li.Find("a")
		if a.Error != nil {
			return nil, a.Error
		}
		span := li.Find("span")
		if span.Error != nil {
			return nil, span.Error
		}
		title := strings.TrimSpace(a.Text())
		url := "https://www.nppa.gov.cn" + strings.TrimSpace(a.Attrs()["href"])
		date, err := time.ParseInLocation("[2006-01-02]", strings.TrimSpace(span.Text()), time.Local)
		if err != nil {
			return nil, err
		}
		board.AppendTitleURLDate(title, url, date)
	}
	return board, nil
}
