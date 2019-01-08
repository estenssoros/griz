package griz

import (
	"reflect"
	"time"
)

const (
	BoolType   = 0
	FloatType  = iota
	StringType = iota
	TimeType   = iota
	UUIDType   = iota
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
	case TimeType:
		return "time"
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
	case TimeType:
		return len(s.TimeMat)
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
	case TimeType:
		return s.TimeString()
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
	case TimeType:
		s.TimeHead(rows)
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
		switch dataType(data) {
		case FloatType:
			return newFloatSeries(data, name)
		case StringType:
			return newStringSeries(data, name)
		case BoolType:
			return newBoolSeries(data, name)
		case TimeType:
			return newTimeSeries(data, name)
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
	case TimeType:
		return s.TimeMat[idx]
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
	case TimeType:
		return s.timeSeriesEqual(other)
	default:
		panicf("series equals: data type not supported %d", s.DataType)
	}
	return nil
}

func (s *Series) EqualsString(val string) *Series {
	if s.DataType != StringType {
		panic("series equals string only supports string series")
	}
	data := make([]bool, s.Len())
	for i := 0; i < s.Len(); i++ {
		data[i] = s.StringMat[i] == val
	}
	return NewSeries(data, s.Name)
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
	case TimeType:
		return s.timeSeriesNotEqual(other)
	default:
		panicf("series equals: data type not supported %d", s.DataType)
	}
	return nil
}

func (s *Series) Unique() interface{} {
	switch s.DataType {
	case FloatType:
		return s.UniqueFloat()
	case StringType:
		return s.UniqueString()
	case BoolType:
		return s.UniqueBool()
	case TimeType:
		return s.UniqueTime()
	}
	return nil
}

func (s *Series) And(other *Series) *Series {
	if s.DataType != BoolType || other.DataType != BoolType {
		panic("series and: both data types must be bool")
	}
	data := make([]bool, s.Len())
	for i := 0; i < s.Len(); i++ {
		data[i] = s.BoolMat[i] && other.BoolMat[i]
	}
	return NewSeries(data, s.Name)
}

func (s *Series) Or(other *Series) *Series {
	if s.DataType != BoolType || other.DataType != BoolType {
		panic("series and: both data types must be bool")
	}
	data := make([]bool, s.Len())
	for i := 0; i < s.Len(); i++ {
		data[i] = s.BoolMat[i] || other.BoolMat[i]
	}
	return NewSeries(data, s.Name)
}
