// Copyright 2009 Bart de Boer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exec

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/iancoleman/strcase"
)

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
