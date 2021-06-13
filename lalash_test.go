package lalash

import (
	"bufio"
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestEval(t *testing.T) {
	tests := []struct {
		name   string
		expr   string
		stdin  string
		stdout string
		stderr string
		err    error
	}{
		{
			name:   "echo1",
			expr:   "echo abc",
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "string literal1",
			expr:   `echo "a b c"`,
			stdin:  "",
			stdout: "a b c\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "stdin1",
			expr:   "wc",
			stdin:  "abc",
			stdout: "      0       1       3\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "stdin2",
			expr:   "grep a",
			stdin:  "abc",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "substitution1",
			expr:   `echo (echo abc)`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "substitution2",
			expr:   `echo (echo abc) (echo def)`,
			stdin:  "",
			stdout: "abc def\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "substitution3",
			expr:   `echo (echo (echo abc))`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "pipe1",
			expr:   `pipe {echo abc} cat`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "pipe2",
			expr:   `pipe {pipe {echo abc} cat} cat`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
	}
	for _, tt := range tests {
		cmd := cmdNew()

		i := strings.NewReader(tt.stdin)
		cmd.Stdin = i

		var out bytes.Buffer
		o := bufio.NewWriter(&out)
		cmd.Stdout = o

		var err bytes.Buffer
		e := bufio.NewWriter(&err)
		cmd.Stderr = e

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		t.Run(tt.name, func(t *testing.T) {
			if err := EvalString(ctx, cmd, tt.expr); err != tt.err {
				t.Errorf(err.Error())
			}

			o.Flush()
			e.Flush()

			if got := out.String(); got != tt.stdout {
				t.Errorf("\n=== Stdout ===\n%q\n---  want  ---\n%q\n--------------", got, tt.stdout)
			}

			if got := err.String(); got != tt.stderr {
				t.Errorf("\n=== Stderr ===\n%q\n---  want  ---\n%q\n--------------", got, tt.stderr)
			}
		})
	}
}