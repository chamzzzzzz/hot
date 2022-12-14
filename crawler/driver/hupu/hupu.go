package hupu

import (
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	Basketball = "basketball"
	Football   = "football"
	Gambia     = "gambia"
	Unknown    = "unknown"
)

const (
	DriverName = "hupu"
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
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.hupu.com", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	dom := soup.HTMLParse(string(data))
	if dom.Error != nil {
		return nil, dom.Error
	}

	board := hot.NewBoard(c.Name())
	for _, class := range classes {
		for _, a := range dom.FindAllStrict("a", "class", class) {
			div := a.FindStrict("div", "class", "hot-title")
			if div.Error != nil {
				return nil, div.Error
			}
			title := strings.TrimSpace(div.Text())
			url := strings.TrimSpace(a.Attrs()["href"])
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
