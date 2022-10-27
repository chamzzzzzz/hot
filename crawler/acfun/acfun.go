package acfun

import (
	"bytes"
	"encoding/json"
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
	"strings"
)

type Crawler struct {
}

func (c *Crawler) Name() string {
	return "acfun"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.acfun.cn/?pagelets=pagelet_header&ajaxpipe=1", nil)
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

	body := &body{}
	if err := json.Unmarshal(bytes.TrimSuffix(data, []byte(`/*<!-- fetch-stream -->*/`)), body); err != nil {
		return nil, err
	}

	dom := soup.HTMLParse(body.HTML)
	if dom.Error != nil {
		return nil, dom.Error
	}

	board := hot.NewBoard(c.Name())
	div := dom.Find("div", "class", "search-result")
	if div.Error != nil {
		return nil, div.Error
	}
	for _, li := range div.FindAllStrict("li") {
		a := li.Find("a")
		if a.Error != nil {
			return nil, a.Error
		}
		b := a.Find("b")
		if b.Error != nil {
			return nil, b.Error
		}
		title := strings.TrimSpace(b.Text())
		url := "https://www.acfun.cn" + strings.TrimSpace(a.Attrs()["href"])
		board.AppendTitleURL(title, url)
	}
	return board, nil
}

type body struct {
	HTML string `json:"html"`
}
