package timecom

import (
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Crawler struct {
	Proxy string
}

func (c *Crawler) Name() string {
	return "timecom"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	if c.Proxy != "" {
		if proxyUrl, err := url.Parse(c.Proxy); err == nil {
			client.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			}
		}
	}

	req, err := http.NewRequest("GET", "https://time.com/", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	html, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	dom := soup.HTMLParse(string(html))
	if dom.Error != nil {
		return nil, dom.Error
	}

	board := hot.NewBoard(c.Name())
	date := time.Now()
	for _, li := range dom.FindAllStrict("li", "class", "most-popular-feed__item") {
		for i, a := range li.FindAll("a") {
			if i == 1 {
				h := a.Find("h3")
				if h.Error != nil {
					return nil, h.Error
				}
				title := strings.TrimSpace(h.Text())
				summary := a.Attrs()["href"]
				board.Append(title, summary, date)
			}
		}
	}
	return board, nil
}
