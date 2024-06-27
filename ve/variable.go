package ve

import (
	"fmt"
	"slices"
)

type Variable struct {
	id       uint16
	outcomes uint16
}

type Variables struct {
	variables []Variable
}

func NewVariables() *Variables {
	return &Variables{}
}

func (v *Variables) Add(outcomes uint16) Variable {
	v.variables = append(v.variables,
		Variable{
			id:       uint16(len(v.variables)),
			outcomes: outcomes,
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

	return Factor{
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

func (v *Variables) Product(f1 *Factor, f2 *Factor) Factor {
	newVars := []Variable{}
	map1 := make([]int, len(f1.variables))
	map2 := make([]int, len(f2.variables))

	for i, v := range f1.variables {
		idx := slices.Index(newVars, v)
		if idx < 0 {
			idx = len(newVars)
			newVars = append(newVars, v)
		}
		map1[i] = idx
	}
	for i, v := range f2.variables {
		idx := slices.Index(newVars, v)
		if idx < 0 {
			idx = len(newVars)
			newVars = append(newVars, v)
		}
		map2[i] = idx
	}

	f := v.CreateFactor(newVars, nil)

	oldIndex1 := make([]int, len(f1.variables))
	oldIndex2 := make([]int, len(f2.variables))
	newIndex := make([]int, len(f.variables))
	for i := range f.data {
		f.Outcomes(i, newIndex)
		for j, idx := range map1 {
			oldIndex1[j] = newIndex[idx]
		}
		for j, idx := range map2 {
			oldIndex2[j] = newIndex[idx]
		}
		f.data[i] = f1.Get(oldIndex1) * f2.Get(oldIndex2)
	}

	return f
}
