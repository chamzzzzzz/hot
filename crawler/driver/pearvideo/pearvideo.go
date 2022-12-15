package pearvideo

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
)

const (
	DriverName  = "pearvideo"
	ProxySwitch = false
	URL         = "https://www.pearvideo.com/userlist_loading.jsp?reqType=1"
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
	for _, li := range dom.FindAllStrict("li", "class", "columns-contem") {
		a := li.Find("a", "class", "actplay")
		if a.Error != nil {
			return nil, a.Error
		}
		div := a.Find("div", "class", "columnsem-title")
		if div.Error != nil {
			return nil, div.Error
		}
		title := strings.TrimSpace(div.Text())
		url := "https://www.pearvideo.com/" + strings.TrimSpace(a.Attrs()["href"])
		board.AppendTitleURL(title, url)
	}
	return board, nil
}
