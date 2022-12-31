package hupu

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
)

const (
	Basketball = "basketball"
	Football   = "football"
	Gambia     = "gambia"
	Unknown    = "unknown"
)

const (
	DriverName  = "hupu"
	ProxySwitch = false
	URL         = "https://www.hupu.com"
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
	switch c.Option.Catalog {
	case Basketball:
		return c.basketball()
	case Football:
		return c.football()
	case Gambia:
		return c.gambia()
	default:
		return c.all()
	}
}

func (c *Crawler) basketball() (*hot.Board, error) {
	return c.withClasses("newpcbasketball")
}

func (c *Crawler) football() (*hot.Board, error) {
	return c.withClasses("newpcsoccer")
}

func (c *Crawler) gambia() (*hot.Board, error) {
	return c.withClasses("newpcbbs")
}

func (c *Crawler) all() (*hot.Board, error) {
	return c.withClasses("newpcbasketball", "newpcsoccer", "newpcbbs")
}

func (c *Crawler) withClasses(classes ...string) (*hot.Board, error) {
	dom := &httputil.DOM{}
	if err := httputil.Request("GET", URL, nil, "dom", dom, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	for _, class := range classes {
		for _, a := range dom.QueryAll("a", "class", class) {
			div, err := a.Find("div", "class", "hot-title")
			if err != nil {
				return nil, err
			}
			title := strings.TrimSpace(div.Text())
			url := strings.TrimSpace(a.Href())
			catalog := classtocatalog(class)
			board.Append4(title, "", url, catalog)
		}
	}
	return board, nil
}

func classtocatalog(class string) string {
	switch class {
	case "newpcbasketball":
		return Basketball
	case "newpcsoccer":
		return Football
	case "newpcbbs":
		return Gambia
	default:
		return Unknown
	}
}
