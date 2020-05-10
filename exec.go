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

func MakeArgs(args ...interface{}) []string {
	ret := []string{}
	for _, arg := range args {
		v := reflect.ValueOf(arg)
		// s := v.Elem()
		// typeOfT := v.Type()

		rt := reflect.TypeOf(arg)
		fmt.Print(rt.Kind())
		fmt.Print("\n")
		fmt.Print(v.Kind())
		fmt.Print("\n")
		// fmt.Print(v.Elem())
		// fmt.Print("\n")
		switch v.Kind() {
		case reflect.Slice, reflect.Array:
			fmt.Print(rt.Elem())
			fmt.Print("\n")
			fmt.Print(rt.Elem().Kind())
			fmt.Print("\n")
			if rt.Elem().Kind() == reflect.String {
				ret = append(ret, arg.([]string)...)
			}
		case reflect.String:
			ret = append(ret, arg.(string))
		case reflect.Struct:
			fmt.Print(arg)
			fmt.Print("\n")
			ret = append(ret, GetArgsFromStruct(arg)...)
		default:
			// fmt.Println(k, "is something else entirely")
		}
	}
	return ret
}

func AppendStructArgs(args []string, input interface{}) []string {
	return append(args, GetArgsFromStruct(input)...)
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

func GetArgsFromSlice(input []interface{}) []string {
	ret := make([]string, len(input))
	for i, v := range input {
		ret[i] = v.(string)
	}
	return ret
}

func (cmd *Cmd) Run() error {
	err := cmd.cmd.Run()
	if err != nil {
		log.Fatalf("Command.Run() failed with %s\n", err)
	}
	return err
}
