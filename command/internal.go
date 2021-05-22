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
	CmdMap sync.Map
	MutVar sync.Map
	Var    sync.Map
}

func NewInternalCmdMap() Internal {
	var m sync.Map

	m.Store("help", InternalCmd{
		Usage: "help",
		Fn: func(e Env, args string, argv ...string) error {
			m.Range(func(key, value interface{}) bool {
				fmt.Fprintln(e.Out, key, ":", value.(InternalCmd).Usage)
				return true
			})
			return nil
		},
	})

	var mutVarMap sync.Map
	var varMap sync.Map

	func() {
		m.Store("l-var", InternalCmd{
			Usage: "l-var <immutable var name> <value>",
			Fn: func(e Env, args string, argv ...string) error {
				_, ok := varMap.Load(argv[0])
				if !ok {
					_, ok = mutVarMap.Load(argv[0])
				}
				if ok {
					return fmt.Errorf("variable is already exists:", argv[0])
				}
				varMap.Store(argv[0], argv[1])
				return nil
			},
		})

		m.Store("l-var-mut", InternalCmd{
			Usage: "l-var-mut <mutable var name> <value>",
			Fn: func(e Env, args string, argv ...string) error {
				_, ok := varMap.Load(argv[0])
				if !ok {
					_, ok = mutVarMap.Load(argv[0])
				}
				if ok {
					return fmt.Errorf("variable is already exists:", argv[0])
				}
				mutVarMap.Store(argv[0], argv[1])
				return nil
			},
		})

		m.Store("l-var-ch", InternalCmd{
			Usage: "l-var-ch <mutable var name> <new value>",
			Fn: func(e Env, args string, argv ...string) error {
				if _, ok := varMap.Load(argv[0]); ok {
					return fmt.Errorf("variable is immutable:", argv[1])
				}
				if _, ok := mutVarMap.Load(argv[0]); !ok {
					return fmt.Errorf("variable is not defined:", argv[1])
				}
				mutVarMap.Store(argv[0], argv[1])
				return nil
			},
		})

		m.Store("l-var-ref", InternalCmd{
			Usage: "l-var-ref <var name>",
			Fn: func(e Env, args string, argv ...string) error {
				v, ok := varMap.Load(argv[0])
				if !ok {
					v, ok = mutVarMap.Load(argv[0])
				}
				if !ok {
					return fmt.Errorf("variable is not defined:", argv[0])
				}
				fmt.Fprintln(e.Out, v)
				return nil
			},
		})

		m.Store("l-var-del", InternalCmd{
			Usage: "l-var-del <var name>",
			Fn: func(e Env, args string, argv ...string) error {
				varMap.Delete(argv[0])
				mutVarMap.Delete(argv[0])
				return nil
			},
		})

		m.Store("l-var-show", InternalCmd{
			Usage: "l-var-show",
			Fn: func(e Env, args string, argv ...string) error {
				fmt.Fprintln(e.Out, "[mutable variables]")
				mutVarMap.Range(func(key, value interface{}) bool {
					fmt.Fprintln(e.Out, key, ":", value)
					return true
				})

				fmt.Fprintln(e.Out, "\n[immutable variables]")
				varMap.Range(func(key, value interface{}) bool {
					fmt.Fprintln(e.Out, key, ":", value)
					return true
				})

				return nil
			},
		})
	}()

	return Internal{
		CmdMap: m,
		MutVar: mutVarMap,
		Var:    varMap,
	}
}

func (i Internal) Get(key string) (InternalCmd, error) {
	v, ok := i.CmdMap.Load(key)
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
