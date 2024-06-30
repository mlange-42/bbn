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
	id       uint16
	outcomes uint16
	nodeType NodeType
}

type Variables struct {
	factorCounter int
	variables     variables
}

func NewVariables() *Variables {
	return &Variables{}
}

func (v *Variables) Add(nodeType NodeType, outcomes uint16) Variable {
	v.variables = append(v.variables,
		Variable{
			id:       uint16(len(v.variables)),
			outcomes: outcomes,
			nodeType: nodeType,
		})
	return v.variables[len(v.variables)-1]
}

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

func (v *Variables) Restrict(f *Factor, variable Variable, observation int) Factor {
	if observation < 0 || observation >= int(variable.outcomes) {
		panic(fmt.Sprintf("observation %d out of range for variable with %d possible observation values", observation, variable.outcomes))
	}
	newVars := []Variable{}
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

func (v *Variables) SumOut(f *Factor, variable Variable) Factor {
	newVars := []Variable{}
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

func (v *Variables) Policy(f *Factor, variable Variable) Factor {
	newVars := []Variable{}
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
		return *f
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

func (v *Variables) Product(factors ...*Factor) Factor {
	if len(factors) == 1 {
		return *factors[0]
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

func (v *Variables) Sum(factors ...*Factor) Factor {
	if len(factors) == 1 {
		return *factors[0]
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

func (v *Variables) Normalize(f *Factor) Factor {
	fNew := v.CreateFactor(f.variables, append([]float64{}, f.data...))

	sum := 0.0
	for _, v := range fNew.data {
		sum += v
	}
	if sum == 0 {
		return fNew
	}

	for i := range fNew.data {
		fNew.data[i] /= sum
	}

	return fNew
}

type variables []Variable

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
