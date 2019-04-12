package main

import (
	"os/exec"
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
	cmd := exec.Command(s.Name, s.Args[1:]...)
	cmd.Env = s.Envs
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd
}
