package griz

import (
	"fmt"
	"reflect"
	"strconv"
)

func isIn(test string, list []string) bool {
	for i := 0; i < len(list); i++ {
		if test == list[i] {
			return true
		}
	}
	return false
}

func index(test string, list []string) int {
	for i := 0; i < len(list); i++ {
		if test == list[i] {
			return i
		}
	}
	panicf("index: could not locate: %s in list", test)
	return -1
}

func valueLen(a reflect.Value) int {
	v := reflect.Indirect(reflect.ValueOf(a.Interface()))
	return v.Len()
}

func dataType(v interface{}) int {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Slice, reflect.Array:
		el := t.Elem()
		if el.Kind() == reflect.Ptr {
			el = el.Elem()
		}
		switch el.Kind() {
		case reflect.Float64:
			return FloatType
		case reflect.String:
			return StringType
		case reflect.Bool:
			return BoolType
		default:
			panicf("data type not supported: %s", t.Kind().String())
		}
	case reflect.String:
		return StringType
	case reflect.Float64:
		return FloatType
	case reflect.Bool:
		return BoolType
	default:
		panicf("data type not supported: %s", t.Kind().String())
	}
	return 0
}

func rowToString(row []interface{}) []string {
	out := make([]string, len(row))
	for i := range row {
		switch dataType(row) {
		case FloatType:
			if row[i].(float64) < 1 {
				out[i] = fmt.Sprintf("%.4f", row[i])
			} else {
				out[i] = fmt.Sprintf("%.2f", row[i])
			}
		case StringType:
			out[i] = fmt.Sprint(row[i])
		}
	}
	return out
}

func panicf(format string, a ...interface{}) {
	panic(fmt.Sprintf(format, a...))
}

func scanArgs(values []interface{}) []interface{} {
	ptrs := make([]interface{}, len(values))
	for i := 0; i < len(values); i++ {
		ptrs[i] = &values[i]
	}
	return ptrs
}

func stringToFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 10)
	if err != nil {
		panicf("string to float: error parsing: %s", s)
	}
	return f
}

func stringToBool(s string) bool {
	switch s {
	case "0", "false", "False":
		return false
	case "1", "true", "True":
		return true
	default:
		panicf("string to bool: error parsing: %s", s)
	}
	return false
}
