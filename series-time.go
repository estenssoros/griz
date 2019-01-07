package griz

import (
	"bufio"
	"bytes"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
)

var timefmt = "2006-01-02 15:04:05MST"

func newTimeSeries(data interface{}, name string) *Series {
	mat, ok := data.([]time.Time)
	if !ok {
		panic("new time series: error converting data to []time.Time")
	}
	return &Series{
		Name:     name,
		TimeMat:  mat,
		DataType: TimeType,
	}
}

// TimeString writes matrix to table
func (s Series) TimeString() string {
	var b bytes.Buffer
	{
		writer := bufio.NewWriter(&b)
		table := tablewriter.NewWriter(writer)
		table.SetHeader([]string{s.Name})
		for _, row := range s.TimeMat {
			table.Append([]string{row.Format(timefmt)})
		}
		table.SetFooter([]string{DataTypeString(s.DataType)})
		table.SetFooterColor(tablewriter.Colors{tablewriter.FgGreenColor})
		table.Render()
		writer.Flush()
	}
	return b.String()
}

// TimeHead prints a table the first x rows of a series
func (s *Series) TimeHead(rows int) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{s.Name})
	for _, row := range s.TimeMat {
		table.Append([]string{row.String()})
	}
	table.Render()
}

func (s *Series) timeSeriesEqual(other *Series) *Series {
	data := make([]bool, s.Len())
	for i := 0; i < s.Len(); i++ {
		data[i] = time.Time(s.TimeMat[i]).Equal(time.Time(other.TimeMat[i]))
	}
	return newBoolSeries(data, s.Name)
}

func (s *Series) timeSeriesNotEqual(other *Series) *Series {
	data := make([]bool, s.Len())
	for i := 0; i < s.Len(); i++ {
		data[i] = !time.Time(s.TimeMat[i]).Equal(time.Time(other.TimeMat[i]))
	}
	return newBoolSeries(data, s.Name)
}
