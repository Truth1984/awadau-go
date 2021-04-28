package u

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

var EMPTY struct{}

func Print(args ...interface{}) {
	for i := range args {
		fmt.Printf("%v ", args[i])
	}
	fmt.Println("")
}

func CP2M(args interface{}) map[int]interface{} {
	arr := reflect.ValueOf(args)
	aMap := make(map[int]interface{})
	for i := 0; i < arr.Len(); i++ {
		aMap[i] = arr.Index(i)
	}
	return aMap
}

func Ternary(boolean bool, trueValue interface{}, falseValue interface{}) interface{} {
	if boolean {
		return trueValue
	}
	return falseValue
}

func Types(source interface{}) string {
	types := reflect.TypeOf(source).Kind().String()
	if types == "struct" {
		return reflect.TypeOf(source).Name()
	}
	return types
}

// @param { "num" } expect
func TypesCheck(source interface{}, expect string) bool {
	stype := Types(source)

	switch strings.ToLower(expect) {
	case "str":
	case "string":
		return stype == "string"
	case "num":
	case "number":
		return strings.Contains(stype, "int") || strings.Contains(stype, "float")
	case "int":
		return strings.Contains(stype, "int")
	case "float":
		return strings.Contains(stype, "float")
	case "arr":
	case "array":
		return stype == "array"
	case "map":
	case "dict":
		return stype == "map"
	case "date":
	case "time":
		return stype == "Time"
	}
	return strings.Contains(stype, expect)
}

func Contains(source interface{}, target interface{}) bool {
	return false
}

func dateLayout(str string) string {
	switch str {
		case "iso":
			return "2020-04-09T06:05:45.290Z"
		case "ANSIC":
			return "Mon Jan _2 15:04:05 2006" 
		case "UnixDate":
			return "Mon Jan _2 15:04:05 MST 2006"
		case "RubyDate":
			return "Mon Jan 02 15:04:05 -0700 2006"
		case "RFC822":
			return "02 Jan 06 15:04 MST"
		case "RFC822Z":
			return "02 Jan 06 15:04 -0700"
		case "RFC850":
			return "Monday, 02-Jan-06 15:04:05 MST"
		case "RFC1123":
			return "Mon, 02 Jan 2006 15:04:05 MST"
		case "RFC1123Z":
			return "Mon, 02 Jan 2006 15:04:05 -0700"
		case "RFC3339":
			return "2006-01-02T15:04:05Z07:00"
		case "RFC3339Nano":
			return "2006-01-02T15:04:05.999999999Z07:00"
		case "Kitchen":
			return "3:04PM"
	}
	return str
}


func dateStrParse(str string) time.Time {

	// t, err := time.Parse(layout, str)

	// * "date":"Thu Apr 09 2020",
	// *
	// * "iso":"2020-04-09T06:05:45.290Z",
	// *
	// * "json":{"year":2020,"month":4,"day":9,"hour":14,"minute":5,"second":45},
	// *
	// * "localedate":"4/9/2020",
	// *
	// * "localetime":"2:05:45 PM",
	// *
	// * "locale":"4/9/2020, 2:05:45 PM",
	// *
	// * "locale24":"4/9/2020, 14:05:45",
	// *
	// * "datetime":"2020-04-09 06:05:45",
	// *
	// * "datetime0":"2020-04-08 16:00:00",
	// *
	// * "string":"Thu Apr 09 2020 14:05:45 GMT+0800 (China Standard Time)",
	// *
	// * "time":"14:05:45 GMT+0800 (China Standard Time)",
	// *
	// * "plain":"2020_4_9_14_5_45",
	// *
	// * "long":1586412345290}
	// utc := "Thu, 09 Apr 2020 06:05:45 GMT"
	return time.Now()
}

// string | ["year", "month", "day", "hour", "minute" , "second"] | {} | number | time | ""
func Date(input interface{}) time.Time {
	if TypesCheck(input, "") {
		return time.Now()
	}

	if TypesCheck(input, "time") {
		return input.(time.Time)
	}

	if TypesCheck(input, "num") {
		tn := input.(int64)
   		return time.Unix(tn / 1000, tn % 1000 *int64(time.Millisecond))
	}

	if TypesCheck(input, "array") {
		tm := CP2M(input)
		year := Ternary(tm[0] == nil, time.Now().Year(), tm[0]).(int)
		month:= Ternary(tm[1] == nil, time.Now(), tm[1]).(time.Month)
		day:= Ternary(tm[2] == nil, time.Now().Day(), tm[2]).(int)
		hour:= Ternary(tm[3] == nil, time.Now().Hour(), tm[3]).(int)
		minute:= Ternary(tm[4] == nil, time.Now().Minute(), tm[4]).(int)
		second:= Ternary(tm[5] == nil, time.Now().Second(), tm[5]).(int)
		nanosecond := Ternary(tm[6] == nil, time.Now().Nanosecond(), tm[6]).(int)
		loc := Ternary(tm[7] == nil, time.Now().Location(), tm[7]).(time.Location)
		return time.Date(year,month,day,hour,minute,second,nanosecond,&loc)
	}

	if TypesCheck(input, "map") {
		tm := input.(map[string] interface{})
		year := Ternary(tm["year"] == nil, time.Now().Year(), tm["year"]).(int)
		month:= Ternary(tm["month"] == nil, time.Now(), tm["month"]).(time.Month)
		day:= Ternary(tm["day"] == nil, time.Now().Day(), tm["day"]).(int)
		hour:= Ternary(tm["hour"] == nil, time.Now().Hour(), tm["hour"]).(int)
		minute:= Ternary(tm["minute"] == nil, time.Now().Minute(), tm["minute"]).(int)
		second:= Ternary(tm["second"] == nil, time.Now().Second(), tm["second"]).(int)
		nanosecond := Ternary(tm["nanosecond"] == nil, time.Now().Nanosecond(), tm["nanosecond"]).(int)
		loc := Ternary(tm["loc"] == nil, time.Now().Location(), tm["loc"]).(time.Location)
		return time.Date(year,month,day,hour,minute,second,nanosecond,&loc)
	}
	
	if TypesCheck(input, "str") {
		return dateStrParse(input.(string));
	}

	panic("unable to convert to Date")
}

// ANSIC       = "Mon Jan _2 15:04:05 2006" 
//
// UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
//
// RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
//
// RFC822      = "02 Jan 06 15:04 MST"
//
// RFC822Z     = "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
//
// RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
//
// RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
//
// RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
//
// RFC3339     = "2006-01-02T15:04:05Z07:00"
//
// RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
//
// Kitchen     = "3:04PM"
func DateFormat(formatThenDate ...interface{}) string {
	args := CP2M(formatThenDate)
	format := dateLayout(args[0].(string))
	dates := Ternary(args[1] == nil, time.Now(), Date(args[1])).(time.Time)

	return dates.Format(format)
}
