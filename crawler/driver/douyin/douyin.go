package douyin

import (
	"fmt"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "douyin"
	ProxySwitch = false
	URL         = "https://aweme.snssdk.com/aweme/v1/hot/search/list/"
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
	body := &body{}
	if err := httputil.Request("GET", URL, nil, "json", body, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	}
	if body.StatusCode != 0 {
		return nil, fmt.Errorf("body status_code: %d", body.StatusCode)
	}

	board := hot.NewBoard(c.Name())
	for _, word := range body.Data.WordList {
		board.Append(&hot.Hot{
			Title: word.Word,
		})
	}
	return board, nil
}

type body struct {
	StatusCode int `json:"status_code"`
	Data       struct {
		WordList []struct {
			Word string `json:"word"`
		} `json:"word_list"`
	} `json:"data"`
}
