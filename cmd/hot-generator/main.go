package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

var (
	logger = log.New(os.Stdout, "generator: ", log.Ldate|log.Lmicroseconds)
)

type Option struct {
	Name      string
	UserAgent string
	Proxy     bool
	Json      bool
	Backtick  string
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
			&cli.StringFlag{
				Name:        "useragent",
				Value:       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36",
				Usage:       "useragent",
				Destination: &o.UserAgent,
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
		},
		Action: func(c *cli.Context) error {
			err := os.Mkdir(filepath.Join("crawler", o.Name), 0750)
			if err != nil && !os.IsExist(err) {
				logger.Printf("generate package dir, error='%s', crawler=%s", err, o.Name)
				return err
			}
			f1, err := os.Create(filepath.Join("crawler", o.Name, o.Name) + ".go")
			if err != nil {
				logger.Printf("generate main code, error='%s', crawler=%s", err, o.Name)
				return err
			}
			defer f1.Close()
			f2, err := os.Create(filepath.Join("crawler", o.Name, o.Name) + "_test.go")
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
	{{- if .Json}}
	"encoding/json"
	{{- end}}
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
	{{- if .Proxy}}
	"net/url"
	{{- end}}
	"strings"
)

type Crawler struct {
}

func (c *Crawler) Name() string {
	return "{{.Name}}"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
{{- if .Proxy}}
	if c.Proxy != "" {
		if proxyUrl, err := url.Parse(c.Proxy); err == nil {
			client.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			}
		}
	}{{"\n"}}
{{- end}}
	req, err := http.NewRequest("GET", "https://www.{{.Name}}.com", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "{{.UserAgent}}")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

{{- if .Json}}

	body := &body{}
	if err := json.Unmarshal(data, body); err != nil {
		return nil, err
	} else if body.Code != "0" {
		return nil, fmt.Errorf("body code: %s", body.Code)
	}

	board := hot.NewBoard(c.Name())
	for _, data := range body.Data {
		board.AppendTitleURL(data.Title, data.URL)
	}
{{- else}}

	dom := soup.HTMLParse(string(data))
	if dom.Error != nil {
		return nil, dom.Error
	}

	board := hot.NewBoard(c.Name())
	div := dom.Find("div", "class", "hot")
	if div.Error != nil {
		return nil, div.Error
	}
	for _, a := range div.FindAllStrict("a") {
		title := strings.TrimSpace(a.Text())
		url := strings.TrimSpace(a.Attrs()["href"])
		board.AppendTitleURL(title, url)
	}
{{- end}}
	return board, nil
}
{{- if .Json}}

type body struct {
	Code    string {{.Backtick}}json:"code"{{.Backtick}}
	Message string {{.Backtick}}json:"message"{{.Backtick}}
	Data    []struct {
		Title string {{.Backtick}}json:"title"{{.Backtick}}
		URL   string {{.Backtick}}json:"url"{{.Backtick}}
	} {{.Backtick}}json:"data"{{.Backtick}}
}
{{- end}}
{{end}}

{{define "test" -}}
package {{.Name}}

import (
	{{- if .Proxy}}
	"os"
	{{- end}}
	"testing"
)

func TestCrawl(t *testing.T) {
	{{- if .Proxy}}
	c := Crawler{
		Proxy: os.Getenv("HOT_CRAWLER_TEST_PROXY"),
	}
	{{- else}}
	c := Crawler{}
	{{- end}}
	if board, err := c.Crawl(); err != nil {
		t.Error(err)
	} else {
		for _, hot := range board.Hots {
			t.Log(hot)
		}
	}
}
{{end}}`))
