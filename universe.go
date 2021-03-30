package main

import (
	"fmt"
)

func Print(args ...interface{}) []string {
	for i := range args {
		fmt.Printf("%v ", args[i])
	}
	fmt.Println("")
	var a []string
	return a
}
