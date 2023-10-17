package pojie52

import (
	"fmt"
	"strings"
	"time"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "pojie52"
	ProxySwitch = false
	URL         = "https://www.52pojie.cn/forum.php?mod=guide&view=hot"
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
	div, err := dom.Find("div", "id", "threadlist")
	if err != nil {
		return nil, err
	}
	div, err = div.Find("div", "class", "bm_c")
	if err != nil {
		return nil, err
	}
	for _, tbody := range div.QueryAll("tbody") {
		a, err := tbody.Find("a", "class", "xst")
		if err != nil {
			return nil, err
		}
		td := tbody.QueryAll("td", "class", "by")
		if len(td) != 3 {
			return nil, fmt.Errorf("td count invalid")
		}
		span, err := td[1].Find("span")
		if err != nil {
			return nil, err
		}

		title := strings.TrimSpace(a.Text())
		url := "https://www.52pojie.cn/" + strings.TrimSpace(a.Href())
		date, err := time.ParseInLocation("2006-1-2 15:04", strings.TrimSpace(span.Text()), time.Local)
		if err != nil {
			return nil, err
		}
		board.Append(&hot.Hot{Title: title, URL: url, PublishDate: date})
	}
	return board, nil
}
