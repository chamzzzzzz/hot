package semiunion

import (
	"fmt"
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
	return "semiunion"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://www.semiunion.com/insight/", nil)
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
	for _, li := range dom.FindAllStrict("li", "class", "each-news") {
		div := li.Find("div", "class", "name")
		if div.Error != nil {
			return nil, div.Error
		}
		a := div.Find("a")
		if a.Error != nil {
			return nil, a.Error
		}
		div2 := li.Find("div", "class", "desc")
		if div2.Error != nil {
			return nil, div2.Error
		}
		title := fmt.Sprintf("%s|%s", strings.TrimSpace(a.Text()), strings.TrimSpace(div2.Text()))
		summary := a.Attrs()["href"]
		board.Append(title, summary, date)
	}
	return board, nil
}
