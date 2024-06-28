package ve

import (
	"fmt"
	"testing"
)

func TestEliminate(t *testing.T) {
	vars := NewVariables()

	rain := vars.Add(2)
	sprinkler := vars.Add(2)
	grass := vars.Add(2)

	fRain := vars.CreateFactor([]Variable{rain}, []float64{
		0.2, 0.8,
	})

	fSprinkler := vars.CreateFactor([]Variable{rain, sprinkler}, []float64{
		0.01, 0.99, // rain+
		0.2, 0.8, // rain-
	})

	fGrass := vars.CreateFactor([]Variable{rain, sprinkler, grass}, []float64{
		0.99, 0.01, // rain+ sprinkler+
		0.8, 0.2, // rain+ sprinkler-
		0.9, 0.1, // rain- sprinkler+
		0.0, 1.0, // rain- sprinkler-
	})

	ve := New(vars, []Factor{fRain, fSprinkler, fGrass})

	for k, f := range ve.factors {
		fmt.Println(k)
		fmt.Println(f)
	}
	fmt.Println("----------------")

	_ = ve.Eliminate([]Evidence{{Variable: grass, Value: 0}, {Variable: sprinkler, Value: 1}}, []Variable{rain})

	for k, f := range ve.factors {
		fmt.Println(k)
		fmt.Println(f)
	}
}
