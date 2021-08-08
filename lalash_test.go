package lalash

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestEvalString(t *testing.T) {
	outputFiles := "./testfiles/out"
	if err := os.RemoveAll(outputFiles); err != nil {
		panic(err.Error())
	}
	if err := os.Mkdir(outputFiles, os.ModePerm); err != nil {
		panic(err.Error())
	}

	tests := []struct {
		name      string
		expr      string
		stdin     string
		stdout    string
		stderr    string
		inExtra   []string
		err       error
		checkFile func() error
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
			stdout: "abc",
			stderr: "",
			err:    nil,
		},
		{
			name:   "cat2",
			expr:   "l-cat --fd 0",
			stdin:  "abc",
			stdout: "abc",
			stderr: "",
			err:    nil,
		},
		{
			name:    "cat3",
			expr:    "l-cat --fd 3",
			stdin:   "",
			stdout:  "abc",
			stderr:  "",
			inExtra: []string{"abc"},
			err:     nil,
		},
		{
			name:    "cat3",
			expr:    "l-cat --fd 6",
			stdin:   "",
			stdout:  "abc",
			stderr:  "",
			inExtra: []string{"", "", "", "abc"},
			err:     nil,
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
		{
			name:   "raw-string literal2",
			expr:   `l-echo {a;b}`,
			stdin:  "",
			stdout: "a;b\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "raw-string literal3",
			expr:   `l-echo {a; c}`,
			stdin:  "",
			stdout: "a; c\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "raw-string literal4",
			expr:   `l-echo {a ;c}`,
			stdin:  "",
			stdout: "a ;c\n",
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
		{
			name:   "substitution4",
			expr:   `l-echo {(l-echo abc)}`,
			stdin:  "",
			stdout: "(l-echo abc)\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "substitution5",
			expr:   `l-echo "(l-echo abc)"`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "substitution6",
			expr:   `l-echo {"(l-echo abc)"}`,
			stdin:  "",
			stdout: "\"(l-echo abc)\"\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "substitution7",
			expr:   `l-eval {l-echo "(l-echo abc)"}`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "substitution8",
			expr:   `l-echo (l-echo aaa); l-echo bbb`,
			stdin:  "",
			stdout: "aaa\nbbb\n",
			stderr: "",
			err:    nil,
		},

		/*
			eval
		*/
		{
			name:   "eval1",
			expr:   `l-eval "echo abc"`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "eval2",
			expr:   `l-eval {echo abc}`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "eval3",
			expr:   `l-eval ""`,
			stdin:  "",
			stdout: "",
			stderr: "",
			err:    nil,
		},
		{
			name:   "eval4",
			expr:   `l-eval {echo} abc`,
			stdin:  "",
			stdout: "\n",
			stderr: "",
			err:    nil,
		}, {
			name:   "eval5",
			expr:   `l-eval {{echo} {abc}}`,
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
			expr:   `l-var aaa xxx; l-echo (l-var --ref aaa)`,
			stdin:  "",
			stdout: "xxx\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "var2",
			expr:   `l-var --mut aaa xxx; l-var --ch aaa yyy; l-echo (l-var --ref aaa)`,
			stdin:  "",
			stdout: "yyy\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "var3",
			expr:   `l-var aaa zzz; l-var --del aaa; l-var --mut aaa xxx; l-var --ch aaa yyy; l-echo (l-var --ref aaa)`,
			stdin:  "",
			stdout: "yyy\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "var4",
			expr:   `l-var aaa xxx; l-var --check aaa`,
			stdin:  "",
			stdout: "true\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "var5",
			expr:   `l-var --mut aaa xxx; l-var --check aaa`,
			stdin:  "",
			stdout: "true\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "var6",
			expr:   `l-var --global aaa xxx; l-var --check aaa`,
			stdin:  "",
			stdout: "true\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "var7",
			expr:   `l-var --global --mut aaa xxx; l-var --check aaa`,
			stdin:  "",
			stdout: "true\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "var8",
			expr:   `l-var --check aaa`,
			stdin:  "",
			stdout: "false\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "var9",
			expr:   `l-var bbb yyy; l-var --check aaa`,
			stdin:  "",
			stdout: "false\n",
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
			expr:   `l-pipe {l-echo abc} l-cat`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "pipe2",
			expr:   `l-pipe {l-pipe {l-echo abc} l-cat} l-cat`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "pipe3",
			expr:   `l-pipe {l-echo --fd=2 abc} l-cat`,
			stdin:  "",
			stdout: "",
			stderr: "abc\n",
			err:    nil,
		},
		{
			name:   "pipe4",
			expr:   `l-pipe -p 2 {l-echo --fd=2 abc} l-cat`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "pipe5",
			expr:   `l-pipe -p 3 {l-echo --fd=3 abc} l-cat`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "pipe6",
			expr:   `l-pipe -p 9 {l-echo --fd=9 abc} l-cat`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "pipe7",
			expr:   `l-pipe -p 100 {l-echo --fd=100 abc} l-cat`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "pipe8",
			expr:   `l-pipe -p 100 {l-echo --fd=100 abc} {l-cat --fd 0}`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "pipe9",
			expr:   `l-pipe -p 1:4 {l-echo abc} {l-cat --fd 4}`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "pipe10",
			expr:   `l-pipe -p 10:11 {l-echo --fd=10 abc} {l-cat --fd 11}`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "pipe11",
			expr:   `l-pipe --in ./testfiles/in/txt l-cat`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "pipe12",
			expr:   `l-pipe --out ./testfiles/out/txt {l-echo abc}`,
			stdin:  "",
			stdout: "",
			stderr: "",
			err:    nil,
			checkFile: func() error {
				txt := "abc\n"
				name := "./testfiles/out/txt"

				data, err := os.ReadFile(name)
				if err != nil {
					return err
				}
				if string(data) != txt {
					return fmt.Errorf("%q\n---  want  ---\n%q\n--------------", string(data), txt)
				}

				if err := os.Remove(name); err != nil {
					return err
				}

				return nil
			},
		},

		/*
			fn
		*/
		{
			name:   "fn1",
			expr:   `l-fn aaa {l-echo abc}; aaa`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "fn2",
			expr:   `l-fn aaa {l-fn bbb {l-echo abc}}; aaa; bbb`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "fn3",
			expr:   `l-fn aaa {l-var --global xxx abc}; aaa; l-var --ref xxx`,
			stdin:  "",
			stdout: "abc\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "fn4",
			expr:   `l-fn aaa {l-var xxx abc}; aaa; l-var --check xxx`,
			stdin:  "",
			stdout: "false\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "fn5",
			expr:   `l-fn aaa {echo (l-arg 0)}; aaa hello`,
			stdin:  "",
			stdout: "hello\n",
			stderr: "",
			err:    nil,
		},
		{
			name:   "fn6",
			expr:   `l-fn aaa {echo (l-arg 0); echo (l-arg 1)}; aaa hello world`,
			stdin:  "",
			stdout: "hello\nworld\n",
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

			if tt.inExtra != nil {
				cmd.ExtraFiles = make([]*os.File, len(tt.inExtra))
				for i, v := range tt.inExtra {
					if v == "" {
						continue
					}

					tmp, err := os.CreateTemp("", fmt.Sprint("lalash_test_extra_input_", i))
					if err != nil {
						panic(err.Error())
					}
					defer tmp.Close()
					defer os.Remove(tmp.Name())

					_, err = io.WriteString(tmp, v)
					if err != nil {
						panic(err.Error())
					}
					tmp.Close()

					tmp, err = os.Open(tmp.Name())
					if err != nil {
						panic(err.Error())
					}

					cmd.ExtraFiles[i] = tmp
				}
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if err := EvalString(ctx, cmd, tt.expr); err != tt.err {
				t.Errorf(err.Error())
			}

			o.Flush()
			e.Flush()

			if got := stdout.String(); got != tt.stdout {
				t.Errorf("%q\n=== Stdout ===\n%q\n---  want  ---\n%q\n--------------", tt.expr, got, tt.stdout)
			}

			if got := stderr.String(); got != tt.stderr {
				t.Errorf("%q\n=== Stderr ===\n%q\n---  want  ---\n%q\n--------------", tt.expr, got, tt.stderr)
			}

			if tt.checkFile != nil {
				if err := tt.checkFile(); err != nil {
					t.Errorf("\n===  File  ===\n%s", err.Error())
				}
			}
		})
	}
}
