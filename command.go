package lalash

import (
	"io"
	"os"
)

type Command struct {
	Stdin    io.Reader
	Stdout   io.Writer
	Stderr   io.Writer
	Internal Internal
}

func cmdNew() Command {
	cmd := Command{
		Stdin:    os.Stdin,
		Stdout:   os.Stdout,
		Stderr:   os.Stderr,
		Internal: NewInternal(),
	}
	cmd.setInternalUtilFamily()
	cmd.setInternalAliasFamily()
	cmd.setInternalVarFamily()
	cmd.setInternalEvalFamily()
	return cmd
}
