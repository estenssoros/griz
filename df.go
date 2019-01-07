package griz

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"reflect"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

// DataFrame tries to provide functionality of pandas
type DataFrame struct {
	Columns   []string
	ColumnMap map[string]int
	DataTypes []int
	Mat       interface{} // is [][]interface{}
}

func (df DataFrame) String() string {
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	table := tablewriter.NewWriter(writer)
	table.SetHeader(df.Columns)
	for row := range df.Iterrows() {
		table.Append(row.ToString())
	}
	footer := make([]string, df.Width())
	for i := 0; i < df.Width(); i++ {
		footer[i] = DataTypeString(df.DataTypes[i])
	}
	table.SetFooter(footer)
	colors := make([]tablewriter.Colors, df.Width())
	for i := 0; i < df.Width(); i++ {
		colors[i] = tablewriter.Colors{tablewriter.FgGreenColor}
	}
	table.SetFooterColor(colors...)
	table.Render()
	writer.Flush()
	return b.String()
}

// Iterrows loops through the rows of a dataframe
func (df *DataFrame) Iterrows() chan *Row {
	ch := make(chan *Row)
	go func() {
		v := reflect.Indirect(reflect.ValueOf(df.Mat))
		for i := 0; i < v.Len(); i++ {
			row := df.Iloc(i)
			ch <- row
		}
		close(ch)
	}()
	return ch
}

// IterSeries loops through the series of a dataframe
func (df *DataFrame) IterSeries() chan *Series {
	//TODO
	ch := make(chan *Series)
	go func() {
		close(ch)
	}()
	return ch
}

// NewDataFrame takes an interface that is realld a [][]intrerface and
// a list of columns names. Returns a pointer to a DataFrame
func NewDataFrame(data interface{}, columns []string, dataTypes []int) *DataFrame {
	if len(columns) != len(dataTypes) {
		panic("new dataframe: columns do not match data types")
	}
	dataV := reflect.Indirect(reflect.ValueOf(data))
	if dataV.Len() == 0 {
		panic(errors.New("new dataframe: no data provided"))
	}
	rowLen := valueLen(dataV.Index(0))
	if rowLen != len(columns) {
		panic("new dataframe: row len does not match column")
	}
	for i := 1; i < dataV.Len(); i++ {
		if valueLen(dataV.Index(i)) != rowLen {
			panic(errors.Errorf("new dataframe: %d does not match inital dimention: %d", i, rowLen))
		}
	}
	columnMap := make(map[string]int)
	for i, c := range columns {
		if _, ok := columnMap[c]; ok {
			panicf("%s already exists in dataframe", c)
		}
		columnMap[c] = i
	}
	return &DataFrame{
		Columns:   columns,
		ColumnMap: columnMap,
		DataTypes: dataTypes,
		Mat:       data,
	}
}

// Len returns the length of the dataframe
func (df *DataFrame) Len() int {
	v := reflect.Indirect(reflect.ValueOf(df.Mat))
	return v.Len()
}

// Width return the width of the dataframe
func (df *DataFrame) Width() int {
	return len(df.Columns)
}

// Dim return the dimensions of the dataframe
func (df *DataFrame) Dim() (r int, c int) {
	return df.Len(), df.Width()
}

