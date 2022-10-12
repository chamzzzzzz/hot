package cyzone

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
	return "cyzone"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.cyzone.cn/hot/", nil)
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
	for _, div := range dom.FindAllStrict("div", "class", "item-intro") {
		a := div.Find("a", "class", "item-title")
		if a.Error != nil {
			return nil, a.Error
		}
		p := div.Find("p", "class", "item-desc")
		if p.Error != nil {
			return nil, p.Error
		}
		span := div.Find("span", "class", "time")
		if span.Error != nil {
			return nil, span.Error
		}

		timestamp, err := strconv.ParseInt(span.Attrs()["data-time"], 10, 64)
		if err != nil {
			return nil, err
		}

		title := strings.TrimSpace(a.Text())
		summary := strings.TrimSpace(p.Text())
		url := "https:" + strings.TrimSpace(a.Attrs()["href"])
		date := time.Unix(timestamp, 0)
		board.Append3x1(title, summary, url, date)
	}
	return board, nil
}
