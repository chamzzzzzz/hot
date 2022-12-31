package csdn

import (
	"encoding/json"
	"fmt"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
)

const (
	DriverName  = "csdn"
	ProxySwitch = false
	URL         = "https://www.csdn.net"
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
		text := strings.TrimSpace(script.Text())
		if strings.HasPrefix(text, "window.__INITIAL_STATE__= ") {
			text = strings.TrimPrefix(text, "window.__INITIAL_STATE__= ")
			text = strings.Trim(text, ";")
			body := &body{}
			if err := json.Unmarshal([]byte(text), body); err != nil {
				return nil, err
			}
			for _, headline := range body.PageData.Data.WwwHeadlines {
				board.AppendTitleSummaryURL(headline.Title, headline.Description, headline.URL)
			}
			for _, headhot := range body.PageData.Data.WwwHeadhot {
				board.AppendTitleSummaryURL(headhot.Title, headhot.Description, headhot.URL)
			}
			return board, nil
		}
	}
	return nil, fmt.Errorf("not found body")
}

type body struct {
	PageData struct {
		Data struct {
			WwwHeadlines []struct {
				Description string `json:"description"`
				Title       string `json:"title"`
				URL         string `json:"url"`
			} `json:"www-Headlines"`
			WwwHeadhot []struct {
				Description string `json:"description"`
				Title       string `json:"title"`
				URL         string `json:"url"`
			} `json:"www-headhot"`
		} `json:"data"`
	} `json:"pageData"`
}
