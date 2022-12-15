package kr36

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
)

const (
	RenQi    = "renqi"
	ZongHe   = "zonghe"
	ShouCang = "shoucang"
	Unknown  = "unknown"
)

const (
	DriverName  = "kr36"
	ProxySwitch = false
	URL         = "https://36kr.com/hot-list/catalog"
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
	for i, div := range dom.FindAllStrict("div", "class", "list-section-wrapper") {
		catalog := itocatalog(i)
		if c.Option.Catalog != "" && c.Option.Catalog != catalog {
			continue
		}

		for _, div := range div.FindAll("div", "class", "article-wrapper") {
			a1 := div.FindStrict("a", "class", "article-item-title weight-bold")
			if a1.Error != nil {
				return nil, a1.Error
			}
			a2 := div.FindStrict("a", "class", "article-item-description ellipsis-2")
			if a2.Error != nil {
				return nil, a2.Error
			}
			title := strings.TrimSpace(a1.Text())
			summary := strings.TrimSpace(a2.Text())
			url := "https://36kr.com" + strings.TrimSpace(a1.Attrs()["href"])
			board.Append4(title, summary, url, catalog)
		}
	}
	return board, nil
}

func itocatalog(i int) string {
	switch i {
	case 0:
		return RenQi
	case 1:
		return ZongHe
	case 2:
		return ShouCang
	default:
		return Unknown
	}
}
