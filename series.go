package grizzly

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"github.com/gonum/floats"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

type Series struct {
	Name string
	Mat  []float64
}

func newSeries(data []float64, name string) *Series {
	return &Series{Name: name, Mat: data}
}

// Len returns the length of the series
func (s *Series) Len() int {
	return len(s.Mat)
}

func (s Series) String() string {
	var b bytes.Buffer
	{
		writer := bufio.NewWriter(&b)
		table := tablewriter.NewWriter(writer)
		table.SetHeader([]string{s.Name})
		for _, row := range s.Mat {
			if row < 1.0 {
				table.Append([]string{fmt.Sprintf("%.4f", row)})
			} else {
				table.Append([]string{fmt.Sprintf("%.2f", row)})
			}
		}
		table.Render()
		writer.Flush()
	}
	return b.String()
}

func (s *Series) Head(rows int) {
	if rows >= s.Len() {
		panic("head index out of bounds")
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{s.Name})
	for i := 0; i < rows; i++ {
		if s.Mat[i] < 1.0 {
			table.Append([]string{fmt.Sprintf("%.4f", s.Mat[i])})
		} else {
			table.Append([]string{fmt.Sprintf("%.2f", s.Mat[i])})
		}
	}
	table.Render()
}

// MulSeriesInplace multiplies a series by another in place
func (s *Series) MulSeriesInplace(other *Series) *Series {
	if s.Len() != other.Len() {
		panic(errors.Errorf("dimension mismatch: this: %d other: %d", len(s.Mat), len(other.Mat)))
	}
	for i := range s.Mat {
		s.Mat[i] = s.Mat[i] * other.Mat[i]
	}
	return s
}

// MulSeries returns a series that is the multiple of another series
func (s *Series) MulSeries(other *Series, name string) *Series {
	if s.Len() != other.Len() {
		panic(errors.Errorf("dimension mismatch: this: %d other: %d", len(s.Mat), len(other.Mat)))
	}
	mat := make([]float64, s.Len())
	for i := range s.Mat {
		mat[i] = s.Mat[i] * other.Mat[i]
	}
	return newSeries(mat, name)
}

// MulValueInplace multiplies a series inplace by a value
func (s *Series) MulValueInplace(value float64) {
	for i := range s.Mat {
		s.Mat[i] = s.Mat[i] * value
	}
}

// MulValue returns a series multiplied by another value
func (s *Series) MulValue(value float64) *Series {
	mat := make([]float64, s.Len())
	for i := range s.Mat {
		mat[i] = s.Mat[i] * value
	}
	return newSeries(mat, s.Name)
}

// Sub subtracts a series from another inplace
func (s *Series) Sub(other *Series) *Series {
	if s.Len() != other.Len() {
		panic(errors.Errorf("dimension mismatch: this: %d other: %d", len(s.Mat), len(other.Mat)))
	}
	for i := range s.Mat {
		s.Mat[i] = s.Mat[i] - other.Mat[i]
	}
	return s
}

// CumSum performs the cumsum of a series
func (s *Series) CumSum(name string) *Series {
	mat := make([]float64, s.Len())
	floats.CumSum(mat, s.Mat)
	return newSeries(mat, name)
}

// Invert returns a new series that is inverted
func (s *Series) Invert(name string) *Series {
	mat := make([]float64, s.Len())
	for i, j := 0, len(mat)-1; i < j; i, j = i+1, j-1 {
		mat[i], mat[j] = s.Mat[j], s.Mat[i]
	}
	return newSeries(mat, name)
}

// Rename renames a series
func (s *Series) Rename(name string) *Series {
	s.Name = name
	return s
}

func (s *Series) Add(other *Series) *Series {
	if s.Len() != other.Len() {
		panic(errors.Errorf("dimension mismatch: this: %d other: %d", len(s.Mat), len(other.Mat)))
	}
	mat := make([]float64, s.Len())
	for i := 0; i < s.Len(); i++ {
		mat[i] = s.Mat[i] + other.Mat[i]
	}
	return newSeries(mat, s.Name)
}
