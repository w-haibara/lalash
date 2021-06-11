package lalash

import (
	"context"
	"fmt"
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

func RunREPL() int {
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
		if err := readAndEval(cmd, line); err != nil {
			log.Println(err.Error())
		}
	}
}

func readAndEval(cmd eval.Command, line *liner.State) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	expr, err := line.Prompt("$ ")
	if err != nil {
		return fmt.Errorf("[read line error] %v", err.Error())
	}
	line.AppendHistory(expr)

	return Eval(ctx, cmd, expr)
}

func Eval(ctx context.Context, cmd eval.Command, expr string) error {
	tokens, err := parser.Parse(expr)
	if err != nil {
		return fmt.Errorf("[parse error] %v", err.Error())
	}

	if tokens == nil || len(tokens) == 0 || tokens[0].Val == "" {
		return nil
	}

	if err := cmd.Eval(ctx, tokens); err != nil {
		return fmt.Errorf("[eval error] %v", err.Error())
	}

	return nil
}
