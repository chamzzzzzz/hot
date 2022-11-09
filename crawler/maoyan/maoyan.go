package maoyan

import (
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	News    = "news"
	Actor   = "actor"
	Unknown = "unknown"
)

type Crawler struct {
	Catalog string
	Cookie  string
}

func (c *Crawler) Name() string {
	return "maoyan"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	err := c.updatecookie()
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.maoyan.com", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")
	req.Header.Set("Cookie", c.Cookie)

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
	for i, div := range dom.FindAllStrict("div", "class", "popular-container") {
		catalog := indextocatalog(i)
		if c.Catalog == "" || c.Catalog == catalog {
			for _, a := range div.FindAllStrict("a") {
				title := strings.TrimSpace(a.Text())
				url := "https://www.maoyan.com" + strings.TrimSpace(a.Attrs()["href"])
				if title == "" {
					continue
				}
				board.AppendTitleURLCatalog(title, url, catalog)
			}
		}
	}
	return board, nil
}

func (c *Crawler) updatecookie() error {
	if c.Cookie != "" {
		return nil
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest("GET", "https://www.maoyan.com/", nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	for _, cookie := range res.Cookies() {
		req.AddCookie(cookie)
	}
	c.Cookie = req.Header.Get("Cookie")
	return nil
}

func indextocatalog(i int) string {
	switch i {
	case 0:
		return News
	case 1:
		return Actor
	default:
		return Unknown
	}
}
