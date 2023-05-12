package taptap

import (
	"fmt"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "taptap"
	ProxySwitch = false
	URL         = `https://www.taptap.cn/webapiv2/sidebar/v2/list?type=rankings&X-UA=V%3D1%26PN%3DWebApp%26LANG%3Dzh_CN%26VN_CODE%3D93%26VN%3D0.1.0%26LOC%3DCN%26PLT%3DPC%26DS%3DAndroid%26UID%3D034e87e7-3cf1-4631-b104-f3abadbc8b40%26DT%3DPC%26OS%3DMac%2520OS%26OSV%3D10.15`
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
	if err := httputil.Request("GET", URL, nil, "json", body, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	}
	if !body.Success {
		return nil, fmt.Errorf("body success: %v", body.Success)
	}

	board := hot.NewBoard(c.Name())
	for _, data := range body.Data {
		for _, data := range data.Data.Data {
			for _, data := range data.List {
				board.Append(&hot.Hot{
					Title: data.Keyword,
				})
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
