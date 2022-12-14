package mddcloud

import (
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	DriverName = "mddcloud"
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
	req, err := http.NewRequest("GET", "https://www.mddcloud.com.cn/", nil)
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
	for _, li := range dom.FindAllStrict("li", "class", "rank-list-item") {
		a := li.FindStrict("a", "class", "g-a-block")
		p1 := li.FindStrict("p", "class", "rank-vod-title")
		p2 := li.FindStrict("p", "class", "rank-vod-desc")
		if a.Error != nil {
			return nil, a.Error
		}
		if p1.Error != nil {
			return nil, p1.Error
		}
		if p2.Error != nil {
			return nil, p2.Error
		}
		title := strings.TrimSpace(p1.Text())
		summary := strings.TrimSpace(p2.Text())
		url := "https://www.mddcloud.com.cn" + strings.TrimSpace(a.Attrs()["href"])
		board.AppendTitleSummaryURL(title, summary, url)
	}
	return board, nil
}
