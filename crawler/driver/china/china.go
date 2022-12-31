package china

import (
	"fmt"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
	"time"
)

const (
	Auto    = "auto"
	Finance = "finance"
	News    = "news"
	Game    = "game"
	Ent     = "ent"
	Mili    = "mili"
)

const (
	DriverName  = "china"
	ProxySwitch = false
	URL         = "https://auto.china.com"
	AutoURL     = URL
	FinanceURL  = "https://finance.china.com"
	RankURL     = "https://rank.china.com/rank/cms/%s/day/rank_all.js"
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
	case Auto:
		return c.auto(nil)
	case Finance:
		return c.finance(nil)
	case News, Game, Ent, Mili:
		return c.rank(nil, c.Option.Catalog)
	default:
		return c.all()
	}
}

func (c *Crawler) all() (*hot.Board, error) {
	board := hot.NewBoard(c.Name())
	if _, err := c.auto(board); err != nil {
		return nil, err
	}
	if _, err := c.finance(board); err != nil {
		return nil, err
	}
	for _, catalog := range []string{News, Game, Ent, Mili} {
		if _, err := c.rank(board, catalog); err != nil {
			return nil, err
		}
	}
	return board, nil
}

func (c *Crawler) auto(board *hot.Board) (*hot.Board, error) {
	dom := &httputil.DOM{}
	if err := httputil.Request("GET", AutoURL, nil, "dom", dom, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	}

	if board == nil {
		board = hot.NewBoard(c.Name())
	}
	div, err := dom.Find("div", "class", "auto-rank mt50")
	if err != nil {
		return nil, err
	}
	for _, a := range div.QueryAll("a") {
		span, err := a.Find("span")
		if err != nil {
			return nil, err
		}
		title := strings.TrimSpace(span.Text())
		url := AutoURL + strings.TrimSpace(a.Href())
		board.Append4(title, "", url, Auto)
	}
	return board, nil
}

func (c *Crawler) finance(board *hot.Board) (*hot.Board, error) {
	dom := &httputil.DOM{}
	if err := httputil.Request("GET", FinanceURL, nil, "dom", dom, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	}

	if board == nil {
		board = hot.NewBoard(c.Name())
	}
	for _, class := range []string{"col-l", "col-r"} {
		div, err := dom.Find("div", "class", class)
		if err != nil {
			return nil, err
		}
		for _, li := range div.QueryAll("li") {
			a, err := li.Find("a")
			if err != nil {
				return nil, err
			}
			title := strings.TrimSpace(a.Text())
			url := strings.TrimSpace(a.Href())
			board.Append4(title, "", url, Finance)
		}
	}
	return board, nil
}

func (c *Crawler) rank(board *hot.Board, catalog string) (*hot.Board, error) {
	body := &body{}
	option := httputil.NewOption(c.Option, ProxySwitch)
	option.TrimPrefix = "var day_top="
	if err := httputil.Request("GET", fmt.Sprintf(RankURL, catalog), nil, "json", body, option); err != nil {
		return nil, err
	}

	if board == nil {
		board = hot.NewBoard(c.Name())
	}
	for _, data := range body.List {
		title := strings.TrimSpace(data.Title)
		summary := strings.TrimSpace(data.Summary)
		url := strings.TrimSpace(data.URL)
		date, err := time.ParseInLocation("2006-01-02 15:04:05", data.CreateTime, time.Local)
		if err != nil {
			return nil, err
		}
		board.Append4x1(title, summary, url, catalog, date)
	}
	return board, nil
}

type body struct {
	Name string `json:"name"`
	List []struct {
		Summary    string `json:"summary"`
		Title      string `json:"title"`
		URL        string `json:"url"`
		CreateTime string `json:"createTime"`
	} `json:"list"`
}
