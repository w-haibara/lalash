package lalash

import (
	"bufio"
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestEvalString(t *testing.T) {
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
			name:   "raw-string literal1",
			expr:   `echo {a b c}`,
			stdin:  "",
			stdout: "a b c\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "split1",
			expr:   `echo abc;`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "split2",
			expr:   `echo abc; echo def`,
			stdin:  "",
			stdout: "abc\ndef\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "split3",
			expr:   `echo abc ; echo def`,
			stdin:  "",
			stdout: "abc\ndef\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "split4",
			expr:   `echo abc;echo def`,
			stdin:  "",
			stdout: "abc;echo def\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "comment1",
			expr:   `echo abc #this is a comment message`,
			stdin:  "",
			stdout: "abc\n",
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
			name:   "alias1",
			expr:   `l-alias -k aaa -v {echo bbb}; aaa`,
			stdin:  "",
			stdout: "bbb\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "alias2",
			expr:   `l-alias -k aaa -v xxx; l-alias -k bbb -v yyy; l-alias --show`,
			stdin:  "",
			stdout: "aaa : xxx\nbbb : yyy\n",
			stderr: "",
			err:    nil,
		},

		{
			name:   "pipe1",
			expr:   `l-pipe {echo abc} cat`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "pipe2",
			expr:   `l-pipe {l-pipe {echo abc} cat} cat`,
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
				t.Errorf("Command: %q\n=== Stdout ===\n%q\n---  want  ---\n%q\n--------------", tt.expr, got, tt.stdout)
			}

			if got := err.String(); got != tt.stderr {
				t.Errorf("Command: %q\n=== Stderr ===\n%q\n---  want  ---\n%q\n--------------", tt.expr, got, tt.stderr)
			}
		})
	}
}
