package zhiguf

import (
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	DriverName = "zhiguf"
)

type Driver struct {
}

func (driver *Driver) Open(option driver.Option) (driver.Crawler, error) {
	return &Crawler{Option: option}, nil
}

func init() {
	driver.Register(DriverName, &Driver{})
}

type Crawler struct {
	Option driver.Option
}

func (c *Crawler) Driver() driver.Driver {
	return &Driver{}
}

func (c *Crawler) Name() string {
	return DriverName
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.zhiguf.com", nil)
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
	for _, a := range dom.FindAllStrict("a", "class", "hot_articles_list") {
		div := a.FindStrict("div", "class", "hot_articles_title")
		if div.Error != nil {
			return nil, div.Error
		}
		p := div.Find("p")
		if p.Error != nil {
			return nil, p.Error
		}
		title := strings.TrimSpace(p.Text())
		url := strings.TrimSpace(a.Attrs()["href"])
		if strings.HasPrefix(url, "/focusnews_detail/") {
			url = "https://www.zhiguf.com" + url
		}
		board.AppendTitleURL(title, url)
	}
	return board, nil
}
