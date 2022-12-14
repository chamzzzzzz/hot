package thepaper

import (
	"encoding/json"
	"fmt"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	DriverName = "thepaper"
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
	req, err := http.NewRequest("GET", "https://www.thepaper.cn/contentapi/wwwIndex/rightSidebar", nil)
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

	body := &body{}
	if err := json.Unmarshal(data, body); err != nil {
		return nil, err
	} else if body.ResultCode != 1 {
		return nil, fmt.Errorf("body result code: %d", body.ResultCode)
	}

	board := hot.NewBoard(c.Name())
	for _, news := range body.Data.HotNews {
		title := strings.TrimSpace(news.Name)
		url := "https://www.thepaper.cn/newsDetail_forward_" + strings.TrimSpace(news.ContID)
		date := time.UnixMilli(news.PubTimeLong)
		board.Append3x1(title, "", url, date)
	}
	return board, nil
}

type body struct {
	ResultCode int    `json:"resultCode"`
	ResultMsg  string `json:"resultMsg"`
	Data       struct {
		HotNews []struct {
			ContID      string `json:"contId"`
			Name        string `json:"name"`
			PubTimeLong int64  `json:"pubTimeLong"`
		} `json:"hotNews"`
	} `json:"data"`
}
