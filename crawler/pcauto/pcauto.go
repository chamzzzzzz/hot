package pcauto

import (
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	Article = "article"
	Topic   = "topic"
	Forum   = "forum"
	Unknown = "unknown"
)

type Crawler struct {
	Catalog string
}

func (c *Crawler) Name() string {
	return "pcauto"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	switch c.Catalog {
	case Article:
		return c.article()
	case Topic:
		return c.topic()
	case Forum:
		return c.forum()
	default:
		return c.all()
	}
}

func (c *Crawler) article() (*hot.Board, error) {
	return c.withClasses("txts")
}

func (c *Crawler) topic() (*hot.Board, error) {
	return c.withClasses("bbs_topics")
}

func (c *Crawler) forum() (*hot.Board, error) {
	return c.withClasses("hotForums")
}

func (c *Crawler) all() (*hot.Board, error) {
	return c.withClasses("txts", "bbs_topics", "hotForums")
}

func (c *Crawler) withClasses(classes ...string) (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.pcauto.com.cn/", nil)
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
	div := dom.FindStrict("div", "class", "section ranking")
	if div.Error != nil {
		return nil, div.Error
	}
	for _, class := range classes {
		ul := div.FindStrict("ul", "class", class)
		if ul.Error != nil {
			return nil, ul.Error
		}
		catalog := classtocatalog(class)
		for _, a := range ul.FindAllStrict("a") {
			title := strings.TrimSpace(a.Text())
			url := "https:" + strings.TrimSpace(a.Attrs()["href"])
			board.AppendTitleURLCatalog(title, url, catalog)
		}
	}
	return board, nil
}

func classtocatalog(class string) string {
	switch class {
	case "txts":
		return Article
	case "bbs_topics":
		return Topic
	case "hotForums":
		return Forum
	default:
		return Unknown
	}
}
