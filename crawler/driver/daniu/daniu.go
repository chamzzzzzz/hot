package daniu

import (
	"fmt"
	"strings"
	"time"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "daniu"
	ProxySwitch = false
	URL         = "https://www.4330.cn/misc.php?mod=ranklist&type=thread&view=heats&orderby=today"
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
	option := httputil.NewOption(c.Option, ProxySwitch)
	option.DetectContentEncoding = true
	if err := httputil.Request("GET", URL, nil, "dom", dom, option); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	div, err := dom.Find("div", "class", "tl")
	if err != nil {
		return nil, err
	}
	tbody := div.QueryAll("tbody")
	if len(tbody) != 2 {
		return nil, fmt.Errorf("tbody count invalid")
	}
	for _, tr := range tbody[1].QueryAll("tr") {
		th, err := tr.Find("th")
		if err != nil {
			return nil, err
		}
		em, err := tr.Find("em")
		if err != nil {
			return nil, err
		}
		a, err := th.Find("a")
		if err != nil {
			return nil, err
		}

		title := strings.TrimSpace(a.Text())
		url := "https://www.4330.cn/" + strings.TrimSpace(a.Href())
		date, err := time.ParseInLocation("2006-1-2 15:04", strings.TrimSpace(em.Text()), time.Local)
		if err != nil {
			return nil, err
		}
		board.Append(&hot.Hot{Title: title, URL: url, PublishDate: date})
	}
	return board, nil
}
