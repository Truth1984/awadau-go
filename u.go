package u

import (
	"fmt"
	"reflect"
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
	return false
}

func Contains(source interface{}, target interface{}) bool {
	return false
}
