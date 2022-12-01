package gk99

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
	return "gk99"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://dota2.gk99.com/rd/", nil)
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
	ul := dom.FindStrict("ul", "class", "cfix fin_newsList")
	if ul.Error != nil {
		return nil, ul.Error
	}
	for _, li := range ul.FindAllStrict("li") {
		a := li.Find("a")
		if a.Error != nil {
			return nil, a.Error
		}
		p := li.Find("p")
		if p.Error != nil {
			return nil, p.Error
		}
		em := li.Find("em", "class", "fRight")
		if em.Error != nil {
			return nil, em.Error
		}
		title := strings.TrimSpace(a.Text())
		summary := strings.TrimSpace(p.Text())
		url := strings.TrimSpace(a.Attrs()["href"])
		date, err := time.ParseInLocation("2006-01-02 15:04:05", strings.TrimSpace(em.Text()), time.Local)
		if err != nil {
			return nil, err
		}
		board.Append3x1(title, summary, url, date)
	}
	return board, nil
}
