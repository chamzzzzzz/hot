package donews

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
	return "donews"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.donews.com/newsflash/index", nil)
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
	for _, div := range dom.FindAllStrict("div", "class", "search-articles-item") {
		a := div.FindStrict("a")
		if a.Error != nil {
			return nil, a.Error
		}
		div2 := a.FindStrict("div", "class", "search-articles-content")
		if div2.Error != nil {
			return nil, div2.Error
		}
		title := strings.TrimSpace(div2.Text())
		summary := a.Attrs()["href"]
		board.Append(title, summary, date)
	}
	return board, nil
}
