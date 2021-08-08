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
	Cmds         *sync.Map
	Alias        *sync.Map
	MutVar       *sync.Map
	Var          *sync.Map
	GlobalMutVar *sync.Map
	GlobalVar    *sync.Map
	Args         *sync.Map
	Return       *sync.Map
}

func NewInternal() Internal {
	in := Internal{
		Cmds:         new(sync.Map),
		Alias:        new(sync.Map),
		MutVar:       new(sync.Map),
		Var:          new(sync.Map),
		GlobalMutVar: new(sync.Map),
		GlobalVar:    new(sync.Map),
		Args:         new(sync.Map),
		Return:       new(sync.Map),
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

func (i Internal) arg(n int) string {
	if v, ok := i.Args.Load(n); ok {
		if v, ok := v.(string); ok {
			return v
		}
	}
	return ""
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

	cmd.Internal.Cmds.Store("l-fn", InternalCmd{
		Usage: "l-fn",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			f := flag.NewFlagSet("fn", flag.ContinueOnError)
			if err := f.Parse(argv); err != nil {
				return err
			}

			if err := checkArgv(f.Args(), 2); err != nil {
				return err
			}

			if f.Arg(0) == "" {
				return fmt.Errorf("function name is blank")
			}

			if f.Arg(1) == "" {
				return fmt.Errorf("function body is blank")
			}

			cmd.Internal.Alias.Store(f.Arg(0), "l-eval {"+f.Arg(1)+"}")

			return nil
		},
	})

	cmd.Internal.Cmds.Store("l-arg", InternalCmd{
		Usage: "l-arg",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 1); err != nil {
				return err
			}

			n, err := strconv.Atoi(argv[0])
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.Stdout, cmd.Internal.arg(n))

			return nil
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

func storeVarToMap(m *sync.Map, name, value string) {
	m.Store(name, value)
}

func loadVarFromMap(m *sync.Map, name string) (string, bool) {
	if v, ok := m.Load(name); ok {
		if v, ok := v.(string); ok {
			return v, true
		}
	}
	return "", false
}

func (cmd Command) setInternalVarFamily() {
	cmd.Internal.SetInternalCmd("l-var", InternalCmd{
		Usage: "l-var",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {

			loadVar := func(name string) (string, bool) {
				if v, ok := loadVarFromMap(cmd.Internal.Var, name); ok {
					return v, true
				}

				if v, ok := loadVarFromMap(cmd.Internal.MutVar, name); ok {
					return v, true
				}

				if v, ok := loadVarFromMap(cmd.Internal.GlobalVar, name); ok {
					return v, true
				}

				if v, ok := loadVarFromMap(cmd.Internal.GlobalMutVar, name); ok {
					return v, true
				}

				return "", false
			}

			f := flag.NewFlagSet("var", flag.ContinueOnError)
			isMut := f.Bool("mut", false, "")
			isRef := f.Bool("ref", false, "")
			isCh := f.Bool("ch", false, "")
			isDel := f.Bool("del", false, "")
			isShow := f.Bool("show", false, "")
			isGlobal := f.Bool("global", false, "")
			isCheck := f.Bool("check", false, "")

			if err := f.Parse(argv); err != nil {
				return err
			}

			if *isCheck {
				if err := checkArgv(f.Args(), 1); err != nil {
					return err
				}
				_, ok := loadVar(f.Arg(0))
				fmt.Fprintln(cmd.Stdout, ok)
				return nil
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

			var varMap = cmd.Internal.Var
			if *isGlobal {
				varMap = cmd.Internal.GlobalVar
			}

			var mutVarMap = cmd.Internal.MutVar
			if *isGlobal {
				mutVarMap = cmd.Internal.GlobalMutVar
			}

			switch {
			case !*isMut && !*isRef && !*isCh && !*isDel && !*isShow:
				if f.Arg(0) == "" {
					return fmt.Errorf("key is blank")
				}
				if f.Arg(1) == "" {
					return fmt.Errorf("value is blank")
				}

				if _, ok := loadVar(f.Arg(0)); ok {
					return fmt.Errorf("variable is already exists: %v", f.Arg(0))
				}

				storeVarToMap(varMap, f.Arg(0), f.Arg(1))
				return nil

			case *isMut && !*isRef && !*isCh && !*isDel && !*isShow:
				if f.Arg(0) == "" {
					return fmt.Errorf("key is blank")
				}

				if f.Arg(1) == "" {
					return fmt.Errorf("value is blank")
				}

				if _, ok := loadVar(f.Arg(0)); ok {
					return fmt.Errorf("variable is already exists: %v", f.Arg(0))
				}

				storeVarToMap(mutVarMap, f.Arg(0), f.Arg(1))

				return nil

			case !*isMut && *isRef && !*isCh && !*isDel && !*isShow:
				if f.Arg(0) == "" {
					return fmt.Errorf("key is blank")
				}

				v, ok := loadVar(f.Arg(0))
				if !ok {
					return fmt.Errorf("variable is not defined: %v", f.Arg(0))
				}
				fmt.Fprintln(cmd.Stdout, v)
				return nil

			case !*isMut && !*isRef && *isCh && !*isDel && !*isShow:
				if f.Arg(0) == "" {
					return fmt.Errorf("key is blank")
				}

				if _, ok := loadVar(f.Arg(0)); !ok {
					return fmt.Errorf("variable is not defined: %v", f.Arg(0))
				}

				if _, ok := loadVarFromMap(cmd.Internal.Var, f.Arg(0)); ok {
					return fmt.Errorf("variable is immutable: %v", f.Arg(0))
				}

				if _, ok := loadVarFromMap(cmd.Internal.Var, f.Arg(0)); ok {
					return fmt.Errorf("variable is immutable: %v", f.Arg(0))
				}

				storeVarToMap(varMap, f.Arg(0), f.Arg(1))

				return nil

			case !*isMut && !*isRef && !*isCh && *isDel && !*isShow:
				if f.Arg(0) == "" {
					return fmt.Errorf("key is blank")
				}

				varMap.Delete(f.Arg(0))
				cmd.Internal.MutVar.Delete(f.Arg(0))
				return nil

			case !*isMut && !*isRef && !*isCh && !*isDel && *isShow:
				sprint := func(m *sync.Map, title string) string {
					res := title
					s := []string{}
					m.Range(func(key, value interface{}) bool {
						k, ok := key.(string)
						if !ok {
							return false
						}

						v, ok := value.(string)
						if !ok {
							return false
						}

						s = append(s, k+" : "+v)
						return true
					})
					return res + "\n" + sortJoin(s)
				}

				fmt.Fprintln(cmd.Stdout, sprint(cmd.Internal.Var, "[variables]"))
				fmt.Fprintln(cmd.Stdout, sprint(cmd.Internal.MutVar, "[mutable variables]"))
				fmt.Fprintln(cmd.Stdout, sprint(cmd.Internal.GlobalVar, "[global variables]"))
				fmt.Fprintln(cmd.Stdout, sprint(cmd.Internal.GlobalMutVar, "[global mutable variables]"))

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

			c := cmd
			c.Internal.Var = new(sync.Map)
			c.Internal.MutVar = new(sync.Map)
			c.Internal.Args = new(sync.Map)
			c.Internal.Return = new(sync.Map)

			for i, v := range argv[1:] {
				c.Internal.Args.Store(i, v)
			}

			if err := EvalString(ctx, c, argv[0]); err != nil {
				return err
			}

			c.Internal.Return.Range(func(key, value interface{}) bool {
				k, ok := key.(string)
				if !ok {
					return false
				}

				v, ok := value.(string)
				if !ok {
					return false
				}

				cmd.Internal.Return.Store(k, v)

				return true
			})

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
