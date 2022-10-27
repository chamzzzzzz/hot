package leikeji

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
	return "leikeji"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.leikeji.com", nil)
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
	ul := dom.FindStrict("ul", "class", "ui-sideArticleList")
	if ul.Error != nil {
		return nil, ul.Error
	}
	for _, a := range ul.FindAllStrict("a", "class", "link") {
		div := a.Find("div", "class", "time")
		if div.Error != nil {
			return nil, div.Error
		}
		title := strings.TrimSpace(a.Attrs()["title"])
		url := "https://www.leikeji.com" + strings.TrimSpace(a.Attrs()["href"])
		date, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(div.Text()), time.Local)
		if err != nil {
			if strings.Contains(strings.TrimSpace(div.Text()), "小时") {
				date = time.Now()
			} else {
				return nil, err
			}
		}
		board.Append3x1(title, "", url, date)
	}
	return board, nil
}
