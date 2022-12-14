package driver

import (
	"fmt"
	"github.com/chamzzzzzz/hot"
	"net/url"
	"os"
	"sort"
	"sync"
)

type Option struct {
	ID         string
	DriverName string
	Proxy      string
	Cookie     string
	Catalog    string
	Raw        string
}

func NewTestOptionFromEnv() (option Option) {
	option.Proxy = os.Getenv("HOT_CRAWLER_DRIVER_TEST_OPTION_PROXY")
	option.Cookie = os.Getenv("HOT_CRAWLER_DRIVER_TEST_OPTION_COOKIE")
	option.Cookie = os.Getenv("HOT_CRAWLER_DRIVER_TEST_OPTION_CATALOG")
	return
}

func ParseOption(raw string) (*Option, error) {
	v, err := url.ParseQuery(raw)
	if err != nil {
		return nil, err
	}
	o := &Option{
		ID:         v.Get("id"),
		DriverName: v.Get("drivername"),
		Proxy:      v.Get("proxy"),
		Cookie:     v.Get("cookie"),
		Catalog:    v.Get("catalog"),
		Raw:        raw,
	}
	if o.ID == "" {
		o.ID = o.DriverName
	}
	return o, nil
}

func (o *Option) Encode() string {
	if o == nil {
		return ""
	}
	v := url.Values{}
	if o.ID != "" && o.ID != o.DriverName {
		v.Set("id", o.ID)
	}
	if o.DriverName != "" {
		v.Set("drivername", o.DriverName)
	}
	if o.Proxy != "" {
		v.Set("proxy", o.Proxy)
	}
	if o.Cookie != "" {
		v.Set("cookie", o.Cookie)
	}
	if o.Catalog != "" {
		v.Set("catalog", o.Catalog)
	}
	return v.Encode()
}

type Crawler interface {
	Driver() Driver
	Name() string
	Crawl() (*hot.Board, error)
}

type Driver interface {
	Open(option Option) (Crawler, error)
}

var (
	drivers = make(map[string]Driver)
	mu      sync.RWMutex
)

func Register(name string, driver Driver) {
	mu.Lock()
	defer mu.Unlock()
	if driver == nil {
		panic("hot/crawler/driver: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("hot/crawler/driver: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

func Drivers() []string {
	mu.RLock()
	defer mu.RUnlock()
	list := make([]string, 0, len(drivers))
	for name := range drivers {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

func Open(option Option) (Crawler, error) {
	mu.RLock()
	driver, ok := drivers[option.DriverName]
	mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unknown driver %q (forgotten import?)", option.DriverName)
	}
	return driver.Open(option)
}
