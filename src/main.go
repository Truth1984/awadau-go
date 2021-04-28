package main

import (
	"time"

	u "github.com/Truth1984/awadau-go"
)

type Foo struct {
	Bar []string
	Baz int
}

func main() {
	s := [8]int{1, 2, 3, 4, 5}
	u.Print(u.Types(s))

	v := map[string]int{
		"d": 32,
		"r": 55,
	}

	v2 := map[string]interface{}{
		"d": 32,
		"s": "str",
	}

	time := time.Now()

	u.Print(v)

	u.Print(u.TypesCheck(s, "string"))

	u.Print(u.Types(time))

	u.Print(len(s))

	u.Print(u.Types(v2["s"]))

}
