package ve

import "fmt"

// Factor definition.
// Create factors with [Variables.CreateFactor] of one of the operations on [Variables].
type Factor struct {
	id        int
	data      []float64
	variables factorVariables
}

// Id of the factor.
//
// Unique for factors created using the same [Variables] instance.
func (f *Factor) Id() int {
	return f.id
}

// Data of the factor.
func (f *Factor) Data() []float64 {
	return f.data
}

// Variables of the factor.
func (f *Factor) Variables() []Variable {
	return []Variable(f.variables)
}

// Index in [Factor.Data] from outcome indices of factor variables.
func (f *Factor) Index(indices []int) int {
	return f.variables.Index(indices)
}

// Index in [Factor.Data] from outcome indices of factor variables.
// Can be used with missing data, represented by -1.
// Second return value is false if there is missing data in indices.
func (f *Factor) IndexWithNoData(indices []int) (int, bool) {
	return f.variables.IndexWithNoData(indices)
}

// Outcomes calculates variable outcomes for a flat index of [Factor.Data].
// Inverse operation of [Factor.Index].
func (f *Factor) Outcomes(index int, indices []int) {
	f.variables.Outcomes(index, indices)
}

// RowIndex returns starting index and length of a "row" of data.
// A row is for fixed outcomes of all but the last variable,
// while the outcome of the last variable is used to index in the row.
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

// Get the factor's value for the given variable outcome indices.
func (f *Factor) Get(indices []int) float64 {
	idx := f.Index(indices)
	return f.data[idx]
}

// Set the factor's value for the given variable outcome indices.
func (f *Factor) Set(indices []int, value float64) {
	idx := f.Index(indices)
	f.data[idx] = value
}

// Row returns a "row" of data.
// A row is for fixed outcomes of all but the last variable,
// while the outcome of the last variable is used to index in the row.
func (f *Factor) Row(indices []int) []float64 {
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

// Index creates a flat [Factor] index from a multi-dimensional index.
// When data is missing in the index (represented by -1),
// the second return value is false.
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
