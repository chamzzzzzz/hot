package v2ex

import (
	"encoding/json"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	DriverName = "v2ex"
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
	if c.Option.Proxy != "" {
		if proxyUrl, err := url.Parse(c.Option.Proxy); err == nil {
			client.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			}
		}
	}

	req, err := http.NewRequest("GET", "https://www.v2ex.com/api/topics/hot.json", nil)
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

	var body []body
	if err := json.Unmarshal(data, &body); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	for _, item := range body {
		date := time.Unix(item.Created, 0)
		board.Append3x1(item.Title, "", item.URL, date)
	}
	return board, nil
}

type body struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Created int64  `json:"created"`
}