// Head prints the for x rows of a database
func (df *DataFrame) Head(rows int) {
	if rows >= df.Len() {
		panic("head index out of bounds")
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(df.Columns)
	for i := 0; i < rows; i++ {
		table.Append(df.Iloc(i).ToString())
	}
	table.Render()
}

func (df *DataFrame) getFloatSeries(column string) *Series {
	idx := df.ColumnMap[column]
	data := make([]float64, df.Len())
	for row := range df.Iterrows() {
		data[row.Idx] = row.Iloc(idx).(float64)
	}
	return newFloatSeries(data, column)
}

func (df *DataFrame) getStringSeries(column string) *Series {
	idx := df.ColumnMap[column]
	data := make([]string, df.Len())
	for row := range df.Iterrows() {
		data[row.Idx] = row.Iloc(idx).(string)
	}
	return newStringSeries(data, column)
}

func (df *DataFrame) getBoolSeries(column string) *Series {
	idx := df.ColumnMap[column]
	data := make([]bool, df.Len())
	for row := range df.Iterrows() {
		data[row.Idx] = row.Iloc(idx).(bool)
	}
	return newBoolSeries(data, column)
}

// GetSeries pulls a single column from a dataframe
func (df *DataFrame) GetSeries(column string) *Series {
	idx, ok := df.ColumnMap[column]
	if !ok {
		panic(errors.Errorf("column: %s does not exist", column))
	}
	switch df.DataTypes[idx] {
	case FloatType:
		return df.getFloatSeries(column)
	case StringType:
		return df.getStringSeries(column)
	case BoolType:
		return df.getBoolSeries(column)
	default:
		panicf("dataframe get: data type %d not supported", df.DataTypes[idx])
	}
	return nil
}

// GetDataFrame pulls a single column from a dataframe
func (df *DataFrame) GetDataFrame(columns ...string) *DataFrame {
	if len(columns) == 0 {
		panic("no columns supplied")
	}
	dataTypes := []int{}
	for _, c := range columns {
		if dIdx, ok := df.ColumnMap[c]; !ok {
			panic(fmt.Sprintf("column: %s does not exists", c))
		} else {
			dataTypes = append(dataTypes, df.DataTypes[dIdx])
		}
	}
	data := make([][]interface{}, df.Len())
	for row := range df.Iterrows() {
		newRow := make([]interface{}, len(columns))
		for i, c := range columns {
			newRow[i] = row.Loc(c)
		}
		data[row.Idx] = newRow
	}
	return NewDataFrame(data, columns, dataTypes)
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
	data := make([][]interface{}, df.Len())
	for row := range df.Iterrows() {
		newRow := make([]interface{}, len(columns))
		for i, c := range columns {
			newRow[i] = row.Loc(c)
		}
		data[row.Idx] = newRow
	}
	// TODO
	// return NewDataFrame(data, columns)
	return &DataFrame{}
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
	for i := 0; i < df.Len(); i++ {
		df.Iloc(i).Append(s.FloatMat[i])
	}
}

// SetLoc sets a location value
func (df *DataFrame) SetLoc(idx int, column string, val interface{}) {
	if _, ok := df.ColumnMap[column]; !ok {
		panic(fmt.Sprintf("column: %s does not exists", column))
	}
	df.Iloc(idx).SetLoc(column, val)
}

// Iloc returns a row by index
func (df *DataFrame) Iloc(idx int) *Row {
	if idx >= df.Len() {
		panic("df iloc: row index out of bounds")
	}
	dataV := reflect.Indirect(reflect.ValueOf(df.Mat))
	return &Row{
		Idx:       idx,
		Columns:   df.Columns,
		Value:     dataV.Index(idx).Interface(),
		DataTypes: df.DataTypes,
	}
}

// NewSeriesFromValue returns a series from a value that is the same length
// as the dataframe
func (df *DataFrame) NewSeriesFromValue(value float64, name string) *Series {
	mat := make([]float64, df.Len())
	for i := 0; i < df.Len(); i++ {
		mat[i] = value
	}
	return NewSeries(mat, name)
}

// Where indexs into a dataframe
func (df *DataFrame) Where(s *Series) *DataFrame {
	if s.DataType != BoolType {
		panic("dataframe where only supports boolean series")
	}
	data := [][]interface{}{}
	matV := reflect.Indirect(reflect.ValueOf(df.Mat))
	for i := 0; i < df.Len(); i++ {
		if s.Bool(i) {
			data = append(data, matV.Index(i).Interface().([]interface{}))
		}
	}
	return NewDataFrame(data, df.Columns, df.DataTypes)
}
