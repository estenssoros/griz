package griz

func (df *DataFrame) Mean() *Series {
	for _, t := range df.DataTypes {
		if t != FloatType {
			panic("average: only supports float type")
		}
	}
	data := []float64{}
	for r := range df.Iterrows() {
		data = append(data, r.Mean())
	}
	return NewSeries(data, "avg")
}

func (df *DataFrame) Min() *Series {
	for _, t := range df.DataTypes {
		if t != FloatType {
			panic("average: only supports float type")
		}
	}
	data := []float64{}
	for r := range df.Iterrows() {
		data = append(data, r.Min())
	}
	return NewSeries(data, "min")
}

func (df *DataFrame) Max() *Series {
	for _, t := range df.DataTypes {
		if t != FloatType {
			panic("average: only supports float type")
		}
	}
	data := []float64{}
	for r := range df.Iterrows() {
		data = append(data, r.Max())
	}
	return NewSeries(data, "max")
}
