package griz

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Row holds row data
type Row struct {
	Idx       int
	Columns   []string
	DataTypes []int
	Value     interface{}
}

func (r Row) String() string {
	row := map[string]interface{}{}
	rowV := reflect.Indirect(reflect.ValueOf(r.Value))
	for i, c := range r.Columns {
		row[c] = rowV.Index(i).Interface()
	}
	ju, _ := json.Marshal(row)
	return string(ju)
}

func (r *Row) Map() map[string]interface{} {
	out := map[string]interface{}{}
	for i, c := range r.Columns {
		out[c] = r.Value.([]interface{})[i]
	}
	return out
}

// ToString converts the row value to a []string
func (r *Row) ToString() []string {
	out := []string{}
	v := reflect.Indirect(reflect.ValueOf(r.Value))
	for i := 0; i < v.Len(); i++ {
		switch dataType(v.Index(i).Interface()) {
		case StringType:
			out = append(out, v.Index(i).Interface().(string))
		case FloatType:
			out = append(out, fmt.Sprintf("%.2f", v.Index(i).Interface().(float64)))
		case BoolType:
			out = append(out, fmt.Sprint(v.Index(i).Interface().(bool)))
		default:
			panicf("row ToString(): unknown dataType: %d", dataType(v.Index(i).Interface()))
		}
	}
	return out
}

// Iloc returns the value at an index
func (r *Row) Iloc(idx int) interface{} {
	v := reflect.Indirect(reflect.ValueOf(r.Value))
	dataType := r.DataTypes[idx]
	switch dataType {
	case FloatType:
		return v.Index(idx).Interface().(float64)
	case StringType:
		return v.Index(idx).Interface().(string)
	case BoolType:
		return v.Index(idx).Interface().(bool)
	default:
		panicf("unknown datatype: %d", dataType)
	}
	return nil
}

// Loc returns the interface value of a column in a row
func (r *Row) Loc(column string) interface{} {
	idx := index(column, r.Columns)
	return r.Iloc(idx)
}

// SetLoc sets the column value of a row to a value
func (r *Row) SetLoc(column string, value interface{}) {
	idx := index(column, r.Columns)
	v := reflect.Indirect(reflect.ValueOf(r.Value))
	v.Index(idx).Set(reflect.ValueOf(value))
}
