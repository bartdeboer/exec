package exec

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"reflect"

	"github.com/iancoleman/strcase"
)

type Cmd struct {
	cmd *exec.Cmd
}

func NewRunCommand(command string, arg ...string) *Cmd {
	runCmd := exec.Command(command, arg...)
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	runCmd.Stdin = os.Stdin
	return &Cmd{cmd: runCmd}
}

func GetArgsFromStruct(input interface{}) []string {
	agrs := []string{}
	// https://blog.golang.org/laws-of-reflection
	v := reflect.ValueOf(input)
	// s := v.Elem()
	typeOfT := v.Type()

	for i := 0; i < v.NumField(); i++ {
		name := typeOfT.Field(i).Name
		f := v.Field(i)
		strVal := fmt.Sprintf("%v", f.Interface())
		// strVal := f.Interface().(string)
		flag := "--" + strcase.ToKebab(name)
		if strVal == "" {
			agrs = append(agrs, flag)
			continue
		}
		agrs = append(agrs, flag, strVal)
	}
	return agrs
}

func (cmd *Cmd) Run() error {
	err := cmd.cmd.Run()
	if err != nil {
		log.Fatalf("Command.Run() failed with %s\n", err)
	}
	return err
}
