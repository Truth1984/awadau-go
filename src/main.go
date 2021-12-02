package main

import u "github.com/Truth1984/awadau-go"

func main() {
	t3()
}

func t1() {
	data, e := u.FetchGet("example.com", nil)
	if e != nil {
		panic(e)
	}
	u.Print(data)
}

func t2() {
	u.Print(u.StringToJson("{}"))
}

func t3() {
	u.Print(u.MapMerge(u.Map(), nil))
}
