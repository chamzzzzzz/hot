package vrtuoluo

import (
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
	"strings"
)

type Crawler struct {
}

func (c *Crawler) Name() string {
	return "vrtuoluo"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.vrtuoluo.cn", nil)
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
	for _, div := range dom.FindAllStrict("div", "class", "hot_article_view") {
		for _, a := range div.FindAllStrict("a", "class", "title_box text_ellipsis") {
			title := strings.TrimSpace(a.Text())
			url := "https://www.vrtuoluo.cn" + strings.TrimSpace(a.Attrs()["href"])
			board.AppendTitleURL(title, url)
		}
		break
	}
	return board, nil
}
