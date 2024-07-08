package ve

import (
	"fmt"
	"math"
	"slices"
)

// Variables provides variable and factor creation functionality, as well as factor operations.
type Variables struct {
	factorCounter int
	variables     factorVariables
	ids           map[int]bool
}

// NewVariables creates a new Variables instance.
func NewVariables() *Variables {
	return &Variables{
		ids: map[int]bool{},
	}
}

// AddVariable creates and add a new [Variable].
func (v *Variables) AddVariable(id int, nodeType NodeType, outcomes uint16) Variable {
	if _, ok := v.ids[id]; ok {
		panic(fmt.Sprintf("there is already a variable with ID %d", id))
	}
	v.variables = append(v.variables,
		Variable{
			id:       id,
			index:    uint16(len(v.variables)),
			outcomes: outcomes,
			nodeType: nodeType,
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
		variables: variables,
		data:      data,
	}
}

// Restrict a factor to the given evidence.
func (v *Variables) Restrict(f *Factor, variable Variable, observation int) Factor {
	if observation < 0 || observation >= int(variable.outcomes) {
		panic(fmt.Sprintf("observation %d out of range for variable with %d possible observation values", observation, variable.outcomes))
	}
	newVars := make([]Variable, 0, len(f.variables)-1)
	idx := -1

	rows := 1
	for i := range f.variables {
		if f.variables[i].id == variable.id {
			idx = i
		} else {
			newVars = append(newVars, f.variables[i])
			rows *= int(f.variables[i].outcomes)
		}
	}

	if idx < 0 {
		panic(fmt.Sprintf("variable %d not in this factor", variable.id))
	}

	fNew := v.CreateFactor(newVars, nil)

	oldIndex := make([]int, len(f.variables))
	newIndex := make([]int, len(newVars))
	for i, v := range f.data {
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
	newVars := make([]Variable, 0, len(f.variables)-1)
	idx := -1

	rows := 1
	for i := range f.variables {
		if f.variables[i].id == variable.id {
			idx = i
		} else {
			newVars = append(newVars, f.variables[i])
			rows *= int(f.variables[i].outcomes)
		}
	}

	if idx < 0 {
		panic(fmt.Sprintf("variable %d not in this factor", variable.id))
	}

	fNew := v.CreateFactor(newVars, nil)

	oldIndex := make([]int, len(f.variables))
	newIndex := make([]int, len(newVars))

	for i, v := range f.data {
		f.Outcomes(i, oldIndex)
		for j := 0; j < idx; j++ {
			newIndex[j] = oldIndex[j]
		}
		for j := idx + 1; j < len(oldIndex); j++ {
			newIndex[j-1] = oldIndex[j]
		}
		idx := fNew.Index(newIndex)
		fNew.data[idx] += v
	}

	return fNew
}

// Policy derives a policy from a [Factor].
func (v *Variables) Policy(f *Factor, variable Variable) Factor {
	newVars := make([]Variable, 0, len(f.variables))
	idx := -1

	for i := range f.variables {
		if f.variables[i].id == variable.id {
			idx = i
		} else {
			newVars = append(newVars, f.variables[i])
		}
	}

	if idx < 0 {
		panic(fmt.Sprintf("variable %d not in this factor", variable.id))
	}
	newVars = append(newVars, f.variables[idx])

	fNew := v.CreateFactor(newVars, nil)

	oldIndex := make([]int, len(f.variables))
	newIndex := make([]int, len(newVars))
	idxNew := len(newVars) - 1

	cols := int(f.variables[idx].outcomes)
	rows := len(f.data) / cols

	rowData := make([]float64, cols)
	maxIndices := make([]int, 0, 8)
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
		for c, u := range rowData {
			if u > maxUtility {
				maxUtility = u
				maxIndices = maxIndices[:1]
				maxIndices[0] = c
			} else if u == maxUtility {
				maxIndices = append(maxIndices, c)
			}
		}

		if len(maxIndices) == 0 {
			panic("no utility values to derive policy")
		}

		probValue := 1.0 / float64(len(maxIndices))
		for _, idx := range maxIndices {
			newIndex[idxNew] = idx
			fNew.Set(newIndex, probValue)
		}

		maxIndices = maxIndices[:0]
	}

	return fNew
}

// Rearrange changes the [Variable] order of a [Factor].
func (v *Variables) Rearrange(f *Factor, variables []Variable) Factor {
	if len(f.variables) != len(variables) {
		panic("number of old and new variables doesn't match")
	}
	varsEqual := true
	for i, vv := range variables {
		if vv != f.variables[i] {
			varsEqual = false
			break
		}
	}
	if varsEqual {
		return v.CreateFactor(f.variables, append([]float64{}, f.data...))
	}

	indices := make([]int, len(variables))
	for i, vv := range variables {
		idx := slices.Index(f.variables, vv)
		if idx < 0 {
			panic(fmt.Sprintf("variable %d not in original factor", vv.id))
		}
		indices[i] = idx
	}

	fNew := v.CreateFactor(variables, nil)
	newIndex := make([]int, len(variables))
	oldIndex := make([]int, len(f.variables))

	for i := range fNew.data {
		fNew.Outcomes(i, newIndex)
		for j, idx := range newIndex {
			oldIndex[indices[j]] = idx
		}
		fNew.data[i] = f.Get(oldIndex)
	}

	return fNew
}

// Product multiplies factors.
func (v *Variables) Product(factors ...*Factor) Factor {
	if len(factors) == 1 {
		return v.CreateFactor(factors[0].variables, append([]float64{}, factors[0].data...))
	}

	newVars := []Variable{}
	maps := make([][]int, len(factors))
	for i, f := range factors {
		m := make([]int, len(f.variables))
		for j, v := range f.variables {
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
		oldIndex[i] = make([]int, len(f.variables))
	}

	newIndex := make([]int, len(f.variables))
	for i := range f.data {
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
		f.data[i] = product
	}

	return f
}

// Sum sums up factors, similar to [Variables.Product].
func (v *Variables) Sum(factors ...*Factor) Factor {
	if len(factors) == 1 {
		return v.CreateFactor(factors[0].variables, append([]float64{}, factors[0].data...))
	}

	newVars := []Variable{}
	maps := make([][]int, len(factors))
	for i, f := range factors {
		m := make([]int, len(f.variables))
		for j, v := range f.variables {
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
		oldIndex[i] = make([]int, len(f.variables))
	}

	newIndex := make([]int, len(f.variables))
	for i := range f.data {
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
		f.data[i] = sum
	}

	return f
}

// Marginal calculates marginal probabilities from a [Factor] for the given [Variable].
func (v *Variables) Marginal(f *Factor, variable Variable) Factor {
	idx := -1
	for i := range f.variables {
		if f.variables[i].id == variable.id {
			idx = i
			break
		}
	}

	if idx < 0 {
		panic(fmt.Sprintf("variable %d not in this factor", variable.id))
	}

	newVars := []Variable{variable}
	fNew := v.CreateFactor(newVars, nil)

	oldIndex := make([]int, len(f.variables))

	for i, v := range f.data {
		f.Outcomes(i, oldIndex)
		fNew.data[oldIndex[idx]] += v
	}
	return fNew
}

// Normalize normalizes a [Factor].
func (v *Variables) Normalize(f *Factor) Factor {
	fNew := v.CreateFactor(f.variables, append([]float64{}, f.data...))

	sum := 0.0
	for _, v := range fNew.data {
		sum += v
	}
	//if sum == 0 {
	//	return fNew
	//}

	for i := range fNew.data {
		fNew.data[i] /= sum
	}

	return fNew
}

// NormalizeFor normalizes a [Factor] for a certain [Variable].
// It also re-arranges the new factor to have the normalized variable as the last one.
func (v *Variables) NormalizeFor(f *Factor, variable Variable) Factor {
	idx := -1
	for i := range f.variables {
		if f.variables[i].id == variable.id {
			idx = i
			break
		}
	}

	if idx < 0 {
		panic(fmt.Sprintf("variable %d not in this factor", variable.id))
	}

	newVars := make([]Variable, len(f.variables))
	for i := 0; i < idx; i++ {
		newVars[i] = f.variables[i]
	}
	for i := idx + 1; i < len(f.variables); i++ {
		newVars[i-1] = f.variables[i]
	}
	newVars[len(newVars)-1] = f.variables[idx]

	fNew := v.Rearrange(f, newVars)
	values := int(variable.outcomes)
	bins := len(fNew.data) / values

	for i := 0; i < bins; i++ {
		sum := 0.0
		for j := 0; j < values; j++ {
			sum += fNew.data[i*values+j]
		}
		if sum == 0 {
			continue
		}
		for j := 0; j < values; j++ {
			fNew.data[i*values+j] /= sum
		}
	}

	return fNew
}

// Invert a [Factor] by applying 1/x for each element (if x != 0).
func (v *Variables) Invert(f *Factor) Factor {
	fNew := v.CreateFactor(f.variables, append([]float64{}, f.data...))

	for i, v := range fNew.data {
		//if v == 0 {
		//	continue
		//}
		fNew.data[i] = 1.0 / v
	}

	return fNew
}
