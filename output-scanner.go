// Copyright 2009 Bart de Boer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exec

import (
	"bufio"
	"os"
	"strings"
)

type OutputScanner struct {
	cmd     *Cmd
	scanner *bufio.Scanner
}

func NewOutputScanner(c *Cmd) *OutputScanner {
	return &OutputScanner{
		cmd: c,
	}
}

func (o *OutputScanner) Start() error {
	stdout, err := o.cmd.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	// sr.stdoutPipe = stdout
	o.cmd.cmd.Stderr = os.Stderr
	// stderr, err := sr.cmd.cmd.StderrPipe()
	// if err != nil {
	// 	return err
	// }
	// sr.stderrPipe = stderr
	if err := o.cmd.Start(); err != nil {
		return err
	}
	o.scanner = bufio.NewScanner(stdout)
	return nil
}

func (o *OutputScanner) Scan() bool {
	return o.scanner.Scan()
}

func (o *OutputScanner) Text() string {
	return o.scanner.Text()
}

// func (sr *OutputScanner) ErrorOutput() ([]byte, error) {
// 	return ioutil.ReadAll(sr.stderrPipe)
// }

func (o *OutputScanner) Wait() (int, error) {
	return o.cmd.Wait()
}

func (o *OutputScanner) Lines() ([]string, error) {
	lines := []string{}
	if o.cmd.cmd.Process == nil {
		err := o.Start()
		if err != nil {
			return nil, err
		}
	}
	for o.Scan() {
		lines = append(lines, o.Text())
	}
	// stderrOut, _ := sr.ErrorOutput()
	// sr.stderrOut = stderrOut
	code, err := o.Wait()
	if code > 0 {
		return lines, err
	}
	return lines, nil
}

func (o *OutputScanner) HasLine(line string) (bool, error) {
	if o.cmd.cmd.Process == nil {
		err := o.Start()
		if err == nil {
			return false, err
		}
	}
	for o.Scan() {
		if strings.Trim(o.Text(), " ") == line {
			return true, nil
		}
	}
	if code, err := o.Wait(); code > 0 {
		return false, err
	}
	return false, nil
}
