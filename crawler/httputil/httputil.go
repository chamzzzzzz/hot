package httputil

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"net/http"
	"net/url"
)

type Option struct {
	Proxy           string
	Header          http.Header
	ContentEncoding string
}

func NewOption(option driver.Option, proxyswitch bool) *Option {
	o := &Option{
		Header: make(http.Header),
	}
	if option.ProxySwitch != "" || proxyswitch {
		o.Proxy = option.Proxy
	}
	return o
}

func NewClient(option *Option) *http.Client {
	client := &http.Client{}
	if option.Proxy != "" {
		if proxyUrl, err := url.Parse(option.Proxy); err == nil {
			client.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			}
		}
	}
	return client
}

func NewRequest(method, url string, body io.Reader, option *Option) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for key, _ := range option.Header {
		req.Header.Del(key)
		for _, value := range option.Header.Values(key) {
			req.Header.Add(key, value)
		}
	}
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")
	}
	return req, nil
}

func Request(method, url string, body io.Reader, unmarshalmethod string, unmarshalbody any, option *Option) error {
	client := NewClient(option)
	req, err := NewRequest(method, url, body, option)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var reader io.Reader = res.Body
	if option.ContentEncoding != "" {
		if option.ContentEncoding == "gbk" {
			reader = transform.NewReader(res.Body, simplifiedchinese.GBK.NewDecoder())
		}
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	switch unmarshalmethod {
	case "json":
		if err := json.Unmarshal(data, unmarshalbody); err != nil {
			return err
		}
	case "xml":
		if err := xml.Unmarshal(data, unmarshalbody); err != nil {
			return err
		}
	case "dom":
		root := soup.HTMLParse(string(data))
		if root.Error != nil {
			return root.Error
		}
		if dom, ok := unmarshalbody.(*DOM); ok {
			dom.Root = root
		} else {
			return fmt.Errorf("no support unmarshal body:%T", unmarshalbody)
		}
	default:
		return fmt.Errorf("no support unmarshal method:%s", unmarshalmethod)
	}
	return nil
}

func RequestData(method, url string, body io.Reader, option *Option) ([]byte, error) {
	client := NewClient(option)
	req, err := NewRequest(method, url, body, option)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func RequestCookie(method, url string, body io.Reader, option *Option) (string, error) {
	client := NewClient(option)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	req, err := NewRequest(method, url, body, option)
	if err != nil {
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	for _, cookie := range res.Cookies() {
		req.AddCookie(cookie)
	}
	return req.Header.Get("Cookie"), nil
}

func UpdateCookie(method, url string, body io.Reader, option *Option, cookie *string) error {
	if cookie == nil {
		return nil
	}
	if *cookie != "" {
		return nil
	}
	if _cookie, err := RequestCookie(method, url, nil, option); err != nil {
		return err
	} else {
		*cookie = _cookie
		return nil
	}
}

type DOM struct {
	soup.Root
}