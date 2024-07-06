package ve

import "fmt"

type Factor struct {
	id        int
	data      []float64
	variables factorVariables
}

func (f *Factor) Data() []float64 {
	return f.data
}

func (f *Factor) Variables() []Variable {
	return []Variable(f.variables)
}

func (f *Factor) Index(indices []int) int {
	return f.variables.Index(indices)
}

func (f *Factor) IndexWithNoData(indices []int) (int, bool) {
	return f.variables.IndexWithNoData(indices)
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

// Helper type for a list of variables for a factor
type factorVariables []Variable

// Index creates a flat [Factor] index from a multi-dimensional index.
func (v factorVariables) Index(indices []int) int {
	if len(indices) != len(v) {
		panic(fmt.Sprintf("factor with %d variables can't use %d indices", len(v), len(indices)))
	}
	if len(v) == 0 {
		return 0
	}

	curr := len(v) - 1
	idx := indices[curr]
	stride := 1

	curr--
	for curr >= 0 {
		stride *= int(v[curr+1].outcomes)
		idx += indices[curr] * stride
		curr--
	}
	return idx
}

func (v factorVariables) IndexWithNoData(indices []int) (int, bool) {
	if len(indices) != len(v) {
		panic(fmt.Sprintf("factor with %d variables can't use %d indices", len(v), len(indices)))
	}
	if len(v) == 0 {
		return 0, true
	}

	curr := len(v) - 1
	idx := indices[curr]
	if idx < 0 {
		return 0, false
	}

	stride := 1

	curr--
	for curr >= 0 {
		currIdx := indices[curr]
		if currIdx < 0 {
			return 0, false
		}
		stride *= int(v[curr+1].outcomes)
		idx += currIdx * stride
		curr--
	}
	return idx, true
}

// Index creates multi-dimensional index from a flat [Factor] index.
func (v factorVariables) Outcomes(index int, indices []int) {
	if len(indices) != len(v) {
		panic(fmt.Sprintf("factor with %d variables can't use %d indices", len(v), len(indices)))
	}
	if len(v) == 0 {
		return
	}

	curr := len(v) - 1

	n := int(v[curr].outcomes)
	indices[curr] = index % n
	index /= n
	curr--
	for curr >= 0 {
		n := int(v[curr].outcomes)
		indices[curr] = index % n
		index /= n
		curr--
	}
}
