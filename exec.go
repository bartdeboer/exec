// Copyright 2009 Bart de Boer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exec

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var ShowCommands bool

var Env map[string]string

type Cmd struct {
	cmd *exec.Cmd
}

func NewCommand(command string, arg ...string) *Cmd {
	return &Cmd{cmd: exec.Command(command, arg...)}
}

func NewIOCommand(command string, arg ...string) *Cmd {
	c := NewCommand(command, arg...)
	c.cmd.Stdin = os.Stdin
	c.cmd.Stdout = os.Stdout
	c.cmd.Stderr = os.Stderr
	return c
}

func (c *Cmd) Start() error {
	if len(Env) > 0 {
		env := os.Environ()
		for k, v := range Env {
			env = append(env, fmt.Sprintf("%s=%s", strings.ToUpper(k), v))
		}
		c.cmd.Env = env
	}
	if ShowCommands {
		fmt.Printf("%+v\n", c.cmd)
		// fmt.Printf("%+v\n", c.cmd.Env)
	}
	return c.cmd.Start()
}

func (c *Cmd) SetEnv(arg ...string) {
	c.cmd.Env = arg
}

func (c *Cmd) Run() error {
	if err := c.Start(); err != nil {
		return err
	}
	return c.cmd.Wait()
}

func (c *Cmd) Wait() (int, error) {
	if err := c.cmd.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode(), err
		}
		return 1, err
	}
	return 0, nil
}
