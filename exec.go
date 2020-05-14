package exec

import (
	"os"
	"os/exec"
)

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
	if err := c.cmd.Start(); err != nil {
		return err
	}
	return nil
}

func (c *Cmd) Run() error {
	if err := c.cmd.Run(); err != nil {
		return err
	}
	return nil
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
