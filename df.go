package griz

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

type DataFrame struct {
	Columns   []string
	ColumnMap map[string]int
	Mat       [][]float64
}

func (df DataFrame) String() string {
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	table := tablewriter.NewWriter(writer)
	table.SetHeader(df.Columns)
	for _, row := range df.Mat {
		table.Append(rowToString(row))
	}
	table.Render()
	writer.Flush()
	return b.String()
}

// Len returns the length of the dataframe
func (df *DataFrame) Len() int {
	return len(df.Mat)
}

// Width return the width of the dataframe
func (df *DataFrame) Width() int {
	return len(df.Columns)
}

// Dim return the dimensions of the dataframe
func (df *DataFrame) Dim() (r int, c int) {
	return df.Len(), df.Width()
}

func (df *DataFrame) Head(rows int) {
	if rows >= df.Len() {
		panic("head index out of bounds")
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(df.Columns)
	for i := 0; i < rows; i++ {
		table.Append(rowToString(df.Row(i)))
	}
	table.Render()
}

// Get pulls a single column from a dataframe
func (df *DataFrame) Get(column string) *Series {
	idx, ok := df.ColumnMap[column]
	if !ok {
		panic(errors.Errorf("column: %s does not exist", column))
	}
	data := make([]float64, len(df.Mat))
	for i, row := range df.Mat {
		data[i] = row[idx]
	}
	return newSeries(data, column)
}

// GetMany returns a dataframe with the selected columns
func (df *DataFrame) GetMany(columns []string) *DataFrame {
	if len(columns) == 0 {
		panic("no columns supplied")
	}
	for _, c := range columns {
		if _, ok := df.ColumnMap[c]; !ok {
			panic(fmt.Sprintf("column: %s does not exists", c))
		}
	}
	data := make([][]float64, df.Len())
	for rowNum, row := range df.Mat {
		newRow := make([]float64, len(columns))
		for i, c := range columns {
			newRow[i] = row[df.ColumnMap[c]]
		}
		data[rowNum] = newRow
	}
	return NewDataFrame(data, columns)
}

// AddColumn adds a column to the dataframe
func (df *DataFrame) AddColumn(s *Series) {
	if s.Len() != df.Len() {
		panic(errors.Errorf("dimension mismatch: dataframe: %d series: %d", df.Len(), s.Len()))
	}
	if _, ok := df.ColumnMap[s.Name]; ok {
		panic(errors.Errorf("%s column already existing in dataframe", s.Name))
	}
	df.Columns = append(df.Columns, s.Name)
	df.ColumnMap[s.Name] = len(df.Columns) - 1
	for i := range df.Mat {
		df.Mat[i] = append(df.Mat[i], s.Mat[i])
	}
}

// Row returns a slice of row
func (df *DataFrame) Row(idx int) []float64 {
	if idx >= df.Len() {
		panic("row index out of bounds")
	}
	return df.Mat[idx]
}

func (df *DataFrame) GetLoc(column string, idx int) float64 {
	colIdx, ok := df.ColumnMap[column]
	if !ok {
		panic(fmt.Sprintf("column: %s does not exists", column))
	}
	return df.Row(idx)[colIdx]
}

func (df *DataFrame) SetLoc(column string, idx int, val float64) {
	colIdx, ok := df.ColumnMap[column]
	if !ok {
		panic(fmt.Sprintf("column: %s does not exists", column))
	}
	df.Row(idx)[colIdx] = val
}

// NewSeriesFromValue returns a series from a value that is the same length
// as the dataframe
func (df *DataFrame) NewSeriesFromValue(value float64, name string) *Series {
	mat := make([]float64, df.Len())
	for i := 0; i < df.Len(); i++ {
		mat[i] = value
	}
	return newSeries(mat, name)
}

func NewDataFrame(raw [][]float64, columns []string) *DataFrame {
	if len(raw) == 0 {
		panic(errors.New("no data present in matrix"))
	}
	// data := []float64{}
	// copy(data, raw[0])
	rowLen := len(raw[0])
	for i := 1; i < len(raw); i++ {
		if len(raw[i]) != rowLen {
			panic(errors.Errorf("row: %d does not match inital dimention: %d", i, rowLen))
		}
		// data = append(data, raw[i]...)
	}
	// mat := mat64.NewDense(len(columns), len(data)/len(columns), data)
	columnMap := make(map[string]int)
	for i, c := range columns {
		columnMap[c] = i
	}
	return &DataFrame{
		Columns:   columns,
		ColumnMap: columnMap,
		Mat:       raw,
	}
}

func (df *DataFrame) ToRow(idx int) *Row {
	return &Row{
		Idx:    idx,
		Names:  df.Columns,
		Values: df.Row(idx),
	}
}

func (df *DataFrame) Iterrows() chan *Row {
	ch := make(chan *Row)
	go func() {
		for i := 0; i < df.Len(); i++ {
			ch <- df.ToRow(i)
		}
		close(ch)
	}()
	return ch
}
