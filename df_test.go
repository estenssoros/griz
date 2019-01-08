package griz

import (
	"testing"
	"time"
)

func TestDataframe(t *testing.T) {
	data := []interface{}{
		[]interface{}{"asdf", 1123.0, 12.5, false, time.Now()},
		[]interface{}{"asdf", 1123.0, 12.5, false, time.Now()},
		[]interface{}{"asdf", 1123.0, 12.5, false, time.Now()},
		[]interface{}{"asdf", 1123.0, 12.5, false, time.Now()},
		[]interface{}{"asdf", 1123.0, 12.5, true, time.Now()},
		[]interface{}{"asdf", 1123.0, 12.5, false, time.Now()},
	}
	columns := []string{"col1", "col2", "col3", "col4", "time"}
	dataTypes := []int{StringType, FloatType, FloatType, BoolType, TimeType}
	df := NewDataFrame(data, columns, dataTypes)

	if want, have := "asdf", df.Iloc(0).Loc("col1"); want != have {
		t.Errorf("have: %s, want: %s", have, want)
	}

	val := 17.5
	col := "col3"
	df.SetLoc(1, col, val)
	if want, have := val, df.Iloc(1).Loc(col); want != have {
		t.Errorf("have: %f, want: %f", have, want)
	}
	df.String()
}

func createTestDF() *DataFrame {
	data := []interface{}{
		[]interface{}{"asdf", 1123.0, 12.5, false, "asdf"},
		[]interface{}{"asdf", 1123.0, 12.5, false, "asdf"},
		[]interface{}{"asdf", 1123.0, 12.5, false, "fdas"},
		[]interface{}{"asdf", 1123.0, 12.5, false, "asdf"},
		[]interface{}{"asdf", 1123.0, 12.5, false, "fdas"},
		[]interface{}{"asdf", 1123.0, 12.5, false, "asdf"},
		[]interface{}{"asdf", 1123.0, 12.5, true, "asdf"},
		[]interface{}{"asdf", 1123.0, 12.5, false, "fdsa"},
	}
	columns := []string{"col1", "col2", "col3", "col4", "col5"}
	dataTypes := []int{StringType, FloatType, FloatType, BoolType, StringType}
	return NewDataFrame(data, columns, dataTypes)
}

func TestAddition(t *testing.T) {
	df := createTestDF()
	s := df.GetSeries("col2").AddSeries(df.GetSeries("col3"))
	val := 12.5
	if want, have := val+1123.0, s.Iloc(1); want != have {
		t.Errorf("have: %f, want: %f", have, want)
	}
	df.String()
}

func TestMultiplication(t *testing.T) {
	df := createTestDF()
	s := df.GetSeries("col2").MulSeries(df.GetSeries("col3"))
	val := 12.5
	if want, have := val*1123.0, s.Iloc(1); want != have {
		t.Errorf("have: %f, want: %f", have, want)
	}
	df.String()
}
func TestSubstraction(t *testing.T) {
	df := createTestDF()
	val := 12.5
	s := df.GetSeries("col2").SubSeries(df.GetSeries("col3"))
	if want, have := 1123.0-val, s.Iloc(1); want != have {
		t.Errorf("have: %f, want: %f", have, want)
	}
}
func TestCumSum(t *testing.T) {
	df := createTestDF()
	s := df.GetSeries("col2").CumSum()
	if want, have := 2246.00, s.Iloc(1); want != have {
		t.Errorf("have: %f, want: %f", have, want)
	}
	df.String()
}

func TestWhere(t *testing.T) {
	df := createTestDF()
	df = df.Where(df.GetSeries("col4"))
	if want, have := 1, df.Len(); want != have {
		t.Errorf("have: %d, want: %d", have, want)
	}
}

