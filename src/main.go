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

func t2() {
	u.Print(u.StringToJson("{}"))
}

func t3() {
	u.Print(u.MapMerge(u.Map(), nil))
}

func t4() {
	u.Print(u.ToInt(1.1), u.ToInt(2), u.ToInt("3"), u.ToInt(u.Map("n", 4)["n"]))
}

func t5() {
	amap := u.Map("a", nil)
	u.Print(amap["a"] == nil)
}

func t6() {
	u.Print(u.JsonToString(u.Map("a", 23), ""))
}
