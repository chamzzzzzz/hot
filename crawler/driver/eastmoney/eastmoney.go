package eastmoney

import (
	"fmt"
	"regexp"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

var re = regexp.MustCompile(`//searchadapter\.eastmoney\.com/api/hotkeyword/get\?count=20&token=([A-Z0-9]+)`)

const (
	DriverName  = "eastmoney"
	ProxySwitch = false
	URL         = "https://so.eastmoney.com/newstatic/js/page/welcome.js"
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
	URL, err := c.getURL()
	if err != nil {
		return nil, err
	}

	body := &body{}
	if err := httputil.Request("GET", URL, nil, "json", body, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	} else if body.Status != 0 {
		return nil, fmt.Errorf("body status: %d", body.Status)
	}

	board := hot.NewBoard(c.Name())
	for _, data := range body.Data {
		board.Append(&hot.Hot{
			Title: data.KeyPhrase,
		})
	}
	return board, nil
}

func (c *Crawler) getURL() (string, error) {
	data, err := httputil.RequestData("GET", URL, nil, httputil.NewOption(c.Option, ProxySwitch))
	if err != nil {
		return "", err
	}

	url := re.FindString(string(data))
	if url == "" {
		return "", fmt.Errorf("not match")
	}
	return "https:" + url, nil
}

type body struct {
	Data []struct {
		KeyPhrase string `json:"KeyPhrase"`
	} `json:"Data"`
	Status  int    `json:"Status"`
	Message string `json:"Message"`
}
