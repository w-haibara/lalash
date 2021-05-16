package commands

import (
	"lalash/env"
)

type InternalCmd map[string]func(env.Env, string, ...string) error

func New() InternalCmd {
	return InternalCmd{
	}
}
