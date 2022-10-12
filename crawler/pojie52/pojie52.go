package pojie52

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Crawler struct {
}

func (c *Crawler) Name() string {
	return "pojie52"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.52pojie.cn/forum.php?mod=guide&view=hot", nil)
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
	div := dom.FindStrict("div", "id", "threadlist")
	if div.Error != nil {
		return nil, div.Error
	}
	div = div.Find("div", "class", "bm_c")
	if div.Error != nil {
		return nil, div.Error
	}
	for _, tbody := range div.FindAll("tbody") {
		a := tbody.Find("a", "class", "xst")
		if a.Error != nil {
			return nil, a.Error
		}
		td := tbody.FindAll("td", "class", "by")
		if len(td) != 3 {
			return nil, fmt.Errorf("td count invalid")
		}
		span := td[1].Find("span")
		if span.Error != nil {
			return nil, span.Error
		}

		title := strings.TrimSpace(a.Text())
		url := "https://www.52pojie.cn/" + strings.TrimSpace(a.Attrs()["href"])
		date, err := time.ParseInLocation("2006-1-2 15:04", strings.TrimSpace(span.Text()), time.Local)
		if err != nil {
			return nil, err
		}
		board.Append3x1(title, "", url, date)
	}
	return board, nil
}
