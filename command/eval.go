package command

import (
	"context"
	"fmt"
	"os/exec"
)

func (c Command) Eval(ctx context.Context, tokens []Token) error {
	argv := []string{}
	for _, v := range tokens {
		argv = append(argv, v.Val)
	}

	if err := c.Exec(ctx, argv); err != nil {
		return err
	}

	return nil
}

func (c Command) Exec(ctx context.Context, argv []string) error {
	if cmd, err := c.Internal.Get(argv[0]); err == nil {
		if err := cmd.Fn(c.Env, argv[0], argv[1:]...); err != nil {
			return fmt.Errorf("[internal exec error]", err)
		}
	}

	if err := Exec(c.Env, ctx, argv[0], argv[1:]...); err != nil {
		return fmt.Errorf("[exec error]", err)
	}

	return nil
}

func Exec(e Env, ctx context.Context, args string, argv ...string) error {
	cmd := exec.CommandContext(ctx, args, argv...)
	cmd.Stdin = e.In
	cmd.Stdout = e.Out
	cmd.Stderr = e.Err
	return cmd.Run()
}
