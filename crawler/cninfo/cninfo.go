package cninfo

import (
	"encoding/json"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Crawler struct {
}

func (c *Crawler) Name() string {
	return "cninfo"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://www.cninfo.com.cn/new/index/getHotAnnouces", nil)
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
	for _, data := range body {
		title := strings.TrimSpace(data.ShortName)
		summary := strings.TrimSpace(data.Name)
		url := "http://static.cninfo.com.cn/" + strings.TrimSpace(data.URL)
		date := time.UnixMilli(data.AnnouncementTime)
		board.Append3x1(title, summary, url, date)
	}
	return board, nil
}

type body struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	URL              string `json:"url"`
	ShortName        string `json:"shortName"`
	AnnouncementTime int64  `json:"announcementTime"`
}
