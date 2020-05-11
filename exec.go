package exec

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
)

type Cmd struct {
	cmd *exec.Cmd
}

func (cmd *Cmd) List() ([]string, error) {
	list := []string{}
	stdout, err := cmd.cmd.StdoutPipe()
	if err != nil {
		return list, err
	}
	stderr, err := cmd.cmd.StderrPipe()
	if err != nil {
		return list, err
	}
	scanner := bufio.NewScanner(stdout)
	if err := cmd.cmd.Start(); err != nil {
		return list, err
	}
	for scanner.Scan() {
		list = append(list, strings.Trim(scanner.Text(), " "))
	}
	errStr, _ := ioutil.ReadAll(stderr)
	if err := cmd.cmd.Wait(); err != nil {
		fmt.Printf("%s\n", strings.Trim(string(errStr), "\n"))
		// Get the exit code
		// if exitError, ok := err.(*exec.ExitError); ok {
		// 	fmt.Printf("ExitCode: %v\n", exitError.ExitCode())
		// }
		return list, err
	}
	return list, scanner.Err()
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

func BuildArgs(args ...interface{}) []string {
	ret := []string{}
	for _, arg := range args {
		rv := reflect.ValueOf(arg)
		rt := reflect.TypeOf(arg)
		switch rv.Kind() {
		case reflect.Slice, reflect.Array:
			if rt.Elem().Kind() == reflect.String {
				ret = append(ret, arg.([]string)...)
			} else {
				ret = append(ret, SliceToString(arg)...)
			}
		case reflect.String:
			ret = append(ret, arg.(string))
		case reflect.Int:
			ret = append(ret, strconv.FormatInt(rv.Int(), 10))
		case reflect.Struct:
			ret = append(ret, BuildArgsFromStruct(arg)...)
		default:
			ret = append(ret, fmt.Sprintf("%v", arg))
		}
	}
	return ret
}

func BuildArgsFromStruct(input interface{}) []string {
	agrs := []string{}
	// https://blog.golang.org/laws-of-reflection
	rv := reflect.ValueOf(input)
	typeOfT := rv.Type()
	if k := rv.Kind(); k != reflect.Struct {
		panic("Value is not a struct")
	}
	for i := 0; i < rv.NumField(); i++ {
		name := typeOfT.Field(i).Name
		value := rv.Field(i)
		opt := "--" + strcase.ToKebab(name)
		if value.Kind() == reflect.Bool {
			if value.Bool() == true {
				agrs = append(agrs, opt)
			}
			continue
		}
		strVal := ValueToString(rv.Field(i))
		// strVal := fmt.Sprintf("%v", f.Interface())
		// strVal := f.Interface().(string)
		if strVal == "" || strVal == "0" {
			continue
		}
		agrs = append(agrs, opt, strVal)
	}
	return agrs
}

func ValueToString(val reflect.Value) string {
	var ret string
	switch val.Kind() {
	case reflect.String:
		ret = val.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ret = strconv.FormatInt(val.Int(), 10)
	}
	return ret
}

func SliceToString(input interface{}) []string {
	var ret []string
	v := reflect.ValueOf(input)
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		ret = make([]string, v.Len())
		for i := 0; i < v.Len(); i++ {
			ret[i] = ValueToString(v.Index(i))
			// ret[i] = fmt.Sprintf("%v", v.Index(i).Interface())
		}
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
