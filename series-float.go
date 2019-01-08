package griz

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func newFloatSeries(data interface{}, name string) *Series {
	mat, ok := data.([]float64)
	if !ok {
		panic("new float series: error converting data to []float64")
	}
	return &Series{
		Name:     name,
		FloatMat: mat,
		DataType: FloatType,
	}
}

// FloatString writes matrix to table
func (s Series) FloatString() string {
	var b bytes.Buffer
	{
		writer := bufio.NewWriter(&b)
		table := tablewriter.NewWriter(writer)
		table.SetHeader([]string{s.Name})
		for _, row := range s.FloatMat {
			if row < 1.0 {
				table.Append([]string{fmt.Sprintf("%.4f", row)})
			} else {
				table.Append([]string{fmt.Sprintf("%.2f", row)})
			}
		}
		table.SetFooter([]string{DataTypeString(s.DataType)})
		table.SetFooterColor(tablewriter.Colors{tablewriter.FgGreenColor})
		table.Render()
		writer.Flush()
	}
	return b.String()
}

// FloatHead prints a table the first x rows of a series
func (s *Series) FloatHead(rows int) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{s.Name})
	for i := 0; i < rows; i++ {
		if s.FloatMat[i] < 1.0 {
			table.Append([]string{fmt.Sprintf("%.4f", s.FloatMat[i])})
		} else {
			table.Append([]string{fmt.Sprintf("%.2f", s.FloatMat[i])})
		}
	}
	table.Render()
}

func (s *Series) floatSeriesEqual(other *Series) *Series {
	data := make([]bool, s.Len())
	for i := 0; i < s.Len(); i++ {
		data[i] = s.FloatMat[i] == other.FloatMat[i]
	}
	return newBoolSeries(data, s.Name)
}

func (s *Series) floatSeriesNotEqual(other *Series) *Series {
	data := make([]bool, s.Len())
	for i := 0; i < s.Len(); i++ {
		data[i] = s.FloatMat[i] != other.FloatMat[i]
	}
	return newBoolSeries(data, s.Name)
}

func (s *Series) UniqueFloat() []float64 {
	u := []float64{}
	m := map[float64]bool{}
	for _, f := range s.FloatMat {
		if ok := m[f]; !ok {
			u = append(u, f)
			m[f] = true
		}
	}
	return u
}
