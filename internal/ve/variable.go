package ve

import (
	"fmt"
	"math"
	"slices"
)

type NodeType uint8

const (
	ChanceNode NodeType = iota
	DecisionNode
	UtilityNode
)

type Variable struct {
	Id       uint16
	outcomes uint16
	NodeType NodeType
}

// Variables provides variable and factor creation functionality, as well as factor operations.
type Variables struct {
	factorCounter int
	variables     variables
}

// NewVariables creates a new Variables instance.
func NewVariables() *Variables {
	return &Variables{}
}

// AddVariable creates and add a new [Variable].
func (v *Variables) AddVariable(nodeType NodeType, outcomes uint16) Variable {
	v.variables = append(v.variables,
		Variable{
			Id:       uint16(len(v.variables)),
			outcomes: outcomes,
			NodeType: nodeType,
		})
	return v.variables[len(v.variables)-1]
}

// CreateFactor creates a new [Factor] for the given variables.
//
// Argument data may be nil.
func (v *Variables) CreateFactor(vars []Variable, data []float64) Factor {
	rows := 1
	variables := make([]Variable, len(vars))
	for i, v := range vars {
		variables[i] = v
		rows *= int(v.outcomes)
	}
	if data == nil {
		data = make([]float64, rows)
	} else {
		if len(data) != rows {
			panic(fmt.Sprintf("wrong data length for factor. expected %d, got %d", rows, len(data)))
		}
	}

	v.factorCounter++

	return Factor{
		id:        v.factorCounter - 1,
		Variables: variables,
		Data:      data,
	}
}

// Restrict a factor to the given evidence.
func (v *Variables) Restrict(f *Factor, variable Variable, observation int) Factor {
	if observation < 0 || observation >= int(variable.outcomes) {
		panic(fmt.Sprintf("observation %d out of range for variable with %d possible observation values", observation, variable.outcomes))
	}
	newVars := []Variable{}
	idx := -1

	rows := 1
	for i := range f.Variables {
		if f.Variables[i].Id == variable.Id {
			idx = i
		} else {
			newVars = append(newVars, f.Variables[i])
			rows *= int(f.Variables[i].outcomes)
		}
	}

	if idx < 0 {
		panic(fmt.Sprintf("variable %d not in this factor", variable.Id))
	}

	fNew := v.CreateFactor(newVars, nil)

	oldIndex := make([]int, len(f.Variables))
	newIndex := make([]int, len(newVars))
	for i, v := range f.Data {
		f.Outcomes(i, oldIndex)
		if oldIndex[idx] != observation {
			continue
		}
		for j := 0; j < idx; j++ {
			newIndex[j] = oldIndex[j]
		}
		for j := idx + 1; j < len(oldIndex); j++ {
			newIndex[j-1] = oldIndex[j]
		}
		fNew.Set(newIndex, v)
	}

	return fNew
}

// SumOut a [Variable] from a [Factor].
func (v *Variables) SumOut(f *Factor, variable Variable) Factor {
	newVars := []Variable{}
	idx := -1

	rows := 1
	for i := range f.Variables {
		if f.Variables[i].Id == variable.Id {
			idx = i
		} else {
			newVars = append(newVars, f.Variables[i])
			rows *= int(f.Variables[i].outcomes)
		}
	}

	if idx < 0 {
		panic(fmt.Sprintf("variable %d not in this factor", variable.Id))
	}

	fNew := v.CreateFactor(newVars, nil)

	oldIndex := make([]int, len(f.Variables))
	newIndex := make([]int, len(newVars))

	for i, v := range f.Data {
		f.Outcomes(i, oldIndex)
		for j := 0; j < idx; j++ {
			newIndex[j] = oldIndex[j]
		}
		for j := idx + 1; j < len(oldIndex); j++ {
			newIndex[j-1] = oldIndex[j]
		}
		idx := fNew.Index(newIndex)
		fNew.Data[idx] += v
	}

	return fNew
}

