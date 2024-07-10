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
	idx := 0
	stride := int(f.variables[len(f.variables)-1].outcomes)

	for i := len(f.variables) - 2; i >= 0; i-- {
		idx += indices[i] * stride
		stride *= int(f.variables[i].outcomes)
	}

	return idx, int(f.variables[len(f.variables)-1].outcomes)
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
	index := 0
	multiplier := 1
	for i := len(indices) - 1; i >= 0; i-- {
		index += indices[i] * multiplier
		multiplier *= int(v[i].outcomes)
	}

	return index
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

	index := 0
	multiplier := 1

	for i := len(indices) - 1; i >= 0; i-- {
		if indices[i] == -1 {
			return 0, false
		}
		index += indices[i] * multiplier
		multiplier *= int(v[i].outcomes)
	}

	return index, true
}

// Index creates multi-dimensional index from a flat [Factor] index.
func (v factorVariables) Outcomes(index int, indices []int) {
	if len(indices) != len(v) {
		panic(fmt.Sprintf("factor with %d variables can't use %d indices", len(v), len(indices)))
	}
	if len(v) == 0 {
		return
	}

	for i := len(indices) - 1; i >= 0; i-- {
		indices[i] = index % int(v[i].outcomes)
		index /= int(v[i].outcomes)
	}
}

// increment increments the multi-dimensional index by one.
// Returns false if the index overflows.
func (v factorVariables) increment(indices []int) bool {
	if len(indices) != len(v) {
		panic(fmt.Sprintf("factor with %d variables can't use %d indices", len(v), len(indices)))
	}

	for i := len(indices) - 1; i >= 0; i-- {
		indices[i]++
		if indices[i] < int(v[i].outcomes) {
			return true
		}
		indices[i] = 0
	}

	return false
}
