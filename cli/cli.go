package cli

import (
	"context"
	"log"
	"os"

	"lalash"

	"github.com/peterh/liner"
)

const (
	history    historyFileName = ".lalash_history"
	exitCodeOK                 = iota
	exitCodeErr
)

func Run() int {
	line := liner.NewLiner()
	defer line.Close()
	line.SetCtrlCAborts(true)
	history.readHistory(line)
	defer history.writeHistory(line)

	env := lalash.Env{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}

	for {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		expr, err := line.Prompt("$ ")
		if err != nil {
			log.Println("[read line error]", err)
			return exitCodeErr
		}
		line.AppendHistory(expr)

		argv, err := env.Parse(expr)
		if err != nil {
			log.Println("[parse error]", err)
			return exitCodeErr
		}

		if argv == nil || argv[0] == "" {
			continue
		}

		if err := env.Exec(ctx, argv[0], argv[1:]...); err != nil {
			log.Println("[exec error]", err)
			return exitCodeErr
		}
	}
}
