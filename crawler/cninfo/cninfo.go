package cninfo

import (
	"encoding/json"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
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

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var bodyJson []bodyJson
	if err := json.Unmarshal(body, &bodyJson); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	date := time.Now()
	for _, data := range bodyJson {
		board.Append(data.Name, data.URL, date)
	}
	return board, nil
}

type bodyJson struct {
	ID                 string      `json:"id"`
	Name               string      `json:"name"`
	CollectionTimes    int         `json:"collectionTimes"`
	ForwardTimes       int         `json:"forwardTimes"`
	ClickTotalTimes    int         `json:"clickTotalTimes"`
	RecentlyClickTimes int         `json:"recentlyClickTimes"`
	LogoURL            interface{} `json:"logoUrl"`
	URL                string      `json:"url"`
	Code               string      `json:"code"`
	ShortName          string      `json:"shortName"`
	PlateCode          string      `json:"plateCode"`
	AdjunctSize        int         `json:"adjunctSize"`
	AdjunctType        string      `json:"adjunctType"`
	AnnouncementTime   int64       `json:"announcementTime"`
}
