package command

import (
	"io"
	"os"
)

type Env struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

type Command struct {
	Env      Env
	Internal InternalCmdMap
}

func New() Command {
	return Command{
		Env: Env{
			In:  os.Stdin,
			Out: os.Stdout,
			Err: os.Stderr,
		},
		Internal: NewInternalCmdMap(),
	}
}
