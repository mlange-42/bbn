package bbn_test

import (
	"math/rand"
	"testing"

	"github.com/mlange-42/bbn"
	"github.com/stretchr/testify/assert"
)

func TestTrainer(t *testing.T) {

	rain := bbn.Node{
		Variable: "Rain",
		Outcomes: []string{"yes", "no"},
		Table:    [][]float64{{0.0, 0.0}},
	}

	sprinkler := bbn.Node{
		Variable: "Sprinkler",
		Given:    []string{"Rain"},
		Outcomes: []string{"yes", "no"},
		Table: [][]float64{
			{0.0, 0.0}, // rain yes
			{0.0, 0.0}, // rain no
		},
	}

	grassWet := bbn.Node{
		Variable: "GrassWet",
		Given:    []string{"Rain", "Sprinkler"},
		Outcomes: []string{"yes", "no"},
		Table: [][]float64{
			{0.0, 0.0}, // rain yes, sprikler yes
			{0.0, 0.0}, // rain yes, sprikler no
			{0.0, 0.0}, // rain no, sprikler yes
			{0.0, 0.0}, // rain no, sprikler no
		},
	}

	net, err := bbn.New("Sprinkler", &rain, &sprinkler, &grassWet)
	assert.Nil(t, err)

	data := [][]int{
		{0, 0, 0},
		{0, 1, 0},
		{1, 0, 0},
		{1, 1, 1},
	}

	trainer := bbn.NewTrainer(net)

	for _, row := range data {
		trainer.AddSample(row)
	}
	net, err = trainer.UpdateNetwork()
	assert.Nil(t, err)

	evidence := map[string]string{
		"Rain":     "no",
		"GrassWet": "yes",
	}
	result, err := net.Sample(evidence, 100_000, rand.New(rand.NewSource(1)))
	assert.Nil(t, err)

	assert.Equal(t, map[string][]float64{
		"Rain":      {0, 1},
		"Sprinkler": {1, 0},
		"GrassWet":  {1, 0},
	}, result)
}
