package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
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
	Name        string
	DisplayName string
	Date        string
	Boards      []*Board
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
	pv       int
	uv       map[string]int
	resetuv  time.Time
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

func (s *Service) addUV(ip net.IP) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.resetuv.IsZero() || s.resetuv.Format("2006-01-02") != time.Now().Format("2006-01-02") {
		if !s.resetuv.IsZero() {
			s.Printf("date[%s] pv[%d] uv[%d]", s.resetuv.Format("2006-01-02"), s.pv, len(s.uv))
		}
		s.resetuv = time.Now()
		s.uv = make(map[string]int)
		s.pv = 0
	}
	s.pv++
	s.uv[ip.String()]++
	return s.uv[ip.String()]
}

func (s *Service) HandleIndex(w http.ResponseWriter, r *http.Request) {
	boardname := r.URL.Query().Get("board")
	if boardname == "" {
		boardname = "default"
	}
	ip := s.ip(r)
	uv := s.addUV(ip)
	maxuv := 10
	if uv > maxuv {
		s.Printf("[%s] uv[%d] limit to access board[%s]", ip, uv, boardname)
		s.limit(w, maxuv)
		return
	}
	s.index(w, s.newArchive(boardname))
	s.Printf("[%s] uv[%d] access board[%s]", ip, uv, boardname)
}

func (s *Service) index(w http.ResponseWriter, archive *Archive) {
	b, err := s.render.Index(archive)
	if err != nil {
		panic(err)
	}
	w.Write(b)
}

func (s *Service) limit(w http.ResponseWriter, maxuv int) {
	b, err := s.render.Limit(maxuv)
	if err != nil {
		panic(err)
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
	archive := &Archive{
		Name:        boardname,
		DisplayName: boardname,
		Date:        date,
	}
	board := s.getBoard(boardname)
	if board != nil {
		if board.DisplayName != "" {
			archive.DisplayName = board.DisplayName
		}
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
		for _, line := range strings.Split(string(b), "\r\n") {
			if line == "" {
				continue
			}
			newboard.Hots = append(newboard.Hots, strings.TrimSpace(line))
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

func (s *Service) newArchive(boardname string) *Archive {
	archive := s.getArchive()
	newarchive := &Archive{
		Name:        boardname,
		DisplayName: boardname,
		Date:        archive.Date,
	}
	board := s.getBoard(boardname)
	if board != nil {
		if board.DisplayName != "" {
			newarchive.DisplayName = board.DisplayName
		}
	}
	var boardnames []string
	if board != nil && board.Boards != nil {
		boardnames = board.Boards
	} else {
		boardnames = []string{boardname}
	}
	for _, name := range boardnames {
		for _, b := range archive.Boards {
			if b.Name == name {
				newboard := b
				if len(b.Hots) > 10 {
					newboard = &Board{
						Name:        b.Name,
						DisplayName: b.DisplayName,
						Hots:        b.Hots[len(b.Hots)-10:],
					}
				}
				newarchive.Boards = append(newarchive.Boards, newboard)
				break
			}
		}
	}
	return newarchive
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
	archive, err := s.loadArchive(time.Now().Format("2006-01-02"), "all")
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

func (s *Service) ip(r *http.Request) net.IP {
	ip := net.ParseIP(r.Header.Get("X-Real-IP"))
	if ip != nil {
		return ip
	}
	if forward := r.Header.Get("X-Forwarded-For"); forward != "" {
		for _, forwardip := range strings.Split(forward, ",") {
			ip = net.ParseIP(forwardip)
			if ip != nil {
				return ip
			}
		}
	}
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return net.ParseIP(host)
}

type Render struct {
}

func (r *Render) Index(archvie *Archive) ([]byte, error) {
	source := `<!DOCTYPE html>
<html>
	<head>
		<title>热门榜单</title>
		<style>
			main {
				max-width: 55rem;
				margin: 0 auto 1.5rem;
				padding: 0 1.5rem;
				font-family: arial, sans-serif;
			}
		</style>
	</head>
	<body>
		<main>
			<h1>热门榜单</h1>
			<p>{{.Date}}</p>
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

func (r *Render) Limit(maxuv int) ([]byte, error) {
	source := `<!DOCTYPE html>
<html>
	<head>
		<title>热门榜单</title>
		<style>
			body {
				background-color: white;
				color: black;
				font-family: arial, sans-serif;
				margin: 0;
				padding: 0;
				display: flex;
				align-items: center;
				justify-content: center;
				min-height: 100vh;
				flex-direction: column;
			}
			h1 {
				font-size: 24px;
				margin: 0;
				padding: 0;
			}
			h2 {
				font-size: 12px;
				margin: 0;
				padding: 0;
			}
		</style>
	</head>
	<body>
		<h1>每日{{.}}次访问额度已用完^_^</h1>
		<h2>HOT WILL EVENTUALLY COOL......</h2>
	</body>
</html>
`
	tpl := template.Must(template.New("index").Parse(source))
	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, maxuv)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func main() {
	addr := flag.String("addr", ":8080", "address")
	pathprefix := flag.String("pathprefix", "", "path prefix")
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
