package main

import (
	u "github.com/Truth1984/awadau-go"
)

type Foo struct {
	Bar []string
	Baz int
}

func vvfunc(args ...interface{}) []interface{} {
	u.Print(u.Types(args), args)
	return args
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

	u.Print(v)

	vvargs := vvfunc("1", 3, 6, "ds")
	u.Print(u.CP2M(vvargs)[0].(string))

	u.Print(u.Date("2021-04-28 15:50:04.593 +0800 CST"))

	u.Print(u.Date(v2), "maptest")

	u.Print(u.DateFormat("ANSIC"))

	u.Print(u.Types(v2["s"]))

}
