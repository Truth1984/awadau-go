package u

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/araddon/dateparse"
)

var EMPTY struct{}

func Print(args ...interface{}) {
	for i := range args {
		fmt.Printf("%v ", args[i])
	}
	fmt.Println("")
}

func ArrayToMap(array []interface{}) map[int]interface{} {
	aMap := make(map[int]interface{})
	for i := 0; i < len(array); i++ {
		aMap[i] = array[i]
	}
	return aMap
}

func ArrayPopRight(array *[]interface{}) interface{} {
	if(len(*array) == 0) {
		return nil
	}
	pop := (*array)[len(*array)-1]
	(*array) = (*array)[:len(*array)-1]
	return pop
}

func ArrayPopLeft(array *[]interface{}) interface{} {
	if(len(*array) == 0) {
		return nil
	}
	shift := (*array)[0]
	(*array) = (*array)[1:]
	return shift
}

func ToString(item interface{}) string {
	return fmt.Sprint(item)
}

// sep = ","
func ArrayToString(arrayThenSep ...interface{}) string {
	ats := CP2M(arrayThenSep)
	array := arrayThenSep[0].([]interface{})
	sep := ","
	if (ats[1] != nil) {
		sep = ats[1].(string)
	}

	result := ""
	for i := 0; i < len(array); i++ { 
		result += ToString(array[i]) + sep
	}
	return strings.TrimSuffix(result,sep)
}

func CP2M(array []interface{}) map[int]interface{} {
	arr := reflect.ValueOf(array)
	aMap := map[int]interface{}{}
	for i := 0; i < arr.Len(); i++ {
		aMap[i] = arr.Index(i).Interface()
	}
	return aMap
}

// trueValue or falseValue must exist / not panic
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
		case "str","string":
			return stype == "string"
		case "num","number":
			return strings.Contains(stype, "int") || strings.Contains(stype, "float")
		case "int":
			return strings.Contains(stype, "int")
		case "float":
			return strings.Contains(stype, "float")
		case "arr","array":
			return stype == "array"
		case "map","dict":
			return stype == "map"
		case "date","time":
			return stype == "Time"
		}

	return strings.Contains(stype, expect)
}

func dateLayout(str string) string {
	switch str {
		case "plain":
			return "2020-04-09-06-05-45"
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
	t,e := dateparse.ParseLocal(str)
	if e != nil {
		panic(e)
	}
	return t
}

// string | ["year", "month", "day", "hour", "minute" , "second"] | {} | number | time | ""
func Date(input interface{}) time.Time {
	if TypesCheck(input, "string") && input.(string) == "" {
		return time.Now()
	}

	if TypesCheck(input, "time") {
		return input.(time.Time)
	}

	if TypesCheck(input, "int") {
		tn := int64(input.(int))
   		return time.Unix(tn / 1000, tn % 1000 *int64(time.Millisecond))
	}

	if TypesCheck(input, "array") {
		tm := ArrayToMap(input.([]interface{}))
		year := Ternary(tm[0] == nil, time.Now().Year(), tm[0]).(int)
		month:= Ternary(tm[1] == nil, time.Now(), tm[1]).(time.Time).Month()
		day:= Ternary(tm[2] == nil, time.Now().Day(), tm[2]).(int)
		hour:= Ternary(tm[3] == nil, time.Now().Hour(), tm[3]).(int)
		minute:= Ternary(tm[4] == nil, time.Now().Minute(), tm[4]).(int)
		second:= Ternary(tm[5] == nil, time.Now().Second(), tm[5]).(int)
		nanosecond := Ternary(tm[6] == nil, time.Now().Nanosecond(), tm[6]).(int)
		loc := Ternary(tm[7] == nil, time.Now().Location(), tm[7]).(*time.Location)
		return time.Date(year,month,day,hour,minute,second,nanosecond,loc)
	}

	if TypesCheck(input, "map") {
		tm := input.(map[string] interface{})
		year := Ternary(tm["year"] == nil, time.Now().Year(), tm["year"]).(int)
		month:= Ternary(tm["month"] == nil, time.Now(), tm["month"]).(time.Time).Month()
		day:= Ternary(tm["day"] == nil, time.Now().Day(), tm["day"]).(int)
		hour:= Ternary(tm["hour"] == nil, time.Now().Hour(), tm["hour"]).(int)
		minute:= Ternary(tm["minute"] == nil, time.Now().Minute(), tm["minute"]).(int)
		second:= Ternary(tm["second"] == nil, time.Now().Second(), tm["second"]).(int)
		nanosecond := Ternary(tm["nanosecond"] == nil, time.Now().Nanosecond(), tm["nanosecond"]).(int)
		loc := Ternary(tm["loc"] == nil, time.Now().Location(), tm["loc"]).(*time.Location)
		return time.Date(year,month,day,hour,minute,second,nanosecond,loc)
	}
	
	if TypesCheck(input, "str") {
		return dateStrParse(input.(string));
	}

	panic("unable to convert to Date")
}


// plain	   = "2020-04-09-06-05-45"
//
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
	format := dateLayout((args[0]).(string))
	dates := time.Now()
	if args[1] != nil {
		dates = Date(args[1])
	}

	return dates.Format(format)
}

func MapGet(aSet map[string]interface{}, keys ...string) map[string]interface{} {
	aMap := make(map[string]interface{})
	for _ , v := range keys {
		aMap[v] = aSet[v]
	}
	return aMap
}

func MapKeys(aSet map[string]interface{}) []string {
	keys := make([]string, 0, len(aSet))
    for k := range aSet {
        keys = append(keys, k)
    }
	return keys
}

func MapValues(aSet map[string]interface{}) []interface{} {
	values := make([]interface{}, 0, len(aSet))
	for v := range aSet {
        values = append(values, aSet[v])
    }
	return values
}

func MapGetExist(aSet map[string]interface{}, keys...string) map[string]interface{} {
	aMap := make(map[string]interface{})
	for _ , v := range keys {
		if aSet[v] != nil {
			aMap[v] = aSet[v]
		}
	}
	return aMap
}

func MapHas(aSet map[string]interface{}, key string) bool {
	return aSet[key] != nil
}

// fallback = interface{}
func MapGetPath(aSet map[string]interface{}, patharrThenFallback ...interface{}) interface{} {
	patf := CP2M(patharrThenFallback)
	patharr := patf[0].([]interface{})
	var fallback interface{}
	if patf[1] != nil {
		fallback = patf[1]
	}

	last := ArrayPopRight(&patharr)
	value := aSet
	for _, v := range patharr{	
		va:= v.(string)
			if MapHas(value , va) && TypesCheck(value[va],"map") {
				value = value[va].(map[string]interface{})
			}else {
				return fallback
			}	
	}

	if MapHas(value, last.(string)) {
		return value[last.(string)]
	} else {
		return fallback
	}
}

// space = "\t"
func JsonToString(aSetThenSpace ...interface{}) (string, error) {
	sts := CP2M(aSetThenSpace)
	aSet := sts[0].(map[string]interface{})
	space := "\t"
	if sts[1] != nil {
		space = sts[1].(string)
	}

	bytes, err := json.MarshalIndent(aSet,"",space)
    if err != nil {
        return "" , err
    }

	return string(bytes), nil
}

func StringToJson(str string) (map[string]interface{} ,error) {
    rawIn := json.RawMessage(str)
    bytes, err := rawIn.MarshalJSON()
	empty := make(map[string]interface{})
    if err != nil {
        return empty , err
    }
	var result map[string]interface{}
    err  = json.Unmarshal(bytes, &result)
    if err != nil {
		return empty , err
    }
	return result, nil
}