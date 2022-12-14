package rfa

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	Cantonese = "cantonese"
	English   = "english"
	Mandarin  = "mandarin"
)

const (
	DriverName = "rfa"
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
	case Cantonese:
		return c.cantonese()
	case English:
		return c.english()
	case Mandarin:
		return c.mandarin()
	default:
		return c.all()
	}
}

func (c *Crawler) all() (*hot.Board, error) {
	board := hot.NewBoard(c.Name())
	if b, err := c.cantonese(); err != nil {
		return nil, err
	} else {
		for _, hot := range b.Hots {
			board.Append0(hot)
		}
	}

	if b, err := c.english(); err != nil {
		return nil, err
	} else {
		for _, hot := range b.Hots {
			board.Append0(hot)
		}
	}

	if b, err := c.mandarin(); err != nil {
		return nil, err
	} else {
		for _, hot := range b.Hots {
			board.Append0(hot)
		}
	}
	return board, nil
}

func (c *Crawler) cantonese() (*hot.Board, error) {
	return c.withLanguage(Cantonese)
}

func (c *Crawler) english() (*hot.Board, error) {
	return c.withLanguage(English)
}

func (c *Crawler) mandarin() (*hot.Board, error) {
	return c.withLanguage(Mandarin)
}

func (c *Crawler) withLanguage(language string) (*hot.Board, error) {
	client := &http.Client{}
	if c.Option.Proxy != "" {
		if proxyUrl, err := url.Parse(c.Option.Proxy); err == nil {
			client.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			}
		}
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.rfa.org/%s", language), nil)
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
	for _, div := range dom.FindAllStrict("div", "class", "most_read_only first_most_bold") {
		for _, li := range div.FindAll("li") {
			a := li.Find("a")
			if a.Error != nil {
				return nil, a.Error
			}
			title := strings.TrimSpace(a.Text())
			url := strings.TrimSpace(a.Attrs()["href"])
			catalog := language
			board.Append4(title, "", url, catalog)
		}
	}
	return board, nil
}
