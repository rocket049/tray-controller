package main

import (
	"os/exec"
	"reflect"
	"syscall"
)

type myCommand struct {
	Name string
	Args []string
	Envs []string
}

func newMyCommand(name string, args []string, envs []string) *myCommand {
	return &myCommand{Name: name, Args: args, Envs: envs}
}

func (s *myCommand) GetCmd() *exec.Cmd {
	f := reflect.ValueOf(exec.Command)
	argv := []reflect.Value{}
	for _, v := range s.Args {
		argv = append(argv, reflect.ValueOf(v))
	}
	ret := f.Call(argv)
	if len(ret) == 0 {
		return nil
	} else {
		cmd := ret[0].Interface().(*exec.Cmd)
		cmd.Env = s.Envs
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		return cmd
	}
}
