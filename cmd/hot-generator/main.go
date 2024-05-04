package main

import (
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/urfave/cli/v2"
)

var (
	logger = log.New(os.Stdout, "generator: ", log.Ldate|log.Lmicroseconds)
)

type Option struct {
	Name     string
	Proxy    bool
	Json     bool
	XML      bool
	Backtick string
}

func main() {
	o := &Option{
		Backtick: "`",
	}

	app := &cli.App{
		Usage: "hot crawler generator",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "name",
				Value:       "hg",
				Destination: &o.Name,
			},
			&cli.BoolFlag{
				Name:        "proxy",
				Value:       false,
				Usage:       "enable proxy",
				Destination: &o.Proxy,
			},
			&cli.BoolFlag{
				Name:        "json",
				Value:       false,
				Usage:       "json response",
				Destination: &o.Json,
			},
			&cli.BoolFlag{
				Name:        "xml",
				Value:       false,
				Usage:       "xml response",
				Destination: &o.XML,
			},
		},
		Action: func(c *cli.Context) error {
			err := os.Mkdir(filepath.Join("crawler", "driver", o.Name), 0750)
			if err != nil && !os.IsExist(err) {
				logger.Printf("generate package dir, error='%s', crawler=%s", err, o.Name)
				return err
			}
			f1, err := os.Create(filepath.Join("crawler", "driver", o.Name, o.Name) + ".go")
			if err != nil {
				logger.Printf("generate main code, error='%s', crawler=%s", err, o.Name)
				return err
			}
			defer f1.Close()
			f2, err := os.Create(filepath.Join("crawler", "driver", o.Name, o.Name) + "_test.go")
			if err != nil {
				logger.Printf("generate test code, error='%s', crawler=%s", err, o.Name)
				return err
			}
			defer f2.Close()
			if err := sources.ExecuteTemplate(f1, "main", o); err != nil {
				logger.Printf("generate main code, error='%s', crawler=%s", err, o.Name)
				return err
			}
			if err := sources.ExecuteTemplate(f2, "test", o); err != nil {
				logger.Printf("generate test code, error='%s', crawler=%s", err, o.Name)
				return err
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Printf("generate, error='%s', crawler=%s", err, o.Name)
	}
}

var sources = template.Must(template.New("sources").Parse(`
{{define "main" -}}
package {{.Name}}

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
)

const (
	DriverName  = "{{.Name}}"
	ProxySwitch = false
	URL         = "https://www.{{.Name}}.com"
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
{{- if .Json}}
	body := &body{}
	if err := httputil.Request("GET", URL, nil, "json", body, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	for _, data := range body.Data {
		board.Append(&hot.Hot{Title: strings.TrimSpace(data.Title), URL: strings.TrimSpace(data.URL)})
	}
{{- else if .XML}}
	body := &body{}
	if err := httputil.Request("GET", URL, nil, "xml", body, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	for _, data := range body.Channel.Item {
		board.Append(&hot.Hot{Title: strings.TrimSpace(data.Title), URL: strings.TrimSpace(data.Link)})
	}
{{- else}}
	dom := &httputil.DOM{}
	if err := httputil.Request("GET", URL, nil, "dom", dom, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	for _, a := range dom.QueryAll("a") {
		title := strings.TrimSpace(a.Text())
		url := strings.TrimSpace(a.Href())
		board.Append(&hot.Hot{Title: title, URL: url})
	}
{{- end}}
	return board, nil
}
{{- if .Json}}

type body struct {
	Code    int    {{.Backtick}}json:"code"{{.Backtick}}
	Message string {{.Backtick}}json:"message"{{.Backtick}}
	Data    []struct {
		Title string {{.Backtick}}json:"title"{{.Backtick}}
		URL   string {{.Backtick}}json:"url"{{.Backtick}}
	} {{.Backtick}}json:"data"{{.Backtick}}
}

func (body *body) NormalizedCode() int {
	return body.Code
}
{{- else if .XML}}

type body struct {
	Channel struct {
		Item []struct {
			Title   string {{.Backtick}}xml:"title"{{.Backtick}}
			Link    string {{.Backtick}}xml:"link"{{.Backtick}}
			PubDate string {{.Backtick}}xml:"pubDate"{{.Backtick}}
		} {{.Backtick}}xml:"item"{{.Backtick}}
	} {{.Backtick}}xml:"channel"{{.Backtick}}
}
{{- end}}
{{end}}

{{define "test" -}}
package {{.Name}}

import (
	"github.com/chamzzzzzz/hot/crawler/driver"
	"testing"
)

func TestCrawl(t *testing.T) {
	c := Crawler{Option: driver.NewTestOptionFromEnv()}
	if board, err := c.Crawl(); err != nil {
		t.Error(err)
	} else {
		for _, hot := range board.Hots {
			t.Log(hot)
		}
	}
}
{{end}}`))
