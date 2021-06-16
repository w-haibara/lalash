package lalash

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"sort"
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

func (cmd Command) setUtilFamily() {
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
			fmt.Fprintln(cmd.Stdout, argv)
			return nil
		},
	})

	cmd.Internal.Cmds.Store("l-exit", InternalCmd{
		Usage: "l-exit",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			os.Exit(0)
			return nil
		},
	})
}

func (cmd Command) setAliasFamily() {
	cmd.Internal.Cmds.Store("l-alias", InternalCmd{
		Usage: "",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {

			f := flag.NewFlagSet("alias", flag.ContinueOnError)
			isUnset := f.Bool("unset", false, "")
			isShow := f.Bool("show", false, "")
			key := f.String("k", "", "")
			val := f.String("v", "", "")
			if err := f.Parse(argv); err != nil {
				return err
			}

			if *isUnset && *isShow {
				return fmt.Errorf("cannot set both --unset and --show.")
			}

			if !*isUnset && !*isShow {
				if *key == "" {
					return fmt.Errorf("key is blank")
				}
				if *val == "" {
					return fmt.Errorf("value is blank")
				}
				cmd.Internal.Alias.Store(*key, *val)
				return nil
			}

			if *isUnset {
				if *key == "" {
					return fmt.Errorf("key is blank")
				}
				cmd.Internal.Alias.Delete(*key)
				return nil
			}

			if *isShow {
				cmd.Internal.Alias.Range(func(key, value interface{}) bool {
					fmt.Fprintln(cmd.Stdout, key, ":", value)
					return true
				})
				return nil
			}

			return nil
		},
	})

}

func (cmd Command) setVarFamily() {
	cmd.Internal.SetInternalCmd("l-var", InternalCmd{
		Usage: "l-var <immutable var name> <value>",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 2); err != nil {
				return err
			}
			_, ok := cmd.Internal.Var.Load(argv[0])
			if !ok {
				_, ok = cmd.Internal.MutVar.Load(argv[0])
			}
			if ok {
				return fmt.Errorf("variable is already exists: %v", argv[0])
			}
			cmd.Internal.Var.Store(argv[0], argv[1])
			return nil
		},
	})

	cmd.Internal.SetInternalCmd("l-var-mut", InternalCmd{
		Usage: "l-var-mut <mutable var name> <value>",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 2); err != nil {
				return err
			}
			_, ok := cmd.Internal.Var.Load(argv[0])
			if !ok {
				_, ok = cmd.Internal.MutVar.Load(argv[0])
			}
			if ok {
				return fmt.Errorf("variable is already exists: %v", argv[0])
			}
			cmd.Internal.MutVar.Store(argv[0], argv[1])
			return nil
		},
	})

	cmd.Internal.SetInternalCmd("l-var-ch", InternalCmd{
		Usage: "l-var-ch <mutable var name> <new value>",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 2); err != nil {
				return err
			}
			if _, ok := cmd.Internal.Var.Load(argv[0]); ok {
				return fmt.Errorf("variable is immutable: %v", argv[1])
			}
			if _, ok := cmd.Internal.MutVar.Load(argv[0]); !ok {
				return fmt.Errorf("variable is not defined: %v", argv[1])
			}
			cmd.Internal.MutVar.Store(argv[0], argv[1])
			return nil
		},
	})

	cmd.Internal.SetInternalCmd("l-var-ch", InternalCmd{
		Usage: "l-var-ch <mutable var name> <new value>",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 2); err != nil {
				return err
			}
			if _, ok := cmd.Internal.Var.Load(argv[0]); ok {
				return fmt.Errorf("variable is immutable: %v", argv[1])
			}
			if _, ok := cmd.Internal.MutVar.Load(argv[0]); !ok {
				return fmt.Errorf("variable is not defined: %v", argv[1])
			}
			cmd.Internal.MutVar.Store(argv[0], argv[1])
			return nil
		},
	})

	cmd.Internal.SetInternalCmd("l-var-ref", InternalCmd{
		Usage: "l-var-ref <var name>",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 1); err != nil {
				return err
			}
			v, ok := cmd.Internal.Var.Load(argv[0])
			if !ok {
				v, ok = cmd.Internal.MutVar.Load(argv[0])
			}
			if !ok {
				return fmt.Errorf("variable is not defined: %v", argv[0])
			}
			fmt.Fprintln(cmd.Stdout, v)
			return nil
		},
	})

	cmd.Internal.SetInternalCmd("l-var-del", InternalCmd{
		Usage: "l-var-del <var name>",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 1); err != nil {
				return err
			}
			cmd.Internal.Var.Delete(argv[0])
			cmd.Internal.MutVar.Delete(argv[0])
			return nil
		},
	})

	cmd.Internal.SetInternalCmd("l-var-show", InternalCmd{
		Usage: "l-var-show",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			fmt.Fprintln(cmd.Stdout, "[mutable variables]")
			cmd.Internal.MutVar.Range(func(key, value interface{}) bool {
				fmt.Fprintln(cmd.Stdout, key, ":", value)
				return true
			})

			fmt.Fprintln(cmd.Stdout, "\n[immutable variables]")
			cmd.Internal.Var.Range(func(key, value interface{}) bool {
				fmt.Fprintln(cmd.Stdout, key, ":", value)
				return true
			})

			return nil
		},
	})
}

func (cmd Command) setEvalFamily() {
	cmd.Internal.Cmds.Store("l-eval", InternalCmd{
		Usage: "l-eval",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 1); err != nil {
				return err
			}
			if err := EvalString(ctx, cmd, argv[0]); err != nil {
				return fmt.Errorf("eval error [%v] : %v", argv[0], err.Error())
			}
			return nil
		},
	})

	cmd.Internal.Cmds.Store("l-pipe", InternalCmd{
		Usage: "l-pipe",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 2); err != nil {
				return err
			}

			var pipe bytes.Buffer

			cmd1 := cmd
			o := bufio.NewWriter(&pipe)
			cmd1.Stdout = o
			if err := EvalString(ctx, cmd1, argv[0]); err != nil {
				return fmt.Errorf("eval error [%v] : %v", argv[0], err.Error())
			}
			o.Flush()

			cmd2 := cmd
			cmd2.Stdin = strings.NewReader(pipe.String())
			if err := EvalString(ctx, cmd2, argv[1]); err != nil {
				return fmt.Errorf("eval error [%v] : %v", argv[1], err.Error())
			}

			return nil
		},
	})
}
