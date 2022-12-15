package pojie52

import (
	"fmt"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
	"time"
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
	option.ContentEncoding = "gbk"
	if err := httputil.Request("GET", URL, nil, "dom", dom, option); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	div := dom.FindStrict("div", "id", "threadlist")
	if div.Error != nil {
		return nil, div.Error
	}
	div = div.Find("div", "class", "bm_c")
	if div.Error != nil {
		return nil, div.Error
	}
	for _, tbody := range div.FindAll("tbody") {
		a := tbody.Find("a", "class", "xst")
		if a.Error != nil {
			return nil, a.Error
		}
		td := tbody.FindAll("td", "class", "by")
		if len(td) != 3 {
			return nil, fmt.Errorf("td count invalid")
		}
		span := td[1].Find("span")
		if span.Error != nil {
			return nil, span.Error
		}

		title := strings.TrimSpace(a.Text())
		url := "https://www.52pojie.cn/" + strings.TrimSpace(a.Attrs()["href"])
		date, err := time.ParseInLocation("2006-1-2 15:04", strings.TrimSpace(span.Text()), time.Local)
		if err != nil {
			return nil, err
		}
		board.Append3x1(title, "", url, date)
	}
	return board, nil
}
