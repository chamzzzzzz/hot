package github

import (
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Crawler struct {
	Proxy string
}

func (c *Crawler) Name() string {
	return "github"
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

	req, err := http.NewRequest("GET", "https://github.com/trending", nil)
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
	for _, article := range dom.FindAllStrict("article", "class", "Box-row") {
		h1 := article.FindStrict("h1", "class", "h3 lh-condensed")
		if h1.Error != nil {
			continue
		}

		a := h1.FindStrict("a")
		if a.Error != nil {
			continue
		}

		p := article.FindStrict("p", "class", "col-9 color-fg-muted my-1 pr-4")
		if p.Error != nil {
			continue
		}

		title := strings.Trim(a.Attrs()["href"], "/")
		summary := strings.Trim(p.Text(), " \n")
		url := "https://github.com/" + title
		board.AppendTitleSummaryURL(title, summary, url)
	}
	return board, nil
}
