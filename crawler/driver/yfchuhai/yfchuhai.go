package yfchuhai

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strconv"
	"strings"
	"time"
)

const (
	DriverName  = "yfchuhai"
	ProxySwitch = false
	URL         = "https://www.yfchuhai.com"
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
	for _, a := range dom.FindAllStrict("a", "class", "list__item flex") {
		h2 := a.Find("h2")
		if h2.Error != nil {
			return nil, h2.Error
		}
		span := a.Find("span", "class", "text")
		if span.Error != nil {
			return nil, span.Error
		}
		title := strings.TrimSpace(h2.Text())
		url := "https://www.yfchuhai.com" + strings.TrimSpace(a.Attrs()["href"])
		date, err := strtotime(strings.TrimSpace(span.Text()))
		if err != nil {
			return nil, err
		}
		board.AppendTitleURLDate(title, url, date)
	}
	return board, nil
}

func strtotime(str string) (time.Time, error) {
	if date, err := time.ParseInLocation("01-02 15:04", str, time.Local); err == nil {
		date = date.AddDate(time.Now().Year(), 0, 0)
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
		} else if strings.Contains(str, "前天") {
			str = strings.TrimSpace(strings.Trim(str, "前天"))
			if date, err := time.ParseInLocation("15:04", str, time.Local); err != nil {
				return date, err
			} else {
				date = date.AddDate(time.Now().Year(), int(time.Now().Month())-1, time.Now().Day()-1)
				date = date.Add(time.Hour * 24 * time.Duration(-2))
				return date, nil
			}
		} else if strings.Contains(str, "昨天") {
			str = strings.TrimSpace(strings.Trim(str, "昨天"))
			if date, err := time.ParseInLocation("15:04", str, time.Local); err != nil {
				return date, err
			} else {
				date = date.AddDate(time.Now().Year(), int(time.Now().Month())-1, time.Now().Day()-1)
				date = date.Add(time.Hour * 24 * time.Duration(-1))
				return date, nil
			}
		} else {
			return date, err
		}
	}
}
