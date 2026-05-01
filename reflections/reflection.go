package main

import (
	"reflect"
)

func walk(x any, fn func(input string)) {
	val := getValue(x)

	switch val.Kind() {
	case reflect.String:
		fn(val.String())

	case reflect.Struct:
		for _, field := range val.Fields() {
			walk(field.Interface(), fn)
		}

	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			walk(val.Index(i).Interface(), fn)
		}

	case reflect.Map:
		iter := val.MapRange()
		for iter.Next() {
			walk(iter.Value().Interface(), fn)
		}
	case reflect.Chan:
		for recv, ok := val.Recv(); ok; recv, ok = val.Recv() {
			walk(recv.Interface(), fn)
		}
	case reflect.Func:
		valFnResult := val.Call(nil)
		for _, res := range valFnResult {
			walk(res.Interface(), fn)
		}
	}
}

func getValue(x any) reflect.Value {
	val := reflect.ValueOf(x)

	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	return val
}
