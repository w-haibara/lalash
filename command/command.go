package command

import (
	"fmt"
	"io"
	"os"
	"sort"
)

type Env struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

type Command struct {
	Env      Env
	Internal Internal
}

func New() Command {
	cmd := Command{
		Env: Env{
			In:  os.Stdin,
			Out: os.Stdout,
			Err: os.Stderr,
		},
		Internal: NewInternal(),
	}
	cmd.setHelp()
	cmd.setAliasFamily()
	cmd.setVarFamily()
	return cmd
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
