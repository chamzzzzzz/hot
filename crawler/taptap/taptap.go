package taptap

import (
	"encoding/json"
	"fmt"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
)

type Crawler struct {
}

func (c *Crawler) Name() string {
	return "taptap"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", `https://www.taptap.cn/webapiv2/sidebar/v2/list?type=rankings&X-UA=V%3D1%26PN%3DWebApp%26LANG%3Dzh_CN%26VN_CODE%3D93%26VN%3D0.1.0%26LOC%3DCN%26PLT%3DPC%26DS%3DAndroid%26UID%3D034e87e7-3cf1-4631-b104-f3abadbc8b40%26DT%3DPC%26OS%3DMac%2520OS%26OSV%3D10.15`, nil)
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
	} else if !body.Success {
		return nil, fmt.Errorf("body success: %v", body.Success)
	}

	board := hot.NewBoard(c.Name())
	for _, data := range body.Data {
		for _, data := range data.Data.Data {
			for _, data := range data.List {
				board.Append1(data.Keyword)
			}
		}
	}
	return board, nil
}

type body struct {
	Data []struct {
		Label string `json:"label"`
		Data  struct {
			Type string `json:"type"`
			Data []struct {
				Title string `json:"title"`
				List  []struct {
					Keyword string `json:"keyword"`
				} `json:"list"`
			} `json:"data"`
		} `json:"data,omitempty"`
	} `json:"data"`
	Now     int  `json:"now"`
	Success bool `json:"success"`
}