// Policy derives a policy from a [Factor].
func (v *Variables) Policy(f *Factor, variable Variable) Factor {
	newVars := []Variable{}
	idx := -1

	for i := range f.Variables {
		if f.Variables[i].Id == variable.Id {
			idx = i
		} else {
			newVars = append(newVars, f.Variables[i])
		}
	}

	if idx < 0 {
		panic(fmt.Sprintf("variable %d not in this factor", variable.Id))
	}
	newVars = append(newVars, f.Variables[idx])

	fNew := v.CreateFactor(newVars, nil)

	oldIndex := make([]int, len(f.Variables))
	newIndex := make([]int, len(newVars))
	idxNew := len(newVars) - 1

	cols := int(f.Variables[idx].outcomes)
	rows := len(f.Data) / cols

	rowData := make([]float64, cols)
	for row := 0; row < rows; row++ {
		fNew.Outcomes(row*cols, newIndex)
		for c := 0; c < cols; c++ {
			newIndex[idxNew] = c

			for j := 0; j < idx; j++ {
				oldIndex[j] = newIndex[j]
			}
			for j := idx + 1; j < len(oldIndex); j++ {
				oldIndex[j] = newIndex[j-1]
			}
			oldIndex[idx] = newIndex[idxNew]
			rowData[c] = f.Get(oldIndex)
		}
		maxUtility := math.Inf(-1)
		maxIdx := -1
		for c, u := range rowData {
			if u > maxUtility {
				maxUtility = u
				maxIdx = c
			}
		}

		if maxIdx < 0 {
			panic("no utility values to derive policy")
		}

		for c := 0; c < cols; c++ {
			newIndex[idxNew] = c
			if c == maxIdx {
				fNew.Set(newIndex, 1)
			} else {
				fNew.Set(newIndex, 0)
			}
		}
	}

	return fNew
}

// Rearrange changes the [Variable] order of a [Factor].
func (v *Variables) Rearrange(f *Factor, variables []Variable) Factor {
	if len(f.Variables) != len(variables) {
		panic("number of old and new variables doesn't match")
	}
	varsEqual := true
	for i, vv := range variables {
		if vv != f.Variables[i] {
			varsEqual = false
			break
		}
	}
	if varsEqual {
		return v.CreateFactor(f.Variables, append([]float64{}, f.Data...))
	}

	indices := make([]int, len(variables))
	for i, vv := range variables {
		idx := slices.Index(f.Variables, vv)
		if idx < 0 {
			panic(fmt.Sprintf("variable %d not in original factor", vv.Id))
		}
		indices[i] = idx
	}

	fNew := v.CreateFactor(variables, nil)
	newIndex := make([]int, len(variables))
	oldIndex := make([]int, len(f.Variables))

	for i := range fNew.Data {
		fNew.Outcomes(i, newIndex)
		for j, idx := range newIndex {
			oldIndex[indices[j]] = idx
		}
		fNew.Data[i] = f.Get(oldIndex)
	}

	return fNew
}

// Product multiplies factors.
func (v *Variables) Product(factors ...*Factor) Factor {
	if len(factors) == 1 {
		return v.CreateFactor(factors[0].Variables, append([]float64{}, factors[0].Data...))
	}

	newVars := []Variable{}
	maps := make([][]int, len(factors))
	for i, f := range factors {
		m := make([]int, len(f.Variables))
		for j, v := range f.Variables {
			idx := slices.Index(newVars, v)
			if idx < 0 {
				idx = len(newVars)
				newVars = append(newVars, v)
			}
			m[j] = idx
		}
		maps[i] = m
	}

	f := v.CreateFactor(newVars, nil)

	oldIndex := make([][]int, len(factors))
	for i, f := range factors {
		oldIndex[i] = make([]int, len(f.Variables))
	}

	newIndex := make([]int, len(f.Variables))
	for i := range f.Data {
		f.Outcomes(i, newIndex)

		product := 1.0
		for j, fOld := range factors {
			m := maps[j]
			oldIdx := oldIndex[j]
			for k, idx := range m {
				oldIdx[k] = newIndex[idx]
			}
			product *= fOld.Get(oldIdx)
		}
		f.Data[i] = product
	}

	return f
}

