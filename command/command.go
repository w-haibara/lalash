package command

import (
	"fmt"

	"lalash/env"
)

type InternalCmd map[string]func(env.Env, string, ...string) error

func New() InternalCmd {
	return InternalCmd{
		"echo": func(env env.Env, args string, argv ...string) error {
			str := ""
			for _, v := range argv {
				str += v + " "
			}
			fmt.Fprintln(env.Out, str)
			return nil
		},
	}
}
