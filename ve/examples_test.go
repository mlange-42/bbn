package ve_test

import (
	"fmt"

	"github.com/mlange-42/bbn/ve"
)

func ExampleVariables() {
	vars := ve.NewVariables()

	a := vars.AddVariable(0, ve.ChanceNode, 2)
	b := vars.AddVariable(1, ve.ChanceNode, 2)
	c := vars.AddVariable(2, ve.ChanceNode, 2)

	f1 := vars.CreateFactor([]ve.Variable{a, c}, []float64{
		// T  F      a
		0.5, 0.5, // T
		1.0, 0.0, // F
	})

	f2 := vars.CreateFactor([]ve.Variable{b, c}, []float64{
		// T  F      b
		0.0, 1.0, // T
		1.0, 0.0, // F
	})

	product := vars.Product(&f1, &f2)

	idx := make([]int, 3)
	for i := 0; i < len(product.Data())/2; i++ {
		product.Outcomes(i*2, idx)
		fmt.Printf("%.1f %.1f  %v\n", product.Data()[i*2], product.Data()[i*2+1], idx[:2])
	}
	// Output:
	//0.0 0.5  [0 0]
	//0.5 0.0  [0 1]
	//0.0 1.0  [1 0]
	//0.0 0.0  [1 1]
}

func ExampleVE() {
	vars := ve.NewVariables()

	rain := vars.AddVariable(0, ve.ChanceNode, 2)
	sprinkler := vars.AddVariable(1, ve.ChanceNode, 2)
	grass := vars.AddVariable(2, ve.ChanceNode, 2)

	fRain := vars.CreateFactor([]ve.Variable{rain}, []float64{
		0.2, 0.8,
	})

	fSprinkler := vars.CreateFactor([]ve.Variable{rain, sprinkler}, []float64{
		0.01, 0.99, // rain+
		0.2, 0.8, // rain-
	})

	fGrass := vars.CreateFactor([]ve.Variable{rain, sprinkler, grass}, []float64{
		0.99, 0.01, // rain+ sprinkler+
		0.8, 0.2, // rain+ sprinkler-
		0.9, 0.1, // rain- sprinkler+
		0.0, 1.0, // rain- sprinkler-
	})

	evidence := []ve.Evidence{
		{Variable: sprinkler, Value: 1},
		{Variable: grass, Value: 0},
	}
	query := []ve.Variable{rain}

	ve := ve.New(vars, []ve.Factor{fRain, fSprinkler, fGrass}, nil, nil)
	result := ve.SolveQuery(evidence, query)

	normalized := vars.Normalize(result)
	fmt.Println(normalized.Data())
	// Output:
	// [1 0]
}
