package command

import (
	"fmt"
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

func NewInternalCmdMap() Internal {
	in := Internal{
		Cmds:   new(sync.Map),
		Alias:  new(sync.Map),
		MutVar: new(sync.Map),
		Var:    new(sync.Map),
	}

	in.Cmds.Store("help", InternalCmd{
		Usage: "help",
		Fn: func(e Env, args string, argv ...string) error {
			in.Cmds.Range(func(key, value interface{}) bool {
				fmt.Fprintln(e.Out, key, ":", value.(InternalCmd).Usage)
				return true
			})
			return nil
		},
	})

	func() {
		in.Cmds.Store("l-alias", InternalCmd{
			Usage: "alias <alias> <command name>",
			Fn: func(e Env, args string, argv ...string) error {
				in.Alias.Store(argv[0], argv[1])
				return nil
			},
		})

		in.Cmds.Store("l-unalias", InternalCmd{
			Usage: "alias <alias> <command name>",
			Fn: func(e Env, args string, argv ...string) error {
				in.Alias.Delete(argv[0])
				return nil
			},
		})

		in.Cmds.Store("l-alias-show", InternalCmd{
			Usage: "l-alias-show",
			Fn: func(e Env, args string, argv ...string) error {
				in.Alias.Range(func(key, value interface{}) bool {
					fmt.Fprintln(e.Out, key, ":", value)
					return true
				})
				return nil
			},
		})
	}()

	func() {
		in.Cmds.Store("l-var", InternalCmd{
			Usage: "l-var <immutable var name> <value>",
			Fn: func(e Env, args string, argv ...string) error {
				_, ok := in.Var.Load(argv[0])
				if !ok {
					_, ok = in.MutVar.Load(argv[0])
				}
				if ok {
					return fmt.Errorf("variable is already exists:", argv[0])
				}
				in.Var.Store(argv[0], argv[1])
				return nil
			},
		})

		in.Cmds.Store("l-var-mut", InternalCmd{
			Usage: "l-var-mut <mutable var name> <value>",
			Fn: func(e Env, args string, argv ...string) error {
				_, ok := in.Var.Load(argv[0])
				if !ok {
					_, ok = in.MutVar.Load(argv[0])
				}
				if ok {
					return fmt.Errorf("variable is already exists:", argv[0])
				}
				in.MutVar.Store(argv[0], argv[1])
				return nil
			},
		})

		in.Cmds.Store("l-var-ch", InternalCmd{
			Usage: "l-var-ch <mutable var name> <new value>",
			Fn: func(e Env, args string, argv ...string) error {
				if _, ok := in.Var.Load(argv[0]); ok {
					return fmt.Errorf("variable is immutable:", argv[1])
				}
				if _, ok := in.MutVar.Load(argv[0]); !ok {
					return fmt.Errorf("variable is not defined:", argv[1])
				}
				in.MutVar.Store(argv[0], argv[1])
				return nil
			},
		})

		in.Cmds.Store("l-var-ref", InternalCmd{
			Usage: "l-var-ref <var name>",
			Fn: func(e Env, args string, argv ...string) error {
				v, ok := in.Var.Load(argv[0])
				if !ok {
					v, ok = in.MutVar.Load(argv[0])
				}
				if !ok {
					return fmt.Errorf("variable is not defined:", argv[0])
				}
				fmt.Fprintln(e.Out, v)
				return nil
			},
		})

		in.Cmds.Store("l-var-del", InternalCmd{
			Usage: "l-var-del <var name>",
			Fn: func(e Env, args string, argv ...string) error {
				in.Var.Delete(argv[0])
				in.MutVar.Delete(argv[0])
				return nil
			},
		})

		in.Cmds.Store("l-var-show", InternalCmd{
			Usage: "l-var-show",
			Fn: func(e Env, args string, argv ...string) error {
				fmt.Fprintln(e.Out, "[mutable variables]")
				in.MutVar.Range(func(key, value interface{}) bool {
					fmt.Fprintln(e.Out, key, ":", value)
					return true
				})

				fmt.Fprintln(e.Out, "\n[immutable variables]")
				in.Var.Range(func(key, value interface{}) bool {
					fmt.Fprintln(e.Out, key, ":", value)
					return true
				})

				return nil
			},
		})
	}()

	return in
}

func (i Internal) GetAlias(args string) string {
	if v, ok := i.Alias.Load(args); ok {
		return v.(string)
	}
	return args
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
