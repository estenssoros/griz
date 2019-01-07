package griz

import (
	"bufio"
	"bytes"
	"os"

	"github.com/olekukonko/tablewriter"
)

func newStringSeries(data interface{}, name string) *Series {
	mat, ok := data.([]string)
	if !ok {
		panic("new float series: error converting data to []string")
	}
	return &Series{
		Name:      name,
		StringMat: mat,
		DataType:  StringType,
	}
}

func (s Series) StringString() string {
	var b bytes.Buffer
	{
		writer := bufio.NewWriter(&b)
		table := tablewriter.NewWriter(writer)
		table.SetHeader([]string{s.Name})
		for _, row := range s.StringMat {
			table.Append([]string{row})
		}
		table.SetFooter([]string{DataTypeString(s.DataType)})
		table.SetFooterColor(tablewriter.Colors{tablewriter.FgGreenColor})
		table.Render()
		writer.Flush()
	}
	return b.String()
}

func (s *Series) StringHead(rows int) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{s.Name})
	for i := 0; i < rows; i++ {
		table.Append([]string{s.StringMat[i]})
	}
	table.Render()
}

func (s *Series) stringSeriesEqual(other *Series) *Series {
	data := make([]bool, s.Len())
	for i := 0; i < s.Len(); i++ {
		data[i] = s.StringMat[i] == other.StringMat[i]
	}
	return newBoolSeries(data, s.Name)
}

func (s *Series) stringSeriesNotEqual(other *Series) *Series {
	data := make([]bool, s.Len())
	for i := 0; i < s.Len(); i++ {
		data[i] = s.StringMat[i] != other.StringMat[i]
	}
	return newBoolSeries(data, s.Name)
}
