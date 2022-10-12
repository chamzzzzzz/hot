package sputniknews

import (
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Crawler struct {
}

func (c *Crawler) Name() string {
	return "sputniknews"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.sputniknews.cn", nil)
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
	div := dom.Find("div", "data-section", "3")
	if div.Error != nil {
		return nil, div.Error
	}
	for _, a := range div.FindAllStrict("a", "class", "cell-list__item m-no-image") {
		span := a.Find("span", "class", "cell__controls-date")
		if span.Error != nil {
			return nil, span.Error
		}
		span = span.Find("span")
		if span.Error != nil {
			return nil, span.Error
		}

		timestamp, err := strconv.ParseInt(span.Attrs()["data-unixtime"], 10, 64)
		if err != nil {
			return nil, err
		}

		title := strings.TrimSpace(a.Attrs()["title"])
		url := "https://www.sputniknews.cn" + strings.TrimSpace(a.Attrs()["href"])
		date := time.Unix(timestamp, 0)
		board.Append3x1(title, "", url, date)
	}
	return board, nil
}
