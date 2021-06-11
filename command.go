package lalash

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
	Internal Internal
}

func cmdNew() Command {
	cmd := Command{
		Env: Env{
			In:  os.Stdin,
			Out: os.Stdout,
			Err: os.Stderr,
		},
		Internal: NewInternal(),
	}
	cmd.setHelp()
	cmd.setExit()
	cmd.setAliasFamily()
	cmd.setVarFamily()
	return cmd
}
