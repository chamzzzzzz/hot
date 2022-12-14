package semiunion

import (
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	DriverName = "semiunion"
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
	req, err := http.NewRequest("GET", "http://www.semiunion.com/insight/", nil)
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
	for _, li := range dom.FindAllStrict("li", "class", "each-news") {
		div := li.Find("div", "class", "name")
		if div.Error != nil {
			return nil, div.Error
		}
		a := div.Find("a")
		if a.Error != nil {
			return nil, a.Error
		}
		div2 := li.Find("div", "class", "desc")
		if div2.Error != nil {
			return nil, div2.Error
		}
		span := li.Find("span", "class", "time")
		if span.Error != nil {
			return nil, span.Error
		}
		title := strings.TrimSpace(a.Text())
		summary := strings.TrimSpace(div2.Text())
		url := "http://www.semiunion.com" + strings.TrimSpace(a.Attrs()["href"])
		date, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(span.Text()), time.Local)
		if err != nil {
			return nil, err
		}
		board.Append3x1(title, summary, url, date)
	}
	return board, nil
}
