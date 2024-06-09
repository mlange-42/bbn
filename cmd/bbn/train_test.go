package main

import (
	"math/rand"
	"testing"

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
