package main

import (
	"bytes"
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"
)

type Board struct {
	Name        string
	DisplayName string
	Hots        []string
}

type Archive struct {
	Date   string
	Boards []*Board
}

type Service struct {
	PathPrefix string
	Manifest   string
	*log.Logger
	mux     http.ServeMux
	render  *Render
	archive *Archive
	mu      sync.RWMutex
}

func (s *Service) Init(ctx context.Context) error {
	if s.Logger == nil {
		s.Logger = log.Default()
	}
	archive, err := s.loadArchive(time.Now().Format("2006-01-02"), nil)
	if err != nil {
		return err
	}
	s.setArchive(archive)

	s.render = &Render{}
	s.mux.HandleFunc(s.PathPrefix+"/", s.HandleIndex)
	s.Printf("service init success. pathprefix=%s", s.PathPrefix)
	return nil
}

func (s *Service) Uninit(ctx context.Context) error {
	return nil
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Service) HandleIndex(w http.ResponseWriter, r *http.Request) {
	b, err := s.render.Index(s.getArchive())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func (s *Service) loadArchive(date string, names []string) (*Archive, error) {
	if names == nil {
		entries, err := os.ReadDir("archives")
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			if entry.IsDir() {
				names = append(names, entry.Name())
			}
		}
	}
	archive := &Archive{
		Date: date,
	}
	for _, name := range names {
		file := filepath.Join("archives", name, date+".txt")
		b, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		board := &Board{
			Name:        name,
			DisplayName: name,
		}
		for _, line := range strings.Split(string(b), "\n") {
			if line == "" {
				continue
			}
			board.Hots = append(board.Hots, line)
		}
		archive.Boards = append(archive.Boards, board)
	}
	return archive, nil
}

func (s *Service) setArchive(archive *Archive) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.archive = archive
}

func (s *Service) getArchive() *Archive {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.archive
}

type Render struct {
}

func (r *Render) Index(archvie *Archive) ([]byte, error) {
	source := `
<!DOCTYPE html>
<html>
	<head>
		<title>热门榜单</title>
	</head>
	<body>
		<header>
			<h1>热门榜单</h1>
			<p>{{.Date}}</p>
		</header>
		<main>
			{{- range $id, $board := .Boards}}
			<section>
				<h2>{{$board.DisplayName}}</h2>
				<ol>
					{{- range $id, $hot := $board.Hots}}
					<li>{{$hot}}</li>
					{{- end}}
				</ol>
			</section>
			{{- end}}
		</main>
		<footer>
			<p>Hot will eventually cool...</p>
		</footer>
	</body>
</html>
`
	tpl := template.Must(template.New("index").Parse(source))
	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, archvie)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func main() {
	addr := flag.String("addr", ":8080", "address")
	pathprefix := flag.String("pathprefix", "", "path prefix")
	flag.Parse()
	service := &Service{
		PathPrefix: *pathprefix,
	}
	if err := service.Init(context.TODO()); err != nil {
		log.Printf("init service fail. err='%s'", err)
		os.Exit(1)
	}
	if err := http.ListenAndServe(*addr, service); err != nil {
		log.Printf("listen and serve fail. err='%s'", err)
		os.Exit(1)
	}
}
