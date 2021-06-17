package lalash

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/peterh/liner"
	"github.com/w-haibara/lalash/history"
)

const (
	historyFileName = ".lalash_history"
	exitCodeOK      = iota
	exitCodeErr
)

func Run(expr string) int {
	cmd := cmdNew()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := EvalString(ctx, cmd, expr); err != nil {
		fmt.Println(err.Error())
		return exitCodeErr
	}

	return exitCodeOK
}

func RunREPL() int {
	cmd := cmdNew()

	line := liner.NewLiner()
	defer line.Close()

	line.SetCompleter(func(line string) (c []string) {
		if len(strings.TrimSpace(line)) <= 0 {
			return nil
		}
		for _, v := range cmd.Internal.GetCmdsAll() {
			if strings.HasPrefix(v, line) {
				c = append(c, v)
			}
		}
		for _, v := range cmd.Internal.GetAliasAll() {
			if strings.HasPrefix(v, line) {
				c = append(c, v)
			}
		}
		return
	})

	line.SetCtrlCAborts(true)

	history := history.New(historyFileName)
	history.ReadHistory(line)
	defer history.WriteHistory(line)

	for {
		if err := func() error {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			expr, err := line.Prompt("$ ")
			if err != nil {
				return fmt.Errorf("[read line error] %v", err.Error())
			}
			line.AppendHistory(expr)

			return EvalString(ctx, cmd, expr)
		}(); err != nil {
			log.Println(err.Error())
		}
	}
}
