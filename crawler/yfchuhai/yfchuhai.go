package yfchuhai

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
	return "yfchuhai"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.yfchuhai.com", nil)
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
	for _, a := range dom.FindAllStrict("a", "class", "list__item flex") {
		h2 := a.Find("h2")
		if h2.Error != nil {
			return nil, h2.Error
		}
		span := a.Find("span", "class", "text")
		if span.Error != nil {
			return nil, span.Error
		}
		title := strings.TrimSpace(h2.Text())
		url := "https://www.yfchuhai.com" + strings.TrimSpace(a.Attrs()["href"])
		date, err := time.ParseInLocation("01-02 15:04", strings.TrimSpace(span.Text()), time.Local)
		if err != nil {
			return nil, err
		}
		date = date.AddDate(time.Now().Year(), 0, 0)
		board.AppendTitleURLDate(title, url, date)
	}
	return board, nil
}
