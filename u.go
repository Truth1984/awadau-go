package u

import (
	"fmt"
	"reflect"
	"strings"
)

func Print(args ...interface{}) {
	for i := range args {
		fmt.Printf("%v ", args[i])
	}
	fmt.Println("")
}

func Types(source interface{}) string {
	return reflect.TypeOf(source).Kind().String()
}

// @param { "num" } expect
func TypesCheck(source interface{}, expect string) bool {
	stype := Types(source)

	switch strings.ToLower(expect) {
	case "str":
	case "string":
		return stype == "string"
	case "num":
	case "number":
		return strings.Contains(stype, "int") || strings.Contains(stype, "float")
	case "int":
		return stype == "int"
	case "float":
		return strings.Contains(stype, "float")
	case "arr":
	case "array":
		return stype == "array"
	case "map":
	case "dict":
		return stype == "map"
	default:
		return strings.Contains(stype, expect)
	}
	return false
}

func Contains(source interface{}, target interface{}) bool {
	return false
}

func Ps(structPtr *struct{}) func() {
	return func() {
		fmt.Println(reflect.TypeOf(structPtr).String())
	}
}
