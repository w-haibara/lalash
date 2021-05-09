package cli

import (
	"log"
	"os"
	"path/filepath"

	"github.com/peterh/liner"
)

type historyFileName string

func (fn historyFileName) readHistory(line *liner.State) {
	if f, err := os.Open(filepath.Join(os.TempDir(), string(fn))); err == nil {
		line.ReadHistory(f)
		f.Close()
	}
}

func (fn historyFileName) writeHistory(line *liner.State) {
	if f, err := os.Create(filepath.Join(os.TempDir(), string(fn))); err != nil {
		log.Print("Error writing history file: ", err)
	} else {
		line.WriteHistory(f)
		f.Close()
	}
}
