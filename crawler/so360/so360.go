package so360

import (
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Crawler struct {
}

func (c *Crawler) Name() string {
	return "so360"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://news.so.com/hotnews", nil)
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
	div := dom.Find("div", "class", "hotnews-main")
	if div.Error != nil {
		return nil, div.Error
	}
	for _, li := range div.FindAll("li") {
		for _, a := range li.FindAll("a") {
			span := a.FindStrict("span", "class", "title")
			if span.Error != nil {
				return nil, span.Error
			}
			title := strings.TrimSpace(span.Text())
			board.Append(title, "", date)
		}
	}
	return board, nil
}