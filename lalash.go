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
	// parce
	argv := strings.Split(expr, " ")

	// exec
	cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
	cmd.Stdout = e.Out
	cmd.Stderr = e.Err
	return cmd.Run()
}
