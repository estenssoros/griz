package griz

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// ReadSQL returns a dataframe from a context and query
func ReadSQL(ctx context.Context, query string) (*DataFrame, error) {
	db := ctx.Value("db").(*sqlx.DB)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	columns := make([]string, len(cols))
	dataTypes := make([]int, len(cols))
	for i, c := range cols {
		var typeName int
		columns[i] = c.Name()
		switch c.DatabaseTypeName() {
		case "VARCHAR", "TEXT":
			typeName = StringType
		case "FLOAT", "INT", "BIGINT", "DOUBLE":
			typeName = FloatType
		case "TINYINT":
			typeName = BoolType
		case "DATE":
			panic("read sql: date not implemented")
		case "DATETIME":
			panic("read sql: datetime not implemented")
		}
		dataTypes[i] = typeName
	}

	data := [][]interface{}{}
	rawResult := make([][]byte, len(cols))
	dest := make([]interface{}, len(cols)) // A temporary interface{} slice
	for i := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}
	for rows.Next() {
		row := make([]interface{}, len(cols))
		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}
		for i, raw := range rawResult {
			switch dataTypes[i] {
			case FloatType:
				row[i] = stringToFloat(string(raw))
			case StringType:
				row[i] = string(raw)
			case BoolType:
				row[i] = stringToBool(string(raw))
			}
		}
		data = append(data, row)
	}
	return NewDataFrame(data, columns, dataTypes), nil
}
