package commands

import (
	"fmt"
	"reflect"
	"sync"

	"lalash/env"
)

type InternalCmd struct {
	Usage string
	Fn    func(env.Env, string, ...string) error
}

type InternalCmdMap sync.Map

func New() InternalCmdMap {
	var m sync.Map

	m.Store("echo1", InternalCmd{
		Usage: "this is echo1 command",
		Fn: func(e env.Env, args string, argv ...string) error {
			fmt.Println("1", argv)
			return nil
		},
	})

	m.Store("echo2", InternalCmd{
		Usage: "this is echo2 command",
		Fn: func(e env.Env, args string, argv ...string) error {
			fmt.Println("2", argv)
			return nil
		},
	})

	m.Store("help", InternalCmd{
		Usage: "help",
		Fn: func(e env.Env, args string, argv ...string) error {
			m.Range(func(key ,value interface{}) bool {
				fmt.Fprintln(e.Out, key, ":", value.(InternalCmd).Usage)
				return true
			})
			return nil
		},
	})

	return InternalCmdMap(m)
}

func (m InternalCmdMap) Get(key string) (InternalCmd, error) {
	cmds := sync.Map(m)
	v, ok := cmds.Load(key)
	if !ok {
		return InternalCmd{}, fmt.Errorf("command not found")
	}
	cmd, ok := v.(InternalCmd)
	if !ok {
		return InternalCmd{}, fmt.Errorf("function of the command is invalid")
	}
	return InternalCmd(cmd), nil
}

func (cmd InternalCmd) Exec(env env.Env, args string, argv ...string) error {
	return cmd.Fn(env, args, argv...)
}

func checkMapType(m sync.Map, key string) bool {
	v, ok := m.Load(key)
	if !ok {
		return false
	}
	return reflect.TypeOf(v) == reflect.TypeOf(InternalCmd{})
}
