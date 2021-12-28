package u

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
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

// array to interface
func ATI(array interface{}) []interface{} {
	slice := reflect.ValueOf(array)
	result := make([]interface{}, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		result[i] = slice.Index(i).Interface()
	}
	return result
}

// map to interface
func MTI(amap interface{}) map[string]interface{} {
	dict := reflect.ValueOf(amap)
	result := make(map[string]interface{})
	iter := dict.MapRange()
	for iter.Next() {
		result[ToString(iter.Key())] = iter.Value().Interface()
	}
	return result
}

func Array(item ...interface{}) []interface{} {
	return item
}

func Map(keyThenValue ...interface{}) map[string]interface{} {
	aMap := make(map[string]interface{})
	if len(keyThenValue)%2 != 0 {
		panic("Map() requires an even number of arguments")
	}
	for i := 0; i < len(keyThenValue); i += 2 {
		aMap[ToString(keyThenValue[i])] = keyThenValue[i+1]
	}
	return aMap
}

func MapToStruct(aMap map[string]interface{}, aStruct interface{}) error {
	aValue := reflect.ValueOf(aStruct)
	if aValue.Kind() != reflect.Ptr {
		return errors.New("MapToStruct() requires a pointer to a struct")
	}
	aValue = aValue.Elem()
	if aValue.Kind() != reflect.Struct {
		return errors.New("MapToStruct() requires a pointer to a struct")
	}
	for key, value := range aMap {
		aValue.FieldByName(key).Set(reflect.ValueOf(value))
	}
	return nil
}

func ArrayAdd(item ...interface{}) []interface{} {
	return append(item[:0:0], item...)
}

func ArrayToMap(array []interface{}) map[int]interface{} {
	aMap := make(map[int]interface{})
	for i := 0; i < len(array); i++ {
		aMap[i] = array[i]
	}
	return aMap
}

func ArrayPopRight(array *[]interface{}) interface{} {
	if len(*array) == 0 {
		return nil
	}
	pop := (*array)[len(*array)-1]
	(*array) = (*array)[:len(*array)-1]
	return pop
}

func ArrayPopLeft(array *[]interface{}) interface{} {
	if len(*array) == 0 {
		return nil
	}
	shift := (*array)[0]
	(*array) = (*array)[1:]
	return shift
}

// END = length
func ArrayExtract(array []interface{}, startThenEnd ...interface{}) []interface{} {
	args := CP2M(startThenEnd)
	start := args[0].(int)
	end := len(array)
	if args[1] != nil {
		end = args[1].(int)
	}

	return array[start:end]
}