// Sum sums up factors, similar to [Variables.Product].
func (v *Variables) Sum(factors ...*Factor) Factor {
	if len(factors) == 1 {
		return v.CreateFactor(factors[0].Variables, append([]float64{}, factors[0].Data...))
	}

	newVars := []Variable{}
	maps := make([][]int, len(factors))
	for i, f := range factors {
		m := make([]int, len(f.Variables))
		for j, v := range f.Variables {
			idx := slices.Index(newVars, v)
			if idx < 0 {
				idx = len(newVars)
				newVars = append(newVars, v)
			}
			m[j] = idx
		}
		maps[i] = m
	}

	f := v.CreateFactor(newVars, nil)

	oldIndex := make([][]int, len(factors))
	for i, f := range factors {
		oldIndex[i] = make([]int, len(f.Variables))
	}

	newIndex := make([]int, len(f.Variables))
	for i := range f.Data {
		f.Outcomes(i, newIndex)

		sum := 0.0
		for j, fOld := range factors {
			m := maps[j]
			oldIdx := oldIndex[j]
			for k, idx := range m {
				oldIdx[k] = newIndex[idx]
			}
			sum += fOld.Get(oldIdx)
		}
		f.Data[i] = sum
	}

	return f
}

// Marginal calculates marginal probabilities from a [Factor] for the given [Variable].
func (v *Variables) Marginal(f *Factor, variable Variable) Factor {
	idx := -1
	for i := range f.Variables {
		if f.Variables[i].Id == variable.Id {
			idx = i
			break
		}
	}

	if idx < 0 {
		panic(fmt.Sprintf("variable %d not in this factor", variable.Id))
	}

	newVars := []Variable{variable}
	fNew := v.CreateFactor(newVars, nil)

	oldIndex := make([]int, len(f.Variables))

	for i, v := range f.Data {
		f.Outcomes(i, oldIndex)
		fNew.Data[oldIndex[idx]] += v
	}
	return fNew
}

// Normalize normalizes a [Factor].
func (v *Variables) Normalize(f *Factor) Factor {
	fNew := v.CreateFactor(f.Variables, append([]float64{}, f.Data...))

	sum := 0.0
	for _, v := range fNew.Data {
		sum += v
	}
	if sum == 0 {
		return fNew
	}

	for i := range fNew.Data {
		fNew.Data[i] /= sum
	}

	return fNew
}

// NormalizeFor normalizes a [Factor] for a certain [Variable].
// It also re-arranges the new factor to have the normalized variable as the last one.
func (v *Variables) NormalizeFor(f *Factor, variable Variable) Factor {
	idx := -1
	for i := range f.Variables {
		if f.Variables[i].Id == variable.Id {
			idx = i
			break
		}
	}

	if idx < 0 {
		panic(fmt.Sprintf("variable %d not in this factor", variable.Id))
	}

	newVars := make([]Variable, len(f.Variables))
	for i := 0; i < idx; i++ {
		newVars[i] = f.Variables[i]
	}
	for i := idx + 1; i < len(f.Variables); i++ {
		newVars[i-1] = f.Variables[i]
	}
	newVars[len(newVars)-1] = f.Variables[idx]

	fNew := v.Rearrange(f, newVars)
	values := int(variable.outcomes)
	bins := len(fNew.Data) / values

	for i := 0; i < bins; i++ {
		sum := 0.0
		for j := 0; j < values; j++ {
			sum += fNew.Data[i*values+j]
		}
		if sum == 0 {
			continue
		}
		for j := 0; j < values; j++ {
			fNew.Data[i*values+j] /= sum
		}
	}

	return fNew
}

// Invert a [Factor] by applying 1/x for each element (if x != 0).
func (v *Variables) Invert(f *Factor) Factor {
	fNew := v.CreateFactor(f.Variables, append([]float64{}, f.Data...))

	for i, v := range fNew.Data {
		if v != 0 {
			fNew.Data[i] = 1.0 / v
		}
	}

	return fNew
}

// Helper type for a list of variables
type variables []Variable

// Index creates a flat [Factor] index from a multi-dimensional index.
func (v variables) Index(indices []int) int {
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

// Index creates multi-dimensional index from a flat [Factor] index.
func (v variables) Outcomes(index int, indices []int) {
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
