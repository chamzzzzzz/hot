package kr36

import (
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	RenQi    = "renqi"
	ZongHe   = "zonghe"
	ShouCang = "shoucang"
	Unknown  = "unknown"
)

type Crawler struct {
	Catalog string
}

func (c *Crawler) Name() string {
	return "kr36"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://36kr.com/hot-list/catalog", nil)
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
	for i, div := range dom.FindAllStrict("div", "class", "list-section-wrapper") {
		catalog := itocatalog(i)
		if c.Catalog != "" && c.Catalog != catalog {
			continue
		}

		for _, div := range div.FindAll("div", "class", "article-wrapper") {
			a1 := div.FindStrict("a", "class", "article-item-title weight-bold")
			if a1.Error != nil {
				return nil, a1.Error
			}
			a2 := div.FindStrict("a", "class", "article-item-description ellipsis-2")
			if a2.Error != nil {
				return nil, a2.Error
			}
			title := strings.TrimSpace(a1.Text())
			summary := strings.TrimSpace(a2.Text())
			url := "https://36kr.com" + strings.TrimSpace(a1.Attrs()["href"])
			board.Append4(title, summary, url, catalog)
		}
	}
	return board, nil
}

func itocatalog(i int) string {
	switch i {
	case 0:
		return RenQi
	case 1:
		return ZongHe
	case 2:
		return ShouCang
	default:
		return Unknown
	}
}
