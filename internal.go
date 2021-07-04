package lalash

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type InternalCmd struct {
	Usage string
	Fn    func(context.Context, Command, string, ...string) error
}

type Internal struct {
	Cmds   *sync.Map
	Alias  *sync.Map
	MutVar *sync.Map
	Var    *sync.Map
}

func NewInternal() Internal {
	in := Internal{
		Cmds:   new(sync.Map),
		Alias:  new(sync.Map),
		MutVar: new(sync.Map),
		Var:    new(sync.Map),
	}
	return in
}

func (in Internal) SetInternalCmd(name string, cmd InternalCmd) {
	in.Cmds.Store(name, cmd)
}

func checkArgv(argv []string, n int) error {
	if len(argv) < n {
		return fmt.Errorf("%d arguments required", n)
	}
	return nil
}

func (i Internal) GetAliasAll() []string {
	var s []string
	i.Alias.Range(func(key, value interface{}) bool {
		s = append(s, key.(string))
		return true
	})
	return s
}

func (i Internal) GetAlias(args string) string {
	if v, ok := i.Alias.Load(args); ok {
		return i.GetAlias(v.(string))
	}
	return args
}

func (i Internal) GetCmdsAll() []string {
	var s []string
	i.Cmds.Range(func(key, value interface{}) bool {
		s = append(s, key.(string))
		return true
	})
	return s
}

func (i Internal) Get(key string) (InternalCmd, error) {
	v, ok := i.Cmds.Load(key)
	if !ok {
		return InternalCmd{}, fmt.Errorf("command not found")
	}
	cmd, ok := v.(InternalCmd)
	if !ok {
		return InternalCmd{}, fmt.Errorf("function of the command is invalid")
	}
	return InternalCmd(cmd), nil
}

func sortJoin(s []string) string {
	sort.Strings(s)
	ret := ""
	for _, v := range s {
		ret += v + "\n"
	}
	return ret
}

func (cmd Command) setInternalUtilFamily() {
	cmd.Internal.Cmds.Store("l-help", InternalCmd{
		Usage: "l-help",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			var s []string
			cmd.Internal.Cmds.Range(func(key, value interface{}) bool {
				s = append(s, fmt.Sprintln(key, ":", value.(InternalCmd).Usage))
				return true
			})
			sort.Strings(s)
			ret := ""
			for _, v := range s {
				ret += v
			}
			fmt.Fprintln(cmd.Stdout, ret)
			return nil
		},
	})

	cmd.Internal.Cmds.Store("l-echo", InternalCmd{
		Usage: "l-echo",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			f := flag.NewFlagSet("echo", flag.ContinueOnError)
			fd := f.Int("fd", 1, "")
			if err := f.Parse(argv); err != nil {
				return err
			}

			out := cmd.Stdout
			switch {
			case *fd == 1:
			case *fd == 2:
				out = cmd.Stderr
			case *fd >= 3:
				out = cmd.ExtraFiles[*fd-3]
			default:
				return fmt.Errorf("invalid fd: %v", *fd)
			}
			fmt.Fprintln(out, strings.Join(f.Args(), " "))
			return nil
		},
	})

	cmd.Internal.Cmds.Store("l-cat", InternalCmd{
		Usage: "l-cat",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			f := flag.NewFlagSet("cat", flag.ContinueOnError)
			fd := f.Int("fd", 0, "")
			if err := f.Parse(argv); err != nil {
				return err
			}

			var src io.Reader
			switch {
			case *fd == 0:
				src = cmd.Stdin
			case *fd >= 3:
				src = cmd.ExtraFiles[*fd-3]
			default:
				return fmt.Errorf("invalid fd: %v", *fd)
			}

			if _, err := io.Copy(cmd.Stdout, src); err != nil {
				return err
			}
			return nil
		},
	})

	cmd.Internal.Cmds.Store("l-cd", InternalCmd{
		Usage: "l-cd <path>",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 1); err != nil {
				return err
			}
			if err := os.Chdir(argv[0]); err != nil {
				return err
			}
			return nil
		},
	})

	cmd.Internal.Cmds.Store("l-exit", InternalCmd{
		Usage: "l-exit",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			return shellExitErr
		},
	})
}

