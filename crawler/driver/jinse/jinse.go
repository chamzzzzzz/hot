package jinse

import (
	"fmt"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "jinse"
	ProxySwitch = false
	URL         = "https://newapi.jinse.com/noah/v1/"
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
	return c.jinse()
}

func (c *Crawler) jinse_x_search() (*hot.Board, error) {
	return c.jinseWithURL(URL + "hot-search")
}

func (c *Crawler) jinse_x_article() (*hot.Board, error) {
	return c.jinseWithURL(URL + "articles/hot?hour_diff=24")
}

func (c *Crawler) jinse() (*hot.Board, error) {
	board, err := c.jinse_x_search()
	if err != nil {
		return nil, err
	}

	board2, err := c.jinse_x_article()
	if err != nil {
		return nil, err
	}

	for _, hot := range board2.Hots {
		board.Append0(hot)
	}
	return board, nil
}

func (c *Crawler) jinseWithURL(url string) (*hot.Board, error) {
	body := &body{}
	if err := httputil.Request("GET", url, nil, "json", body, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	} else if body.StatusCode != 0 {
		return nil, fmt.Errorf("body status_code: %d", body.StatusCode)
	}

	board := hot.NewBoard(c.Name())
	for _, data := range body.Data {
		board.AppendTitleURL(data.Title, data.JumpURL)
	}
	return board, nil
}

type body struct {
	StatusCode int `json:"status_code"`
	Data       []struct {
		Title   string `json:"title"`
		JumpURL string `json:"jump_url"`
	} `json:"data"`
}
