package bbn_test

import (
	"testing"

	"github.com/mlange-42/bbn"
	"github.com/mlange-42/bbn/internal/ve"
	"github.com/stretchr/testify/assert"
)

func TestTrainer(t *testing.T) {
	vars := []bbn.Variable{
		{
			Name:     "Rain",
			Outcomes: []string{"yes", "no"},
		},
		{
			Name:     "Sprinkler",
			Outcomes: []string{"yes", "no"},
		},
		{
			Name:     "GrassWet",
			Outcomes: []string{"yes", "no"},
		},
	}

	factors := []bbn.Factor{
		{
			For:   "Rain",
			Table: []float64{0.0, 0.0},
		},
		{
			For:   "Sprinkler",
			Given: []string{"Rain"},
			Table: []float64{
				0.0, 0.0, // rain yes
				0.0, 0.0, // rain no
			},
		},
		{
			For:   "GrassWet",
			Given: []string{"Rain", "Sprinkler"},
			Table: []float64{
				0.0, 0.0, // rain yes, sprikler yes
				0.0, 0.0, // rain yes, sprikler no
				0.0, 0.0, // rain no, sprikler yes
				0.0, 0.0, // rain no, sprikler no
			},
		},
	}

	net := bbn.New("Sprinkler", vars, factors)

	data := [][]int{
		{0, 0, 0},
		{0, 1, 0},
		{1, 0, 0},
		{1, 1, 1},
	}

	trainer := bbn.NewTrainer(net)

	for _, row := range data {
		trainer.AddSample(row, nil)
	}

	net, err := trainer.UpdateNetwork()
	assert.Nil(t, err)

	evidence := map[string]string{
		"Rain":     "no",
		"GrassWet": "yes",
	}
	result, _, err := net.SolveQuery(evidence, []string{"Sprinkler"}, false)
	assert.Nil(t, err)

	assert.Equal(t, map[string][]float64{
		"Sprinkler": {1, 0},
	}, result)
}

func TestTrainerDecision(t *testing.T) {

	vars := []bbn.Variable{
		{Name: "weather", Type: ve.ChanceNode, Outcomes: []string{"rainy", "sunny"}},
		{Name: "forecast", Type: ve.ChanceNode, Outcomes: []string{"sunny", "cloudy", "rainy"}},
		{Name: "umbrella", Type: ve.DecisionNode, Outcomes: []string{"yes", "no"}},
		{Name: "utility", Type: ve.UtilityNode, Outcomes: []string{"utility"}},
	}

	factors := []bbn.Factor{
		{For: "weather", Table: []float64{
			// rain+, rain-
			0.0, 0.0,
		}},
		{For: "forecast", Given: []string{"weather"}, Table: []float64{
			// sunny, cloudy, rainy
			0.0, 0.0, 0.0, // rain+
			0.0, 0.0, 0.0, // rain-
		}},
		{For: "umbrella", Given: []string{"forecast"}},

		{For: "utility", Given: []string{"weather", "umbrella"}, Table: []float64{
			0, // rain+, umbrella+
			0, // rain+, umbrella-
			0, // rain-, umbrella+
			0, // rain-, umbrella-
		}},
	}

	net := bbn.New("umbrella", vars, factors)
	trainer := bbn.NewTrainer(net)
	_ = trainer

	samples := [][]int{
		{0, 2, 0, 0},
		{0, 2, 1, 0},
		{1, 0, 0, 0},
		{1, 0, 1, 0},
	}
	utility := [][]float64{
		{0, 0, 0, 70},
		{0, 0, 0, 0},
		{0, 0, 0, 20},
		{0, 0, 0, 100},
	}

	for i := range samples {
		trainer.AddSample(samples[i], utility[i])
	}

	net, err := trainer.UpdateNetwork()
	assert.Nil(t, err)

	policy, err := net.SolvePolicies(false)
	assert.Nil(t, err)

	assert.Equal(t, []float64{0, 1, 0.5, 0.5, 1, 0}, policy["umbrella"].Table)
}
