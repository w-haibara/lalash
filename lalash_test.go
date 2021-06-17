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
		/*
			basic echo
		*/
		{
			name:   "echo1",
			expr:   "l-echo abc",
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},

		/*
			basic cat
		*/
		{
			name:   "cat1",
			expr:   "l-cat",
			stdin:  "abc",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},

		/*
			string literal
		*/
		{
			name:   "string literal1",
			expr:   `l-echo "a b c"`,
			stdin:  "",
			stdout: "a b c\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "raw-string literal1",
			expr:   `l-echo {a b c}`,
			stdin:  "",
			stdout: "a b c\n",
			stderr: "",
			err:    nil,
		},

		/*
			separate
		*/
		{
			name:   "separate1",
			expr:   `l-echo abc;`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "separate2",
			expr:   `l-echo abc; l-echo def`,
			stdin:  "",
			stdout: "abc\ndef\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "separate3",
			expr:   `l-echo abc ; l-echo def`,
			stdin:  "",
			stdout: "abc\ndef\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "separate4",
			expr:   `l-echo abc;l-echo def`,
			stdin:  "",
			stdout: "abc;l-echo def\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "separate5",
			expr:   `l-echo abc; l-echo def; l-echo ghi`,
			stdin:  "",
			stdout: "abc\ndef\nghi\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "separate6",
			expr:   `l-echo abc def; l-echo ghi jkl; l-echo mno`,
			stdin:  "",
			stdout: "abc def\nghi jkl\nmno\n",
			stderr: "",
			err:    nil,
		},

		/*
			comment
		*/
		{
			name:   "comment1",
			expr:   `l-echo abc #this is a comment message`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},

		/*
			command substitution
		*/
		{
			name:   "substitution1",
			expr:   `l-echo (l-echo abc)`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "substitution2",
			expr:   `l-echo (l-echo abc) (l-echo def)`,
			stdin:  "",
			stdout: "abc def\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "substitution3",
			expr:   `l-echo (l-echo (l-echo abc))`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},

		/*
			var
		*/
		{
			name:   "var1",
			expr:   `l-var aaa xxx; echo (l-var --ref aaa)`,
			stdin:  "",
			stdout: "xxx\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "var2",
			expr:   `l-var --mut aaa xxx; l-var --ch aaa yyy; echo (l-var --ref aaa)`,
			stdin:  "",
			stdout: "yyy\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "var3",
			expr:   `l-var aaa zzz; l-var --del aaa; l-var --mut aaa xxx; l-var --ch aaa yyy; echo (l-var --ref aaa)`,
			stdin:  "",
			stdout: "yyy\n",
			stderr: "",
			err:    nil,
		},

		/*
			alias
		*/
		{
			name:   "alias1",
			expr:   `l-alias aaa {l-echo bbb}; aaa`,
			stdin:  "",
			stdout: "bbb\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "alias2",
			expr:   `l-alias aaa xxx; l-alias bbb yyy; l-alias --show`,
			stdin:  "",
			stdout: "aaa : xxx\nbbb : yyy\n",
			stderr: "",
			err:    nil,
		},

		/*
			pipe
		*/
		{
			name:   "pipe1",
			expr:   `l-pipe {l-echo abc} cat`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "pipe2",
			expr:   `l-pipe {l-pipe {l-echo abc} cat} cat`,
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
