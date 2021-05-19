package main

import (
	"context"
	"log"
	"os"

	"lalash/command"
	"lalash/history"

	"github.com/peterh/liner"
)

const (
	historyFileName = ".lalash_history"
	exitCodeOK      = iota
	exitCodeErr
)

func main() {
	os.Exit(Run())
}

func Run() int {
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	history := history.New(historyFileName)
	history.ReadHistory(line)
	defer history.WriteHistory(line)

	cmd := command.New()

	for {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		expr, err := line.Prompt("$ ")
		if err != nil {
			log.Println("[read line error]", err)
			return exitCodeErr
		}
		line.AppendHistory(expr)

		tokens, err := command.Parse(expr)
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
