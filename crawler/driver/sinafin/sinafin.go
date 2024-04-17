package sinafin

import (
	"fmt"
	"strings"
	"time"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "sinafin"
	ProxySwitch = false
	URL         = "https://top.finance.sina.com.cn/ws/GetTopDataList.php?top_type=day&top_cat=finance_0_suda&top_time=%s&top_show_num=20&top_order=DESC&get_new=1"
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
	url := fmt.Sprintf(URL, time.Now().Format("20060102"))
	option := httputil.NewOption(c.Option, ProxySwitch)
	option.TrimPrefix = "var data = "
	option.TrimSuffix = ";"
	option.TrimSpace = true
	if err := httputil.Request("GET", url, nil, "json", body, option); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	for _, data := range body.Data {
		board.Append(&hot.Hot{Title: strings.TrimSpace(data.Title), URL: strings.TrimSpace(data.URL)})
	}
	return board, nil
}

type body struct {
	Data []struct {
		Title string `json:"title,omitempty"`
		URL   string `json:"url,omitempty"`
	} `json:"data,omitempty"`
}
