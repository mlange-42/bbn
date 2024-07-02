package ve

import "fmt"

type Factor struct {
	id        int
	Data      []float64
	Variables variables
}

func (f *Factor) Index(indices []int) int {
	return f.Variables.Index(indices)
}

func (f *Factor) IndexWithNoData(indices []int) (int, bool) {
	return f.Variables.IndexWithNoData(indices)
}

func (f *Factor) Outcomes(index int, indices []int) {
	f.Variables.Outcomes(index, indices)
}

func (f *Factor) RowIndex(indices []int) (int, int) {
	if len(indices) != len(f.Variables)-1 {
		panic(fmt.Sprintf("factor with %d variables can't use %d row indices", len(f.Variables), len(indices)))
	}
	cols := int(f.Variables[len(f.Variables)-1].outcomes)
	idx := 0
	stride := 1
	curr := len(f.Variables) - 2
	for curr >= 0 {
		stride *= int(f.Variables[curr+1].outcomes)
		idx += indices[curr] * stride
		curr--
	}
	return idx, cols
}

func (f *Factor) Get(indices []int) float64 {
	idx := f.Index(indices)
	return f.Data[idx]
}

func (f *Factor) Set(indices []int, value float64) {
	idx := f.Index(indices)
	f.Data[idx] = value
}

func (f *Factor) GetRow(indices []int) []float64 {
	idx, ln := f.RowIndex(indices)
	return f.Data[idx : idx+ln]
}
