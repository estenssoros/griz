package griz

import "fmt"

func rowToString(row []float64) []string {
	out := make([]string, len(row))
	for i := range row {
		if row[i] < 1 {
			out[i] = fmt.Sprintf("%.4f", row[i])
		} else {
			out[i] = fmt.Sprintf("%.2f", row[i])
		}
	}
	return out
}
