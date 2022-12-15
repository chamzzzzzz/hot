package ifeng

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
)

const (
	DriverName  = "ifeng"
	ProxySwitch = false
	URL         = "https://www.ifeng.com/"
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
	div := dom.Find("div", "class", "hot_box-1yXFLW7e")
	if div.Error != nil {
		return nil, div.Error
	}
	h3 := div.FindStrict("h3", "class", "list_title-nOvTJ00k big-3QPOjsEI")
	if h3.Error != nil {
		return nil, div.Error
	}
	a := h3.Find("a")
	if a.Error != nil {
		return nil, div.Error
	}
	title := strings.TrimSpace(a.Text())
	url := strings.TrimSpace(a.Attrs()["href"])
	board.AppendTitleURL(title, url)
	for _, p := range div.FindAllStrict("p", "class", "news_list_p-3EcL2Tvk ") {
		a = p.Find("a")
		if a.Error != nil {
			return nil, div.Error
		}
		title = strings.TrimSpace(a.Text())
		url = strings.TrimSpace(a.Attrs()["href"])
		board.AppendTitleURL(title, url)
	}
	return board, nil
}
