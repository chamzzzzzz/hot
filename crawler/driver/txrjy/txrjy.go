package txrjy

import (
	"strings"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "txrjy"
	ProxySwitch = false
	URL         = "https://www.txrjy.com/forum.php"
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
	for _, id := range []string{"review", "day"} {
		for _, a := range dom.Query("div", "id", id).QueryAll("a") {
			title := strings.TrimSpace(a.Title())
			if title == "" {
				title = strings.TrimSpace(a.Text())
			}
			url := "https://www.txrjy.com/" + strings.TrimSpace(a.Href())
			board.Append(&hot.Hot{Title: title, URL: url})
		}
	}
	return board, nil
}
