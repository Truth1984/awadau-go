package u

import (
	"fmt"
)

func Print(args ...interface{}) {
	for i := range args {
		fmt.Printf("%v ", args[i])
	}
	fmt.Println("")
}