func (cmd Command) setInternalAliasFamily() {
	cmd.Internal.Cmds.Store("l-alias", InternalCmd{
		Usage: "",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			f := flag.NewFlagSet("alias", flag.ContinueOnError)
			isUnset := f.Bool("unset", false, "")
			isShow := f.Bool("show", false, "")
			if err := f.Parse(argv); err != nil {
				return err
			}

			if *isUnset && *isShow {
				return fmt.Errorf("cannot set both --unset and --show.")
			}

			switch {
			case !*isUnset && !*isShow:
				if f.Arg(0) == "" {
					return fmt.Errorf("key is blank")
				}
				if f.Arg(1) == "" {
					return fmt.Errorf("value is blank")
				}
				cmd.Internal.Alias.Store(f.Arg(0), f.Arg(1))
				return nil

			case *isUnset:
				if f.Arg(0) == "" {
					return fmt.Errorf("key is blank")
				}
				cmd.Internal.Alias.Delete(f.Arg(0))
				return nil

			case *isShow:
				s := []string{}
				cmd.Internal.Alias.Range(func(key, value interface{}) bool {
					s = append(s, fmt.Sprintf("%v : %v", key, value))
					return true
				})
				fmt.Fprint(cmd.Stdout, sortJoin(s))
				return nil
			}

			return nil
		},
	})

}

func (cmd Command) setInternalVarFamily() {
	cmd.Internal.SetInternalCmd("l-var", InternalCmd{
		Usage: "l-var",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			f := flag.NewFlagSet("var", flag.ContinueOnError)
			isMut := f.Bool("mut", false, "")
			isRef := f.Bool("ref", false, "")
			isCh := f.Bool("ch", false, "")
			isDel := f.Bool("del", false, "")
			isShow := f.Bool("show", false, "")
			if err := f.Parse(argv); err != nil {
				return err
			}

			if *isMut && (*isRef || *isCh || *isDel || *isShow) {
				return fmt.Errorf("cannot set --mut and others")
			}

			if *isRef && (*isCh || *isDel || *isShow) {
				return fmt.Errorf("cannot set --ref and others")
			}

			if *isCh && (*isDel || *isShow) {
				return fmt.Errorf("cannot set --ch and others")
			}
			if *isDel && *isShow {
				return fmt.Errorf("cannot set both --del and --show")
			}

			switch {
			case !*isMut && !*isRef && !*isCh && !*isDel && !*isShow:
				if f.Arg(0) == "" {
					return fmt.Errorf("key is blank")
				}
				if f.Arg(1) == "" {
					return fmt.Errorf("value is blank")
				}

				_, ok := cmd.Internal.Var.Load(f.Arg(0))
				if !ok {
					_, ok = cmd.Internal.MutVar.Load(f.Arg(0))
				}
				if ok {
					return fmt.Errorf("variable is already exists: %v", f.Arg(0))
				}
				cmd.Internal.Var.Store(f.Arg(0), f.Arg(1))
				return nil

			case *isMut && !*isRef && !*isCh && !*isDel && !*isShow:
				if f.Arg(0) == "" {
					return fmt.Errorf("key is blank")
				}
				if f.Arg(1) == "" {
					return fmt.Errorf("value is blank")
				}

				_, ok := cmd.Internal.Var.Load(f.Arg(0))
				if !ok {
					_, ok = cmd.Internal.MutVar.Load(f.Arg(0))
				}
				if ok {
					return fmt.Errorf("variable is already exists: %v", f.Arg(0))
				}
				cmd.Internal.MutVar.Store(f.Arg(0), f.Arg(1))
				return nil

			case !*isMut && *isRef && !*isCh && !*isDel && !*isShow:
				if f.Arg(0) == "" {
					return fmt.Errorf("key is blank")
				}

				v, ok := cmd.Internal.Var.Load(f.Arg(0))
				if !ok {
					v, ok = cmd.Internal.MutVar.Load(f.Arg(0))
				}
				if !ok {
					return fmt.Errorf("variable is not defined: %v", f.Arg(0))
				}
				fmt.Fprintln(cmd.Stdout, v)
				return nil

			case !*isMut && !*isRef && *isCh && !*isDel && !*isShow:
				if f.Arg(0) == "" {
					return fmt.Errorf("key is blank")
				}

				if _, ok := cmd.Internal.Var.Load(f.Arg(0)); ok {
					return fmt.Errorf("variable is immutable: %v", f.Arg(0))
				}
				if _, ok := cmd.Internal.MutVar.Load(f.Arg(0)); !ok {
					return fmt.Errorf("variable is not defined: %v", f.Arg(0))
				}
				cmd.Internal.MutVar.Store(f.Arg(0), f.Arg(1))
				return nil

			case !*isMut && !*isRef && !*isCh && *isDel && !*isShow:
				if f.Arg(0) == "" {
					return fmt.Errorf("key is blank")
				}

				cmd.Internal.Var.Delete(f.Arg(0))
				cmd.Internal.MutVar.Delete(f.Arg(0))
				return nil

			case !*isMut && !*isRef && !*isCh && !*isDel && *isShow:
				fmt.Fprintln(cmd.Stdout, "[mutable variables]")
				s1 := []string{}
				cmd.Internal.MutVar.Range(func(key, value interface{}) bool {
					fmt.Fprintln(cmd.Stdout, key, ":", value)
					return true
				})
				fmt.Fprint(cmd.Stdout, sortJoin(s1))

				fmt.Fprintln(cmd.Stdout, "\n[immutable variables]")
				s2 := []string{}
				cmd.Internal.Var.Range(func(key, value interface{}) bool {
					fmt.Fprintln(cmd.Stdout, key, ":", value)
					return true
				})
				fmt.Fprint(cmd.Stdout, sortJoin(s2))

				return nil
			}

			return nil
		},
	})
}

