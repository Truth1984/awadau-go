package main

import u "github.com/Truth1984/awadau-go"

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
		"d":      32,
		"s":      "str",
		"year":   2,
		"month":  3,
		"minute": -5,
	}

	v3 := map[string]interface{}{
		"year": 2018,
		"ds":   v2,
	}

	v4 := map[string]interface{}{
		"df": v3,
	}

	u.Print(v)

	vvargs := vvfunc("1", 3, 6, "ds")
	u.Print(u.CP2M(vvargs)[0].(string))

	u.Print(u.Date("2021-04-28 15:50:04.593 +0800 CST"))

	u.Print("DA", u.DateAdd(v2))

	u.Print(u.Date(v3), "maptest")

	u.Print(u.DateFormat("ANSIC"))

	u.Print("values", u.MapGetExist(v2, "d", "f"))

	ps := []interface{}{"df", "ds", "p", "c"}
	ps2 := []interface{}{"dd", "21", 15}
	u.Print("path", u.MapGetPath(v4, ps))

	u.Print("ats", u.ArrayToString(ps2, "-"))

	jstring := "{\"data\":{\"base\":\"BTC\",\"currency\":\"USD\",\"amount\":40000.48}}"
	jmap, _ := u.StringToJson(jstring)

	psjmap := []interface{}{"data", "amounts"}

	u.Print(jmap, u.MapGetPath(jmap, psjmap))

	jsvalue, _ := u.JsonToString(v3)
	u.Print("v3", jsvalue)

	u.Print(u.Types(v2["s"]))

	u.Print("extract", u.ArrayExtract(u.ATI(s), 6))
	// u.Print(u.ATI(s))

	u.Print("mti", u.MTI(v))
}