func ArrayEqual(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func ToString(item interface{}) string {
	return fmt.Sprint(item)
}

// sep = ","
func ArrayToString(arrayThenSep ...interface{}) string {
	ats := CP2M(arrayThenSep)
	array := arrayThenSep[0].([]interface{})
	sep := ","
	if ats[1] != nil {
		sep = ats[1].(string)
	}

	result := ""
	for i := 0; i < len(array); i++ {
		result += ToString(array[i]) + sep
	}
	return strings.TrimSuffix(result, sep)
}

// Change Param To Map, implementation of kwargs
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

// expect : str | string | num | number | int | float | arr | array | map | dict | date | time
func TypesCheck(source interface{}, expect string) bool {
	stype := Types(source)

	switch strings.ToLower(expect) {
	case "str", "string":
		return stype == "string"
	case "num", "number":
		return strings.Contains(stype, "int") || strings.Contains(stype, "float")
	case "int":
		return strings.Contains(stype, "int")
	case "float":
		return strings.Contains(stype, "float")
	case "arr", "array":
		return stype == "array"
	case "map", "dict":
		return stype == "map"
	case "date", "time":
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
	t, e := dateparse.ParseLocal(str)
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
		return time.Unix(tn/1000, tn%1000*int64(time.Millisecond))
	}

	if TypesCheck(input, "array") {
		tm := ArrayToMap(input.([]interface{}))
		year := Ternary(tm[0] == nil, time.Now().Year(), tm[0]).(int)
		month := Ternary(tm[1] == nil, int(time.Now().Month()), tm[1]).(int)
		day := Ternary(tm[2] == nil, time.Now().Day(), tm[2]).(int)
		hour := Ternary(tm[3] == nil, time.Now().Hour(), tm[3]).(int)
		minute := Ternary(tm[4] == nil, time.Now().Minute(), tm[4]).(int)
		second := Ternary(tm[5] == nil, time.Now().Second(), tm[5]).(int)
		nanosecond := Ternary(tm[6] == nil, time.Now().Nanosecond(), tm[6]).(int)
		loc := Ternary(tm[7] == nil, time.Now().Location(), tm[7]).(*time.Location)
		return time.Date(year, time.Month(month), day, hour, minute, second, nanosecond, loc)
	}

	if TypesCheck(input, "map") {
		tm := input.(map[string]interface{})
		year := Ternary(tm["year"] == nil, time.Now().Year(), tm["year"]).(int)
		month := Ternary(tm["month"] == nil, int(time.Now().Month()), tm["month"]).(int)
		day := Ternary(tm["day"] == nil, time.Now().Day(), tm["day"]).(int)
		hour := Ternary(tm["hour"] == nil, time.Now().Hour(), tm["hour"]).(int)
		minute := Ternary(tm["minute"] == nil, time.Now().Minute(), tm["minute"]).(int)
		second := Ternary(tm["second"] == nil, time.Now().Second(), tm["second"]).(int)
		nanosecond := Ternary(tm["nanosecond"] == nil, time.Now().Nanosecond(), tm["nanosecond"]).(int)
		loc := Ternary(tm["loc"] == nil, time.Now().Location(), tm["loc"]).(*time.Location)
		return time.Date(year, time.Month(month), day, hour, minute, second, nanosecond, loc)
	}

	if TypesCheck(input, "str") {
		return dateStrParse(input.(string))
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

// diff : map[string]interface{} key: year?, month?, day?, hour?, minute?, second?, nanosecond?
//
// Date: time.Time, default to now
func DateAdd(diffThenDate ...interface{}) time.Time {
	args := CP2M(diffThenDate)
	diff := (args[0]).(map[string]interface{})
	date := time.Now()
	if args[1] != nil {
		date = args[1].(time.Time)
	}

	diff["year"] = Ternary(diff["year"] != nil, diff["year"], 0)
	diff["month"] = Ternary(diff["month"] != nil, diff["month"], 0)
	diff["day"] = Ternary(diff["day"] != nil, diff["day"], 0)
	diff["hour"] = Ternary(diff["hour"] != nil, diff["hour"], 0)
	diff["minute"] = Ternary(diff["minute"] != nil, diff["minute"], 0)
	diff["second"] = Ternary(diff["second"] != nil, diff["second"], 0)
	diff["nanosecond"] = Ternary(diff["nanosecond"] != nil, diff["nanosecond"], 0)

	return date.Local().AddDate(diff["year"].(int), diff["month"].(int), diff["day"].(int)).Add(
		time.Hour*time.Duration(diff["hour"].(int)) +
			time.Minute*time.Duration(diff["minute"].(int)) +
			time.Second*time.Duration(diff["second"].(int)) +
			time.Nanosecond*time.Duration(diff["nanosecond"].(int)))
}

func MapGet(aSet map[string]interface{}, keys ...string) map[string]interface{} {
	aMap := make(map[string]interface{})
	for _, v := range keys {
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

func MapGetExist(aSet map[string]interface{}, keys ...string) map[string]interface{} {
	aMap := make(map[string]interface{})
	for _, v := range keys {
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
	for _, v := range patharr {
		va := v.(string)
		if MapHas(value, va) && TypesCheck(value[va], "map") {
			value = value[va].(map[string]interface{})
		} else {
			return fallback
		}
	}

	if MapHas(value, last.(string)) {
		return value[last.(string)]
	} else {
		return fallback
	}
}

func MapMerge(sets ...map[string]interface{}) map[string]interface{} {
	aMap := make(map[string]interface{})
	for _, v := range sets {
		for k, vv := range v {
			aMap[k] = vv
		}
	}
	return aMap
}

func MapMergeDeep(sets ...map[string]interface{}) map[string]interface{} {
	aMap := make(map[string]interface{})
	for _, v := range sets {
		for k, vv := range v {
			if TypesCheck(aMap[k], "map") && TypesCheck(vv, "map") {
				aMap[k] = MapMergeDeep(aMap[k].(map[string]interface{}), vv.(map[string]interface{}))
			} else {
				aMap[k] = vv
			}
		}
	}
	return aMap
}

func MapEqual(a, b map[string]interface{}) bool {
	return reflect.DeepEqual(a, b)
}

func Serialize(aMap map[string]interface{}) string {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(aMap)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func Deserialize(aString string) map[string]interface{} {
	var buf bytes.Buffer
	buf.WriteString(aString)
	dec := gob.NewDecoder(&buf)
	var aMap map[string]interface{}
	err := dec.Decode(&aMap)
	if err != nil {
		panic(err)
	}
	return aMap
}

// space = "\t"
func JsonToString(aSetThenSpace ...interface{}) (string, error) {
	sts := CP2M(aSetThenSpace)
	aSet := sts[0].(map[string]interface{})
	space := "\t"
	if sts[1] != nil {
		space = sts[1].(string)
	}

	bytes, err := json.MarshalIndent(aSet, "", space)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func StringToJson(str string) (map[string]interface{}, error) {
	rawIn := json.RawMessage(str)
	bytes, err := rawIn.MarshalJSON()
	empty := make(map[string]interface{})
	if err != nil {
		return empty, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return empty, err
	}
	return result, nil
}

func StringContains(str string, substr string) bool {
	return strings.Contains(str, substr)
}

func StringReplace(str string, repl map[string]string, recursive bool) string {
	times := Ternary(recursive, -1, 1).(int)
	for k, v := range repl {
		str = strings.Replace(str, k, v, times)
	}
	return str
}

// use `` instead of // for regex
func ReFind(str string, pattern string) string {
	re := regexp.MustCompile(pattern)
	return re.FindString(str)
}

func ReHas(str string, pattern string) bool {
	re := regexp.MustCompile(pattern)
	return re.MatchString(str)
}

func ReFindAll(str string, pattern string) []string {
	re := regexp.MustCompile(pattern)
	return re.FindAllString(str, -1)
}

func ReSub(str string, pattern string, repl string) string {
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(str, repl)
}

func ReSubMap(str string, pattern string, repl map[string]interface{}) string {
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllStringFunc(str, func(s string) string {
		return repl[s].(string)
	})
}

func Password(length int, strong bool) string {
	var charset = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	if strong {
		charset = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+{}[]|\\:;'<>,.?/")
	}
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func Url(url string) string {
	if ReHas(url, `^localhost`) {
		return "http://" + url
	}
	if url == "about:blank" {
		return url
	}
	if ReHas(url, `^www\.`) {
		return "https://" + url
	}
	if ReHas(url, `^http`) {
		return url
	}
	if ReHas(url, `^\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`) {
		return "http://" + url
	}
	return "https://www." + url
}

func Fetch(url string, method string, header map[string]interface{}, data map[string]interface{}) (string, error) {
	var client = &http.Client{}
	url = Url(url)
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	var rBody io.Reader = nil
	if method != "GET" {
		rBody = bytes.NewReader(jsonBytes)
	}
	req, err := http.NewRequest(method, url, rBody)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")

	for k, v := range header {
		req.Header.Set(k, v.(string))
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		err := Map("status", resp.StatusCode, "body", string(body), "url", url, "method", method, "header", req.Header)
		result, _ := JsonToString(err)
		return string(body), errors.New(result)
	}

	return string(body), nil
}

func FetchGet(url string, header map[string]interface{}) (string, error) {
	return Fetch(url, "GET", header, nil)
}

func FetchPost(url string, data map[string]interface{}, header map[string]interface{}) (string, error) {
	return Fetch(url, "POST", header, data)
}

func FetchUploadFile(url string, files map[string]interface{}, body map[string]interface{}, header map[string]interface{}) (string, error) {
	var client = &http.Client{}
	url = Url(url)
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")

	for k, v := range header {
		req.Header.Set(k, v.(string))
	}

	// multipart form
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	for k, v := range body {
		bodyWriter.WriteField(k, v.(string))
	}

	for k, v := range files {
		f, err := os.Open(v.(string))
		if err != nil {
			return "", err
		}
		defer f.Close()

		fileWriter, err := bodyWriter.CreateFormFile(k, filepath.Base(v.(string)))
		if err != nil {
			return "", err
		}

		_, err = io.Copy(fileWriter, f)
		if err != nil {
			return "", err
		}
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-Length", strconv.Itoa(bodyBuf.Len()))

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	resbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		err := Map("status", resp.StatusCode, "body", string(resbody), "url", url, "header", req.Header)
		result, _ := JsonToString(err)
		return string(resbody), errors.New(result)
	}

	return string(resbody), nil
}

// use global variable to pass data
func Retry(fn func() error, retryCount int) error {
	var err error

	for i := 0; i < retryCount; i++ {
		err = fn()
		if err == nil {
			return nil
		}

		time.Sleep(time.Second * 2)
	}
	return err
}

func RetryEH(fn func() error, retryCount int, eh func(error)) error {
	var err error

	for i := 0; i < retryCount; i++ {
		err = fn()
		if err == nil {
			return nil
		}

		eh(err)
		time.Sleep(time.Second * 2)
	}
	return err
}
