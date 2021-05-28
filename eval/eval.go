package eval

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"lalash/command"
	"lalash/parser"
)

type Command command.Command

func (c Command) Eval(ctx context.Context, tokens []parser.Token) error {
	argv := []string{}
	for i, v := range tokens {
		if v.Kind == parser.SubstitutionToken {
			res, err := func() (string, error) {
				var b bytes.Buffer
				w := bufio.NewWriter(&b)
				c := c
				c.Env.Out = w

				tokens, err := parser.Parse(v.Val)
				if err != nil {
					return "", fmt.Errorf("[parse error]", err)
				}

				if tokens == nil || len(tokens) == 0 || tokens[0].Val == "" {
					return "", nil
				}

				if err := c.Eval(ctx, tokens); err != nil {
					return "", fmt.Errorf("[eval error]", err)
				}

				w.Flush()

				return strings.TrimSpace(strings.ReplaceAll(b.String(), "\n", " ")), nil
			}()
			if err != nil {
				return err
			}

			tokens[i].Val = res
			tokens[i].Kind = parser.StringToken
		}
		argv = append(argv, tokens[i].Val)
	}

	if err := c.Exec(ctx, argv); err != nil {
		return err
	}

	return nil
}

func (c Command) Exec(ctx context.Context, argv []string) error {
	argv[0] = c.Internal.GetAlias(argv[0])

	if cmd, err := c.Internal.Get(argv[0]); err == nil {
		if err := cmd.Exec(c.Env, argv[0], argv[1:]...); err != nil {
			return fmt.Errorf("[internal exec error]", err)
		}
		return nil
	}

	if err := Exec(c.Env, ctx, argv[0], argv[1:]...); err != nil {
		return fmt.Errorf("[exec error]", err)
	}

	return nil
}

func Exec(e command.Env, ctx context.Context, args string, argv ...string) error {
	cmd := exec.CommandContext(ctx, args, argv...)
	cmd.Stdin = e.In
	cmd.Stdout = e.Out
	cmd.Stderr = e.Err
	return cmd.Run()
}
