package rfa

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
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

type Crawler struct {
	BoardName string
	Proxy     string
}

func (c *Crawler) Name() string {
	switch c.BoardName {
	case Cantonese:
		return "rfa_x_cantonese"
	case English:
		return "rfa_x_english"
	case Mandarin:
		return "rfa"
	default:
		return "rfa"
	}
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	switch c.BoardName {
	case Cantonese:
		return c.rfa_x_cantonese()
	case English:
		return c.rfa_x_english()
	case Mandarin:
		return c.rfa()
	default:
		return c.rfa()
	}
}

func (c *Crawler) rfa_x_cantonese() (*hot.Board, error) {
	return c.rfaWithLanguage("cantonese")
}

func (c *Crawler) rfa_x_english() (*hot.Board, error) {
	return c.rfaWithLanguage("english")
}

func (c *Crawler) rfa() (*hot.Board, error) {
	return c.rfaWithLanguage("mandarin")
}

func (c *Crawler) rfaWithLanguage(language string) (*hot.Board, error) {
	client := &http.Client{}
	if c.Proxy != "" {
		if proxyUrl, err := url.Parse(c.Proxy); err == nil {
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
			board.AppendTitleURL(title, url)
		}
	}
	return board, nil
}
