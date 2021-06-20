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
		{
			name:   "echo2",
			expr:   "l-echo --fd=2 abc",
			stdin:  "",
			stdout: "",
			stderr: "abc\n",
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
		{
			name:   "pipe3",
			expr:   `l-pipe {l-echo --fd=2 abc} cat`,
			stdin:  "",
			stdout: "",
			stderr: "abc\n",
			err:    nil,
		},
		{
			name:   "pipe4",
			expr:   `l-pipe -p 2 {l-echo --fd=2 abc} cat`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "pipe5",
			expr:   `l-pipe -p 3 {l-echo --fd=3 abc} {cat}`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "pipe6",
			expr:   `l-pipe -p 9 {l-echo --fd=9 abc} {cat}`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "pipe7",
			expr:   `l-pipe -p 100 {l-echo --fd=100 abc} {cat}`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := cmdNew()

			i := strings.NewReader(tt.stdin)
			cmd.Stdin = i

			var stdout bytes.Buffer
			o := bufio.NewWriter(&stdout)
			cmd.Stdout = o

			var stderr bytes.Buffer
			e := bufio.NewWriter(&stderr)
			cmd.Stderr = e

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if err := EvalString(ctx, cmd, tt.expr); err != tt.err {
				t.Errorf(err.Error())
			}

			o.Flush()
			e.Flush()

			if got := stdout.String(); got != tt.stdout {
				t.Errorf("\n=== Stdout ===\n%q\n---  want  ---\n%q\n--------------", got, tt.stdout)
			}

			if got := stderr.String(); got != tt.stderr {
				t.Errorf("\n=== Stderr ===\n%q\n---  want  ---\n%q\n--------------", got, tt.stderr)
			}
		})
	}
}
