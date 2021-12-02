package main

import u "github.com/Truth1984/awadau-go"

func main() {
	data, e := u.FetchGet("example.com", nil)
	if e != nil {
		panic(e)
	}
	u.Print(data)
}
