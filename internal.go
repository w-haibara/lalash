package lalash

import (
	"fmt"
	"os"
	"sort"
	"sync"
)

type InternalCmd struct {
	Usage string
	Fn    func(Env, string, ...string) error
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

func (cmd InternalCmd) Exec(env Env, args string, argv ...string) error {
	return cmd.Fn(env, args, argv...)
}

func (cmd Command) setHelp() {
	cmd.Internal.Cmds.Store("help", InternalCmd{
		Usage: "help",
		Fn: func(e Env, args string, argv ...string) error {
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
			fmt.Fprintln(e.Out, ret)
			return nil
		},
	})
}

func (cmd Command) setExit() {
	cmd.Internal.Cmds.Store("exit", InternalCmd{
		Usage: "exit",
		Fn: func(e Env, args string, argv ...string) error {
			os.Exit(0)
			return nil
		},
	})
}

func (cmd Command) setAliasFamily() {
	cmd.Internal.Cmds.Store("l-alias", InternalCmd{
		Usage: "alias <alias> <command name>",
		Fn: func(e Env, args string, argv ...string) error {
			if err := checkArgv(argv, 2); err != nil {
				return err
			}
			cmd.Internal.Alias.Store(argv[0], argv[1])
			return nil
		},
	})

	cmd.Internal.Cmds.Store("l-unalias", InternalCmd{
		Usage: "alias <alias> <command name>",
		Fn: func(e Env, args string, argv ...string) error {
			if err := checkArgv(argv, 1); err != nil {
				return err
			}
			cmd.Internal.Alias.Delete(argv[0])
			return nil
		},
	})

	cmd.Internal.Cmds.Store("l-alias-show", InternalCmd{
		Usage: "l-alias-show",
		Fn: func(e Env, args string, argv ...string) error {
			cmd.Internal.Alias.Range(func(key, value interface{}) bool {
				fmt.Fprintln(e.Out, key, ":", value)
				return true
			})
			return nil
		},
	})
}

func (cmd Command) setVarFamily() {
	cmd.Internal.SetInternalCmd("l-var", InternalCmd{
		Usage: "l-var <immutable var name> <value>",
		Fn: func(e Env, args string, argv ...string) error {
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
		Fn: func(e Env, args string, argv ...string) error {
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
		Fn: func(e Env, args string, argv ...string) error {
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
		Fn: func(e Env, args string, argv ...string) error {
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
		Fn: func(e Env, args string, argv ...string) error {
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
			fmt.Fprintln(e.Out, v)
			return nil
		},
	})

	cmd.Internal.SetInternalCmd("l-var-del", InternalCmd{
		Usage: "l-var-del <var name>",
		Fn: func(e Env, args string, argv ...string) error {
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
		Fn: func(e Env, args string, argv ...string) error {
			fmt.Fprintln(e.Out, "[mutable variables]")
			cmd.Internal.MutVar.Range(func(key, value interface{}) bool {
				fmt.Fprintln(e.Out, key, ":", value)
				return true
			})

			fmt.Fprintln(e.Out, "\n[immutable variables]")
			cmd.Internal.Var.Range(func(key, value interface{}) bool {
				fmt.Fprintln(e.Out, key, ":", value)
				return true
			})

			return nil
		},
	})
}
