package tianya

import (
	"encoding/json"
	"fmt"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	DriverName = "tianya"
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
	req, err := http.NewRequest("GET", "https://bbs.tianya.cn/api?method=bbs.ice.getHotArticleList&params.pageSize=40&params.pageNum=1", nil)
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
	} else if body.Code != "1" {
		return nil, fmt.Errorf("body code: %s", body.Code)
	}

	board := hot.NewBoard(c.Name())
	for _, row := range body.Data.Rows {
		title := strings.TrimSpace(row.Title)
		url := strings.TrimSpace(row.URL)
		board.AppendTitleURL(title, url)
	}
	return board, nil
}

type body struct {
	Success string `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Rows []struct {
			Title string `json:"title"`
			URL   string `json:"url"`
		} `json:"rows"`
	} `json:"data"`
}
