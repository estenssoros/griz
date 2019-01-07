package griz

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func newBoolSeries(data interface{}, name string) *Series {
	mat, ok := data.([]bool)
	if !ok {
		panic("new float series: error converting data to []bool")
	}
	return &Series{
		Name:     name,
		BoolMat:  mat,
		DataType: BoolType,
	}
}

func (s Series) BoolString() string {
	var b bytes.Buffer
	{
		writer := bufio.NewWriter(&b)
		table := tablewriter.NewWriter(writer)
		table.SetHeader([]string{s.Name})
		for _, row := range s.BoolMat {
			table.Append([]string{fmt.Sprint(row)})
		}
		table.SetFooter([]string{DataTypeString(s.DataType)})
		table.SetFooterColor(tablewriter.Colors{tablewriter.FgGreenColor})
		table.Render()
		writer.Flush()
	}
	return b.String()
}

func (s *Series) BoolHead(rows int) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{s.Name})
	for i := 0; i < rows; i++ {
		for _, row := range s.BoolMat {
			table.Append([]string{fmt.Sprint(row)})
		}
	}
	table.Render()
}

func (s *Series) boolSeriesEqual(other *Series) *Series {
	data := make([]bool, s.Len())
	for i := 0; i < s.Len(); i++ {
		data[i] = s.BoolMat[i] == other.BoolMat[i]
	}
	return newBoolSeries(data, s.Name)
}

func (s *Series) boolSeriesNotEqual(other *Series) *Series {
	data := make([]bool, s.Len())
	for i := 0; i < s.Len(); i++ {
		data[i] = s.BoolMat[i] != other.BoolMat[i]
	}
	return newBoolSeries(data, s.Name)
}
