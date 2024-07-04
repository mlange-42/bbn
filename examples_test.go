package bbn_test

import (
	"fmt"

	"github.com/mlange-42/bbn"
)

func Example_sprinkler() {
	variables := []bbn.Variable{
		{Name: "Rain", Outcomes: []string{"yes", "no"}},
		{Name: "Sprinkler", Outcomes: []string{"yes", "no"}},
		{Name: "GrassWet", Outcomes: []string{"yes", "no"}},
	}

	factors := []bbn.Factor{
		{
			For: "Rain",
			Table: []float64{
				// rain+, rain-
				0.2, 0.8,
			},
		},
		{
			For:   "Sprinkler",
			Given: []string{"Rain"},
			Table: []float64{
				// yes   no
				0.01, 0.99, // rain+
				0.2, 0.8, // rain-
			},
		},
		{
			For:   "GrassWet",
			Given: []string{"Rain", "Sprinkler"},
			Table: []float64{
				//   yes   no
				0.99, 0.01, // rain+, sprikler+s
				0.8, 0.2, // rain+, sprikler-
				0.9, 0.1, // rain-, sprikler+
				0.0, 1.0, // rain-, sprikler-
			},
		},
	}

	net, err := bbn.New("Sprinkler", "", variables, factors)
	if err != nil {
		panic(err)
	}

	evidence := map[string]string{
		"GrassWet": "yes",
		"Rain":     "no",
	}
	query := []string{
		"Sprinkler",
	}

	result, _, err := net.SolveQuery(evidence, query, false)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
	// Output:
	//map[Sprinkler:[1 0]]
}
