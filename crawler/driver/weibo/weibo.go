package weibo

import (
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	DriverName = "weibo"
	Cookie     = "WBPSESS=durPiJxsbzq5XDaI2wW0NxQldYOrBwQzLVlPfvpcy3mQ3XQonV49sfubFlqvuI_rBrarQ2dZHLfrOVaRKnvrm9130Jsv26CGHwu2LjHl3RrnHDHKIUtZPYEi9qKk6n-K; SUB=_2AkMU1LJTf8NxqwJRmPAQymrhaYl_yg_EieKiiEOIJRMxHRl-yT92qkI6tRB6P1ScvMDt8JtdZqvVJlNftBcRg-WjvODv; SUBP=0033WrSXqPxfM72-Ws9jqgMF55529P9D9WFSWJI0b_93sKJGpCc_.aOL; XSRF-TOKEN=Z3qrKi3V9M0TVao6eTMMmpRC"
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
	cookie string
}

func (c *Crawler) Driver() driver.Driver {
	return &Driver{}
}

func (c *Crawler) Name() string {
	return DriverName
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	err := c.updatecookie()
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://s.weibo.com/top/summary", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")
	req.Header.Set("Cookie", c.cookie)

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

	div := dom.FindStrict("div", "id", "pl_top_realtimehot")
	if div.Error != nil {
		return nil, div.Error
	}

	board := hot.NewBoard(c.Name())
	for _, tr := range div.FindAllStrict("tr", "class", "") {
		td01 := tr.Find("td", "class", "td-01")
		if td01.Error != nil {
			return nil, td01.Error
		}
		if _, err := strconv.Atoi(td01.Text()); err != nil {
			continue
		}

		td02 := tr.Find("td", "class", "td-02")
		if td02.Error != nil {
			return nil, td02.Error
		}
		a := td02.Find("a")
		if a.Error != nil {
			return nil, a.Error
		}
		board.Append1(a.Text())
	}
	return board, nil
}

func (c *Crawler) updatecookie() error {
	if c.cookie != "" {
		return nil
	}
	if c.Option.Cookie != "" {
		c.cookie = c.Option.Cookie
		return nil
	}
	c.cookie = Cookie
	return nil
}
