package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"strings"

	"lalash/command"
	"lalash/env"
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

		if argv == nil || argv[0] == "" {
			continue
		}

		if fn, ok := cmd[argv[0]]; ok {
			if err := fn(env, argv[0], argv[1:]...); err != nil {
				log.Println("[internal exec error]", err)
			}
			continue
		}

		if err := Exec(env, ctx, argv[0], argv[1:]...); err != nil {
			log.Println("[exec error]", err)
			continue
		}
	}
}

func Parse(expr string) ([]string, error) {
	return strings.Split(expr, " "), nil
}

func Exec(e env.Env, ctx context.Context, args string, argv ...string) error {
	cmd := exec.CommandContext(ctx, args, argv...)
	cmd.Stdin = e.In
	cmd.Stdout = e.Out
	cmd.Stderr = e.Err
	return cmd.Run()
}
