package zol

import (
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
	"strings"
)

type Crawler struct {
}

func (c *Crawler) Name() string {
	return "zol"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://news.zol.com.cn/", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(transform.NewReader(res.Body, simplifiedchinese.GBK.NewDecoder()))
	if err != nil {
		return nil, err
	}

	dom := soup.HTMLParse(string(data))
	if dom.Error != nil {
		return nil, dom.Error
	}

	board := hot.NewBoard(c.Name())
	for _, id := range []string{"list-v-1", "list-v-2"} {
		div := dom.Find("div", "id", id)
		if div.Error != nil {
			return nil, div.Error
		}
		for _, a := range div.FindAllStrict("a") {
			title := strings.TrimSpace(a.Text())
			url := "https:" + strings.TrimSpace(a.Attrs()["href"])
			board.AppendTitleURL(title, url)
		}
	}
	return board, nil
}
