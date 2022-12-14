package leikeji

import (
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	DriverName = "leikeji"
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
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.leikeji.com", nil)
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
	ul := dom.FindStrict("ul", "class", "ui-sideArticleList")
	if ul.Error != nil {
		return nil, ul.Error
	}
	for _, a := range ul.FindAllStrict("a", "class", "link") {
		div := a.Find("div", "class", "time")
		if div.Error != nil {
			return nil, div.Error
		}
		title := strings.TrimSpace(a.Attrs()["title"])
		url := "https://www.leikeji.com" + strings.TrimSpace(a.Attrs()["href"])
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
