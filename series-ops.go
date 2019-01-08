package griz

import (
	"github.com/montanaflynn/stats"
	"github.com/pkg/errors"
)

func (s *Series) checkAreFloat(other *Series) {
	if s.DataType != FloatType || other.DataType != FloatType {
		panicf("data types must be float. this: %d, other: %d", s.DataType, other.DataType)
	}
}

// MulValue returns a series multiplied by another value
func (s *Series) MulValue(value float64) *Series {
	mat := make([]float64, s.Len())
	for i := range s.FloatMat {
		mat[i] = s.FloatMat[i] * value
	}
	return newFloatSeries(mat, s.Name)
}

// CumSum performs the cumsum of a series
func (s *Series) CumSum() *Series {
	mat := make([]float64, s.Len())
	mat, _ = stats.CumulativeSum(s.FloatMat)
	return newFloatSeries(mat, s.Name)
}

// Invert returns a new series that is inverted
func (s *Series) Invert() *Series {
	mat := make([]float64, s.Len())
	for i, j := 0, len(mat)-1; i < j; i, j = i+1, j-1 {
		mat[i], mat[j] = s.FloatMat[j], s.FloatMat[i]
	}
	return newFloatSeries(mat, s.Name)
}

// Rename renames a series
func (s *Series) Rename(name string) *Series {
	s.Name = name
	return s
}

// AddSeries adds a series to another
func (s *Series) AddSeries(other *Series) *Series {
	s.checkAreFloat(other)
	if s.Len() != other.Len() {
		panic(errors.Errorf("dimension mismatch: this: %d other: %d", len(s.FloatMat), len(other.FloatMat)))
	}
	mat := make([]float64, s.Len())
	for i := 0; i < s.Len(); i++ {
		mat[i] = s.FloatMat[i] + other.FloatMat[i]
	}
	return newFloatSeries(mat, s.Name)
}

// SubSeries subtracts a series from another inplace
func (s *Series) SubSeries(other *Series) *Series {
	s.checkAreFloat(other)
	if s.Len() != other.Len() {
		panic(errors.Errorf("dimension mismatch: this: %d other: %d", len(s.FloatMat), len(other.FloatMat)))
	}
	mat := make([]float64, s.Len())
	for i := range s.FloatMat {
		mat[i] = s.FloatMat[i] - other.FloatMat[i]
	}
	return newFloatSeries(mat, s.Name)
}

// MulSeries returns a series that is the multiple of another series
func (s *Series) MulSeries(other *Series) *Series {
	s.checkAreFloat(other)
	if s.Len() != other.Len() {
		panicf("dimension mismatch: this: %d other: %d", len(s.FloatMat), len(other.FloatMat))
	}
	mat := make([]float64, s.Len())
	for i := range s.FloatMat {
		mat[i] = s.FloatMat[i] * other.FloatMat[i]
	}
	return newFloatSeries(mat, s.Name)
}

func (s *Series) Mean() float64 {
	if s.DataType != FloatType {
		panic("series mean: only float types supported")
	}
	m, _ := stats.Mean(s.FloatMat)
	return m
}

func (s *Series) Min() float64 {
	if s.DataType != FloatType {
		panic("series min: only float types supported")
	}
	m, _ := stats.Min(s.FloatMat)
	return m
}

func (s *Series) Max() float64 {
	if s.DataType != FloatType {
		panic("series max: only float types supported")
	}
	m, _ := stats.Max(s.FloatMat)
	return m
}

func (s *Series) Quartile(quartile int) float64 {
	if s.DataType != FloatType {
		panic("series quartile: only float types supported")
	}
	if quartile < 1 || quartile > 3 {
		panicf("quartile out of range: %d", quartile)
	}
	q, _ := stats.Quartile(s.FloatMat)
	switch quartile {
	case 1:
		return q.Q1
	case 2:
		return q.Q2
	case 3:
		return q.Q3
	}
	return 0.0
}
