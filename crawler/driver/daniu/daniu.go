package daniu

import (
	"fmt"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
	"time"
)

const (
	DriverName  = "daniu"
	ProxySwitch = false
	URL         = "https://www.daniu523.com/misc.php?mod=ranklist&type=thread&view=heats&orderby=today"
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
	option.ContentEncoding = "gbk"
	if err := httputil.Request("GET", URL, nil, "dom", dom, option); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	div := dom.FindStrict("div", "class", "tl")
	if div.Error != nil {
		return nil, div.Error
	}
	tbody := div.FindAllStrict("tbody")
	if len(tbody) != 2 {
		return nil, fmt.Errorf("tbody count invalid")
	}
	for _, tr := range tbody[1].FindAllStrict("tr") {
		th := tr.Find("th")
		if th.Error != nil {
			return nil, th.Error
		}
		em := tr.FindStrict("em")
		if em.Error != nil {
			fmt.Println(tr.FullText())
			return nil, em.Error
		}
		a := th.Find("a")
		if a.Error != nil {
			return nil, a.Error
		}

		title := strings.TrimSpace(a.Text())
		url := "https://www.daniu523.com/" + strings.TrimSpace(a.Attrs()["href"])
		date, err := time.ParseInLocation("2006-1-2 15:04", strings.TrimSpace(em.Text()), time.Local)
		if err != nil {
			return nil, err
		}
		board.Append3x1(title, "", url, date)
	}
	return board, nil
}
