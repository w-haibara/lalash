package lalash

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/w-haibara/lalash/parser"
)

func EvalString(ctx context.Context, cmd Command, expr string) error {
	tokens, err := parser.Parse(expr)
	if err != nil {
		return fmt.Errorf("[parse error] %v", err.Error())
	}

	start := 0
	for i := 1; i <= len(tokens); i++ {
		if tokens[i-1].Kind == parser.SeparateToken {
			if err := eval(ctx, cmd, tokens[start:i-1]); err != nil {
				return fmt.Errorf("[eval error] %v", err.Error())
			}
			start = i
			if start >= len(tokens) {
				break
			}
			continue
		}

		if i == len(tokens) {
			if err := eval(ctx, cmd, tokens[start:i]); err != nil {
				return fmt.Errorf("[eval error] %v", err.Error())
			}
		}
	}

	return nil
}

func eval(ctx context.Context, cmd Command, tokens []parser.Token) error {
	if tokens == nil || len(tokens) == 0 || tokens[0].Val == "" {
		return nil
	}

	argv := []string{}
	for i, v := range tokens {
		if v.Kind == parser.SubstitutionToken {
			res, err := func() (string, error) {
				var b bytes.Buffer
				w := bufio.NewWriter(&b)
				cmd := cmd
				cmd.Stdout = w

				tokens, err := parser.Parse(v.Val)
				if err != nil {
					return "", fmt.Errorf("[parse error] %v", err.Error())
				}

				if tokens == nil || len(tokens) == 0 || tokens[0].Val == "" {
					return "", nil
				}

				if err := eval(ctx, cmd, tokens); err != nil {
					return "", fmt.Errorf("[eval error] %v", err.Error())
				}

				w.Flush()

				return strings.TrimSpace(strings.ReplaceAll(b.String(), "\n", " ")), nil
			}()
			if err != nil {
				return err
			}

			tokens[i].Val = res
			tokens[i].Kind = parser.CommandToken
		}
		argv = append(argv, tokens[i].Val)
	}

	if err := Exec(ctx, cmd, argv); err != nil {
		return err
	}

	return nil
}

func Exec(ctx context.Context, cmd Command, argv []string) error {
	if alias := cmd.Internal.GetAlias(argv[0]); alias != argv[0] {
		return EvalString(ctx, cmd, alias)
	}

	if c, err := cmd.Internal.Get(argv[0]); err == nil {
		if err := c.Fn(ctx, cmd, argv[0], argv[1:]...); err != nil {
			return fmt.Errorf("[internal exec error] %v", err.Error())
		}
		return nil
	}

	if err := func() error {
		c := exec.CommandContext(ctx, argv[0], argv[1:]...)
		c.Stdin = cmd.Stdin
		c.Stdout = cmd.Stdout
		c.Stderr = cmd.Stderr
		return c.Run()
	}(); err != nil {
		return fmt.Errorf("[exec error] %v", err.Error())
	}
	return nil
}
