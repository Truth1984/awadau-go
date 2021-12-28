package main

import u "github.com/Truth1984/awadau-go"

type User struct {
	Name string
	Age  int
}

func main() {
	// age := float64(24)
	// amap := u.Map("Name", "Arron", "Age", 12)
	amap := u.Map("Name", "Arron", "Age", 27, "c", 16)

	user := User{}
	u.MapToStructHandled(amap, &user)

	val, err := u.StructToMap(&user)
	u.Print(err, val)
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
