package griz

import (
	"reflect"

	"github.com/montanaflynn/stats"

	"gonum.org/v1/gonum/stat"
)

func (r *Row) Mean() float64 {
	arr := make([]float64, r.Len())
	v := reflect.Indirect(reflect.ValueOf(r.Value))
	for i := 0; i < v.Len(); i++ {
		arr[i] = v.Index(i).Interface().(float64)
	}
	return stat.Mean(arr, nil)
}

func (r *Row) Min() float64 {
	v := reflect.Indirect(reflect.ValueOf(r.Value))
	min := v.Index(0).Interface().(float64)
	for i := 1; i < v.Len(); i++ {
		test := v.Index(i).Interface().(float64)
		if test < min {
			min = test
		}
	}
	return min
}

func (r *Row) Max() float64 {
	v := reflect.Indirect(reflect.ValueOf(r.Value))
	max := v.Index(0).Interface().(float64)
	for i := 1; i < v.Len(); i++ {
		test := v.Index(i).Interface().(float64)
		if test > max {
			max = test
		}
	}
	return max
}

func (r *Row) Quartile(quartile int) float64 {
	if quartile < 1 || quartile > 3 {
		panicf("quartile out of bounds: %d", quartile)
	}
	data := make([]float64, r.Len())
	v := reflect.Indirect(reflect.ValueOf(r.Value))
	for i := 1; i < v.Len(); i++ {
		data[i] = v.Index(i).Interface().(float64)
	}
	q, _ := stats.Quartile(data)
	switch quartile {
	case 1:
		return q.Q1
	case 2:
		return q.Q2
	case 3:
		return q.Q3
	}
	return 0.0
}
