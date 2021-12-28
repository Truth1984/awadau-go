package main

import u "github.com/Truth1984/awadau-go"

type User struct {
	Name string
	Age  int
}

func main() {
	amap := u.Map("Name", "Arron")
	user := User{}
	u.MapToStruct(amap, &user)
	u.Print(user)
}
