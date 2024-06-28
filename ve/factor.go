package ve

import "fmt"

type Factor struct {
	id        int
	data      []float64
	variables variables
}

func (f *Factor) Index(indices []int) int {
	return f.variables.Index(indices)
}

func (f *Factor) Outcomes(index int, indices []int) {
	f.variables.Outcomes(index, indices)
}

func (f *Factor) RowIndex(indices []int) (int, int) {
	if len(indices) != len(f.variables)-1 {
		panic(fmt.Sprintf("factor with %d variables can't use %d row indices", len(f.variables), len(indices)))
	}
	cols := int(f.variables[len(f.variables)-1].outcomes)
	idx := 0
	stride := 1
	curr := len(f.variables) - 2
	for curr >= 0 {
		stride *= int(f.variables[curr+1].outcomes)
		idx += indices[curr] * stride
		curr--
	}
	return idx, cols
}

func (f *Factor) Get(indices []int) float64 {
	idx := f.Index(indices)
	return f.data[idx]
}

func (f *Factor) Set(indices []int, value float64) {
	idx := f.Index(indices)
	f.data[idx] = value
}

func (f *Factor) GetRow(indices []int) []float64 {
	idx, ln := f.RowIndex(indices)
	return f.data[idx : idx+ln]
}

func (f *Factor) Normalize() {
	sum := 0.0
	for _, v := range f.data {
		sum += v
	}
	if sum == 0 {
		return
	}

	for i := range f.data {
		f.data[i] /= sum
	}
}