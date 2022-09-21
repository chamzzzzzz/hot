package daniu

import (
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
	"time"
)

type Crawler struct {
}

func (c *Crawler) Name() string {
	return "daniu"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.daniu522.com/misc.php?mod=ranklist&type=thread&view=heats&orderby=today", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	html, err := ioutil.ReadAll(transform.NewReader(res.Body, simplifiedchinese.GBK.NewDecoder()))
	if err != nil {
		return nil, err
	}

	dom := soup.HTMLParse(string(html))
	if dom.Error != nil {
		return nil, dom.Error
	}

	board := hot.NewBoard(c.Name())
	date := time.Now()
	div := dom.FindStrict("div", "class", "tl")
	if div.Error != nil {
		return nil, div.Error
	}
	for _, th := range div.FindAllStrict("th") {
		a := th.FindStrict("a")
		if a.Error != nil {
			continue
		}
		title := a.Text()
		summary := a.Attrs()["href"]
		board.Append(title, summary, date)
	}
	return board, nil
}
