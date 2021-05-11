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

func (e Env) Parse(expr string) ([]string, error) {
	return strings.Split(expr, " "), nil
}

func (e Env) Exec(ctx context.Context, args string, argv ...string) error {
	cmd := exec.CommandContext(ctx, args, argv...)

	cmd.Stdin = e.In
	cmd.Stdout = e.Out
	cmd.Stderr = e.Err

	return cmd.Run()
}
