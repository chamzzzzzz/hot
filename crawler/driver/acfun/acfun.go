package acfun

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
	DriverName = "acfun"
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
	req, err := http.NewRequest("GET", "https://www.acfun.cn/rest/pc-direct/rank/channel?channelId=&subChannelId=&rankLimit=30&rankPeriod=DAY", nil)
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
		fmt.Println(string(data))
		return nil, err
	} else if body.Result != 0 {
		return nil, fmt.Errorf("body result: %d", body.Result)
	}

	board := hot.NewBoard(c.Name())
	for _, data := range body.RankList {
		title := strings.TrimSpace(data.ContentTitle)
		summary := strings.TrimSpace(data.ContentDesc)
		url := strings.TrimSpace(data.ShareURL)
		date := time.UnixMilli(data.CreateTimeMillis)
		board.Append3x1(title, summary, url, date)
	}
	return board, nil
}

type body struct {
	Result   int `json:"result"`
	RankList []struct {
		ContentID        int    `json:"contentId"`
		ContributeTime   int64  `json:"contributeTime"`
		ContentTitle     string `json:"contentTitle"`
		ContentDesc      string `json:"contentDesc"`
		CreateTimeMillis int64  `json:"createTimeMillis"`
		ShareURL         string `json:"shareUrl"`
		Title            string `json:"title"`
	} `json:"rankList"`
}
