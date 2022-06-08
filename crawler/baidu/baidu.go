package baidu

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
	return "baidu"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://top.baidu.com/board?tab=realtime", nil)
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
	for _, div := range dom.FindAllStrict("div", "class", "content_1YWBm") {
		div01 := div.FindStrict("div", "class", "c-single-text-ellipsis")
		div02 := div.Find("div", "class", "small_Uvkd3")
		if div01.Error != nil || div02.Error != nil {
			continue
		}

		title := strings.ReplaceAll(strings.Trim(div01.Text(), " "), "\n", "")
		summary := strings.ReplaceAll(strings.Trim(div02.Text(), " "), "\n", "")
		board.Append(title, summary, date)
	}
	return board, nil
}
