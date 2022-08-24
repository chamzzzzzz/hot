package sohu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
	"time"
)

type Crawler struct {
}

func (c *Crawler) Name() string {
	return "sohu"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://cis.sohu.com/cisv4/feeds", bytes.NewReader([]byte(`{"suv":"1620464912841ka7wm8","pvId":"1621619924793AvVUE3F","clientType":1,"resourceParam":[{"requestId":"1661230451471_1620464912841k_ihx","resourceId":"1661230451472509684","page":1,"size":10,"spm":"smpc.csrpage.0.0.166123044062415lMEH3","context":{"feedType":"XTOPIC_SYNTHETICAL"},"resProductParam":{"productId":268,"productType":14},"productParam":{"productId":268,"productType":14,"categoryId":"19","mediaId":1}}]}`)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var bodyJsonWrap bodyJsonWrap
	if err := json.Unmarshal(body, &bodyJsonWrap); err != nil {
		return nil, err
	}

	bodyJson := &bodyJson{}
	if err := json.Unmarshal(bodyJsonWrap.Body(), bodyJson); err != nil {
		return nil, err
	} else if bodyJson.Status != 1 {
		return nil, fmt.Errorf("body status: %d", bodyJson.Status)
	}

	board := hot.NewBoard(c.Name())
	date := time.Now()
	for _, data := range bodyJson.Data {
		board.Append(data.ResourceData.ContentData.Title, "", date)
	}
	return board, nil
}

type bodyJsonWrap map[string]json.RawMessage

func (w bodyJsonWrap) Body() json.RawMessage {
	for _, v := range w {
		return v
	}
	return nil
}

type bodyJson struct {
	Data []struct {
		ResourceData struct {
			ContentData struct {
				Title string `json:"title"`
			} `json:"contentData"`
		} `json:"resourceData"`
	} `json:"data"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}
