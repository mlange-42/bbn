package ve

import "fmt"

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

func (v *Variables) Add(outcomes uint16) *Variable {
	v.variables = append(v.variables,
		Variable{
			id:       uint16(len(v.variables)),
			outcomes: outcomes,
		})
	return &v.variables[len(v.variables)-1]
}

func (v *Variables) CreateFactor(vars []*Variable, data []float64) Factor {
	rows := 1
	variables := make([]Variable, len(vars))
	for i, v := range vars {
		variables[i] = *v
		rows *= int(v.outcomes)
	}
	if len(data) != rows {
		panic(fmt.Sprintf("wrong data length for factor. expected %d, got %d", rows, len(data)))
	}

	return Factor{
		variables: variables,
		data:      data,
	}
}
