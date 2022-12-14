package sohu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"io/ioutil"
	"net/http"
)

const (
	DriverName = "sohu"
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

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var bodyWrap bodyWrap
	if err := json.Unmarshal(data, &bodyWrap); err != nil {
		return nil, err
	}

	body := &body{}
	if err := json.Unmarshal(bodyWrap.Body(), body); err != nil {
		return nil, err
	} else if body.Status != 1 {
		return nil, fmt.Errorf("body status: %d", body.Status)
	}

	board := hot.NewBoard(c.Name())
	for _, data := range body.Data {
		board.Append1(data.ResourceData.ContentData.Title)
	}
	return board, nil
}

type bodyWrap map[string]json.RawMessage

func (w bodyWrap) Body() json.RawMessage {
	for _, v := range w {
		return v
	}
	return nil
}

type body struct {
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
