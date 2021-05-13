package env

import (
	"io"
)

type Env struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

