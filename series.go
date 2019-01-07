package griz

import (
	"reflect"
	"time"
)

const (
	BoolType   = 0
	FloatType  = iota
	StringType = iota
)

// DataTypeString converts a data type to string
func DataTypeString(i int) string {
	switch i {
	case BoolType:
		return "bool"
	case FloatType:
		return "float"
	case StringType:
		return "string"
	default:
		panicf("unknown datatype: %d", i)
	}
	return ""
}

// Series holds data in a columnar format
type Series struct {
	Name      string
	FloatMat  []float64
	StringMat []string
	BoolMat   []bool
	DateMat   []Date
	TimeMat   []time.Time
	DataType  int
}

// Len returns the length of the series
func (s *Series) Len() int {
	switch s.DataType {
	case FloatType:
		return len(s.FloatMat)
	case StringType:
		return len(s.StringMat)
	case BoolType:
		return len(s.BoolMat)
	default:
		panicf("unknown dataType: %d", s.DataType)
	}
	return 0
}

func (s Series) String() string {
	switch s.DataType {
	case FloatType:
		return s.FloatString()
	case StringType:
		return s.StringString()
	case BoolType:
		return s.BoolString()
	default:
		panicf("series String(): unknown dataType: %d", s.DataType)
	}
	return ""
}

// Head returns the top rows rows of a series
func (s *Series) Head(rows int) {
	if rows >= s.Len() {
		panic("head index out of bounds")
	}
	switch s.DataType {
	case FloatType:
		s.FloatHead(rows)
	case StringType:
		s.StringHead(rows)
	case BoolType:
		s.BoolHead(rows)
	default:
		panicf("unknown dataType: %d", s.DataType)
	}
}

// NewSeries creates a new series
func NewSeries(data interface{}, name string) *Series {
	t := reflect.TypeOf(data)
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
			return newFloatSeries(data, name)
		case reflect.String:
			return newStringSeries(data, name)
		case reflect.Bool:
			return newBoolSeries(data, name)
		default:
			panicf("new series: data type not supported: %s", el.Kind().String())
		}
	default:
		panicf("new series: data must be array not %s", t.Kind().String())
	}
	return &Series{}
}

// Iloc returns the interface value at a location
func (s *Series) Iloc(idx int) interface{} {
	switch s.DataType {
	case FloatType:
		return s.FloatMat[idx]
	case StringType:
		return s.StringMat[idx]
	case BoolType:
		return s.BoolMat[idx]
	default:
		panicf("series iloc: data type not supported: %d", s.DataType)
	}
	return nil
}

// Bool gets the bool value at an index
func (s *Series) Bool(idx int) bool {
	if s.DataType != BoolType {
		panic("bool only supports bool")
	}
	return s.BoolMat[idx]
}

// Equals returns a boolean series of the comparison between two series
func (s *Series) Equals(other *Series) *Series {
	if s.Len() != other.Len() {
		panicf("series equal: dimension mismatch: this: %d other: %d", s.Len(), other.Len())
	}
	if s.DataType != other.DataType {
		panicf("series equal: can only compare same datatypes")
	}
	switch s.DataType {
	case FloatType:
		return s.floatSeriesEqual(other)
	case StringType:
		return s.stringSeriesEqual(other)
	case BoolType:
		return s.boolSeriesEqual(other)
	default:
		panicf("series equals: data type not supported %d", s.DataType)
	}
	return nil
}

// NotEquals returns a bool series where two series are equal
func (s *Series) NotEquals(other *Series) *Series {
	if s.Len() != other.Len() {
		panicf("series equal: dimension mismatch: this: %d other: %d", s.Len(), other.Len())
	}
	if s.DataType != other.DataType {
		panicf("series equal: can only compare same datatypes")
	}
	switch s.DataType {
	case FloatType:
		return s.floatSeriesNotEqual(other)
	case StringType:
		return s.stringSeriesNotEqual(other)
	case BoolType:
		return s.boolSeriesNotEqual(other)
	default:
		panicf("series equals: data type not supported %d", s.DataType)
	}
	return nil
}
