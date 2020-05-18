// Copyright 2009 Bart de Boer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exec

import (
	"fmt"
	"strings"
)

func NewPowerShellCommand(command string, arg ...string) *Cmd {
	return NewCommand("powershell", "-Command",
		fmt.Sprintf("& {%s %s; If (!$?) { exit 1 }}", command, strings.Join(arg, " ")),
	)
}

func NewPowerShellIOCommand(command string, arg ...string) *Cmd {
	return NewIOCommand("powershell", "-Command",
		fmt.Sprintf("& {%s %s; If (!$?) { exit 1 }}", command, strings.Join(arg, " ")),
	)
}
