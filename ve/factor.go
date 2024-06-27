package ve

import "fmt"

type Factor struct {
	data      []float64
	variables []Variable
}

func (f *Factor) Index(indices []int) int {
	if len(indices) != len(f.variables) {
		panic(fmt.Sprintf("factor with %d variables can't use %d indices", len(f.variables), len(indices)))
	}
	curr := len(f.variables) - 1
	idx := indices[curr]
	stride := 1

	curr--
	for curr >= 0 {
		stride *= int(f.variables[curr+1].outcomes)
		idx += indices[curr] * stride
		curr--
	}
	return idx
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

func (f *Factor) GetRow(indices []int) []float64 {
	idx, ln := f.RowIndex(indices)
	return f.data[idx : idx+ln]
}
