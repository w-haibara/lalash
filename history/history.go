package history

import (
	"os"
	"path/filepath"

	"github.com/peterh/liner"
)

type HistoryFileName string

func New(name string) HistoryFileName {
	return HistoryFileName(name)
}

func (fn HistoryFileName) ReadHistory(line *liner.State) error {
	f, err := os.Open(filepath.Join(os.TempDir(), string(fn)))
	if err != nil {
		return err
	}
	line.ReadHistory(f)
	f.Close()
	return nil
}

func (fn HistoryFileName) WriteHistory(line *liner.State) error {
	f, err := os.Create(filepath.Join(os.TempDir(), string(fn)))
	if err != nil {
		return err
	} 
		line.WriteHistory(f)
		f.Close()
	return nil
}
