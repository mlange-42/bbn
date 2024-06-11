package main

import (
	"math/rand"
	"os"
	"path"
	"testing"

	"github.com/mlange-42/bbn"
	"github.com/stretchr/testify/assert"
)

func TestRunTrainCommand(t *testing.T) {
	net, err := runTrainCommand("../../_examples/fruits.yml", "../../_examples/fruits.csv", "", ',')
	assert.Nil(t, err)

	evidence := map[string]string{
		"Fruit": "banana",
		"Size":  "small",
	}
	result, err := net.Sample(evidence, 100_000, rand.New(rand.NewSource(1)))
	assert.Nil(t, err)

	assert.Equal(t, []float64{0.0, 1.0}, result["Tasty"])
}

func TestTrainUtility(t *testing.T) {
	dir := t.TempDir()

	a := bbn.Node{
		Variable: "A",
		Outcomes: []string{"yes", "no"},
		Table:    [][]float64{{0.0, 0.0}},
	}

	b := bbn.Node{
		Variable: "B",
		Outcomes: []string{"yes", "no"},
		Table:    [][]float64{{0.0, 0.0}},
	}

	utility := bbn.Node{
		Variable: "U",
		Type:     "utility",
		Given:    []string{"A", "B"},
		Outcomes: []string{"U"},
		Table: [][]float64{
			{0.0},
			{0.0},
			{0.0},
			{0.0},
		},
	}

	net, err := bbn.New("Utility", &a, &b, &utility)
	assert.Nil(t, err)

	yml, err := bbn.ToYAML(net)
	assert.Nil(t, err)

	err = os.WriteFile(path.Join(dir, "net.yml"), yml, 0644)
	assert.Nil(t, err)

	data := `A,B,U
yes,yes,0
yes,yes,0
yes,no,30
yes,no,50
no,yes,50
no,yes,70
no,no,100
no,no,100`

	err = os.WriteFile(path.Join(dir, "data.csv"), []byte(data), 0644)
	assert.Nil(t, err)

	net, err = runTrainCommand(path.Join(dir, "net.yml"), path.Join(dir, "data.csv"), "", ',')
	assert.Nil(t, err)

	evidence := map[string]string{
		"A": "no",
		"B": "yes",
	}
	result, err := net.Sample(evidence, 100_000, rand.New(rand.NewSource(1)))
	assert.Nil(t, err)

	assert.Equal(t, []float64{60}, result["U"])
}