func (cmd Command) setInternalEvalFamily() {
	cmd.Internal.Cmds.Store("l-eval", InternalCmd{
		Usage: "l-eval",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 1); err != nil {
				return err
			}
			if err := EvalString(ctx, cmd, argv[0]); err != nil {
				return err
			}
			return nil
		},
	})

	cmd.Internal.Cmds.Store("l-pipe", InternalCmd{
		Usage: "l-pipe",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			f := flag.NewFlagSet("pipe", flag.ContinueOnError)
			p := f.String("p", "", "")
			in := f.String("in", "", "")
			out := f.String("out", "", "")
			if err := f.Parse(argv); err != nil {
				return err
			}

			if *in != "" && *out != "" {
				return fmt.Errorf("can not set both --in and --out")
			}

			if (*in != "" || *out != "") && *p != "" {
				return fmt.Errorf("can not set both -p and --in / --out")
			}

			if *in != "" {
				cmd1 := cmd
				i, err := filepath.Abs(*in)
				if err != nil {
					return err
				}
				file, err := os.Open(i)
				if err != nil {
					return err
				}
				cmd1.Stdin = file
				if err := EvalString(ctx, cmd1, f.Arg(0)); err != nil {
					return err
				}
				if err := file.Close(); err != nil {
					return err
				}
				return nil
			}

			if *out != "" {
				cmd1 := cmd
				o, err := filepath.Abs(*out)
				if err != nil {
					return err
				}
				file, err := os.OpenFile(o, os.O_WRONLY|os.O_CREATE, 0666)
				if err != nil {
					return err
				}
				cmd1.Stdout = file
				if err := EvalString(ctx, cmd1, f.Arg(0)); err != nil {
					return err
				}
				if err := file.Close(); err != nil {
					return err
				}
				return nil
			}

			if *p == "" {
				*p = "1:0"
			}

			type pair struct {
				in, out int64
			}

			v, err := func(expr string) ([]pair, error) {
				res := []pair{}

				for _, v1 := range strings.Split(expr, ",") {
					if v1 == "" {
						continue
					}

					n := strings.SplitN(v1, ":", 2)

					for i, v := range n {
						n[i] = strings.TrimSpace(v)
					}

					tmp := pair{}

					if in, err := strconv.ParseInt(n[0], 10, 32); err == nil {
						tmp.in = in
					} else {
						return nil, err
					}

					if len(n) >= 2 {
						if out, err := strconv.ParseInt(n[1], 10, 32); err == nil {
							tmp.out = out
						} else {
							return nil, err
						}
					}

					res = append(res, tmp)
				}

				return res, nil
			}(*p)
			if err != nil {
				return err
			}

			for _, v := range v {
				if err := func() error {
					r, w, err := os.Pipe()
					if err != nil {
						panic(err)
					}
					defer r.Close()
					defer w.Close()

					w.Chmod(os.ModeNamedPipe)

					cmd1 := cmd
					switch {
					case v.in == 1:
						cmd1.Stdout = w
					case v.in == 2:
						cmd1.Stderr = w
					case v.in >= 3:
						cmd1.ExtraFiles = make([]*os.File, v.in-2)
						cmd1.ExtraFiles[v.in-3] = w
					default:
						return fmt.Errorf("invalid input fd: %v", v.in)
					}

					if err := EvalString(ctx, cmd1, f.Arg(0)); err != nil {
						return err
					}
					w.Close()

					cmd2 := cmd
					if v.out < 3 {
						if v.out != 0 {
							return fmt.Errorf("invalid output fd: %v", v.out)
						}

						out := new(bytes.Buffer)
						io.Copy(out, r)
						r.Close()
						cmd2.Stdin = out
					} else {
						cmd2.ExtraFiles = make([]*os.File, v.out-2)
						cmd2.ExtraFiles[v.out-3] = r
					}

					if err := EvalString(ctx, cmd2, f.Arg(1)); err != nil {
						return err
					}

					return nil
				}(); err != nil {
					return err
				}
			}
			return nil
		},
	})
}
