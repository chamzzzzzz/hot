package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"
)

type Board struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"displayname"`
	Boards      []string `json:"boards"`
	Hots        []string
}

type Boardset struct {
	Boards map[string]*Board `json:"boards"`
}

type Archive struct {
	Name   string
	Date   string
	Boards []*Board
}

type Service struct {
	PathPrefix  string
	ArchivesDir string
	GitPull     bool
	*log.Logger
	mux      http.ServeMux
	render   *Render
	boardset *Boardset
	archive  *Archive
	mu       sync.RWMutex
}

func (s *Service) Init(ctx context.Context) error {
	if s.Logger == nil {
		s.Logger = log.Default()
	}
	if err := s.updating(); err != nil {
		return err
	}
	s.render = &Render{}
	s.mux.HandleFunc(s.PathPrefix+"/", s.HandleIndex)
	s.Printf("service init success. pathprefix=%s", s.PathPrefix)
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

func (s *Service) loadBoardset() (*Boardset, error) {
	boardset := &Boardset{}
	if b, err := os.ReadFile("website-board.json"); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	} else {
		if err := json.Unmarshal(b, &boardset); err != nil {
			return nil, err
		}
	}
	if boardset.Boards["all"] == nil {
		board := &Board{}
		entries, err := os.ReadDir(s.ArchivesDir)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			if entry.IsDir() {
				board.Boards = append(board.Boards, entry.Name())
			}
		}
		if boardset.Boards == nil {
			boardset.Boards = make(map[string]*Board)
		}
		boardset.Boards["all"] = board
	}
	for name, board := range boardset.Boards {
		if board.Name == "" {
			board.Name = name
		}
		if board.DisplayName == "" {
			board.DisplayName = board.Name
		}
	}
	return boardset, nil
}

func (s *Service) setBoardset(boardset *Boardset) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.boardset = boardset
}

func (s *Service) getBoardset() *Boardset {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.boardset
}

func (s *Service) getBoard(name string) *Board {
	boardset := s.getBoardset()
	if boardset == nil {
		return nil
	}
	return boardset.Boards[name]
}

func (s *Service) loadArchive(date string, boardname string) (*Archive, error) {
	board := s.getBoard(boardname)
	archivename := boardname
	if board != nil {
		archivename = board.DisplayName
	}
	archive := &Archive{
		Name: archivename,
		Date: date,
	}
	var boardnames []string
	if board != nil && board.Boards != nil {
		boardnames = board.Boards
	} else {
		boardnames = []string{boardname}
	}
	for _, name := range boardnames {
		file := filepath.Join(s.ArchivesDir, name, date+".txt")
		b, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		newboard := &Board{
			Name:        name,
			DisplayName: name,
		}
		for _, line := range strings.Split(string(b), "\n") {
			if line == "" {
				continue
			}
			newboard.Hots = append(newboard.Hots, line)
		}
		if board := s.getBoard(name); board != nil {
			if board.DisplayName != "" {
				newboard.DisplayName = board.DisplayName
			}
		}
		archive.Boards = append(archive.Boards, newboard)
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

func (s *Service) gitPull() ([]byte, error) {
	if !s.GitPull {
		return nil, nil
	}
	return exec.Command("git", "pull").CombinedOutput()
}

func (s *Service) update(pull bool) error {
	if pull {
		if output, err := s.gitPull(); err != nil {
			out := string(output)
			out = strings.ReplaceAll(out, "\n", "\\n")
			out = strings.ReplaceAll(out, "\r", "\\r")
			return fmt.Errorf("gitpull out:'%s' err:%w", out, err)
		}
	}
	boardset, err := s.loadBoardset()
	if err != nil {
		return err
	}
	s.setBoardset(boardset)
	archive, err := s.loadArchive(time.Now().Format("2006-01-02"), "default")
	if err != nil {
		return err
	}
	s.setArchive(archive)
	return nil
}

func (s *Service) updating() error {
	if err := s.update(false); err != nil {
		return err
	}
	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 30, 0, 0, now.Location()).Add(time.Hour)
			s.Printf("next update at %s", next)
			time.Sleep(next.Sub(now))
			if err := s.update(true); err != nil {
				s.Printf("update fail. err='%s'", err)
			} else {
				s.Printf("update success.")
			}
		}
	}()
	return nil
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
	pathprefix := flag.String("pathprefix", "/hot", "path prefix")
	archivesdir := flag.String("archivesdir", "archives", "archives dir")
	gitpull := flag.Bool("gitpull", true, "git pull")
	flag.Parse()
	service := &Service{
		PathPrefix:  *pathprefix,
		ArchivesDir: *archivesdir,
		GitPull:     *gitpull,
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
