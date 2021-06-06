package lalash

import (
	"context"
	"log"
	"strings"

	"lalash/command"
	"lalash/eval"
	"lalash/history"
	"lalash/parser"

	"github.com/peterh/liner"
)

const (
	historyFileName = ".lalash_history"
	exitCodeOK      = iota
	exitCodeErr
)

func Run() int {
	cmd := eval.Command(command.New())

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
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		expr, err := line.Prompt("$ ")
		if err != nil {
			log.Println("[read line error]", err)
			return exitCodeErr
		}
		line.AppendHistory(expr)

		tokens, err := parser.Parse(expr)
		if err != nil {
			log.Println("[parse error]", err)
		}

		if tokens == nil || len(tokens) == 0 || tokens[0].Val == "" {
			continue
		}

		if err := cmd.Eval(ctx, tokens); err != nil {
			log.Println("[eval error]", err)
			continue
		}
	}
}
