package ifeng

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
	return "ifeng"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.ifeng.com/", nil)
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
	div := dom.Find("div", "class", "hot_box-1yXFLW7e")
	if div.Error != nil {
		return nil, div.Error
	}
	h3 := div.FindStrict("h3", "class", "list_title-nOvTJ00k big-3QPOjsEI")
	if h3.Error != nil {
		return nil, div.Error
	}
	a := h3.Find("a")
	if a.Error != nil {
		return nil, div.Error
	}
	title := strings.TrimSpace(a.Text())
	summary := a.Attrs()["href"]
	board.Append(title, summary, date)
	for _, p := range div.FindAllStrict("p", "class", "news_list_p-3EcL2Tvk ") {
		a = p.Find("a")
		if a.Error != nil {
			return nil, div.Error
		}
		title = strings.TrimSpace(a.Text())
		summary = a.Attrs()["href"]
		board.Append(title, summary, date)
	}
	return board, nil
}
