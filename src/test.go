package main

import (
	"fmt"
	"reflect"
	"time"
)

var EMPTY struct{}

func print(args ...interface{}) {
	for i := range args {
		fmt.Printf("%v ", args[i])
	}
	fmt.Println("")
}

func CP2M(args interface{}) map[int]interface{} {
	arr := reflect.ValueOf(args)
	aMap := make(map[int]interface{})
	for i := 0; i < arr.Len(); i++ {
		aMap[i] = arr.Index(i)
	}
	return aMap
}

func main() {
	// s := [8]int{1, 2, 3, 4, 5}
	p := [7]interface{}{1, EMPTY, EMPTY, 8}
	aMap := map[int]interface{}{1: "12", 3: "18"}

	var tm int64 = 1619587539715
	var sec = (tm / 1000)
	var msec = (tm % 1000)
	var t = time.Unix(sec, msec*int64(time.Millisecond))

	print(t.Format("Mon Jan _2 15:04:05 2006"))

	print(aMap[5] == nil)
	print(CP2M(p)[2])
}
