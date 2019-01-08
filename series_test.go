package griz

import (
	"strings"
	"testing"
	"time"
)

func TestStringSeries(t *testing.T) {
	data := []string{"asdf", "fdsa", "asfd", "fdas"}
	s := NewSeries(data, "asdf")
	if want, have := StringType, s.DataType; have != want {
		t.Errorf("have: %d, want: %d", have, want)
	}
	str := s.String()
	if strings.Contains(str, "panic") {
		t.Errorf(str)
	}
}

func TestFloatSeries(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	s := NewSeries(data, "asdf")
	if want, have := FloatType, s.DataType; have != want {
		t.Errorf("have: %d, want: %d", have, want)
	}
	str := s.String()
	if strings.Contains(str, "panic") {
		t.Errorf(str)
	}
}
func TestBoolSeries(t *testing.T) {
	data := []bool{true, false, true, false}
	s := NewSeries(data, "asdf")
	if want, have := BoolType, s.DataType; have != want {
		t.Errorf("have: %d, want: %d", have, want)
	}
	str := s.String()
	if strings.Contains(str, "panic") {
		t.Errorf(str)
	}
}

func TestTimeSeries(t *testing.T) {
	data := []time.Time{time.Now(), time.Now(), time.Now()}
	s := NewSeries(data, "asdf")
	if want, have := TimeType, s.DataType; have != want {
		t.Errorf("have: %d, want: %d", have, want)
	}
	str := s.String()
	if strings.Contains(str, "panic") {
		t.Errorf(str)
	}
}

func TestMeanSeries(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	s := NewSeries(data, "avg")
	if want, have := 5.5, s.Mean(); have != want {
		t.Errorf("have: %v, want: %v", have, want)
	}
}
func TestMinSeries(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	s := NewSeries(data, "avg")
	if want, have := 1.0, s.Min(); have != want {
		t.Errorf("have: %v, want: %v", have, want)
	}
}
func TestMaxSeries(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	s := NewSeries(data, "avg")
	if want, have := 10.0, s.Max(); have != want {
		t.Errorf("have: %v, want: %v", have, want)
	}
}
