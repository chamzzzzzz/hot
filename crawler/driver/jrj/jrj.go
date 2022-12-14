package jrj

import (
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	Finance = "finance"
	Tech    = "tech"
	House   = "house"
)

const (
	DriverName = "jrj"
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
	case Finance:
		return c.finance()
	case Tech:
		return c.tech()
	case House:
		return c.house()
	default:
		return c.all()
	}
}

func (c *Crawler) all() (*hot.Board, error) {
	board := hot.NewBoard(c.Name())
	if b, err := c.finance(); err != nil {
		return nil, err
	} else {
		for _, hot := range b.Hots {
			board.Append0(hot)
		}
	}

	if b, err := c.tech(); err != nil {
		return nil, err
	} else {
		for _, hot := range b.Hots {
			board.Append0(hot)
		}
	}

	if b, err := c.house(); err != nil {
		return nil, err
	} else {
		for _, hot := range b.Hots {
			board.Append0(hot)
		}
	}
	return board, nil
}

func (c *Crawler) tech() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://tech.jrj.com.cn", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(transform.NewReader(res.Body, simplifiedchinese.GBK.NewDecoder()))
	if err != nil {
		return nil, err
	}

	dom := soup.HTMLParse(string(data))
	if dom.Error != nil {
		return nil, dom.Error
	}

	board := hot.NewBoard(c.Name())
	div := dom.FindStrict("div", "class", "hotart hotnews")
	if div.Error != nil {
		return nil, div.Error
	}
	for _, a := range div.FindAllStrict("a") {
		title := strings.TrimSpace(a.Text())
		url := strings.TrimSpace(a.Attrs()["href"])
		board.Append4(title, "", url, Tech)
	}
	return board, nil
}

func (c *Crawler) house() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://house.jrj.com.cn", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(transform.NewReader(res.Body, simplifiedchinese.GBK.NewDecoder()))
	if err != nil {
		return nil, err
	}

	dom := soup.HTMLParse(string(data))
	if dom.Error != nil {
		return nil, dom.Error
	}

	board := hot.NewBoard(c.Name())
	div := dom.FindStrict("div", "class", "hotart")
	if div.Error != nil {
		return nil, div.Error
	}
	for _, a := range div.FindAllStrict("a") {
		title := strings.TrimSpace(a.Text())
		url := strings.TrimSpace(a.Attrs()["href"])
		board.Append4(title, "", url, House)
	}
	return board, nil
}

func (c *Crawler) finance() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://finance.jrj.com.cn/list/industrynews.shtml", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(transform.NewReader(res.Body, simplifiedchinese.GBK.NewDecoder()))
	if err != nil {
		return nil, err
	}

	dom := soup.HTMLParse(string(data))
	if dom.Error != nil {
		return nil, dom.Error
	}

	board := hot.NewBoard(c.Name())
	for _, id := range []string{"con_1", "con_2"} {
		ul := dom.Find("ul", "id", id)
		if ul.Error != nil {
			return nil, ul.Error
		}
		for _, a := range ul.FindAll("a") {
			title := strings.TrimSpace(a.Text())
			url := strings.TrimSpace(a.Attrs()["href"])
			board.Append4(title, "", url, Finance)
		}
	}
	return board, nil
}
