package odaily

import (
	"encoding/json"
	"fmt"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
)

const (
	DriverName  = "odaily"
	ProxySwitch = false
	URL         = "https://www.odaily.news"
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
	for _, script := range dom.QueryAll("script") {
		_body := strings.TrimSpace(script.Text())
		idx := strings.Index(_body, "window.__INITIAL_STATE__ = ")
		if idx < 0 {
			continue
		}
		_body = _body[idx:]
		_body = strings.TrimPrefix(_body, "window.__INITIAL_STATE__ = ")
		_body = strings.Trim(_body, ";")
		body := &body{}
		if err := json.Unmarshal([]byte(_body), body); err != nil {
			return nil, err
		}
		for _, topPost := range body.Home.TopPost {
			url := fmt.Sprintf("https://www.odaily.news/post/%d", topPost.ID)
			board.AppendTitleSummaryURL(topPost.Title, topPost.Summary, url)
		}
		return board, nil
	}
	return nil, fmt.Errorf("not found body")
}

type body struct {
	Home struct {
		TopPost []struct {
			ID      int    `json:"id"`
			Title   string `json:"title"`
			Summary string `json:"summary"`
		} `json:"topPost"`
	} `json:"home"`
}
