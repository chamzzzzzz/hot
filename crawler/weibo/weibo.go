package weibo

import (
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Crawler struct {
	Cookie string
}

func (c *Crawler) Name() string {
	return "weibo"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://s.weibo.com/top/summary", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")
	req.Header.Set("Cookie", c.Cookie)

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

	div := dom.FindStrict("div", "id", "pl_top_realtimehot")
	if div.Error != nil {
		return nil, div.Error
	}

	board := hot.NewBoard(c.Name())
	for _, tr := range div.FindAllStrict("tr", "class", "") {
		td01 := tr.Find("td", "class", "td-01")
		if td01.Error != nil {
			return nil, td01.Error
		}
		if _, err := strconv.Atoi(td01.Text()); err != nil {
			continue
		}

		td02 := tr.Find("td", "class", "td-02")
		if td02.Error != nil {
			return nil, td02.Error
		}
		a := td02.Find("a")
		if a.Error != nil {
			return nil, a.Error
		}
		board.Append1(a.Text())
	}
	return board, nil
}
