package leikeji

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strconv"
	"strings"
	"time"
)

const (
	DriverName  = "leikeji"
	ProxySwitch = false
	URL         = "https://www.leikeji.com"
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
	ul, err := dom.Find("ul", "class", "ui-sideArticleList")
	if err != nil {
		return nil, err
	}
	for _, a := range ul.QueryAll("a", "class", "link") {
		div, err := a.Find("div", "class", "time")
		if err != nil {
			return nil, err
		}
		title := strings.TrimSpace(a.Title())
		url := "https://www.leikeji.com" + strings.TrimSpace(a.Href())
		date, err := strtotime(strings.TrimSpace(div.Text()))
		if err != nil {
			return nil, err
		}
		board.Append3x1(title, "", url, date)
	}
	return board, nil
}

func strtotime(str string) (time.Time, error) {
	if date, err := time.ParseInLocation("2006-01-02", str, time.Local); err == nil {
		return date, nil
	} else {
		if strings.Contains(str, "小时前") {
			str = strings.TrimSpace(strings.Trim(str, "小时前"))
			if hour, err := strconv.Atoi(str); err != nil {
				return date, err
			} else {
				date = time.Now().Add(time.Hour * time.Duration(-hour))
				return date, nil
			}
		} else if strings.Contains(str, "天前") {
			str = strings.TrimSpace(strings.Trim(str, "天前"))
			if day, err := strconv.Atoi(str); err != nil {
				return date, err
			} else {
				date = time.Now().Add(time.Hour * 24 * time.Duration(-day))
				return date, nil
			}
		} else {
			return date, err
		}
	}
}
