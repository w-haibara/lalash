package lalash

import (
	"context"
	"io"
	"os/exec"
	"strings"
)

type Env struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

func (e Env) Eval(ctx context.Context, expr string) error {
	argv, err := e.parce(expr)
	if err != nil {
		return err
	}

	if err := e.exec(ctx, argv[0], argv[1:]...); err != nil {
		return err
	}

	return nil
}

func (e Env) parce(expr string) ([]string, error) {
	return strings.Split(expr, " "), nil
}

func (e Env) exec(ctx context.Context, args string, argv ...string) error {
	cmd := exec.CommandContext(ctx, args, argv...)

	cmd.Stdin = e.In
	cmd.Stdout = e.Out
	cmd.Stderr = e.Err

	return cmd.Run()
}
