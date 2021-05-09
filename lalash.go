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

func Eval(env Env, expr string, ctx context.Context) error {
	// parce
	argv := strings.Split(expr, " ")

	// exec
	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Stdout = env.Out
	cmd.Stderr = env.Err
	return cmd.Run()
}
