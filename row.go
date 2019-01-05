package grizzly

type Row struct {
	Idx    int
	Names  []string
	Values []float64
}

func (r *Row) Map() map[string]float64 {
	out := map[string]float64{}
	for i, n := range r.Names {
		out[n] = r.Values[i]
	}
	return out
}
