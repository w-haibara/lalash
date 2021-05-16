package main

import (
	"context"
	"log"
	"os"

	"lalash/command"
	"lalash/env"
	"lalash/history"
	"lalash/process"

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

	env := env.Env{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}

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

		argv, err := Parse(expr)
		if err != nil {
			log.Println("[parse error]", err)
			return exitCodeErr
		}

		if argv == nil || len(argv) == 0 || argv[0] == "" {
			continue
		}

		if fn, ok := cmd[argv[0]]; ok {
			if err := fn(env, argv[0], argv[1:]...); err != nil {
				log.Println("[internal exec error]", err)
			}
			continue
		}

		if err := process.Exec(env, ctx, argv[0], argv[1:]...); err != nil {
			log.Println("[exec error]", err)
			continue
		}
	}
}
