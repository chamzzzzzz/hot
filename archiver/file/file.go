package file

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/chamzzzzzz/hot"
)

type Archiver struct {
}

func (a *Archiver) Name() string {
	return "file-archiver"
}

func (a *Archiver) Archive(board *hot.Board) (archived int, err error) {
	os.MkdirAll(fmt.Sprintf("archives/%s", board.Name), 0755)
	name := fmt.Sprintf("archives/%s/%s.txt", board.Name, time.Now().Format("2006-01-02"))
	b, err := os.ReadFile(name)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("[%s] read archive file failed, err:%v\n", board.Name, err)
			return
		}
	}

	var words []string
	if len(b) > 0 {
		words = strings.Split(string(b), "\r\n")
	}

	n := 0
	for _, hot := range board.Hots {
		word := hot.Title
		word = strings.TrimSpace(word)
		word = strings.ReplaceAll(word, "\r\n", "")
		has := false
		for _, w := range words {
			if w == word {
				has = true
				break
			}
		}
		if !has {
			words = append(words, word)
			n++
		}
	}

	err = os.WriteFile(name, []byte(strings.Join(words, "\r\n")), 0755)
	if err != nil {
		log.Printf("[%s] write archive file failed, err:%v\n", board.Name, err)
		return
	}

	log.Printf("[%s] archived %d/%d new words\n", board.Name, n, len(board.Hots))
	return n, nil
}
