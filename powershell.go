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