func TestGetDataFrame(t *testing.T) {
	df := createTestDF()
	columns := []string{"col1", "col4", "col3"}
	df = df.GetDataFrame(columns...)
	if want, have := len(columns), df.Width(); want != have {
		t.Errorf("have: %d, want: %d", have, want)
	}
	for i := 0; i < len(columns); i++ {
		if want, have := columns[i], df.Columns[i]; want != have {
			t.Errorf("have: %v, want: %v", have, want)
		}
	}
	columns = []string{"col4"}
	df = df.GetDataFrame(columns...)
	for i := 0; i < len(columns); i++ {
		if want, have := columns[i], df.Columns[i]; want != have {
			t.Errorf("have: %v, want: %v", have, want)
		}
	}
	df = df.Where(df.GetSeries("col4"))
	if want, have := 1, df.Len(); want != have {
		t.Errorf("have: %v, want: %v", have, want)
	}
	if want, have := 1, df.Width(); want != have {
		t.Errorf("have: %v, want: %v", have, want)
	}
	df.String()
}

func TestEquals(t *testing.T) {
	df := createTestDF()
	df = df.Where(df.GetSeries("col1").Equals(df.GetSeries("col5")))
	if want, have := 5, df.Len(); want != have {
		t.Errorf("have: %v, want: %v", have, want)
	}
}

func TestNotEquals(t *testing.T) {
	df := createTestDF()
	df = df.Where(df.GetSeries("col1").NotEquals(df.GetSeries("col5")))
	if want, have := 3, df.Len(); want != have {
		t.Errorf("have: %v, want: %v", have, want)
	}
}

func TestAppend(t *testing.T) {
	df := createTestDF()
	data := make([]string, df.Len())
	for i := 0; i < df.Len(); i++ {
		data[i] = "astring"
	}
	before := df.Width()
	s := NewSeries(data, "append_me")
	df = df.AddColumn(s)
	after := df.Width()
	if want, have := after, before+1; want != have {
		t.Errorf("have: %v, want: %v", have, want)
	}
}

func TestNewSeriesFromValue(t *testing.T) {
	df := createTestDF()
	s := df.NewSeriesFromValue("asdf", "new_series")
	before := df.Width()
	df = df.AddColumn(s)
	after := df.Width()
	if want, have := after, before+1; want != have {
		t.Errorf("have: %v, want: %v", have, want)
	}
}

func createMathDF() *DataFrame {
	data := []interface{}{
		[]interface{}{1.0, 1123.0, 12.5, 5.0, 10.0},
		[]interface{}{1.0, 1123.0, 12.5, 5.0, 10.0},
		[]interface{}{1.0, 1123.0, 12.5, 5.0, 10.0},
		[]interface{}{1.0, 1123.0, 12.5, 5.0, 10.0},
		[]interface{}{1.0, 1123.0, 12.5, 5.0, 10.0},
		[]interface{}{1.0, 1123.0, 12.5, 5.0, 10.0},
		[]interface{}{1.0, 1123.0, 12.5, 5.0, 10.0},
	}
	columns := []string{"col1", "col2", "col3", "col4", "time"}
	dataTypes := []int{FloatType, FloatType, FloatType, FloatType, FloatType}
	return NewDataFrame(data, columns, dataTypes)
}

func TestDFAverage(t *testing.T) {
	df := createMathDF()
	df = df.AddColumn(df.Mean())
	if want, have := 230.30, df.Iloc(0).Loc("avg").(float64); want != have {
		t.Errorf("have: %v, want: %v", have, want)
	}
}
func TestDFMin(t *testing.T) {
	df := createMathDF()
	df = df.AddColumn(df.Min())
	if want, have := 1.0, df.Iloc(0).Loc("min").(float64); want != have {
		t.Errorf("have: %v, want: %v", have, want)
	}
}
func TestDFMax(t *testing.T) {
	df := createMathDF()
	df = df.AddColumn(df.Max())
	if want, have := 1123.0, df.Iloc(0).Loc("max").(float64); want != have {
		t.Errorf("have: %v, want: %v", have, want)
	}
}
