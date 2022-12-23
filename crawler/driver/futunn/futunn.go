package futunn

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
)

const (
	DriverName  = "futunn"
	ProxySwitch = false
	URL         = "https://news.futunn.com/node/api/hot-news/get-info"
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
	token  string
}

func (c *Crawler) Driver() driver.Driver {
	return &Driver{}
}

func (c *Crawler) Name() string {
	return DriverName
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	if err := c.updatetoken(); err != nil {
		return nil, err
	}

	body := &body{}
	option := httputil.NewOption(c.Option, ProxySwitch)
	option.Header.Set("futu-x-csrf-token-v2", c.token)
	if err := httputil.Request("GET", URL, nil, "json", body, option); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	for _, data := range body.Data.HotNews {
		board.AppendTitleURL(strings.TrimSpace(data.Title), strings.TrimSpace((data.URL)))
	}
	return board, nil
}

func (c *Crawler) updatetoken() error {
	_, cookies, err := httputil.RequestCookie("GET", "https://news.futunn.com/node/hot-news", nil, httputil.NewOption(c.Option, ProxySwitch))
	if err != nil {
		return err
	}
	for _, cookie := range cookies {
		if cookie.Name == "_news_csrf_token" {
			c.token = cookie.Value
			return nil
		}
	}
	return nil
}

type body struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		HotNews []struct {
			Title string `json:"title"`
			URL   string `json:"url"`
		} `json:"hotNews"`
		UpdatedAt string `json:"updatedAt"`
	} `json:"data"`
}

func (body *body) NormalizedCode() int {
	return body.Code
}
